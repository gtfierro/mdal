// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mdal.proto

package mdalgrpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Selector int32

const (
	Selector_RAW   Selector = 0
	Selector_MEAN  Selector = 1
	Selector_MIN   Selector = 2
	Selector_MAX   Selector = 3
	Selector_COUNT Selector = 4
	Selector_SUM   Selector = 5
)

var Selector_name = map[int32]string{
	0: "RAW",
	1: "MEAN",
	2: "MIN",
	3: "MAX",
	4: "COUNT",
	5: "SUM",
}
var Selector_value = map[string]int32{
	"RAW":   0,
	"MEAN":  1,
	"MIN":   2,
	"MAX":   3,
	"COUNT": 4,
	"SUM":   5,
}

func (x Selector) String() string {
	return proto.EnumName(Selector_name, int32(x))
}
func (Selector) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_mdal_5f24f6fd28c73679, []int{0}
}

type DataQueryRequest struct {
	Composition          []string             `protobuf:"bytes,1,rep,name=composition" json:"composition,omitempty"`
	Selectors            []Selector           `protobuf:"varint,2,rep,packed,name=selectors,enum=mdalgrpc.Selector" json:"selectors,omitempty"`
	Variables            map[string]*Variable `protobuf:"bytes,3,rep,name=variables" json:"variables,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Time                 *TimeParams          `protobuf:"bytes,4,opt,name=time" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *DataQueryRequest) Reset()         { *m = DataQueryRequest{} }
