// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/pbconnect/connect.proto

package pbconnect

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	pbcommon "github.com/hashicorp/consul/proto/pbcommon"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// CARoots is the list of all currently trusted CA Roots.
//
// mog annotation:
//
// target=github.com/hashicorp/consul/agent/structs.IndexedCARoots
// output=connect.gen.go
// name=StructsIndexedCARoots
type CARoots struct {
	// ActiveRootID is the ID of a root in Roots that is the active CA root.
	// Other roots are still valid if they're in the Roots list but are in
	// the process of being rotated out.
	ActiveRootID string `protobuf:"bytes,1,opt,name=ActiveRootID,proto3" json:"ActiveRootID,omitempty"`
	// TrustDomain is the identification root for this Consul cluster. All
	// certificates signed by the cluster's CA must have their identifying URI in
	// this domain.
	//
	// This does not include the protocol (currently spiffe://) since we may
	// implement other protocols in future with equivalent semantics. It should be
	// compared against the "authority" section of a URI (i.e. host:port).
	//
	// We need to support migrating a cluster between trust domains to support
	// Multi-DC migration in Enterprise. In this case the current trust domain is
	// here but entries in Roots may also have ExternalTrustDomain set to a
	// non-empty value implying they were previous roots that are still trusted
	// but under a different trust domain.
	//
	// Note that we DON'T validate trust domain during AuthZ since it causes
	// issues of loss of connectivity during migration between trust domains. The
	// only time the additional validation adds value is where the cluster shares
	// an external root (e.g. organization-wide root) with another distinct Consul
	// cluster or PKI system. In this case, x509 Name Constraints can be added to
	// enforce that Consul's CA can only validly sign or trust certs within the
	// same trust-domain. Name constraints as enforced by TLS handshake also allow
	// seamless rotation between trust domains thanks to cross-signing.
	TrustDomain string `protobuf:"bytes,2,opt,name=TrustDomain,proto3" json:"TrustDomain,omitempty"`
	// Roots is a list of root CA certs to trust.
	Roots []*CARoot `protobuf:"bytes,3,rep,name=Roots,proto3" json:"Roots,omitempty"`
	// QueryMeta here is mainly used to contain the latest Raft Index that could
	// be used to perform a blocking query.
	// mog: func-to=QueryMetaTo func-from=QueryMetaFrom
	QueryMeta            *pbcommon.QueryMeta `protobuf:"bytes,4,opt,name=QueryMeta,proto3" json:"QueryMeta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *CARoots) Reset()         { *m = CARoots{} }
func (m *CARoots) String() string { return proto.CompactTextString(m) }
func (*CARoots) ProtoMessage()    {}
func (*CARoots) Descriptor() ([]byte, []int) {
	return fileDescriptor_80627e709958eb04, []int{0}
}

func (m *CARoots) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CARoots.Unmarshal(m, b)
}
func (m *CARoots) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CARoots.Marshal(b, m, deterministic)
}
func (m *CARoots) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CARoots.Merge(m, src)
}
func (m *CARoots) XXX_Size() int {
	return xxx_messageInfo_CARoots.Size(m)
}
func (m *CARoots) XXX_DiscardUnknown() {
	xxx_messageInfo_CARoots.DiscardUnknown(m)
}

var xxx_messageInfo_CARoots proto.InternalMessageInfo

func (m *CARoots) GetActiveRootID() string {
	if m != nil {
		return m.ActiveRootID
	}
	return ""
}

func (m *CARoots) GetTrustDomain() string {
	if m != nil {
		return m.TrustDomain
	}
	return ""
}

func (m *CARoots) GetRoots() []*CARoot {
	if m != nil {
		return m.Roots
	}
	return nil
}

func (m *CARoots) GetQueryMeta() *pbcommon.QueryMeta {
	if m != nil {
		return m.QueryMeta
	}
	return nil
}

// CARoot is the trusted CA Root.
//
// mog annotation:
//
// target=github.com/hashicorp/consul/agent/structs.CARoot
// output=connect.gen.go
// name=StructsCARoot
type CARoot struct {
	// ID is a globally unique ID (UUID) representing this CA root.
	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	// Name is a human-friendly name for this CA root. This value is
	// opaque to Consul and is not used for anything internally.
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	// SerialNumber is the x509 serial number of the certificate.
	SerialNumber uint64 `protobuf:"varint,3,opt,name=SerialNumber,proto3" json:"SerialNumber,omitempty"`
	// SigningKeyID is the ID of the public key that corresponds to the private
	// key used to sign leaf certificates. Is is the HexString format of the
	// raw AuthorityKeyID bytes.
	SigningKeyID string `protobuf:"bytes,4,opt,name=SigningKeyID,proto3" json:"SigningKeyID,omitempty"`
	// ExternalTrustDomain is the trust domain this root was generated under. It
	// is usually empty implying "the current cluster trust-domain". It is set
	// only in the case that a cluster changes trust domain and then all old roots
	// that are still trusted have the old trust domain set here.
	//
	// We currently DON'T validate these trust domains explicitly anywhere, see
	// IndexedRoots.TrustDomain doc. We retain this information for debugging and
	// future flexibility.
	ExternalTrustDomain string `protobuf:"bytes,5,opt,name=ExternalTrustDomain,proto3" json:"ExternalTrustDomain,omitempty"`
	// Time validity bounds.
	// mog: func-to=structs.TimeFromProto func-from=structs.TimeToProto
	NotBefore *timestamp.Timestamp `protobuf:"bytes,6,opt,name=NotBefore,proto3" json:"NotBefore,omitempty"`
	// mog: func-to=structs.TimeFromProto func-from=structs.TimeToProto
	NotAfter *timestamp.Timestamp `protobuf:"bytes,7,opt,name=NotAfter,proto3" json:"NotAfter,omitempty"`
	// RootCert is the PEM-encoded public certificate.
	RootCert string `protobuf:"bytes,8,opt,name=RootCert,proto3" json:"RootCert,omitempty"`
	// IntermediateCerts is a list of PEM-encoded intermediate certs to
	// attach to any leaf certs signed by this CA.
	IntermediateCerts []string `protobuf:"bytes,9,rep,name=IntermediateCerts,proto3" json:"IntermediateCerts,omitempty"`
	// SigningCert is the PEM-encoded signing certificate and SigningKey
	// is the PEM-encoded private key for the signing certificate. These
	// may actually be empty if the CA plugin in use manages these for us.
	SigningCert string `protobuf:"bytes,10,opt,name=SigningCert,proto3" json:"SigningCert,omitempty"`
	SigningKey  string `protobuf:"bytes,11,opt,name=SigningKey,proto3" json:"SigningKey,omitempty"`
	// Active is true if this is the current active CA. This must only
	// be true for exactly one CA. For any method that modifies roots in the
	// state store, tests should be written to verify that multiple roots
	// cannot be active.
	Active bool `protobuf:"varint,12,opt,name=Active,proto3" json:"Active,omitempty"`
	// RotatedOutAt is the time at which this CA was removed from the state.
	// This will only be set on roots that have been rotated out from being the
	// active root.
	// mog: func-to=structs.TimeFromProto func-from=structs.TimeToProto
	RotatedOutAt *timestamp.Timestamp `protobuf:"bytes,13,opt,name=RotatedOutAt,proto3" json:"RotatedOutAt,omitempty"`
	// PrivateKeyType is the type of the private key used to sign certificates. It
	// may be "rsa" or "ec". This is provided as a convenience to avoid parsing
	// the public key to from the certificate to infer the type.
	PrivateKeyType string `protobuf:"bytes,14,opt,name=PrivateKeyType,proto3" json:"PrivateKeyType,omitempty"`
	// PrivateKeyBits is the length of the private key used to sign certificates.
	// This is provided as a convenience to avoid parsing the public key from the
	// certificate to infer the type.
	// mog: func-to=int func-from=int32
	PrivateKeyBits int32 `protobuf:"varint,15,opt,name=PrivateKeyBits,proto3" json:"PrivateKeyBits,omitempty"`
	// mog: func-to=RaftIndexTo func-from=RaftIndexFrom
	RaftIndex            *pbcommon.RaftIndex `protobuf:"bytes,16,opt,name=RaftIndex,proto3" json:"RaftIndex,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *CARoot) Reset()         { *m = CARoot{} }
func (m *CARoot) String() string { return proto.CompactTextString(m) }
func (*CARoot) ProtoMessage()    {}
func (*CARoot) Descriptor() ([]byte, []int) {
	return fileDescriptor_80627e709958eb04, []int{1}
}

func (m *CARoot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CARoot.Unmarshal(m, b)
}
func (m *CARoot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CARoot.Marshal(b, m, deterministic)
}
func (m *CARoot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CARoot.Merge(m, src)
}
func (m *CARoot) XXX_Size() int {
	return xxx_messageInfo_CARoot.Size(m)
}
func (m *CARoot) XXX_DiscardUnknown() {
	xxx_messageInfo_CARoot.DiscardUnknown(m)
}

var xxx_messageInfo_CARoot proto.InternalMessageInfo

func (m *CARoot) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *CARoot) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CARoot) GetSerialNumber() uint64 {
	if m != nil {
		return m.SerialNumber
	}
	return 0
}

