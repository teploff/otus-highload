// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messenger.proto

package messenger

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type CreateChatRequest struct {
	MasterToken          string   `protobuf:"bytes,1,opt,name=master_token,json=masterToken,proto3" json:"master_token,omitempty"`
	SlaveId              string   `protobuf:"bytes,2,opt,name=slave_id,json=slaveId,proto3" json:"slave_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateChatRequest) Reset()         { *m = CreateChatRequest{} }
func (m *CreateChatRequest) String() string { return proto.CompactTextString(m) }
func (*CreateChatRequest) ProtoMessage()    {}
func (*CreateChatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{0}
}

func (m *CreateChatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateChatRequest.Unmarshal(m, b)
}
func (m *CreateChatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateChatRequest.Marshal(b, m, deterministic)
}
func (m *CreateChatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateChatRequest.Merge(m, src)
}
func (m *CreateChatRequest) XXX_Size() int {
	return xxx_messageInfo_CreateChatRequest.Size(m)
}
func (m *CreateChatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateChatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateChatRequest proto.InternalMessageInfo

func (m *CreateChatRequest) GetMasterToken() string {
	if m != nil {
		return m.MasterToken
	}
	return ""
}

func (m *CreateChatRequest) GetSlaveId() string {
	if m != nil {
		return m.SlaveId
	}
	return ""
}

type CreateChatResponse struct {
	ChatId               string   `protobuf:"bytes,1,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateChatResponse) Reset()         { *m = CreateChatResponse{} }
func (m *CreateChatResponse) String() string { return proto.CompactTextString(m) }
func (*CreateChatResponse) ProtoMessage()    {}
func (*CreateChatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{1}
}

func (m *CreateChatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateChatResponse.Unmarshal(m, b)
}
func (m *CreateChatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateChatResponse.Marshal(b, m, deterministic)
}
func (m *CreateChatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateChatResponse.Merge(m, src)
}
func (m *CreateChatResponse) XXX_Size() int {
	return xxx_messageInfo_CreateChatResponse.Size(m)
}
func (m *CreateChatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateChatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateChatResponse proto.InternalMessageInfo

func (m *CreateChatResponse) GetChatId() string {
	if m != nil {
		return m.ChatId
	}
	return ""
}

type GetChatsRequest struct {
	UserToken            string               `protobuf:"bytes,1,opt,name=user_token,json=userToken,proto3" json:"user_token,omitempty"`
	Offset               *wrappers.Int32Value `protobuf:"bytes,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                *wrappers.Int32Value `protobuf:"bytes,3,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GetChatsRequest) Reset()         { *m = GetChatsRequest{} }
func (m *GetChatsRequest) String() string { return proto.CompactTextString(m) }
func (*GetChatsRequest) ProtoMessage()    {}
func (*GetChatsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{2}
}

func (m *GetChatsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetChatsRequest.Unmarshal(m, b)
}
func (m *GetChatsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetChatsRequest.Marshal(b, m, deterministic)
}
func (m *GetChatsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetChatsRequest.Merge(m, src)
}
func (m *GetChatsRequest) XXX_Size() int {
	return xxx_messageInfo_GetChatsRequest.Size(m)
}
func (m *GetChatsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetChatsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetChatsRequest proto.InternalMessageInfo

func (m *GetChatsRequest) GetUserToken() string {
	if m != nil {
		return m.UserToken
	}
	return ""
}

func (m *GetChatsRequest) GetOffset() *wrappers.Int32Value {
	if m != nil {
		return m.Offset
	}
	return nil
}

func (m *GetChatsRequest) GetLimit() *wrappers.Int32Value {
	if m != nil {
		return m.Limit
	}
	return nil
}

type Participant struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Surname              string   `protobuf:"bytes,3,opt,name=surname,proto3" json:"surname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Participant) Reset()         { *m = Participant{} }
func (m *Participant) String() string { return proto.CompactTextString(m) }
func (*Participant) ProtoMessage()    {}
func (*Participant) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{3}
}

func (m *Participant) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Participant.Unmarshal(m, b)
}
func (m *Participant) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Participant.Marshal(b, m, deterministic)
}
func (m *Participant) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Participant.Merge(m, src)
}
func (m *Participant) XXX_Size() int {
	return xxx_messageInfo_Participant.Size(m)
}
func (m *Participant) XXX_DiscardUnknown() {
	xxx_messageInfo_Participant.DiscardUnknown(m)
}

var xxx_messageInfo_Participant proto.InternalMessageInfo

func (m *Participant) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Participant) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Participant) GetSurname() string {
	if m != nil {
		return m.Surname
	}
	return ""
}

type Chat struct {
	Id                   string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreateTime           int64          `protobuf:"varint,2,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	Participants         []*Participant `protobuf:"bytes,3,rep,name=participants,proto3" json:"participants,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Chat) Reset()         { *m = Chat{} }
func (m *Chat) String() string { return proto.CompactTextString(m) }
func (*Chat) ProtoMessage()    {}
func (*Chat) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{4}
}

func (m *Chat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Chat.Unmarshal(m, b)
}
func (m *Chat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Chat.Marshal(b, m, deterministic)
}
func (m *Chat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Chat.Merge(m, src)
}
func (m *Chat) XXX_Size() int {
	return xxx_messageInfo_Chat.Size(m)
}
func (m *Chat) XXX_DiscardUnknown() {
	xxx_messageInfo_Chat.DiscardUnknown(m)
}

var xxx_messageInfo_Chat proto.InternalMessageInfo

func (m *Chat) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Chat) GetCreateTime() int64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *Chat) GetParticipants() []*Participant {
	if m != nil {
		return m.Participants
	}
	return nil
}

