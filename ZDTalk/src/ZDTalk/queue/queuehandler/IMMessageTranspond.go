package queuehandler

import (
	//	"ZDTalk/config"
	"ZDTalk/manager/memory"
	msg "ZDTalk/queue/customMsg"
	"ZDTalk/queue/transmitter"
	//	"ZDTalk/queue/publisher"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
)

//转发协议对应的MessageContent 协议编号 0x0702 需要转发
type IMMessage_Transpond struct { //需要用到转发内容中的某些字段 不定义结构体 无法引用内部字段
	ClassRoomId       int32  //教室编号
	MessageBody       string //消息内容(与IM协议中的MessageContent加以区分)
	MessageBodyFormat int16  //消息体类型(与IM中已有的一致) 1为文本信息，2为图片信息，3为音频信息，4为视频信息 ，5为通知类型信息， 6为上传文件，7为json类信息，8为地址类信息
	SendUserId        int    //发送者id
	SendUserName      string //发送者姓名
	SendRole          int    //发送者角色
	SendUserIcon      string //发送者头像
}

//转发消息
func (self IMMessage_Transpond) TransportMessage(senderImUserId int32, clientMsg msg.SendCustomClientMessageRequest) {
	transpondStruct := IMMessage_Transpond{}
	err := json.Unmarshal(clientMsg.MessageContent, &transpondStruct)
	if err != nil {
		logs.GetLogger().Info("json parse MessageContent is wrong", err)
		return
	}
	logs.GetLogger().Info("转发消息 当前房间 %d 接受者IMUserId %d", transpondStruct.ClassRoomId, clientMsg.ReceiverImUserId)
	if transpondStruct.ClassRoomId != 0 {
		//获取当前教室的在线列表结合
		onLineUsers := GetOnLineUsers(transpondStruct.ClassRoomId)
		var userIdStr string
		//向每一个成员转发消息

		for _, onLineUserId := range onLineUsers {
			receiverImUserId := memory.GetUserInfoMemoryManager().GetUserIMId(onLineUserId) //接收消息的用户ImUserId
			logs.GetLogger().Info("转发 消息 receiveIMUserId", receiverImUserId)
			logs.GetLogger().Info("转发 消息 content", string(clientMsg.MessageContent))
			if receiverImUserId != senderImUserId {

				transmitter.SendMessage2IMServer(clientMsg.ServerId, receiverImUserId, receiverImUserId, clientMsg.MessageFormat, clientMsg.MessageContent, clientMsg.MessageId, clientMsg.SendTime)

				//				transmitter.SendMessage2IMServer(cilentMsg.MessageId, cilentMsg.SendTime, cilentMsg.ServerId, cilentMsg.MessageFormat, receiverImUserId, receiverImUserId, cilentMsg.MessageContent)
				userIdStr += "["
				userIdStr += fmt.Sprintf("%d", onLineUserId)
				userIdStr += "] "
			}
		}
		logs.GetLogger().Info("当前教室 %d 的在线数量 %d 用户Id集合为 %s", transpondStruct.ClassRoomId, len(onLineUsers), userIdStr)
	} else {
		logs.GetLogger().Info("classRoomid is null:", transpondStruct.ClassRoomId)
		return
	}
}

//获取教室的在线列表结合
func GetOnLineUsers(classRoomId int32) []int32 {

	if classRoom, ok := memory.GetClassRoomMemoryManager().ClassRooms[classRoomId]; ok {
		onLineMembers := classRoom.OnLineMemberList
		return onLineMembers
	} else {
		logs.GetLogger().Error("ClassRoomId is not Exist", classRoomId)
		return nil
	}
}
