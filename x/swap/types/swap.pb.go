// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: swap/swap.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Swap struct {
	TransferType TransferType                                  `protobuf:"varint,1,opt,name=TransferType,proto3,enum=swap.TransferType" json:"transfer_type"`
	HashedSecret Hash                                          `protobuf:"bytes,2,opt,name=HashedSecret,proto3,customtype=Hash" json:"hashed_secret"`
	From         github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,3,opt,name=From,proto3,customtype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"from"`
	Recipient    string                                        `protobuf:"bytes,4,opt,name=Recipient,proto3" json:"recipient"`
	Amount       github_com_cosmos_cosmos_sdk_types.Coins      `protobuf:"bytes,5,rep,name=Amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	Timestamp    uint64                                        `protobuf:"varint,6,opt,name=Timestamp,proto3" json:"timestamp"`
	Redeemed     bool                                          `protobuf:"varint,7,opt,name=Redeemed,proto3" json:"redeemed"`
	Refunded     bool                                          `protobuf:"varint,8,opt,name=Refunded,proto3" json:"refunded"`
}

func (m *Swap) Reset()         { *m = Swap{} }
func (m *Swap) String() string { return proto.CompactTextString(m) }
func (*Swap) ProtoMessage()    {}
func (*Swap) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4906e0bf1273377, []int{0}
}
func (m *Swap) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Swap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Swap.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Swap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Swap.Merge(m, src)
}
func (m *Swap) XXX_Size() int {
	return m.Size()
}
func (m *Swap) XXX_DiscardUnknown() {
	xxx_messageInfo_Swap.DiscardUnknown(m)
}

var xxx_messageInfo_Swap proto.InternalMessageInfo

func (m *Swap) GetTransferType() TransferType {
	if m != nil {
		return m.TransferType
	}
	return TransferType_TransferTypeOut
}

func (m *Swap) GetRecipient() string {
	if m != nil {
		return m.Recipient
	}
	return ""
}

func (m *Swap) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func (m *Swap) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Swap) GetRedeemed() bool {
	if m != nil {
		return m.Redeemed
	}
	return false
}

func (m *Swap) GetRefunded() bool {
	if m != nil {
		return m.Refunded
	}
	return false
}

