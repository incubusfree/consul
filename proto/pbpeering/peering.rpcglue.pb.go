// Code generated by proto-gen-rpc-glue. DO NOT EDIT.

package pbpeering

import (
	"time"

	"github.com/hashicorp/consul/agent/structs"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ structs.RPCInfo
var _ time.Month

// RequestDatacenter implements structs.RPCInfo
func (msg *PeeringReadRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}

// IsRead implements structs.RPCInfo
func (msg *PeeringReadRequest) IsRead() bool {
	// TODO(peering): figure out read semantics here
	return true
}

// AllowStaleRead implements structs.RPCInfo
func (msg *PeeringReadRequest) AllowStaleRead() bool {
	// TODO(peering): figure out read semantics here
	// TODO(peering): this needs to stay false for calls to head to the leader until we sync stream tracker information
	// like ImportedServicesCount, ExportedServicesCount, as well as general Status fields thru raft to make available
	// to followers as well
	return false
}

// HasTimedOut implements structs.RPCInfo
func (msg *PeeringReadRequest) HasTimedOut(start time.Time, rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) (bool, error) {
	// TODO(peering): figure out read semantics here
	return time.Since(start) > rpcHoldTimeout, nil
}

// Timeout implements structs.RPCInfo
func (msg *PeeringReadRequest) Timeout(rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) time.Duration {
	// TODO(peering): figure out read semantics here
	return rpcHoldTimeout
}

// SetTokenSecret implements structs.RPCInfo
func (msg *PeeringReadRequest) SetTokenSecret(s string) {
	// TODO(peering): figure out read semantics here
}

// TokenSecret implements structs.RPCInfo
func (msg *PeeringReadRequest) TokenSecret() string {
	// TODO(peering): figure out read semantics here
	return ""
}

// Token implements structs.RPCInfo
func (msg *PeeringReadRequest) Token() string {
	// TODO(peering): figure out read semantics here
	return ""
}

// RequestDatacenter implements structs.RPCInfo
func (msg *PeeringListRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}

// IsRead implements structs.RPCInfo
func (msg *PeeringListRequest) IsRead() bool {
	// TODO(peering): figure out read semantics here
	return true
}

// AllowStaleRead implements structs.RPCInfo
func (msg *PeeringListRequest) AllowStaleRead() bool {
	// TODO(peering): figure out read semantics here
	// TODO(peering): this needs to stay false for calls to head to the leader until we sync stream tracker information
	// like ImportedServicesCount, ExportedServicesCount, as well as general Status fields thru raft to make available
	// to followers as well
	return false
}

// HasTimedOut implements structs.RPCInfo
func (msg *PeeringListRequest) HasTimedOut(start time.Time, rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) (bool, error) {
	// TODO(peering): figure out read semantics here
	return time.Since(start) > rpcHoldTimeout, nil
}

// Timeout implements structs.RPCInfo
func (msg *PeeringListRequest) Timeout(rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) time.Duration {
	// TODO(peering): figure out read semantics here
	return rpcHoldTimeout
}

// SetTokenSecret implements structs.RPCInfo
func (msg *PeeringListRequest) SetTokenSecret(s string) {
	// TODO(peering): figure out read semantics here
}

// TokenSecret implements structs.RPCInfo
func (msg *PeeringListRequest) TokenSecret() string {
	// TODO(peering): figure out read semantics here
	return ""
}

// Token implements structs.RPCInfo
func (msg *PeeringListRequest) Token() string {
	// TODO(peering): figure out read semantics here
	return ""
}

// RequestDatacenter implements structs.RPCInfo
func (msg *PeeringWriteRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}

// IsRead implements structs.RPCInfo
func (msg *PeeringWriteRequest) IsRead() bool {
	// TODO(peering): figure out write semantics here
	return false
}

// AllowStaleRead implements structs.RPCInfo
func (msg *PeeringWriteRequest) AllowStaleRead() bool {
	// TODO(peering): figure out write semantics here
	return false
}

// HasTimedOut implements structs.RPCInfo
func (msg *PeeringWriteRequest) HasTimedOut(start time.Time, rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) (bool, error) {
	// TODO(peering): figure out write semantics here
	return time.Since(start) > rpcHoldTimeout, nil
}

// Timeout implements structs.RPCInfo
func (msg *PeeringWriteRequest) Timeout(rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) time.Duration {
	// TODO(peering): figure out write semantics here
	return rpcHoldTimeout
}

// SetTokenSecret implements structs.RPCInfo
func (msg *PeeringWriteRequest) SetTokenSecret(s string) {
	// TODO(peering): figure out write semantics here
}

// TokenSecret implements structs.RPCInfo
func (msg *PeeringWriteRequest) TokenSecret() string {
	// TODO(peering): figure out write semantics here
	return ""
}

// RequestDatacenter implements structs.RPCInfo
func (msg *PeeringDeleteRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}

// IsRead implements structs.RPCInfo
func (msg *PeeringDeleteRequest) IsRead() bool {
	// TODO(peering): figure out write semantics here
	return false
}

// AllowStaleRead implements structs.RPCInfo
func (msg *PeeringDeleteRequest) AllowStaleRead() bool {
	// TODO(peering): figure out write semantics here
	return false
}

// HasTimedOut implements structs.RPCInfo
func (msg *PeeringDeleteRequest) HasTimedOut(start time.Time, rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) (bool, error) {
	// TODO(peering): figure out write semantics here
	return time.Since(start) > rpcHoldTimeout, nil
}

// Timeout implements structs.RPCInfo
func (msg *PeeringDeleteRequest) Timeout(rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) time.Duration {
	// TODO(peering): figure out write semantics here
	return rpcHoldTimeout
}

// SetTokenSecret implements structs.RPCInfo
func (msg *PeeringDeleteRequest) SetTokenSecret(s string) {
	// TODO(peering): figure out write semantics here
}

// TokenSecret implements structs.RPCInfo
func (msg *PeeringDeleteRequest) TokenSecret() string {
	// TODO(peering): figure out write semantics here
	return ""
}

// RequestDatacenter implements structs.RPCInfo
func (msg *TrustBundleListByServiceRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}

// IsRead implements structs.RPCInfo
func (msg *TrustBundleListByServiceRequest) IsRead() bool {
	// TODO(peering): figure out read semantics here
	return true
}

// AllowStaleRead implements structs.RPCInfo
func (msg *TrustBundleListByServiceRequest) AllowStaleRead() bool {
	// TODO(peering): figure out read semantics here
	return false
}

// HasTimedOut implements structs.RPCInfo
func (msg *TrustBundleListByServiceRequest) HasTimedOut(start time.Time, rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) (bool, error) {
	// TODO(peering): figure out read semantics here
	return time.Since(start) > rpcHoldTimeout, nil
}

// Timeout implements structs.RPCInfo
func (msg *TrustBundleListByServiceRequest) Timeout(rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) time.Duration {
	// TODO(peering): figure out read semantics here
	return rpcHoldTimeout
}

// SetTokenSecret implements structs.RPCInfo
func (msg *TrustBundleListByServiceRequest) SetTokenSecret(s string) {
	// TODO(peering): figure out read semantics here
}

// TokenSecret implements structs.RPCInfo
func (msg *TrustBundleListByServiceRequest) TokenSecret() string {
	// TODO(peering): figure out read semantics here
	return ""
}

// Token implements structs.RPCInfo
func (msg *TrustBundleListByServiceRequest) Token() string {
	// TODO(peering): figure out read semantics here
	return ""
}

// RequestDatacenter implements structs.RPCInfo
func (msg *TrustBundleReadRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}

// IsRead implements structs.RPCInfo
func (msg *TrustBundleReadRequest) IsRead() bool {
	// TODO(peering): figure out read semantics here
	return true
}

// AllowStaleRead implements structs.RPCInfo
func (msg *TrustBundleReadRequest) AllowStaleRead() bool {
	// TODO(peering): figure out read semantics here
	return false
}

// HasTimedOut implements structs.RPCInfo
func (msg *TrustBundleReadRequest) HasTimedOut(start time.Time, rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) (bool, error) {
	// TODO(peering): figure out read semantics here
	return time.Since(start) > rpcHoldTimeout, nil
}

// Timeout implements structs.RPCInfo
func (msg *TrustBundleReadRequest) Timeout(rpcHoldTimeout time.Duration, a time.Duration, b time.Duration) time.Duration {
	// TODO(peering): figure out read semantics here
	return rpcHoldTimeout
}

// SetTokenSecret implements structs.RPCInfo
func (msg *TrustBundleReadRequest) SetTokenSecret(s string) {
	// TODO(peering): figure out read semantics here
}

// TokenSecret implements structs.RPCInfo
func (msg *TrustBundleReadRequest) TokenSecret() string {
	// TODO(peering): figure out read semantics here
	return ""
}

// Token implements structs.RPCInfo
func (msg *TrustBundleReadRequest) Token() string {
	// TODO(peering): figure out read semantics here
	return ""
}

// RequestDatacenter implements structs.RPCInfo
func (msg *PeeringTrustBundleWriteRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}

// RequestDatacenter implements structs.RPCInfo
func (msg *PeeringTrustBundleDeleteRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}

// RequestDatacenter implements structs.RPCInfo
func (msg *EstablishRequest) RequestDatacenter() string {
	if msg == nil {
		return ""
	}
	return msg.Datacenter
}
