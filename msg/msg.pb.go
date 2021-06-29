// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg.proto

package msg

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

// 消息ID
type MessageID int32

const (
	MessageID_MSG_Ping              MessageID = 0
	MessageID_MSG_Pong              MessageID = 1
	MessageID_MSG_Login_C2S         MessageID = 2
	MessageID_MSG_Login_S2C         MessageID = 3
	MessageID_MSG_Logout_C2S        MessageID = 4
	MessageID_MSG_Logout_S2C        MessageID = 5
	MessageID_MSG_JoinRoom_C2S      MessageID = 6
	MessageID_MSG_JoinRoom_S2C      MessageID = 7
	MessageID_MSG_EnterRoom_S2C     MessageID = 8
	MessageID_MSG_PlayerAction_C2S  MessageID = 9
	MessageID_MSG_PlayerAction_SC2  MessageID = 10
	MessageID_MSG_SendWinMoney_C2S  MessageID = 11
	MessageID_MSG_SendWinMoney_S2C  MessageID = 12
	MessageID_MSG_GetRewards_C2S    MessageID = 13
	MessageID_MSG_GetRewards_S2C    MessageID = 14
	MessageID_MSG_ChangeRoomCfg_C2S MessageID = 15
	MessageID_MSG_ChangeRoomCfg_S2C MessageID = 16
	MessageID_MSG_ErrorMsg_S2C      MessageID = 17
)

var MessageID_name = map[int32]string{
	0:  "MSG_Ping",
	1:  "MSG_Pong",
	2:  "MSG_Login_C2S",
	3:  "MSG_Login_S2C",
	4:  "MSG_Logout_C2S",
	5:  "MSG_Logout_S2C",
	6:  "MSG_JoinRoom_C2S",
	7:  "MSG_JoinRoom_S2C",
	8:  "MSG_EnterRoom_S2C",
	9:  "MSG_PlayerAction_C2S",
	10: "MSG_PlayerAction_SC2",
	11: "MSG_SendWinMoney_C2S",
	12: "MSG_SendWinMoney_S2C",
	13: "MSG_GetRewards_C2S",
	14: "MSG_GetRewards_S2C",
	15: "MSG_ChangeRoomCfg_C2S",
	16: "MSG_ChangeRoomCfg_S2C",
	17: "MSG_ErrorMsg_S2C",
}

var MessageID_value = map[string]int32{
	"MSG_Ping":              0,
	"MSG_Pong":              1,
	"MSG_Login_C2S":         2,
	"MSG_Login_S2C":         3,
	"MSG_Logout_C2S":        4,
	"MSG_Logout_S2C":        5,
	"MSG_JoinRoom_C2S":      6,
	"MSG_JoinRoom_S2C":      7,
	"MSG_EnterRoom_S2C":     8,
	"MSG_PlayerAction_C2S":  9,
	"MSG_PlayerAction_SC2":  10,
	"MSG_SendWinMoney_C2S":  11,
	"MSG_SendWinMoney_S2C":  12,
	"MSG_GetRewards_C2S":    13,
	"MSG_GetRewards_S2C":    14,
	"MSG_ChangeRoomCfg_C2S": 15,
	"MSG_ChangeRoomCfg_S2C": 16,
	"MSG_ErrorMsg_S2C":      17,
}

func (x MessageID) String() string {
	return proto.EnumName(MessageID_name, int32(x))
}

func (MessageID) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