type Params struct {
	LockedTimeOut time.Duration `protobuf:"bytes,1,opt,name=LockedTimeOut,proto3,customtype=time.Duration" json:"locked_time_out"`
	LockedTimeIn  time.Duration `protobuf:"bytes,2,opt,name=LockedTimeIn,proto3,customtype=time.Duration" json:"locked_time_in"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4906e0bf1273377, []int{1}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Swap)(nil), "swap.Swap")
	proto.RegisterType((*Params)(nil), "swap.Params")
}

func init() { proto.RegisterFile("swap/swap.proto", fileDescriptor_b4906e0bf1273377) }

var fileDescriptor_b4906e0bf1273377 = []byte{
	// 498 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x4f, 0x8b, 0xd3, 0x40,
	0x14, 0x6f, 0x6c, 0xac, 0xed, 0x90, 0xee, 0xe2, 0xe0, 0x21, 0x2c, 0x98, 0x84, 0xbd, 0x18, 0x90,
	0x26, 0xa0, 0xe0, 0xcd, 0x43, 0xa3, 0xc8, 0x2e, 0xca, 0xaa, 0xb3, 0x3d, 0x79, 0x29, 0xd3, 0xcc,
	0xb4, 0x0d, 0xdd, 0xc9, 0x84, 0x99, 0x09, 0xbb, 0xfb, 0x2d, 0xfc, 0x12, 0x5e, 0xfc, 0x24, 0x7b,
	0xdc, 0xa3, 0x78, 0x88, 0xd2, 0xe2, 0x25, 0x9f, 0x42, 0x66, 0x92, 0xdd, 0xb6, 0x07, 0xc1, 0x4b,
	0xf2, 0xf2, 0xfb, 0xf7, 0x1e, 0x79, 0x0f, 0x1c, 0xca, 0x4b, 0x5c, 0xc4, 0xfa, 0x11, 0x15, 0x82,
	0x2b, 0x0e, 0x6d, 0x5d, 0x1f, 0x3d, 0x59, 0xf0, 0x05, 0x37, 0x40, 0xac, 0xab, 0x86, 0x3b, 0x1a,
	0x1a, 0xb1, 0xba, 0x6a, 0x3e, 0x8f, 0xff, 0x74, 0x81, 0x7d, 0x7e, 0x89, 0x0b, 0x78, 0x0a, 0x9c,
	0x89, 0xc0, 0xb9, 0x9c, 0x53, 0x31, 0xb9, 0x2e, 0xa8, 0x6b, 0x05, 0x56, 0x78, 0xf0, 0x02, 0x46,
	0x26, 0x76, 0x97, 0x49, 0x1e, 0xd7, 0x95, 0x3f, 0x54, 0x2d, 0x32, 0x55, 0xd7, 0x05, 0x45, 0x7b,
	0x56, 0x38, 0x06, 0xce, 0x09, 0x96, 0x4b, 0x4a, 0xce, 0x69, 0x2a, 0xa8, 0x72, 0x1f, 0x04, 0x56,
	0xe8, 0x24, 0x4f, 0x6f, 0x2a, 0xbf, 0xf3, 0xb3, 0xf2, 0x6d, 0xcd, 0xe9, 0x88, 0xa5, 0xd1, 0x4c,
	0xa5, 0x11, 0xa1, 0x3d, 0x0b, 0xfc, 0x0c, 0xec, 0x77, 0x82, 0x33, 0xb7, 0x6b, 0xac, 0xaf, 0x5b,
	0xeb, 0x68, 0x91, 0xa9, 0x65, 0x39, 0x8b, 0x52, 0xce, 0xe2, 0x94, 0x4b, 0xc6, 0x65, 0xfb, 0x1a,
	0x49, 0xb2, 0x8a, 0xf5, 0x24, 0x32, 0x1a, 0xa7, 0xe9, 0x98, 0x10, 0x41, 0xa5, 0xac, 0x2b, 0xdf,
	0x9e, 0x0b, 0xce, 0x90, 0x89, 0x82, 0xcf, 0xc1, 0x00, 0xd1, 0x34, 0x2b, 0x32, 0x9a, 0x2b, 0xd7,
	0x0e, 0xac, 0x70, 0x90, 0x0c, 0xeb, 0xca, 0x1f, 0x88, 0x3b, 0x10, 0x6d, 0x79, 0x78, 0x06, 0x7a,
	0x63, 0xc6, 0xcb, 0x5c, 0xb9, 0x0f, 0x83, 0x6e, 0xe8, 0x24, 0xaf, 0xea, 0xca, 0xef, 0x61, 0x83,
	0x7c, 0xff, 0xe5, 0x87, 0xff, 0x31, 0xc7, 0x1b, 0x9e, 0xe5, 0x12, 0xb5, 0x29, 0xba, 0xf9, 0x24,
	0x63, 0x54, 0x2a, 0xcc, 0x0a, 0xb7, 0x17, 0x58, 0xa1, 0xdd, 0x34, 0x57, 0x77, 0x20, 0xda, 0xf2,
	0x30, 0x04, 0x7d, 0x44, 0x09, 0xa5, 0x8c, 0x12, 0xf7, 0x51, 0x60, 0x85, 0xfd, 0xc4, 0xa9, 0x2b,
	0xbf, 0x2f, 0x5a, 0x0c, 0xdd, 0xb3, 0x8d, 0x72, 0x5e, 0xe6, 0x84, 0x12, 0xb7, 0xbf, 0xab, 0x6c,
	0x30, 0x74, 0xcf, 0x1e, 0x7f, 0xb3, 0x40, 0xef, 0x13, 0x16, 0x98, 0x49, 0x78, 0x06, 0x86, 0x1f,
	0x78, 0xba, 0xa2, 0x44, 0x77, 0xfc, 0x58, 0x2a, 0xb3, 0x6a, 0x27, 0x09, 0xdb, 0x9f, 0x3c, 0xd4,
	0x33, 0x45, 0x6f, 0x4b, 0x81, 0x55, 0xc6, 0xf3, 0xba, 0xf2, 0x0f, 0x2f, 0x8c, 0x7a, 0xaa, 0xf1,
	0x29, 0x2f, 0x15, 0xda, 0xb7, 0xc3, 0xf7, 0xc0, 0xd9, 0x02, 0xa7, 0x79, 0xbb, 0xee, 0x67, 0xff,
	0x8a, 0x3b, 0xd8, 0x8d, 0xcb, 0x72, 0xb4, 0x67, 0x4e, 0x4e, 0x6e, 0xd6, 0x9e, 0x75, 0xbb, 0xf6,
	0xac, 0xdf, 0x6b, 0xcf, 0xfa, 0xba, 0xf1, 0x3a, 0xb7, 0x1b, 0xaf, 0xf3, 0x63, 0xe3, 0x75, 0xbe,
	0x44, 0xb3, 0x4c, 0xcd, 0xca, 0x74, 0x45, 0x55, 0xc4, 0xc5, 0x22, 0x26, 0x34, 0xcd, 0x18, 0xbe,
	0x50, 0x14, 0xb3, 0x78, 0xc1, 0x47, 0x39, 0x27, 0x34, 0xbe, 0x8a, 0x9b, 0xe3, 0xd6, 0x0b, 0x98,
	0xf5, 0xcc, 0x81, 0xbf, 0xfc, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x58, 0x89, 0x3d, 0x82, 0x1e, 0x03,
	0x00, 0x00,
}

func (m *Swap) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Swap) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Swap) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Refunded {
		i--
		if m.Refunded {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if m.Redeemed {
		i--
		if m.Redeemed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.Timestamp != 0 {
		i = encodeVarintSwap(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x30
	}
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Amount[iNdEx])
			copy(dAtA[i:], m.Amount[iNdEx])
			i = encodeVarintSwap(dAtA, i, uint64(len(m.Amount[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintSwap(dAtA, i, uint64(len(m.Recipient)))
		i--
		dAtA[i] = 0x22
	}
	{
		size := m.From.Size()
		i -= size
		if _, err := m.From.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSwap(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.HashedSecret.Size()
		i -= size
		if _, err := m.HashedSecret.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSwap(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.TransferType != 0 {
		i = encodeVarintSwap(dAtA, i, uint64(m.TransferType))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.LockedTimeIn.Size()
		i -= size
		if _, err := m.LockedTimeIn.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSwap(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.LockedTimeOut.Size()
		i -= size
		if _, err := m.LockedTimeOut.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintSwap(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintSwap(dAtA []byte, offset int, v uint64) int {
	offset -= sovSwap(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Swap) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.TransferType != 0 {
		n += 1 + sovSwap(uint64(m.TransferType))
	}
	l = m.HashedSecret.Size()
	n += 1 + l + sovSwap(uint64(l))
	l = m.From.Size()
	n += 1 + l + sovSwap(uint64(l))
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovSwap(uint64(l))
	}
	if len(m.Amount) > 0 {
		for _, b := range m.Amount {
			l = len(b)
			n += 1 + l + sovSwap(uint64(l))
		}
	}
	if m.Timestamp != 0 {
		n += 1 + sovSwap(uint64(m.Timestamp))
	}
	if m.Redeemed {
		n += 2
	}
	if m.Refunded {
		n += 2
	}
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LockedTimeOut.Size()
	n += 1 + l + sovSwap(uint64(l))
	l = m.LockedTimeIn.Size()
	n += 1 + l + sovSwap(uint64(l))
	return n
}

func sovSwap(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSwap(x uint64) (n int) {
	return sovSwap(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Swap) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSwap
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
			return fmt.Errorf("proto: Swap: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Swap: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferType", wireType)
			}
			m.TransferType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TransferType |= TransferType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HashedSecret", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSwap
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.HashedSecret.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSwap
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.From.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
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
				return ErrInvalidLengthSwap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Recipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSwap
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, make([]byte, postIndex-iNdEx))
			copy(m.Amount[len(m.Amount)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Redeemed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Redeemed = bool(v != 0)
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Refunded", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Refunded = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipSwap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSwap
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
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSwap
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockedTimeOut", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSwap
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LockedTimeOut.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockedTimeIn", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSwap
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSwap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LockedTimeIn.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSwap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSwap
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
func skipSwap(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSwap
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
					return 0, ErrIntOverflowSwap
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
					return 0, ErrIntOverflowSwap
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
				return 0, ErrInvalidLengthSwap
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSwap
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSwap
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSwap        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSwap          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSwap = fmt.Errorf("proto: unexpected end of group")
)