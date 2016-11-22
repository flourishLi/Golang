package queuehandler

import (
	"ZDTalk/manager/memory"
	msg "ZDTalk/queue/customMsg"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
)

func DrawTranspond(senderImUserId int32, clientMessage msg.SendCustomClientMessageRequest) {

	//将接收的byte内容转化成结构体
	drawContent := DrawContent{}
	messageContent := clientMessage.MessageContent
	logs.GetLogger().Info("绘制Json " + string(messageContent))
	err := json.Unmarshal(messageContent, &drawContent)
	if err != nil {
		logs.GetLogger().Error("messageContent unmarshal to struct is err", err)
		return
	}

	if drawContent.ClassroomId == 0 {
		logs.GetLogger().Error("教室Id 为 0 无法转发绘制命令")
		return
	}
	if !memory.GetClassRoomMemoryManager().IsExistLockClassRoom(drawContent.ClassroomId) {
		logs.GetLogger().Error("教室 %d 不存在， 无法转发绘制命令", drawContent.ClassroomId)
		return
	}
	//内存绘制 命令存储 供第一次进入教室的学员绘制命令
	drawCommandMemoryManager := memory.GetDrawCommandMemoryManager()
	//该教室已存在命令修改
	if drawCommandInfo, ok := drawCommandMemoryManager.DrawContent[drawContent.ClassroomId]; ok {
		drawCommandInfo.MessageContent = messageContent
		drawCommandInfo.ServerId = clientMessage.ServerId
		drawCommandInfo.SendTime = clientMessage.SendTime
		drawCommandInfo.MessageFormat = clientMessage.MessageFormat
		drawCommandInfo.MessageId = clientMessage.MessageId
	} else {
		//添加命令
		drawCommandInfo := &memory.DrawCommandInfo{}
		logs.GetLogger().Info("------------- 添加命令 ------------- classRoom ", drawCommandInfo)
		drawCommandInfo.MessageContent = messageContent
		drawCommandInfo.ServerId = clientMessage.ServerId
		drawCommandInfo.SendTime = clientMessage.SendTime
		drawCommandInfo.MessageContent = clientMessage.MessageContent
		drawCommandInfo.MessageFormat = clientMessage.MessageFormat
		drawCommandInfo.MessageId = clientMessage.MessageId
		drawCommandMemoryManager.DrawContent[drawContent.ClassroomId] = drawCommandInfo
	}
	//定义接口对象
	drawCommand := DrawCommand{}
	//转发消息
	drawCommand.TransportMessage(senderImUserId, clientMessage)
}
