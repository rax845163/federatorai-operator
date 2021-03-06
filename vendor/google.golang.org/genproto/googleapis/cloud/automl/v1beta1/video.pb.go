// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/automl/v1beta1/video.proto

package automl

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Dataset metadata specific to video classification.
// All Video Classification datasets are treated as multi label.
type VideoClassificationDatasetMetadata struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VideoClassificationDatasetMetadata) Reset()         { *m = VideoClassificationDatasetMetadata{} }
func (m *VideoClassificationDatasetMetadata) String() string { return proto.CompactTextString(m) }
func (*VideoClassificationDatasetMetadata) ProtoMessage()    {}
func (*VideoClassificationDatasetMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8c6b74d94d2916, []int{0}
}

func (m *VideoClassificationDatasetMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VideoClassificationDatasetMetadata.Unmarshal(m, b)
}
func (m *VideoClassificationDatasetMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VideoClassificationDatasetMetadata.Marshal(b, m, deterministic)
}
func (m *VideoClassificationDatasetMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VideoClassificationDatasetMetadata.Merge(m, src)
}
func (m *VideoClassificationDatasetMetadata) XXX_Size() int {
	return xxx_messageInfo_VideoClassificationDatasetMetadata.Size(m)
}
func (m *VideoClassificationDatasetMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_VideoClassificationDatasetMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_VideoClassificationDatasetMetadata proto.InternalMessageInfo

// Dataset metadata specific to video object tracking.
type VideoObjectTrackingDatasetMetadata struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VideoObjectTrackingDatasetMetadata) Reset()         { *m = VideoObjectTrackingDatasetMetadata{} }
func (m *VideoObjectTrackingDatasetMetadata) String() string { return proto.CompactTextString(m) }
func (*VideoObjectTrackingDatasetMetadata) ProtoMessage()    {}
func (*VideoObjectTrackingDatasetMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8c6b74d94d2916, []int{1}
}

func (m *VideoObjectTrackingDatasetMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VideoObjectTrackingDatasetMetadata.Unmarshal(m, b)
}
func (m *VideoObjectTrackingDatasetMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VideoObjectTrackingDatasetMetadata.Marshal(b, m, deterministic)
}
func (m *VideoObjectTrackingDatasetMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VideoObjectTrackingDatasetMetadata.Merge(m, src)
}
func (m *VideoObjectTrackingDatasetMetadata) XXX_Size() int {
	return xxx_messageInfo_VideoObjectTrackingDatasetMetadata.Size(m)
}
func (m *VideoObjectTrackingDatasetMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_VideoObjectTrackingDatasetMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_VideoObjectTrackingDatasetMetadata proto.InternalMessageInfo

// Model metadata specific to video classification.
type VideoClassificationModelMetadata struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VideoClassificationModelMetadata) Reset()         { *m = VideoClassificationModelMetadata{} }
func (m *VideoClassificationModelMetadata) String() string { return proto.CompactTextString(m) }
func (*VideoClassificationModelMetadata) ProtoMessage()    {}
func (*VideoClassificationModelMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8c6b74d94d2916, []int{2}
}

func (m *VideoClassificationModelMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VideoClassificationModelMetadata.Unmarshal(m, b)
}
func (m *VideoClassificationModelMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VideoClassificationModelMetadata.Marshal(b, m, deterministic)
}
func (m *VideoClassificationModelMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VideoClassificationModelMetadata.Merge(m, src)
}
func (m *VideoClassificationModelMetadata) XXX_Size() int {
	return xxx_messageInfo_VideoClassificationModelMetadata.Size(m)
}
func (m *VideoClassificationModelMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_VideoClassificationModelMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_VideoClassificationModelMetadata proto.InternalMessageInfo

// Model metadata specific to video object tracking.
type VideoObjectTrackingModelMetadata struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VideoObjectTrackingModelMetadata) Reset()         { *m = VideoObjectTrackingModelMetadata{} }
func (m *VideoObjectTrackingModelMetadata) String() string { return proto.CompactTextString(m) }
func (*VideoObjectTrackingModelMetadata) ProtoMessage()    {}
func (*VideoObjectTrackingModelMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8c6b74d94d2916, []int{3}
}

func (m *VideoObjectTrackingModelMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VideoObjectTrackingModelMetadata.Unmarshal(m, b)
}
func (m *VideoObjectTrackingModelMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VideoObjectTrackingModelMetadata.Marshal(b, m, deterministic)
}
func (m *VideoObjectTrackingModelMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VideoObjectTrackingModelMetadata.Merge(m, src)
}
func (m *VideoObjectTrackingModelMetadata) XXX_Size() int {
	return xxx_messageInfo_VideoObjectTrackingModelMetadata.Size(m)
}
func (m *VideoObjectTrackingModelMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_VideoObjectTrackingModelMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_VideoObjectTrackingModelMetadata proto.InternalMessageInfo

func init() {
	proto.RegisterType((*VideoClassificationDatasetMetadata)(nil), "google.cloud.automl.v1beta1.VideoClassificationDatasetMetadata")
	proto.RegisterType((*VideoObjectTrackingDatasetMetadata)(nil), "google.cloud.automl.v1beta1.VideoObjectTrackingDatasetMetadata")
	proto.RegisterType((*VideoClassificationModelMetadata)(nil), "google.cloud.automl.v1beta1.VideoClassificationModelMetadata")
	proto.RegisterType((*VideoObjectTrackingModelMetadata)(nil), "google.cloud.automl.v1beta1.VideoObjectTrackingModelMetadata")
}

func init() {
	proto.RegisterFile("google/cloud/automl/v1beta1/video.proto", fileDescriptor_ec8c6b74d94d2916)
}

var fileDescriptor_ec8c6b74d94d2916 = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0xd0, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x06, 0x60, 0x95, 0x81, 0xc1, 0x23, 0x63, 0x8b, 0x04, 0x8a, 0x90, 0xd8, 0x6c, 0x2a, 0x46,
	0xa6, 0xb6, 0x48, 0x4c, 0x11, 0x1d, 0x50, 0x07, 0x94, 0xe5, 0x62, 0x1f, 0x96, 0xc1, 0xf5, 0x45,
	0xf1, 0xa5, 0xcf, 0xc1, 0x73, 0xf1, 0x54, 0x28, 0x76, 0x84, 0x14, 0x1a, 0x65, 0xb4, 0xfe, 0xef,
	0x7e, 0xeb, 0x4e, 0xdc, 0x5b, 0x22, 0xeb, 0x51, 0x69, 0x4f, 0x9d, 0x51, 0xd0, 0x31, 0x1d, 0xbd,
	0x3a, 0xad, 0x6b, 0x64, 0x58, 0xab, 0x93, 0x33, 0x48, 0xb2, 0x69, 0x89, 0xe9, 0x6a, 0x95, 0xa1,
	0x4c, 0x50, 0x66, 0x28, 0x07, 0xb8, 0x7c, 0x98, 0x6b, 0xd1, 0x1e, 0x62, 0x74, 0x1f, 0x4e, 0x03,
	0x3b, 0x0a, 0xb9, 0x6e, 0x79, 0x3d, 0x4c, 0x40, 0xe3, 0x14, 0x84, 0x40, 0x9c, 0xc2, 0x98, 0xd3,
	0xe2, 0x4e, 0x14, 0x87, 0xfe, 0xef, 0xdd, 0x68, 0xf4, 0x19, 0x18, 0x22, 0x72, 0x89, 0x0c, 0x06,
	0x18, 0xfe, 0xd4, 0x6b, 0xfd, 0x89, 0x9a, 0xdf, 0x5a, 0xd0, 0x5f, 0x2e, 0xd8, 0xff, 0xaa, 0x10,
	0xb7, 0x13, 0x5d, 0x25, 0x19, 0xf4, 0x67, 0x66, 0xdc, 0x34, 0x32, 0xdb, 0xef, 0x85, 0xb8, 0xd1,
	0x74, 0x94, 0x33, 0x77, 0xd8, 0x8a, 0xd4, 0xb2, 0xef, 0x77, 0xd8, 0x2f, 0xde, 0x37, 0x03, 0xb5,
	0xe4, 0x21, 0x58, 0x49, 0xad, 0x55, 0x16, 0x43, 0xda, 0x50, 0xe5, 0x08, 0x1a, 0x17, 0x27, 0x8f,
	0xf6, 0x94, 0x9f, 0x3f, 0x17, 0xab, 0x97, 0x04, 0xab, 0x5d, 0x8f, 0xaa, 0x4d, 0xc7, 0x54, 0xfa,
	0xea, 0x90, 0x51, 0x7d, 0x99, 0xba, 0x1e, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x1a, 0xeb, 0xe8,
	0x65, 0xc5, 0x01, 0x00, 0x00,
}