func (m *DataQueryRequest) String() string { return proto.CompactTextString(m) }
func (*DataQueryRequest) ProtoMessage()    {}
func (*DataQueryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_mdal_5f24f6fd28c73679, []int{0}
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

func (m *DataQueryRequest) GetSelectors() []Selector {
	if m != nil {
		return m.Selectors
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
	return fileDescriptor_mdal_5f24f6fd28c73679, []int{1}
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
	return fileDescriptor_mdal_5f24f6fd28c73679, []int{2}
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

type DataQueryResponse struct {
	// TODO: need to transform this from list of triples to a map
	// keyed by UUID of the triples relevant to that UUID
	Triples []*Triple `protobuf:"bytes,1,rep,name=triples" json:"triples,omitempty"`
	Uuids   [][]byte  `protobuf:"bytes,2,rep,name=uuids,proto3" json:"uuids,omitempty"`
	// apache arrow serialized response
	Arrow                [][]byte `protobuf:"bytes,3,rep,name=arrow,proto3" json:"arrow,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DataQueryResponse) Reset()         { *m = DataQueryResponse{} }
func (m *DataQueryResponse) String() string { return proto.CompactTextString(m) }
func (*DataQueryResponse) ProtoMessage()    {}
func (*DataQueryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_mdal_5f24f6fd28c73679, []int{3}
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

func (m *DataQueryResponse) GetArrow() [][]byte {
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
	return fileDescriptor_mdal_5f24f6fd28c73679, []int{4}
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
	proto.RegisterMapType((map[string]*Variable)(nil), "mdalgrpc.DataQueryRequest.VariablesEntry")
	proto.RegisterType((*Variable)(nil), "mdalgrpc.Variable")
	proto.RegisterType((*TimeParams)(nil), "mdalgrpc.TimeParams")
	proto.RegisterType((*DataQueryResponse)(nil), "mdalgrpc.DataQueryResponse")
	proto.RegisterType((*Triple)(nil), "mdalgrpc.Triple")
	proto.RegisterEnum("mdalgrpc.Selector", Selector_name, Selector_value)
}

func init() { proto.RegisterFile("mdal.proto", fileDescriptor_mdal_5f24f6fd28c73679) }

var fileDescriptor_mdal_5f24f6fd28c73679 = []byte{
	// 484 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x53, 0x5f, 0x6f, 0xd3, 0x3e,
	0x14, 0xfd, 0xe5, 0x4f, 0xdb, 0xe4, 0xf6, 0xa7, 0x29, 0x5c, 0x4d, 0x28, 0x2a, 0x12, 0x8a, 0xfa,
	0x14, 0xf6, 0x50, 0x4d, 0xe5, 0x05, 0xf1, 0x56, 0x6d, 0x03, 0x21, 0xd1, 0x32, 0xbc, 0x0e, 0x78,
	0x75, 0x1b, 0x6f, 0xf2, 0x96, 0x7f, 0xd8, 0xce, 0xa6, 0x7e, 0x26, 0xbe, 0x24, 0xb2, 0x9d, 0x2c,
	0x61, 0x82, 0xb7, 0x7b, 0xce, 0xbd, 0xf7, 0xd8, 0xe7, 0x38, 0x01, 0x28, 0x32, 0x9a, 0x2f, 0x6a,
	0x51, 0xa9, 0x0a, 0x03, 0x5d, 0xdf, 0x8a, 0x7a, 0x3f, 0xff, 0xe5, 0x42, 0x74, 0x4e, 0x15, 0xfd,
	0xda, 0x30, 0x71, 0x20, 0xec, 0x67, 0xc3, 0xa4, 0xc2, 0x04, 0xa6, 0xfb, 0xaa, 0xa8, 0x2b, 0xc9,
	0x15, 0xaf, 0xca, 0xd8, 0x49, 0xbc, 0x34, 0x24, 0x43, 0x0a, 0x4f, 0x21, 0x94, 0x2c, 0x67, 0x7b,
	0x55, 0x09, 0x19, 0xbb, 0x89, 0x97, 0x1e, 0x2d, 0x71, 0xd1, 0x89, 0x2e, 0xae, 0xda, 0x16, 0xe9,
	0x87, 0xf0, 0x23, 0x84, 0x0f, 0x54, 0x70, 0xba, 0xcb, 0x99, 0x8c, 0xbd, 0xc4, 0x4b, 0xa7, 0xcb,
	0x37, 0xfd, 0xc6, 0xf3, 0x2b, 0x2c, 0xbe, 0x75, 0xb3, 0x17, 0xa5, 0x12, 0x07, 0xd2, 0xef, 0x62,
	0x0a, 0xbe, 0xe2, 0x05, 0x8b, 0xfd, 0xc4, 0x49, 0xa7, 0xcb, 0xe3, 0x5e, 0x63, 0xcb, 0x0b, 0x76,
	0x49, 0x05, 0x2d, 0x24, 0x31, 0x13, 0xb3, 0x4b, 0x38, 0xfa, 0x53, 0x06, 0x23, 0xf0, 0xee, 0xd9,
	0x21, 0x76, 0x12, 0x27, 0x0d, 0x89, 0x2e, 0x31, 0x85, 0xd1, 0x03, 0xcd, 0x1b, 0x16, 0xbb, 0x46,
	0x6e, 0x60, 0xa2, 0x5b, 0x25, 0x76, 0xe0, 0xbd, 0xfb, 0xce, 0x99, 0xdf, 0x41, 0xd0, 0xd1, 0x88,
	0xe0, 0x97, 0xb4, 0x60, 0xad, 0x98, 0xa9, 0xf1, 0x35, 0x40, 0xc6, 0x6e, 0x78, 0x69, 0x73, 0x73,
	0x4d, 0x67, 0xc0, 0xe0, 0x31, 0x8c, 0x9a, 0x86, 0x67, 0x36, 0x80, 0xff, 0x89, 0x05, 0x86, 0x2d,
	0xb9, 0x92, 0xc6, 0x52, 0x48, 0x2c, 0x98, 0xdf, 0x00, 0xf4, 0x8e, 0xf4, 0x8c, 0x54, 0x54, 0xa8,
	0xf6, 0x38, 0x0b, 0xb4, 0x1f, 0x56, 0x66, 0xed, 0x41, 0xba, 0xc4, 0x97, 0x30, 0x7e, 0xe4, 0x65,
	0x56, 0x3d, 0xc6, 0x9e, 0x21, 0x5b, 0x84, 0x31, 0x4c, 0x68, 0xce, 0x6f, 0x4b, 0x96, 0x99, 0x53,
	0x02, 0xd2, 0xc1, 0xf9, 0x3d, 0xbc, 0x18, 0xa4, 0x2f, 0xeb, 0xaa, 0x94, 0x0c, 0x4f, 0x60, 0xa2,
	0x04, 0xaf, 0xf5, 0x5b, 0x39, 0xe6, 0xad, 0xa2, 0x41, 0xce, 0xa6, 0x41, 0xba, 0x81, 0xde, 0x94,
	0xfb, 0xcc, 0x14, 0x15, 0xc2, 0xdc, 0xc3, 0xb0, 0x06, 0xcc, 0xb7, 0x30, 0xb6, 0xeb, 0xfa, 0x42,
	0xb2, 0xd9, 0xdd, 0xb1, 0x7d, 0x67, 0xa9, 0x83, 0x38, 0x83, 0xa0, 0x16, 0x55, 0xcd, 0x84, 0x3a,
	0xb4, 0xce, 0x9e, 0xb0, 0x56, 0xb5, 0xcf, 0x65, 0xdd, 0x59, 0x70, 0x72, 0x06, 0x41, 0xf7, 0xc9,
	0xe1, 0x04, 0x3c, 0xb2, 0xfa, 0x1e, 0xfd, 0x87, 0x01, 0xf8, 0xeb, 0x8b, 0xd5, 0x26, 0x72, 0x34,
	0xb5, 0xfe, 0xb4, 0x89, 0x5c, 0x53, 0xac, 0x7e, 0x44, 0x1e, 0x86, 0x30, 0x3a, 0xfb, 0x72, 0xbd,
	0xd9, 0x46, 0xbe, 0xe6, 0xae, 0xae, 0xd7, 0xd1, 0x68, 0xb9, 0x01, 0x7f, 0x7d, 0xbe, 0xfa, 0x8c,
	0x1f, 0x20, 0x7c, 0xca, 0x03, 0x67, 0xff, 0xfe, 0x44, 0x67, 0xaf, 0xfe, 0xda, 0xb3, 0x01, 0x9e,
	0x3a, 0xbb, 0xb1, 0xf9, 0xd5, 0xde, 0xfe, 0x0e, 0x00, 0x00, 0xff, 0xff, 0xa9, 0x7e, 0x06, 0x92,
	0x78, 0x03, 0x00, 0x00,
}
