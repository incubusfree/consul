package xds

import (
	"errors"
	"fmt"

	envoy_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"

	"github.com/golang/protobuf/proto"
	bexpr "github.com/hashicorp/go-bexpr"

	"github.com/hashicorp/consul/agent/connect"
	"github.com/hashicorp/consul/agent/proxycfg"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/hashicorp/consul/api"
)

const (
	UnnamedSubset = ""
)

// endpointsFromSnapshot returns the xDS API representation of the "endpoints"
func (s *ResourceGenerator) endpointsFromSnapshot(cfgSnap *proxycfg.ConfigSnapshot) ([]proto.Message, error) {
	if cfgSnap == nil {
		return nil, errors.New("nil config given")
	}

	switch cfgSnap.Kind {
	case structs.ServiceKindConnectProxy:
		return s.endpointsFromSnapshotConnectProxy(cfgSnap)
	case structs.ServiceKindTerminatingGateway:
		return s.endpointsFromSnapshotTerminatingGateway(cfgSnap)
	case structs.ServiceKindMeshGateway:
		return s.endpointsFromSnapshotMeshGateway(cfgSnap)
	case structs.ServiceKindIngressGateway:
		return s.endpointsFromSnapshotIngressGateway(cfgSnap)
	default:
		return nil, fmt.Errorf("Invalid service kind: %v", cfgSnap.Kind)
	}
}

// endpointsFromSnapshotConnectProxy returns the xDS API representation of the "endpoints"
// (upstream instances) in the snapshot.
func (s *ResourceGenerator) endpointsFromSnapshotConnectProxy(cfgSnap *proxycfg.ConfigSnapshot) ([]proto.Message, error) {
	// TODO: this estimate is wrong
	resources := make([]proto.Message, 0,
		len(cfgSnap.ConnectProxy.PreparedQueryEndpoints)+
			cfgSnap.ConnectProxy.PeerUpstreamEndpoints.Len()+
			len(cfgSnap.ConnectProxy.WatchedUpstreamEndpoints))

	// NOTE: Any time we skip a chain below we MUST also skip that discovery chain in clusters.go
	// so that the sets of endpoints generated matches the sets of clusters.
	for uid, chain := range cfgSnap.ConnectProxy.DiscoveryChain {
		upstream := cfgSnap.ConnectProxy.UpstreamConfig[uid]

		explicit := upstream.HasLocalPortOrSocket()
		implicit := cfgSnap.ConnectProxy.IsImplicitUpstream(uid)
		if !implicit && !explicit {
			// Discovery chain is not associated with a known explicit or implicit upstream so it is skipped.
			continue
		}

		var upstreamConfigMap map[string]interface{}
		if upstream != nil {
			upstreamConfigMap = upstream.Config
		}

		es, err := s.endpointsFromDiscoveryChain(
			uid,
			chain,
			cfgSnap.Locality,
			upstreamConfigMap,
			cfgSnap.ConnectProxy.WatchedUpstreamEndpoints[uid],
			cfgSnap.ConnectProxy.WatchedGatewayEndpoints[uid],
			false,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, es...)
	}

	// NOTE: Any time we skip an upstream below we MUST also skip that same
	// upstream in clusters.go so that the sets of endpoints generated matches
	// the sets of clusters.
	for _, uid := range cfgSnap.ConnectProxy.PeeredUpstreamIDs() {
		upstreamCfg := cfgSnap.ConnectProxy.UpstreamConfig[uid]

		explicit := upstreamCfg.HasLocalPortOrSocket()
		implicit := cfgSnap.ConnectProxy.IsImplicitUpstream(uid)
		if !implicit && !explicit {
			// Not associated with a known explicit or implicit upstream so it is skipped.
			continue
		}

		peerMeta := cfgSnap.ConnectProxy.UpstreamPeerMeta(uid)

		// TODO(peering): if we replicated service metadata separately from the
		// instances we wouldn't have to flip/flop this cluster name like this.
		clusterName := peerMeta.PrimarySNI()
		if clusterName == "" {
			clusterName = uid.EnvoyID()
		}

		// Also skip peer instances with a hostname as their address. EDS
		// cannot resolve hostnames, so we provide them through CDS instead.
		if _, ok := cfgSnap.ConnectProxy.PeerUpstreamEndpointsUseHostnames[uid]; ok {
			continue
		}

		endpoints, ok := cfgSnap.ConnectProxy.PeerUpstreamEndpoints.Get(uid)
		if ok {
			la := makeLoadAssignment(
				clusterName,
				[]loadAssignmentEndpointGroup{
					{Endpoints: endpoints},
				},
				proxycfg.GatewayKey{ /*empty so it never matches*/ },
			)
			resources = append(resources, la)
		}
	}

	// Looping over explicit upstreams is only needed for prepared queries because they do not have discovery chains
	for _, u := range cfgSnap.Proxy.Upstreams {
		if u.DestinationType != structs.UpstreamDestTypePreparedQuery {
			continue
		}
		uid := proxycfg.NewUpstreamID(&u)

		dc := u.Datacenter
		if dc == "" {
			dc = cfgSnap.Datacenter
		}
		clusterName := connect.UpstreamSNI(&u, "", dc, cfgSnap.Roots.TrustDomain)

		endpoints, ok := cfgSnap.ConnectProxy.PreparedQueryEndpoints[uid]
		if ok {
			la := makeLoadAssignment(
				clusterName,
				[]loadAssignmentEndpointGroup{
					{Endpoints: endpoints},
				},
				cfgSnap.Locality,
			)
			resources = append(resources, la)
		}
	}

	// Loop over potential destinations in the mesh, then grab the gateway nodes associated with each
	cfgSnap.ConnectProxy.DestinationsUpstream.ForEachKey(func(uid proxycfg.UpstreamID) bool {
		name := clusterNameForDestination(cfgSnap, uid.Name, uid.NamespaceOrDefault(), uid.PartitionOrDefault())

		endpoints, ok := cfgSnap.ConnectProxy.DestinationGateways.Get(uid)
		if ok {
			la := makeLoadAssignment(
				name,
				[]loadAssignmentEndpointGroup{
					{Endpoints: endpoints},
				},
				proxycfg.GatewayKey{ /*empty so it never matches*/ },
			)
			resources = append(resources, la)
		}
		return true
	})

	return resources, nil
}