type GetChatsResponse struct {
	Total                int32    `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Offset               int32    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int32    `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Chats                []*Chat  `protobuf:"bytes,4,rep,name=chats,proto3" json:"chats,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetChatsResponse) Reset()         { *m = GetChatsResponse{} }
func (m *GetChatsResponse) String() string { return proto.CompactTextString(m) }
func (*GetChatsResponse) ProtoMessage()    {}
func (*GetChatsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{5}
}

func (m *GetChatsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetChatsResponse.Unmarshal(m, b)
}
func (m *GetChatsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetChatsResponse.Marshal(b, m, deterministic)
}
func (m *GetChatsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetChatsResponse.Merge(m, src)
}
func (m *GetChatsResponse) XXX_Size() int {
	return xxx_messageInfo_GetChatsResponse.Size(m)
}
func (m *GetChatsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetChatsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetChatsResponse proto.InternalMessageInfo

func (m *GetChatsResponse) GetTotal() int32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *GetChatsResponse) GetOffset() int32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *GetChatsResponse) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *GetChatsResponse) GetChats() []*Chat {
	if m != nil {
		return m.Chats
	}
	return nil
}

type GetMessagesRequest struct {
	UserToken            string               `protobuf:"bytes,1,opt,name=user_token,json=userToken,proto3" json:"user_token,omitempty"`
	ChatId               string               `protobuf:"bytes,2,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
	Offset               *wrappers.Int32Value `protobuf:"bytes,3,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                *wrappers.Int32Value `protobuf:"bytes,4,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GetMessagesRequest) Reset()         { *m = GetMessagesRequest{} }
func (m *GetMessagesRequest) String() string { return proto.CompactTextString(m) }
func (*GetMessagesRequest) ProtoMessage()    {}
func (*GetMessagesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{6}
}

func (m *GetMessagesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMessagesRequest.Unmarshal(m, b)
}
func (m *GetMessagesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMessagesRequest.Marshal(b, m, deterministic)
}
func (m *GetMessagesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMessagesRequest.Merge(m, src)
}
func (m *GetMessagesRequest) XXX_Size() int {
	return xxx_messageInfo_GetMessagesRequest.Size(m)
}
func (m *GetMessagesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMessagesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetMessagesRequest proto.InternalMessageInfo

func (m *GetMessagesRequest) GetUserToken() string {
	if m != nil {
		return m.UserToken
	}
	return ""
}

func (m *GetMessagesRequest) GetChatId() string {
	if m != nil {
		return m.ChatId
	}
	return ""
}

func (m *GetMessagesRequest) GetOffset() *wrappers.Int32Value {
	if m != nil {
		return m.Offset
	}
	return nil
}

func (m *GetMessagesRequest) GetLimit() *wrappers.Int32Value {
	if m != nil {
		return m.Limit
	}
	return nil
}

type Message struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Text                 string   `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Status               string   `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	CreateTime           int64    `protobuf:"varint,4,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	UserId               string   `protobuf:"bytes,5,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ChatId               string   `protobuf:"bytes,6,opt,name=chat_id,json=chatId,proto3" json:"chat_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{7}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *Message) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Message) GetCreateTime() int64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *Message) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *Message) GetChatId() string {
	if m != nil {
		return m.ChatId
	}
	return ""
}

