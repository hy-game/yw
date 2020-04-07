// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg_public.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 一些公共结构
type MsgInt struct {
	Value int32 `protobuf:"zigzag32,1,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgInt) Reset()                    { *m = MsgInt{} }
func (m *MsgInt) String() string            { return proto.CompactTextString(m) }
func (*MsgInt) ProtoMessage()               {}
func (*MsgInt) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *MsgInt) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type MsgUint struct {
	Value uint32 `protobuf:"varint,1,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgUint) Reset()                    { *m = MsgUint{} }
func (m *MsgUint) String() string            { return proto.CompactTextString(m) }
func (*MsgUint) ProtoMessage()               {}
func (*MsgUint) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *MsgUint) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type MsgBigInt struct {
	Value int64 `protobuf:"zigzag64,1,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgBigInt) Reset()                    { *m = MsgBigInt{} }
func (m *MsgBigInt) String() string            { return proto.CompactTextString(m) }
func (*MsgBigInt) ProtoMessage()               {}
func (*MsgBigInt) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *MsgBigInt) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type MsgStr struct {
	Value string `protobuf:"bytes,1,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgStr) Reset()                    { *m = MsgStr{} }
func (m *MsgStr) String() string            { return proto.CompactTextString(m) }
func (*MsgStr) ProtoMessage()               {}
func (*MsgStr) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *MsgStr) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type MsgBool struct {
	Value bool `protobuf:"varint,1,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgBool) Reset()                    { *m = MsgBool{} }
func (m *MsgBool) String() string            { return proto.CompactTextString(m) }
func (*MsgBool) ProtoMessage()               {}
func (*MsgBool) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *MsgBool) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

type MsgFloat struct {
	Value float32 `protobuf:"fixed32,1,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgFloat) Reset()                    { *m = MsgFloat{} }
func (m *MsgFloat) String() string            { return proto.CompactTextString(m) }
func (*MsgFloat) ProtoMessage()               {}
func (*MsgFloat) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *MsgFloat) GetValue() float32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type MsgKeyValue struct {
	Key   int32 `protobuf:"zigzag32,1,opt,name=Key,json=key" json:"Key,omitempty"`
	Value int32 `protobuf:"zigzag32,2,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgKeyValue) Reset()                    { *m = MsgKeyValue{} }
func (m *MsgKeyValue) String() string            { return proto.CompactTextString(m) }
func (*MsgKeyValue) ProtoMessage()               {}
func (*MsgKeyValue) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

func (m *MsgKeyValue) GetKey() int32 {
	if m != nil {
		return m.Key
	}
	return 0
}

func (m *MsgKeyValue) GetValue() int32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type MsgKeyValueU struct {
	Key   uint32 `protobuf:"varint,1,opt,name=Key,json=key" json:"Key,omitempty"`
	Value uint32 `protobuf:"varint,2,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgKeyValueU) Reset()                    { *m = MsgKeyValueU{} }
func (m *MsgKeyValueU) String() string            { return proto.CompactTextString(m) }
func (*MsgKeyValueU) ProtoMessage()               {}
func (*MsgKeyValueU) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

func (m *MsgKeyValueU) GetKey() uint32 {
	if m != nil {
		return m.Key
	}
	return 0
}

func (m *MsgKeyValueU) GetValue() uint32 {
	if m != nil {
		return m.Value
	}
	return 0
}

type MsgStrKeyValue struct {
	Key   string `protobuf:"bytes,1,opt,name=Key,json=key" json:"Key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgStrKeyValue) Reset()                    { *m = MsgStrKeyValue{} }
func (m *MsgStrKeyValue) String() string            { return proto.CompactTextString(m) }
func (*MsgStrKeyValue) ProtoMessage()               {}
func (*MsgStrKeyValue) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

func (m *MsgStrKeyValue) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *MsgStrKeyValue) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type MsgStrKeyValueArray struct {
	Key   string   `protobuf:"bytes,1,opt,name=Key,json=key" json:"Key,omitempty"`
	Value []string `protobuf:"bytes,2,rep,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgStrKeyValueArray) Reset()                    { *m = MsgStrKeyValueArray{} }
func (m *MsgStrKeyValueArray) String() string            { return proto.CompactTextString(m) }
func (*MsgStrKeyValueArray) ProtoMessage()               {}
func (*MsgStrKeyValueArray) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{9} }

func (m *MsgStrKeyValueArray) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *MsgStrKeyValueArray) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

type MsgBoolKeyValue struct {
	Key   int32 `protobuf:"varint,1,opt,name=Key,json=key" json:"Key,omitempty"`
	Value bool  `protobuf:"varint,2,opt,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgBoolKeyValue) Reset()                    { *m = MsgBoolKeyValue{} }
func (m *MsgBoolKeyValue) String() string            { return proto.CompactTextString(m) }
func (*MsgBoolKeyValue) ProtoMessage()               {}
func (*MsgBoolKeyValue) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{10} }

func (m *MsgBoolKeyValue) GetKey() int32 {
	if m != nil {
		return m.Key
	}
	return 0
}

func (m *MsgBoolKeyValue) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

type MsgIntArrary struct {
	Value []int32 `protobuf:"zigzag32,1,rep,packed,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgIntArrary) Reset()                    { *m = MsgIntArrary{} }
func (m *MsgIntArrary) String() string            { return proto.CompactTextString(m) }
func (*MsgIntArrary) ProtoMessage()               {}
func (*MsgIntArrary) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{11} }