type Ping struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ping) Reset()         { *m = Ping{} }
func (m *Ping) String() string { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()    {}
func (*Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

func (m *Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ping.Unmarshal(m, b)
}
func (m *Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ping.Marshal(b, m, deterministic)
}
func (m *Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ping.Merge(m, src)
}
func (m *Ping) XXX_Size() int {
	return xxx_messageInfo_Ping.Size(m)
}
func (m *Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Ping proto.InternalMessageInfo

type Pong struct {
	ServerTime           int64    `protobuf:"varint,1,opt,name=serverTime,proto3" json:"serverTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pong) Reset()         { *m = Pong{} }
func (m *Pong) String() string { return proto.CompactTextString(m) }
func (*Pong) ProtoMessage()    {}
func (*Pong) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{1}
}

func (m *Pong) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pong.Unmarshal(m, b)
}
func (m *Pong) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pong.Marshal(b, m, deterministic)
}
func (m *Pong) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pong.Merge(m, src)
}
func (m *Pong) XXX_Size() int {
	return xxx_messageInfo_Pong.Size(m)
}
func (m *Pong) XXX_DiscardUnknown() {
	xxx_messageInfo_Pong.DiscardUnknown(m)
}

var xxx_messageInfo_Pong proto.InternalMessageInfo

func (m *Pong) GetServerTime() int64 {
	if m != nil {
		return m.ServerTime
	}
	return 0
}

type PlayerInfo struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	NickName             string   `protobuf:"bytes,2,opt,name=nickName,proto3" json:"nickName,omitempty"`
	HeadImg              string   `protobuf:"bytes,3,opt,name=headImg,proto3" json:"headImg,omitempty"`
	Account              float64  `protobuf:"fixed64,4,opt,name=account,proto3" json:"account,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerInfo) Reset()         { *m = PlayerInfo{} }
func (m *PlayerInfo) String() string { return proto.CompactTextString(m) }
func (*PlayerInfo) ProtoMessage()    {}
func (*PlayerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{2}
}

func (m *PlayerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerInfo.Unmarshal(m, b)
}
func (m *PlayerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerInfo.Marshal(b, m, deterministic)
}
func (m *PlayerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerInfo.Merge(m, src)
}
func (m *PlayerInfo) XXX_Size() int {
	return xxx_messageInfo_PlayerInfo.Size(m)
}
func (m *PlayerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerInfo proto.InternalMessageInfo

func (m *PlayerInfo) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *PlayerInfo) GetNickName() string {
	if m != nil {
		return m.NickName
	}
	return ""
}

func (m *PlayerInfo) GetHeadImg() string {
	if m != nil {
		return m.HeadImg
	}
	return ""
}

func (m *PlayerInfo) GetAccount() float64 {
	if m != nil {
		return m.Account
	}
	return 0
}

type Login_C2S struct {
	Id                   string   `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	PassWord             string   `protobuf:"bytes,2,opt,name=PassWord,proto3" json:"PassWord,omitempty"`
	Token                string   `protobuf:"bytes,3,opt,name=Token,proto3" json:"Token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Login_C2S) Reset()         { *m = Login_C2S{} }
func (m *Login_C2S) String() string { return proto.CompactTextString(m) }
func (*Login_C2S) ProtoMessage()    {}
func (*Login_C2S) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{3}
}

func (m *Login_C2S) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Login_C2S.Unmarshal(m, b)
}
func (m *Login_C2S) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Login_C2S.Marshal(b, m, deterministic)
}
func (m *Login_C2S) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Login_C2S.Merge(m, src)
}
func (m *Login_C2S) XXX_Size() int {
	return xxx_messageInfo_Login_C2S.Size(m)
}
func (m *Login_C2S) XXX_DiscardUnknown() {
	xxx_messageInfo_Login_C2S.DiscardUnknown(m)
}

var xxx_messageInfo_Login_C2S proto.InternalMessageInfo

func (m *Login_C2S) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Login_C2S) GetPassWord() string {
	if m != nil {
		return m.PassWord
	}
	return ""
}

func (m *Login_C2S) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type Login_S2C struct {
	PlayerInfo           *PlayerInfo `protobuf:"bytes,1,opt,name=playerInfo,proto3" json:"playerInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Login_S2C) Reset()         { *m = Login_S2C{} }
func (m *Login_S2C) String() string { return proto.CompactTextString(m) }
func (*Login_S2C) ProtoMessage()    {}
func (*Login_S2C) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{4}
}

func (m *Login_S2C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Login_S2C.Unmarshal(m, b)
}
func (m *Login_S2C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Login_S2C.Marshal(b, m, deterministic)
}
func (m *Login_S2C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Login_S2C.Merge(m, src)
}
func (m *Login_S2C) XXX_Size() int {
	return xxx_messageInfo_Login_S2C.Size(m)
}
func (m *Login_S2C) XXX_DiscardUnknown() {
	xxx_messageInfo_Login_S2C.DiscardUnknown(m)
}

var xxx_messageInfo_Login_S2C proto.InternalMessageInfo

func (m *Login_S2C) GetPlayerInfo() *PlayerInfo {
	if m != nil {
		return m.PlayerInfo
	}
	return nil
}

type Logout_C2S struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Logout_C2S) Reset()         { *m = Logout_C2S{} }
func (m *Logout_C2S) String() string { return proto.CompactTextString(m) }
func (*Logout_C2S) ProtoMessage()    {}
func (*Logout_C2S) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{5}
}

func (m *Logout_C2S) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Logout_C2S.Unmarshal(m, b)
}
func (m *Logout_C2S) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Logout_C2S.Marshal(b, m, deterministic)
}
func (m *Logout_C2S) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Logout_C2S.Merge(m, src)
}
func (m *Logout_C2S) XXX_Size() int {
	return xxx_messageInfo_Logout_C2S.Size(m)
}
func (m *Logout_C2S) XXX_DiscardUnknown() {
	xxx_messageInfo_Logout_C2S.DiscardUnknown(m)
}

var xxx_messageInfo_Logout_C2S proto.InternalMessageInfo

type Logout_S2C struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Logout_S2C) Reset()         { *m = Logout_S2C{} }
func (m *Logout_S2C) String() string { return proto.CompactTextString(m) }
func (*Logout_S2C) ProtoMessage()    {}
func (*Logout_S2C) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{6}
}

func (m *Logout_S2C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Logout_S2C.Unmarshal(m, b)
}
func (m *Logout_S2C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Logout_S2C.Marshal(b, m, deterministic)
}
func (m *Logout_S2C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Logout_S2C.Merge(m, src)
}
func (m *Logout_S2C) XXX_Size() int {
	return xxx_messageInfo_Logout_S2C.Size(m)
}
func (m *Logout_S2C) XXX_DiscardUnknown() {
	xxx_messageInfo_Logout_S2C.DiscardUnknown(m)
}

var xxx_messageInfo_Logout_S2C proto.InternalMessageInfo

type PlayerData struct {
	PlayerInfo           *PlayerInfo `protobuf:"bytes,1,opt,name=playerInfo,proto3" json:"playerInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *PlayerData) Reset()         { *m = PlayerData{} }
func (m *PlayerData) String() string { return proto.CompactTextString(m) }
func (*PlayerData) ProtoMessage()    {}
func (*PlayerData) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{7}
}

