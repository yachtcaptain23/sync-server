// Code generated by protoc-gen-go. DO NOT EDIT.
// source: web_app_specifics.proto

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

// This enum should be a subset of the DisplayMode enum in
// chrome/browser/web_applications/proto/web_app.proto and
// third_party/blink/public/mojom/manifest/display_mode.mojom
type WebAppSpecifics_UserDisplayMode int32

const (
	// UNDEFINED is never serialized.
	WebAppSpecifics_BROWSER WebAppSpecifics_UserDisplayMode = 1
	// MINIMAL_UI is never serialized.
	WebAppSpecifics_STANDALONE WebAppSpecifics_UserDisplayMode = 3
)

var WebAppSpecifics_UserDisplayMode_name = map[int32]string{
	1: "BROWSER",
	3: "STANDALONE",
}

var WebAppSpecifics_UserDisplayMode_value = map[string]int32{
	"BROWSER":    1,
	"STANDALONE": 3,
}

func (x WebAppSpecifics_UserDisplayMode) Enum() *WebAppSpecifics_UserDisplayMode {
	p := new(WebAppSpecifics_UserDisplayMode)
	*p = x
	return p
}

func (x WebAppSpecifics_UserDisplayMode) String() string {
	return proto.EnumName(WebAppSpecifics_UserDisplayMode_name, int32(x))
}

func (x *WebAppSpecifics_UserDisplayMode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(WebAppSpecifics_UserDisplayMode_value, data, "WebAppSpecifics_UserDisplayMode")
	if err != nil {
		return err
	}
	*x = WebAppSpecifics_UserDisplayMode(value)
	return nil
}

func (WebAppSpecifics_UserDisplayMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e4fa27fdb53b8766, []int{0, 0}
}

// WebApp data. This is a synced part of
// chrome/browser/web_applications/proto/web_app.proto data.
type WebAppSpecifics struct {
	LaunchUrl            *string                          `protobuf:"bytes,1,opt,name=launch_url,json=launchUrl" json:"launch_url,omitempty"`
	Name                 *string                          `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	UserDisplayMode      *WebAppSpecifics_UserDisplayMode `protobuf:"varint,3,opt,name=user_display_mode,json=userDisplayMode,enum=sync_pb.WebAppSpecifics_UserDisplayMode" json:"user_display_mode,omitempty"`
	ThemeColor           *uint32                          `protobuf:"varint,4,opt,name=theme_color,json=themeColor" json:"theme_color,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *WebAppSpecifics) Reset()         { *m = WebAppSpecifics{} }
func (m *WebAppSpecifics) String() string { return proto.CompactTextString(m) }
func (*WebAppSpecifics) ProtoMessage()    {}
func (*WebAppSpecifics) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4fa27fdb53b8766, []int{0}
}

func (m *WebAppSpecifics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WebAppSpecifics.Unmarshal(m, b)
}
func (m *WebAppSpecifics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WebAppSpecifics.Marshal(b, m, deterministic)
}
func (m *WebAppSpecifics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WebAppSpecifics.Merge(m, src)
}
func (m *WebAppSpecifics) XXX_Size() int {
	return xxx_messageInfo_WebAppSpecifics.Size(m)
}
func (m *WebAppSpecifics) XXX_DiscardUnknown() {
	xxx_messageInfo_WebAppSpecifics.DiscardUnknown(m)
}

var xxx_messageInfo_WebAppSpecifics proto.InternalMessageInfo

func (m *WebAppSpecifics) GetLaunchUrl() string {
	if m != nil && m.LaunchUrl != nil {
		return *m.LaunchUrl
	}
	return ""
}

func (m *WebAppSpecifics) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *WebAppSpecifics) GetUserDisplayMode() WebAppSpecifics_UserDisplayMode {
	if m != nil && m.UserDisplayMode != nil {
		return *m.UserDisplayMode
	}
	return WebAppSpecifics_BROWSER
}

func (m *WebAppSpecifics) GetThemeColor() uint32 {
	if m != nil && m.ThemeColor != nil {
		return *m.ThemeColor
	}
	return 0
}

