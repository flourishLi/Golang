package onLine

import (
	MessageIdFactory "ZDTalk/creator/messageId"
	"ZDTalk/manager/memory"
	msg "ZDTalk/queue/customMsg"
	queueHandler "ZDTalk/queue/queuehandler"
	//	httpUtils "ZDTalk/utils/httputil"
	"time"
	//	byteUtils "ZDTalk/utils/byteoperator"
	logs "ZDTalk/utils/log4go"
	"ZDTalk/utils/timeutils"
	"encoding/json"
	//	"fmt"
	//	"ZDTalk/ZDTalk_http/bean"
)

func EntryClassRoomAndSendHeartBeat() {
	logs.GetLogger().Info("------------------- EntryClassRoomAndSendHeartBeat------------------")
	//初始化消息Id创建工厂
	MessageIdFactory.InitMessageIdCreator(MessageIdFactory.ZDTALK_MESSAGE_ID_START_NUM)

	//	获取Memory对象
	classRoomMemory := memory.GetClassRoomMemoryManager()
	queueOnlineMemory := memory.GetQueueOnlineMemoryManager()

	userIdSlice := []int32{10033, 10034, 10035, 10036, 10037, 10038}
	requestSlice := make([]*msg.SendCustomServerMessageRequest, 0)
	var classRoomId int32 = 1
	logs.GetLogger().Info("------------------- OnLineTest for 循环前------------------")
	for index, value := range userIdSlice {

		//		entryClassRequest := &bean.EntryClassRoomRequest{}
		//		entryClassRequest.Command = "ENTRY_CLASSROOM"
		//		entryClassRequest.ClassRoomId = classRoomId
		//		entryClassRequest.RequestUserId = value
		//		jsonBuf, _ := json.Marshal(entryClassRequest)
		//		resultJsonBuf, _ := httpUtils.PostData("127.0.0.1:8015/ZTalkInterface", jsonBuf)

		//		entryClassResponse := &bean.EntryClassRoomResponse{}
		//		json.Unmarshal(resultJsonBuf, &entryClassResponse)
		//		if entryClassResponse.Code == 1 {
		//			logs.GetLogger().Info("用户 %d Http接口进入教室 %d 成功 ", value, classRoomId)
		//			time.Sleep(time.Second * 1)
		//		} else {
		//			logs.GetLogger().Info("用户 %d Http接口进入教室 %d 失败 错误码 %d 错误信息 %s ", value, classRoomId, entryClassResponse.Code, entryClassResponse.ErrMsg)
		//			continue
		//		}

		//发心跳
		messageId := MessageIdFactory.CreateId()
		//当前时间戳
		currentTime := timeutils.GetUnix13NowTime()
		msgRequest := &msg.SendCustomServerMessageRequest{}
		msgRequest.MessageId = messageId
		msgRequest.SourceId = value
		msgRequest.CustomServerId = 1
		msgRequest.SendTime = currentTime
		msgRequest.MessageFormat = msg.ONLINE_HEARTBEAT

		heartBeatStruct := queueHandler.HeartBeat{}
		if index == 4 || index == 5 || index == 6 {
			heartBeatStruct.ClassRoomId = 2
		}
		heartBeatStruct.ClassRoomId = classRoomId
		jsonBuf, _ := json.Marshal(heartBeatStruct)
		msgRequest.MessageContent = jsonBuf
		requestSlice = append(requestSlice, msgRequest)

	}

	for {

		for _, value := range requestSlice {

			queueHandler.HandleServerMessage(value.MessageId, value.SendTime, value.SendTime, value.MessageFormat, value.SourceId, value.MessageContent, classRoomMemory, queueOnlineMemory)

		}
		//发送一次心跳
		time.Sleep(time.Second * 10)
	}
}