func (m *CARoot) GetSigningKeyID() string {
	if m != nil {
		return m.SigningKeyID
	}
	return ""
}

func (m *CARoot) GetExternalTrustDomain() string {
	if m != nil {
		return m.ExternalTrustDomain
	}
	return ""
}

func (m *CARoot) GetNotBefore() *timestamp.Timestamp {
	if m != nil {
		return m.NotBefore
	}
	return nil
}

func (m *CARoot) GetNotAfter() *timestamp.Timestamp {
	if m != nil {
		return m.NotAfter
	}
	return nil
}

func (m *CARoot) GetRootCert() string {
	if m != nil {
		return m.RootCert
	}
	return ""
}

func (m *CARoot) GetIntermediateCerts() []string {
	if m != nil {
		return m.IntermediateCerts
	}
	return nil
}

func (m *CARoot) GetSigningCert() string {
	if m != nil {
		return m.SigningCert
	}
	return ""
}

func (m *CARoot) GetSigningKey() string {
	if m != nil {
		return m.SigningKey
	}
	return ""
}

func (m *CARoot) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
}

func (m *CARoot) GetRotatedOutAt() *timestamp.Timestamp {
	if m != nil {
		return m.RotatedOutAt
	}
	return nil
}

func (m *CARoot) GetPrivateKeyType() string {
	if m != nil {
		return m.PrivateKeyType
	}
	return ""
}

