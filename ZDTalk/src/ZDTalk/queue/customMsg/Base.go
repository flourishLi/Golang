package customMsg

import (
	b "ZDTalk/utils/byteoperator"
	"bytes"

	"fmt"
)

const (
	CLIENT_MESSAGE = 1
	SERVER_MESSAGE = 2

	PUSH_OFF_LINE     int16 = 1 //发送离线消息
	NOT_PUSH_OFF_LINE int16 = 0 //不发送离线消息
)

const (
	USER_OFF_ONLINE = -100 //用户不在线
)

type MessageHeader struct {
	MessageLength int16
	CommandId     int16

	////////////////////////服务器端使用
	UserId    int32 //请求者
	GatewayId int16 //请求者网关ID

	AppId      int16
	SubAppId   int16
	DeviceType int16

	TargetId int32 //消息发送给谁，如果为0时则执行具体消息中的TargetId，
	//主要用来同步用户在一个设备中发送消息时，同步到另一台设备上。

	RequestId      int64
	IS_pushOffLine int16 //是否推送离线消息 1为true , 0为false
}

func GetHeaderLength(key int) int {
	if key == SERVER_MESSAGE {
		return 12
	} else {
		return 4
	}
}

var ids *MessageIdManager

func getIdManager() *MessageIdManager {
	if ids == nil {
		ids = CreateMessageIdManager(0, 0)

	}

	return ids
}

func CreateMessageHeader(commandId int16) *MessageHeader {
	s := &MessageHeader{
		CommandId: commandId}

	s.RequestId = getIdManager().CreateId()
	s.IS_pushOffLine = PUSH_OFF_LINE
	return s
}

//参数key：1标识发给客户端用，2标识服务器内部使用
func (m *MessageHeader) ToBytes(key int) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(b.Int16ToBytes(m.CommandId))

	if key == SERVER_MESSAGE {
		buf.Write(b.Int32ToBytes(m.UserId))
		buf.Write(b.Int16ToBytes(m.GatewayId))
		buf.Write(b.Int16ToBytes(m.AppId))
		buf.Write(b.Int16ToBytes(m.SubAppId))
		buf.Write(b.Int16ToBytes(m.DeviceType))
		buf.Write(b.Int32ToBytes(m.TargetId))
		buf.Write(b.Int64ToBytes(m.RequestId))
		buf.Write(b.Int16ToBytes(m.IS_pushOffLine))
	}
	return buf.Bytes()
}

//参数key：1标识发给客户端用，2标识服务器内部使用
func (m *MessageHeader) FromBytes(data []byte, key int) (bool, int) {
	if len(data) < GetHeaderLength(key) {
		return false, 0
	}

	var start int
	start = 0

	_, v := b.Bytes2Int16(data, start)
	m.MessageLength = v
	start += 2
	_, m.CommandId = b.Bytes2Int16(data, start)
	start += 2

	if key == SERVER_MESSAGE {
		_, v2 := b.Bytes2Int32(data, start)
		m.UserId = v2
		start += 4

		_, v = b.Bytes2Int16(data, start)
		m.GatewayId = v
		start += 2

		_, v = b.Bytes2Int16(data, start)
		m.AppId = v
		start += 2

		_, v = b.Bytes2Int16(data, start)
		m.SubAppId = v
		start += 2

		_, v = b.Bytes2Int16(data, start)
		m.DeviceType = v
		start += 2

		_, v2 = b.Bytes2Int32(data, start)
		m.TargetId = v2
		start += 4

		_, v3 := b.Bytes2Int64(data, start)
		m.RequestId = v3
		start += 8

		_, v = b.Bytes2Int16(data, start)
		m.IS_pushOffLine = v
		start += 2

	}

	return true, start
}

func (self *MessageHeader) SetHeader(userId int32, gatewayId,
	appId, subAppId, deviceType int16) {
	self.AppId = appId
	self.UserId = userId
	self.GatewayId = gatewayId
	self.SubAppId = subAppId
	self.DeviceType = deviceType
}

func (self MessageHeader) Clone() *MessageHeader {

	d := MessageHeader{
		MessageLength:  self.MessageLength,
		CommandId:      self.CommandId,
		UserId:         self.UserId,
		GatewayId:      self.GatewayId,
		AppId:          self.AppId,
		SubAppId:       self.SubAppId,
		DeviceType:     self.DeviceType,
		IS_pushOffLine: self.IS_pushOffLine}
	return &d
}

type ArrowMessage interface {
	ToBytes(key int) []byte

	GetHeader() *MessageHeader

	SetHeader(header *MessageHeader)

	IsGroupMessage() bool
}

type ArrowRequest interface {
	GetRequestId() int64
	GetTargetId() int32
	GetSendTime() int64
	SetSendTime(currentTime int64)
	ToBytes(key int) []byte
}

//是否推送离线消息方法
func Is_PushOffLineMsg(param int16) bool {

	if param == PUSH_OFF_LINE {
		return true
	} else if param == NOT_PUSH_OFF_LINE {
		return false
	} else {
		fmt.Println("服务器内部的是否发送离线消息值错误为%d，默认返回发送离线消息", param)
		return true
	}
	return true
}
