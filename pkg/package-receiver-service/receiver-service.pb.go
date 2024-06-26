// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: receiver-service.proto

package package_receiver

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetTPRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tp  string `protobuf:"bytes,1,opt,name=tp,proto3" json:"tp,omitempty"`
	Doc string `protobuf:"bytes,2,opt,name=doc,proto3" json:"doc,omitempty"`
}

func (x *GetTPRequest) Reset() {
	*x = GetTPRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_receiver_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTPRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTPRequest) ProtoMessage() {}

func (x *GetTPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTPRequest.ProtoReflect.Descriptor instead.
func (*GetTPRequest) Descriptor() ([]byte, []int) {
	return file_receiver_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetTPRequest) GetTp() string {
	if x != nil {
		return x.Tp
	}
	return ""
}

func (x *GetTPRequest) GetDoc() string {
	if x != nil {
		return x.Doc
	}
	return ""
}

type GetTPResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Founded           bool         `protobuf:"varint,1,opt,name=founded,proto3" json:"founded,omitempty"`
	Tp                string       `protobuf:"bytes,2,opt,name=tp,proto3" json:"tp,omitempty"`
	Origin            string       `protobuf:"bytes,3,opt,name=origin,proto3" json:"origin,omitempty"`
	IsReceipt         bool         `protobuf:"varint,4,opt,name=isReceipt,proto3" json:"isReceipt,omitempty"`
	ReceiptUrl        string       `protobuf:"bytes,5,opt,name=receipt_url,json=receiptUrl,proto3" json:"receipt_url,omitempty"`
	IsSuccess         bool         `protobuf:"varint,6,opt,name=isSuccess,proto3" json:"isSuccess,omitempty"`
	IsValidationError bool         `protobuf:"varint,7,opt,name=isValidationError,proto3" json:"isValidationError,omitempty"`
	IsInternalError   bool         `protobuf:"varint,8,opt,name=isInternalError,proto3" json:"isInternalError,omitempty"`
	IsNew             bool         `protobuf:"varint,9,opt,name=isNew,proto3" json:"isNew,omitempty"`
	CreatedAt         string       `protobuf:"bytes,10,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	SendTaskNextAt    string       `protobuf:"bytes,11,opt,name=sendTaskNextAt,proto3" json:"sendTaskNextAt,omitempty"`
	Content           []*Directory `protobuf:"bytes,12,rep,name=content,proto3" json:"content,omitempty"`
	ErrorText         string       `protobuf:"bytes,13,opt,name=errorText,proto3" json:"errorText,omitempty"`
	ErrorCode         string       `protobuf:"bytes,14,opt,name=errorCode,proto3" json:"errorCode,omitempty"`
	TimeLayout        string       `protobuf:"bytes,15,opt,name=timeLayout,proto3" json:"timeLayout,omitempty"`
}

func (x *GetTPResponse) Reset() {
	*x = GetTPResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_receiver_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTPResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTPResponse) ProtoMessage() {}

func (x *GetTPResponse) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTPResponse.ProtoReflect.Descriptor instead.
func (*GetTPResponse) Descriptor() ([]byte, []int) {
	return file_receiver_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetTPResponse) GetFounded() bool {
	if x != nil {
		return x.Founded
	}
	return false
}

func (x *GetTPResponse) GetTp() string {
	if x != nil {
		return x.Tp
	}
	return ""
}

func (x *GetTPResponse) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

func (x *GetTPResponse) GetIsReceipt() bool {
	if x != nil {
		return x.IsReceipt
	}
	return false
}

func (x *GetTPResponse) GetReceiptUrl() string {
	if x != nil {
		return x.ReceiptUrl
	}
	return ""
}

func (x *GetTPResponse) GetIsSuccess() bool {
	if x != nil {
		return x.IsSuccess
	}
	return false
}

func (x *GetTPResponse) GetIsValidationError() bool {
	if x != nil {
		return x.IsValidationError
	}
	return false
}

func (x *GetTPResponse) GetIsInternalError() bool {
	if x != nil {
		return x.IsInternalError
	}
	return false
}

func (x *GetTPResponse) GetIsNew() bool {
	if x != nil {
		return x.IsNew
	}
	return false
}

func (x *GetTPResponse) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *GetTPResponse) GetSendTaskNextAt() string {
	if x != nil {
		return x.SendTaskNextAt
	}
	return ""
}

func (x *GetTPResponse) GetContent() []*Directory {
	if x != nil {
		return x.Content
	}
	return nil
}

func (x *GetTPResponse) GetErrorText() string {
	if x != nil {
		return x.ErrorText
	}
	return ""
}

func (x *GetTPResponse) GetErrorCode() string {
	if x != nil {
		return x.ErrorCode
	}
	return ""
}

func (x *GetTPResponse) GetTimeLayout() string {
	if x != nil {
		return x.TimeLayout
	}
	return ""
}

type Directory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Files []string `protobuf:"bytes,2,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *Directory) Reset() {
	*x = Directory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_receiver_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Directory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Directory) ProtoMessage() {}

