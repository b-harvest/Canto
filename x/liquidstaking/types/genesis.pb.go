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
	LiquidUnstakeUnbondingDelegations []*LiquidUnstakeUnbondingDelegationInfo `protobuf:"bytes,7,rep,name=liquid_unstake_unbonding_delegations,json=liquidUnstakeUnbondingDelegations,proto3" json:"liquid_unstake_unbonding_delegations,omitempty"`
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
	// 433 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0x41, 0x6b, 0x14, 0x31,
	0x14, 0xc7, 0x67, 0xdc, 0xed, 0x0a, 0xa9, 0x22, 0x06, 0x91, 0xb1, 0x60, 0xba, 0x2d, 0x3d, 0x94,
	0x82, 0x13, 0xaa, 0xa8, 0x08, 0xbd, 0xd8, 0x2a, 0x32, 0x20, 0x22, 0x23, 0x7b, 0xf1, 0x32, 0x64,
	0x66, 0xd2, 0x6c, 0xd8, 0x99, 0x64, 0x9c, 0x24, 0xab, 0x7e, 0x08, 0xc1, 0x8f, 0xb5, 0xc7, 0x1e,
	0x3d, 0x89, 0xdd, 0xfd, 0x22, 0x92, 0x64, 0xba, 0xd8, 0xc5, 0xd1, 0xdb, 0xbc, 0xf7, 0x7e, 0xff,
	0xdf, 0xbc, 0x90, 0x80, 0x83, 0x82, 0x08, 0x2d, 0x71, 0xc5, 0x3f, 0x19, 0x5e, 0x2a, 0x4d, 0x66,
	0x5c, 0x30, 0x3c, 0x3f, 0xc6, 0x8c, 0x0a, 0xaa, 0xb8, 0x8a, 0x9b, 0x56, 0x6a, 0x09, 0xef, 0x3b,
	0x2a, 0xbe, 0x46, 0xc5, 0xf3, 0xe3, 0x9d, 0x7b, 0x4c, 0x32, 0xe9, 0x10, 0x6c, 0xbf, 0x3c, 0xbd,
	0x83, 0x98, 0x94, 0xac, 0xa2, 0xd8, 0x55, 0xb9, 0x39, 0xc7, 0xa5, 0x69, 0x89, 0xe6, 0x52, 0x74,
	0xf3, 0xdd, 0xcd, 0xb9, 0xe6, 0x35, 0x55, 0x9a, 0xd4, 0x4d, 0x07, 0x3c, 0x28, 0xa4, 0xaa, 0xa5,
	0xca, 0xbc, 0xd9, 0x17, 0xdd, 0xe8, 0xa8, 0x67, 0xdf, 0xeb, 0xab, 0x39, 0x76, 0xff, 0x72, 0x00,
	0x6e, 0xbd, 0xf1, 0xe7, 0xf8, 0xa0, 0x89, 0xa6, 0xf0, 0x04, 0x8c, 0x1a, 0xd2, 0x92, 0x5a, 0x45,
	0xe1, 0x38, 0x3c, 0xdc, 0x7e, 0x8c, 0xe2, 0xbf, 0x9f, 0x2b, 0x7e, 0xef, 0xa8, 0xd3, 0xe1, 0xe2,
	0xe7, 0x6e, 0x90, 0x76, 0x19, 0xf8, 0x02, 0x6c, 0xd1, 0x46, 0x16, 0xd3, 0xe8, 0x86, 0x0b, 0x3f,
	0xec, 0x0b, 0xbf, 0xb6, 0x50, 0x97, 0xf5, 0x09, 0xb8, 0x0f, 0x6e, 0x57, 0x44, 0xe9, 0xac, 0x98,
	0x1a, 0x31, 0xcb, 0x78, 0x19, 0x0d, 0xc6, 0xe1, 0xe1, 0x30, 0xdd, 0xb6, 0xcd, 0x33, 0xdb, 0x4b,
	0x4a, 0x78, 0x04, 0xee, 0x3a, 0x86, 0x0b, 0x65, 0x5a, 0x22, 0x0a, 0x6a, 0xb9, 0xa1, 0xe3, 0xee,
	0xd8, 0x41, 0x72, 0xd5, 0x4f, 0x4a, 0xf8, 0x14, 0x8c, 0x9c, 0x4a, 0x45, 0x5b, 0xe3, 0xc1, 0xbf,
	0x76, 0x71, 0xf2, 0xb4, 0x83, 0xe1, 0x4b, 0x00, 0xd6, 0x76, 0x15, 0x8d, 0x5c, 0x74, 0xaf, 0x2f,
	0xba, 0xfe, 0x5f, 0xfa, 0x47, 0x08, 0x7e, 0x0b, 0xc1, 0x81, 0x47, 0x33, 0x23, 0x2c, 0x4d, 0x33,
	0x23, 0x72, 0x29, 0x4a, 0x2e, 0x58, 0x56, 0xd2, 0x8a, 0x32, 0x77, 0xd3, 0x2a, 0xba, 0xe9, 0xec,
	0x27, 0x7d, 0xf6, 0xb7, 0xae, 0x31, 0xf1, 0x8a, 0xc9, 0x95, 0xe1, 0xd5, 0x5a, 0x90, 0x88, 0x73,
	0x99, 0xee, 0x55, 0xff, 0xa1, 0xd4, 0xe9, 0x64, 0x71, 0x89, 0x82, 0xc5, 0x12, 0x85, 0x17, 0x4b,
	0x14, 0xfe, 0x5a, 0xa2, 0xf0, 0xfb, 0x0a, 0x05, 0x17, 0x2b, 0x14, 0xfc, 0x58, 0xa1, 0xe0, 0xe3,
	0x73, 0xc6, 0xf5, 0xd4, 0xe4, 0x71, 0x21, 0x6b, 0x7c, 0x66, 0x17, 0x79, 0xf4, 0x8e, 0xea, 0xcf,
	0xb2, 0x9d, 0xf9, 0x0a, 0xcf, 0x9f, 0xe1, 0x2f, 0x1b, 0x6f, 0x49, 0x7f, 0x6d, 0xa8, 0xca, 0x47,
	0xee, 0x05, 0x3d, 0xf9, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x82, 0x8f, 0xbb, 0x02, 0x1f, 0x03, 0x00,
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
