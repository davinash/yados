// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.14.0
// source: messages.proto

package gen

import (
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

type Member_MemberStatus int32

const (
	Member_Healthy   Member_MemberStatus = 0
	Member_UnHealthy Member_MemberStatus = 1
)

// Enum value maps for Member_MemberStatus.
var (
	Member_MemberStatus_name = map[int32]string{
		0: "Healthy",
		1: "UnHealthy",
	}
	Member_MemberStatus_value = map[string]int32{
		"Healthy":   0,
		"UnHealthy": 1,
	}
)

func (x Member_MemberStatus) Enum() *Member_MemberStatus {
	p := new(Member_MemberStatus)
	*p = x
	return p
}

func (x Member_MemberStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Member_MemberStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_messages_proto_enumTypes[0].Descriptor()
}

func (Member_MemberStatus) Type() protoreflect.EnumType {
	return &file_messages_proto_enumTypes[0]
}

func (x Member_MemberStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Member_MemberStatus.Descriptor instead.
func (Member_MemberStatus) EnumDescriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0, 0}
}

type Member struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name             string              `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Address          string              `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Port             int32               `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	Status           Member_MemberStatus `protobuf:"varint,4,opt,name=status,proto3,enum=Member_MemberStatus" json:"status,omitempty"`
	ReceivedHCLastAt int64               `protobuf:"varint,5,opt,name=ReceivedHCLastAt,proto3" json:"ReceivedHCLastAt,omitempty"`
}

func (x *Member) Reset() {
	*x = Member{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Member) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Member) ProtoMessage() {}

func (x *Member) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Member.ProtoReflect.Descriptor instead.
func (*Member) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0}
}

func (x *Member) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Member) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Member) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *Member) GetStatus() Member_MemberStatus {
	if x != nil {
		return x.Status
	}
	return Member_Healthy
}

func (x *Member) GetReceivedHCLastAt() int64 {
	if x != nil {
		return x.ReceivedHCLastAt
	}
	return 0
}

type NewMemberRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Member *Member `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
}

func (x *NewMemberRequest) Reset() {
	*x = NewMemberRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewMemberRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewMemberRequest) ProtoMessage() {}

func (x *NewMemberRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewMemberRequest.ProtoReflect.Descriptor instead.
func (*NewMemberRequest) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{1}
}

func (x *NewMemberRequest) GetMember() *Member {
	if x != nil {
		return x.Member
	}
	return nil
}

type NewMemberReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Member *Member `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
}

func (x *NewMemberReply) Reset() {
	*x = NewMemberReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewMemberReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewMemberReply) ProtoMessage() {}

func (x *NewMemberReply) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewMemberReply.ProtoReflect.Descriptor instead.
func (*NewMemberReply) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{2}
}

func (x *NewMemberReply) GetMember() *Member {
	if x != nil {
		return x.Member
	}
	return nil
}

type ListOfPeersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListOfPeersRequest) Reset() {
	*x = ListOfPeersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOfPeersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOfPeersRequest) ProtoMessage() {}

func (x *ListOfPeersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOfPeersRequest.ProtoReflect.Descriptor instead.
func (*ListOfPeersRequest) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{3}
}

type ListOfPeersReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Member []*Member `protobuf:"bytes,1,rep,name=member,proto3" json:"member,omitempty"`
}

func (x *ListOfPeersReply) Reset() {
	*x = ListOfPeersReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListOfPeersReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListOfPeersReply) ProtoMessage() {}

func (x *ListOfPeersReply) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListOfPeersReply.ProtoReflect.Descriptor instead.
func (*ListOfPeersReply) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{4}
}

func (x *ListOfPeersReply) GetMember() []*Member {
	if x != nil {
		return x.Member
	}
	return nil
}

type StopServerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StopServerRequest) Reset() {
	*x = StopServerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopServerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopServerRequest) ProtoMessage() {}

func (x *StopServerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopServerRequest.ProtoReflect.Descriptor instead.
func (*StopServerRequest) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{5}
}

type StopServerReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StopServerReply) Reset() {
	*x = StopServerReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopServerReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopServerReply) ProtoMessage() {}

func (x *StopServerReply) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopServerReply.ProtoReflect.Descriptor instead.
func (*StopServerReply) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{6}
}

type RemoveServerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Member *Member `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
}

