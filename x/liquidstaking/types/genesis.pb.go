// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: canto/liquidstaking/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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
	Params                            Params                                  `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	Epoch                             Epoch                                   `protobuf:"bytes,2,opt,name=epoch,proto3" json:"epoch"`
	LastChunkId                       uint64                                  `protobuf:"varint,3,opt,name=last_chunk_id,json=lastChunkId,proto3" json:"last_chunk_id,omitempty"`
	LastInsuranceId                   uint64                                  `protobuf:"varint,4,opt,name=last_insurance_id,json=lastInsuranceId,proto3" json:"last_insurance_id,omitempty"`
	Chunks                            []*Chunk                                `protobuf:"bytes,5,rep,name=chunks,proto3" json:"chunks,omitempty"`
	Insurances                        []*Insurance                            `protobuf:"bytes,6,rep,name=insurances,proto3" json:"insurances,omitempty"`
	WithdrawingInsurances             []*WithdrawingInsurance                 `protobuf:"bytes,7,rep,name=withdrawing_insurances,json=withdrawingInsurances,proto3" json:"withdrawing_insurances,omitempty"`
	LiquidUnstakeUnbondingDelegations []*LiquidUnstakeUnbondingDelegationInfo `protobuf:"bytes,8,rep,name=liquid_unstake_unbonding_delegations,json=liquidUnstakeUnbondingDelegations,proto3" json:"liquid_unstake_unbonding_delegations,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8c4913de4c15036, []int{0}
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

func init() {
	proto.RegisterType((*GenesisState)(nil), "canto.liquidstaking.v1.GenesisState")
}

func init() {
	proto.RegisterFile("canto/liquidstaking/v1/genesis.proto", fileDescriptor_b8c4913de4c15036)
}

