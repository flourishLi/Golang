package info

import (
	b "ZDTalk/utils/byteoperator"
	"bytes"
)

//教室Db 结构体
type ClassRoomDbInfo struct {
	ClassRoomId         int32
	ClassRoomIMId       int32
	CreatorUserId       int32
	ClassRoomStatus     int32   //教室当前状态0 刚创建, 1 初始状态, 2 上课状态, 3 下课状态
	SettingStatus       []int32 //教室设置状态 1 允许所有人打字，2 允许学生录音，3 禁止游客打字，4 禁止游客举手
	CreateTime          int64
	ClassRoomName       string
	ClassRoomLogo       string
	Description         string
	ClassRoomCourse     string
	MemberList          []int32       //成员Id集合
	OnLineMemberList    []int32       //在线成员Id集合
	HandMemberList      []HandsMember //举手成员Id集合
	ForbidSayMemberList []int32       //被禁言的成员Id集合
	ForbidHandStatus    int32         //教室的禁止举手状态 教室的禁止举手状态 0可举手 1禁止举手
	SayingMemberList    []int32       //正在发言的成员Id集合
}

//举手成员列表需要依据时间排序
type HandsMember struct {
	UserId    int32 //举手成员的Id
	HandsTime int64 //举手时的时间
}

//[]int32转byte[]
func (self *ClassRoomDbInfo) ToBytes(data []int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	length := len(data)
	buf.Write(b.Int32ToBytes(int32(length))) //数据长度

	for _, v := range data { //数据
		buf.Write(b.Int32ToBytes(v))
	}

	return buf.Bytes()
}

func (self *ClassRoomDbInfo) StuctToBytes(data []HandsMember) []byte {
	buf := bytes.NewBuffer([]byte{})
	length := len(data)
	buf.Write(b.Int32ToBytes(int32(length))) //数据长度

	for _, v := range data { //数据
		buf.Write(b.Int32ToBytes(v.UserId))
		buf.Write(b.Int64ToBytes(v.HandsTime))

	}

	return buf.Bytes()
}

//[]byte转int32[]
func (self *ClassRoomDbInfo) FromBytes(data []byte) []int32 {
	start := 0
	//BLOB 实际占用 = 存储的字节字符串长度 + 值
	_, datalen := b.Bytes2Int32(data, start)

	start += 4
	result := make([]int32, 0)
	var i int32 = 0
	for i = 0; i < datalen; i++ {
		_, v32 := b.Bytes2Int32(data, start)
		start += 4
		result = append(result, v32)
	}
	return result
}

func (self *ClassRoomDbInfo) StructFromBytes(data []byte) []HandsMember {
	start := 0
	//BLOB 实际占用 = 存储的字节字符串长度 + 值
	_, datalen := b.Bytes2Int32(data, start)

	start += 4
	result := make([]HandsMember, 0)
	var i int32 = 0
	for i = 0; i < datalen; i++ {
		member := HandsMember{}
		_, v32 := b.Bytes2Int32(data, start)
		member.UserId = v32
		start += 4
		_, v64 := b.Bytes2Int64(data, start)
		member.HandsTime = v64
		start += 8
		result = append(result, member)
	}
	return result
}
