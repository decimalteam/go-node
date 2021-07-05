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
	TransferType TransferType                             `protobuf:"varint,1,opt,name=TransferType,proto3,enum=swap.TransferType" json:"transfer_type"`
	HashedSecret Hash                                     `protobuf:"bytes,2,opt,name=HashedSecret,proto3,casttype=Hash" json:"hashed_secret"`
	From         string                                   `protobuf:"bytes,3,opt,name=From,proto3" json:"from"`
	Recipient    string                                   `protobuf:"bytes,4,opt,name=Recipient,proto3" json:"recipient"`
	Amount       github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,5,rep,name=Amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
	Timestamp    uint64                                   `protobuf:"varint,6,opt,name=Timestamp,proto3" json:"timestamp"`
	Redeemed     bool                                     `protobuf:"varint,7,opt,name=Redeemed,proto3" json:"redeemed"`
	Refunded     bool                                     `protobuf:"varint,8,opt,name=Refunded,proto3" json:"refunded"`
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

func (m *Swap) GetHashedSecret() Hash {
	if m != nil {
		return m.HashedSecret
	}
	return nil
}

func (m *Swap) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
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
	LockedTimeOut time.Duration `protobuf:"varint,1,opt,name=LockedTimeOut,proto3,casttype=time.Duration" json:"locked_time_out"`
	LockedTimeIn  time.Duration `protobuf:"varint,2,opt,name=LockedTimeIn,proto3,casttype=time.Duration" json:"locked_time_in"`
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

func (m *Params) GetLockedTimeOut() time.Duration {
	if m != nil {
		return m.LockedTimeOut
	}
	return 0
}

func (m *Params) GetLockedTimeIn() time.Duration {
	if m != nil {
		return m.LockedTimeIn
	}
	return 0
}

func init() {
	proto.RegisterType((*Swap)(nil), "swap.Swap")
	proto.RegisterType((*Params)(nil), "swap.Params")
}

func init() { proto.RegisterFile("swap/swap.proto", fileDescriptor_b4906e0bf1273377) }