func (m *PlayerData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerData.Unmarshal(m, b)
}
func (m *PlayerData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerData.Marshal(b, m, deterministic)
}
func (m *PlayerData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerData.Merge(m, src)
}
func (m *PlayerData) XXX_Size() int {
	return xxx_messageInfo_PlayerData.Size(m)
}
func (m *PlayerData) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerData.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerData proto.InternalMessageInfo

func (m *PlayerData) GetPlayerInfo() *PlayerInfo {
	if m != nil {
		return m.PlayerInfo
	}
	return nil
}

type RoomData struct {
	RoomId               string        `protobuf:"bytes,1,opt,name=roomId,proto3" json:"roomId,omitempty"`
	PlayerData           []*PlayerData `protobuf:"bytes,2,rep,name=playerData,proto3" json:"playerData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RoomData) Reset()         { *m = RoomData{} }
func (m *RoomData) String() string { return proto.CompactTextString(m) }
func (*RoomData) ProtoMessage()    {}
func (*RoomData) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{8}
}

func (m *RoomData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RoomData.Unmarshal(m, b)
}
func (m *RoomData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RoomData.Marshal(b, m, deterministic)
}
func (m *RoomData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoomData.Merge(m, src)
}
func (m *RoomData) XXX_Size() int {
	return xxx_messageInfo_RoomData.Size(m)
}
func (m *RoomData) XXX_DiscardUnknown() {
	xxx_messageInfo_RoomData.DiscardUnknown(m)
}

var xxx_messageInfo_RoomData proto.InternalMessageInfo

func (m *RoomData) GetRoomId() string {
	if m != nil {
		return m.RoomId
	}
	return ""
}

func (m *RoomData) GetPlayerData() []*PlayerData {
	if m != nil {
		return m.PlayerData
	}
	return nil
}

type JoinRoom_C2S struct {
	Cfg                  string   `protobuf:"bytes,1,opt,name=cfg,proto3" json:"cfg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinRoom_C2S) Reset()         { *m = JoinRoom_C2S{} }
func (m *JoinRoom_C2S) String() string { return proto.CompactTextString(m) }
func (*JoinRoom_C2S) ProtoMessage()    {}
func (*JoinRoom_C2S) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{9}
}

func (m *JoinRoom_C2S) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinRoom_C2S.Unmarshal(m, b)
}
func (m *JoinRoom_C2S) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinRoom_C2S.Marshal(b, m, deterministic)
}
func (m *JoinRoom_C2S) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinRoom_C2S.Merge(m, src)
}
func (m *JoinRoom_C2S) XXX_Size() int {
	return xxx_messageInfo_JoinRoom_C2S.Size(m)
}
func (m *JoinRoom_C2S) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinRoom_C2S.DiscardUnknown(m)
}

var xxx_messageInfo_JoinRoom_C2S proto.InternalMessageInfo

func (m *JoinRoom_C2S) GetCfg() string {
	if m != nil {
		return m.Cfg
	}
	return ""
}

