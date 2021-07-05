// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: validator/genesis.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

type GenesisState struct {
	Params               Params                                 `protobuf:"bytes,1,opt,name=Params,proto3" json:"params" yaml:"params"`
	LastTotalPower       github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=LastTotalPower,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"last_total_power" yaml:"last_total_power"`
	LastValidatorPowers  []LastValidatorPower                   `protobuf:"bytes,3,rep,name=LastValidatorPowers,proto3" json:"last_validator_powers" yaml:"last_validator_powers"`
	Validators           Validators                             `protobuf:"bytes,4,opt,name=Validators,proto3,customtype=Validators" json:"validators" yaml:"validators"`
	Delegations          Delegations                            `protobuf:"bytes,5,opt,name=Delegations,proto3,customtype=Delegations" json:"delegations" yaml:"delegations"`
	UnbondingDelegations []UnbondingDelegation                  `protobuf:"bytes,6,rep,name=UnbondingDelegations,proto3" json:"unbonding_delegations" yaml:"unbonding_delegations"`
	DelegationsNFT       DelegationsNFT                         `protobuf:"bytes,7,opt,name=DelegationsNFT,proto3,customtype=DelegationsNFT" json:"delegations_nft" yaml:"delegations_nft"`
	Exported             bool                                   `protobuf:"varint,8,opt,name=Exported,proto3" json:"exported" yaml:"exported"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_8143c6ee7ddaa59a, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

type LastValidatorPower struct {
	Address string `protobuf:"bytes,1,opt,name=Address,proto3" json:"address"`
	Power   int64  `protobuf:"varint,2,opt,name=Power,proto3" json:"Power,omitempty"`
}

func (m *LastValidatorPower) Reset()         { *m = LastValidatorPower{} }
func (m *LastValidatorPower) String() string { return proto.CompactTextString(m) }
func (*LastValidatorPower) ProtoMessage()    {}
func (*LastValidatorPower) Descriptor() ([]byte, []int) {
	return fileDescriptor_8143c6ee7ddaa59a, []int{1}
}
func (m *LastValidatorPower) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LastValidatorPower) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LastValidatorPower.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LastValidatorPower) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LastValidatorPower.Merge(m, src)
}
func (m *LastValidatorPower) XXX_Size() int {
	return m.Size()
}
func (m *LastValidatorPower) XXX_DiscardUnknown() {
	xxx_messageInfo_LastValidatorPower.DiscardUnknown(m)
}

var xxx_messageInfo_LastValidatorPower proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GenesisState)(nil), "validator.GenesisState")
	proto.RegisterType((*LastValidatorPower)(nil), "validator.LastValidatorPower")
}

func init() { proto.RegisterFile("validator/genesis.proto", fileDescriptor_8143c6ee7ddaa59a) }