var fileDescriptor_b4906e0bf1273377 = []byte{
	// 486 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0x67, 0x1a, 0x4a, 0x6a, 0x25, 0x9b, 0xb0, 0x38, 0x84, 0x09, 0x25, 0xd1, 0xe0, 0x10,
	0x09, 0x2d, 0x91, 0x40, 0xe2, 0xc6, 0x81, 0x80, 0x60, 0x43, 0x08, 0x90, 0xd7, 0x13, 0x97, 0xca,
	0x4d, 0xdc, 0x34, 0xea, 0x1c, 0x47, 0xb6, 0xa3, 0x6d, 0xdf, 0x82, 0x2f, 0xc0, 0x17, 0x40, 0x7c,
	0x10, 0x8e, 0x3b, 0x72, 0x0a, 0xa8, 0xbd, 0xe5, 0x23, 0x70, 0x42, 0x76, 0xba, 0xb5, 0x15, 0x97,
	0xe4, 0xcd, 0xef, 0x79, 0xde, 0xe7, 0x8d, 0xff, 0xc0, 0x03, 0x79, 0x41, 0xea, 0x44, 0x3f, 0xe2,
	0x5a, 0x70, 0xc5, 0x91, 0xa5, 0xeb, 0xc3, 0x07, 0x05, 0x2f, 0xb8, 0x01, 0x89, 0xae, 0x7a, 0xed,
	0xd0, 0x35, 0x66, 0x75, 0xd9, 0x7f, 0x1e, 0xfd, 0x18, 0x40, 0xeb, 0xec, 0x82, 0xd4, 0xe8, 0x14,
	0x3a, 0x63, 0x41, 0x2a, 0x39, 0xa3, 0x62, 0x7c, 0x55, 0x53, 0x0f, 0x84, 0x20, 0xda, 0x7f, 0x86,
	0x62, 0x13, 0xbb, 0xad, 0xa4, 0xf7, 0xbb, 0x36, 0x70, 0xd5, 0x9a, 0x4c, 0xd4, 0x55, 0x4d, 0xf1,
	0x4e, 0x2b, 0x7a, 0x09, 0x9d, 0x13, 0x22, 0xe7, 0x34, 0x3f, 0xa3, 0x99, 0xa0, 0xca, 0xbb, 0x13,
	0x82, 0xc8, 0x49, 0x1f, 0xea, 0xb6, 0xb9, 0xe1, 0x13, 0x69, 0x84, 0xbf, 0x6d, 0x60, 0x69, 0x23,
	0xde, 0xb1, 0xa3, 0x47, 0xd0, 0x7a, 0x2b, 0x38, 0xf3, 0x06, 0x21, 0x88, 0x46, 0xa9, 0xdd, 0xb5,
	0x81, 0x35, 0x13, 0x9c, 0x61, 0x43, 0xd1, 0x53, 0x38, 0xc2, 0x34, 0x2b, 0xeb, 0x92, 0x56, 0xca,
	0xb3, 0x8c, 0xc5, 0xed, 0xda, 0x60, 0x24, 0x6e, 0x20, 0xde, 0xe8, 0xe8, 0x23, 0x1c, 0xbe, 0x62,
	0xbc, 0xa9, 0x94, 0x77, 0x37, 0x1c, 0x44, 0x4e, 0xfa, 0xa2, 0x6b, 0x83, 0x21, 0x31, 0xe4, 0xfb,
	0xef, 0x20, 0x2a, 0x4a, 0x35, 0x6f, 0xa6, 0x71, 0xc6, 0x59, 0x92, 0x71, 0xc9, 0xb8, 0x5c, 0xbf,
	0x8e, 0x65, 0xbe, 0x48, 0xf4, 0xc2, 0x64, 0xfc, 0x9a, 0x97, 0x95, 0xc4, 0xeb, 0x14, 0x3d, 0x7c,
	0x5c, 0x32, 0x2a, 0x15, 0x61, 0xb5, 0x37, 0x0c, 0x41, 0x64, 0xf5, 0xc3, 0xd5, 0x0d, 0xc4, 0x1b,
	0x1d, 0x45, 0xd0, 0xc6, 0x34, 0xa7, 0x94, 0xd1, 0xdc, 0xbb, 0x17, 0x82, 0xc8, 0x4e, 0x9d, 0xae,
	0x0d, 0x6c, 0xb1, 0x66, 0xf8, 0x56, 0xed, 0x9d, 0xb3, 0xa6, 0xca, 0x69, 0xee, 0xd9, 0xdb, 0xce,
	0x9e, 0xe1, 0x5b, 0xf5, 0xe8, 0x1b, 0x80, 0xc3, 0xcf, 0x44, 0x10, 0x26, 0xd1, 0x7b, 0xe8, 0x7e,
	0xe0, 0xd9, 0x82, 0xe6, 0x7a, 0xe2, 0xa7, 0x46, 0x99, 0x13, 0x1b, 0xa4, 0x4f, 0xba, 0x36, 0x38,
	0x38, 0x37, 0xc2, 0x44, 0xff, 0xd6, 0x84, 0x37, 0x7a, 0xa3, 0x5d, 0x5d, 0xc7, 0x6f, 0x1a, 0x41,
	0x54, 0xc9, 0x2b, 0xbc, 0xdb, 0x8a, 0xde, 0x41, 0x67, 0x03, 0x4e, 0x2b, 0x73, 0x62, 0x83, 0xf4,
	0x71, 0xd7, 0x06, 0xfb, 0xdb, 0x51, 0x65, 0xf5, 0x7f, 0xd2, 0x4e, 0x63, 0x7a, 0xf2, 0x73, 0xe9,
	0x83, 0xeb, 0xa5, 0x0f, 0xfe, 0x2c, 0x7d, 0xf0, 0x75, 0xe5, 0xef, 0x5d, 0xaf, 0xfc, 0xbd, 0x5f,
	0x2b, 0x7f, 0xef, 0x4b, 0x3c, 0x2d, 0xd5, 0xb4, 0xc9, 0x16, 0x54, 0xc5, 0x5c, 0x14, 0x49, 0x4e,
	0xb3, 0x92, 0x91, 0x73, 0x45, 0x09, 0x4b, 0x0a, 0x7e, 0x5c, 0xf1, 0x9c, 0x26, 0x97, 0x49, 0x7f,
	0x37, 0xf5, 0xc6, 0x4f, 0x87, 0xe6, 0x7e, 0x3e, 0xff, 0x17, 0x00, 0x00, 0xff, 0xff, 0x58, 0xc2,
	0xf2, 0xf7, 0xdd, 0x02, 0x00, 0x00,
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
	if len(m.From) > 0 {
		i -= len(m.From)
		copy(dAtA[i:], m.From)
		i = encodeVarintSwap(dAtA, i, uint64(len(m.From)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.HashedSecret) > 0 {
		i -= len(m.HashedSecret)
		copy(dAtA[i:], m.HashedSecret)
		i = encodeVarintSwap(dAtA, i, uint64(len(m.HashedSecret)))
		i--
		dAtA[i] = 0x12
	}
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
	if m.LockedTimeIn != 0 {
		i = encodeVarintSwap(dAtA, i, uint64(m.LockedTimeIn))
		i--
		dAtA[i] = 0x10
	}
	if m.LockedTimeOut != 0 {
		i = encodeVarintSwap(dAtA, i, uint64(m.LockedTimeOut))
		i--
		dAtA[i] = 0x8
	}
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
	l = len(m.HashedSecret)
	if l > 0 {
		n += 1 + l + sovSwap(uint64(l))
	}
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovSwap(uint64(l))
	}
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
	if m.LockedTimeOut != 0 {
		n += 1 + sovSwap(uint64(m.LockedTimeOut))
	}
	if m.LockedTimeIn != 0 {
		n += 1 + sovSwap(uint64(m.LockedTimeIn))
	}
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
			m.HashedSecret = append(m.HashedSecret[:0], dAtA[iNdEx:postIndex]...)
			if m.HashedSecret == nil {
				m.HashedSecret = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
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
			m.From = string(dAtA[iNdEx:postIndex])
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockedTimeOut", wireType)
			}
			m.LockedTimeOut = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LockedTimeOut |= time.Duration(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LockedTimeIn", wireType)
			}
			m.LockedTimeIn = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSwap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LockedTimeIn |= time.Duration(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
