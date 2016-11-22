// QueueOnlineManager
package memory

//在线用户管理说明
//数据中存储了在线用户Id集合([]int32)
//教室中在线用户数据有[输入输出]时，会修改内存数据，同时修改数据库；
//在线用户列表输入输出 1、进入(退出)教室接口[输入输出] 2、用户发送心跳(或心跳超时)[输入输出] 3、教师将用户移除教室[输出]
//服务启动时，ClassRoomMemoryManager会先去加载数据库中所有教室数据，数据库中在线用户有数据时，QueueOnlineMemoryManager从ClassRoomMemoryManager中获取数据到内存中，获取之后开启QueueOnlineMember的心跳超时检查

//心跳检查说明
//心跳检查机制由QueueOnlineMemoryManager与QueueOnlineMember[结构体]共同完成，QueueOnlineMember对外提供开启超时处理与更新心跳时间的方法
//QueueOnlineMemoryManager 接收到心跳时，先检查用户是否存在，不存在时，创建用户加入到内存修改数据库，并且开启QueueOnlineMember的心跳检查；存在时，更新QueueOnlineMember心跳时间；当QueueOnlineMember检查的心跳超时时，会调用QueueOnlineMemoryManager的TimeOut方法，然后由QueueOnlineMemoryManager来处理用户超时(将用户移除在线列表)

import (
	"ZDTalk/manager/db"
	"ZDTalk/manager/db/info"
	logs "ZDTalk/utils/log4go"
	"ZDTalk/utils/sliceutils"
	"fmt"
	//	"ZDTalk/queue/queuehandler"
	"ZDTalk/config"
	timeUtils "ZDTalk/utils/timeutils"
	"sync"
)

//维护心跳的Chanel 产生阻塞 避免超时处理死循环
type HeartBeatChannelManager struct {
	ChHeartBeat chan int32
}

//在线用户管理
type QueueOnlineMemoryManager struct {
	ClassRoomOnLineMap map[int32]map[int32]*QueueOnlineMember //教室的在线用户 key：roomId value:在线用户列表  每一个教室对应一个在线用户列表
	Lock               sync.Mutex
}

var queueOnlineMemoryManager *QueueOnlineMemoryManager

//初始化QueueOnlineMemoryManager 加载所有教室
func GetQueueOnlineMemoryManager() *QueueOnlineMemoryManager {

	//获取当前时间 初始化每一位新成员的举手时间
	currentTime := timeUtils.GetUnix13NowTime()
	if queueOnlineMemoryManager == nil {
		//定义queueOnlineMemoryManager对象
		queueOnlineMemoryManager = &QueueOnlineMemoryManager{}
		queueOnlineMemoryManager.ClassRoomOnLineMap = make(map[int32]map[int32]*QueueOnlineMember, 0)

		//实时更新 每次调用重新构造
		// 利用数据库的在线用户列表初始化数据
		//获取教室Memory对象 教室集合
		classRoomMemory := GetClassRoomMemoryManager()
		Classrooms := classRoomMemory.GetClassRoomList()

		//遍历数据库所有的教室
		for _, roomInfo := range Classrooms {
			logs.GetLogger().Info("初始化QueueOnlineMemoryManager时 教室 %d 的用户集合为 %d", roomInfo.ClassId, roomInfo.OnLineMemberList)

			//queueOnlineMemoryManager没有创建此教室
			if onlineMembers, ok := queueOnlineMemoryManager.ClassRoomOnLineMap[roomInfo.ClassId]; !ok {
				//构造该教室的queue线成员列表
				queueOnlineMembers := make(map[int32]*QueueOnlineMember)
				//遍历该教室下的所有在线成员 并添加到queue
				for _, userId := range roomInfo.OnLineMemberList {
					queueOnlineMember := &QueueOnlineMember{}
					queueOnlineMember.UserId = userId
					queueOnlineMember.ClassRoomId = roomInfo.ClassId
					queueOnlineMember.StartTime = currentTime
					queueOnlineMember.SetKeepAliving(true)
					queueOnlineMembers[userId] = queueOnlineMember
					//打开超时监听
					go queueOnlineMember.DetectHeartBeatTimeOut(queueOnlineMemoryManager, config.ConfigNodeInfo.OnLineTimeSpan)
				}
				//添加到queueManager
				queueOnlineMemoryManager.ClassRoomOnLineMap[roomInfo.ClassId] = queueOnlineMembers
			} else { //queueOnlineMemoryManager已有此教室
				//遍历该教室下的所有在线成员
				for _, userId := range roomInfo.OnLineMemberList {
					//不在queueOnlineMembers 添加
					if _, ok := onlineMembers[userId]; !ok {
						queueOnlineMember := QueueOnlineMember{}
						queueOnlineMember.UserId = userId
						queueOnlineMember.ClassRoomId = roomInfo.ClassId
						queueOnlineMember.StartTime = currentTime
						queueOnlineMember.SetKeepAliving(true)
						onlineMembers[userId] = &queueOnlineMember
						//打开超时监听
						go queueOnlineMember.DetectHeartBeatTimeOut(queueOnlineMemoryManager, config.ConfigNodeInfo.OnLineTimeSpan)
					}
				}
			}
		}
	}
	return queueOnlineMemoryManager
}