type JoinRoom_S2C struct {
	RoomData             *RoomData `protobuf:"bytes,1,opt,name=roomData,proto3" json:"roomData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *JoinRoom_S2C) Reset()         { *m = JoinRoom_S2C{} }
func (m *JoinRoom_S2C) String() string { return proto.CompactTextString(m) }
func (*JoinRoom_S2C) ProtoMessage()    {}
func (*JoinRoom_S2C) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{10}
}

func (m *JoinRoom_S2C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinRoom_S2C.Unmarshal(m, b)
}
func (m *JoinRoom_S2C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinRoom_S2C.Marshal(b, m, deterministic)
}
func (m *JoinRoom_S2C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinRoom_S2C.Merge(m, src)
}
func (m *JoinRoom_S2C) XXX_Size() int {
	return xxx_messageInfo_JoinRoom_S2C.Size(m)
}
func (m *JoinRoom_S2C) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinRoom_S2C.DiscardUnknown(m)
}

var xxx_messageInfo_JoinRoom_S2C proto.InternalMessageInfo

func (m *JoinRoom_S2C) GetRoomData() *RoomData {
	if m != nil {
		return m.RoomData
	}
	return nil
}

type EnterRoom_S2C struct {
	RoomData             *RoomData `protobuf:"bytes,1,opt,name=roomData,proto3" json:"roomData,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *EnterRoom_S2C) Reset()         { *m = EnterRoom_S2C{} }
func (m *EnterRoom_S2C) String() string { return proto.CompactTextString(m) }
func (*EnterRoom_S2C) ProtoMessage()    {}
func (*EnterRoom_S2C) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{11}
}

func (m *EnterRoom_S2C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnterRoom_S2C.Unmarshal(m, b)
}
func (m *EnterRoom_S2C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnterRoom_S2C.Marshal(b, m, deterministic)
}
func (m *EnterRoom_S2C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnterRoom_S2C.Merge(m, src)
}
func (m *EnterRoom_S2C) XXX_Size() int {
	return xxx_messageInfo_EnterRoom_S2C.Size(m)
}
func (m *EnterRoom_S2C) XXX_DiscardUnknown() {
	xxx_messageInfo_EnterRoom_S2C.DiscardUnknown(m)
}

var xxx_messageInfo_EnterRoom_S2C proto.InternalMessageInfo

func (m *EnterRoom_S2C) GetRoomData() *RoomData {
	if m != nil {
		return m.RoomData
	}
	return nil
}

type PlayerAction_C2S struct {
	DownBet              float64  `protobuf:"fixed64,1,opt,name=downBet,proto3" json:"downBet,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerAction_C2S) Reset()         { *m = PlayerAction_C2S{} }
func (m *PlayerAction_C2S) String() string { return proto.CompactTextString(m) }
func (*PlayerAction_C2S) ProtoMessage()    {}
func (*PlayerAction_C2S) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{12}
}

func (m *PlayerAction_C2S) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerAction_C2S.Unmarshal(m, b)
}
func (m *PlayerAction_C2S) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerAction_C2S.Marshal(b, m, deterministic)
}
func (m *PlayerAction_C2S) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerAction_C2S.Merge(m, src)
}
func (m *PlayerAction_C2S) XXX_Size() int {
	return xxx_messageInfo_PlayerAction_C2S.Size(m)
}
func (m *PlayerAction_C2S) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerAction_C2S.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerAction_C2S proto.InternalMessageInfo

func (m *PlayerAction_C2S) GetDownBet() float64 {
	if m != nil {
		return m.DownBet
	}
	return 0
}

type PlayerAction_S2C struct {
	IsWin                bool     `protobuf:"varint,1,opt,name=IsWin,proto3" json:"IsWin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PlayerAction_S2C) Reset()         { *m = PlayerAction_S2C{} }
func (m *PlayerAction_S2C) String() string { return proto.CompactTextString(m) }
func (*PlayerAction_S2C) ProtoMessage()    {}
func (*PlayerAction_S2C) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{13}
}

func (m *PlayerAction_S2C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PlayerAction_S2C.Unmarshal(m, b)
}
func (m *PlayerAction_S2C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PlayerAction_S2C.Marshal(b, m, deterministic)
}
func (m *PlayerAction_S2C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PlayerAction_S2C.Merge(m, src)
}
func (m *PlayerAction_S2C) XXX_Size() int {
	return xxx_messageInfo_PlayerAction_S2C.Size(m)
}
func (m *PlayerAction_S2C) XXX_DiscardUnknown() {
	xxx_messageInfo_PlayerAction_S2C.DiscardUnknown(m)
}

var xxx_messageInfo_PlayerAction_S2C proto.InternalMessageInfo

func (m *PlayerAction_S2C) GetIsWin() bool {
	if m != nil {
		return m.IsWin
	}
	return false
}

