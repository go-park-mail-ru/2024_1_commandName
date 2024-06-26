// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: contacts.proto

package chats

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserIDContacts struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID uint64 `protobuf:"varint,1,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *UserIDContacts) Reset() {
	*x = UserIDContacts{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contacts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIDContacts) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIDContacts) ProtoMessage() {}

func (x *UserIDContacts) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIDContacts.ProtoReflect.Descriptor instead.
func (*UserIDContacts) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{0}
}

func (x *UserIDContacts) GetUserID() uint64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

type UserIDArray struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*UserIDContacts `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *UserIDArray) Reset() {
	*x = UserIDArray{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contacts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserIDArray) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserIDArray) ProtoMessage() {}

func (x *UserIDArray) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserIDArray.ProtoReflect.Descriptor instead.
func (*UserIDArray) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{1}
}

func (x *UserIDArray) GetUsers() []*UserIDContacts {
	if x != nil {
		return x.Users
	}
	return nil
}

type AddToAllReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users        *UserIDArray `protobuf:"bytes,1,opt,name=users,proto3" json:"users,omitempty"`
	UserAddingID uint64       `protobuf:"varint,2,opt,name=userAddingID,proto3" json:"userAddingID,omitempty"`
}

func (x *AddToAllReq) Reset() {
	*x = AddToAllReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contacts_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddToAllReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddToAllReq) ProtoMessage() {}

func (x *AddToAllReq) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddToAllReq.ProtoReflect.Descriptor instead.
func (*AddToAllReq) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{2}
}

func (x *AddToAllReq) GetUsers() *UserIDArray {
	if x != nil {
		return x.Users
	}
	return nil
}

func (x *AddToAllReq) GetUserAddingID() uint64 {
	if x != nil {
		return x.UserAddingID
	}
	return 0
}

type Person struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID           uint64                 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Username     string                 `protobuf:"bytes,2,opt,name=Username,proto3" json:"Username,omitempty"`
	Email        string                 `protobuf:"bytes,3,opt,name=Email,proto3" json:"Email,omitempty"`
	Name         string                 `protobuf:"bytes,4,opt,name=Name,proto3" json:"Name,omitempty"`
	Surname      string                 `protobuf:"bytes,5,opt,name=Surname,proto3" json:"Surname,omitempty"`
	About        string                 `protobuf:"bytes,6,opt,name=About,proto3" json:"About,omitempty"`
	Password     string                 `protobuf:"bytes,7,opt,name=Password,proto3" json:"Password,omitempty"`
	CreateTime   *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=CreateTime,proto3" json:"CreateTime,omitempty"`
	LastSeenDate *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=LastSeenDate,proto3" json:"LastSeenDate,omitempty"`
	AvatarPath   string                 `protobuf:"bytes,10,opt,name=AvatarPath,proto3" json:"AvatarPath,omitempty"`
	PasswordSalt string                 `protobuf:"bytes,11,opt,name=PasswordSalt,proto3" json:"PasswordSalt,omitempty"`
}

func (x *Person) Reset() {
	*x = Person{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contacts_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Person) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Person) ProtoMessage() {}

func (x *Person) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Person.ProtoReflect.Descriptor instead.
func (*Person) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{3}
}

func (x *Person) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Person) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *Person) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Person) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Person) GetSurname() string {
	if x != nil {
		return x.Surname
	}
	return ""
}

func (x *Person) GetAbout() string {
	if x != nil {
		return x.About
	}
	return ""
}

func (x *Person) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *Person) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Person) GetLastSeenDate() *timestamppb.Timestamp {
	if x != nil {
		return x.LastSeenDate
	}
	return nil
}

func (x *Person) GetAvatarPath() string {
	if x != nil {
		return x.AvatarPath
	}
	return ""
}

func (x *Person) GetPasswordSalt() string {
	if x != nil {
		return x.PasswordSalt
	}
	return ""
}

type PersonArray struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Persons []*Person `protobuf:"bytes,1,rep,name=Persons,proto3" json:"Persons,omitempty"`
}

func (x *PersonArray) Reset() {
	*x = PersonArray{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contacts_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersonArray) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersonArray) ProtoMessage() {}

func (x *PersonArray) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersonArray.ProtoReflect.Descriptor instead.
func (*PersonArray) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{4}
}

func (x *PersonArray) GetPersons() []*Person {
	if x != nil {
		return x.Persons
	}
	return nil
}

type AddByUsernameReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserAddingID  uint64 `protobuf:"varint,1,opt,name=UserAddingID,proto3" json:"UserAddingID,omitempty"`
	UsernameToAdd string `protobuf:"bytes,2,opt,name=UsernameToAdd,proto3" json:"UsernameToAdd,omitempty"`
	UserToAddID   uint64 `protobuf:"varint,3,opt,name=UserToAddID,proto3" json:"UserToAddID,omitempty"`
}

func (x *AddByUsernameReq) Reset() {
	*x = AddByUsernameReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contacts_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddByUsernameReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddByUsernameReq) ProtoMessage() {}

func (x *AddByUsernameReq) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddByUsernameReq.ProtoReflect.Descriptor instead.
func (*AddByUsernameReq) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{5}
}

func (x *AddByUsernameReq) GetUserAddingID() uint64 {
	if x != nil {
		return x.UserAddingID
	}
	return 0
}

func (x *AddByUsernameReq) GetUsernameToAdd() string {
	if x != nil {
		return x.UsernameToAdd
	}
	return ""
}

func (x *AddByUsernameReq) GetUserToAddID() uint64 {
	if x != nil {
		return x.UserToAddID
	}
	return 0
}

type EmptyContacts struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dummy bool `protobuf:"varint,1,opt,name=Dummy,proto3" json:"Dummy,omitempty"`
}

func (x *EmptyContacts) Reset() {
	*x = EmptyContacts{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contacts_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyContacts) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyContacts) ProtoMessage() {}

func (x *EmptyContacts) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyContacts.ProtoReflect.Descriptor instead.
func (*EmptyContacts) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{6}
}

func (x *EmptyContacts) GetDummy() bool {
	if x != nil {
		return x.Dummy
	}
	return false
}

type BoolResponseContacts struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *BoolResponseContacts) Reset() {
	*x = BoolResponseContacts{}
	if protoimpl.UnsafeEnabled {
		mi := &file_contacts_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoolResponseContacts) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoolResponseContacts) ProtoMessage() {}

func (x *BoolResponseContacts) ProtoReflect() protoreflect.Message {
	mi := &file_contacts_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoolResponseContacts.ProtoReflect.Descriptor instead.
func (*BoolResponseContacts) Descriptor() ([]byte, []int) {
	return file_contacts_proto_rawDescGZIP(), []int{7}
}

func (x *BoolResponseContacts) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_contacts_proto protoreflect.FileDescriptor

var file_contacts_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x63, 0x68, 0x61, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x28, 0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x44, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x44, 0x22, 0x3a, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x41, 0x72, 0x72, 0x61,
	0x79, 0x12, 0x2b, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x5b,
	0x0a, 0x0b, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x12, 0x28, 0x0a,
	0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63,
	0x68, 0x61, 0x74, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x41, 0x72, 0x72, 0x61, 0x79,
	0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x41,
	0x64, 0x64, 0x69, 0x6e, 0x67, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x75,
	0x73, 0x65, 0x72, 0x41, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x49, 0x44, 0x22, 0xea, 0x02, 0x0a, 0x06,
	0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x53, 0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53,
	0x75, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x41, 0x62, 0x6f, 0x75, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x41, 0x62, 0x6f, 0x75, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x3a, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x3e, 0x0a, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x53, 0x65, 0x65, 0x6e,
	0x44, 0x61, 0x74, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x53, 0x65, 0x65, 0x6e,
	0x44, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x50, 0x61,
	0x74, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72,
	0x50, 0x61, 0x74, 0x68, 0x12, 0x22, 0x0a, 0x0c, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x53, 0x61, 0x6c, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x50, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x53, 0x61, 0x6c, 0x74, 0x22, 0x36, 0x0a, 0x0b, 0x50, 0x65, 0x72, 0x73,
	0x6f, 0x6e, 0x41, 0x72, 0x72, 0x61, 0x79, 0x12, 0x27, 0x0a, 0x07, 0x50, 0x65, 0x72, 0x73, 0x6f,
	0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73,
	0x2e, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x52, 0x07, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x73,
	0x22, 0x7e, 0x0a, 0x10, 0x41, 0x64, 0x64, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x52, 0x65, 0x71, 0x12, 0x22, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x41, 0x64, 0x64, 0x69,
	0x6e, 0x67, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x55, 0x73, 0x65, 0x72,
	0x41, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x49, 0x44, 0x12, 0x24, 0x0a, 0x0d, 0x55, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x12, 0x20,
	0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x49, 0x44, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x54, 0x6f, 0x41, 0x64, 0x64, 0x49, 0x44,
	0x22, 0x25, 0x0a, 0x0d, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x44, 0x75, 0x6d, 0x6d, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x05, 0x44, 0x75, 0x6d, 0x6d, 0x79, 0x22, 0x26, 0x0a, 0x14, 0x42, 0x6f, 0x6f, 0x6c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x32,
	0xd0, 0x01, 0x0a, 0x08, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x12, 0x38, 0x0a, 0x0b,
	0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x12, 0x15, 0x2e, 0x63, 0x68,
	0x61, 0x74, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x73, 0x1a, 0x12, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x6f,
	0x6e, 0x41, 0x72, 0x72, 0x61, 0x79, 0x12, 0x45, 0x0a, 0x14, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x6e,
	0x74, 0x61, 0x63, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x12, 0x43, 0x0a,
	0x10, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x41, 0x6c, 0x6c, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74,
	0x73, 0x12, 0x12, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x41,
	0x6c, 0x6c, 0x52, 0x65, 0x71, 0x1a, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x73, 0x2e, 0x42, 0x6f,
	0x6f, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x63,
	0x74, 0x73, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x63, 0x68, 0x61, 0x74, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_contacts_proto_rawDescOnce sync.Once
	file_contacts_proto_rawDescData = file_contacts_proto_rawDesc
)