func (m *MsgIntArrary) GetValue() []int32 {
	if m != nil {
		return m.Value
	}
	return nil
}

type MsgStringArray struct {
	Value []string `protobuf:"bytes,1,rep,name=Value,json=value" json:"Value,omitempty"`
}

func (m *MsgStringArray) Reset()                    { *m = MsgStringArray{} }
func (m *MsgStringArray) String() string            { return proto.CompactTextString(m) }
func (*MsgStringArray) ProtoMessage()               {}
func (*MsgStringArray) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{12} }

func (m *MsgStringArray) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgInt)(nil), "pb.MsgInt")
	proto.RegisterType((*MsgUint)(nil), "pb.MsgUint")
	proto.RegisterType((*MsgBigInt)(nil), "pb.MsgBigInt")
	proto.RegisterType((*MsgStr)(nil), "pb.MsgStr")
	proto.RegisterType((*MsgBool)(nil), "pb.MsgBool")
	proto.RegisterType((*MsgFloat)(nil), "pb.MsgFloat")
	proto.RegisterType((*MsgKeyValue)(nil), "pb.MsgKeyValue")
	proto.RegisterType((*MsgKeyValueU)(nil), "pb.MsgKeyValueU")
	proto.RegisterType((*MsgStrKeyValue)(nil), "pb.MsgStrKeyValue")
	proto.RegisterType((*MsgStrKeyValueArray)(nil), "pb.MsgStrKeyValueArray")
	proto.RegisterType((*MsgBoolKeyValue)(nil), "pb.MsgBoolKeyValue")
	proto.RegisterType((*MsgIntArrary)(nil), "pb.MsgIntArrary")
	proto.RegisterType((*MsgStringArray)(nil), "pb.MsgStringArray")
}

func init() { proto.RegisterFile("msg_public.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 262 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xc8, 0x2d, 0x4e, 0x8f,
	0x2f, 0x28, 0x4d, 0xca, 0xc9, 0x4c, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48,
	0x52, 0x92, 0xe3, 0x62, 0xf3, 0x2d, 0x4e, 0xf7, 0xcc, 0x2b, 0x11, 0x12, 0xe1, 0x62, 0x0d, 0x4b,
	0xcc, 0x29, 0x4d, 0x95, 0x60, 0x54, 0x60, 0xd4, 0x10, 0x0c, 0x62, 0x2d, 0x03, 0x71, 0x94, 0xe4,
	0xb9, 0xd8, 0x7d, 0x8b, 0xd3, 0x43, 0x33, 0xd1, 0x15, 0xf0, 0xc2, 0x14, 0x28, 0x72, 0x71, 0xfa,
	0x16, 0xa7, 0x3b, 0x65, 0x62, 0x9a, 0x21, 0x04, 0x53, 0x02, 0xb1, 0x23, 0xb8, 0xa4, 0x08, 0x55,
	0x9e, 0x13, 0xd5, 0x0e, 0xa7, 0xfc, 0xfc, 0x1c, 0x54, 0x05, 0x1c, 0x30, 0x05, 0x0a, 0x5c, 0x1c,
	0xbe, 0xc5, 0xe9, 0x6e, 0x39, 0xf9, 0x89, 0x68, 0x56, 0x30, 0xc1, 0x54, 0x98, 0x72, 0x71, 0xfb,
	0x16, 0xa7, 0x7b, 0xa7, 0x56, 0x82, 0xe5, 0x84, 0x04, 0xb8, 0x98, 0xbd, 0x53, 0x2b, 0xa1, 0x3e,
	0x61, 0xce, 0x4e, 0xad, 0x44, 0x68, 0x63, 0x42, 0xf6, 0x9d, 0x19, 0x17, 0x0f, 0x92, 0xb6, 0x50,
	0x64, 0x7d, 0xbc, 0x58, 0xf4, 0xc1, 0x3d, 0x6d, 0xc1, 0xc5, 0x07, 0xf1, 0x11, 0x36, 0x1b, 0x39,
	0xb1, 0xe8, 0x84, 0xfb, 0xd5, 0x96, 0x4b, 0x18, 0x55, 0xa7, 0x63, 0x51, 0x51, 0x62, 0x25, 0x7e,
	0xed, 0xcc, 0x08, 0xed, 0x96, 0x5c, 0xfc, 0xd0, 0xa0, 0xc2, 0x66, 0x33, 0x2b, 0x16, 0x9b, 0xe1,
	0x81, 0xa8, 0x02, 0xf6, 0xab, 0x67, 0x5e, 0x09, 0xc8, 0xc6, 0xa2, 0x4a, 0xe4, 0x80, 0x64, 0x46,
	0x84, 0x88, 0x1a, 0xcc, 0x67, 0x99, 0x79, 0xe9, 0x10, 0xa7, 0xa1, 0xa8, 0x83, 0x39, 0xc4, 0x89,
	0xc9, 0x83, 0x39, 0x89, 0x0d, 0x9c, 0x8c, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa7, 0xdf,
	0x13, 0xed, 0x5a, 0x02, 0x00, 0x00,
}
