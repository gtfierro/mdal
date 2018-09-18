// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mdal.proto

package mdalgrpc

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	math "math"
)

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
	return fileDescriptor_20e1752e2497525f, []int{0}
}

type DataQueryRequest struct {
	Composition          []string                `protobuf:"bytes,1,rep,name=composition,proto3" json:"composition,omitempty"`
	Aggregation          map[string]*Aggregation `protobuf:"bytes,2,rep,name=aggregation,proto3" json:"aggregation,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Variables            map[string]*Variable    `protobuf:"bytes,3,rep,name=variables,proto3" json:"variables,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Time                 *TimeParams             `protobuf:"bytes,4,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *DataQueryRequest) Reset()         { *m = DataQueryRequest{} }
func (m *DataQueryRequest) String() string { return proto.CompactTextString(m) }
func (*DataQueryRequest) ProtoMessage()    {}
func (*DataQueryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_20e1752e2497525f, []int{0}
}
func (m *DataQueryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataQueryRequest.Unmarshal(m, b)
}
func (m *DataQueryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataQueryRequest.Marshal(b, m, deterministic)
}
func (m *DataQueryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataQueryRequest.Merge(m, src)
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
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Definition           string   `protobuf:"bytes,2,opt,name=definition,proto3" json:"definition,omitempty"`
	Uuids                []string `protobuf:"bytes,3,rep,name=uuids,proto3" json:"uuids,omitempty"`
	Units                string   `protobuf:"bytes,4,opt,name=units,proto3" json:"units,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Variable) Reset()         { *m = Variable{} }
func (m *Variable) String() string { return proto.CompactTextString(m) }
func (*Variable) ProtoMessage()    {}
func (*Variable) Descriptor() ([]byte, []int) {
	return fileDescriptor_20e1752e2497525f, []int{1}
}
func (m *Variable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Variable.Unmarshal(m, b)
}
func (m *Variable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Variable.Marshal(b, m, deterministic)
}
func (m *Variable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Variable.Merge(m, src)
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

func (m *Variable) GetUuids() []string {
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
	Start                string   `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End                  string   `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	Window               string   `protobuf:"bytes,3,opt,name=window,proto3" json:"window,omitempty"`
	Aligned              bool     `protobuf:"varint,4,opt,name=aligned,proto3" json:"aligned,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TimeParams) Reset()         { *m = TimeParams{} }
func (m *TimeParams) String() string { return proto.CompactTextString(m) }
func (*TimeParams) ProtoMessage()    {}
func (*TimeParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_20e1752e2497525f, []int{2}
}
func (m *TimeParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TimeParams.Unmarshal(m, b)
}
func (m *TimeParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TimeParams.Marshal(b, m, deterministic)
}
func (m *TimeParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TimeParams.Merge(m, src)
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
	Funcs                []AggFunc `protobuf:"varint,1,rep,packed,name=funcs,proto3,enum=mdalgrpc.AggFunc" json:"funcs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Aggregation) Reset()         { *m = Aggregation{} }
func (m *Aggregation) String() string { return proto.CompactTextString(m) }
func (*Aggregation) ProtoMessage()    {}
func (*Aggregation) Descriptor() ([]byte, []int) {
	return fileDescriptor_20e1752e2497525f, []int{3}
}
func (m *Aggregation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Aggregation.Unmarshal(m, b)
}
func (m *Aggregation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Aggregation.Marshal(b, m, deterministic)
}
func (m *Aggregation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Aggregation.Merge(m, src)
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
	Rows []*Row `protobuf:"bytes,1,rep,name=rows,proto3" json:"rows,omitempty"`
	// variable name -> list of uuids
	Mapping map[string]*VarMap `protobuf:"bytes,2,rep,name=mapping,proto3" json:"mapping,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// uuid -> triple context
	Context []*Row   `protobuf:"bytes,3,rep,name=context,proto3" json:"context,omitempty"`
	Uuids   []string `protobuf:"bytes,4,rep,name=uuids,proto3" json:"uuids,omitempty"`
	// apache arrow serialized response
	Arrow                []byte        `protobuf:"bytes,5,opt,name=arrow,proto3" json:"arrow,omitempty"`
	Error                string        `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
	Times                []int64       `protobuf:"varint,7,rep,packed,name=times,proto3" json:"times,omitempty"`
	Values               []*ValueArray `protobuf:"bytes,8,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *DataQueryResponse) Reset()         { *m = DataQueryResponse{} }
func (m *DataQueryResponse) String() string { return proto.CompactTextString(m) }
func (*DataQueryResponse) ProtoMessage()    {}
func (*DataQueryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_20e1752e2497525f, []int{4}
}
func (m *DataQueryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataQueryResponse.Unmarshal(m, b)
}
func (m *DataQueryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataQueryResponse.Marshal(b, m, deterministic)
}
func (m *DataQueryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataQueryResponse.Merge(m, src)
}
func (m *DataQueryResponse) XXX_Size() int {
	return xxx_messageInfo_DataQueryResponse.Size(m)
}
func (m *DataQueryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DataQueryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DataQueryResponse proto.InternalMessageInfo

func (m *DataQueryResponse) GetRows() []*Row {
	if m != nil {
		return m.Rows
	}
	return nil
}

func (m *DataQueryResponse) GetMapping() map[string]*VarMap {
	if m != nil {
		return m.Mapping
	}
	return nil
}

func (m *DataQueryResponse) GetContext() []*Row {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *DataQueryResponse) GetUuids() []string {
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

func (m *DataQueryResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *DataQueryResponse) GetTimes() []int64 {
	if m != nil {
		return m.Times
	}
	return nil
}

func (m *DataQueryResponse) GetValues() []*ValueArray {
	if m != nil {
		return m.Values
	}
	return nil
}

type Row struct {
	Uuid                 []byte            `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Row                  map[string]string `protobuf:"bytes,2,rep,name=row,proto3" json:"row,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Row) Reset()         { *m = Row{} }
func (m *Row) String() string { return proto.CompactTextString(m) }
func (*Row) ProtoMessage()    {}
func (*Row) Descriptor() ([]byte, []int) {
	return fileDescriptor_20e1752e2497525f, []int{5}
}
func (m *Row) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Row.Unmarshal(m, b)
}
func (m *Row) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Row.Marshal(b, m, deterministic)
}
func (m *Row) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Row.Merge(m, src)
}
func (m *Row) XXX_Size() int {
	return xxx_messageInfo_Row.Size(m)
}
func (m *Row) XXX_DiscardUnknown() {
	xxx_messageInfo_Row.DiscardUnknown(m)
}

var xxx_messageInfo_Row proto.InternalMessageInfo

func (m *Row) GetUuid() []byte {
	if m != nil {
		return m.Uuid
	}
	return nil
}

func (m *Row) GetRow() map[string]string {
	if m != nil {
		return m.Row
	}
	return nil
}

type VarMap struct {
	Uuids                [][]byte `protobuf:"bytes,1,rep,name=uuids,proto3" json:"uuids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VarMap) Reset()         { *m = VarMap{} }
func (m *VarMap) String() string { return proto.CompactTextString(m) }
func (*VarMap) ProtoMessage()    {}
func (*VarMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_20e1752e2497525f, []int{6}
}
func (m *VarMap) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VarMap.Unmarshal(m, b)
}
func (m *VarMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VarMap.Marshal(b, m, deterministic)
}
func (m *VarMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VarMap.Merge(m, src)
}
func (m *VarMap) XXX_Size() int {
	return xxx_messageInfo_VarMap.Size(m)
}
func (m *VarMap) XXX_DiscardUnknown() {
	xxx_messageInfo_VarMap.DiscardUnknown(m)
}

var xxx_messageInfo_VarMap proto.InternalMessageInfo

func (m *VarMap) GetUuids() [][]byte {
	if m != nil {
		return m.Uuids
	}
	return nil
}

type ValueArray struct {
	Value                []float64 `protobuf:"fixed64,1,rep,packed,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ValueArray) Reset()         { *m = ValueArray{} }
func (m *ValueArray) String() string { return proto.CompactTextString(m) }
func (*ValueArray) ProtoMessage()    {}
func (*ValueArray) Descriptor() ([]byte, []int) {
	return fileDescriptor_20e1752e2497525f, []int{7}
}
func (m *ValueArray) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValueArray.Unmarshal(m, b)
}
func (m *ValueArray) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValueArray.Marshal(b, m, deterministic)
}
func (m *ValueArray) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValueArray.Merge(m, src)
}
func (m *ValueArray) XXX_Size() int {
	return xxx_messageInfo_ValueArray.Size(m)
}
func (m *ValueArray) XXX_DiscardUnknown() {
	xxx_messageInfo_ValueArray.DiscardUnknown(m)
}

