// Code generated by protoc-gen-go. DO NOT EDIT.
// source: device_info_specifics.proto

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

// Enum defining available Sharing features.
type SharingSpecificFields_EnabledFeatures int32

const (
	SharingSpecificFields_UNKNOWN          SharingSpecificFields_EnabledFeatures = 0
	SharingSpecificFields_CLICK_TO_CALL    SharingSpecificFields_EnabledFeatures = 1
	SharingSpecificFields_SHARED_CLIPBOARD SharingSpecificFields_EnabledFeatures = 2
	SharingSpecificFields_SMS_FETCHER      SharingSpecificFields_EnabledFeatures = 3
	SharingSpecificFields_REMOTE_COPY      SharingSpecificFields_EnabledFeatures = 4
	SharingSpecificFields_PEER_CONNECTION  SharingSpecificFields_EnabledFeatures = 5
)

var SharingSpecificFields_EnabledFeatures_name = map[int32]string{
	0: "UNKNOWN",
	1: "CLICK_TO_CALL",
	2: "SHARED_CLIPBOARD",
	3: "SMS_FETCHER",
	4: "REMOTE_COPY",
	5: "PEER_CONNECTION",
}

var SharingSpecificFields_EnabledFeatures_value = map[string]int32{
	"UNKNOWN":          0,
	"CLICK_TO_CALL":    1,
	"SHARED_CLIPBOARD": 2,
	"SMS_FETCHER":      3,
	"REMOTE_COPY":      4,
	"PEER_CONNECTION":  5,
}

func (x SharingSpecificFields_EnabledFeatures) Enum() *SharingSpecificFields_EnabledFeatures {
	p := new(SharingSpecificFields_EnabledFeatures)
	*p = x
	return p
}

func (x SharingSpecificFields_EnabledFeatures) String() string {
	return proto.EnumName(SharingSpecificFields_EnabledFeatures_name, int32(x))
}

func (x *SharingSpecificFields_EnabledFeatures) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SharingSpecificFields_EnabledFeatures_value, data, "SharingSpecificFields_EnabledFeatures")
	if err != nil {
		return err
	}
	*x = SharingSpecificFields_EnabledFeatures(value)
	return nil
}

func (SharingSpecificFields_EnabledFeatures) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_199d98dfb54dc818, []int{2, 0}
}

