package info

import (
	b "ZDTalk/utils/byteoperator"
	"bytes"
)

type ClassRoomMemberDbInfo struct {
	UserId        int32   //用户Id
	ChatId        int32   //用户Id对应的IM中的userId
	Role          int32   //用户身份 1 学生 3老师
	DeviceType    int32   //设备登录类型
	LoginName     string  //登录账号
	UserName      string  //用户昵称
	UserIcon      string  //用户头像
	Password      string  //密码
	YYToken       string  //YY提供的Token
	ClassRoomList []int32 //教室列表(Json转的Blob)
}

func (self *ClassRoomMemberDbInfo) ToBytes(data []int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	length := len(data)
	buf.Write(b.Int32ToBytes(int32(length)))

	for _, v := range data {
		buf.Write(b.Int32ToBytes(v))
	}

	return buf.Bytes()
}

//[]byte转int32[]
func (self *ClassRoomMemberDbInfo) FromBytes(data []byte) []int32 {
	start := 0
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