var xxx_messageInfo_ValueArray proto.InternalMessageInfo

func (m *ValueArray) GetValue() []float64 {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*DataQueryRequest)(nil), "mdalgrpc.DataQueryRequest")
	proto.RegisterMapType((map[string]*Aggregation)(nil), "mdalgrpc.DataQueryRequest.AggregationEntry")
	proto.RegisterMapType((map[string]*Variable)(nil), "mdalgrpc.DataQueryRequest.VariablesEntry")
	proto.RegisterType((*Variable)(nil), "mdalgrpc.Variable")
	proto.RegisterType((*TimeParams)(nil), "mdalgrpc.TimeParams")
	proto.RegisterType((*Aggregation)(nil), "mdalgrpc.Aggregation")
	proto.RegisterType((*DataQueryResponse)(nil), "mdalgrpc.DataQueryResponse")
	proto.RegisterMapType((map[string]*VarMap)(nil), "mdalgrpc.DataQueryResponse.MappingEntry")
	proto.RegisterType((*Row)(nil), "mdalgrpc.Row")
	proto.RegisterMapType((map[string]string)(nil), "mdalgrpc.Row.RowEntry")
	proto.RegisterType((*VarMap)(nil), "mdalgrpc.VarMap")
	proto.RegisterType((*ValueArray)(nil), "mdalgrpc.ValueArray")
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

