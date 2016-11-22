package memory

import (
	logs "ZDTalk/utils/log4go"
	timeUtils "ZDTalk/utils/timeutils"
	"sync"
	"time"
)

//在线用户结构体
type QueueOnlineMember struct {
	UserId             int32 //用户Id
	ClassRoomId        int32 //所在教室Id
	StartTime          int64 //最新的心跳时间戳
	isKeepAliving      bool  //是否正在保持心跳
	isDetectingTimeOut bool  //是否正在检测心跳是否超时
	lock               sync.Mutex
}

func (self *QueueOnlineMember) SetKeepAliving(keepAliving bool) {
	self.isKeepAliving = keepAliving
}

func (self *QueueOnlineMember) IsKeepAliving() bool {
	return self.isKeepAliving
}

func (self *QueueOnlineMember) SetDetectingTimeOut(isDetectingTimeOut bool) {
	self.isDetectingTimeOut = isDetectingTimeOut
}

func (self *QueueOnlineMember) IsDetectingTimeOut() bool {
	return self.isDetectingTimeOut
}

//开启在线用户的超时处理
func (self *QueueOnlineMember) DetectHeartBeatTimeOut(onLineMemoryManager *QueueOnlineMemoryManager, timeOutSpan int64) {
	for self.IsKeepAliving() {
		if !self.isDetectingTimeOut {
			self.SetDetectingTimeOut(true)
		}
		newHeartBeatTime := timeUtils.GetUnix13NowTime()
		//超时
		if newHeartBeatTime-self.StartTime > timeOutSpan {
			logs.GetLogger().Info("---------------------------------------------------------------")
			logs.GetLogger().Info("检测时 教室 %d 中的 用户 %d 心跳超时, 移除在线用户列表", self.ClassRoomId, self.UserId)
			logs.GetLogger().Info("---------------------------------------------------------------\n")
			self.SetKeepAliving(false)
			onLineMemoryManager.TimeOut(self)
		} else {
			//			logs.GetLogger().Info("**************************************************************************************")
			//			logs.GetLogger().Info("未检测到超时 教室 %d 用户 %d", self.ClassRoomId, self.UserId)
			//			logs.GetLogger().Info("检测时 当前时间为 %d", newHeartBeatTime)
			//			logs.GetLogger().Info("检测时 教室 %d 中的用户 %d 上一次心跳时间为 %d", self.ClassRoomId, self.UserId, self.StartTime)
			//			logs.GetLogger().Info("检测时 心跳超时设置的时间间隔为 %d", timeOutSpan)
			//			logs.GetLogger().Info("检测时 教室 %d 中的用户 %d 的与上一次心跳间隔为 %d", self.ClassRoomId, self.UserId, (newHeartBeatTime - self.StartTime))
			//			logs.GetLogger().Info("**************************************************************************************\n")
		}
		time.Sleep(time.Second * 6)
	}
	self.SetDetectingTimeOut(false)
}

func (self *QueueOnlineMember) UpdateHeartBeat(onLineMemoryManager *QueueOnlineMemoryManager, newHeartBeatTime int64, timeOutSpan int64) {

	self.lock.Lock()
	defer self.lock.Unlock()
	if self.IsKeepAliving() {
		if newHeartBeatTime-self.StartTime > timeOutSpan {
			//没在检测心跳超时(服务器重启时，从数据库中加载出来的在线用户列表并没有启动检测)
			if !self.IsDetectingTimeOut() {
				logs.GetLogger().Info("××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××")
				logs.GetLogger().Info("更新心跳时 教室 %d 中的 用户 %d 没有开启心跳超时检测，正在开启", self.ClassRoomId, self.UserId)
				logs.GetLogger().Info("××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××××")
				if !self.IsKeepAliving() {
					self.SetKeepAliving(true)
				}
				self.DetectHeartBeatTimeOut(onLineMemoryManager, timeOutSpan)
			}
			//		else {
			//				//检测心跳时，心跳发送时间超时
			//				logs.GetLogger().Info("---------------------------------------------------------------")
			//				logs.GetLogger().Info("更新心跳时 教室 %d 中的 用户 %d 心跳超时, 移除在线用户列表", self.ClassRoomId, self.UserId)
			//				logs.GetLogger().Info("---------------------------------------------------------------\n")
			//				self.SetKeepAliving(false)
			//				onLineMemoryManager.TimeOut(self)
			//			}

		} else {

			//未超时，更新心跳时间
			//			logs.GetLogger().Info("**************************************************************************************")
			//			logs.GetLogger().Info("更新心跳 教室 %d 中的 用户 %d 上一次心跳时间为 %d ", self.ClassRoomId, self.UserId, self.StartTime)
			//			logs.GetLogger().Info("更新心跳 教室 %d 中的 用户 %d 更新最新心跳时间为 %d ", self.ClassRoomId, self.UserId, newHeartBeatTime)
			//			logs.GetLogger().Info("更新心跳 教室 %d 中的 用户 %d 与上一次心跳间隔为 %d ", self.ClassRoomId, self.UserId, (newHeartBeatTime - self.StartTime))
			//			logs.GetLogger().Info("**************************************************************************************\n")
			self.StartTime = newHeartBeatTime
		}
	} else {
		logs.GetLogger().Error("更新心跳 用户 %d 更新最新心跳时间时 IsKeepAliving状态为false ", self.UserId)
	}

}
