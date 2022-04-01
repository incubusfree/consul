//go:build !consulent
// +build !consulent

package xds

import (
	"path/filepath"
	"sort"
	"testing"

	envoy_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/golang/protobuf/proto"
	testinf "github.com/mitchellh/go-testing-interface"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/consul/agent/proxycfg"
	"github.com/hashicorp/consul/agent/xds/proxysupport"
	"github.com/hashicorp/consul/agent/xds/serverlessplugin"
	"github.com/hashicorp/consul/agent/xds/xdscommon"
	"github.com/hashicorp/consul/sdk/testutil"
)

func TestServerlessPluginFromSnapshot(t *testing.T) {
	tests := []struct {
		name   string
		create func(t testinf.T) *proxycfg.ConfigSnapshot
	}{
		{
			name:   "lambda-terminating-gateway",
			create: proxycfg.TestConfigSnapshotTerminatingGatewayWithLambdaService,
		},
	}

	latestEnvoyVersion := proxysupport.EnvoyVersions[0]
	for _, envoyVersion := range proxysupport.EnvoyVersions {
		sf, err := determineSupportedProxyFeaturesFromString(envoyVersion)
		require.NoError(t, err)
		t.Run("envoy-"+envoyVersion, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					// Sanity check default with no overrides first
					snap := tt.create(t)

					// We need to replace the TLS certs with deterministic ones to make golden
					// files workable. Note we don't update these otherwise they'd change
					// golden files for every test case and so not be any use!
					setupTLSRootsAndLeaf(t, snap)

					g := newResourceGenerator(testutil.Logger(t), nil, nil, false)
					g.ProxyFeatures = sf

					res, err := g.allResourcesFromSnapshot(snap)
					require.NoError(t, err)

					indexedResources := indexResources(g.Logger, res)
					newResourceMap, err := serverlessplugin.MutateIndexedResources(indexedResources, xdscommon.MakePluginConfiguration(snap))
					require.NoError(t, err)

					entities := []struct {
						name   string
						key    string
						sorter func([]proto.Message) func(int, int) bool
					}{
						{
							name: "clusters",
							key:  xdscommon.ClusterType,
							sorter: func(msgs []proto.Message) func(int, int) bool {
								return func(i, j int) bool {
									return msgs[i].(*envoy_cluster_v3.Cluster).Name < msgs[j].(*envoy_cluster_v3.Cluster).Name
								}
							},
						},
						{
							name: "listeners",
							key:  xdscommon.ListenerType,
							sorter: func(msgs []proto.Message) func(int, int) bool {
								return func(i, j int) bool {
									return msgs[i].(*envoy_listener_v3.Listener).Name < msgs[j].(*envoy_listener_v3.Listener).Name
								}
							},
						},
					}

					for _, entity := range entities {
						var msgs []proto.Message
						for _, e := range newResourceMap.Index[entity.key] {
							msgs = append(msgs, e)
						}

						sort.Slice(msgs, entity.sorter(msgs))
						r, err := createResponse(entity.key, "00000001", "00000001", msgs)
						require.NoError(t, err)

						t.Run(entity.name, func(t *testing.T) {
							gotJSON := protoToJSON(t, r)

							require.JSONEq(t, goldenEnvoy(t,
								filepath.Join("serverless_plugin", entity.name, tt.name),
								envoyVersion, latestEnvoyVersion, gotJSON), gotJSON)
						})
					}
				})
			}
		})
	}
}
