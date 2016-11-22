package queuehandler

import (
	"ZDTalk/config"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
)

//在线状态协议对应的MessageContent 协议编号 0x0701 只接收 不需转发
type HeartBeat struct {
	ClassRoomId int32 //教室编号
}

//更新在线状态
func UpdateOnLine(imUserId int32, messageContent []byte) {
	heartBeatStruct := HeartBeat{}
	err := json.Unmarshal(messageContent, &heartBeatStruct)
	if err != nil {
		logs.GetLogger().Error("json parse MessageContent is wrong", err)
	}
	//	logs.GetLogger().Info("接收到心跳 heartBeat begin 教室 %d 发送者 %d", heartBeatStruct.ClassRoomId, imUserId)

	if heartBeatStruct.ClassRoomId == 0 {
		logs.GetLogger().Error("接收到心跳 classRoomId is 0")
		return
	}
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()
	if !classRoomMemoryManager.IsExistLockClassRoom(heartBeatStruct.ClassRoomId) {
		logs.GetLogger().Error("接收心跳时 教室 %d 不存在", heartBeatStruct.ClassRoomId)
		return
	}

	//在线状态处理
	//获取用户id
	userIdZTalk := memory.GetUserInfoMemoryManager().GetUserId(imUserId) //sourceId为IM中的userId

	if userIdZTalk == 0 {
		logs.GetLogger().Error("接收心跳时 未能根据用户 %d 的IMUserId 查到 早道用户Id", imUserId)
		return
	}
	//更新在线状态
	//	logs.GetLogger().Info("接收到心跳 更新用户在线状态 教室 %d 发送者UserId %d 发送者IMUserId %d", heartBeatStruct.ClassRoomId, userIdZTalk, imUserId)
	onLineMember := &memory.QueueOnlineMember{}
	onLineMember.UserId = userIdZTalk
	onLineMember.ClassRoomId = heartBeatStruct.ClassRoomId

	memory.GetQueueOnlineMemoryManager().DetectOnLineHeartBeat(onLineMember, config.ConfigNodeInfo.OnLineTimeSpan)

}