func (s *ResourceGenerator) filterSubsetEndpoints(subset *structs.ServiceResolverSubset, endpoints structs.CheckServiceNodes) (structs.CheckServiceNodes, error) {
	// locally execute the subsets filter
	if subset.Filter != "" {
		filter, err := bexpr.CreateFilter(subset.Filter, nil, endpoints)
		if err != nil {
			return nil, err
		}

		raw, err := filter.Execute(endpoints)
		if err != nil {
			return nil, err
		}
		return raw.(structs.CheckServiceNodes), nil
	}
	return endpoints, nil
}

func (s *ResourceGenerator) endpointsFromSnapshotTerminatingGateway(cfgSnap *proxycfg.ConfigSnapshot) ([]proto.Message, error) {
	return s.endpointsFromServicesAndResolvers(cfgSnap, cfgSnap.TerminatingGateway.ServiceGroups, cfgSnap.TerminatingGateway.ServiceResolvers)
}

func (s *ResourceGenerator) endpointsFromSnapshotMeshGateway(cfgSnap *proxycfg.ConfigSnapshot) ([]proto.Message, error) {
	keys := cfgSnap.MeshGateway.GatewayKeys()
	resources := make([]proto.Message, 0, len(keys)+len(cfgSnap.MeshGateway.ServiceGroups))

	for _, key := range keys {
		if key.Matches(cfgSnap.Datacenter, cfgSnap.ProxyID.PartitionOrDefault()) {
			continue // skip local
		}
		// Also skip gateways with a hostname as their address. EDS cannot resolve hostnames,
		// so we provide them through CDS instead.
		if len(cfgSnap.MeshGateway.HostnameDatacenters[key.String()]) > 0 {
			continue
		}

		endpoints := cfgSnap.GetMeshGatewayEndpoints(key)
		if len(endpoints) == 0 {
			s.Logger.Error("skipping mesh gateway endpoints because no definition found", "datacenter", key)
			continue
		}

		{ // standard connect
			clusterName := connect.GatewaySNI(key.Datacenter, key.Partition, cfgSnap.Roots.TrustDomain)

			la := makeLoadAssignment(
				clusterName,
				[]loadAssignmentEndpointGroup{
					{Endpoints: endpoints},
				},
				cfgSnap.Locality,
			)
			resources = append(resources, la)
		}

		if cfgSnap.ProxyID.InDefaultPartition() &&
			cfgSnap.ServiceMeta[structs.MetaWANFederationKey] == "1" &&
			cfgSnap.ServerSNIFn != nil {

			clusterName := cfgSnap.ServerSNIFn(key.Datacenter, "")
			la := makeLoadAssignment(
				clusterName,
				[]loadAssignmentEndpointGroup{
					{Endpoints: endpoints},
				},
				cfgSnap.Locality,
			)
			resources = append(resources, la)
		}
	}

	// generate endpoints for our servers if WAN federation is enabled
	if cfgSnap.ProxyID.InDefaultPartition() &&
		cfgSnap.ServiceMeta[structs.MetaWANFederationKey] == "1" &&
		cfgSnap.ServerSNIFn != nil {
		var allServersLbEndpoints []*envoy_endpoint_v3.LbEndpoint

		for _, srv := range cfgSnap.MeshGateway.ConsulServers {
			clusterName := cfgSnap.ServerSNIFn(cfgSnap.Datacenter, srv.Node.Node)

			_, addr, port := srv.BestAddress(false /*wan*/)

			lbEndpoint := &envoy_endpoint_v3.LbEndpoint{
				HostIdentifier: &envoy_endpoint_v3.LbEndpoint_Endpoint{
					Endpoint: &envoy_endpoint_v3.Endpoint{
						Address: makeAddress(addr, port),
					},
				},
				HealthStatus: envoy_core_v3.HealthStatus_UNKNOWN,
			}

			cla := &envoy_endpoint_v3.ClusterLoadAssignment{
				ClusterName: clusterName,
				Endpoints: []*envoy_endpoint_v3.LocalityLbEndpoints{{
					LbEndpoints: []*envoy_endpoint_v3.LbEndpoint{lbEndpoint},
				}},
			}
			allServersLbEndpoints = append(allServersLbEndpoints, lbEndpoint)

			resources = append(resources, cla)
		}

		// And add one catch all so that remote datacenters can dial ANY server
		// in this datacenter without knowing its name.
		resources = append(resources, &envoy_endpoint_v3.ClusterLoadAssignment{
			ClusterName: cfgSnap.ServerSNIFn(cfgSnap.Datacenter, ""),
			Endpoints: []*envoy_endpoint_v3.LocalityLbEndpoints{{
				LbEndpoints: allServersLbEndpoints,
			}},
		})
	}

	// Generate the endpoints for each service and its subsets
	e, err := s.endpointsFromServicesAndResolvers(cfgSnap, cfgSnap.MeshGateway.ServiceGroups, cfgSnap.MeshGateway.ServiceResolvers)
	if err != nil {
		return nil, err
	}
	resources = append(resources, e...)

	// Generate the endpoints for exported discovery chain targets.
	e, err = s.makeExportedUpstreamEndpointsForMeshGateway(cfgSnap)
	if err != nil {
		return nil, err
	}
	resources = append(resources, e...)

	return resources, nil
}