func (x *RemoveServerRequest) Reset() {
	*x = RemoveServerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveServerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveServerRequest) ProtoMessage() {}

func (x *RemoveServerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveServerRequest.ProtoReflect.Descriptor instead.
func (*RemoveServerRequest) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{7}
}

func (x *RemoveServerRequest) GetMember() *Member {
	if x != nil {
		return x.Member
	}
	return nil
}

type RemoveServerReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveServerReply) Reset() {
	*x = RemoveServerReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveServerReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveServerReply) ProtoMessage() {}

func (x *RemoveServerReply) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveServerReply.ProtoReflect.Descriptor instead.
func (*RemoveServerReply) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{8}
}

type HealthStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Member *Member `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
}

func (x *HealthStatusRequest) Reset() {
	*x = HealthStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthStatusRequest) ProtoMessage() {}

func (x *HealthStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthStatusRequest.ProtoReflect.Descriptor instead.
func (*HealthStatusRequest) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{9}
}

func (x *HealthStatusRequest) GetMember() *Member {
	if x != nil {
		return x.Member
	}
	return nil
}

type HealthStatusReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HealthStatusReply) Reset() {
	*x = HealthStatusReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthStatusReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthStatusReply) ProtoMessage() {}

func (x *HealthStatusReply) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthStatusReply.ProtoReflect.Descriptor instead.
func (*HealthStatusReply) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{10}
}

type StoreCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Replication int32  `protobuf:"varint,2,opt,name=replication,proto3" json:"replication,omitempty"`
}

func (x *StoreCreateRequest) Reset() {
	*x = StoreCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreCreateRequest) ProtoMessage() {}

func (x *StoreCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreCreateRequest.ProtoReflect.Descriptor instead.
func (*StoreCreateRequest) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{11}
}

func (x *StoreCreateRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *StoreCreateRequest) GetReplication() int32 {
	if x != nil {
		return x.Replication
	}
	return 0
}

type StoreCreateReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StoreCreateReply) Reset() {
	*x = StoreCreateReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreCreateReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreCreateReply) ProtoMessage() {}

func (x *StoreCreateReply) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreCreateReply.ProtoReflect.Descriptor instead.
func (*StoreCreateReply) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{12}
}

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term          uint64 `protobuf:"varint,1,opt,name=Term,proto3" json:"Term,omitempty"`
	CandidateName string `protobuf:"bytes,2,opt,name=CandidateName,proto3" json:"CandidateName,omitempty"`
	StoreName     string `protobuf:"bytes,3,opt,name=storeName,proto3" json:"storeName,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{13}
}

func (x *VoteRequest) GetTerm() uint64 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *VoteRequest) GetCandidateName() string {
	if x != nil {
		return x.CandidateName
	}
	return ""
}

func (x *VoteRequest) GetStoreName() string {
	if x != nil {
		return x.StoreName
	}
	return ""
}

type VoteRequestReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *VoteRequestReply) Reset() {
	*x = VoteRequestReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequestReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequestReply) ProtoMessage() {}

func (x *VoteRequestReply) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequestReply.ProtoReflect.Descriptor instead.
func (*VoteRequestReply) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{14}
}

var File_messages_proto protoreflect.FileDescriptor

