// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: signerMessages.proto

package pb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ProtocolMessage_Type int32

const (
	ProtocolMessage_DEFAULT       ProtocolMessage_Type = 0
	ProtocolMessage_SIGN_REQUEST  ProtocolMessage_Type = 100
	ProtocolMessage_SIGN_RESPONSE ProtocolMessage_Type = 101
)

// Enum value maps for ProtocolMessage_Type.
var (
	ProtocolMessage_Type_name = map[int32]string{
		0:   "DEFAULT",
		100: "SIGN_REQUEST",
		101: "SIGN_RESPONSE",
	}
	ProtocolMessage_Type_value = map[string]int32{
		"DEFAULT":       0,
		"SIGN_REQUEST":  100,
		"SIGN_RESPONSE": 101,
	}
)

func (x ProtocolMessage_Type) Enum() *ProtocolMessage_Type {
	p := new(ProtocolMessage_Type)
	*p = x
	return p
}

func (x ProtocolMessage_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProtocolMessage_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_signerMessages_proto_enumTypes[0].Descriptor()
}

func (ProtocolMessage_Type) Type() protoreflect.EnumType {
	return &file_signerMessages_proto_enumTypes[0]
}

func (x ProtocolMessage_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProtocolMessage_Type.Descriptor instead.
func (ProtocolMessage_Type) EnumDescriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{0, 0}
}

type VerifyResponse_Status int32

const (
	VerifyResponse_STATUS_UNSET VerifyResponse_Status = 0
	VerifyResponse_OK           VerifyResponse_Status = 1
	VerifyResponse_INVALID      VerifyResponse_Status = 2
)

// Enum value maps for VerifyResponse_Status.
var (
	VerifyResponse_Status_name = map[int32]string{
		0: "STATUS_UNSET",
		1: "OK",
		2: "INVALID",
	}
	VerifyResponse_Status_value = map[string]int32{
		"STATUS_UNSET": 0,
		"OK":           1,
		"INVALID":      2,
	}
)

func (x VerifyResponse_Status) Enum() *VerifyResponse_Status {
	p := new(VerifyResponse_Status)
	*p = x
	return p
}

func (x VerifyResponse_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VerifyResponse_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_signerMessages_proto_enumTypes[1].Descriptor()
}

func (VerifyResponse_Status) Type() protoreflect.EnumType {
	return &file_signerMessages_proto_enumTypes[1]
}

func (x VerifyResponse_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VerifyResponse_Status.Descriptor instead.
func (VerifyResponse_Status) EnumDescriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{5, 0}
}

type MembershipResponse_Status int32

const (
	MembershipResponse_STATUS_UNSET MembershipResponse_Status = 0
	MembershipResponse_OK           MembershipResponse_Status = 1
	MembershipResponse_INVALID      MembershipResponse_Status = 2
)

// Enum value maps for MembershipResponse_Status.
var (
	MembershipResponse_Status_name = map[int32]string{
		0: "STATUS_UNSET",
		1: "OK",
		2: "INVALID",
	}
	MembershipResponse_Status_value = map[string]int32{
		"STATUS_UNSET": 0,
		"OK":           1,
		"INVALID":      2,
	}
)

func (x MembershipResponse_Status) Enum() *MembershipResponse_Status {
	p := new(MembershipResponse_Status)
	*p = x
	return p
}

func (x MembershipResponse_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MembershipResponse_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_signerMessages_proto_enumTypes[2].Descriptor()
}

func (MembershipResponse_Status) Type() protoreflect.EnumType {
	return &file_signerMessages_proto_enumTypes[2]
}

func (x MembershipResponse_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MembershipResponse_Status.Descriptor instead.
func (MembershipResponse_Status) EnumDescriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{6, 0}
}

type ProtocolMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type          ProtocolMessage_Type `protobuf:"varint,1,opt,name=type,proto3,enum=pb.ProtocolMessage_Type" json:"type,omitempty"`
	CorrelationId string               `protobuf:"bytes,2,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	Content       []byte               `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ProtocolMessage) Reset() {
	*x = ProtocolMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_signerMessages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProtocolMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProtocolMessage) ProtoMessage() {}

func (x *ProtocolMessage) ProtoReflect() protoreflect.Message {
	mi := &file_signerMessages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProtocolMessage.ProtoReflect.Descriptor instead.
func (*ProtocolMessage) Descriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{0}
}

func (x *ProtocolMessage) GetType() ProtocolMessage_Type {
	if x != nil {
		return x.Type
	}
	return ProtocolMessage_DEFAULT
}

func (x *ProtocolMessage) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *ProtocolMessage) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type ClientSignMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UUID                 string `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	SmartContractAddress string `protobuf:"bytes,2,opt,name=SmartContractAddress,proto3" json:"SmartContractAddress,omitempty"`
	//uint32 t = 3;
	//uint32 n = 4;
	//string Scheme = 5;
	Content []byte `protobuf:"bytes,6,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ClientSignMessage) Reset() {
	*x = ClientSignMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_signerMessages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientSignMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientSignMessage) ProtoMessage() {}