func (s *ResourceGenerator) endpointsFromServicesAndResolvers(
	cfgSnap *proxycfg.ConfigSnapshot,
	services map[structs.ServiceName]structs.CheckServiceNodes,
	resolvers map[structs.ServiceName]*structs.ServiceResolverConfigEntry,
) ([]proto.Message, error) {
	resources := make([]proto.Message, 0, len(services))

	// generate the endpoints for the linked service groups
	for svc, endpoints := range services {
		// Skip creating endpoints for services that have hostnames as addresses
		// EDS cannot resolve hostnames so we provide them through CDS instead
		if cfgSnap.Kind == structs.ServiceKindTerminatingGateway && len(cfgSnap.TerminatingGateway.HostnameServices[svc]) > 0 {
			continue
		}

		clusterEndpoints := make(map[string][]loadAssignmentEndpointGroup)
		clusterEndpoints[UnnamedSubset] = []loadAssignmentEndpointGroup{{Endpoints: endpoints, OnlyPassing: false}}

		// Collect all of the loadAssignmentEndpointGroups for the various subsets. We do this before generating
		// the endpoints for the default/unnamed subset so that we can take into account the DefaultSubset on the
		// service-resolver which may prevent the default/unnamed cluster from creating endpoints for all service
		// instances.
		if resolver, hasResolver := resolvers[svc]; hasResolver {
			for subsetName, subset := range resolver.Subsets {
				subsetEndpoints, err := s.filterSubsetEndpoints(&subset, endpoints)
				if err != nil {
					return nil, err
				}
				groups := []loadAssignmentEndpointGroup{{Endpoints: subsetEndpoints, OnlyPassing: subset.OnlyPassing}}
				clusterEndpoints[subsetName] = groups

				// if this subset is the default then override the unnamed subset with this configuration
				if subsetName == resolver.DefaultSubset {
					clusterEndpoints[UnnamedSubset] = groups
				}
			}
		}

		// now generate the load assignment for all subsets
		for subsetName, groups := range clusterEndpoints {
			clusterName := connect.ServiceSNI(svc.Name, subsetName, svc.NamespaceOrDefault(), svc.PartitionOrDefault(), cfgSnap.Datacenter, cfgSnap.Roots.TrustDomain)
			la := makeLoadAssignment(
				clusterName,
				groups,
				cfgSnap.Locality,
			)
			resources = append(resources, la)
		}
	}

	return resources, nil
}