// Information about a device that is running a sync-enabled Chrome browser.  We
// are mapping the per-client cache guid to more specific information about the
// device.
type DeviceInfoSpecifics struct {
	// The cache_guid created to identify a sync client on this device.
	CacheGuid *string `protobuf:"bytes,1,opt,name=cache_guid,json=cacheGuid" json:"cache_guid,omitempty"`
	// A non-unique but human readable name to describe this client.
	ClientName *string `protobuf:"bytes,2,opt,name=client_name,json=clientName" json:"client_name,omitempty"`
	// The platform of the device.
	DeviceType *SyncEnums_DeviceType `protobuf:"varint,3,opt,name=device_type,json=deviceType,enum=sync_pb.SyncEnums_DeviceType" json:"device_type,omitempty"`
	// The UserAgent used when contacting the Chrome Sync server.
	SyncUserAgent *string `protobuf:"bytes,4,opt,name=sync_user_agent,json=syncUserAgent" json:"sync_user_agent,omitempty"`
	// The Chrome instance's version.  Updated (if necessary) on every startup.
	ChromeVersion *string `protobuf:"bytes,5,opt,name=chrome_version,json=chromeVersion" json:"chrome_version,omitempty"`
	// Last time when pre-sync data on the device was saved. The device can be
	// restored to state back to this time. In millisecond since UNIX epoch.
	// DEPRECATED in M50.
	DeprecatedBackupTimestamp *int64 `protobuf:"varint,6,opt,name=deprecated_backup_timestamp,json=deprecatedBackupTimestamp" json:"deprecated_backup_timestamp,omitempty"` // Deprecated: Do not use.
	// Device_id that is stable until user signs out. This device_id is used for
	// annotating login scoped refresh token.
	SigninScopedDeviceId *string `protobuf:"bytes,7,opt,name=signin_scoped_device_id,json=signinScopedDeviceId" json:"signin_scoped_device_id,omitempty"`
	// This field is updated to be the current time periodically, and is also set
	// to the current time whenever another field changes. By examining the
	// difference between this field and the current time, it should be possible
	// to reason about the inactivity of any device that was syncing at one time.
	LastUpdatedTimestamp *int64 `protobuf:"varint,8,opt,name=last_updated_timestamp,json=lastUpdatedTimestamp" json:"last_updated_timestamp,omitempty"`
	// Device info fields that are specific to a feature. This is information that
	// can not be derived from the other fields in the proto and are not general
	// enough to be used by another feature.
	FeatureFields *FeatureSpecificFields `protobuf:"bytes,9,opt,name=feature_fields,json=featureFields" json:"feature_fields,omitempty"`
	// Device specific information for Sharing feature.
	SharingFields *SharingSpecificFields `protobuf:"bytes,10,opt,name=sharing_fields,json=sharingFields" json:"sharing_fields,omitempty"`
	// Model of device.
	Model *string `protobuf:"bytes,11,opt,name=model" json:"model,omitempty"`
	// Name of device manufacturer.
	Manufacturer         *string  `protobuf:"bytes,12,opt,name=manufacturer" json:"manufacturer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeviceInfoSpecifics) Reset()         { *m = DeviceInfoSpecifics{} }
func (m *DeviceInfoSpecifics) String() string { return proto.CompactTextString(m) }
func (*DeviceInfoSpecifics) ProtoMessage()    {}
func (*DeviceInfoSpecifics) Descriptor() ([]byte, []int) {
	return fileDescriptor_199d98dfb54dc818, []int{0}
}

func (m *DeviceInfoSpecifics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeviceInfoSpecifics.Unmarshal(m, b)
}
func (m *DeviceInfoSpecifics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeviceInfoSpecifics.Marshal(b, m, deterministic)
}
func (m *DeviceInfoSpecifics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeviceInfoSpecifics.Merge(m, src)
}
func (m *DeviceInfoSpecifics) XXX_Size() int {
	return xxx_messageInfo_DeviceInfoSpecifics.Size(m)
}
func (m *DeviceInfoSpecifics) XXX_DiscardUnknown() {
	xxx_messageInfo_DeviceInfoSpecifics.DiscardUnknown(m)
}

var xxx_messageInfo_DeviceInfoSpecifics proto.InternalMessageInfo

func (m *DeviceInfoSpecifics) GetCacheGuid() string {
	if m != nil && m.CacheGuid != nil {
		return *m.CacheGuid
	}
	return ""
}

func (m *DeviceInfoSpecifics) GetClientName() string {
	if m != nil && m.ClientName != nil {
		return *m.ClientName
	}
	return ""
}

func (m *DeviceInfoSpecifics) GetDeviceType() SyncEnums_DeviceType {
	if m != nil && m.DeviceType != nil {
		return *m.DeviceType
	}
	return SyncEnums_TYPE_UNSET
}

func (m *DeviceInfoSpecifics) GetSyncUserAgent() string {
	if m != nil && m.SyncUserAgent != nil {
		return *m.SyncUserAgent
	}
	return ""
}

func (m *DeviceInfoSpecifics) GetChromeVersion() string {
	if m != nil && m.ChromeVersion != nil {
		return *m.ChromeVersion
	}
	return ""
}

// Deprecated: Do not use.
func (m *DeviceInfoSpecifics) GetDeprecatedBackupTimestamp() int64 {
	if m != nil && m.DeprecatedBackupTimestamp != nil {
		return *m.DeprecatedBackupTimestamp
	}
	return 0
}

func (m *DeviceInfoSpecifics) GetSigninScopedDeviceId() string {
	if m != nil && m.SigninScopedDeviceId != nil {
		return *m.SigninScopedDeviceId
	}
	return ""
}

func (m *DeviceInfoSpecifics) GetLastUpdatedTimestamp() int64 {
	if m != nil && m.LastUpdatedTimestamp != nil {
		return *m.LastUpdatedTimestamp
	}
	return 0
}

func (m *DeviceInfoSpecifics) GetFeatureFields() *FeatureSpecificFields {
	if m != nil {
		return m.FeatureFields
	}
	return nil
}

func (m *DeviceInfoSpecifics) GetSharingFields() *SharingSpecificFields {
	if m != nil {
		return m.SharingFields
	}
	return nil
}

func (m *DeviceInfoSpecifics) GetModel() string {
	if m != nil && m.Model != nil {
		return *m.Model
	}
	return ""
}

func (m *DeviceInfoSpecifics) GetManufacturer() string {
	if m != nil && m.Manufacturer != nil {
		return *m.Manufacturer
	}
	return ""
}

// Feature specific information about the device that is running a sync-enabled
// Chrome browser. Adding to this proto is discouraged and should only be added
// when the information can not be derived more generally.
type FeatureSpecificFields struct {
	// Tracks whether the SendTabToSelf feature is enabled on the device. For this
	// to be true, two things must be true: (1) The receiving side of the feature
	// must be enabled on the device (2) The user has enabled sync for this
	// feature
	SendTabToSelfReceivingEnabled *bool    `protobuf:"varint,1,opt,name=send_tab_to_self_receiving_enabled,json=sendTabToSelfReceivingEnabled" json:"send_tab_to_self_receiving_enabled,omitempty"`
	XXX_NoUnkeyedLiteral          struct{} `json:"-"`
	XXX_unrecognized              []byte   `json:"-"`
	XXX_sizecache                 int32    `json:"-"`
}

func (m *FeatureSpecificFields) Reset()         { *m = FeatureSpecificFields{} }
func (m *FeatureSpecificFields) String() string { return proto.CompactTextString(m) }
func (*FeatureSpecificFields) ProtoMessage()    {}
func (*FeatureSpecificFields) Descriptor() ([]byte, []int) {
	return fileDescriptor_199d98dfb54dc818, []int{1}
}

func (m *FeatureSpecificFields) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureSpecificFields.Unmarshal(m, b)
}
func (m *FeatureSpecificFields) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureSpecificFields.Marshal(b, m, deterministic)
}
func (m *FeatureSpecificFields) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureSpecificFields.Merge(m, src)
}
func (m *FeatureSpecificFields) XXX_Size() int {
	return xxx_messageInfo_FeatureSpecificFields.Size(m)
}
func (m *FeatureSpecificFields) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureSpecificFields.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureSpecificFields proto.InternalMessageInfo

func (m *FeatureSpecificFields) GetSendTabToSelfReceivingEnabled() bool {
	if m != nil && m.SendTabToSelfReceivingEnabled != nil {
		return *m.SendTabToSelfReceivingEnabled
	}
	return false
}

// Device specific information for Sharing feature. Used to send end-to-end
// encrypted message through FCM to other devices.
type SharingSpecificFields struct {
	// FCM registration token of device subscribed using VAPID key.
	// TODO(crbug.com/1012226): Deprecate when VAPID migration is over.
	VapidFcmToken *string `protobuf:"bytes,1,opt,name=vapid_fcm_token,json=vapidFcmToken" json:"vapid_fcm_token,omitempty"`
	// Public key for message encryption [RFC8291] using VAPID key.
	// TODO(crbug.com/1012226): Deprecate when VAPID migration is over.
	VapidP256Dh []byte `protobuf:"bytes,2,opt,name=vapid_p256dh,json=vapidP256dh" json:"vapid_p256dh,omitempty"`
	// Auth secret for message encryption [RFC8291] using VAPID key.
	// TODO(crbug.com/1012226): Deprecate when VAPID migration is over.
	VapidAuthSecret []byte `protobuf:"bytes,3,opt,name=vapid_auth_secret,json=vapidAuthSecret" json:"vapid_auth_secret,omitempty"`
	// A list of enabled Sharing features.
	EnabledFeatures []SharingSpecificFields_EnabledFeatures `protobuf:"varint,4,rep,name=enabled_features,json=enabledFeatures,enum=sync_pb.SharingSpecificFields_EnabledFeatures" json:"enabled_features,omitempty"`
	// FCM registration token of device subscribed using Sharing sender ID.
	SenderIdFcmToken *string `protobuf:"bytes,5,opt,name=sender_id_fcm_token,json=senderIdFcmToken" json:"sender_id_fcm_token,omitempty"`
	// Public key for message encryption [RFC8291] using Sharing sender ID.
	SenderIdP256Dh []byte `protobuf:"bytes,6,opt,name=sender_id_p256dh,json=senderIdP256dh" json:"sender_id_p256dh,omitempty"`
	// Auth secret for message encryption [RFC8291] using Sharing sender ID.
	SenderIdAuthSecret   []byte   `protobuf:"bytes,7,opt,name=sender_id_auth_secret,json=senderIdAuthSecret" json:"sender_id_auth_secret,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SharingSpecificFields) Reset()         { *m = SharingSpecificFields{} }
func (m *SharingSpecificFields) String() string { return proto.CompactTextString(m) }
func (*SharingSpecificFields) ProtoMessage()    {}
func (*SharingSpecificFields) Descriptor() ([]byte, []int) {
	return fileDescriptor_199d98dfb54dc818, []int{2}
}

func (m *SharingSpecificFields) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SharingSpecificFields.Unmarshal(m, b)
}
func (m *SharingSpecificFields) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SharingSpecificFields.Marshal(b, m, deterministic)
}
func (m *SharingSpecificFields) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SharingSpecificFields.Merge(m, src)
}
func (m *SharingSpecificFields) XXX_Size() int {
	return xxx_messageInfo_SharingSpecificFields.Size(m)
}
func (m *SharingSpecificFields) XXX_DiscardUnknown() {
	xxx_messageInfo_SharingSpecificFields.DiscardUnknown(m)
}

