package queuehandler

import (
	msg "ZDTalk/queue/customMsg"
	logs "ZDTalk/utils/log4go"

	MessageIdFactory "ZDTalk/creator/messageId"
	"ZDTalk/utils/timeutils"

	"github.com/nsqio/go-nsq"
)

//接收IMServer 发送过来的自定义消息 句柄
type IMServerCustomMessageReceiveHandler struct {
	NsqConsumer *nsq.Consumer
}

//接收IMServer 发送过来的自定义消息
func (h *IMServerCustomMessageReceiveHandler) HandleMessage(msg_ *nsq.Message) error {

	go func(messages *nsq.Message) {
		defer func() {
			if err := recover(); err != nil {
				logs.GetLogger().Error("ZDTalk messagedispatcher 接收from IMServer nsq转发消息 error :", err)
			}
		}()

		logs.GetLogger().Info("********* ZDTalk messagedispatcher 接收from IMServer nsq转发消息 *********")
		arrowmessage := msg.CreateMessageFromBytes(msg_.Body, msg.SERVER_MESSAGE)

		switch v := arrowmessage.(type) {
		case *msg.SendCustomServerMessageRequest:

			//处理收到的IM消息
			messageId := MessageIdFactory.CreateId()
			//当前时间戳
			currentTime := timeutils.GetUnix13NowTime()

			logs.GetLogger().Info("********* from IMServer nsq转发消息 ServerId --> %d *********", v.CustomServerId)

			HandleServerMessage(messageId, currentTime, v.MessageFormat, v.CustomServerId, v.SenderImUserId, v.MessageContent)
		case *msg.SendCustomServerMessageResponse:
			logs.GetLogger().Info("服务端发给客户端的自定义消息的响应")
			logs.GetLogger().Info(v)
		default:
			logs.GetLogger().Error("收到错误消息:", v)
		}
	}(msg_)
	return nil
}

//处理收到的server消息
func HandleServerMessage(messageId, sendTime int64, messageFormat int16, serverId int16, senderImUserId int32, messageContent []byte) {
	if messageFormat == msg.ONLINE_HEARTBEAT { //更新在线状态
		//调用接口
		UpdateOnLine(senderImUserId, messageContent)
	} else {
		//创建发送消息结构体
		logs.GetLogger().Info("imUserId: ", senderImUserId)
		clientMessage := msg.SendCustomClientMessageRequest{}
		clientMessage.ServerId = serverId
		clientMessage.MessageFormat = messageFormat
		clientMessage.MessageContent = messageContent
		clientMessage.MessageId = messageId
		clientMessage.SendTime = sendTime
		if messageFormat == msg.ONLINE_TRANSPOND { //转发消息 向在线列表转发消息
			//定义接口对象
			iMMessage_Transpond := IMMessage_Transpond{}
			//转发消息
			iMMessage_Transpond.TransportMessage(senderImUserId, clientMessage)
		} else if messageFormat == msg.DRAW_TRANSPOND {
			//绘制命令转发
			DrawTranspond(senderImUserId, clientMessage)
		} else { //协议错误
			logs.GetLogger().Info("messageFormat is wrong", messageFormat)
			return
		}
	}
}