func (s *ResourceGenerator) endpointsFromSnapshotIngressGateway(cfgSnap *proxycfg.ConfigSnapshot) ([]proto.Message, error) {
	var resources []proto.Message
	createdClusters := make(map[proxycfg.UpstreamID]bool)
	for _, upstreams := range cfgSnap.IngressGateway.Upstreams {
		for _, u := range upstreams {
			uid := proxycfg.NewUpstreamID(&u)

			// If we've already created endpoints for this upstream, skip it. Multiple listeners may
			// reference the same upstream, so we don't need to create duplicate endpoints in that case.
			if createdClusters[uid] {
				continue
			}

			es, err := s.endpointsFromDiscoveryChain(
				uid,
				cfgSnap.IngressGateway.DiscoveryChain[uid],
				proxycfg.GatewayKey{Datacenter: cfgSnap.Datacenter, Partition: u.DestinationPartition},
				u.Config,
				cfgSnap.IngressGateway.WatchedUpstreamEndpoints[uid],
				cfgSnap.IngressGateway.WatchedGatewayEndpoints[uid],
				false,
			)
			if err != nil {
				return nil, err
			}
			resources = append(resources, es...)
			createdClusters[uid] = true
		}
	}
	return resources, nil
}

// used in clusters.go
func makeEndpoint(host string, port int) *envoy_endpoint_v3.LbEndpoint {
	return &envoy_endpoint_v3.LbEndpoint{
		HostIdentifier: &envoy_endpoint_v3.LbEndpoint_Endpoint{
			Endpoint: &envoy_endpoint_v3.Endpoint{
				Address: makeAddress(host, port),
			},
		},
	}
}

func makePipeEndpoint(path string) *envoy_endpoint_v3.LbEndpoint {
	return &envoy_endpoint_v3.LbEndpoint{
		HostIdentifier: &envoy_endpoint_v3.LbEndpoint_Endpoint{
			Endpoint: &envoy_endpoint_v3.Endpoint{
				Address: makePipeAddress(path, 0),
			},
		},
	}
}