func init() { proto.RegisterFile("mdal.proto", fileDescriptor_20e1752e2497525f) }

var fileDescriptor_20e1752e2497525f = []byte{
	// 712 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0xcd, 0x6e, 0xd3, 0x4c,
	0x14, 0xfd, 0x1c, 0x3b, 0x7f, 0x37, 0xf9, 0x2a, 0x77, 0xbe, 0x7c, 0x95, 0x09, 0xa8, 0x0a, 0x5e,
	0x50, 0xd3, 0xa2, 0x04, 0x82, 0x54, 0xa1, 0xee, 0x52, 0x5a, 0x10, 0x52, 0x53, 0xca, 0xd0, 0x02,
	0xdb, 0x69, 0x32, 0xb5, 0x0c, 0xf1, 0x8c, 0x3b, 0xb6, 0x1b, 0x22, 0x76, 0xbc, 0x02, 0x0f, 0xc0,
	0x43, 0xf1, 0x0a, 0xbc, 0x03, 0x5b, 0x34, 0x33, 0x36, 0x76, 0xa2, 0xd0, 0x45, 0xa4, 0x7b, 0xcf,
	0x9c, 0x7b, 0xee, 0xcc, 0xc9, 0x91, 0x01, 0xc2, 0x29, 0x99, 0xf5, 0x23, 0xc1, 0x13, 0x8e, 0x1a,
	0xb2, 0xf6, 0x45, 0x34, 0xe9, 0xde, 0xf3, 0x39, 0xf7, 0x67, 0x74, 0x40, 0xa2, 0x60, 0x40, 0x18,
	0xe3, 0x09, 0x49, 0x02, 0xce, 0x62, 0xcd, 0x73, 0xbf, 0x9b, 0x60, 0x1f, 0x91, 0x84, 0xbc, 0x49,
	0xa9, 0x58, 0x60, 0x7a, 0x9d, 0xd2, 0x38, 0x41, 0x3d, 0x68, 0x4d, 0x78, 0x18, 0xf1, 0x38, 0x90,
	0x54, 0xc7, 0xe8, 0x99, 0x5e, 0x13, 0x97, 0x21, 0x34, 0x86, 0x16, 0xf1, 0x7d, 0x41, 0x7d, 0x25,
	0xe6, 0x54, 0x7a, 0xa6, 0xd7, 0x1a, 0xee, 0xf5, 0xf3, 0xa5, 0xfd, 0x55, 0xc9, 0xfe, 0xa8, 0x60,
	0x1f, 0xb3, 0x44, 0x2c, 0x70, 0x79, 0x1e, 0xbd, 0x84, 0xe6, 0x0d, 0x11, 0x01, 0xb9, 0x9c, 0xd1,
	0xd8, 0x31, 0x95, 0xd8, 0xc3, 0x5b, 0xc4, 0xde, 0xe5, 0x5c, 0x2d, 0x55, 0xcc, 0x22, 0x0f, 0xac,
	0x24, 0x08, 0xa9, 0x63, 0xf5, 0x0c, 0xaf, 0x35, 0xec, 0x14, 0x1a, 0xe7, 0x41, 0x48, 0xcf, 0x88,
	0x20, 0x61, 0x8c, 0x15, 0xa3, 0x7b, 0x01, 0xf6, 0xea, 0x9d, 0x90, 0x0d, 0xe6, 0x27, 0xba, 0x70,
	0x8c, 0x9e, 0xe1, 0x35, 0xb1, 0x2c, 0xd1, 0x1e, 0x54, 0x6f, 0xc8, 0x2c, 0xa5, 0x4e, 0x45, 0x09,
	0xfe, 0x5f, 0x08, 0x96, 0x86, 0xb1, 0xe6, 0x1c, 0x54, 0x9e, 0x19, 0xdd, 0x33, 0xd8, 0x58, 0xbe,
	0xdd, 0x1a, 0x51, 0x6f, 0x59, 0x14, 0x15, 0xa2, 0xf9, 0x68, 0x49, 0xd1, 0xfd, 0x08, 0x8d, 0x1c,
	0x46, 0x08, 0x2c, 0x46, 0x42, 0x9a, 0x89, 0xa9, 0x1a, 0x6d, 0x03, 0x4c, 0xe9, 0x55, 0xc0, 0x82,
	0xec, 0x9f, 0x90, 0x27, 0x25, 0x04, 0x75, 0xa0, 0x9a, 0xa6, 0xc1, 0x54, 0xfb, 0xda, 0xc4, 0xba,
	0x51, 0x28, 0x0b, 0x92, 0x58, 0x39, 0x25, 0x51, 0xd9, 0xb8, 0x57, 0x00, 0x85, 0x51, 0x92, 0x13,
	0x27, 0x44, 0x24, 0xd9, 0x3a, 0xdd, 0xc8, 0xf7, 0x50, 0x36, 0xcd, 0x16, 0xc9, 0x12, 0x6d, 0x41,
	0x6d, 0x1e, 0xb0, 0x29, 0x9f, 0x3b, 0xa6, 0x02, 0xb3, 0x0e, 0x39, 0x50, 0x27, 0xb3, 0xc0, 0x67,
	0x74, 0xaa, 0xb6, 0x34, 0x70, 0xde, 0xba, 0xfb, 0xd0, 0x2a, 0xf9, 0x87, 0x76, 0xa0, 0x7a, 0x95,
	0xb2, 0x49, 0xac, 0x92, 0xb6, 0x31, 0xdc, 0x5c, 0x72, 0xf9, 0x45, 0xca, 0x26, 0x58, 0x9f, 0xbb,
	0xbf, 0x2a, 0xb0, 0x59, 0x4a, 0x43, 0x1c, 0x71, 0x16, 0x53, 0x74, 0x1f, 0x2c, 0xc1, 0xe7, 0x7a,
	0xba, 0x35, 0xfc, 0xb7, 0x98, 0xc6, 0x7c, 0x8e, 0xd5, 0x11, 0x3a, 0x84, 0x7a, 0x48, 0xa2, 0x28,
	0x60, 0x7e, 0x96, 0x55, 0x6f, 0x6d, 0xbc, 0xb4, 0x60, 0x7f, 0xac, 0xa9, 0x3a, 0x5d, 0xf9, 0x20,
	0xda, 0x81, 0xfa, 0x84, 0xb3, 0x84, 0x7e, 0x4e, 0xb2, 0x88, 0xae, 0x6c, 0xca, 0x4f, 0x0b, 0xc7,
	0xad, 0x15, 0xc7, 0x89, 0x10, 0x7c, 0xee, 0x54, 0x7b, 0x86, 0xd7, 0xc6, 0xba, 0x91, 0x28, 0x15,
	0x82, 0x0b, 0xa7, 0xa6, 0x3d, 0x56, 0x8d, 0x44, 0x65, 0x48, 0x63, 0xa7, 0xde, 0x33, 0x3d, 0x13,
	0xeb, 0x06, 0x3d, 0x82, 0x9a, 0x8a, 0x45, 0xec, 0x34, 0xd4, 0xfe, 0x4e, 0x39, 0x38, 0xb3, 0x94,
	0x8e, 0x84, 0x20, 0x0b, 0x9c, 0x71, 0xba, 0x27, 0xd0, 0x2e, 0xbf, 0x63, 0x4d, 0x0e, 0x1f, 0x2c,
	0xe7, 0xd0, 0x5e, 0xca, 0xe1, 0x98, 0x44, 0xe5, 0x14, 0x7e, 0x01, 0x13, 0xf3, 0xb9, 0x0c, 0xa0,
	0x7c, 0x8d, 0x52, 0x69, 0x63, 0x55, 0x23, 0x0f, 0x4c, 0xf9, 0x2c, 0xed, 0xeb, 0xd6, 0x92, 0x27,
	0xf2, 0xa7, 0x5d, 0x94, 0x94, 0xee, 0x3e, 0x34, 0x72, 0x60, 0xcd, 0x75, 0x3a, 0xe5, 0xeb, 0x34,
	0xcb, 0xcb, 0xb7, 0xa1, 0xa6, 0x6f, 0x54, 0x58, 0x2b, 0xff, 0xeb, 0x76, 0x66, 0xad, 0xeb, 0x02,
	0x14, 0x06, 0x14, 0x3a, 0x92, 0x63, 0x64, 0x3a, 0xbb, 0x87, 0x50, 0xcf, 0xc2, 0x84, 0xea, 0x60,
	0xe2, 0xd1, 0x7b, 0xfb, 0x1f, 0xd4, 0x00, 0x6b, 0x7c, 0x3c, 0x3a, 0xb5, 0x0d, 0x09, 0x8d, 0x5f,
	0x9d, 0xda, 0x15, 0x55, 0x8c, 0x3e, 0xd8, 0x26, 0x6a, 0x42, 0xf5, 0xf9, 0xeb, 0x8b, 0xd3, 0x73,
	0xdb, 0x92, 0xd8, 0xdb, 0x8b, 0xb1, 0x5d, 0x1d, 0x06, 0x60, 0x8d, 0x8f, 0x46, 0x27, 0x88, 0x40,
	0xf3, 0x4f, 0x68, 0x50, 0xf7, 0xef, 0x1f, 0xaa, 0xee, 0xdd, 0x5b, 0x52, 0xe6, 0xde, 0xf9, 0xfa,
	0xe3, 0xe7, 0xb7, 0xca, 0x7f, 0xee, 0xc6, 0xe0, 0xe6, 0xc9, 0x40, 0xf2, 0x06, 0xd7, 0xf2, 0xfc,
	0xc0, 0xd8, 0x7d, 0x6c, 0x5c, 0xd6, 0xd4, 0xe7, 0xf9, 0xe9, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xa1, 0xad, 0x29, 0x21, 0xd4, 0x05, 0x00, 0x00,
}
