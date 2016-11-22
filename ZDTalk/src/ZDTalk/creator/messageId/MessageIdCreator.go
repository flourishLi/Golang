package messageId

import (
	"ZDTalk/utils/timeutils"
	"sync"
)

type MessageIdCreator struct {
	lastId         int64
	mutext         sync.Mutex
	currentSecond  int64
	startHeaderNum int32
}

var inst *MessageIdCreator

func CreateMessageIdCreator(startHeaderNum int32) *MessageIdCreator {
	v := MessageIdCreator{}
	v.startHeaderNum = startHeaderNum
	inst = &v
	return &v
}

func getTimeSecond() int64 {
	v := timeutils.GetUnix13NowTime() / 1000

	return v
}

func (self *MessageIdCreator) CreateId() int64 {
	self.mutext.Lock()
	defer self.mutext.Unlock()

	v := getTimeSecond()

	if self.currentSecond == v {
		self.lastId++
		return self.lastId
	}

	self.currentSecond = v
	self.lastId = (int64(self.startHeaderNum)*10000000000+v)*10000 + 111

	return self.lastId

}

func CreateMessageId() int64 {
	return inst.CreateId()
}