//检测在线用户心跳
func (self *QueueOnlineMemoryManager) DetectOnLineHeartBeat(paramOnlineMember *QueueOnlineMember, timeOutSpan int64) {

	if paramOnlineMember == nil {
		logs.GetLogger().Error("检测心跳时 参数 queueOnlineMember == nil")
		return
	}
	if !self.ContainOnlineClassRoom(paramOnlineMember.ClassRoomId) {
		logs.GetLogger().Error("检测心跳时 教室 %d 不存在", paramOnlineMember.ClassRoomId)
		return
	}

	self.Lock.Lock()
	defer self.Lock.Unlock()

	notExist := !self.ContainOnlineUser(paramOnlineMember.ClassRoomId, paramOnlineMember.UserId)

	//	logs.GetLogger().Info("检测心跳时 教室 %d 是否存在用户 %d  %s", paramOnlineMember.ClassRoomId, paramOnlineMember.UserId, !notExist)
	var userIdStr string
	var length int32 = 0
	if room, ok := self.ClassRoomOnLineMap[paramOnlineMember.ClassRoomId]; ok {
		for _, v := range room {
			userIdStr += "["
			userIdStr += fmt.Sprintf("UserId --> %d , ClassRoomId --> %d", v.UserId, v.ClassRoomId)
			userIdStr += "] "
			length++
		}
	}
	logs.GetLogger().Info("检测心跳时当前教室 %d 的在线数量 %d 用户Id集合为 %s ", paramOnlineMember.ClassRoomId, length, userIdStr)
	newHeartBeatTime := timeUtils.GetUnix13NowTime()
	//用户不在在线列表中时，发送了心跳时，将其加入到在线列表中，并且修改数据库
	if notExist {

		//在线用户中不存在该用户 添加用户到在线列表
		classRoom := self.ClassRoomOnLineMap[paramOnlineMember.ClassRoomId]
		logs.GetLogger().Info("检测心跳时 当前教室 %d 不存在 用户%d roomMap 为 %d", paramOnlineMember.ClassRoomId, userIdStr, classRoom)
		//根据参数 创建新用户添加到在线成员列表中
		onLineMember := &QueueOnlineMember{}
		onLineMember.UserId = paramOnlineMember.UserId
		onLineMember.ClassRoomId = paramOnlineMember.ClassRoomId
		onLineMember.SetKeepAliving(true)
		onLineMember.StartTime = newHeartBeatTime

		//给其他成员发送某人进入教室的通知
		onLineMemberSlice := []int32{}
		for _, value := range self.ClassRoomOnLineMap[onLineMember.ClassRoomId] {
			imUserId := GetUserInfoMemoryManager().GetUserIMId(value.UserId)
			onLineMemberSlice = append(onLineMemberSlice, imUserId)
		}
		SendMessage_EntryClassRoom(onLineMember.UserId, onLineMember.ClassRoomId, onLineMemberSlice)

		//		logs.GetLogger().Info("**** 用户不存在 添加到教室 用户信息为 userId %d ,classRoomId %d", onLineMember.UserId, onLineMember.ClassRoomId)

		self.addMemoryOnlineUser(onLineMember)
		//		logs.GetLogger().Info("检测心跳时 当前教室 %d 添加 用户%d 之后roomMap 为 %d", paramOnlineMember.ClassRoomId, length, userIdStr, classRoom)
		//某人第一次进入教室时，加入到内存中之后，开启检测，修改数据库
		go onLineMember.DetectHeartBeatTimeOut(self, timeOutSpan)
		go self.updateOnLineMemoryDb()

	} else {

		classRoom := self.ClassRoomOnLineMap[paramOnlineMember.ClassRoomId]
		onLineMember := classRoom[paramOnlineMember.UserId]
		//		logs.GetLogger().Info("**** 用户存在 用户信息为 userId %d ,classRoomId %d", onLineMember.UserId, onLineMember.ClassRoomId)
		if onLineMember == nil {
			logs.GetLogger().Error("更新心跳时 onLineMember == nil")
		}
		if !onLineMember.IsKeepAliving() {
			onLineMember.SetKeepAliving(true)
		}

		//更新心跳时间
		go onLineMember.UpdateHeartBeat(self, newHeartBeatTime, timeOutSpan)
	}

}

