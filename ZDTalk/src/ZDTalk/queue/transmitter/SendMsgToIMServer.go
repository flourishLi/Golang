package transmitter

import (
	"ZDTalk/config"
	msg "ZDTalk/queue/customMsg"
	"ZDTalk/queue/publisher"
	logs "ZDTalk/utils/log4go"
)

//发送自定义消息给IMServer()
func SendMessage2IMServer(serverId int16, headUserId int32, receiverImUserId int32, messageFormat int16, content []byte, messageId int64, sendTime int64) {
	//func SendMessage2IMServer(messageId, sendTime int64, serverId int16, messageFormat int16, receiverImUserId int32, headUserId int32, content []byte) {
	clientMsg := msg.CreateSendCustomCliendMessageRequest()
	clientMsg.ReceiverImUserId = receiverImUserId
	clientMsg.ServerId = serverId
	clientMsg.MessageContent = content
	clientMsg.Header.UserId = headUserId
	clientMsg.MessageFormat = messageFormat
	clientMsg.MessageId = messageId
	clientMsg.SendTime = sendTime
	pool := publisher.GetQueuePool()
	queue, err := pool.GetQueue(config.ConfigNodeInfo.IMServerSendCustomMsgNsqUrl)
	if err != nil {
		logs.GetLogger().Error(err.Error())
	}
	queue.SendMessage(config.GetNsqSendTopic(), clientMsg.ToBytes(msg.SERVER_MESSAGE))
}
