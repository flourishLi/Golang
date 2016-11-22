package interfac

import (
	"ZDTalk/manager/memory"
)

//在线用户接口
type OnLineUserinterfac interface {
	//检测心跳
	//newHeartBeatTime 最新的心跳时间,timeOutSpan 超时时间
	DetectHeartBeat(onLineMemoryManager memory.QueueOnlineMemoryManager, newHeartBeatTime int64, timeOutSpan int64)

	//更新最新心跳时间
	UpdateHeartBeat(newHeartBeatTime int64)

	//超时时调用
	TimeOutHeartBeat(onLineMemoryManager memory.QueueOnlineMemoryManager, queueOnlineMember memory.QueueOnlineMember)
}