type SendWinMoney_C2S struct {
	Money                float64  `protobuf:"fixed64,1,opt,name=money,proto3" json:"money,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendWinMoney_C2S) Reset()         { *m = SendWinMoney_C2S{} }
func (m *SendWinMoney_C2S) String() string { return proto.CompactTextString(m) }
func (*SendWinMoney_C2S) ProtoMessage()    {}
func (*SendWinMoney_C2S) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{14}
}

func (m *SendWinMoney_C2S) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendWinMoney_C2S.Unmarshal(m, b)
}
func (m *SendWinMoney_C2S) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendWinMoney_C2S.Marshal(b, m, deterministic)
}
func (m *SendWinMoney_C2S) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendWinMoney_C2S.Merge(m, src)
}
func (m *SendWinMoney_C2S) XXX_Size() int {
	return xxx_messageInfo_SendWinMoney_C2S.Size(m)
}
func (m *SendWinMoney_C2S) XXX_DiscardUnknown() {
	xxx_messageInfo_SendWinMoney_C2S.DiscardUnknown(m)
}

var xxx_messageInfo_SendWinMoney_C2S proto.InternalMessageInfo

func (m *SendWinMoney_C2S) GetMoney() float64 {
	if m != nil {
		return m.Money
	}
	return 0
}

type SendWinMoney_S2S struct {
	Account              float64  `protobuf:"fixed64,1,opt,name=account,proto3" json:"account,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SendWinMoney_S2S) Reset()         { *m = SendWinMoney_S2S{} }
func (m *SendWinMoney_S2S) String() string { return proto.CompactTextString(m) }
func (*SendWinMoney_S2S) ProtoMessage()    {}
func (*SendWinMoney_S2S) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{15}
}

func (m *SendWinMoney_S2S) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SendWinMoney_S2S.Unmarshal(m, b)
}
func (m *SendWinMoney_S2S) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SendWinMoney_S2S.Marshal(b, m, deterministic)
}
func (m *SendWinMoney_S2S) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SendWinMoney_S2S.Merge(m, src)
}
func (m *SendWinMoney_S2S) XXX_Size() int {
	return xxx_messageInfo_SendWinMoney_S2S.Size(m)
}
func (m *SendWinMoney_S2S) XXX_DiscardUnknown() {
	xxx_messageInfo_SendWinMoney_S2S.DiscardUnknown(m)
}

var xxx_messageInfo_SendWinMoney_S2S proto.InternalMessageInfo

func (m *SendWinMoney_S2S) GetAccount() float64 {
	if m != nil {
		return m.Account
	}
	return 0
}

// 发送获奖
type GetRewards_C2S struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRewards_C2S) Reset()         { *m = GetRewards_C2S{} }
func (m *GetRewards_C2S) String() string { return proto.CompactTextString(m) }
func (*GetRewards_C2S) ProtoMessage()    {}
func (*GetRewards_C2S) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{16}
}

func (m *GetRewards_C2S) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRewards_C2S.Unmarshal(m, b)
}
func (m *GetRewards_C2S) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRewards_C2S.Marshal(b, m, deterministic)
}
func (m *GetRewards_C2S) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRewards_C2S.Merge(m, src)
}
func (m *GetRewards_C2S) XXX_Size() int {
	return xxx_messageInfo_GetRewards_C2S.Size(m)
}
func (m *GetRewards_C2S) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRewards_C2S.DiscardUnknown(m)
}

var xxx_messageInfo_GetRewards_C2S proto.InternalMessageInfo

// 回传奖励
type GetRewards_S2C struct {
	RewardsNum           int32    `protobuf:"varint,1,opt,name=rewardsNum,proto3" json:"rewardsNum,omitempty"`
	RewardsMoney         int32    `protobuf:"varint,2,opt,name=rewardsMoney,proto3" json:"rewardsMoney,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRewards_S2C) Reset()         { *m = GetRewards_S2C{} }
func (m *GetRewards_S2C) String() string { return proto.CompactTextString(m) }
func (*GetRewards_S2C) ProtoMessage()    {}
func (*GetRewards_S2C) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{17}
}

func (m *GetRewards_S2C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRewards_S2C.Unmarshal(m, b)
}
func (m *GetRewards_S2C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRewards_S2C.Marshal(b, m, deterministic)
}
func (m *GetRewards_S2C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRewards_S2C.Merge(m, src)
}
func (m *GetRewards_S2C) XXX_Size() int {
	return xxx_messageInfo_GetRewards_S2C.Size(m)
}
func (m *GetRewards_S2C) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRewards_S2C.DiscardUnknown(m)
}

var xxx_messageInfo_GetRewards_S2C proto.InternalMessageInfo

func (m *GetRewards_S2C) GetRewardsNum() int32 {
	if m != nil {
		return m.RewardsNum
	}
	return 0
}

func (m *GetRewards_S2C) GetRewardsMoney() int32 {
	if m != nil {
		return m.RewardsMoney
	}
	return 0
}

// 修改房间区分配置
type ChangeRoomCfg_C2S struct {
	Config               string   `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangeRoomCfg_C2S) Reset()         { *m = ChangeRoomCfg_C2S{} }