func (s *ResourceGenerator) endpointsFromDiscoveryChain(
	uid proxycfg.UpstreamID,
	chain *structs.CompiledDiscoveryChain,
	gatewayKey proxycfg.GatewayKey,
	upstreamConfigMap map[string]interface{},
	upstreamEndpoints map[string]structs.CheckServiceNodes,
	gatewayEndpoints map[string]structs.CheckServiceNodes,
	forMeshGateway bool,
) ([]proto.Message, error) {
	if chain == nil {
		if forMeshGateway {
			return nil, fmt.Errorf("missing discovery chain for %s", uid)
		}
		return nil, nil
	}

	if upstreamConfigMap == nil {
		upstreamConfigMap = make(map[string]interface{}) // TODO:needed?
	}

	var resources []proto.Message

	var escapeHatchCluster *envoy_cluster_v3.Cluster
	if !forMeshGateway {
		cfg, err := structs.ParseUpstreamConfigNoDefaults(upstreamConfigMap)
		if err != nil {
			// Don't hard fail on a config typo, just warn. The parse func returns
			// default config if there is an error so it's safe to continue.
			s.Logger.Warn("failed to parse", "upstream", uid,
				"error", err)
		}

		if cfg.EnvoyClusterJSON != "" {
			if chain.Default {
				// If you haven't done anything to setup the discovery chain, then
				// you can use the envoy_cluster_json escape hatch.
				escapeHatchCluster, err = makeClusterFromUserConfig(cfg.EnvoyClusterJSON)
				if err != nil {
					return resources, nil
				}
			} else {
				s.Logger.Warn("ignoring escape hatch setting, because a discovery chain is configued for",
					"discovery chain", chain.ServiceName, "upstream", uid,
					"envoy_cluster_json", chain.ServiceName)
			}
		}
	}

	// Find all resolver nodes.
	for _, node := range chain.Nodes {
		if node.Type != structs.DiscoveryGraphNodeTypeResolver {
			continue
		}
		failover := node.Resolver.Failover
		targetID := node.Resolver.Target

		target := chain.Targets[targetID]

		clusterName := CustomizeClusterName(target.Name, chain)
		if escapeHatchCluster != nil {
			clusterName = escapeHatchCluster.Name
		}
		if forMeshGateway {
			clusterName = meshGatewayExportedClusterNamePrefix + clusterName
		}
		s.Logger.Debug("generating endpoints for", "cluster", clusterName)

		// Determine if we have to generate the entire cluster differently.
		failoverThroughMeshGateway := chain.WillFailoverThroughMeshGateway(node) && !forMeshGateway

		if failoverThroughMeshGateway {
			actualTargetID := firstHealthyTarget(
				chain.Targets,
				upstreamEndpoints,
				targetID,
				failover.Targets,
			)
			if actualTargetID != targetID {
				targetID = actualTargetID
			}

			failover = nil
		}

		primaryGroup, valid := makeLoadAssignmentEndpointGroup(
			chain.Targets,
			upstreamEndpoints,
			gatewayEndpoints,
			targetID,
			gatewayKey,
			forMeshGateway,
		)
		if !valid {
			continue // skip the cluster if we're still populating the snapshot
		}

		var numFailoverTargets int
		if failover != nil {
			numFailoverTargets = len(failover.Targets)
		}

		endpointGroups := make([]loadAssignmentEndpointGroup, 0, numFailoverTargets+1)
		endpointGroups = append(endpointGroups, primaryGroup)

		if failover != nil && len(failover.Targets) > 0 {
			for _, failTargetID := range failover.Targets {
				failoverGroup, valid := makeLoadAssignmentEndpointGroup(
					chain.Targets,
					upstreamEndpoints,
					gatewayEndpoints,
					failTargetID,
					gatewayKey,
					forMeshGateway,
				)
				if !valid {
					continue // skip the failover target if we're still populating the snapshot
				}
				endpointGroups = append(endpointGroups, failoverGroup)
			}
		}

		la := makeLoadAssignment(
			clusterName,
			endpointGroups,
			gatewayKey,
		)
		resources = append(resources, la)
	}

	return resources, nil
}

