package customMsg

import (
	"ZDTalk/utils/timeutils"
	"sync"
)

type MessageIdManager struct {
	lastId         int64
	mutext         sync.Mutex
	currentSecond  int64
	startHeaderNum int32
}

const (
	GATEWAY_SERVER = 1
	IM_SERVER      = 2
	GROUP_SERVER   = 3
)

var inst *MessageIdManager

func CreateMessageIdManager(serverType int32, serverNum int32) *MessageIdManager {
	v := MessageIdManager{}
	v.startHeaderNum = serverType*1000 + serverNum
	inst = &v
	return &v
}

func getTimeSecond() int64 {
	v := timeutils.GetUnix13NowTime() / 1000

	return v
}

func (self *MessageIdManager) CreateId() int64 {
	self.mutext.Lock()
	defer self.mutext.Unlock()

	v := getTimeSecond()

	if self.currentSecond == v {
		self.lastId++
		return self.lastId
	}

	self.currentSecond = v
	self.lastId = (int64(self.startHeaderNum)*10000000000+v)*100000 + 1

	return self.lastId

}

func CreateMessageId() int64 {
	return inst.CreateId()
}