func (m *ChangeRoomCfg_C2S) String() string { return proto.CompactTextString(m) }
func (*ChangeRoomCfg_C2S) ProtoMessage()    {}
func (*ChangeRoomCfg_C2S) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{18}
}

func (m *ChangeRoomCfg_C2S) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangeRoomCfg_C2S.Unmarshal(m, b)
}
func (m *ChangeRoomCfg_C2S) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangeRoomCfg_C2S.Marshal(b, m, deterministic)
}
func (m *ChangeRoomCfg_C2S) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeRoomCfg_C2S.Merge(m, src)
}
func (m *ChangeRoomCfg_C2S) XXX_Size() int {
	return xxx_messageInfo_ChangeRoomCfg_C2S.Size(m)
}
func (m *ChangeRoomCfg_C2S) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeRoomCfg_C2S.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeRoomCfg_C2S proto.InternalMessageInfo

func (m *ChangeRoomCfg_C2S) GetConfig() string {
	if m != nil {
		return m.Config
	}
	return ""
}

// 返回配置信息
type ChangeRoomCfg_S2C struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangeRoomCfg_S2C) Reset()         { *m = ChangeRoomCfg_S2C{} }
func (m *ChangeRoomCfg_S2C) String() string { return proto.CompactTextString(m) }
func (*ChangeRoomCfg_S2C) ProtoMessage()    {}
func (*ChangeRoomCfg_S2C) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{19}
}

func (m *ChangeRoomCfg_S2C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangeRoomCfg_S2C.Unmarshal(m, b)
}
func (m *ChangeRoomCfg_S2C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangeRoomCfg_S2C.Marshal(b, m, deterministic)
}
func (m *ChangeRoomCfg_S2C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeRoomCfg_S2C.Merge(m, src)
}
func (m *ChangeRoomCfg_S2C) XXX_Size() int {
	return xxx_messageInfo_ChangeRoomCfg_S2C.Size(m)
}
func (m *ChangeRoomCfg_S2C) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeRoomCfg_S2C.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeRoomCfg_S2C proto.InternalMessageInfo

