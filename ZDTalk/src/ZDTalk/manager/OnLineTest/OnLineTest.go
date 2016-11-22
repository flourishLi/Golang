// OnLineTest
package main

import (
	"ZDTalk/manager/memory"
	"ZDTalk/manager/onLine"
	queueHandler "ZDTalk/queue/queuehandler"
	"time"
	//	byteUtils "ZDTalk/utils/byteoperator"
	logs "ZDTalk/utils/log4go"
	"ZDTalk/utils/timeutils"
	"encoding/json"
	//	"fmt"
)

type UserOnLine struct {
	SourceId       int32
	MessageContent []byte
	CurrentTime    int64
}

func main() {
	logs.GetLogger().Info("------------------- OnLineTest ------------------")
	onLine.EntryClassRoomAndSendHeartBeat()
}

func UpdateOnLineTest(sourceIMId, sourceNum, classroomNum int32) {
	//chan
	ch := make(chan int32)
	//接口获取对象
	classRoomManager := memory.GetClassRoomMemoryManager()
	queueManager := memory.GetQueueOnlineMemoryManager()

	heartBeat := queueHandler.HeartBeat{}

	var classRoomId int32
	var sourceId int32 //IM的userId
	for classRoomId = 1; classRoomId < classroomNum; classRoomId++ {

		if _, ok := classRoomManager.ClassRooms[classRoomId]; ok {
			/*输出更新之前的在线用户*/
			logs.GetLogger().Info("classRoomId is%d current Online Members", classRoomId)
			outOnlineUser(classRoomId, classRoomManager)

			//模拟登录
			sourceidNum := sourceIMId + sourceNum
			for sourceId = sourceIMId; sourceId < sourceidNum; sourceId++ {
				beat := queueHandler.HeartBeat{classRoomId}
				messages, err := json.Marshal(beat)
				if err != nil {
					logs.GetLogger().Info("json parse err", err)
				}
				go UserProduce(sourceId, messages, ch, queueManager, classRoomManager, heartBeat)
				id := <-ch
				logs.GetLogger().Info("登录的id", id)
				/*输出实时更新之后的在线用户*/
				logs.GetLogger().Info("classRoomId is%d CurrentUpdate Online Members", classRoomId)
				outOnlineUser(classRoomId, classRoomManager)
			}
		} else {
			logs.GetLogger().Info("classRoomId is not exist", classRoomId)
		}
		/*输出最终更新之后的在线用户*/
		logs.GetLogger().Info("classRoomId is%d LastUpdate Online Members", classRoomId)
		outOnlineUser(classRoomId, classRoomManager)
	}
}

func UserProduce(sourceId int32, messageContent []byte, ch chan int32, queueManager *memory.QueueOnlineMemoryManager, classRoomManager *memory.ClassRoomMemoryManager, heartBeat queueHandler.HeartBeat) {
	logs.GetLogger().Info("User update")
	currentTime := timeutils.GetUnix13NowTime()
	userOnLine := UserOnLine{sourceId, messageContent, currentTime}
	heartBeat.UpdateOnLine(userOnLine.SourceId, userOnLine.MessageContent, userOnLine.CurrentTime, queueManager, classRoomManager)
	time.Sleep(time.Second * 10) //增加执行的时间 看出超时效果
	ch <- sourceId               //一直堵塞 直到登录成功
}

func outOnlineUser(classRoomId int32, classRoomManager *memory.ClassRoomMemoryManager) {
	if room, ok := classRoomManager.ClassRooms[classRoomId]; ok {
		for _, v := range room.OnLineMemberList {
			logs.GetLogger().Info("userId:", v)
		}
	}
}
