package customMsg

import (
	b "ZDTalk/utils/byteoperator"
	"bytes"
	"encoding/binary"
)

//发送给客户端自定义消息 转发结构体
type SendCustomClientMessageRequest struct {
	Header           *MessageHeader
	ServerId         int16
	ReceiverImUserId int32
	MessageId        int64

	MessageContent []byte
	SendTime       int64
	MessageFormat  int16 //消息类型
}

func (self *SendCustomClientMessageRequest) GetHeader() *MessageHeader {
	return self.Header
}

func (self *SendCustomClientMessageRequest) SetHeader(header *MessageHeader) {
	self.Header = header
}

func (self *SendCustomClientMessageRequest) IsGroupMessage() bool {
	return false
}

func CreateSendCustomCliendMessageRequest() *SendCustomClientMessageRequest {
	return &SendCustomClientMessageRequest{
		Header: CreateMessageHeader(SEND_CUSTOM_CLIENT_MESSAGE_REQUEST)}
}

func (self *SendCustomClientMessageRequest) CreateSendCustomClientMessageResponse() *SendCustomClientMessageResponse {
	respone := CreateSendCustomClientMessageResponse()
	respone.Header.AppId = self.Header.AppId
	respone.Header.UserId = self.Header.UserId
	respone.Header.SubAppId = self.Header.SubAppId
	respone.Header.DeviceType = self.Header.DeviceType
	return respone
}

func (self *SendCustomClientMessageRequest) ToBytes(key int) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(self.Header.ToBytes(key))
	buf.Write(b.Int16ToBytes(self.ServerId))

	buf.Write(b.Int32ToBytes(self.ReceiverImUserId))

	buf.Write(b.Int64ToBytes(self.MessageId))

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
func (self *SendCustomClientMessageRequest) FromBytes(data []byte, key int) (bool, int) {
	if bo, start := self.Header.FromBytes(data, key); bo {

		_, v16 := b.Bytes2Int16(data, start)
		self.ServerId = v16
		start += 2

		_, v32 := b.Bytes2Int32(data, start)
		self.ReceiverImUserId = v32
		start += 4

		_, v64 := b.Bytes2Int64(data, start)
		self.MessageId = v64
		start += 8

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