func (x *ClientSignMessage) ProtoReflect() protoreflect.Message {
	mi := &file_signerMessages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientSignMessage.ProtoReflect.Descriptor instead.
func (*ClientSignMessage) Descriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{1}
}

func (x *ClientSignMessage) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

func (x *ClientSignMessage) GetSmartContractAddress() string {
	if x != nil {
		return x.SmartContractAddress
	}
	return ""
}

func (x *ClientSignMessage) GetContent() []byte {
	if x != nil {
		return x.Content
	}
	return nil
}

type ClientSignResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scheme    string `protobuf:"bytes,1,opt,name=scheme,proto3" json:"scheme,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *ClientSignResponse) Reset() {
	*x = ClientSignResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_signerMessages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientSignResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientSignResponse) ProtoMessage() {}

func (x *ClientSignResponse) ProtoReflect() protoreflect.Message {
	mi := &file_signerMessages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientSignResponse.ProtoReflect.Descriptor instead.
func (*ClientSignResponse) Descriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{2}
}

func (x *ClientSignResponse) GetScheme() string {
	if x != nil {
		return x.Scheme
	}
	return ""
}

func (x *ClientSignResponse) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type SignResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UUID      string `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *SignResponse) Reset() {
	*x = SignResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_signerMessages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignResponse) ProtoMessage() {}

func (x *SignResponse) ProtoReflect() protoreflect.Message {
	mi := &file_signerMessages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignResponse.ProtoReflect.Descriptor instead.
func (*SignResponse) Descriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{3}
}

func (x *SignResponse) GetUUID() string {
	if x != nil {
		return x.UUID
	}
	return ""
}

func (x *SignResponse) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type ClientVerifyMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scheme    string `protobuf:"bytes,1,opt,name=Scheme,proto3" json:"Scheme,omitempty"`
	PublicKey []byte `protobuf:"bytes,2,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Digest    []byte `protobuf:"bytes,3,opt,name=digest,proto3" json:"digest,omitempty"`
	Signature []byte `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *ClientVerifyMessage) Reset() {
	*x = ClientVerifyMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_signerMessages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientVerifyMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientVerifyMessage) ProtoMessage() {}

func (x *ClientVerifyMessage) ProtoReflect() protoreflect.Message {
	mi := &file_signerMessages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientVerifyMessage.ProtoReflect.Descriptor instead.
func (*ClientVerifyMessage) Descriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{4}
}

func (x *ClientVerifyMessage) GetScheme() string {
	if x != nil {
		return x.Scheme
	}
	return ""
}

func (x *ClientVerifyMessage) GetPublicKey() []byte {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (x *ClientVerifyMessage) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

func (x *ClientVerifyMessage) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type VerifyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status VerifyResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=pb.VerifyResponse_Status" json:"status,omitempty"`
}

func (x *VerifyResponse) Reset() {
	*x = VerifyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_signerMessages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyResponse) ProtoMessage() {}

func (x *VerifyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_signerMessages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyResponse.ProtoReflect.Descriptor instead.
func (*VerifyResponse) Descriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{5}
}

func (x *VerifyResponse) GetStatus() VerifyResponse_Status {
	if x != nil {
		return x.Status
	}
	return VerifyResponse_STATUS_UNSET
}

type MembershipResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status MembershipResponse_Status `protobuf:"varint,1,opt,name=status,proto3,enum=pb.MembershipResponse_Status" json:"status,omitempty"`
	Peers  []*MembershipResponsePeer `protobuf:"bytes,2,rep,name=peers,proto3" json:"peers,omitempty"`
}

func (x *MembershipResponse) Reset() {
	*x = MembershipResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_signerMessages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MembershipResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MembershipResponse) ProtoMessage() {}

func (x *MembershipResponse) ProtoReflect() protoreflect.Message {
	mi := &file_signerMessages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MembershipResponse.ProtoReflect.Descriptor instead.
func (*MembershipResponse) Descriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{6}
}

func (x *MembershipResponse) GetStatus() MembershipResponse_Status {
	if x != nil {
		return x.Status
	}
	return MembershipResponse_STATUS_UNSET
}

func (x *MembershipResponse) GetPeers() []*MembershipResponsePeer {
	if x != nil {
		return x.Peers
	}
	return nil
}

type MembershipResponsePeer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr []string `protobuf:"bytes,2,rep,name=addr,proto3" json:"addr,omitempty"`
}

func (x *MembershipResponsePeer) Reset() {
	*x = MembershipResponsePeer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_signerMessages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MembershipResponsePeer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MembershipResponsePeer) ProtoMessage() {}

func (x *MembershipResponsePeer) ProtoReflect() protoreflect.Message {
	mi := &file_signerMessages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MembershipResponsePeer.ProtoReflect.Descriptor instead.
func (*MembershipResponsePeer) Descriptor() ([]byte, []int) {
	return file_signerMessages_proto_rawDescGZIP(), []int{6, 0}
}

func (x *MembershipResponsePeer) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *MembershipResponsePeer) GetAddr() []string {
	if x != nil {
		return x.Addr
	}
	return nil
}

var File_signerMessages_proto protoreflect.FileDescriptor

