// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ias/ias.proto

package ias // import "github.com/oasislabs/ekiden/go/grpc/ias"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "github.com/oasislabs/ekiden/go/grpc/common"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type GetSPIDInfoRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSPIDInfoRequest) Reset()         { *m = GetSPIDInfoRequest{} }
func (m *GetSPIDInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetSPIDInfoRequest) ProtoMessage()    {}
func (*GetSPIDInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ias_8ae2fc785feacc5f, []int{0}
}
func (m *GetSPIDInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSPIDInfoRequest.Unmarshal(m, b)
}
func (m *GetSPIDInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSPIDInfoRequest.Marshal(b, m, deterministic)
}
func (dst *GetSPIDInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSPIDInfoRequest.Merge(dst, src)
}
func (m *GetSPIDInfoRequest) XXX_Size() int {
	return xxx_messageInfo_GetSPIDInfoRequest.Size(m)
}
func (m *GetSPIDInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSPIDInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSPIDInfoRequest proto.InternalMessageInfo

type GetSPIDInfoResponse struct {
	Spid                 []byte   `protobuf:"bytes,1,opt,name=spid,proto3" json:"spid,omitempty"`
	QuoteSignatureType   uint32   `protobuf:"varint,2,opt,name=quote_signature_type,json=quoteSignatureType,proto3" json:"quote_signature_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSPIDInfoResponse) Reset()         { *m = GetSPIDInfoResponse{} }
func (m *GetSPIDInfoResponse) String() string { return proto.CompactTextString(m) }
func (*GetSPIDInfoResponse) ProtoMessage()    {}
func (*GetSPIDInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ias_8ae2fc785feacc5f, []int{1}
}
func (m *GetSPIDInfoResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSPIDInfoResponse.Unmarshal(m, b)
}
func (m *GetSPIDInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSPIDInfoResponse.Marshal(b, m, deterministic)
}
func (dst *GetSPIDInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSPIDInfoResponse.Merge(dst, src)
}
func (m *GetSPIDInfoResponse) XXX_Size() int {
	return xxx_messageInfo_GetSPIDInfoResponse.Size(m)
}
func (m *GetSPIDInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSPIDInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSPIDInfoResponse proto.InternalMessageInfo

func (m *GetSPIDInfoResponse) GetSpid() []byte {
	if m != nil {
		return m.Spid
	}
	return nil
}

func (m *GetSPIDInfoResponse) GetQuoteSignatureType() uint32 {
	if m != nil {
		return m.QuoteSignatureType
	}
	return 0
}

type VerifyEvidenceRequest struct {
	// Signed blob should be a CBOR-serialized quote + PSE manifest.
	Evidence             *common.Signed `protobuf:"bytes,1,opt,name=evidence,proto3" json:"evidence,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *VerifyEvidenceRequest) Reset()         { *m = VerifyEvidenceRequest{} }
func (m *VerifyEvidenceRequest) String() string { return proto.CompactTextString(m) }
func (*VerifyEvidenceRequest) ProtoMessage()    {}
func (*VerifyEvidenceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ias_8ae2fc785feacc5f, []int{2}
}
func (m *VerifyEvidenceRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifyEvidenceRequest.Unmarshal(m, b)
}
func (m *VerifyEvidenceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifyEvidenceRequest.Marshal(b, m, deterministic)
}
func (dst *VerifyEvidenceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifyEvidenceRequest.Merge(dst, src)
}
func (m *VerifyEvidenceRequest) XXX_Size() int {
	return xxx_messageInfo_VerifyEvidenceRequest.Size(m)
}
func (m *VerifyEvidenceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifyEvidenceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_VerifyEvidenceRequest proto.InternalMessageInfo

func (m *VerifyEvidenceRequest) GetEvidence() *common.Signed {
	if m != nil {
		return m.Evidence
	}
	return nil
}

type VerifyEvidenceResponse struct {
	Avr                  []byte   `protobuf:"bytes,1,opt,name=avr,proto3" json:"avr,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	CertificateChain     []byte   `protobuf:"bytes,3,opt,name=certificate_chain,json=certificateChain,proto3" json:"certificate_chain,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerifyEvidenceResponse) Reset()         { *m = VerifyEvidenceResponse{} }
func (m *VerifyEvidenceResponse) String() string { return proto.CompactTextString(m) }
func (*VerifyEvidenceResponse) ProtoMessage()    {}
func (*VerifyEvidenceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ias_8ae2fc785feacc5f, []int{3}
}
func (m *VerifyEvidenceResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifyEvidenceResponse.Unmarshal(m, b)
}
func (m *VerifyEvidenceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifyEvidenceResponse.Marshal(b, m, deterministic)
}
func (dst *VerifyEvidenceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifyEvidenceResponse.Merge(dst, src)
}
func (m *VerifyEvidenceResponse) XXX_Size() int {
	return xxx_messageInfo_VerifyEvidenceResponse.Size(m)
}
func (m *VerifyEvidenceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifyEvidenceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VerifyEvidenceResponse proto.InternalMessageInfo

func (m *VerifyEvidenceResponse) GetAvr() []byte {
	if m != nil {
		return m.Avr
	}
	return nil
}

func (m *VerifyEvidenceResponse) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *VerifyEvidenceResponse) GetCertificateChain() []byte {
	if m != nil {
		return m.CertificateChain
	}
	return nil
}

type GetSigRLRequest struct {
	EpidGid              uint32   `protobuf:"varint,1,opt,name=epid_gid,json=epidGid,proto3" json:"epid_gid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSigRLRequest) Reset()         { *m = GetSigRLRequest{} }
func (m *GetSigRLRequest) String() string { return proto.CompactTextString(m) }
func (*GetSigRLRequest) ProtoMessage()    {}
func (*GetSigRLRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ias_8ae2fc785feacc5f, []int{4}
}
func (m *GetSigRLRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSigRLRequest.Unmarshal(m, b)
}
func (m *GetSigRLRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSigRLRequest.Marshal(b, m, deterministic)
}
func (dst *GetSigRLRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSigRLRequest.Merge(dst, src)
}
func (m *GetSigRLRequest) XXX_Size() int {
	return xxx_messageInfo_GetSigRLRequest.Size(m)
}
func (m *GetSigRLRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSigRLRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetSigRLRequest proto.InternalMessageInfo

func (m *GetSigRLRequest) GetEpidGid() uint32 {
	if m != nil {
		return m.EpidGid
	}
	return 0
}

type GetSigRLResponse struct {
	SigRl                []byte   `protobuf:"bytes,1,opt,name=sig_rl,json=sigRl,proto3" json:"sig_rl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetSigRLResponse) Reset()         { *m = GetSigRLResponse{} }
func (m *GetSigRLResponse) String() string { return proto.CompactTextString(m) }
func (*GetSigRLResponse) ProtoMessage()    {}
func (*GetSigRLResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ias_8ae2fc785feacc5f, []int{5}
}
func (m *GetSigRLResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetSigRLResponse.Unmarshal(m, b)
}
func (m *GetSigRLResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetSigRLResponse.Marshal(b, m, deterministic)
}
func (dst *GetSigRLResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetSigRLResponse.Merge(dst, src)
}
func (m *GetSigRLResponse) XXX_Size() int {
	return xxx_messageInfo_GetSigRLResponse.Size(m)
}
func (m *GetSigRLResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetSigRLResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetSigRLResponse proto.InternalMessageInfo

func (m *GetSigRLResponse) GetSigRl() []byte {
	if m != nil {
		return m.SigRl
	}
	return nil
}

func init() {
	proto.RegisterType((*GetSPIDInfoRequest)(nil), "ias.GetSPIDInfoRequest")
	proto.RegisterType((*GetSPIDInfoResponse)(nil), "ias.GetSPIDInfoResponse")
	proto.RegisterType((*VerifyEvidenceRequest)(nil), "ias.VerifyEvidenceRequest")
	proto.RegisterType((*VerifyEvidenceResponse)(nil), "ias.VerifyEvidenceResponse")
	proto.RegisterType((*GetSigRLRequest)(nil), "ias.GetSigRLRequest")
	proto.RegisterType((*GetSigRLResponse)(nil), "ias.GetSigRLResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// IASClient is the client API for IAS service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IASClient interface {
	// Get SPID-related information.
	GetSPIDInfo(ctx context.Context, in *GetSPIDInfoRequest, opts ...grpc.CallOption) (*GetSPIDInfoResponse, error)
	// Verify attestation evidence.
	VerifyEvidence(ctx context.Context, in *VerifyEvidenceRequest, opts ...grpc.CallOption) (*VerifyEvidenceResponse, error)
	// Get SigRL.
	GetSigRL(ctx context.Context, in *GetSigRLRequest, opts ...grpc.CallOption) (*GetSigRLResponse, error)
}

type iASClient struct {
	cc *grpc.ClientConn
}

func NewIASClient(cc *grpc.ClientConn) IASClient {
	return &iASClient{cc}
}

func (c *iASClient) GetSPIDInfo(ctx context.Context, in *GetSPIDInfoRequest, opts ...grpc.CallOption) (*GetSPIDInfoResponse, error) {
	out := new(GetSPIDInfoResponse)
	err := c.cc.Invoke(ctx, "/ias.IAS/GetSPIDInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iASClient) VerifyEvidence(ctx context.Context, in *VerifyEvidenceRequest, opts ...grpc.CallOption) (*VerifyEvidenceResponse, error) {
	out := new(VerifyEvidenceResponse)
	err := c.cc.Invoke(ctx, "/ias.IAS/VerifyEvidence", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iASClient) GetSigRL(ctx context.Context, in *GetSigRLRequest, opts ...grpc.CallOption) (*GetSigRLResponse, error) {
	out := new(GetSigRLResponse)
	err := c.cc.Invoke(ctx, "/ias.IAS/GetSigRL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IASServer is the server API for IAS service.
type IASServer interface {
	// Get SPID-related information.
	GetSPIDInfo(context.Context, *GetSPIDInfoRequest) (*GetSPIDInfoResponse, error)
	// Verify attestation evidence.
	VerifyEvidence(context.Context, *VerifyEvidenceRequest) (*VerifyEvidenceResponse, error)
	// Get SigRL.
	GetSigRL(context.Context, *GetSigRLRequest) (*GetSigRLResponse, error)
}

func RegisterIASServer(s *grpc.Server, srv IASServer) {
	s.RegisterService(&_IAS_serviceDesc, srv)
}

func _IAS_GetSPIDInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSPIDInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IASServer).GetSPIDInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ias.IAS/GetSPIDInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IASServer).GetSPIDInfo(ctx, req.(*GetSPIDInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAS_VerifyEvidence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyEvidenceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IASServer).VerifyEvidence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ias.IAS/VerifyEvidence",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IASServer).VerifyEvidence(ctx, req.(*VerifyEvidenceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAS_GetSigRL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSigRLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IASServer).GetSigRL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ias.IAS/GetSigRL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IASServer).GetSigRL(ctx, req.(*GetSigRLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IAS_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ias.IAS",
	HandlerType: (*IASServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSPIDInfo",
			Handler:    _IAS_GetSPIDInfo_Handler,
		},
		{
			MethodName: "VerifyEvidence",
			Handler:    _IAS_VerifyEvidence_Handler,
		},
		{
			MethodName: "GetSigRL",
			Handler:    _IAS_GetSigRL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ias/ias.proto",
}

func init() { proto.RegisterFile("ias/ias.proto", fileDescriptor_ias_8ae2fc785feacc5f) }

var fileDescriptor_ias_8ae2fc785feacc5f = []byte{
	// 396 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xc5, 0x18, 0x4a, 0x98, 0x36, 0x25, 0x6c, 0x13, 0x30, 0x86, 0x43, 0xe5, 0x0b, 0x2d, 0xa0,
	0x18, 0x95, 0x13, 0x47, 0x5a, 0x50, 0x15, 0xc1, 0x01, 0x6d, 0x10, 0x07, 0x38, 0x58, 0x1b, 0x7b,
	0xb2, 0x1d, 0x91, 0x7a, 0x5d, 0xef, 0x3a, 0x52, 0x7e, 0x26, 0xff, 0xa8, 0x5a, 0x67, 0xed, 0x7c,
	0x9e, 0x3c, 0x7e, 0x6f, 0x3c, 0xef, 0xf9, 0xcd, 0x40, 0x97, 0x84, 0x8e, 0x49, 0xe8, 0x61, 0x51,
	0x2a, 0xa3, 0x98, 0x4f, 0x42, 0x87, 0x27, 0xa9, 0xba, 0xbd, 0x55, 0x79, 0xbc, 0x7c, 0x2c, 0x99,
	0xa8, 0x0f, 0xec, 0x1a, 0xcd, 0xf8, 0xe7, 0xe8, 0xeb, 0x28, 0x9f, 0x2a, 0x8e, 0x77, 0x15, 0x6a,
	0x13, 0xfd, 0x85, 0x93, 0x0d, 0x54, 0x17, 0x2a, 0xd7, 0xc8, 0x18, 0x3c, 0xd2, 0x05, 0x65, 0x81,
	0x77, 0xea, 0x9d, 0x1d, 0xf1, 0xba, 0x66, 0x1f, 0xa1, 0x7f, 0x57, 0x29, 0x83, 0x89, 0x26, 0x99,
	0x0b, 0x53, 0x95, 0x98, 0x98, 0x45, 0x81, 0xc1, 0xc3, 0x53, 0xef, 0xac, 0xcb, 0x59, 0xcd, 0x8d,
	0x1b, 0xea, 0xd7, 0xa2, 0xc0, 0xe8, 0x0a, 0x06, 0xbf, 0xb1, 0xa4, 0xe9, 0xe2, 0xdb, 0x9c, 0x32,
	0xcc, 0x53, 0x74, 0xaa, 0xec, 0x1d, 0x74, 0xd0, 0x41, 0xb5, 0xc4, 0xe1, 0xc5, 0xf1, 0xd0, 0x99,
	0xb5, 0x13, 0x30, 0xe3, 0x2d, 0x1f, 0x55, 0xf0, 0x62, 0x7b, 0x88, 0x33, 0xd9, 0x03, 0x5f, 0xcc,
	0x4b, 0xe7, 0xd1, 0x96, 0xec, 0x0d, 0x3c, 0x6d, 0xcd, 0xd5, 0xbe, 0x8e, 0xf8, 0x0a, 0x60, 0xef,
	0xe1, 0x79, 0x8a, 0xa5, 0xa1, 0x29, 0xa5, 0xc2, 0x60, 0x92, 0xde, 0x08, 0xca, 0x03, 0xbf, 0xee,
	0xea, 0xad, 0x11, 0x57, 0x16, 0x8f, 0x3e, 0xc0, 0x33, 0x1b, 0x0c, 0x49, 0xfe, 0xa3, 0x71, 0xfd,
	0x0a, 0x3a, 0x58, 0x50, 0x96, 0x48, 0x17, 0x4c, 0x97, 0x3f, 0xb1, 0xef, 0xd7, 0x94, 0x45, 0xe7,
	0xd0, 0x5b, 0x75, 0x3b, 0x7b, 0x03, 0x38, 0xd0, 0x24, 0x93, 0x72, 0xe6, 0x1c, 0x3e, 0xd6, 0x24,
	0xf9, 0xec, 0xe2, 0xbf, 0x07, 0xfe, 0xe8, 0xcb, 0x98, 0x5d, 0xc2, 0xe1, 0x5a, 0xf2, 0xec, 0xe5,
	0xd0, 0x2e, 0x71, 0x77, 0x43, 0x61, 0xb0, 0x4b, 0x2c, 0x05, 0xa2, 0x07, 0xec, 0x3b, 0x1c, 0x6f,
	0x66, 0xc3, 0xc2, 0xba, 0x7b, 0x6f, 0xea, 0xe1, 0xeb, 0xbd, 0x5c, 0x3b, 0xec, 0x33, 0x74, 0x9a,
	0x7f, 0x60, 0xfd, 0x56, 0x74, 0x2d, 0x80, 0x70, 0xb0, 0x85, 0x36, 0x9f, 0x5e, 0x9e, 0xff, 0x79,
	0x2b, 0xc9, 0xdc, 0x54, 0x13, 0xbb, 0xc5, 0x58, 0x09, 0x4d, 0x7a, 0x26, 0x26, 0x3a, 0xc6, 0x7f,
	0x56, 0x25, 0x96, 0x2a, 0x96, 0x65, 0x91, 0xda, 0x33, 0x9d, 0x1c, 0xd4, 0xd7, 0xf8, 0xe9, 0x3e,
	0x00, 0x00, 0xff, 0xff, 0x1b, 0xc4, 0xa8, 0x13, 0xb8, 0x02, 0x00, 0x00,
}