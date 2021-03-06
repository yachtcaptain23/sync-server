// Code generated by protoc-gen-go. DO NOT EDIT.
// source: unique_position.proto

package sync_pb

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

// A UniquePosition is a string of bytes.
//
// Unique positions are unique per-item, since they are guaranteed to end with a
// fixed-length suffix that is unique per-item.  The position string may not end
// with a '\0' byte.
//
// Prior to the suffix is a series of arbitrary bytes of arbitrary length.
// Items under the same parent are positioned relative to each other by a
// lexicographic comparison of their UniquePosition values.
type UniquePosition struct {
	// The uncompressed string of bytes representing the position.
	//
	// Deprecated.  See history note above.
	Value []byte `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
	// The client may choose to write a compressed position to this field instead
	// of populating the 'value' above.  If it chooses to use compression, the
	// 'value' field above must be empty.  The position value will be compressed
	// with gzip and stored in the compressed_value field.  The position's
	// uncompressed length must be specified and written to the
	// uncompressed_length field.
	//
	// Deprecated.  See history note above.
	CompressedValue    []byte  `protobuf:"bytes,2,opt,name=compressed_value,json=compressedValue" json:"compressed_value,omitempty"`
	UncompressedLength *uint64 `protobuf:"varint,3,opt,name=uncompressed_length,json=uncompressedLength" json:"uncompressed_length,omitempty"`
	// This encoding uses compression scheme designed especially for unique
	// positions.  It has the property that X < Y precisely when Compressed(X) <
	// Compressed(Y), which is very useful when the most common operation is to
	// compare these positions against each other.  Their values may remain
	// compressed in memory.
	//
	// The compression scheme is implemented and documented in
	// sync/core_impl/base/unique_position.cc.
	//
	// As of M30, this is the preferred encoding.  Newer clients may continue to
	// populate the 'value' and 'compressed_value' fields to ensure backwards
	// compatibility, but they will always try to read from this field first.
	CustomCompressedV1   []byte   `protobuf:"bytes,4,opt,name=custom_compressed_v1,json=customCompressedV1" json:"custom_compressed_v1,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UniquePosition) Reset()         { *m = UniquePosition{} }
func (m *UniquePosition) String() string { return proto.CompactTextString(m) }
func (*UniquePosition) ProtoMessage()    {}
func (*UniquePosition) Descriptor() ([]byte, []int) {
	return fileDescriptor_a6e5f9f179b89638, []int{0}
}

func (m *UniquePosition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UniquePosition.Unmarshal(m, b)
}
func (m *UniquePosition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UniquePosition.Marshal(b, m, deterministic)
}
func (m *UniquePosition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UniquePosition.Merge(m, src)
}
func (m *UniquePosition) XXX_Size() int {
	return xxx_messageInfo_UniquePosition.Size(m)
}
func (m *UniquePosition) XXX_DiscardUnknown() {
	xxx_messageInfo_UniquePosition.DiscardUnknown(m)
}

var xxx_messageInfo_UniquePosition proto.InternalMessageInfo

func (m *UniquePosition) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *UniquePosition) GetCompressedValue() []byte {
	if m != nil {
		return m.CompressedValue
	}
	return nil
}

func (m *UniquePosition) GetUncompressedLength() uint64 {
	if m != nil && m.UncompressedLength != nil {
		return *m.UncompressedLength
	}
	return 0
}

func (m *UniquePosition) GetCustomCompressedV1() []byte {
	if m != nil {
		return m.CustomCompressedV1
	}
	return nil
}

func init() {
	proto.RegisterType((*UniquePosition)(nil), "sync_pb.UniquePosition")
}

func init() {
	proto.RegisterFile("unique_position.proto", fileDescriptor_a6e5f9f179b89638)
}

var fileDescriptor_a6e5f9f179b89638 = []byte{
	// 200 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0xcd, 0xcb, 0x2c,
	0x2c, 0x4d, 0x8d, 0x2f, 0xc8, 0x2f, 0xce, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0x62, 0x2f, 0xae, 0xcc, 0x4b, 0x8e, 0x2f, 0x48, 0x52, 0xda, 0xc2, 0xc8, 0xc5, 0x17,
	0x0a, 0x56, 0x12, 0x00, 0x55, 0x21, 0x24, 0xc2, 0xc5, 0x5a, 0x96, 0x98, 0x53, 0x9a, 0x2a, 0xc1,
	0xa8, 0xc0, 0xa8, 0xc1, 0x13, 0x04, 0xe1, 0x08, 0x69, 0x72, 0x09, 0x24, 0xe7, 0xe7, 0x16, 0x14,
	0xa5, 0x16, 0x17, 0xa7, 0xa6, 0xc4, 0x43, 0x14, 0x30, 0x81, 0x15, 0xf0, 0x23, 0xc4, 0xc3, 0xc0,
	0x4a, 0xf5, 0xb9, 0x84, 0x4b, 0xf3, 0x90, 0x14, 0xe7, 0xa4, 0xe6, 0xa5, 0x97, 0x64, 0x48, 0x30,
	0x2b, 0x30, 0x6a, 0xb0, 0x04, 0x09, 0x21, 0x4b, 0xf9, 0x80, 0x65, 0x84, 0x0c, 0xb8, 0x44, 0x92,
	0x4b, 0x8b, 0x4b, 0xf2, 0x73, 0xe3, 0x91, 0xad, 0x30, 0x94, 0x60, 0x01, 0x9b, 0x2f, 0x04, 0x91,
	0x73, 0x46, 0xd8, 0x62, 0xe8, 0xa4, 0xcd, 0xa5, 0x9a, 0x5f, 0x94, 0xae, 0x97, 0x9c, 0x51, 0x94,
	0x9f, 0x9b, 0x59, 0x9a, 0xab, 0x07, 0xd2, 0x97, 0x9f, 0x97, 0x9a, 0x57, 0x52, 0xac, 0x07, 0xf2,
	0x19, 0xc4, 0x97, 0xc9, 0xf9, 0x39, 0x1e, 0xcc, 0x01, 0x8c, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x57, 0xda, 0x62, 0x8e, 0x04, 0x01, 0x00, 0x00,
}