//开启心跳检测
func (self *QueueOnlineMemoryManager) StartDetectOnLineHeartBeat(paramOnlineMember *QueueOnlineMember) {
	if !self.ContainOnlineUser(paramOnlineMember.ClassRoomId, paramOnlineMember.UserId) {
		logs.GetLogger().Error("教室 %d 中开启用户超时检测时，用户 %d 不在教室内", paramOnlineMember.ClassRoomId, paramOnlineMember.UserId)
		return
	}
	classRoom := self.ClassRoomOnLineMap[paramOnlineMember.ClassRoomId]
	onLineMember := classRoom[paramOnlineMember.UserId]
	if !onLineMember.IsKeepAliving() {
		onLineMember.SetKeepAliving(true)
	}
	onLineMember.StartTime = timeUtils.GetUnix13NowTime()
	logs.GetLogger().Info("教室 %d 中的用户 %d 开启心跳超时检测", paramOnlineMember.ClassRoomId, paramOnlineMember.UserId)
	go onLineMember.DetectHeartBeatTimeOut(self, config.ConfigNodeInfo.OnLineTimeSpan)
}

//根据QueueOnLineMemoryManager中的数据 更新ClassRoomMemoryManager内存及数据库的在线用户列表
func (self *QueueOnlineMemoryManager) updateOnLineMemoryDb() {
	//	数据库操作对象
	dbManager := db.ClassRoomDbManager{}
	//	内存操作对象
	memoryManager := GetClassRoomMemoryManager()

	//	遍历queueManager每个教室
	for classRoomId, OnLineMemberSlice := range self.ClassRoomOnLineMap {
		if r, ok := memoryManager.ClassRooms[classRoomId]; ok {
			//重新构造内存的在线用户列表
			r.OnLineMemberList = make([]int32, 0)
			//遍历该教室下queue的在线成员 重新添加到内存
			for userId, _ := range OnLineMemberSlice {
				r.OnLineMemberList = append(r.OnLineMemberList, userId)
			}
			//更新数据库
			dbManager.EntryClassroom(r.OnLineMemberList, classRoomId)
		} else {
			logs.GetLogger().Info("ClassRoomId is not Exit", classRoomId)
		}
	}
}

// 更新数据库的发言用户列表
func (self *QueueOnlineMemoryManager) updateSayListDb() {
	//	数据库操作对象
	dbManager := db.ClassRoomDbManager{}
	//	内存操作对象
	memoryManager := GetClassRoomMemoryManager()
	//	遍历每个教室
	for classRoomId, roomInfo := range memoryManager.ClassRooms {
		//更新数据库
		dbManager.DeleteAddToSpeakArea(roomInfo.SayingMemberList, classRoomId)
	}
}

