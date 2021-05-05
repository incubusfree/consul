// Code generated by mog. DO NOT EDIT.

package pbservice

import structs "github.com/hashicorp/consul/agent/structs"

func ConnectProxyConfigToStructs(s ConnectProxyConfig) structs.ConnectProxyConfig {
	var t structs.ConnectProxyConfig
	t.DestinationServiceName = s.DestinationServiceName
	t.DestinationServiceID = s.DestinationServiceID
	t.LocalServiceAddress = s.LocalServiceAddress
	t.LocalServicePort = int(s.LocalServicePort)
	t.LocalServiceSocketPath = s.LocalServiceSocketPath
	t.Mode = s.Mode
	t.Config = ProtobufTypesStructToMapStringInterface(s.Config)
	t.Upstreams = UpstreamsToStructs(s.Upstreams)
	t.MeshGateway = MeshGatewayConfigToStructs(s.MeshGateway)
	t.Expose = ExposeConfigToStructs(s.Expose)
	t.TransparentProxy = TransparentProxyConfigToStructs(s.TransparentProxy)
	return t
}
func NewConnectProxyConfigFromStructs(t structs.ConnectProxyConfig) ConnectProxyConfig {
	var s ConnectProxyConfig
	s.DestinationServiceName = t.DestinationServiceName
	s.DestinationServiceID = t.DestinationServiceID
	s.LocalServiceAddress = t.LocalServiceAddress
	s.LocalServicePort = int32(t.LocalServicePort)
	s.LocalServiceSocketPath = t.LocalServiceSocketPath
	s.Mode = t.Mode
	s.Config = MapStringInterfaceToProtobufTypesStruct(t.Config)
	s.Upstreams = NewUpstreamsFromStructs(t.Upstreams)
	s.MeshGateway = NewMeshGatewayConfigFromStructs(t.MeshGateway)
	s.Expose = NewExposeConfigFromStructs(t.Expose)
	s.TransparentProxy = NewTransparentProxyConfigFromStructs(t.TransparentProxy)
	return s
}
func ExposeConfigToStructs(s ExposeConfig) structs.ExposeConfig {
	var t structs.ExposeConfig
	t.Checks = s.Checks
	t.Paths = ExposePathSliceToStructs(s.Paths)
	return t
}
func NewExposeConfigFromStructs(t structs.ExposeConfig) ExposeConfig {
	var s ExposeConfig
	s.Checks = t.Checks
	s.Paths = NewExposePathSliceFromStructs(t.Paths)
	return s
}
func ExposePathToStructs(s ExposePath) structs.ExposePath {
	var t structs.ExposePath
	t.ListenerPort = int(s.ListenerPort)
	t.Path = s.Path
	t.LocalPathPort = int(s.LocalPathPort)
	t.Protocol = s.Protocol
	t.ParsedFromCheck = s.ParsedFromCheck
	return t
}
func NewExposePathFromStructs(t structs.ExposePath) ExposePath {
	var s ExposePath
	s.ListenerPort = int32(t.ListenerPort)
	s.Path = t.Path
	s.LocalPathPort = int32(t.LocalPathPort)
	s.Protocol = t.Protocol
	s.ParsedFromCheck = t.ParsedFromCheck
	return s
}
func MeshGatewayConfigToStructs(s MeshGatewayConfig) structs.MeshGatewayConfig {
	var t structs.MeshGatewayConfig
	t.Mode = s.Mode
	return t
}
func NewMeshGatewayConfigFromStructs(t structs.MeshGatewayConfig) MeshGatewayConfig {
	var s MeshGatewayConfig
	s.Mode = t.Mode
	return s
}
func ServiceConnectToStructs(s ServiceConnect) structs.ServiceConnect {
	var t structs.ServiceConnect
	t.Native = s.Native
	t.SidecarService = ServiceDefinitionPtrToStructs(s.SidecarService)
	return t
}
func NewServiceConnectFromStructs(t structs.ServiceConnect) ServiceConnect {
	var s ServiceConnect
	s.Native = t.Native
	s.SidecarService = NewServiceDefinitionPtrFromStructs(t.SidecarService)
	return s
}
func ServiceDefinitionToStructs(s ServiceDefinition) structs.ServiceDefinition {
	var t structs.ServiceDefinition
	t.Kind = s.Kind
	t.ID = s.ID
	t.Name = s.Name
	t.Tags = s.Tags
	t.Address = s.Address
	t.TaggedAddresses = MapStringServiceAddressToStructs(s.TaggedAddresses)
	t.Meta = s.Meta
	t.Port = int(s.Port)
	t.SocketPath = s.SocketPath
	t.Check = CheckTypeToStructs(s.Check)
	t.Checks = CheckTypesToStructs(s.Checks)
	t.Weights = WeightsPtrToStructs(s.Weights)
	t.Token = s.Token
	t.EnableTagOverride = s.EnableTagOverride
	t.Proxy = ConnectProxyConfigPtrToStructs(s.Proxy)
	t.EnterpriseMeta = EnterpriseMetaToStructs(s.EnterpriseMeta)
	t.Connect = ServiceConnectPtrToStructs(s.Connect)
	return t
}
func NewServiceDefinitionFromStructs(t structs.ServiceDefinition) ServiceDefinition {
	var s ServiceDefinition
	s.Kind = t.Kind
	s.ID = t.ID
	s.Name = t.Name
	s.Tags = t.Tags
	s.Address = t.Address
	s.TaggedAddresses = NewMapStringServiceAddressFromStructs(t.TaggedAddresses)
	s.Meta = t.Meta
	s.Port = int32(t.Port)
	s.SocketPath = t.SocketPath
	s.Check = NewCheckTypeFromStructs(t.Check)
	s.Checks = NewCheckTypesFromStructs(t.Checks)
	s.Weights = NewWeightsPtrFromStructs(t.Weights)
	s.Token = t.Token
	s.EnableTagOverride = t.EnableTagOverride
	s.Proxy = NewConnectProxyConfigPtrFromStructs(t.Proxy)
	s.EnterpriseMeta = NewEnterpriseMetaFromStructs(t.EnterpriseMeta)
	s.Connect = NewServiceConnectPtrFromStructs(t.Connect)
	return s
}
func TransparentProxyConfigToStructs(s TransparentProxyConfig) structs.TransparentProxyConfig {
	var t structs.TransparentProxyConfig
	t.OutboundListenerPort = int(s.OutboundListenerPort)
	return t
}
func NewTransparentProxyConfigFromStructs(t structs.TransparentProxyConfig) TransparentProxyConfig {
	var s TransparentProxyConfig
	s.OutboundListenerPort = int32(t.OutboundListenerPort)
	return s
}
func UpstreamToStructs(s Upstream) structs.Upstream {
	var t structs.Upstream
	t.DestinationType = s.DestinationType
	t.DestinationNamespace = s.DestinationNamespace
	t.DestinationName = s.DestinationName
	t.Datacenter = s.Datacenter
	t.LocalBindAddress = s.LocalBindAddress
	t.LocalBindPort = int(s.LocalBindPort)
	t.LocalBindSocketPath = s.LocalBindSocketPath
	t.LocalBindSocketMode = s.LocalBindSocketMode
	t.Config = ProtobufTypesStructToMapStringInterface(s.Config)
	t.MeshGateway = MeshGatewayConfigToStructs(s.MeshGateway)
	t.CentrallyConfigured = s.CentrallyConfigured
	return t
}
func NewUpstreamFromStructs(t structs.Upstream) Upstream {
	var s Upstream
	s.DestinationType = t.DestinationType
	s.DestinationNamespace = t.DestinationNamespace
	s.DestinationName = t.DestinationName
	s.Datacenter = t.Datacenter
	s.LocalBindAddress = t.LocalBindAddress
	s.LocalBindPort = int32(t.LocalBindPort)
	s.LocalBindSocketPath = t.LocalBindSocketPath
	s.LocalBindSocketMode = t.LocalBindSocketMode
	s.Config = MapStringInterfaceToProtobufTypesStruct(t.Config)
	s.MeshGateway = NewMeshGatewayConfigFromStructs(t.MeshGateway)
	s.CentrallyConfigured = t.CentrallyConfigured
	return s
}