func (s *ResourceGenerator) makeExportedUpstreamEndpointsForMeshGateway(cfgSnap *proxycfg.ConfigSnapshot) ([]proto.Message, error) {
	var resources []proto.Message

	populatedExportedClusters := make(map[string]struct{}) // key=clusterName
	for _, svc := range cfgSnap.MeshGatewayValidExportedServices() {
		chain := cfgSnap.MeshGateway.DiscoveryChain[svc]

		chainEndpoints := make(map[string]structs.CheckServiceNodes)
		for _, target := range chain.Targets {
			if !cfgSnap.Locality.Matches(target.Datacenter, target.Partition) {
				s.Logger.Warn("ignoring discovery chain target that crosses a datacenter or partition boundary in a mesh gateway",
					"target", target,
					"gatewayLocality", cfgSnap.Locality,
				)
				continue
			}

			targetSvc := target.ServiceName()

			endpoints, ok := cfgSnap.MeshGateway.ServiceGroups[targetSvc]
			if !ok {
				continue // ignore; not ready
			}

			if target.ServiceSubset == "" {
				chainEndpoints[target.ID] = endpoints
			} else {
				resolver, ok := cfgSnap.MeshGateway.ServiceResolvers[targetSvc]
				if !ok {
					continue // ignore; not ready
				}
				subset, ok := resolver.Subsets[target.ServiceSubset]
				if !ok {
					continue // ignore; not ready
				}

				subsetEndpoints, err := s.filterSubsetEndpoints(&subset, endpoints)
				if err != nil {
					return nil, err
				}
				chainEndpoints[target.ID] = subsetEndpoints
			}
		}

		clusterEndpoints, err := s.endpointsFromDiscoveryChain(
			proxycfg.NewUpstreamIDFromServiceName(svc),
			chain,
			cfgSnap.Locality,
			nil,
			chainEndpoints,
			nil,
			true,
		)
		if err != nil {
			return nil, err
		}
		for _, endpoints := range clusterEndpoints {
			clusterName := getResourceName(endpoints)
			if _, ok := populatedExportedClusters[clusterName]; ok {
				continue
			}
			populatedExportedClusters[clusterName] = struct{}{}
			resources = append(resources, endpoints)
		}
	}
	return resources, nil
}

type loadAssignmentEndpointGroup struct {
	Endpoints      structs.CheckServiceNodes
	OnlyPassing    bool
	OverrideHealth envoy_core_v3.HealthStatus
}

func makeLoadAssignment(clusterName string, endpointGroups []loadAssignmentEndpointGroup, localKey proxycfg.GatewayKey) *envoy_endpoint_v3.ClusterLoadAssignment {
	cla := &envoy_endpoint_v3.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints:   make([]*envoy_endpoint_v3.LocalityLbEndpoints, 0, len(endpointGroups)),
	}

	if len(endpointGroups) > 1 {
		cla.Policy = &envoy_endpoint_v3.ClusterLoadAssignment_Policy{
			// We choose such a large value here that the failover math should
			// in effect not happen until zero instances are healthy.
			OverprovisioningFactor: makeUint32Value(100000),
		}
	}

	for priority, endpointGroup := range endpointGroups {
		endpoints := endpointGroup.Endpoints
		es := make([]*envoy_endpoint_v3.LbEndpoint, 0, len(endpoints))

		for _, ep := range endpoints {
			// TODO (mesh-gateway) - should we respect the translate_wan_addrs configuration here or just always use the wan for cross-dc?
			_, addr, port := ep.BestAddress(!localKey.Matches(ep.Node.Datacenter, ep.Node.PartitionOrDefault()))
			healthStatus, weight := calculateEndpointHealthAndWeight(ep, endpointGroup.OnlyPassing)

			if endpointGroup.OverrideHealth != envoy_core_v3.HealthStatus_UNKNOWN {
				healthStatus = endpointGroup.OverrideHealth
			}

			es = append(es, &envoy_endpoint_v3.LbEndpoint{
				HostIdentifier: &envoy_endpoint_v3.LbEndpoint_Endpoint{
					Endpoint: &envoy_endpoint_v3.Endpoint{
						Address: makeAddress(addr, port),
					},
				},
				HealthStatus:        healthStatus,
				LoadBalancingWeight: makeUint32Value(weight),
			})
		}

		cla.Endpoints = append(cla.Endpoints, &envoy_endpoint_v3.LocalityLbEndpoints{
			Priority:    uint32(priority),
			LbEndpoints: es,
		})
	}

	return cla
}

