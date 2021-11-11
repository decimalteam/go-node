// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: validator/querier.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type QueryValidatorsParams struct {
	Page   int32  `protobuf:"varint,1,opt,name=Page,proto3" json:"Page,omitempty"`
	Limit  int32  `protobuf:"varint,2,opt,name=Limit,proto3" json:"Limit,omitempty"`
	Status string `protobuf:"bytes,3,opt,name=Status,proto3" json:"Status,omitempty"`
}

func (m *QueryValidatorsParams) Reset()         { *m = QueryValidatorsParams{} }
func (m *QueryValidatorsParams) String() string { return proto.CompactTextString(m) }
func (*QueryValidatorsParams) ProtoMessage()    {}
func (*QueryValidatorsParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_b72a1400752c4f69, []int{0}
}
func (m *QueryValidatorsParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryValidatorsParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryValidatorsParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryValidatorsParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryValidatorsParams.Merge(m, src)
}
func (m *QueryValidatorsParams) XXX_Size() int {
	return m.Size()
}
func (m *QueryValidatorsParams) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryValidatorsParams.DiscardUnknown(m)
}

var xxx_messageInfo_QueryValidatorsParams proto.InternalMessageInfo

func (m *QueryValidatorsParams) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *QueryValidatorsParams) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *QueryValidatorsParams) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type QueryValidatorParams struct {
	ValidatorAddr string `protobuf:"bytes,1,opt,name=ValidatorAddr,proto3" json:"validator_addr" yaml:"validator_addr"`
}

func (m *QueryValidatorParams) Reset()         { *m = QueryValidatorParams{} }
func (m *QueryValidatorParams) String() string { return proto.CompactTextString(m) }
func (*QueryValidatorParams) ProtoMessage()    {}
func (*QueryValidatorParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_b72a1400752c4f69, []int{1}
}
func (m *QueryValidatorParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryValidatorParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryValidatorParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryValidatorParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryValidatorParams.Merge(m, src)
}
func (m *QueryValidatorParams) XXX_Size() int {
	return m.Size()
}
func (m *QueryValidatorParams) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryValidatorParams.DiscardUnknown(m)
}

var xxx_messageInfo_QueryValidatorParams proto.InternalMessageInfo

func (m *QueryValidatorParams) GetValidatorAddr() string {
	if m != nil {
		return m.ValidatorAddr
	}
	return ""
}

func init() {
	proto.RegisterType((*QueryValidatorsParams)(nil), "validator.QueryValidatorsParams")
	proto.RegisterType((*QueryValidatorParams)(nil), "validator.QueryValidatorParams")
}

func init() { proto.RegisterFile("validator/querier.proto", fileDescriptor_b72a1400752c4f69) }