// 更新数据库的举手用户列表
func (self *QueueOnlineMemoryManager) updateHandMemberDb() {
	//	数据库操作对象
	dbManager := db.ClassRoomDbManager{}
	//	内存操作对象
	memoryManager := GetClassRoomMemoryManager()
	//	遍历每个教室
	for classRoomId, roomInfo := range memoryManager.ClassRooms {
		handList := []info.HandsMember{}
		for _, handMem := range roomInfo.HandMemberList {
			handList = append(handList, handMem)
		}
		//更新数据库
		dbManager.HandsUp(handList, classRoomId)

	}
}

func (self *QueueOnlineMemoryManager) TimeOut(queueOnlineMember *QueueOnlineMember) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	//从内存中移除在线成员
	self.removeMemoryOnlineUser(queueOnlineMember.ClassRoomId, queueOnlineMember.UserId)
	//从内存中移除举手成员
	self.removeMemoryHandList(queueOnlineMember.ClassRoomId, queueOnlineMember.UserId)
	//从内存中移除发言 成员
	self.removeMemorySayList(queueOnlineMember.ClassRoomId, queueOnlineMember.UserId)

	//从数据库中移除在线成员
	go self.updateOnLineMemoryDb()
	//从数据库中移除举手成员
	go self.updateHandMemberDb()
	//从数据库中移除发言成员
	go self.updateSayListDb()

	onLineMemberSlice := []int32{}
	for _, value := range self.ClassRoomOnLineMap[queueOnlineMember.ClassRoomId] {
		imUserId := GetUserInfoMemoryManager().GetUserIMId(value.UserId)
		onLineMemberSlice = append(onLineMemberSlice, imUserId)
	}
	length := len(onLineMemberSlice)
	if length == 0 { //没有用户时清除该教室的绘制命令
		drawCommandMemoryManager := GetDrawCommandMemoryManager()
		if _, ok := drawCommandMemoryManager.DrawContent[queueOnlineMember.ClassRoomId]; ok {
			delete(drawCommandMemoryManager.DrawContent, queueOnlineMember.ClassRoomId)
		}

	}
	//发送退出教室通知
	SendMessage_ExitClassRoom(queueOnlineMember.UserId, queueOnlineMember.ClassRoomId, onLineMemberSlice)
	//发送退出发言列表通知
	SendMessage_DeleteAddSpeaker(queueOnlineMember.UserId, queueOnlineMember.ClassRoomId, onLineMemberSlice)
	//发送退出举手列表通知
	SendMessage_HandsUp(queueOnlineMember.UserId, queueOnlineMember.ClassRoomId, onLineMemberSlice)
}

//判断是否包含用户
func (self *QueueOnlineMemoryManager) ContainOnlineUser(classRoomId int32, userId int32) bool {
	if classRoom, ok := self.ClassRoomOnLineMap[classRoomId]; !ok {
		return false
	} else {
		_, ok := classRoom[userId]
		return ok
	}
}

//判断是否包含教室
func (self *QueueOnlineMemoryManager) ContainOnlineClassRoom(classRoomId int32) bool {

	_, ok := self.ClassRoomOnLineMap[classRoomId]
	return ok

}

//内存中添加在线用户(加锁，供外部调用)
func (self *QueueOnlineMemoryManager) AddMemoryLockOnlineUser(queueOnlineMember *QueueOnlineMember) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	self.addMemoryOnlineUser(queueOnlineMember)
}

//内存中添加在线用户(不加锁，供内部调用)
func (self *QueueOnlineMemoryManager) addMemoryOnlineUser(queueOnlineMember *QueueOnlineMember) {
	if !self.ContainOnlineClassRoom(queueOnlineMember.ClassRoomId) {
		logs.GetLogger().Error("将要添加的在线用户的 教室 %d 不存在", queueOnlineMember.ClassRoomId)
		return
	}

	if self.ContainOnlineUser(queueOnlineMember.ClassRoomId, queueOnlineMember.UserId) {
		logs.GetLogger().Error("将要添加到教室 %d 的在线用户 %d 已经在教室中", queueOnlineMember.ClassRoomId, queueOnlineMember.UserId)
		return
	}

	if !queueOnlineMember.IsKeepAliving() {
		queueOnlineMember.SetKeepAliving(true)
	}
	if queueOnlineMember.StartTime == 0 {
		queueOnlineMember.StartTime = timeUtils.GetUnix13NowTime()
	}
	classRoom := self.ClassRoomOnLineMap[queueOnlineMember.ClassRoomId]
	classRoom[queueOnlineMember.UserId] = queueOnlineMember
}

