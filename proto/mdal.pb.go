// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mdal.proto

package mdalgrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type AggFunc int32

const (
	AggFunc_RAW   AggFunc = 0
	AggFunc_MEAN  AggFunc = 1
	AggFunc_MIN   AggFunc = 2
	AggFunc_MAX   AggFunc = 3
	AggFunc_COUNT AggFunc = 4
	AggFunc_SUM   AggFunc = 5
)

var AggFunc_name = map[int32]string{
	0: "RAW",
	1: "MEAN",
	2: "MIN",
	3: "MAX",
	4: "COUNT",
	5: "SUM",
}
var AggFunc_value = map[string]int32{
	"RAW":   0,
	"MEAN":  1,
	"MIN":   2,
	"MAX":   3,
	"COUNT": 4,
	"SUM":   5,
}

func (x AggFunc) String() string {
	return proto.EnumName(AggFunc_name, int32(x))
}
func (AggFunc) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_mdal_086aa3dd14e90695, []int{0}
}

type DataQueryRequest struct {
	Composition          []string                `protobuf:"bytes,1,rep,name=composition" json:"composition,omitempty"`
	Aggregation          map[string]*Aggregation `protobuf:"bytes,2,rep,name=aggregation" json:"aggregation,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Variables            map[string]*Variable    `protobuf:"bytes,3,rep,name=variables" json:"variables,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Time                 *TimeParams             `protobuf:"bytes,4,opt,name=time" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *DataQueryRequest) Reset()         { *m = DataQueryRequest{} }
func (m *DataQueryRequest) String() string { return proto.CompactTextString(m) }
func (*DataQueryRequest) ProtoMessage()    {}
func (*DataQueryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mdal_086aa3dd14e90695, []int{0}
}
func (m *DataQueryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataQueryRequest.Unmarshal(m, b)
}
func (m *DataQueryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataQueryRequest.Marshal(b, m, deterministic)
}
func (dst *DataQueryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataQueryRequest.Merge(dst, src)
}
func (m *DataQueryRequest) XXX_Size() int {
	return xxx_messageInfo_DataQueryRequest.Size(m)
}
func (m *DataQueryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DataQueryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DataQueryRequest proto.InternalMessageInfo

func (m *DataQueryRequest) GetComposition() []string {
	if m != nil {
		return m.Composition
	}
	return nil
}

func (m *DataQueryRequest) GetAggregation() map[string]*Aggregation {
	if m != nil {
		return m.Aggregation
	}
	return nil
}

func (m *DataQueryRequest) GetVariables() map[string]*Variable {
	if m != nil {
		return m.Variables
	}
	return nil
}

func (m *DataQueryRequest) GetTime() *TimeParams {
	if m != nil {
		return m.Time
	}
	return nil
}

type Variable struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Definition           string   `protobuf:"bytes,2,opt,name=definition" json:"definition,omitempty"`
	Uuids                [][]byte `protobuf:"bytes,3,rep,name=uuids,proto3" json:"uuids,omitempty"`
	Units                string   `protobuf:"bytes,4,opt,name=units" json:"units,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Variable) Reset()         { *m = Variable{} }
func (m *Variable) String() string { return proto.CompactTextString(m) }
func (*Variable) ProtoMessage()    {}
func (*Variable) Descriptor() ([]byte, []int) {
	return fileDescriptor_mdal_086aa3dd14e90695, []int{1}
}
func (m *Variable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Variable.Unmarshal(m, b)
}
func (m *Variable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Variable.Marshal(b, m, deterministic)
}
func (dst *Variable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Variable.Merge(dst, src)
}
func (m *Variable) XXX_Size() int {
	return xxx_messageInfo_Variable.Size(m)
}
func (m *Variable) XXX_DiscardUnknown() {
	xxx_messageInfo_Variable.DiscardUnknown(m)
}

var xxx_messageInfo_Variable proto.InternalMessageInfo

func (m *Variable) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Variable) GetDefinition() string {
	if m != nil {
		return m.Definition
	}
	return ""
}

func (m *Variable) GetUuids() [][]byte {
	if m != nil {
		return m.Uuids
	}
	return nil
}

func (m *Variable) GetUnits() string {
	if m != nil {
		return m.Units
	}
	return ""
}

type TimeParams struct {
	Start                string   `protobuf:"bytes,1,opt,name=start" json:"start,omitempty"`
	End                  string   `protobuf:"bytes,2,opt,name=end" json:"end,omitempty"`
	Window               string   `protobuf:"bytes,3,opt,name=window" json:"window,omitempty"`
	Aligned              bool     `protobuf:"varint,4,opt,name=aligned" json:"aligned,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TimeParams) Reset()         { *m = TimeParams{} }
func (m *TimeParams) String() string { return proto.CompactTextString(m) }
func (*TimeParams) ProtoMessage()    {}
func (*TimeParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_mdal_086aa3dd14e90695, []int{2}
}
func (m *TimeParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TimeParams.Unmarshal(m, b)
}
func (m *TimeParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TimeParams.Marshal(b, m, deterministic)
}
func (dst *TimeParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TimeParams.Merge(dst, src)
}
func (m *TimeParams) XXX_Size() int {
	return xxx_messageInfo_TimeParams.Size(m)
}
func (m *TimeParams) XXX_DiscardUnknown() {
	xxx_messageInfo_TimeParams.DiscardUnknown(m)
}

var xxx_messageInfo_TimeParams proto.InternalMessageInfo

func (m *TimeParams) GetStart() string {
	if m != nil {
		return m.Start
	}
	return ""
}

func (m *TimeParams) GetEnd() string {
	if m != nil {
		return m.End
	}
	return ""
}

func (m *TimeParams) GetWindow() string {
	if m != nil {
		return m.Window
	}
	return ""
}

func (m *TimeParams) GetAligned() bool {
	if m != nil {
		return m.Aligned
	}
	return false
}

type Aggregation struct {
	Funcs                []AggFunc `protobuf:"varint,1,rep,packed,name=funcs,enum=mdalgrpc.AggFunc" json:"funcs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Aggregation) Reset()         { *m = Aggregation{} }
func (m *Aggregation) String() string { return proto.CompactTextString(m) }
func (*Aggregation) ProtoMessage()    {}
func (*Aggregation) Descriptor() ([]byte, []int) {
	return fileDescriptor_mdal_086aa3dd14e90695, []int{3}
}
func (m *Aggregation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Aggregation.Unmarshal(m, b)
}
func (m *Aggregation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Aggregation.Marshal(b, m, deterministic)
}
func (dst *Aggregation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Aggregation.Merge(dst, src)
}
func (m *Aggregation) XXX_Size() int {
	return xxx_messageInfo_Aggregation.Size(m)
}
func (m *Aggregation) XXX_DiscardUnknown() {
	xxx_messageInfo_Aggregation.DiscardUnknown(m)
}

var xxx_messageInfo_Aggregation proto.InternalMessageInfo

func (m *Aggregation) GetFuncs() []AggFunc {
	if m != nil {
		return m.Funcs
	}
	return nil
}

type DataQueryResponse struct {
	// TODO: need to transform this from list of triples to a map
	// keyed by UUID of the triples relevant to that UUID
	Triples []*Triple `protobuf:"bytes,1,rep,name=triples" json:"triples,omitempty"`
	Uuids   [][]byte  `protobuf:"bytes,2,rep,name=uuids,proto3" json:"uuids,omitempty"`
	// apache arrow serialized response
	Arrow                []byte   `protobuf:"bytes,3,opt,name=arrow,proto3" json:"arrow,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DataQueryResponse) Reset()         { *m = DataQueryResponse{} }
func (m *DataQueryResponse) String() string { return proto.CompactTextString(m) }
func (*DataQueryResponse) ProtoMessage()    {}
func (*DataQueryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mdal_086aa3dd14e90695, []int{4}
}
func (m *DataQueryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataQueryResponse.Unmarshal(m, b)
}
func (m *DataQueryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataQueryResponse.Marshal(b, m, deterministic)
}
func (dst *DataQueryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataQueryResponse.Merge(dst, src)
}
func (m *DataQueryResponse) XXX_Size() int {
	return xxx_messageInfo_DataQueryResponse.Size(m)
}
func (m *DataQueryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DataQueryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DataQueryResponse proto.InternalMessageInfo

func (m *DataQueryResponse) GetTriples() []*Triple {
	if m != nil {
		return m.Triples
	}
	return nil
}

func (m *DataQueryResponse) GetUuids() [][]byte {
	if m != nil {
		return m.Uuids
	}
	return nil
}

func (m *DataQueryResponse) GetArrow() []byte {
	if m != nil {
		return m.Arrow
	}
	return nil
}

type Triple struct {
	Subject              string   `protobuf:"bytes,1,opt,name=subject" json:"subject,omitempty"`
	Property             string   `protobuf:"bytes,2,opt,name=property" json:"property,omitempty"`
	Value                string   `protobuf:"bytes,3,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Triple) Reset()         { *m = Triple{} }
func (m *Triple) String() string { return proto.CompactTextString(m) }
func (*Triple) ProtoMessage()    {}
func (*Triple) Descriptor() ([]byte, []int) {
	return fileDescriptor_mdal_086aa3dd14e90695, []int{5}
}
func (m *Triple) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Triple.Unmarshal(m, b)
}
func (m *Triple) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Triple.Marshal(b, m, deterministic)
}
func (dst *Triple) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Triple.Merge(dst, src)
}
func (m *Triple) XXX_Size() int {
	return xxx_messageInfo_Triple.Size(m)
}
func (m *Triple) XXX_DiscardUnknown() {
	xxx_messageInfo_Triple.DiscardUnknown(m)
}

var xxx_messageInfo_Triple proto.InternalMessageInfo

func (m *Triple) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Triple) GetProperty() string {
	if m != nil {
		return m.Property
	}
	return ""
}

func (m *Triple) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*DataQueryRequest)(nil), "mdalgrpc.DataQueryRequest")
	proto.RegisterMapType((map[string]*Aggregation)(nil), "mdalgrpc.DataQueryRequest.AggregationEntry")
	proto.RegisterMapType((map[string]*Variable)(nil), "mdalgrpc.DataQueryRequest.VariablesEntry")
	proto.RegisterType((*Variable)(nil), "mdalgrpc.Variable")
	proto.RegisterType((*TimeParams)(nil), "mdalgrpc.TimeParams")
	proto.RegisterType((*Aggregation)(nil), "mdalgrpc.Aggregation")
	proto.RegisterType((*DataQueryResponse)(nil), "mdalgrpc.DataQueryResponse")
	proto.RegisterType((*Triple)(nil), "mdalgrpc.Triple")
	proto.RegisterEnum("mdalgrpc.AggFunc", AggFunc_name, AggFunc_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MDALClient is the client API for MDAL service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MDALClient interface {
	DataQuery(ctx context.Context, in *DataQueryRequest, opts ...grpc.CallOption) (MDAL_DataQueryClient, error)
}

type mDALClient struct {
	cc *grpc.ClientConn
}

func NewMDALClient(cc *grpc.ClientConn) MDALClient {
	return &mDALClient{cc}
}

func (c *mDALClient) DataQuery(ctx context.Context, in *DataQueryRequest, opts ...grpc.CallOption) (MDAL_DataQueryClient, error) {
	stream, err := c.cc.NewStream(ctx, &_MDAL_serviceDesc.Streams[0], "/mdalgrpc.MDAL/DataQuery", opts...)
	if err != nil {
		return nil, err
	}
	x := &mDALDataQueryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MDAL_DataQueryClient interface {
	Recv() (*DataQueryResponse, error)
	grpc.ClientStream
}

type mDALDataQueryClient struct {
	grpc.ClientStream
}

func (x *mDALDataQueryClient) Recv() (*DataQueryResponse, error) {
	m := new(DataQueryResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MDALServer is the server API for MDAL service.
type MDALServer interface {
	DataQuery(*DataQueryRequest, MDAL_DataQueryServer) error
}

func RegisterMDALServer(s *grpc.Server, srv MDALServer) {
	s.RegisterService(&_MDAL_serviceDesc, srv)
}

func _MDAL_DataQuery_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DataQueryRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MDALServer).DataQuery(m, &mDALDataQueryServer{stream})
}

type MDAL_DataQueryServer interface {
	Send(*DataQueryResponse) error
	grpc.ServerStream
}

type mDALDataQueryServer struct {
	grpc.ServerStream
}

func (x *mDALDataQueryServer) Send(m *DataQueryResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _MDAL_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mdalgrpc.MDAL",
	HandlerType: (*MDALServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DataQuery",
			Handler:       _MDAL_DataQuery_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mdal.proto",
}

func init() { proto.RegisterFile("mdal.proto", fileDescriptor_mdal_086aa3dd14e90695) }

var fileDescriptor_mdal_086aa3dd14e90695 = []byte{
	// 531 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0xdf, 0x6f, 0xd3, 0x30,
	0x10, 0xc7, 0x49, 0x93, 0xfe, 0xc8, 0x75, 0x9a, 0xb2, 0xd3, 0x40, 0x51, 0x91, 0x50, 0x95, 0x17,
	0xc2, 0x26, 0x55, 0xa8, 0x48, 0x08, 0xf1, 0x56, 0xd8, 0x86, 0x90, 0x68, 0x19, 0xa6, 0x05, 0x5e,
	0xdd, 0xc6, 0xad, 0xbc, 0x35, 0x3f, 0x70, 0x9c, 0x4d, 0xfd, 0x6b, 0xf8, 0x57, 0x91, 0xed, 0x64,
	0x49, 0xab, 0xc1, 0xdb, 0x7d, 0xcf, 0x77, 0x1f, 0xdb, 0x5f, 0x5f, 0x02, 0x10, 0x47, 0x74, 0x3b,
	0xca, 0x44, 0x2a, 0x53, 0xec, 0xa9, 0x78, 0x23, 0xb2, 0x55, 0xf0, 0xc7, 0x06, 0xef, 0x82, 0x4a,
	0xfa, 0xad, 0x60, 0x62, 0x47, 0xd8, 0xef, 0x82, 0xe5, 0x12, 0x87, 0xd0, 0x5f, 0xa5, 0x71, 0x96,
	0xe6, 0x5c, 0xf2, 0x34, 0xf1, 0xad, 0xa1, 0x1d, 0xba, 0xa4, 0x99, 0xc2, 0x29, 0xf4, 0xe9, 0x66,
	0x23, 0xd8, 0x86, 0xea, 0x8a, 0xd6, 0xd0, 0x0e, 0xfb, 0xe3, 0xf3, 0x51, 0x85, 0x1d, 0x1d, 0x22,
	0x47, 0x93, 0xba, 0xfa, 0x32, 0x91, 0x62, 0x47, 0x9a, 0xfd, 0xf8, 0x09, 0xdc, 0x3b, 0x2a, 0x38,
	0x5d, 0x6e, 0x59, 0xee, 0xdb, 0x1a, 0xf6, 0xea, 0x3f, 0xb0, 0x1f, 0x55, 0xad, 0x41, 0xd5, 0xbd,
	0x18, 0x82, 0x23, 0x79, 0xcc, 0x7c, 0x67, 0x68, 0x85, 0xfd, 0xf1, 0x69, 0xcd, 0x98, 0xf3, 0x98,
	0x5d, 0x53, 0x41, 0xe3, 0x9c, 0xe8, 0x8a, 0xc1, 0x02, 0xbc, 0xc3, 0x33, 0xa1, 0x07, 0xf6, 0x2d,
	0xdb, 0xf9, 0xd6, 0xd0, 0x0a, 0x5d, 0xa2, 0x42, 0x3c, 0x87, 0xf6, 0x1d, 0xdd, 0x16, 0xcc, 0x6f,
	0x69, 0xe0, 0xd3, 0x1a, 0xd8, 0x68, 0x26, 0xa6, 0xe6, 0x7d, 0xeb, 0x9d, 0x35, 0xb8, 0x86, 0xe3,
	0xfd, 0xd3, 0x3d, 0x02, 0x0d, 0xf7, 0xa1, 0x58, 0x43, 0xab, 0xd6, 0x06, 0x31, 0xb8, 0x81, 0x5e,
	0x95, 0x46, 0x04, 0x27, 0xa1, 0x31, 0x2b, 0x61, 0x3a, 0xc6, 0x17, 0x00, 0x11, 0x5b, 0xf3, 0x84,
	0x97, 0x2f, 0xa1, 0x56, 0x1a, 0x19, 0x3c, 0x85, 0x76, 0x51, 0xf0, 0xc8, 0xf8, 0x7a, 0x44, 0x8c,
	0xd0, 0xd9, 0x84, 0xcb, 0x5c, 0x3b, 0xe5, 0x12, 0x23, 0x82, 0x35, 0x40, 0x6d, 0x94, 0xaa, 0xc9,
	0x25, 0x15, 0xb2, 0xdc, 0xce, 0x08, 0x75, 0x1f, 0x96, 0x44, 0xe5, 0x46, 0x2a, 0xc4, 0x67, 0xd0,
	0xb9, 0xe7, 0x49, 0x94, 0xde, 0xfb, 0xb6, 0x4e, 0x96, 0x0a, 0x7d, 0xe8, 0xd2, 0x2d, 0xdf, 0x24,
	0x2c, 0xd2, 0xbb, 0xf4, 0x48, 0x25, 0x83, 0xb7, 0xd0, 0x6f, 0xf8, 0x87, 0x2f, 0xa1, 0xbd, 0x2e,
	0x92, 0x55, 0xae, 0x27, 0xed, 0x78, 0x7c, 0xb2, 0xe7, 0xf2, 0x55, 0x91, 0xac, 0x88, 0x59, 0x0f,
	0x6e, 0xe1, 0xa4, 0x31, 0x0c, 0x79, 0x96, 0x26, 0x39, 0xc3, 0x33, 0xe8, 0x4a, 0xc1, 0x33, 0x35,
	0x3a, 0x96, 0x1e, 0x1d, 0xaf, 0xf1, 0xec, 0x7a, 0x81, 0x54, 0x05, 0xb5, 0x19, 0xad, 0x03, 0x33,
	0xa8, 0x10, 0xe5, 0xf9, 0x8f, 0x88, 0x11, 0xc1, 0x1c, 0x3a, 0xa6, 0x5d, 0x5d, 0x24, 0x2f, 0x96,
	0x37, 0x6c, 0x55, 0x59, 0x51, 0x49, 0x1c, 0x40, 0x2f, 0x13, 0x69, 0xc6, 0x84, 0xdc, 0x95, 0x8e,
	0x3c, 0x68, 0x45, 0x35, 0xcf, 0x6c, 0x5c, 0x31, 0xe2, 0xec, 0x03, 0x74, 0xcb, 0x4b, 0x61, 0x17,
	0x6c, 0x32, 0xf9, 0xe9, 0x3d, 0xc1, 0x1e, 0x38, 0xd3, 0xcb, 0xc9, 0xcc, 0xb3, 0x54, 0x6a, 0xfa,
	0x79, 0xe6, 0xb5, 0x74, 0x30, 0xf9, 0xe5, 0xd9, 0xe8, 0x42, 0xfb, 0xe3, 0xd7, 0xc5, 0x6c, 0xee,
	0x39, 0x2a, 0xf7, 0x7d, 0x31, 0xf5, 0xda, 0xe3, 0x19, 0x38, 0xd3, 0x8b, 0xc9, 0x17, 0xbc, 0x02,
	0xf7, 0xc1, 0x0e, 0x1c, 0xfc, 0xfb, 0x83, 0x19, 0x3c, 0x7f, 0x74, 0xcd, 0xf8, 0xf7, 0xda, 0x5a,
	0x76, 0xf4, 0x5f, 0xe1, 0xcd, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9b, 0x31, 0xc1, 0x23, 0x23,
	0x04, 0x00, 0x00,
}
