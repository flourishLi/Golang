package memory

import (
	//	"ZDTalk/errorcode"
	//	"ZDTalk/manager/db/mysqldb"
	logs "ZDTalk/utils/log4go"
	"sync"
)

//绘制内存管理结构体

type DrawCommandInfo struct {
	ServerId         int16
	ReceiverImUserId int32
	MessageId        int64

	MessageContent []byte
	SendTime       int64
	MessageFormat  int16 //消息类型
}

type DrawCommandMemoryManager struct {
	DrawContent map[int32]*DrawCommandInfo //多个房间 key:roomid value:DrawCommandInfo
	Lock        sync.Mutex
}

//全局内存变量
var drawCommandMemoryManager *DrawCommandMemoryManager

//初始化全局内存变量drawCommandMemoryManager
func GetDrawCommandMemoryManager() *DrawCommandMemoryManager {
	logs.Logs.Info("------------- ZDTalk Memory DrawCommandMemoryManager Initial end-------------")
	if drawCommandMemoryManager == nil {
		drawCommandMemoryManager = &DrawCommandMemoryManager{}
		drawCommandMemoryManager.DrawContent = make(map[int32]*DrawCommandInfo)
	}
	logs.Logs.Info("------------- ZDTalk Memory DrawCommandMemoryManager Initial end-------------")

	return drawCommandMemoryManager
}
