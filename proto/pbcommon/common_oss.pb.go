// +build !consulent

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: proto/pbcommon/common_oss.proto

package pbcommon

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EnterpriseMeta struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnterpriseMeta) Reset()         { *m = EnterpriseMeta{} }
func (m *EnterpriseMeta) String() string { return proto.CompactTextString(m) }
func (*EnterpriseMeta) ProtoMessage()    {}
func (*EnterpriseMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_8f9d7cd54dd4e173, []int{0}
}
func (m *EnterpriseMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EnterpriseMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EnterpriseMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EnterpriseMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnterpriseMeta.Merge(m, src)
}
func (m *EnterpriseMeta) XXX_Size() int {
	return m.Size()
}
func (m *EnterpriseMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_EnterpriseMeta.DiscardUnknown(m)
}

var xxx_messageInfo_EnterpriseMeta proto.InternalMessageInfo

func init() {
	proto.RegisterType((*EnterpriseMeta)(nil), "common.EnterpriseMeta")
}

func init() { proto.RegisterFile("proto/pbcommon/common_oss.proto", fileDescriptor_8f9d7cd54dd4e173) }

var fileDescriptor_8f9d7cd54dd4e173 = []byte{
	// 127 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2f, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x2f, 0x48, 0x4a, 0xce, 0xcf, 0xcd, 0xcd, 0xcf, 0xd3, 0x87, 0x50, 0xf1, 0xf9, 0xc5,
	0xc5, 0x7a, 0x60, 0x19, 0x21, 0x36, 0x88, 0x88, 0x92, 0x00, 0x17, 0x9f, 0x6b, 0x5e, 0x49, 0x6a,
	0x51, 0x41, 0x51, 0x66, 0x71, 0xaa, 0x6f, 0x6a, 0x49, 0xa2, 0x93, 0xcd, 0x89, 0x47, 0x72, 0x8c,
	0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x56, 0x7a,
	0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x46, 0x62, 0x71, 0x46, 0x66, 0x72,
	0x7e, 0x51, 0x81, 0x7e, 0x72, 0x7e, 0x5e, 0x71, 0x69, 0x8e, 0x3e, 0xaa, 0x45, 0x49, 0x6c, 0x60,
	0xbe, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x2f, 0x9b, 0x6f, 0x83, 0x81, 0x00, 0x00, 0x00,
}

func (m *EnterpriseMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EnterpriseMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EnterpriseMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommonOss(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommonOss(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EnterpriseMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovCommonOss(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommonOss(x uint64) (n int) {
	return sovCommonOss(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EnterpriseMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommonOss
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: EnterpriseMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EnterpriseMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipCommonOss(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthCommonOss
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthCommonOss
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipCommonOss(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommonOss
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommonOss
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowCommonOss
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthCommonOss
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthCommonOss
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowCommonOss
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipCommonOss(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthCommonOss
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthCommonOss = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommonOss   = fmt.Errorf("proto: integer overflow")
)
