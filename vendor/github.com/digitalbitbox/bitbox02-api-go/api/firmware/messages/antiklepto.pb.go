// Code generated by protoc-gen-go. DO NOT EDIT.
// source: antiklepto.proto

package messages

import (
	fmt "fmt"
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

type AntiKleptoHostNonceCommitment struct {
	Commitment           []byte   `protobuf:"bytes,1,opt,name=commitment,proto3" json:"commitment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AntiKleptoHostNonceCommitment) Reset()         { *m = AntiKleptoHostNonceCommitment{} }
func (m *AntiKleptoHostNonceCommitment) String() string { return proto.CompactTextString(m) }
func (*AntiKleptoHostNonceCommitment) ProtoMessage()    {}
func (*AntiKleptoHostNonceCommitment) Descriptor() ([]byte, []int) {
	return fileDescriptor_39dadd9b6cadc7c0, []int{0}
}

func (m *AntiKleptoHostNonceCommitment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AntiKleptoHostNonceCommitment.Unmarshal(m, b)
}
func (m *AntiKleptoHostNonceCommitment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AntiKleptoHostNonceCommitment.Marshal(b, m, deterministic)
}
func (m *AntiKleptoHostNonceCommitment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AntiKleptoHostNonceCommitment.Merge(m, src)
}
func (m *AntiKleptoHostNonceCommitment) XXX_Size() int {
	return xxx_messageInfo_AntiKleptoHostNonceCommitment.Size(m)
}
func (m *AntiKleptoHostNonceCommitment) XXX_DiscardUnknown() {
	xxx_messageInfo_AntiKleptoHostNonceCommitment.DiscardUnknown(m)
}

var xxx_messageInfo_AntiKleptoHostNonceCommitment proto.InternalMessageInfo

func (m *AntiKleptoHostNonceCommitment) GetCommitment() []byte {
	if m != nil {
		return m.Commitment
	}
	return nil
}

type AntiKleptoSignerCommitment struct {
	Commitment           []byte   `protobuf:"bytes,1,opt,name=commitment,proto3" json:"commitment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AntiKleptoSignerCommitment) Reset()         { *m = AntiKleptoSignerCommitment{} }
func (m *AntiKleptoSignerCommitment) String() string { return proto.CompactTextString(m) }
func (*AntiKleptoSignerCommitment) ProtoMessage()    {}
func (*AntiKleptoSignerCommitment) Descriptor() ([]byte, []int) {
	return fileDescriptor_39dadd9b6cadc7c0, []int{1}
}

func (m *AntiKleptoSignerCommitment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AntiKleptoSignerCommitment.Unmarshal(m, b)
}
func (m *AntiKleptoSignerCommitment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AntiKleptoSignerCommitment.Marshal(b, m, deterministic)
}
func (m *AntiKleptoSignerCommitment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AntiKleptoSignerCommitment.Merge(m, src)
}
func (m *AntiKleptoSignerCommitment) XXX_Size() int {
	return xxx_messageInfo_AntiKleptoSignerCommitment.Size(m)
}
func (m *AntiKleptoSignerCommitment) XXX_DiscardUnknown() {
	xxx_messageInfo_AntiKleptoSignerCommitment.DiscardUnknown(m)
}

var xxx_messageInfo_AntiKleptoSignerCommitment proto.InternalMessageInfo

func (m *AntiKleptoSignerCommitment) GetCommitment() []byte {
	if m != nil {
		return m.Commitment
	}
	return nil
}

type AntiKleptoSignatureRequest struct {
	HostNonce            []byte   `protobuf:"bytes,1,opt,name=host_nonce,json=hostNonce,proto3" json:"host_nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AntiKleptoSignatureRequest) Reset()         { *m = AntiKleptoSignatureRequest{} }
func (m *AntiKleptoSignatureRequest) String() string { return proto.CompactTextString(m) }
func (*AntiKleptoSignatureRequest) ProtoMessage()    {}
func (*AntiKleptoSignatureRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_39dadd9b6cadc7c0, []int{2}
}

func (m *AntiKleptoSignatureRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AntiKleptoSignatureRequest.Unmarshal(m, b)
}
func (m *AntiKleptoSignatureRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AntiKleptoSignatureRequest.Marshal(b, m, deterministic)
}
func (m *AntiKleptoSignatureRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AntiKleptoSignatureRequest.Merge(m, src)
}
func (m *AntiKleptoSignatureRequest) XXX_Size() int {
	return xxx_messageInfo_AntiKleptoSignatureRequest.Size(m)
}
func (m *AntiKleptoSignatureRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AntiKleptoSignatureRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AntiKleptoSignatureRequest proto.InternalMessageInfo

func (m *AntiKleptoSignatureRequest) GetHostNonce() []byte {
	if m != nil {
		return m.HostNonce
	}
	return nil
}

func init() {
	proto.RegisterType((*AntiKleptoHostNonceCommitment)(nil), "shiftcrypto.bitbox02.AntiKleptoHostNonceCommitment")
	proto.RegisterType((*AntiKleptoSignerCommitment)(nil), "shiftcrypto.bitbox02.AntiKleptoSignerCommitment")
	proto.RegisterType((*AntiKleptoSignatureRequest)(nil), "shiftcrypto.bitbox02.AntiKleptoSignatureRequest")
}

func init() { proto.RegisterFile("antiklepto.proto", fileDescriptor_39dadd9b6cadc7c0) }

var fileDescriptor_39dadd9b6cadc7c0 = []byte{
	// 161 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0xcc, 0x2b, 0xc9,
	0xcc, 0xce, 0x49, 0x2d, 0x28, 0xc9, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x29, 0xce,
	0xc8, 0x4c, 0x2b, 0x49, 0x2e, 0xaa, 0x04, 0x09, 0x25, 0x65, 0x96, 0x24, 0xe5, 0x57, 0x18, 0x18,
	0x29, 0xd9, 0x73, 0xc9, 0x3a, 0xe6, 0x95, 0x64, 0x7a, 0x83, 0x55, 0x7a, 0xe4, 0x17, 0x97, 0xf8,
	0xe5, 0xe7, 0x25, 0xa7, 0x3a, 0xe7, 0xe7, 0xe6, 0x66, 0x96, 0xe4, 0xa6, 0xe6, 0x95, 0x08, 0xc9,
	0x71, 0x71, 0x25, 0xc3, 0x79, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x48, 0x22, 0x4a, 0x36,
	0x5c, 0x52, 0x08, 0x03, 0x82, 0x33, 0xd3, 0xf3, 0x52, 0x8b, 0x48, 0xd0, 0x6d, 0x8d, 0xae, 0x3b,
	0xb1, 0xa4, 0xb4, 0x28, 0x35, 0x28, 0xb5, 0xb0, 0x34, 0xb5, 0xb8, 0x44, 0x48, 0x96, 0x8b, 0x2b,
	0x23, 0xbf, 0xb8, 0x24, 0x3e, 0x0f, 0xe4, 0x26, 0xa8, 0x6e, 0xce, 0x0c, 0x98, 0x23, 0x93, 0xd8,
	0xc0, 0x1e, 0x33, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x5d, 0x99, 0x02, 0x07, 0xec, 0x00, 0x00,
	0x00,
}