func makeLoadAssignmentEndpointGroup(
	targets map[string]*structs.DiscoveryTarget,
	targetHealth map[string]structs.CheckServiceNodes,
	gatewayHealth map[string]structs.CheckServiceNodes,
	targetID string,
	localKey proxycfg.GatewayKey,
	forMeshGateway bool,
) (loadAssignmentEndpointGroup, bool) {
	realEndpoints, ok := targetHealth[targetID]
	if !ok {
		// skip the cluster if we're still populating the snapshot
		return loadAssignmentEndpointGroup{}, false
	}
	target := targets[targetID]

	var gatewayKey proxycfg.GatewayKey

	switch target.MeshGateway.Mode {
	case structs.MeshGatewayModeRemote:
		gatewayKey.Datacenter = target.Datacenter
		gatewayKey.Partition = target.Partition
	case structs.MeshGatewayModeLocal:
		gatewayKey = localKey
	}

	if forMeshGateway || gatewayKey.IsEmpty() || localKey.Matches(target.Datacenter, target.Partition) {
		// Gateways are not needed if the request isn't for a remote DC or partition.
		return loadAssignmentEndpointGroup{
			Endpoints:   realEndpoints,
			OnlyPassing: target.Subset.OnlyPassing,
		}, true
	}

	// If using a mesh gateway we need to pull those endpoints instead.
	gatewayEndpoints, ok := gatewayHealth[gatewayKey.String()]
	if !ok {
		// skip the cluster if we're still populating the snapshot
		return loadAssignmentEndpointGroup{}, false
	}

	// But we will use the health from the actual backend service.
	overallHealth := envoy_core_v3.HealthStatus_UNHEALTHY
	for _, ep := range realEndpoints {
		health, _ := calculateEndpointHealthAndWeight(ep, target.Subset.OnlyPassing)
		if health == envoy_core_v3.HealthStatus_HEALTHY {
			overallHealth = envoy_core_v3.HealthStatus_HEALTHY
			break
		}
	}

	return loadAssignmentEndpointGroup{
		Endpoints:      gatewayEndpoints,
		OverrideHealth: overallHealth,
	}, true
}

func calculateEndpointHealthAndWeight(
	ep structs.CheckServiceNode,
	onlyPassing bool,
) (envoy_core_v3.HealthStatus, int) {
	healthStatus := envoy_core_v3.HealthStatus_HEALTHY
	weight := 1
	if ep.Service.Weights != nil {
		weight = ep.Service.Weights.Passing
	}

	for _, chk := range ep.Checks {
		if chk.Status == api.HealthCritical {
			healthStatus = envoy_core_v3.HealthStatus_UNHEALTHY
		}
		if onlyPassing && chk.Status != api.HealthPassing {
			healthStatus = envoy_core_v3.HealthStatus_UNHEALTHY
		}
		if chk.Status == api.HealthWarning && ep.Service.Weights != nil {
			weight = ep.Service.Weights.Warning
		}
	}
	// Make weights fit Envoy's limits. A zero weight means that either Warning
	// (likely) or Passing (weirdly) weight has been set to 0 effectively making
	// this instance unhealthy and should not be sent traffic.
	if weight < 1 {
		healthStatus = envoy_core_v3.HealthStatus_UNHEALTHY
		weight = 1
	}
	if weight > 128 {
		weight = 128
	}
	return healthStatus, weight
}