type GetMessagesResponse struct {
	Total                int32      `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Offset               int32      `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit                int32      `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Messages             []*Message `protobuf:"bytes,4,rep,name=messages,proto3" json:"messages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *GetMessagesResponse) Reset()         { *m = GetMessagesResponse{} }
func (m *GetMessagesResponse) String() string { return proto.CompactTextString(m) }
func (*GetMessagesResponse) ProtoMessage()    {}
func (*GetMessagesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b99aba0cbf4e4b91, []int{8}
}

func (m *GetMessagesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMessagesResponse.Unmarshal(m, b)
}
func (m *GetMessagesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMessagesResponse.Marshal(b, m, deterministic)
}
func (m *GetMessagesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMessagesResponse.Merge(m, src)
}
func (m *GetMessagesResponse) XXX_Size() int {
	return xxx_messageInfo_GetMessagesResponse.Size(m)
}
func (m *GetMessagesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMessagesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetMessagesResponse proto.InternalMessageInfo

func (m *GetMessagesResponse) GetTotal() int32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *GetMessagesResponse) GetOffset() int32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *GetMessagesResponse) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *GetMessagesResponse) GetMessages() []*Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

func init() {
	proto.RegisterType((*CreateChatRequest)(nil), "messenger.CreateChatRequest")
	proto.RegisterType((*CreateChatResponse)(nil), "messenger.CreateChatResponse")
	proto.RegisterType((*GetChatsRequest)(nil), "messenger.GetChatsRequest")
	proto.RegisterType((*Participant)(nil), "messenger.Participant")
	proto.RegisterType((*Chat)(nil), "messenger.Chat")
	proto.RegisterType((*GetChatsResponse)(nil), "messenger.GetChatsResponse")
	proto.RegisterType((*GetMessagesRequest)(nil), "messenger.GetMessagesRequest")
	proto.RegisterType((*Message)(nil), "messenger.Message")
	proto.RegisterType((*GetMessagesResponse)(nil), "messenger.GetMessagesResponse")
}

func init() {
	proto.RegisterFile("messenger.proto", fileDescriptor_b99aba0cbf4e4b91)
}