var file_messages_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xd0, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x2c, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e,
	0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2a, 0x0a, 0x10, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x48, 0x43, 0x4c, 0x61, 0x73, 0x74, 0x41, 0x74, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x48,
	0x43, 0x4c, 0x61, 0x73, 0x74, 0x41, 0x74, 0x22, 0x2a, 0x0a, 0x0c, 0x4d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x79, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x6e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68,
	0x79, 0x10, 0x01, 0x22, 0x33, 0x0a, 0x10, 0x4e, 0x65, 0x77, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72,
	0x52, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x31, 0x0a, 0x0e, 0x4e, 0x65, 0x77, 0x4d,
	0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1f, 0x0a, 0x06, 0x6d, 0x65,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x14, 0x0a, 0x12, 0x4c,
	0x69, 0x73, 0x74, 0x4f, 0x66, 0x50, 0x65, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x33, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x66, 0x50, 0x65, 0x65, 0x72, 0x73,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1f, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x06,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x13, 0x0a, 0x11, 0x53, 0x74, 0x6f, 0x70, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x11, 0x0a, 0x0f, 0x53,
	0x74, 0x6f, 0x70, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x36,
	0x0a, 0x13, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x06,
	0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x13, 0x0a, 0x11, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x36, 0x0a, 0x13, 0x48,
	0x65, 0x61, 0x6c, 0x74, 0x68, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x07, 0x2e, 0x4d, 0x65, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x06, 0x6d, 0x65, 0x6d,
	0x62, 0x65, 0x72, 0x22, 0x13, 0x0a, 0x11, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x4a, 0x0a, 0x12, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x12, 0x0a, 0x10, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x65, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x24, 0x0a, 0x0d, 0x43,
	0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x12, 0x0a, 0x10, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x42, 0x06, 0x5a, 0x04, 0x2f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_messages_proto_rawDescOnce sync.Once
	file_messages_proto_rawDescData = file_messages_proto_rawDesc
)

func file_messages_proto_rawDescGZIP() []byte {
	file_messages_proto_rawDescOnce.Do(func() {
		file_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_proto_rawDescData)
	})
	return file_messages_proto_rawDescData
}

var file_messages_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_messages_proto_goTypes = []interface{}{
	(Member_MemberStatus)(0),    // 0: Member.MemberStatus
	(*Member)(nil),              // 1: Member
	(*NewMemberRequest)(nil),    // 2: NewMemberRequest
	(*NewMemberReply)(nil),      // 3: NewMemberReply
	(*ListOfPeersRequest)(nil),  // 4: ListOfPeersRequest
	(*ListOfPeersReply)(nil),    // 5: ListOfPeersReply
	(*StopServerRequest)(nil),   // 6: StopServerRequest
	(*StopServerReply)(nil),     // 7: StopServerReply
	(*RemoveServerRequest)(nil), // 8: RemoveServerRequest
	(*RemoveServerReply)(nil),   // 9: RemoveServerReply
	(*HealthStatusRequest)(nil), // 10: HealthStatusRequest
	(*HealthStatusReply)(nil),   // 11: HealthStatusReply
	(*StoreCreateRequest)(nil),  // 12: StoreCreateRequest
	(*StoreCreateReply)(nil),    // 13: StoreCreateReply
	(*VoteRequest)(nil),         // 14: VoteRequest
	(*VoteRequestReply)(nil),    // 15: VoteRequestReply
}
var file_messages_proto_depIdxs = []int32{
	0, // 0: Member.status:type_name -> Member.MemberStatus
	1, // 1: NewMemberRequest.member:type_name -> Member
	1, // 2: NewMemberReply.member:type_name -> Member
	1, // 3: ListOfPeersReply.member:type_name -> Member
	1, // 4: RemoveServerRequest.member:type_name -> Member
	1, // 5: HealthStatusRequest.member:type_name -> Member
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_messages_proto_init() }
func file_messages_proto_init() {
	if File_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Member); i {
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
		file_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewMemberRequest); i {
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
		file_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewMemberReply); i {
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
		file_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListOfPeersRequest); i {
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
		file_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListOfPeersReply); i {
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
		file_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopServerRequest); i {
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
		file_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StopServerReply); i {
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
		file_messages_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveServerRequest); i {
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
		file_messages_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveServerReply); i {
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
		file_messages_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthStatusRequest); i {
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
		file_messages_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthStatusReply); i {
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
		file_messages_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreCreateRequest); i {
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
		file_messages_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreCreateReply); i {
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
		file_messages_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
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
		file_messages_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequestReply); i {
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
			RawDescriptor: file_messages_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messages_proto_goTypes,
		DependencyIndexes: file_messages_proto_depIdxs,
		EnumInfos:         file_messages_proto_enumTypes,
		MessageInfos:      file_messages_proto_msgTypes,
	}.Build()
	File_messages_proto = out.File
	file_messages_proto_rawDesc = nil
	file_messages_proto_goTypes = nil
	file_messages_proto_depIdxs = nil
}
