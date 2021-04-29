// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/config/filter/listener/original_dst/v2/original_dst.proto

package envoy_config_filter_listener_original_dst_v2

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	proto "github.com/golang/protobuf/proto"
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

type OriginalDst struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OriginalDst) Reset()         { *m = OriginalDst{} }
func (m *OriginalDst) String() string { return proto.CompactTextString(m) }
func (*OriginalDst) ProtoMessage()    {}
func (*OriginalDst) Descriptor() ([]byte, []int) {
	return fileDescriptor_1b6ce05288da843c, []int{0}
}

func (m *OriginalDst) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OriginalDst.Unmarshal(m, b)
}
func (m *OriginalDst) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OriginalDst.Marshal(b, m, deterministic)
}
func (m *OriginalDst) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OriginalDst.Merge(m, src)
}
func (m *OriginalDst) XXX_Size() int {
	return xxx_messageInfo_OriginalDst.Size(m)
}
func (m *OriginalDst) XXX_DiscardUnknown() {
	xxx_messageInfo_OriginalDst.DiscardUnknown(m)
}

var xxx_messageInfo_OriginalDst proto.InternalMessageInfo

func init() {
	proto.RegisterType((*OriginalDst)(nil), "envoy.config.filter.listener.original_dst.v2.OriginalDst")
}

func init() {
	proto.RegisterFile("envoy/config/filter/listener/original_dst/v2/original_dst.proto", fileDescriptor_1b6ce05288da843c)
}

var fileDescriptor_1b6ce05288da843c = []byte{
	// 214 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xb2, 0x4f, 0xcd, 0x2b, 0xcb,
	0xaf, 0xd4, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x4f, 0xcb, 0xcc, 0x29, 0x49, 0x2d, 0xd2,
	0xcf, 0xc9, 0x2c, 0x2e, 0x49, 0xcd, 0x4b, 0x2d, 0xd2, 0xcf, 0x2f, 0xca, 0x4c, 0xcf, 0xcc, 0x4b,
	0xcc, 0x89, 0x4f, 0x29, 0x2e, 0xd1, 0x2f, 0x33, 0x42, 0xe1, 0xeb, 0x15, 0x14, 0xe5, 0x97, 0xe4,
	0x0b, 0xe9, 0x80, 0x0d, 0xd0, 0x83, 0x18, 0xa0, 0x07, 0x31, 0x40, 0x0f, 0x66, 0x80, 0x1e, 0x8a,
	0x86, 0x32, 0x23, 0x29, 0xb9, 0xd2, 0x94, 0x82, 0x44, 0xfd, 0xc4, 0xbc, 0xbc, 0xfc, 0x92, 0xc4,
	0x92, 0xcc, 0xfc, 0xbc, 0x62, 0xfd, 0xdc, 0xcc, 0xf4, 0xa2, 0xc4, 0x92, 0x54, 0x88, 0x69, 0x52,
	0xb2, 0x18, 0xf2, 0xc5, 0x25, 0x89, 0x25, 0xa5, 0xc5, 0x10, 0x69, 0x25, 0x5e, 0x2e, 0x6e, 0x7f,
	0xa8, 0x89, 0x2e, 0xc5, 0x25, 0x4e, 0x13, 0x19, 0x3f, 0xcd, 0xf8, 0xd7, 0xcf, 0x6a, 0x2c, 0x64,
	0x08, 0x71, 0x44, 0x6a, 0x45, 0x49, 0x6a, 0x5e, 0x31, 0x48, 0x1b, 0xd4, 0x21, 0xc5, 0xb8, 0x5c,
	0x62, 0xbc, 0xab, 0xe1, 0xc4, 0x45, 0x36, 0x26, 0x01, 0x46, 0x2e, 0xab, 0xcc, 0x7c, 0x3d, 0xb0,
	0xee, 0x82, 0xa2, 0xfc, 0x8a, 0x4a, 0x3d, 0x52, 0x7c, 0xe3, 0x24, 0x80, 0xe4, 0x98, 0x00, 0x90,
	0x03, 0x03, 0x18, 0x93, 0xd8, 0xc0, 0x2e, 0x35, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x80, 0xc9,
	0x36, 0x5e, 0x59, 0x01, 0x00, 0x00,
}