var fileDescriptor_b8c4913de4c15036 = []byte{
	// 465 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xdf, 0x6a, 0xd4, 0x40,
	0x14, 0x87, 0x13, 0xbb, 0x5d, 0x65, 0xaa, 0x88, 0x83, 0x96, 0x58, 0x70, 0xba, 0x2d, 0xbd, 0x28,
	0x45, 0x33, 0x54, 0x51, 0x11, 0x7a, 0x63, 0xab, 0xc8, 0x82, 0x88, 0x44, 0x16, 0xc1, 0x9b, 0x30,
	0x9b, 0x4c, 0x67, 0x87, 0x4d, 0x66, 0x62, 0x66, 0xb2, 0xab, 0x0f, 0x21, 0xf8, 0x54, 0xb2, 0x97,
	0xbd, 0xf4, 0x4a, 0x74, 0xf7, 0x45, 0x64, 0x4e, 0xd2, 0xed, 0x1f, 0x1a, 0xbd, 0xcb, 0x9c, 0xf3,
	0xfd, 0xbe, 0x73, 0x02, 0x07, 0xed, 0x24, 0x4c, 0x59, 0x4d, 0x33, 0xf9, 0xb9, 0x92, 0xa9, 0xb1,
	0x6c, 0x2c, 0x95, 0xa0, 0x93, 0x7d, 0x2a, 0xb8, 0xe2, 0x46, 0x9a, 0xb0, 0x28, 0xb5, 0xd5, 0x78,
	0x1d, 0xa8, 0xf0, 0x02, 0x15, 0x4e, 0xf6, 0x37, 0xee, 0x0a, 0x2d, 0x34, 0x20, 0xd4, 0x7d, 0xd5,
	0xf4, 0x06, 0x11, 0x5a, 0x8b, 0x8c, 0x53, 0x78, 0x0d, 0xab, 0x63, 0x9a, 0x56, 0x25, 0xb3, 0x52,
	0xab, 0xa6, 0xbf, 0x79, 0xb9, 0x6f, 0x65, 0xce, 0x8d, 0x65, 0x79, 0xd1, 0x00, 0xf7, 0x13, 0x6d,
	0x72, 0x6d, 0xe2, 0xda, 0x5c, 0x3f, 0x9a, 0xd6, 0x5e, 0xcb, 0xbe, 0x17, 0x57, 0x03, 0x76, 0xfb,
	0x47, 0x07, 0xdd, 0x7c, 0x53, 0xff, 0xc7, 0x07, 0xcb, 0x2c, 0xc7, 0x07, 0xa8, 0x5b, 0xb0, 0x92,
	0xe5, 0x26, 0xf0, 0x7b, 0xfe, 0xee, 0xda, 0x63, 0x12, 0x5e, 0xfd, 0x5f, 0xe1, 0x7b, 0xa0, 0x0e,
	0x3b, 0xb3, 0x5f, 0x9b, 0x5e, 0xd4, 0x64, 0xf0, 0x0b, 0xb4, 0xca, 0x0b, 0x9d, 0x8c, 0x82, 0x6b,
	0x10, 0x7e, 0xd0, 0x16, 0x7e, 0xed, 0xa0, 0x26, 0x5b, 0x27, 0xf0, 0x36, 0xba, 0x95, 0x31, 0x63,
	0xe3, 0x64, 0x54, 0xa9, 0x71, 0x2c, 0xd3, 0x60, 0xa5, 0xe7, 0xef, 0x76, 0xa2, 0x35, 0x57, 0x3c,
	0x72, 0xb5, 0x7e, 0x8a, 0xf7, 0xd0, 0x1d, 0x60, 0xa4, 0x32, 0x55, 0xc9, 0x54, 0xc2, 0x1d, 0xd7,
	0x01, 0xee, 0xb6, 0x6b, 0xf4, 0x4f, 0xeb, 0xfd, 0x14, 0x3f, 0x45, 0x5d, 0x50, 0x99, 0x60, 0xb5,
	0xb7, 0xf2, 0xaf, 0x5d, 0x40, 0x1e, 0x35, 0x30, 0x7e, 0x89, 0xd0, 0xd2, 0x6e, 0x82, 0x2e, 0x44,
	0xb7, 0xda, 0xa2, 0xcb, 0x79, 0xd1, 0xb9, 0x10, 0x4e, 0xd0, 0xfa, 0x54, 0xda, 0x51, 0x5a, 0xb2,
	0xa9, 0x54, 0x22, 0x3e, 0xa7, 0xbb, 0x0e, 0xba, 0x87, 0x6d, 0xba, 0x8f, 0x67, 0xa9, 0x33, 0xf3,
	0xbd, 0xe9, 0x15, 0x55, 0x83, 0xbf, 0xf9, 0x68, 0xa7, 0x16, 0xc4, 0x95, 0x72, 0x0e, 0x1e, 0x57,
	0x6a, 0xa8, 0x55, 0xea, 0x46, 0xa6, 0x3c, 0xe3, 0x02, 0xce, 0xc9, 0x04, 0x37, 0x60, 0xe6, 0x41,
	0xdb, 0xcc, 0xb7, 0x50, 0x18, 0xd4, 0x8a, 0xc1, 0xa9, 0xe1, 0xd5, 0x52, 0xd0, 0x57, 0xc7, 0x3a,
	0xda, 0xca, 0xfe, 0x43, 0x99, 0xc3, 0xc1, 0xec, 0x0f, 0xf1, 0x66, 0x73, 0xe2, 0x9f, 0xcc, 0x89,
	0xff, 0x7b, 0x4e, 0xfc, 0xef, 0x0b, 0xe2, 0x9d, 0x2c, 0x88, 0xf7, 0x73, 0x41, 0xbc, 0x4f, 0xcf,
	0x85, 0xb4, 0xa3, 0x6a, 0x18, 0x26, 0x3a, 0xa7, 0x47, 0x6e, 0x91, 0x47, 0xef, 0xb8, 0x9d, 0xea,
	0x72, 0x5c, 0xbf, 0xe8, 0xe4, 0x19, 0xfd, 0x72, 0xe9, 0x60, 0xed, 0xd7, 0x82, 0x9b, 0x61, 0x17,
	0xce, 0xf4, 0xc9, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x46, 0xc2, 0x4a, 0xdc, 0x84, 0x03, 0x00,
	0x00,
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
	if len(m.LiquidUnstakeUnbondingDelegations) > 0 {
		for iNdEx := len(m.LiquidUnstakeUnbondingDelegations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.LiquidUnstakeUnbondingDelegations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.WithdrawingInsurances) > 0 {
		for iNdEx := len(m.WithdrawingInsurances) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.WithdrawingInsurances[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.Insurances) > 0 {
		for iNdEx := len(m.Insurances) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Insurances[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Chunks) > 0 {
		for iNdEx := len(m.Chunks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Chunks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if m.LastInsuranceId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.LastInsuranceId))
		i--
		dAtA[i] = 0x20
	}
	if m.LastChunkId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.LastChunkId))
		i--
		dAtA[i] = 0x18
	}
	{
		size, err := m.Epoch.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
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
	l = m.Epoch.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.LastChunkId != 0 {
		n += 1 + sovGenesis(uint64(m.LastChunkId))
	}
	if m.LastInsuranceId != 0 {
		n += 1 + sovGenesis(uint64(m.LastInsuranceId))
	}
	if len(m.Chunks) > 0 {
		for _, e := range m.Chunks {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Insurances) > 0 {
		for _, e := range m.Insurances {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.WithdrawingInsurances) > 0 {
		for _, e := range m.WithdrawingInsurances {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.LiquidUnstakeUnbondingDelegations) > 0 {
		for _, e := range m.LiquidUnstakeUnbondingDelegations {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
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
				return fmt.Errorf("proto: wrong wireType = %d for field Epoch", wireType)
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
			if err := m.Epoch.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastChunkId", wireType)
			}
			m.LastChunkId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastChunkId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastInsuranceId", wireType)
			}
			m.LastInsuranceId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastInsuranceId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chunks", wireType)
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
			m.Chunks = append(m.Chunks, &Chunk{})
			if err := m.Chunks[len(m.Chunks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Insurances", wireType)
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
			m.Insurances = append(m.Insurances, &Insurance{})
			if err := m.Insurances[len(m.Insurances)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WithdrawingInsurances", wireType)
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
			m.WithdrawingInsurances = append(m.WithdrawingInsurances, &WithdrawingInsurance{})
			if err := m.WithdrawingInsurances[len(m.WithdrawingInsurances)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiquidUnstakeUnbondingDelegations", wireType)
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
			m.LiquidUnstakeUnbondingDelegations = append(m.LiquidUnstakeUnbondingDelegations, &LiquidUnstakeUnbondingDelegationInfo{})
			if err := m.LiquidUnstakeUnbondingDelegations[len(m.LiquidUnstakeUnbondingDelegations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
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
