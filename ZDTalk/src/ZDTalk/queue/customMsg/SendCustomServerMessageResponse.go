package customMsg

import (
	"bytes"

	b "ZDTalk/utils/byteoperator"
	"encoding/binary"
)

type SendCustomServerMessageResponse struct {
	Header      *MessageHeader
	Result      int32 //响应结果
	MsgId       int64 //后台生成的消息的ID
	SendTime    int64 //消息发送的时间
	ClientMsgId int64 //客户端所传信息ID
}

func (self *SendCustomServerMessageResponse) GetHeader() *MessageHeader {
	return self.Header
}

func (self *SendCustomServerMessageResponse) SetHeader(header *MessageHeader) {
	self.Header = header
}

func (self *SendCustomServerMessageResponse) IsGroupMessage() bool {
	return false
}

func CreateSendCustomServerMessageResponse() *SendCustomServerMessageResponse {
	return &SendCustomServerMessageResponse{
		Header: CreateMessageHeader(SEND_CUSTOM_SERVER_MESSAGE_RESPONSE)}
}

//参数key：1标识发给客户端用，2标识服务器内部使用
func (self SendCustomServerMessageResponse) ToBytes(key int) []byte {
	//	logs.Debug("send message response commandId--->", self.Header.CommandId)
	buf := bytes.NewBuffer([]byte{})
	buf.Write(self.Header.ToBytes(key))

	buf.Write(b.Int32ToBytes(self.Result))

	buf.Write(b.Int64ToBytes(self.MsgId))

	buf.Write(b.Int64ToBytes(self.SendTime))

	buf.Write(b.Int64ToBytes(self.ClientMsgId))

	bs := buf.Bytes()
	bufs := bytes.NewBuffer([]byte{})
	messageLen := len(bs) + 2
	binary.Write(bufs, binary.BigEndian, int16(messageLen))
	bufs.Write(bs)
	return bufs.Bytes()
}

//参数key：1标识发给客户端用，2标识服务器内部使用
func (self *SendCustomServerMessageResponse) FromBytes(data []byte, key int) (bool, int) {
	if bo, start := self.Header.FromBytes(data, key); bo {

		_, v32 := b.Bytes2Int32(data, start)
		self.Result = v32
		start += 4

		_, v64 := b.Bytes2Int64(data, start)
		self.MsgId = v64
		start += 8

		_, v64 = b.Bytes2Int64(data, start)
		self.SendTime = v64
		start += 8

		_, v64 = b.Bytes2Int64(data, start)
		self.ClientMsgId = v64
		start += 8

		return bo, start
	} else {
		return bo, start
	}
}