var file_signerMessages_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0xba, 0x01, 0x0a, 0x0f, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2c,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x70,
	0x62, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x25, 0x0a, 0x0e,
	0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x38, 0x0a,
	0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x44, 0x45, 0x46, 0x41, 0x55, 0x4c, 0x54,
	0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x49, 0x47, 0x4e, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45,
	0x53, 0x54, 0x10, 0x64, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x49, 0x47, 0x4e, 0x5f, 0x52, 0x45, 0x53,
	0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10, 0x65, 0x22, 0x75, 0x0a, 0x11, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x53, 0x69, 0x67, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x55, 0x55, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x55, 0x49, 0x44,
	0x12, 0x32, 0x0a, 0x14, 0x53, 0x6d, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63,
	0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14,
	0x53, 0x6d, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x4a,
	0x0a, 0x12, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x69, 0x67, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x40, 0x0a, 0x0c, 0x53, 0x69,
	0x67, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x55, 0x55,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x55, 0x55, 0x49, 0x44, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x82, 0x01, 0x0a,
	0x13, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x64, 0x69, 0x67,
	0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x22, 0x74, 0x0a, 0x0e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2f, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x10, 0x0a, 0x0c, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x45, 0x54,
	0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e,
	0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x02, 0x22, 0xdb, 0x01, 0x0a, 0x12, 0x4d, 0x65, 0x6d, 0x62,
	0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d,
	0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x73, 0x68, 0x69, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x31, 0x0a, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x73, 0x68, 0x69, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x70, 0x65, 0x65,
	0x72, 0x52, 0x05, 0x70, 0x65, 0x65, 0x72, 0x73, 0x1a, 0x2a, 0x0a, 0x04, 0x70, 0x65, 0x65, 0x72,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04,
	0x61, 0x64, 0x64, 0x72, 0x22, 0x2f, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x10,
	0x0a, 0x0c, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x45, 0x54, 0x10, 0x00,
	0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e, 0x56, 0x41,
	0x4c, 0x49, 0x44, 0x10, 0x02, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_signerMessages_proto_rawDescOnce sync.Once
	file_signerMessages_proto_rawDescData = file_signerMessages_proto_rawDesc
)

func file_signerMessages_proto_rawDescGZIP() []byte {
	file_signerMessages_proto_rawDescOnce.Do(func() {
		file_signerMessages_proto_rawDescData = protoimpl.X.CompressGZIP(file_signerMessages_proto_rawDescData)
	})
	return file_signerMessages_proto_rawDescData
}

var file_signerMessages_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_signerMessages_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_signerMessages_proto_goTypes = []interface{}{
	(ProtocolMessage_Type)(0),      // 0: pb.ProtocolMessage.Type
	(VerifyResponse_Status)(0),     // 1: pb.VerifyResponse.Status
	(MembershipResponse_Status)(0), // 2: pb.MembershipResponse.Status
	(*ProtocolMessage)(nil),        // 3: pb.ProtocolMessage
	(*ClientSignMessage)(nil),      // 4: pb.ClientSignMessage
	(*ClientSignResponse)(nil),     // 5: pb.ClientSignResponse
	(*SignResponse)(nil),           // 6: pb.SignResponse
	(*ClientVerifyMessage)(nil),    // 7: pb.ClientVerifyMessage
	(*VerifyResponse)(nil),         // 8: pb.VerifyResponse
	(*MembershipResponse)(nil),     // 9: pb.MembershipResponse
	(*MembershipResponsePeer)(nil), // 10: pb.MembershipResponse.peer
}
var file_signerMessages_proto_depIdxs = []int32{
	0,  // 0: pb.ProtocolMessage.type:type_name -> pb.ProtocolMessage.Type
	1,  // 1: pb.VerifyResponse.status:type_name -> pb.VerifyResponse.Status
	2,  // 2: pb.MembershipResponse.status:type_name -> pb.MembershipResponse.Status
	10, // 3: pb.MembershipResponse.peers:type_name -> pb.MembershipResponse.peer
	4,  // [4:4] is the sub-list for method output_type
	4,  // [4:4] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_signerMessages_proto_init() }
func file_signerMessages_proto_init() {
	if File_signerMessages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_signerMessages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProtocolMessage); i {
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
		file_signerMessages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientSignMessage); i {
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
		file_signerMessages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientSignResponse); i {
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
		file_signerMessages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignResponse); i {
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
		file_signerMessages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientVerifyMessage); i {
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
		file_signerMessages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifyResponse); i {
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
		file_signerMessages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MembershipResponse); i {
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
		file_signerMessages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MembershipResponsePeer); i {
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
			RawDescriptor: file_signerMessages_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_signerMessages_proto_goTypes,
		DependencyIndexes: file_signerMessages_proto_depIdxs,
		EnumInfos:         file_signerMessages_proto_enumTypes,
		MessageInfos:      file_signerMessages_proto_msgTypes,
	}.Build()
	File_signerMessages_proto = out.File
	file_signerMessages_proto_rawDesc = nil
	file_signerMessages_proto_goTypes = nil
	file_signerMessages_proto_depIdxs = nil
}