var xxx_messageInfo_SharingSpecificFields proto.InternalMessageInfo

func (m *SharingSpecificFields) GetVapidFcmToken() string {
	if m != nil && m.VapidFcmToken != nil {
		return *m.VapidFcmToken
	}
	return ""
}

func (m *SharingSpecificFields) GetVapidP256Dh() []byte {
	if m != nil {
		return m.VapidP256Dh
	}
	return nil
}

func (m *SharingSpecificFields) GetVapidAuthSecret() []byte {
	if m != nil {
		return m.VapidAuthSecret
	}
	return nil
}

func (m *SharingSpecificFields) GetEnabledFeatures() []SharingSpecificFields_EnabledFeatures {
	if m != nil {
		return m.EnabledFeatures
	}
	return nil
}

func (m *SharingSpecificFields) GetSenderIdFcmToken() string {
	if m != nil && m.SenderIdFcmToken != nil {
		return *m.SenderIdFcmToken
	}
	return ""
}

func (m *SharingSpecificFields) GetSenderIdP256Dh() []byte {
	if m != nil {
		return m.SenderIdP256Dh
	}
	return nil
}

func (m *SharingSpecificFields) GetSenderIdAuthSecret() []byte {
	if m != nil {
		return m.SenderIdAuthSecret
	}
	return nil
}