var fileDescriptor_b72a1400752c4f69 = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2f, 0x4b, 0xcc, 0xc9,
	0x4c, 0x49, 0x2c, 0xc9, 0x2f, 0xd2, 0x2f, 0x2c, 0x4d, 0x2d, 0xca, 0x4c, 0x2d, 0xd2, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x84, 0x4b, 0x48, 0x89, 0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0x45, 0xf5,
	0x41, 0x2c, 0x88, 0x02, 0xa5, 0x48, 0x2e, 0xd1, 0xc0, 0xd2, 0xd4, 0xa2, 0xca, 0x30, 0x98, 0xba,
	0xe2, 0x80, 0xc4, 0xa2, 0xc4, 0xdc, 0x62, 0x21, 0x21, 0x2e, 0x96, 0x80, 0xc4, 0xf4, 0x54, 0x09,
	0x46, 0x05, 0x46, 0x0d, 0xd6, 0x20, 0x30, 0x5b, 0x48, 0x84, 0x8b, 0xd5, 0x27, 0x33, 0x37, 0xb3,
	0x44, 0x82, 0x09, 0x2c, 0x08, 0xe1, 0x08, 0x89, 0x71, 0xb1, 0x05, 0x97, 0x24, 0x96, 0x94, 0x16,
	0x4b, 0x30, 0x2b, 0x30, 0x6a, 0x70, 0x06, 0x41, 0x79, 0x4a, 0x99, 0x5c, 0x22, 0xa8, 0x46, 0x43,
	0x4d, 0x0e, 0xe4, 0xe2, 0x85, 0x0b, 0x39, 0xa6, 0xa4, 0x14, 0x81, 0xad, 0xe0, 0x74, 0xd2, 0x7e,
	0x75, 0x4f, 0x9e, 0x0f, 0xee, 0xdc, 0xf8, 0xc4, 0x94, 0x94, 0xa2, 0x4f, 0xf7, 0xe4, 0x45, 0x2b,
	0x13, 0x73, 0x73, 0xac, 0x94, 0x50, 0xc5, 0x95, 0x82, 0x50, 0x4d, 0x70, 0xf2, 0x3d, 0xf1, 0x48,
	0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0,
	0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xe3, 0xa4, 0xcc, 0x92, 0xa4, 0xd2, 0xe4, 0xec,
	0xd4, 0x12, 0xbd, 0xfc, 0xa2, 0x74, 0xfd, 0x94, 0xd4, 0xe4, 0xcc, 0xdc, 0xc4, 0x9c, 0x92, 0xd4,
	0xc4, 0x5c, 0xfd, 0xf4, 0x7c, 0xdd, 0xbc, 0xfc, 0x94, 0x54, 0xfd, 0x0a, 0x7d, 0x44, 0xe8, 0x95,
	0x54, 0x16, 0xa4, 0x16, 0x27, 0xb1, 0x81, 0xc3, 0xc6, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x22,
	0x00, 0xb8, 0x99, 0x57, 0x01, 0x00, 0x00,
}

func (m *QueryValidatorsParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryValidatorsParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryValidatorsParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Status) > 0 {
		i -= len(m.Status)
		copy(dAtA[i:], m.Status)
		i = encodeVarintQuerier(dAtA, i, uint64(len(m.Status)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Limit != 0 {
		i = encodeVarintQuerier(dAtA, i, uint64(m.Limit))
		i--
		dAtA[i] = 0x10
	}
	if m.Page != 0 {
		i = encodeVarintQuerier(dAtA, i, uint64(m.Page))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryValidatorParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryValidatorParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryValidatorParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ValidatorAddr) > 0 {
		i -= len(m.ValidatorAddr)
		copy(dAtA[i:], m.ValidatorAddr)
		i = encodeVarintQuerier(dAtA, i, uint64(len(m.ValidatorAddr)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuerier(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuerier(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryValidatorsParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Page != 0 {
		n += 1 + sovQuerier(uint64(m.Page))
	}
	if m.Limit != 0 {
		n += 1 + sovQuerier(uint64(m.Limit))
	}
	l = len(m.Status)
	if l > 0 {
		n += 1 + l + sovQuerier(uint64(l))
	}
	return n
}

func (m *QueryValidatorParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ValidatorAddr)
	if l > 0 {
		n += 1 + l + sovQuerier(uint64(l))
	}
	return n
}

func sovQuerier(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuerier(x uint64) (n int) {
	return sovQuerier(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryValidatorsParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuerier
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
			return fmt.Errorf("proto: QueryValidatorsParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryValidatorsParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Page", wireType)
			}
			m.Page = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuerier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Page |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			m.Limit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuerier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Limit |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuerier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuerier
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuerier
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuerier(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuerier
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryValidatorParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuerier
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
			return fmt.Errorf("proto: QueryValidatorParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryValidatorParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuerier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuerier
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuerier
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ValidatorAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuerier(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuerier
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuerier(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuerier
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
					return 0, ErrIntOverflowQuerier
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuerier
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
				return 0, ErrInvalidLengthQuerier
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuerier
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuerier
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuerier        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuerier          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuerier = fmt.Errorf("proto: unexpected end of group")
)