func init() {
	proto.RegisterEnum("sync_pb.WebAppSpecifics_UserDisplayMode", WebAppSpecifics_UserDisplayMode_name, WebAppSpecifics_UserDisplayMode_value)
	proto.RegisterType((*WebAppSpecifics)(nil), "sync_pb.WebAppSpecifics")
}

func init() {
	proto.RegisterFile("web_app_specifics.proto", fileDescriptor_e4fa27fdb53b8766)
}

var fileDescriptor_e4fa27fdb53b8766 = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x4f, 0x4b, 0xf3, 0x40,
	0x10, 0x87, 0xd9, 0x37, 0x85, 0xd2, 0x29, 0x6f, 0x53, 0xf7, 0x62, 0x2e, 0x62, 0x28, 0x08, 0x01,
	0x61, 0x0f, 0x7e, 0x83, 0xd4, 0x16, 0x3c, 0x68, 0x2b, 0x49, 0x4b, 0x8f, 0x4b, 0xb2, 0x19, 0x4d,
	0x60, 0x37, 0xb3, 0xec, 0x26, 0x48, 0x3f, 0xb3, 0x5f, 0x42, 0x4c, 0xf4, 0x60, 0x6f, 0xc3, 0x33,
	0x7f, 0x9e, 0xe1, 0x07, 0xd7, 0x1f, 0x58, 0xca, 0xc2, 0x5a, 0xe9, 0x2d, 0xaa, 0xe6, 0xad, 0x51,
	0x5e, 0x58, 0x47, 0x1d, 0xf1, 0xa9, 0x3f, 0xb7, 0x4a, 0xda, 0x72, 0xf5, 0xc9, 0x20, 0x3c, 0x61,
	0x99, 0x5a, 0x9b, 0xff, 0x8e, 0xf0, 0x1b, 0x00, 0x5d, 0xf4, 0xad, 0xaa, 0x65, 0xef, 0x74, 0xc4,
	0x62, 0x96, 0xcc, 0xb2, 0xd9, 0x48, 0x8e, 0x4e, 0x73, 0x0e, 0x93, 0xb6, 0x30, 0x18, 0xfd, 0x1b,
	0x1a, 0x43, 0xcd, 0x0f, 0x70, 0xd5, 0x7b, 0x74, 0xb2, 0x6a, 0xbc, 0xd5, 0xc5, 0x59, 0x1a, 0xaa,
	0x30, 0x0a, 0x62, 0x96, 0x2c, 0x1e, 0x12, 0xf1, 0xe3, 0x12, 0x17, 0x1e, 0x71, 0xf4, 0xe8, 0x36,
	0xe3, 0xc2, 0x0b, 0x55, 0x98, 0x85, 0xfd, 0x5f, 0xc0, 0x6f, 0x61, 0xde, 0xd5, 0x68, 0x50, 0x2a,
	0xd2, 0xe4, 0xa2, 0x49, 0xcc, 0x92, 0xff, 0x19, 0x0c, 0xe8, 0xf1, 0x9b, 0xac, 0x04, 0x84, 0x17,
	0x47, 0xf8, 0x1c, 0xa6, 0xeb, 0x6c, 0x7f, 0xca, 0xb7, 0xd9, 0x92, 0xf1, 0x05, 0x40, 0x7e, 0x48,
	0x77, 0x9b, 0xf4, 0x79, 0xbf, 0xdb, 0x2e, 0x83, 0xf5, 0x3d, 0xdc, 0x91, 0x7b, 0x17, 0xaa, 0x76,
	0x64, 0x9a, 0xde, 0x08, 0x45, 0xc6, 0x52, 0x8b, 0x6d, 0xe7, 0x87, 0x27, 0xc7, 0x70, 0x14, 0xe9,
	0xa7, 0xe0, 0x95, 0x7d, 0x05, 0x00, 0x00, 0xff, 0xff, 0x13, 0x8a, 0xb4, 0xa8, 0x3d, 0x01, 0x00,
	0x00,
}