//内存中移除在线用户(加锁，供外部调用)
func (self *QueueOnlineMemoryManager) RemoveMemoryLockOnlineUser(classRoomId int32, userId int32) {
	logs.GetLogger().Info("-------------------- RemoveMemoryOnlineUser -----------------------------------")
	self.Lock.Lock()
	defer self.Lock.Unlock()
	self.removeMemoryOnlineUser(classRoomId, userId)
}

//内存中移除在线用户(不加锁，供内部调用)
func (self *QueueOnlineMemoryManager) removeMemoryOnlineUser(classRoomId int32, userId int32) {
	logs.GetLogger().Info("Remove前 教室 %d 当前在线用户集合为 ", classRoomId, self.ClassRoomOnLineMap)

	if !self.ContainOnlineClassRoom(classRoomId) {
		logs.GetLogger().Info("将要被移除的 教室 %d 不存在", classRoomId)
		return
	}

	if !self.ContainOnlineUser(classRoomId, userId) {
		logs.GetLogger().Info("将要被移除的 教室 %d 中的用户 %d 不存在", classRoomId, userId)
		return
	}

	//从在线列表中移除用户
	classRoom := self.ClassRoomOnLineMap[classRoomId]
	onLineMember := classRoom[userId]
	onLineMember.SetKeepAliving(false)
	delete(classRoom, userId)
	logs.GetLogger().Info("Remove后后后 教室 %d 当前在线用户集合为 ", classRoomId, self.ClassRoomOnLineMap)

}

//内存中移除发言用户
//Param classRoomId userId
func (self *QueueOnlineMemoryManager) removeMemorySayList(classRoomId, userId int32) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	classRoomMemory := GetClassRoomMemoryManager()

	if r, ok := classRoomMemory.ClassRooms[classRoomId]; ok {
		//获取当前教室的发言列表
		sayUserIdList := r.SayingMemberList

		//用户是否在当前发言列表中
		exit, errOne := sliceutils.Containts(sayUserIdList, userId)
		if errOne != nil {
			logs.GetLogger().Error("sliceutils Containts err", errOne)
		}
		//在 移除
		if exit {
			sayUserIdList = sliceutils.RemoveInt32(sayUserIdList, userId)
			//更新内存发言列表
			r.SayingMemberList = sayUserIdList

		} else {
			logs.GetLogger().Info("将要被移除的 教室 %d 中的用户 %d 不存在", classRoomId, userId)
		}

	} else {
		logs.GetLogger().Error("教室不存在", classRoomId)
		return
	}
	return
}

//从内存中移除举手用户
//Param classRoomId userId
func (self *QueueOnlineMemoryManager) removeMemoryHandList(classRoomId, userId int32) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	classRoomMemory := GetClassRoomMemoryManager()

	if r, ok := classRoomMemory.ClassRooms[classRoomId]; ok {
		//用户在举手列表中 移除
		if _, ok := r.HandMemberList[userId]; ok {
			delete(r.HandMemberList, userId)
		} else {
			logs.GetLogger().Info("将要被移除的 教室 %d 中的用户 %d 不存在", classRoomId, userId)
		}
	} else {
		logs.GetLogger().Error("教室不存在", classRoomId)
		return
	}
	return
}

//获取在线用户信息
func (self *QueueOnlineMemoryManager) GetMemoryOnlineUser(classRoomId int32, userId int32) *QueueOnlineMember {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if ok := self.ContainOnlineUser(classRoomId, userId); !ok {
		return nil
	}
	classRoom := self.ClassRoomOnLineMap[classRoomId]
	onLineMember := classRoom[userId]
	return onLineMember
}