func (x *Directory) ProtoReflect() protoreflect.Message {
	mi := &file_receiver_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Directory.ProtoReflect.Descriptor instead.
func (*Directory) Descriptor() ([]byte, []int) {
	return file_receiver_service_proto_rawDescGZIP(), []int{2}
}

func (x *Directory) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Directory) GetFiles() []string {
	if x != nil {
		return x.Files
	}
	return nil
}

var File_receiver_service_proto protoreflect.FileDescriptor

var file_receiver_service_proto_rawDesc = []byte{
	0x0a, 0x16, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67,
	0x65, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32,
	0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x30, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x54, 0x50, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x74, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6f, 0x63, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x64, 0x6f, 0x63, 0x22, 0xf5, 0x03, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x54,
	0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x66, 0x6f, 0x75,
	0x6e, 0x64, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x66, 0x6f, 0x75, 0x6e,
	0x64, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x74, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x74, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x69,
	0x73, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09,
	0x69, 0x73, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x70, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x73,
	0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69,
	0x73, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x2c, 0x0a, 0x11, 0x69, 0x73, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x11, 0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x28, 0x0a, 0x0f, 0x69, 0x73, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0f, 0x69, 0x73, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x73, 0x4e, 0x65, 0x77, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x05, 0x69, 0x73, 0x4e, 0x65, 0x77, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x65, 0x6e, 0x64, 0x54, 0x61, 0x73, 0x6b,
	0x4e, 0x65, 0x78, 0x74, 0x41, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x65,
	0x6e, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x4e, 0x65, 0x78, 0x74, 0x41, 0x74, 0x12, 0x35, 0x0a, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x65, 0x78, 0x74,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x65, 0x78,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x0e,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x18, 0x0f, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x22,
	0x35, 0x0a, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x32, 0xca, 0x01, 0x0a, 0x0f, 0x50, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x57, 0x0a, 0x0b, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x12, 0x12, 0x10, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2f, 0x68, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x12, 0x5e, 0x0a, 0x05, 0x47, 0x65, 0x74, 0x54, 0x50, 0x12, 0x1e, 0x2e, 0x70,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e,
	0x47, 0x65, 0x74, 0x54, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2e,
	0x47, 0x65, 0x74, 0x54, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x14, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x12, 0x0c, 0x2f, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72,
	0x2f, 0x74, 0x70, 0x42, 0xc2, 0x01, 0x92, 0x41, 0x8d, 0x01, 0x12, 0x1d, 0x0a, 0x14, 0x50, 0x61,
	0x63, 0x6b, 0x61, 0x67, 0x65, 0x20, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x20, 0x41,
	0x50, 0x49, 0x32, 0x05, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x1a, 0x0e, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x68, 0x6f, 0x73, 0x74, 0x3a, 0x38, 0x30, 0x38, 0x30, 0x2a, 0x01, 0x01, 0x32, 0x10, 0x61, 0x70,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x10,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x6a, 0x73, 0x6f, 0x6e,
	0x5a, 0x23, 0x0a, 0x21, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x74, 0x68,
	0x12, 0x13, 0x08, 0x02, 0x1a, 0x0d, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x20, 0x02, 0x62, 0x10, 0x0a, 0x0e, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x4b, 0x65,
	0x79, 0x41, 0x75, 0x74, 0x68, 0x12, 0x00, 0x5a, 0x2f, 0x2e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x2d, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x2d,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x3b, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x5f,
	0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_receiver_service_proto_rawDescOnce sync.Once
	file_receiver_service_proto_rawDescData = file_receiver_service_proto_rawDesc
)

func file_receiver_service_proto_rawDescGZIP() []byte {
	file_receiver_service_proto_rawDescOnce.Do(func() {
		file_receiver_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_receiver_service_proto_rawDescData)
	})
	return file_receiver_service_proto_rawDescData
}

var file_receiver_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_receiver_service_proto_goTypes = []interface{}{
	(*GetTPRequest)(nil),  // 0: package_receiver.GetTPRequest
	(*GetTPResponse)(nil), // 1: package_receiver.GetTPResponse
	(*Directory)(nil),     // 2: package_receiver.Directory
	(*emptypb.Empty)(nil), // 3: google.protobuf.Empty
}
var file_receiver_service_proto_depIdxs = []int32{
	2, // 0: package_receiver.GetTPResponse.content:type_name -> package_receiver.Directory
	3, // 1: package_receiver.PackageReceiver.CheckHealth:input_type -> google.protobuf.Empty
	0, // 2: package_receiver.PackageReceiver.GetTP:input_type -> package_receiver.GetTPRequest
	3, // 3: package_receiver.PackageReceiver.CheckHealth:output_type -> google.protobuf.Empty
	1, // 4: package_receiver.PackageReceiver.GetTP:output_type -> package_receiver.GetTPResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_receiver_service_proto_init() }
func file_receiver_service_proto_init() {
	if File_receiver_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_receiver_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTPRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_receiver_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTPResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_receiver_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Directory); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_receiver_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_receiver_service_proto_goTypes,
		DependencyIndexes: file_receiver_service_proto_depIdxs,
		MessageInfos:      file_receiver_service_proto_msgTypes,
	}.Build()
	File_receiver_service_proto = out.File
	file_receiver_service_proto_rawDesc = nil
	file_receiver_service_proto_goTypes = nil
	file_receiver_service_proto_depIdxs = nil
}