var fileDescriptor_b99aba0cbf4e4b91 = []byte{
	// 552 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0xfd, 0x1c, 0xc7, 0xf9, 0x19, 0x57, 0x5f, 0x60, 0xa9, 0x5a, 0x93, 0x92, 0x52, 0x2c, 0x21,
	0xf5, 0x06, 0x57, 0x24, 0x77, 0xdc, 0x56, 0xa8, 0x8a, 0xaa, 0x22, 0xb0, 0x2a, 0x6e, 0xa3, 0x6d,
	0x3c, 0x49, 0x2d, 0xe2, 0x1f, 0xbc, 0x63, 0xe0, 0x86, 0x07, 0xe0, 0x0d, 0xfa, 0x20, 0xbc, 0x17,
	0xaf, 0x80, 0xbc, 0xbb, 0x4e, 0x36, 0x6e, 0x2b, 0x2a, 0xc4, 0x9d, 0x67, 0xe6, 0x78, 0xe6, 0xcc,
	0x39, 0x63, 0xc3, 0x20, 0x41, 0x21, 0x30, 0x5d, 0x62, 0x11, 0xe4, 0x45, 0x46, 0x19, 0xeb, 0xaf,
	0x13, 0xc3, 0xc3, 0x65, 0x96, 0x2d, 0x57, 0x78, 0x22, 0x0b, 0x57, 0xe5, 0xe2, 0xe4, 0x6b, 0xc1,
	0xf3, 0x1c, 0x0b, 0xa1, 0xa0, 0xfe, 0x07, 0x78, 0x7c, 0x5a, 0x20, 0x27, 0x3c, 0xbd, 0xe6, 0x14,
	0xe2, 0xe7, 0x12, 0x05, 0xb1, 0x17, 0xb0, 0x93, 0x70, 0x41, 0x58, 0xcc, 0x28, 0xfb, 0x84, 0xa9,
	0x67, 0x1d, 0x59, 0xc7, 0xfd, 0xd0, 0x55, 0xb9, 0xcb, 0x2a, 0xc5, 0x9e, 0x42, 0x4f, 0xac, 0xf8,
	0x17, 0x9c, 0xc5, 0x91, 0xd7, 0x92, 0xe5, 0xae, 0x8c, 0xa7, 0x91, 0xff, 0x0a, 0x98, 0xd9, 0x52,
	0xe4, 0x59, 0x2a, 0x90, 0xed, 0x43, 0x77, 0x7e, 0xcd, 0xa9, 0xc2, 0xab, 0x76, 0x9d, 0x2a, 0x9c,
	0x46, 0xfe, 0x8d, 0x05, 0x83, 0x33, 0xa4, 0x0a, 0x2c, 0x6a, 0x02, 0x23, 0x80, 0x52, 0x34, 0xc6,
	0xf7, 0xab, 0x8c, 0x1a, 0x3e, 0x81, 0x4e, 0xb6, 0x58, 0x08, 0x24, 0x39, 0xda, 0x1d, 0x1f, 0x04,
	0x6a, 0xcb, 0xa0, 0xde, 0x32, 0x98, 0xa6, 0x34, 0x19, 0x7f, 0xe4, 0xab, 0x12, 0x43, 0x0d, 0x65,
	0xaf, 0xc1, 0x59, 0xc5, 0x49, 0x4c, 0x9e, 0xfd, 0xe7, 0x77, 0x14, 0xd2, 0x3f, 0x07, 0xf7, 0x3d,
	0x2f, 0x28, 0x9e, 0xc7, 0x39, 0x4f, 0x89, 0xfd, 0x0f, 0xad, 0x35, 0xfb, 0x56, 0x1c, 0x31, 0x06,
	0xed, 0x94, 0x27, 0xa8, 0xf7, 0x97, 0xcf, 0xcc, 0x83, 0xae, 0x28, 0x0b, 0x99, 0xb6, 0xb5, 0x2c,
	0x2a, 0xf4, 0x05, 0xb4, 0xab, 0x1d, 0x6f, 0x75, 0x79, 0x0e, 0xee, 0x5c, 0xca, 0x35, 0xa3, 0x58,
	0x37, 0xb3, 0x43, 0x50, 0xa9, 0xcb, 0x38, 0x41, 0xf6, 0x06, 0x76, 0xf2, 0x0d, 0x0b, 0xe1, 0xd9,
	0x47, 0xf6, 0xb1, 0x3b, 0xde, 0x0b, 0x36, 0xae, 0x1b, 0x24, 0xc3, 0x2d, 0xac, 0xff, 0x1d, 0x1e,
	0x6d, 0xb4, 0xd5, 0x4e, 0xec, 0x82, 0x43, 0x19, 0xf1, 0x95, 0xe4, 0xe0, 0x84, 0x2a, 0x60, 0x7b,
	0x5b, 0x9a, 0x3a, 0x6b, 0xd9, 0x76, 0x4d, 0xd9, 0x1c, 0xad, 0x0c, 0x7b, 0x09, 0x4e, 0x65, 0x9f,
	0xf0, 0xda, 0x92, 0xcc, 0xc0, 0x20, 0x23, 0x5d, 0x57, 0x55, 0xff, 0xa7, 0x05, 0xec, 0x0c, 0xe9,
	0x02, 0x85, 0xe0, 0x4b, 0x7c, 0xa8, 0xbd, 0xc6, 0xa9, 0xb4, 0xcc, 0x53, 0x31, 0x7c, 0xb7, 0xff,
	0xc2, 0xf7, 0xf6, 0x83, 0x7d, 0xbf, 0xb1, 0xa0, 0xab, 0x39, 0xdf, 0x65, 0x3a, 0xe1, 0x37, 0xaa,
	0x4d, 0xaf, 0x9e, 0x2b, 0xed, 0x04, 0x71, 0x2a, 0x85, 0xf6, 0x5c, 0x47, 0x4d, 0x6b, 0xdb, 0xb7,
	0xac, 0xdd, 0x87, 0xae, 0x14, 0x22, 0x8e, 0x3c, 0x47, 0xbd, 0x59, 0x85, 0xd3, 0xc8, 0x94, 0xa0,
	0xb3, 0xf5, 0xb5, 0xfc, 0xb0, 0xe0, 0xc9, 0x96, 0xa2, 0xff, 0xd0, 0xd4, 0x00, 0x7a, 0x89, 0xee,
	0xab, 0x7d, 0x65, 0x86, 0xaf, 0x7a, 0x64, 0xb8, 0xc6, 0x8c, 0x7f, 0x59, 0xd0, 0xbf, 0xa8, 0xeb,
	0xec, 0x1c, 0x60, 0xf3, 0xd9, 0xb3, 0x67, 0xe6, 0x45, 0x34, 0x7f, 0x30, 0xc3, 0xd1, 0x3d, 0x55,
	0xb5, 0x8c, 0xff, 0x1f, 0x7b, 0x0b, 0xbd, 0xfa, 0x6e, 0xd9, 0xd0, 0x00, 0x37, 0x7e, 0x14, 0xc3,
	0x83, 0x3b, 0x6b, 0xeb, 0x36, 0xef, 0xc0, 0x35, 0xc4, 0x62, 0xa3, 0x6d, 0x74, 0xe3, 0x2c, 0x87,
	0x87, 0xf7, 0x95, 0xeb, 0x7e, 0x57, 0x1d, 0x79, 0x34, 0x93, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x61, 0xf6, 0xf9, 0x0a, 0x72, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MessengerClient is the client API for Messenger service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MessengerClient interface {
	CreateChat(ctx context.Context, in *CreateChatRequest, opts ...grpc.CallOption) (*CreateChatResponse, error)
	GetChats(ctx context.Context, in *GetChatsRequest, opts ...grpc.CallOption) (*GetChatsResponse, error)
	GetMessages(ctx context.Context, in *GetMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error)
}

type messengerClient struct {
	cc grpc.ClientConnInterface
}

func NewMessengerClient(cc grpc.ClientConnInterface) MessengerClient {
	return &messengerClient{cc}
}

func (c *messengerClient) CreateChat(ctx context.Context, in *CreateChatRequest, opts ...grpc.CallOption) (*CreateChatResponse, error) {
	out := new(CreateChatResponse)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/CreateChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerClient) GetChats(ctx context.Context, in *GetChatsRequest, opts ...grpc.CallOption) (*GetChatsResponse, error) {
	out := new(GetChatsResponse)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/GetChats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerClient) GetMessages(ctx context.Context, in *GetMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error) {
	out := new(GetMessagesResponse)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/GetMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessengerServer is the server API for Messenger service.
type MessengerServer interface {
	CreateChat(context.Context, *CreateChatRequest) (*CreateChatResponse, error)
	GetChats(context.Context, *GetChatsRequest) (*GetChatsResponse, error)
	GetMessages(context.Context, *GetMessagesRequest) (*GetMessagesResponse, error)
}

// UnimplementedMessengerServer can be embedded to have forward compatible implementations.
type UnimplementedMessengerServer struct {
}

func (*UnimplementedMessengerServer) CreateChat(ctx context.Context, req *CreateChatRequest) (*CreateChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateChat not implemented")
}
func (*UnimplementedMessengerServer) GetChats(ctx context.Context, req *GetChatsRequest) (*GetChatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChats not implemented")
}
func (*UnimplementedMessengerServer) GetMessages(ctx context.Context, req *GetMessagesRequest) (*GetMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessages not implemented")
}

func RegisterMessengerServer(s *grpc.Server, srv MessengerServer) {
	s.RegisterService(&_Messenger_serviceDesc, srv)
}

func _Messenger_CreateChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).CreateChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/CreateChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).CreateChat(ctx, req.(*CreateChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messenger_GetChats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).GetChats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/GetChats",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).GetChats(ctx, req.(*GetChatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messenger_GetMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).GetMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/GetMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).GetMessages(ctx, req.(*GetMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Messenger_serviceDesc = grpc.ServiceDesc{
	ServiceName: "messenger.Messenger",
	HandlerType: (*MessengerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateChat",
			Handler:    _Messenger_CreateChat_Handler,
		},
		{
			MethodName: "GetChats",
			Handler:    _Messenger_GetChats_Handler,
		},
		{
			MethodName: "GetMessages",
			Handler:    _Messenger_GetMessages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messenger.proto",
}