func (m *CARoot) GetPrivateKeyBits() int32 {
	if m != nil {
		return m.PrivateKeyBits
	}
	return 0
}

func (m *CARoot) GetRaftIndex() *pbcommon.RaftIndex {
	if m != nil {
		return m.RaftIndex
	}
	return nil
}

// RaftIndex is used to track the index used while creating
// or modifying a given struct type.
//
// mog annotation:
//
// target=github.com/hashicorp/consul/agent/structs.IssuedCert
// output=connect.gen.go
// name=StructsIssuedCert
type IssuedCert struct {
	// SerialNumber is the unique serial number for this certificate.
	// This is encoded in standard hex separated by :.
	SerialNumber string `protobuf:"bytes,1,opt,name=SerialNumber,proto3" json:"SerialNumber,omitempty"`
	// CertPEM and PrivateKeyPEM are the PEM-encoded certificate and private
	// key for that cert, respectively. This should not be stored in the
	// state store, but is present in the sign API response.
	CertPEM       string `protobuf:"bytes,2,opt,name=CertPEM,proto3" json:"CertPEM,omitempty"`
	PrivateKeyPEM string `protobuf:"bytes,3,opt,name=PrivateKeyPEM,proto3" json:"PrivateKeyPEM,omitempty"`
	// Service is the name of the service for which the cert was issued.
	// ServiceURI is the cert URI value.
	Service    string `protobuf:"bytes,4,opt,name=Service,proto3" json:"Service,omitempty"`
	ServiceURI string `protobuf:"bytes,5,opt,name=ServiceURI,proto3" json:"ServiceURI,omitempty"`
	// Agent is the name of the node for which the cert was issued.
	// AgentURI is the cert URI value.
	Agent    string `protobuf:"bytes,6,opt,name=Agent,proto3" json:"Agent,omitempty"`
	AgentURI string `protobuf:"bytes,7,opt,name=AgentURI,proto3" json:"AgentURI,omitempty"`
	// ValidAfter and ValidBefore are the validity periods for the
	// certificate.
	// mog: func-to=structs.TimeFromProto func-from=structs.TimeToProto
	ValidAfter *timestamp.Timestamp `protobuf:"bytes,8,opt,name=ValidAfter,proto3" json:"ValidAfter,omitempty"`
	// mog: func-to=structs.TimeFromProto func-from=structs.TimeToProto
	ValidBefore *timestamp.Timestamp `protobuf:"bytes,9,opt,name=ValidBefore,proto3" json:"ValidBefore,omitempty"`
	// EnterpriseMeta is the Consul Enterprise specific metadata
	// mog: func-to=EnterpriseMetaTo func-from=EnterpriseMetaFrom
	EnterpriseMeta *pbcommon.EnterpriseMeta `protobuf:"bytes,10,opt,name=EnterpriseMeta,proto3" json:"EnterpriseMeta,omitempty"`
	// mog: func-to=RaftIndexTo func-from=RaftIndexFrom
	RaftIndex            *pbcommon.RaftIndex `protobuf:"bytes,11,opt,name=RaftIndex,proto3" json:"RaftIndex,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *IssuedCert) Reset()         { *m = IssuedCert{} }
func (m *IssuedCert) String() string { return proto.CompactTextString(m) }
func (*IssuedCert) ProtoMessage()    {}
func (*IssuedCert) Descriptor() ([]byte, []int) {
	return fileDescriptor_80627e709958eb04, []int{2}
}

func (m *IssuedCert) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IssuedCert.Unmarshal(m, b)
}
func (m *IssuedCert) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IssuedCert.Marshal(b, m, deterministic)
}
func (m *IssuedCert) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IssuedCert.Merge(m, src)
}
func (m *IssuedCert) XXX_Size() int {
	return xxx_messageInfo_IssuedCert.Size(m)
}
func (m *IssuedCert) XXX_DiscardUnknown() {
	xxx_messageInfo_IssuedCert.DiscardUnknown(m)
}

var xxx_messageInfo_IssuedCert proto.InternalMessageInfo

func (m *IssuedCert) GetSerialNumber() string {
	if m != nil {
		return m.SerialNumber
	}
	return ""
}

func (m *IssuedCert) GetCertPEM() string {
	if m != nil {
		return m.CertPEM
	}
	return ""
}

func (m *IssuedCert) GetPrivateKeyPEM() string {
	if m != nil {
		return m.PrivateKeyPEM
	}
	return ""
}

func (m *IssuedCert) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *IssuedCert) GetServiceURI() string {
	if m != nil {
		return m.ServiceURI
	}
	return ""
}

func (m *IssuedCert) GetAgent() string {
	if m != nil {
		return m.Agent
	}
	return ""
}

func (m *IssuedCert) GetAgentURI() string {
	if m != nil {
		return m.AgentURI
	}
	return ""
}

func (m *IssuedCert) GetValidAfter() *timestamp.Timestamp {
	if m != nil {
		return m.ValidAfter
	}
	return nil
}

func (m *IssuedCert) GetValidBefore() *timestamp.Timestamp {
	if m != nil {
		return m.ValidBefore
	}
	return nil
}

func (m *IssuedCert) GetEnterpriseMeta() *pbcommon.EnterpriseMeta {
	if m != nil {
		return m.EnterpriseMeta
	}
	return nil
}

func (m *IssuedCert) GetRaftIndex() *pbcommon.RaftIndex {
	if m != nil {
		return m.RaftIndex
	}
	return nil
}

func init() {
	proto.RegisterType((*CARoots)(nil), "connect.CARoots")
	proto.RegisterType((*CARoot)(nil), "connect.CARoot")
	proto.RegisterType((*IssuedCert)(nil), "connect.IssuedCert")
}

func init() {
	proto.RegisterFile("proto/pbconnect/connect.proto", fileDescriptor_80627e709958eb04)
}

var fileDescriptor_80627e709958eb04 = []byte{
	// 632 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x55, 0xd7, 0xcf, 0xdc, 0x6e, 0x1d, 0x33, 0x68, 0xb2, 0x8a, 0x80, 0xa8, 0x02, 0x14, 0x09,
	0x68, 0xd0, 0x90, 0x10, 0x42, 0x68, 0x52, 0xb7, 0xee, 0x21, 0x9a, 0x56, 0x86, 0x37, 0x78, 0xe0,
	0x2d, 0x6d, 0x6f, 0x3b, 0x4b, 0x4d, 0x5c, 0x39, 0xce, 0xb4, 0xfe, 0x22, 0x7e, 0x0a, 0xbf, 0x0a,
	0x09, 0xd9, 0x4e, 0xda, 0xa4, 0x20, 0xf5, 0x29, 0xbe, 0xe7, 0x1e, 0x5f, 0xdf, 0xeb, 0x73, 0x62,
	0x78, 0xb6, 0x94, 0x42, 0x09, 0x7f, 0x39, 0x9e, 0x88, 0x38, 0xc6, 0x89, 0xf2, 0xb3, 0x6f, 0xdf,
	0xe0, 0xa4, 0x99, 0x85, 0xdd, 0x17, 0x73, 0x21, 0xe6, 0x0b, 0xf4, 0x0d, 0x3c, 0x4e, 0x67, 0xbe,
	0xe2, 0x11, 0x26, 0x2a, 0x8c, 0x96, 0x96, 0xd9, 0x7d, 0xba, 0x29, 0x14, 0x45, 0x22, 0xf6, 0xed,
	0xc7, 0x26, 0x7b, 0xbf, 0x2a, 0xd0, 0x3c, 0x1f, 0x30, 0x21, 0x54, 0x42, 0x7a, 0xb0, 0x3f, 0x98,
	0x28, 0x7e, 0x8f, 0x3a, 0x0c, 0x86, 0xb4, 0xe2, 0x56, 0x3c, 0x87, 0x95, 0x30, 0xe2, 0x42, 0xfb,
	0x56, 0xa6, 0x89, 0x1a, 0x8a, 0x28, 0xe4, 0x31, 0xdd, 0x33, 0x94, 0x22, 0x44, 0x5e, 0x41, 0xdd,
	0x94, 0xa3, 0x55, 0xb7, 0xea, 0xb5, 0x4f, 0x0e, 0xfb, 0x79, 0xdf, 0xf6, 0x18, 0x66, 0xb3, 0xc4,
	0x07, 0xe7, 0x5b, 0x8a, 0x72, 0x75, 0x85, 0x2a, 0xa4, 0x35, 0xb7, 0xe2, 0xb5, 0x4f, 0x8e, 0xfa,
	0x59, 0x6b, 0xeb, 0x04, 0xdb, 0x70, 0x7a, 0x7f, 0x6a, 0xd0, 0xb0, 0x25, 0x48, 0x07, 0xf6, 0xd6,
	0xed, 0xed, 0x05, 0x43, 0x42, 0xa0, 0x36, 0x0a, 0x23, 0xcc, 0xba, 0x31, 0x6b, 0x3d, 0xcc, 0x0d,
	0x4a, 0x1e, 0x2e, 0x46, 0x69, 0x34, 0x46, 0x49, 0xab, 0x6e, 0xc5, 0xab, 0xb1, 0x12, 0x66, 0x38,
	0x7c, 0x1e, 0xf3, 0x78, 0x7e, 0x89, 0xab, 0x60, 0x68, 0xda, 0x70, 0x58, 0x09, 0x23, 0xef, 0xe1,
	0xf1, 0xc5, 0x83, 0x42, 0x19, 0x87, 0x8b, 0xe2, 0xe0, 0x75, 0x43, 0xfd, 0x5f, 0x8a, 0x7c, 0x02,
	0x67, 0x24, 0xd4, 0x19, 0xce, 0x84, 0x44, 0xda, 0x30, 0x93, 0x75, 0xfb, 0x56, 0xa4, 0x7e, 0x2e,
	0x52, 0xff, 0x36, 0x17, 0x89, 0x6d, 0xc8, 0xe4, 0x23, 0xb4, 0x46, 0x42, 0x0d, 0x66, 0x0a, 0x25,
	0x6d, 0xee, 0xdc, 0xb8, 0xe6, 0x92, 0x2e, 0xb4, 0xf4, 0xbd, 0x9c, 0xa3, 0x54, 0xb4, 0x65, 0x1a,
	0x5b, 0xc7, 0xe4, 0x2d, 0x1c, 0x05, 0xb1, 0x42, 0x19, 0xe1, 0x94, 0x87, 0x0a, 0x35, 0x96, 0x50,
	0xc7, 0xad, 0x7a, 0x0e, 0xfb, 0x37, 0xa1, 0xe5, 0xcd, 0xa6, 0x37, 0xc5, 0xc0, 0xca, 0x5b, 0x80,
	0xc8, 0x73, 0x80, 0xcd, 0xfd, 0xd0, 0xb6, 0x21, 0x14, 0x10, 0x72, 0x0c, 0x0d, 0x6b, 0x18, 0xba,
	0xef, 0x56, 0xbc, 0x16, 0xcb, 0x22, 0x72, 0x0a, 0xfb, 0x4c, 0xa8, 0x50, 0xe1, 0xf4, 0x6b, 0xaa,
	0x06, 0x8a, 0x1e, 0xec, 0x9c, 0xaf, 0xc4, 0x27, 0xaf, 0xa1, 0x73, 0x2d, 0xf9, 0x7d, 0xa8, 0xf0,
	0x12, 0x57, 0xb7, 0xab, 0x25, 0xd2, 0x8e, 0x39, 0x7b, 0x0b, 0x2d, 0xf3, 0xce, 0xb8, 0x4a, 0xe8,
	0xa1, 0x5b, 0xf1, 0xea, 0x6c, 0x0b, 0xd5, 0xfe, 0x63, 0xe1, 0x4c, 0x05, 0xf1, 0x14, 0x1f, 0xe8,
	0xa3, 0xb2, 0xff, 0xd6, 0x09, 0xb6, 0xe1, 0xf4, 0x7e, 0x57, 0x01, 0x82, 0x24, 0x49, 0x71, 0x6a,
	0xee, 0x61, 0xdb, 0x5f, 0xd9, 0xcf, 0x52, 0xf2, 0x17, 0x85, 0xa6, 0xe6, 0x5e, 0x5f, 0x5c, 0x65,
	0xd6, 0xcc, 0x43, 0xf2, 0x12, 0x0e, 0x36, 0xfd, 0xe8, 0x7c, 0xd5, 0xe4, 0xcb, 0xa0, 0xde, 0x7f,
	0x83, 0xf2, 0x9e, 0x4f, 0x30, 0xb3, 0x66, 0x1e, 0x1a, 0x15, 0xec, 0xf2, 0x3b, 0x0b, 0x32, 0x33,
	0x16, 0x10, 0xf2, 0x04, 0xea, 0x83, 0x39, 0xc6, 0xca, 0xf8, 0xcf, 0x61, 0x36, 0xd0, 0x3e, 0x31,
	0x0b, 0xbd, 0xa7, 0x69, 0x7d, 0x92, 0xc7, 0xe4, 0x33, 0xc0, 0x8f, 0x70, 0xc1, 0xa7, 0xd6, 0x7d,
	0xad, 0x9d, 0xea, 0x14, 0xd8, 0xe4, 0x0b, 0xb4, 0x4d, 0x94, 0x79, 0xde, 0xd9, 0xb9, 0xb9, 0x48,
	0x27, 0xa7, 0xd0, 0xb9, 0xd0, 0x46, 0x5c, 0x4a, 0x9e, 0xa0, 0x79, 0x0e, 0xc0, 0x14, 0x38, 0xce,
	0xe5, 0x28, 0x67, 0xd9, 0x16, 0xbb, 0xac, 0x64, 0x7b, 0xb7, 0x92, 0x67, 0xef, 0x7e, 0xbe, 0x99,
	0x73, 0x75, 0x97, 0x8e, 0x35, 0xcb, 0xbf, 0x0b, 0x93, 0x3b, 0x3e, 0x11, 0x72, 0xa9, 0x1f, 0xd8,
	0x24, 0x5d, 0xf8, 0x5b, 0xef, 0xee, 0xb8, 0x61, 0x80, 0x0f, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff,
	0x77, 0x18, 0x20, 0xcd, 0x91, 0x05, 0x00, 0x00,
}