func file_contacts_proto_rawDescGZIP() []byte {
	file_contacts_proto_rawDescOnce.Do(func() {
		file_contacts_proto_rawDescData = protoimpl.X.CompressGZIP(file_contacts_proto_rawDescData)
	})
	return file_contacts_proto_rawDescData
}

var file_contacts_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_contacts_proto_goTypes = []interface{}{
	(*UserIDContacts)(nil),        // 0: chats.UserIDContacts
	(*UserIDArray)(nil),           // 1: chats.UserIDArray
	(*AddToAllReq)(nil),           // 2: chats.AddToAllReq
	(*Person)(nil),                // 3: chats.Person
	(*PersonArray)(nil),           // 4: chats.PersonArray
	(*AddByUsernameReq)(nil),      // 5: chats.AddByUsernameReq
	(*EmptyContacts)(nil),         // 6: chats.EmptyContacts
	(*BoolResponseContacts)(nil),  // 7: chats.BoolResponseContacts
	(*timestamppb.Timestamp)(nil), // 8: google.protobuf.Timestamp
}
var file_contacts_proto_depIdxs = []int32{
	0, // 0: chats.UserIDArray.users:type_name -> chats.UserIDContacts
	1, // 1: chats.AddToAllReq.users:type_name -> chats.UserIDArray
	8, // 2: chats.Person.CreateTime:type_name -> google.protobuf.Timestamp
	8, // 3: chats.Person.LastSeenDate:type_name -> google.protobuf.Timestamp
	3, // 4: chats.PersonArray.Persons:type_name -> chats.Person
	0, // 5: chats.Contacts.GetContacts:input_type -> chats.UserIDContacts
	5, // 6: chats.Contacts.AddContactByUsername:input_type -> chats.AddByUsernameReq
	2, // 7: chats.Contacts.AddToAllContacts:input_type -> chats.AddToAllReq
	4, // 8: chats.Contacts.GetContacts:output_type -> chats.PersonArray
	6, // 9: chats.Contacts.AddContactByUsername:output_type -> chats.EmptyContacts
	7, // 10: chats.Contacts.AddToAllContacts:output_type -> chats.BoolResponseContacts
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_contacts_proto_init() }
func file_contacts_proto_init() {
	if File_contacts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_contacts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIDContacts); i {
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
		file_contacts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserIDArray); i {
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
		file_contacts_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddToAllReq); i {
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
		file_contacts_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Person); i {
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
		file_contacts_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersonArray); i {
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
		file_contacts_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddByUsernameReq); i {
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
		file_contacts_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyContacts); i {
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
		file_contacts_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoolResponseContacts); i {
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
			RawDescriptor: file_contacts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_contacts_proto_goTypes,
		DependencyIndexes: file_contacts_proto_depIdxs,
		MessageInfos:      file_contacts_proto_msgTypes,
	}.Build()
	File_contacts_proto = out.File
	file_contacts_proto_rawDesc = nil
	file_contacts_proto_goTypes = nil
	file_contacts_proto_depIdxs = nil
}