type ErrorMsg_S2C struct {
	MsgData              string   `protobuf:"bytes,1,opt,name=msgData,proto3" json:"msgData,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ErrorMsg_S2C) Reset()         { *m = ErrorMsg_S2C{} }
func (m *ErrorMsg_S2C) String() string { return proto.CompactTextString(m) }
func (*ErrorMsg_S2C) ProtoMessage()    {}
func (*ErrorMsg_S2C) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{20}
}

func (m *ErrorMsg_S2C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ErrorMsg_S2C.Unmarshal(m, b)
}
func (m *ErrorMsg_S2C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ErrorMsg_S2C.Marshal(b, m, deterministic)
}
func (m *ErrorMsg_S2C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ErrorMsg_S2C.Merge(m, src)
}
func (m *ErrorMsg_S2C) XXX_Size() int {
	return xxx_messageInfo_ErrorMsg_S2C.Size(m)
}
func (m *ErrorMsg_S2C) XXX_DiscardUnknown() {
	xxx_messageInfo_ErrorMsg_S2C.DiscardUnknown(m)
}

var xxx_messageInfo_ErrorMsg_S2C proto.InternalMessageInfo

func (m *ErrorMsg_S2C) GetMsgData() string {
	if m != nil {
		return m.MsgData
	}
	return ""
}

func init() {
	proto.RegisterEnum("msg.MessageID", MessageID_name, MessageID_value)
	proto.RegisterType((*Ping)(nil), "msg.Ping")
	proto.RegisterType((*Pong)(nil), "msg.Pong")
	proto.RegisterType((*PlayerInfo)(nil), "msg.PlayerInfo")
	proto.RegisterType((*Login_C2S)(nil), "msg.Login_C2S")
	proto.RegisterType((*Login_S2C)(nil), "msg.Login_S2C")
	proto.RegisterType((*Logout_C2S)(nil), "msg.Logout_C2S")
	proto.RegisterType((*Logout_S2C)(nil), "msg.Logout_S2C")
	proto.RegisterType((*PlayerData)(nil), "msg.PlayerData")
	proto.RegisterType((*RoomData)(nil), "msg.RoomData")
	proto.RegisterType((*JoinRoom_C2S)(nil), "msg.JoinRoom_C2S")
	proto.RegisterType((*JoinRoom_S2C)(nil), "msg.JoinRoom_S2C")
	proto.RegisterType((*EnterRoom_S2C)(nil), "msg.EnterRoom_S2C")
	proto.RegisterType((*PlayerAction_C2S)(nil), "msg.PlayerAction_C2S")
	proto.RegisterType((*PlayerAction_S2C)(nil), "msg.PlayerAction_S2C")
	proto.RegisterType((*SendWinMoney_C2S)(nil), "msg.SendWinMoney_C2S")
	proto.RegisterType((*SendWinMoney_S2S)(nil), "msg.SendWinMoney_S2S")
	proto.RegisterType((*GetRewards_C2S)(nil), "msg.GetRewards_C2S")
	proto.RegisterType((*GetRewards_S2C)(nil), "msg.GetRewards_S2C")
	proto.RegisterType((*ChangeRoomCfg_C2S)(nil), "msg.ChangeRoomCfg_C2S")
	proto.RegisterType((*ChangeRoomCfg_S2C)(nil), "msg.ChangeRoomCfg_S2C")
	proto.RegisterType((*ErrorMsg_S2C)(nil), "msg.ErrorMsg_S2C")
}

func init() {
	proto.RegisterFile("msg.proto", fileDescriptor_c06e4cca6c2cc899)
}

var fileDescriptor_c06e4cca6c2cc899 = []byte{
	// 629 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x5d, 0x4f, 0xdb, 0x30,
	0x14, 0x5d, 0x9b, 0xb6, 0xb4, 0x97, 0xb6, 0xb8, 0x5e, 0x41, 0xd9, 0x1e, 0x50, 0xe5, 0x87, 0xa9,
	0xfb, 0x10, 0x93, 0xba, 0xa7, 0x4d, 0xdb, 0xc3, 0x16, 0x10, 0xea, 0x44, 0x51, 0x95, 0x20, 0xf1,
	0x88, 0xb2, 0xc6, 0x35, 0x11, 0xc4, 0x46, 0x4e, 0x18, 0xe2, 0xbf, 0xec, 0xc7, 0x4e, 0xd7, 0x89,
	0x43, 0x52, 0xfa, 0xb2, 0xbd, 0xf5, 0x9c, 0x7b, 0x7c, 0xee, 0x67, 0x03, 0xbd, 0x24, 0x15, 0x47,
	0x77, 0x5a, 0x65, 0x8a, 0x3a, 0x49, 0x2a, 0x58, 0x07, 0x5a, 0xcb, 0x58, 0x0a, 0xf6, 0x06, 0x5a,
	0x4b, 0x25, 0x05, 0x3d, 0x04, 0x48, 0xb9, 0xfe, 0xcd, 0xf5, 0x45, 0x9c, 0x70, 0xb7, 0x31, 0x69,
	0x4c, 0x1d, 0xbf, 0xc2, 0xb0, 0x5b, 0x80, 0xe5, 0x6d, 0xf8, 0xc8, 0xf5, 0x5c, 0xae, 0x15, 0x1d,
	0x42, 0x73, 0x1e, 0x19, 0x55, 0xcf, 0x6f, 0xce, 0x23, 0xfa, 0x1a, 0xba, 0x32, 0x5e, 0xdd, 0x9c,
	0x87, 0x09, 0x77, 0x9b, 0x86, 0x2d, 0x31, 0x75, 0x61, 0xe7, 0x9a, 0x87, 0xd1, 0x3c, 0x11, 0xae,
	0x63, 0x42, 0x16, 0x62, 0x24, 0x5c, 0xad, 0xd4, 0xbd, 0xcc, 0xdc, 0xd6, 0xa4, 0x31, 0x6d, 0xf8,
	0x16, 0xb2, 0x05, 0xf4, 0xce, 0x94, 0x88, 0xe5, 0x95, 0x37, 0x0b, 0xb6, 0x25, 0x5b, 0x86, 0x69,
	0x7a, 0xa9, 0x74, 0x64, 0x93, 0x59, 0x4c, 0xc7, 0xd0, 0xbe, 0x50, 0x37, 0x5c, 0x16, 0xa9, 0x72,
	0xc0, 0xbe, 0x5a, 0xbb, 0x60, 0xe6, 0xd1, 0x8f, 0x00, 0x77, 0x65, 0x27, 0xc6, 0x76, 0x77, 0xb6,
	0x77, 0x84, 0xe3, 0x79, 0x6a, 0xd0, 0xaf, 0x48, 0x58, 0x1f, 0xe0, 0x4c, 0x09, 0x75, 0x9f, 0x61,
	0x35, 0x15, 0x14, 0xcc, 0x3c, 0xf6, 0xcd, 0x8e, 0xe5, 0x38, 0xcc, 0xc2, 0x7f, 0xb7, 0x0e, 0xa0,
	0xeb, 0x2b, 0x95, 0x98, 0xc7, 0x07, 0xd0, 0xd1, 0x4a, 0x25, 0x65, 0xab, 0x05, 0x7a, 0x32, 0x45,
	0x95, 0xdb, 0x9c, 0x38, 0x1b, 0xa6, 0x48, 0xfb, 0x15, 0x09, 0x9b, 0x40, 0xff, 0xa7, 0x8a, 0x25,
	0x1a, 0x9b, 0xf9, 0x11, 0x70, 0x56, 0x6b, 0x51, 0xb8, 0xe2, 0x4f, 0xf6, 0xb9, 0xa2, 0xc0, 0x91,
	0xbc, 0x85, 0xae, 0x2e, 0xca, 0x28, 0xaa, 0x1e, 0x98, 0x04, 0xb6, 0x36, 0xbf, 0x0c, 0xb3, 0x2f,
	0x30, 0x38, 0x91, 0x19, 0xd7, 0xff, 0xf3, 0xf6, 0x03, 0x90, 0xbc, 0xe4, 0xef, 0xab, 0x2c, 0x56,
	0xf9, 0x72, 0x5d, 0xd8, 0x89, 0xd4, 0x83, 0xfc, 0xc1, 0x33, 0xf3, 0xba, 0xe1, 0x5b, 0xc8, 0xa6,
	0x1b, 0x6a, 0x4c, 0x36, 0x86, 0xf6, 0x3c, 0xbd, 0x8c, 0xa5, 0xd1, 0x76, 0xfd, 0x1c, 0xa0, 0x32,
	0xe0, 0x32, 0xba, 0x8c, 0xe5, 0x42, 0x49, 0xfe, 0x68, 0x7c, 0xc7, 0xd0, 0x4e, 0x10, 0x14, 0xae,
	0x39, 0xc0, 0x0a, 0x6a, 0xca, 0x20, 0xaf, 0xc0, 0x5e, 0x61, 0xa3, 0x7e, 0x85, 0x04, 0x86, 0xa7,
	0x3c, 0xf3, 0xf9, 0x43, 0xa8, 0xa3, 0xd4, 0x2c, 0xff, 0xa2, 0xc6, 0x60, 0x45, 0x87, 0x00, 0x3a,
	0x87, 0xe7, 0xf7, 0x89, 0x31, 0x68, 0xfb, 0x15, 0x86, 0x32, 0xe8, 0x17, 0xc8, 0x64, 0x34, 0x07,
	0xdb, 0xf6, 0x6b, 0x1c, 0x7b, 0x0f, 0x23, 0xef, 0x3a, 0x94, 0x82, 0xe3, 0xcc, 0xbc, 0xb5, 0x30,
	0x0d, 0x1c, 0x40, 0x67, 0xa5, 0xe4, 0x3a, 0xb6, 0x8b, 0x2b, 0x10, 0x7b, 0xb9, 0x29, 0xc6, 0x33,
	0x9c, 0x42, 0xff, 0x44, 0x6b, 0xa5, 0x17, 0xa9, 0xc1, 0xd8, 0x53, 0x92, 0x8a, 0x72, 0x27, 0x3d,
	0xdf, 0xc2, 0x77, 0x7f, 0x1c, 0xe8, 0x2d, 0x78, 0x9a, 0x86, 0x82, 0xcf, 0x8f, 0x69, 0x1f, 0xba,
	0x8b, 0xe0, 0xf4, 0x0a, 0xbf, 0x04, 0xe4, 0x45, 0x89, 0x94, 0x14, 0xa4, 0x41, 0x47, 0x30, 0x40,
	0x54, 0xfe, 0x0f, 0x49, 0xb3, 0x4e, 0x05, 0x33, 0x8f, 0x38, 0x94, 0xc2, 0xb0, 0xa0, 0x8a, 0x3f,
	0x08, 0x69, 0x6d, 0x70, 0xa8, 0x6b, 0xd3, 0x31, 0x10, 0xe4, 0xaa, 0x87, 0x49, 0x3a, 0xcf, 0x58,
	0xd4, 0xee, 0xd0, 0x7d, 0x18, 0x21, 0x5b, 0xbb, 0x33, 0xd2, 0xa5, 0x2e, 0x8c, 0x4d, 0x79, 0x1b,
	0x27, 0x44, 0x7a, 0x5b, 0x23, 0x81, 0x37, 0x23, 0x60, 0x23, 0x9b, 0xe7, 0x41, 0x76, 0xb7, 0x46,
	0x30, 0x4f, 0x9f, 0x1e, 0x00, 0xc5, 0x48, 0x7d, 0xf5, 0x64, 0xb0, 0x85, 0x47, 0xfd, 0x90, 0xbe,
	0x82, 0x7d, 0xe4, 0x9f, 0xad, 0x90, 0xec, 0x6d, 0x0f, 0xe1, 0x2b, 0x62, 0x5b, 0xaf, 0xae, 0x8d,
	0x8c, 0x7e, 0x75, 0xcc, 0x27, 0xfa, 0xd3, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xee, 0x1f, 0x58,
	0x90, 0xaf, 0x05, 0x00, 0x00,
}
