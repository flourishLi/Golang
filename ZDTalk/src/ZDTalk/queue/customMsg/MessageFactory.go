package customMsg

/**
消息创建工厂，通过 commandId 来创建各种request 与 response
*/
import (
	b "ZDTalk/utils/byteoperator"
	l4g "ZDTalk/utils/log4go"

	//	"strconv"
)

var logs = l4g.GetLogger()

//获取命令编号
func getCommandId(data []byte) int16 {
	if len(data) == 0 {
		return 0
	}
	_, v16 := b.Bytes2Int16(data, 2)
	return v16
}

//获取序列号
func getSequnceId(data []byte) int32 {
	_, v32 := b.Bytes2Int32(data, 4)
	return v32
}

//创建消息体
func CreateMessageFromBytes(data []byte, key int) ArrowMessage {
	commandId := getCommandId(data)
	switch commandId {
	//IMServer 发来的自定义消息请求
	case SEND_CUSTOM_CLIENT_MESSAGE_REQUEST:
		msg := CreateSendCustomCliendMessageRequest()
		msg.FromBytes(data, key)
		return msg

	//IMServer 发来的自定义消息响应
	case SEND_CUSTOM_CLIENT_MESSAGE_RESPONSE:
		msg := CreateSendCustomClientMessageResponse()
		msg.FromBytes(data, key)
		return msg

	//发送给IMServer 的自定义消息请求
	case SEND_CUSTOM_SERVER_MESSAGE_REQUEST:
		msg := CreateSendCustomServerMessageRequest()
		msg.FromBytes(data, key)
		return msg

	//发送给IMServer 的自定义消息响应
	case SEND_CUSTOM_SERVER_MESSAGE_RESPONSE:
		msg := CreateSendCustomServerMessageResponse()
		msg.FromBytes(data, key)
		return msg

	default:

		return nil
	}
}