var fileDescriptor_8143c6ee7ddaa59a = []byte{
	// 606 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x94, 0x4f, 0x6f, 0xd3, 0x3e,
	0x18, 0xc7, 0xe3, 0xdf, 0x7e, 0xeb, 0x3a, 0x77, 0x6c, 0xcc, 0x0c, 0x16, 0x26, 0x88, 0xab, 0x48,
	0xa0, 0x5e, 0xd6, 0x48, 0x9b, 0x90, 0xd0, 0x10, 0x07, 0xc2, 0x3f, 0x81, 0x00, 0x4d, 0x61, 0xe3,
	0x30, 0x09, 0x55, 0x6e, 0x63, 0x42, 0xb4, 0x24, 0xae, 0x62, 0x0f, 0xb6, 0x17, 0x80, 0xc4, 0x09,
	0xc1, 0x8d, 0x63, 0x5f, 0xce, 0x8e, 0x3b, 0x22, 0x0e, 0x16, 0x6a, 0x2f, 0x28, 0xc7, 0xbe, 0x02,
	0x94, 0x38, 0x4d, 0xb2, 0xb6, 0xa7, 0xc6, 0xdf, 0xcf, 0xe3, 0xef, 0xd7, 0xcf, 0xd3, 0x38, 0x70,
	0xf3, 0x13, 0x09, 0x7c, 0x97, 0x08, 0x16, 0x5b, 0x1e, 0x8d, 0x28, 0xf7, 0x79, 0xbb, 0x1f, 0x33,
	0xc1, 0xd0, 0x72, 0x01, 0xb6, 0x36, 0x3c, 0xe6, 0xb1, 0x4c, 0xb5, 0xd2, 0x27, 0x55, 0xb0, 0x75,
	0xb3, 0xdc, 0x59, 0x3c, 0xe5, 0x68, 0xab, 0x44, 0x2e, 0x0d, 0xa8, 0x47, 0x84, 0xcf, 0xa2, 0x9c,
	0x19, 0xf3, 0x58, 0x27, 0xfa, 0x20, 0x14, 0x37, 0xc7, 0x35, 0xb8, 0xf2, 0x5c, 0x9d, 0xe4, 0xad,
	0x20, 0x82, 0xa2, 0x97, 0xb0, 0xb6, 0x4f, 0x62, 0x12, 0x72, 0x1d, 0x34, 0x41, 0xab, 0xb1, 0xb3,
	0xde, 0x2e, 0xe3, 0x14, 0xb0, 0xf1, 0xb9, 0xc4, 0x5a, 0x22, 0x71, 0xad, 0x9f, 0xad, 0xc7, 0x12,
	0x5f, 0x39, 0x23, 0x61, 0xb0, 0x67, 0xaa, 0xb5, 0xe9, 0xe4, 0x0e, 0xe8, 0x0b, 0x80, 0xab, 0xaf,
	0x08, 0x17, 0x07, 0x4c, 0x90, 0x60, 0x9f, 0x7d, 0xa6, 0xb1, 0xfe, 0x5f, 0x13, 0xb4, 0x96, 0xed,
	0xf7, 0xa9, 0xc3, 0x6f, 0x89, 0xef, 0x7a, 0xbe, 0xf8, 0x78, 0xd2, 0x6d, 0xf7, 0x58, 0x68, 0xf5,
	0x18, 0x0f, 0x19, 0xcf, 0x7f, 0xb6, 0xb9, 0x7b, 0x6c, 0x89, 0xb3, 0x3e, 0xe5, 0xed, 0x17, 0x91,
	0x48, 0x24, 0xbe, 0x1a, 0x10, 0x2e, 0x3a, 0x22, 0x35, 0xea, 0xf4, 0x53, 0xa7, 0xb1, 0xc4, 0x9b,
	0x2a, 0x75, 0x9a, 0x98, 0xce, 0x54, 0x28, 0xfa, 0x06, 0xe0, 0xb5, 0x54, 0x7a, 0x37, 0xe9, 0x24,
	0x93, 0xb9, 0xbe, 0xd0, 0x5c, 0x68, 0x35, 0x76, 0x6e, 0x57, 0x3a, 0x9c, 0xad, 0xb2, 0x1f, 0xe6,
	0xdd, 0x5e, 0xcf, 0x72, 0x8a, 0x52, 0x95, 0x95, 0x36, 0x7f, 0xab, 0x72, 0x8c, 0x69, 0x6c, 0x3a,
	0xf3, 0x82, 0xd1, 0x21, 0x84, 0x85, 0xc4, 0xf5, 0xff, 0x9b, 0xa0, 0xb5, 0x62, 0xdf, 0xcb, 0x67,
	0x52, 0x21, 0x89, 0xc4, 0xb0, 0x70, 0x4c, 0xa3, 0xd6, 0x55, 0x54, 0xa9, 0x99, 0x4e, 0xa5, 0x1c,
	0x1d, 0xc1, 0xc6, 0x93, 0xe2, 0x4f, 0xe6, 0xfa, 0x62, 0xe6, 0x7b, 0x3f, 0xf7, 0xad, 0xa2, 0x44,
	0xe2, 0x46, 0xf9, 0x3a, 0xa4, 0xce, 0x48, 0x39, 0x57, 0x44, 0xd3, 0xa9, 0xee, 0x40, 0x3f, 0x00,
	0xdc, 0x38, 0x8c, 0xba, 0x2c, 0x72, 0xfd, 0xc8, 0xab, 0xa6, 0xd4, 0xb2, 0x21, 0x1a, 0x95, 0x21,
	0xce, 0x29, 0x2b, 0xa7, 0x78, 0x32, 0x81, 0x9d, 0xcb, 0x07, 0xc8, 0xa7, 0x38, 0x17, 0x9b, 0xce,
	0xdc, 0x68, 0x74, 0x0c, 0x57, 0x2b, 0xcb, 0x37, 0xcf, 0x0e, 0xf4, 0xa5, 0xac, 0xe5, 0xc7, 0x79,
	0xcb, 0x53, 0x34, 0x91, 0x78, 0xad, 0xe2, 0x9a, 0xde, 0x82, 0xb1, 0xc4, 0x37, 0x66, 0x3a, 0x4f,
	0x81, 0xe9, 0x4c, 0x6d, 0x46, 0x0f, 0x60, 0xfd, 0xe9, 0x69, 0x9f, 0xc5, 0x82, 0xba, 0x7a, 0xbd,
	0x09, 0x5a, 0x75, 0x1b, 0x27, 0x12, 0xd7, 0x69, 0xae, 0x8d, 0x25, 0x5e, 0x53, 0x4e, 0x13, 0xc5,
	0x74, 0x8a, 0x0d, 0x7b, 0x2b, 0x5f, 0x07, 0x58, 0xfb, 0x39, 0xc0, 0xe0, 0xef, 0x00, 0x6b, 0x66,
	0x07, 0xa2, 0xd9, 0xb7, 0x02, 0xdd, 0x81, 0x4b, 0x8f, 0x5c, 0x37, 0xa6, 0x5c, 0x5d, 0xbd, 0x65,
	0xbb, 0x91, 0x48, 0xbc, 0x44, 0x94, 0xe4, 0x4c, 0x18, 0xda, 0x80, 0x8b, 0xe5, 0x55, 0x5a, 0x70,
	0xd4, 0xe2, 0x72, 0x80, 0xfd, 0xfa, 0x7c, 0x68, 0x80, 0x8b, 0xa1, 0x01, 0xfe, 0x0c, 0x0d, 0xf0,
	0x7d, 0x64, 0x68, 0x17, 0x23, 0x43, 0xfb, 0x35, 0x32, 0xb4, 0xa3, 0xdd, 0xae, 0x2f, 0xba, 0x27,
	0xbd, 0x63, 0x2a, 0xda, 0x2c, 0xf6, 0x2c, 0x97, 0xf6, 0xfc, 0x90, 0x04, 0x82, 0x92, 0xd0, 0xf2,
	0xd8, 0x76, 0xc4, 0x5c, 0x6a, 0x9d, 0x96, 0x1f, 0x18, 0x75, 0x05, 0xbb, 0xb5, 0xec, 0x5b, 0xb1,
	0xfb, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x5c, 0x7f, 0xdc, 0x8d, 0xbe, 0x04, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Exported {
		i--
		if m.Exported {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	{
		size := m.DelegationsNFT.Size()
		i -= size
		if _, err := m.DelegationsNFT.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if len(m.UnbondingDelegations) > 0 {
		for iNdEx := len(m.UnbondingDelegations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.UnbondingDelegations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	{
		size := m.Delegations.Size()
		i -= size
		if _, err := m.Delegations.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.Validators.Size()
		i -= size
		if _, err := m.Validators.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.LastValidatorPowers) > 0 {
		for iNdEx := len(m.LastValidatorPowers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LastValidatorPowers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	{
		size := m.LastTotalPower.Size()
		i -= size
		if _, err := m.LastTotalPower.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *LastValidatorPower) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LastValidatorPower) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LastValidatorPower) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Power != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Power))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.LastTotalPower.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.LastValidatorPowers) > 0 {
		for _, e := range m.LastValidatorPowers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.Validators.Size()
	n += 1 + l + sovGenesis(uint64(l))
	l = m.Delegations.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.UnbondingDelegations) > 0 {
		for _, e := range m.UnbondingDelegations {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.DelegationsNFT.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.Exported {
		n += 2
	}
	return n
}

func (m *LastValidatorPower) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Power != 0 {
		n += 1 + sovGenesis(uint64(m.Power))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastTotalPower", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LastTotalPower.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastValidatorPowers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LastValidatorPowers = append(m.LastValidatorPowers, LastValidatorPower{})
			if err := m.LastValidatorPowers[len(m.LastValidatorPowers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validators", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Validators.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Delegations", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Delegations.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnbondingDelegations", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UnbondingDelegations = append(m.UnbondingDelegations, UnbondingDelegation{})
			if err := m.UnbondingDelegations[len(m.UnbondingDelegations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegationsNFT", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DelegationsNFT.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exported", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
			m.Exported = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *LastValidatorPower) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: LastValidatorPower: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LastValidatorPower: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Power", wireType)
			}
			m.Power = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Power |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
