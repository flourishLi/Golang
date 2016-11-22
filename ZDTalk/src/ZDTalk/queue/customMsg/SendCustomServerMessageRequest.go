package customMsg

import (
	"bytes"

	b "ZDTalk/utils/byteoperator"
	"encoding/binary"
)

const (
	MESSAGE_TYPE_NOTIFY = 0x0301
)

//接收的IM的心跳反馈信息
type SendCustomServerMessageRequest struct {
	Header         *MessageHeader
	SenderImUserId int32  //发送消息者ID
	MessageId      int64  //消息ID
	CustomServerId int16  //服务端Id
	MessageContent []byte //消息内容
	SendTime       int64  //发送消息时间
	MessageFormat  int16  //消息类型

}

func (self *SendCustomServerMessageRequest) GetHeader() *MessageHeader {
	return self.Header
}

func (self *SendCustomServerMessageRequest) SetHeader(header *MessageHeader) {
	self.Header = header
}

func (self *SendCustomServerMessageRequest) IsGroupMessage() bool {
	return false
}

func CreateSendCustomServerMessageRequest() *SendCustomServerMessageRequest {
	return &SendCustomServerMessageRequest{
		Header: CreateMessageHeader(SEND_CUSTOM_SERVER_MESSAGE_REQUEST)}
}

func (self *SendCustomServerMessageRequest) CreateResposne() *SendCustomServerMessageResponse {
	response := CreateSendCustomServerMessageResponse()

	response.Header.AppId = self.Header.AppId
	response.Header.UserId = self.Header.UserId
	response.Header.SubAppId = self.Header.SubAppId
	response.Header.DeviceType = self.Header.DeviceType
	response.ClientMsgId = self.MessageId

	return response
}

//参数key：1标识发给客户端用，2标识服务器内部使用
func (self *SendCustomServerMessageRequest) ToBytes(key int) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(self.Header.ToBytes(key))
	buf.Write(b.Int32ToBytes(self.SenderImUserId))

	buf.Write(b.Int64ToBytes(self.MessageId))
	buf.Write(b.Int16ToBytes(self.CustomServerId))
	buf.Write(b.MsgBytesToBytes(self.MessageContent))

	buf.Write(b.Int64ToBytes(self.SendTime))
	buf.Write(b.Int16ToBytes(self.MessageFormat))
	bs := buf.Bytes()
	bufs := bytes.NewBuffer([]byte{})
	messageLen := len(bs) + 2
	binary.Write(bufs, binary.BigEndian, int16(messageLen))
	bufs.Write(bs)
	return bufs.Bytes()
}

//参数key：1标识发给客户端用，2标识服务器内部使用
func (self *SendCustomServerMessageRequest) FromBytes(data []byte, key int) (bool, int) {
	if bo, start := self.Header.FromBytes(data, key); bo {
		_, v32 := b.Bytes2Int32(data, start)
		self.SenderImUserId = v32
		start += 4

		_, v64 := b.Bytes2Int64(data, start)
		self.MessageId = v64
		start += 8

		_, v16 := b.Bytes2Int16(data, start)
		self.CustomServerId = v16
		start += 2

		err, l, s := b.BytesToMsgBytesIncludeDataLength(data, start)

		if err != nil {
			return false, 0
		}
		self.MessageContent = s
		start += 2 + int(l)

		_, v64 = b.Bytes2Int64(data, start)
		self.SendTime = v64
		start += 8

		_, v16 = b.Bytes2Int16(data, start)
		self.MessageFormat = v16
		start += 2

		return bo, start
	} else {
		return bo, start
	}
}
