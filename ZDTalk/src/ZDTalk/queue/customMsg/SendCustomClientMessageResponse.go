package customMsg

import (
	"bytes"

	b "ZDTalk/utils/byteoperator"
	"encoding/binary"
)

type SendCustomClientMessageResponse struct {
	Header   *MessageHeader
	Result   int32
	MsgId    int64
	SendTime int64
}

func (self *SendCustomClientMessageResponse) GetHeader() *MessageHeader {
	return self.Header
}

func (self *SendCustomClientMessageResponse) SetHeader(header *MessageHeader) {
	self.Header = header
}

func (self *SendCustomClientMessageResponse) IsGroupMessage() bool {
	return false
}

func CreateSendCustomClientMessageResponse() *SendCustomClientMessageResponse {
	return &SendCustomClientMessageResponse{
		Header: CreateMessageHeader(SEND_CUSTOM_CLIENT_MESSAGE_RESPONSE)}
}

func (self *SendCustomClientMessageResponse) ToBytes(key int) []byte {
	buf := bytes.NewBuffer([]byte{})
	buf.Write(self.Header.ToBytes(key))

	buf.Write(b.Int32ToBytes(self.Result))

	buf.Write(b.Int64ToBytes(self.MsgId))

	buf.Write(b.Int64ToBytes(self.SendTime))

	bs := buf.Bytes()
	bufs := bytes.NewBuffer([]byte{})
	messageLen := len(bs) + 2
	binary.Write(bufs, binary.BigEndian, int16(messageLen))
	bufs.Write(bs)
	return bufs.Bytes()
}

//参数key：1标识发给客户端用，2标识服务器内部使用
func (self *SendCustomClientMessageResponse) FromBytes(data []byte, key int) (bool, int) {
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

		return bo, start
	} else {
		return bo, start
	}
}