func init() {
	proto.RegisterEnum("sync_pb.SharingSpecificFields_EnabledFeatures", SharingSpecificFields_EnabledFeatures_name, SharingSpecificFields_EnabledFeatures_value)
	proto.RegisterType((*DeviceInfoSpecifics)(nil), "sync_pb.DeviceInfoSpecifics")
	proto.RegisterType((*FeatureSpecificFields)(nil), "sync_pb.FeatureSpecificFields")
	proto.RegisterType((*SharingSpecificFields)(nil), "sync_pb.SharingSpecificFields")
}

func init() {
	proto.RegisterFile("device_info_specifics.proto", fileDescriptor_199d98dfb54dc818)
}

var fileDescriptor_199d98dfb54dc818 = []byte{
	// 739 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0x51, 0x6f, 0xe3, 0x44,
	0x10, 0x26, 0x4d, 0x7b, 0xbd, 0x4e, 0xd2, 0xc4, 0xb7, 0x6d, 0xc1, 0x80, 0x0a, 0x21, 0xd2, 0xa1,
	0x08, 0x44, 0x24, 0x2a, 0x8e, 0x47, 0xa4, 0xc4, 0x75, 0x69, 0x74, 0xbd, 0x24, 0xb2, 0x5d, 0xd0,
	0x3d, 0xad, 0x36, 0xeb, 0x71, 0xb2, 0x3a, 0x7b, 0x6d, 0x79, 0xd7, 0x95, 0xfa, 0xc2, 0x3f, 0xe5,
	0x67, 0xf0, 0x8e, 0xbc, 0x6b, 0x5f, 0x68, 0x75, 0xe2, 0xcd, 0xfe, 0xbe, 0x6f, 0xe6, 0x9b, 0x6f,
	0x3c, 0x09, 0x7c, 0x1d, 0xe3, 0x83, 0xe0, 0x48, 0x85, 0x4c, 0x72, 0xaa, 0x0a, 0xe4, 0x22, 0x11,
	0x5c, 0x4d, 0x8b, 0x32, 0xd7, 0x39, 0x39, 0x56, 0x8f, 0x92, 0xd3, 0x62, 0xf3, 0x95, 0x63, 0x1e,
	0x50, 0x56, 0x59, 0x43, 0x8d, 0xff, 0x3e, 0x84, 0xb3, 0x6b, 0x53, 0xba, 0x90, 0x49, 0x1e, 0xb6,
	0x85, 0xe4, 0x12, 0x80, 0x33, 0xbe, 0x43, 0xba, 0xad, 0x44, 0xec, 0x76, 0x46, 0x9d, 0xc9, 0x49,
	0x70, 0x62, 0x90, 0xdf, 0x2b, 0x11, 0x93, 0x6f, 0xa1, 0xc7, 0x53, 0x81, 0x52, 0x53, 0xc9, 0x32,
	0x74, 0x0f, 0x0c, 0x0f, 0x16, 0x5a, 0xb2, 0x0c, 0xc9, 0x6f, 0xd0, 0x6b, 0x26, 0xd2, 0x8f, 0x05,
	0xba, 0xdd, 0x51, 0x67, 0x32, 0xb8, 0xba, 0x9c, 0x36, 0x83, 0x4c, 0xc3, 0x47, 0xc9, 0x7d, 0x33,
	0x86, 0x35, 0x8f, 0x1e, 0x0b, 0x0c, 0x20, 0xfe, 0xf8, 0x4c, 0xbe, 0x87, 0xa1, 0xd1, 0x56, 0x0a,
	0x4b, 0xca, 0xb6, 0x28, 0xb5, 0x7b, 0x68, 0x4c, 0x4e, 0x6b, 0xf8, 0x5e, 0x61, 0x39, 0xab, 0x41,
	0xf2, 0x1a, 0x06, 0x7c, 0x57, 0xe6, 0x19, 0xd2, 0x07, 0x2c, 0x95, 0xc8, 0xa5, 0x7b, 0x64, 0x65,
	0x16, 0xfd, 0xc3, 0x82, 0x64, 0x5e, 0x2f, 0xa8, 0x28, 0x91, 0x33, 0x8d, 0x31, 0xdd, 0x30, 0xfe,
	0xa1, 0x2a, 0xa8, 0x16, 0x19, 0x2a, 0xcd, 0xb2, 0xc2, 0x7d, 0x31, 0xea, 0x4c, 0xba, 0xf3, 0x03,
	0xb7, 0x13, 0x7c, 0xb9, 0x97, 0xcd, 0x8d, 0x2a, 0x6a, 0x45, 0xe4, 0x0d, 0x7c, 0xa1, 0xc4, 0x56,
	0x0a, 0x49, 0x15, 0xcf, 0x0b, 0x8c, 0x69, 0xbb, 0xf2, 0xd8, 0x3d, 0x36, 0x9e, 0xe7, 0x96, 0x0e,
	0x0d, 0xdb, 0x2c, 0x35, 0x26, 0xbf, 0xc0, 0xe7, 0x29, 0x53, 0x9a, 0x56, 0x45, 0x6c, 0xcc, 0xf7,
	0xae, 0x2f, 0x6b, 0xd7, 0xe0, 0xbc, 0x66, 0xef, 0x2d, 0xb9, 0x37, 0xf3, 0x61, 0x90, 0x20, 0xd3,
	0x55, 0x89, 0x34, 0x11, 0x98, 0xc6, 0xca, 0x3d, 0x19, 0x75, 0x26, 0xbd, 0xab, 0x6f, 0x3e, 0xae,
	0xf0, 0xc6, 0xd2, 0xed, 0x27, 0xbb, 0x31, 0xaa, 0xe0, 0xb4, 0xa9, 0xb2, 0xaf, 0x75, 0x1b, 0xb5,
	0x63, 0xa5, 0x90, 0xdb, 0xb6, 0x0d, 0x3c, 0x6b, 0x13, 0x5a, 0xfa, 0x79, 0x9b, 0xa6, 0xaa, 0x69,
	0x73, 0x0e, 0x47, 0x59, 0x1e, 0x63, 0xea, 0xf6, 0x4c, 0x50, 0xfb, 0x42, 0xc6, 0xd0, 0xcf, 0x98,
	0xac, 0x12, 0xc6, 0x6b, 0xc7, 0xd2, 0xed, 0x1b, 0xf2, 0x09, 0x36, 0xde, 0xc0, 0xc5, 0x27, 0x07,
	0x25, 0x0b, 0x18, 0x2b, 0x94, 0x31, 0xd5, 0x6c, 0x43, 0x75, 0x4e, 0x15, 0xa6, 0x09, 0x2d, 0x91,
	0xa3, 0x78, 0xa8, 0x87, 0x45, 0xc9, 0x36, 0x29, 0xda, 0xc3, 0x7b, 0x19, 0x5c, 0xd6, 0xca, 0x88,
	0x6d, 0xa2, 0x3c, 0xc4, 0x34, 0x09, 0x5a, 0x95, 0x6f, 0x45, 0xe3, 0x7f, 0xba, 0x70, 0xf1, 0xc9,
	0x18, 0xf5, 0x15, 0x3d, 0xb0, 0x42, 0xc4, 0x34, 0xe1, 0x19, 0xd5, 0xf9, 0x07, 0x94, 0xcd, 0x29,
	0x9f, 0x1a, 0xf8, 0x86, 0x67, 0x51, 0x0d, 0x92, 0xef, 0xa0, 0x6f, 0x75, 0xc5, 0xd5, 0x9b, 0x5f,
	0xe3, 0x9d, 0xb9, 0xe7, 0x7e, 0xd0, 0x33, 0xd8, 0xda, 0x40, 0xe4, 0x07, 0x78, 0x65, 0x25, 0xac,
	0xd2, 0x3b, 0xaa, 0x90, 0x97, 0xa8, 0xcd, 0x59, 0xf7, 0x03, 0xeb, 0x31, 0xab, 0xf4, 0x2e, 0x34,
	0x30, 0x79, 0x0f, 0x4e, 0x13, 0x80, 0x36, 0x9f, 0x43, 0xb9, 0x87, 0xa3, 0xee, 0x64, 0x70, 0x35,
	0xfd, 0xff, 0xbd, 0x4f, 0x9b, 0x48, 0xcd, 0xca, 0x54, 0x30, 0xc4, 0xa7, 0x00, 0xf9, 0x09, 0xce,
	0xea, 0x65, 0x60, 0x49, 0x9f, 0xa4, 0xb2, 0x47, 0xef, 0x58, 0x6a, 0xb1, 0x0f, 0x36, 0x01, 0x67,
	0x2f, 0x6f, 0xc2, 0xbd, 0x30, 0x43, 0x0f, 0x5a, 0x6d, 0x93, 0xef, 0x67, 0xb8, 0xd8, 0x2b, 0xff,
	0x9b, 0xf1, 0xd8, 0xc8, 0x49, 0x2b, 0xdf, 0xc7, 0x1c, 0xff, 0x05, 0xc3, 0x67, 0xf3, 0x92, 0x1e,
	0x1c, 0xdf, 0x2f, 0xdf, 0x2e, 0x57, 0x7f, 0x2e, 0x9d, 0xcf, 0xc8, 0x2b, 0x38, 0xf5, 0xee, 0x16,
	0xde, 0x5b, 0x1a, 0xad, 0xa8, 0x37, 0xbb, 0xbb, 0x73, 0x3a, 0xe4, 0x1c, 0x9c, 0xf0, 0x76, 0x16,
	0xf8, 0xd7, 0xd4, 0xbb, 0x5b, 0xac, 0xe7, 0xab, 0x59, 0x70, 0xed, 0x1c, 0x90, 0x21, 0xf4, 0xc2,
	0x77, 0x21, 0xbd, 0xf1, 0x23, 0xef, 0xd6, 0x0f, 0x9c, 0x6e, 0x0d, 0x04, 0xfe, 0xbb, 0x55, 0xe4,
	0x53, 0x6f, 0xb5, 0x7e, 0xef, 0x1c, 0x92, 0x33, 0x18, 0xae, 0x7d, 0x3f, 0xa0, 0xde, 0x6a, 0xb9,
	0xf4, 0xbd, 0x68, 0xb1, 0x5a, 0x3a, 0x47, 0xf3, 0x1f, 0xe1, 0x75, 0x5e, 0x6e, 0xa7, 0xe6, 0x97,
	0x2e, 0xaa, 0x6c, 0xca, 0xf3, 0xac, 0xc8, 0x25, 0x4a, 0xad, 0xcc, 0x96, 0xed, 0x3f, 0x1c, 0xcf,
	0xd3, 0xdb, 0xee, 0xba, 0xf3, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x29, 0xd2, 0xdc, 0xd9, 0x21,
	0x05, 0x00, 0x00,
}