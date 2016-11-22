package memory

import (
	"ZDTalk/errorcode"
	"ZDTalk/manager/db"
	dbInfo "ZDTalk/manager/db/info"
	logs "ZDTalk/utils/log4go"
	"ZDTalk/utils/sliceutils"
	"sync"
)

//教室信息表
type ClassRoomMemoryInfo struct {
	ClassId             int32                        //教室ID
	ClassRoomIMId       int32                        //教室对应IM中的ID
	CreatorUserId       int32                        //教室创建者ID
	ClassRoomStatus     int32                        //教室状态
	SettingStatus       []int32                      //成员设置
	CreateTime          int64                        //创建时间
	ClassName           string                       //教室名称
	Description         string                       //教室描述
	ClassLogo           string                       //教室头像
	ClassRoomCourse     string                       //课程信息
	MemberList          []int32                      //成员信息集合
	OnLineMemberList    []int32                      //在线成员信息集合
	HandMemberList      map[int32]dbInfo.HandsMember //举手成员信息集合
	ForbidSayMemberList []int32                      //被禁言成员信息集合
	ForbidHandStatus    int32                        //教室的禁止举手状态

	SayingMemberList []int32 //正在发言成员信息集合
}

//教室信息内存管理结构体
type ClassRoomMemoryManager struct {
	ClassRooms     map[int32]*ClassRoomMemoryInfo
	ClassRoomNames map[string]*ClassRoomMemoryInfo
	Lock           sync.Mutex
}

//全局内存变量
// 对应数据库表classRoomInfo
var classRoomMemoryManager *ClassRoomMemoryManager

//初始化全局内存变量ClassRoomMemoryManager 即读取将数据表classRoomInfo到内存
func GetClassRoomMemoryManager() *ClassRoomMemoryManager {
	if classRoomMemoryManager == nil {
		logs.Logs.Info("------------- ZDTalk classRoomInfo Memory Initial Begin-------------")

		//初始化roomManager
		classRoomMemoryManager = &ClassRoomMemoryManager{}
		classRoomMemoryManager.ClassRoomNames = make(map[string]*ClassRoomMemoryInfo)
		classRoomMemoryManager.ClassRooms = make(map[int32]*ClassRoomMemoryInfo)

		//加载数据库中classRoomInfo数据表到内存 客户端的请求操作直接去内存读取数据 即声明的全局变量memory.classRoomMemoryManager
		manager := new(db.ClassRoomDbManager)
		result, err := manager.LoadAllClassRoom()
		if err != nil {
			logs.GetLogger().Error("LoadAllClassRoom Error:" + err.Error())
			return nil
		}

		for _, roomInfo := range result {
			//初始化举手对象(ID 举手时间)
			handMemberMap := make(map[int32]dbInfo.HandsMember)
			for _, handMember := range roomInfo.HandMemberList {
				handMemberMap[handMember.UserId] = handMember
			}
			//初始化memory.ClassRoomMemoryInfo
			room := &ClassRoomMemoryInfo{ClassId: roomInfo.ClassRoomId, ClassRoomIMId: roomInfo.ClassRoomIMId, CreatorUserId: roomInfo.CreatorUserId, ClassRoomStatus: roomInfo.ClassRoomStatus, SettingStatus: roomInfo.SettingStatus, ClassName: roomInfo.ClassRoomName, Description: roomInfo.Description, ClassLogo: roomInfo.ClassRoomLogo, ClassRoomCourse: roomInfo.ClassRoomCourse, MemberList: roomInfo.MemberList, OnLineMemberList: roomInfo.OnLineMemberList, HandMemberList: handMemberMap, ForbidSayMemberList: roomInfo.ForbidSayMemberList, SayingMemberList: roomInfo.SayingMemberList, ForbidHandStatus: roomInfo.ForbidHandStatus}
			//设置roomManager的ClassRooms属性
			classRoomMemoryManager.ClassRooms[roomInfo.ClassRoomId] = room
			//设置roomManager的ClassRoomNames属性
			classRoomMemoryManager.ClassRoomNames[roomInfo.ClassRoomName] = room

		}
		logs.Logs.Info("------------- ZDTalk classRoomInfo Memory Initial end-------------")
	}
	return classRoomMemoryManager
}

//删除在线用户
func (self *ClassRoomMemoryManager) RemoveOnLineUser(classRoomId int32, userId int32) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if !self.isExistClassRoom(classRoomId) {
		return
	}

	classRoom := self.ClassRooms[classRoomId]
	onLineMemberList := classRoom.OnLineMemberList
	//删除用户
	sliceutils.RemoveInt32(onLineMemberList, userId)
}

//获取教室信息
func (self *ClassRoomMemoryManager) GetClassRoom(roomId int32) *ClassRoomMemoryInfo {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	if v, ok := self.ClassRooms[roomId]; ok {
		return v
	} else {
		logs.GetLogger().Error("教室不存在")
		return nil
	}

	return nil
}

//获取教室列表
func (self *ClassRoomMemoryManager) GetClassRoomList() []*ClassRoomMemoryInfo {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	data := []*ClassRoomMemoryInfo{}

	for _, v := range self.ClassRooms {
		data = append(data, v)
	}

	return data
}

//创建教室
//Param classRoomName教室名称 classRoomLogo教室头像 description教室说明 ClassRoomCourse课程信息 CreatorUserId创建者ID CreateTime创建时间 ClassRoomIMId 教室对应IM的Id
//return code ErrMsg roomID
func (self *ClassRoomMemoryManager) CreateClassroom(roomID int32, icon string, name string, description string, course string, userID int32, creatTime int64, roomIMID int32, settingStatus []int32) (int32, string, int32) {
	logs.GetLogger().Info("create classroom In Memory begins")
	self.Lock.Lock()
	defer self.Lock.Unlock()

	//每一个举手对象
	//	handMember := dbInfo.HandsMember{}
	//	handMember.UserId = userID
	//	handMember.HandsTime = creatTime

	r := &ClassRoomMemoryInfo{}
	r.ClassId = roomID
	r.ClassRoomIMId = roomIMID
	r.CreatorUserId = userID
	r.ClassRoomStatus = 0 //刚创建
	r.SettingStatus = settingStatus
	r.CreateTime = creatTime
	r.ClassName = name
	r.ClassLogo = icon
	r.Description = description
	r.ClassRoomCourse = course
	//	r.ForbidSayMemberList = []int32{userID}
	//	r.HandMemberList = make(map[int32]dbInfo.HandsMember)
	//	r.HandMemberList[handMember.UserId] = handMember
	//	r.MemberList = []int32{userID}
	//	r.OnLineMemberList = []int32{userID}
	//	r.SayingMemberList = []int32{userID}

	self.ClassRooms[r.ClassId] = r
	self.ClassRoomNames[name] = r
	logs.GetLogger().Info("CreateClassRoom In Memory Success:", r)
	return errorcode.SUCCESS, "", r.ClassId
}

//删除教室
func (self *ClassRoomMemoryManager) DeleteClassroom(RoomId int32) (int32, string) {
	logs.GetLogger().Info("Delete classroom In Memory begins")

	self.Lock.Lock()
	defer self.Lock.Unlock()

	if r, ok := self.ClassRooms[RoomId]; ok {
		delete(self.ClassRoomNames, r.ClassName)
	}

	if _, ok := self.ClassRooms[RoomId]; !ok {
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}

	delete(self.ClassRooms, RoomId)
	logs.GetLogger().Info("DeleteClassRoom In Memory Success:")
	return errorcode.SUCCESS, "DeleteClassroom SUCCESS"
}

//修改教室信息
func (self *ClassRoomMemoryManager) UpdateClassroom(RoomId int32, name, logo, description, course string) (int32, string) {
	logs.GetLogger().Info("Update classroom In Memory begins")

	self.Lock.Lock()
	defer self.Lock.Unlock()
	if r, ok := self.ClassRooms[RoomId]; ok {
		if logo != "" {
			r.ClassLogo = logo
		}
		if description != "" {
			r.Description = description
		}
		if course != "" {
			r.ClassRoomCourse = course
		}
		if name != "" {
			r.ClassName = name
		}
	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("UpdateClassRoom In Memory Success:")
	return errorcode.SUCCESS, "UpdateClassroom SUCCESS"
}

//教室设置 SettingStatus设置
func (self *ClassRoomMemoryManager) UpdateClassRoomSettingStatus(RoomId, teacherUserID int32, settingStatus []int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if r, ok := self.ClassRooms[RoomId]; ok {
		r.SettingStatus = settingStatus
	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("ClassRoomSettingStatus In Memory Success:")
	return errorcode.SUCCESS, "ClassRoomSettingStatus SUCCESS"
}

//老师上下课设置 classRoomStatus设置
func (self *ClassRoomMemoryManager) UpdateClassRoomStatus(RoomId, teacherUserID, classRoomStatus int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("ClassRoomStatus IS:", classRoomStatus)
	if r, ok := self.ClassRooms[RoomId]; ok {
		r.ClassRoomStatus = classRoomStatus

	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("ClassRoomStatus In Memory Success:")
	return errorcode.SUCCESS, "ClassRoomStatus SUCCESS"
}

//教室禁止举手状态设置 ForbidHandStatus 0可举手 1禁止举手
func (self *ClassRoomMemoryManager) UpdateHandForbidStatus(RoomId, forbidHandStatus int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("UpdateHandForbidStatus IS:", forbidHandStatus)
	if r, ok := self.ClassRooms[RoomId]; ok {
		r.ForbidHandStatus = forbidHandStatus

	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("UpdateHandForbidStatus In Memory Success:")
	return errorcode.SUCCESS, "UpdateHandForbidStatus SUCCESS"
}

//查看内存中是否存在 教室
func (self *ClassRoomMemoryManager) IsExistLockClassRoom(classRoomId int32) bool {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	ok := self.isExistClassRoom(classRoomId)
	return ok

}

func (self *ClassRoomMemoryManager) isExistClassRoom(classRoomId int32) bool {
	_, ok := self.ClassRooms[classRoomId]
	return ok
}

//Param classRoomId userId handsType 1=举手，2=取消举手
func (self *ClassRoomMemoryManager) HandsUp(classRoomId, userId, handsType int32, handsTime int64) (int32, string) {

	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("HandsUp In Memory begins")

	if r, ok := self.ClassRooms[classRoomId]; ok {
		//获取当前教室的举手列表
		handMemberMap := r.HandMemberList
		logs.GetLogger().Info("OriginMemory HandMemberList Is:", handMemberMap)

		//创建举手对象
		handMember := dbInfo.HandsMember{}
		handMember.UserId = userId
		handMember.HandsTime = handsTime

		//添加到举手列表
		if handsType == HandsType_Up {
			//原列表不存在该用户
			if _, ok := handMemberMap[handMember.UserId]; !ok {
				handMemberMap[handMember.UserId] = handMember
			} else {
				logs.GetLogger().Info("UserId Is Exist:", handMember.UserId)
				return errorcode.USER_ID_IS_EXIST, errorcode.USER_ID_IS_EXIST_MSG
			}
		} else if handsType == HandsType_Down { //从举手列表删除
			//原列表存在该用户
			if _, ok := handMemberMap[handMember.UserId]; ok {
				delete(handMemberMap, handMember.UserId)
			} else {
				logs.GetLogger().Info("UserId Is not Exit:", handMember.UserId)
				return errorcode.USER_ID_IS_NOT_EXIT, errorcode.USER_ID_IS_NOT_EXIT_MSG
			}
		}

		logs.GetLogger().Info("CurrentMemory HandMemberList Is:", handMemberMap)
		//根据新的举手列表更新内存
		r.HandMemberList = handMemberMap
	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("HandsUp In Memory Success")
	return errorcode.SUCCESS, "HandsUp SUCCESS"
}

//禁止 解禁动作请求 禁言区
//Param classRoomId userIdList
func (self *ClassRoomMemoryManager) ForbidArea(classRoomId, crudType int32, userIdList []int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("ForbidArea In Memory begins")

	if r, ok := self.ClassRooms[classRoomId]; ok {
		//获取当前教室的禁止举手列表
		currentUserIdList := r.ForbidSayMemberList
		logs.GetLogger().Info("OriginMemory ForbidSayMemberList Is:", currentUserIdList)

		if CrudType_ADD == crudType { //禁止举手
			//遍历要添加的禁止列表成员
			for _, userID := range userIdList {
				//用户是否在当前列表中
				isExit, err := sliceutils.Containts(currentUserIdList, userID)
				if err != nil {
					logs.GetLogger().Error("sliceutils Containts err", err)
				}
				//不在 添加
				if !isExit {
					currentUserIdList = append(currentUserIdList, userID)
				} else {
					logs.GetLogger().Info("UserID: is exit", userID)
				}
			}
		} else if CrudType_DELETE == crudType { //解禁举手
			//遍历要解禁列表成员
			for _, userID := range userIdList {
				//用户是否在当前列表中
				isExit, err := sliceutils.Containts(currentUserIdList, userID)
				if err != nil {
					logs.GetLogger().Error("sliceutils Containts err", err)
				}
				//在 删除
				if isExit {
					currentUserIdList = sliceutils.RemoveInt32(currentUserIdList, userID)
				} else {
					logs.GetLogger().Info("UserID: not exit", userID)
				}
			}
		}
		logs.GetLogger().Info("CurrentMemory ForbidSayMemberList Is:", currentUserIdList)
		//根据新的举手列表更新内存
		r.ForbidSayMemberList = currentUserIdList
	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("ForbidArea In Memory Success")
	return errorcode.SUCCESS, "ForbidArea SUCCESS"
}

//添加 移除举手到发言区动作请求
//Param classRoomId userIdList
func (self *ClassRoomMemoryManager) DeleteAddToSpeakArea(classRoomId, curdType int32, userIdList []int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("AddToSpeakArea In Memory begins")
	if r, ok := self.ClassRooms[classRoomId]; ok {
		//获取当前教室的发言列表
		currentUserIdList := r.SayingMemberList
		logs.GetLogger().Info("OriginMemory SayingMemberList Is:", currentUserIdList)

		//添加到发言区
		if curdType == CrudType_ADD {
			//遍历要添加的举手发言列表成员
			for _, userID := range userIdList {
				//用户是否在当前列表中
				isExit, err := sliceutils.Containts(currentUserIdList, userID)
				if err != nil {
					logs.GetLogger().Error("sliceutils Containts err", err)
				}
				//不在 添加
				if !isExit {
					currentUserIdList = append(currentUserIdList, userID)
				}
				//添加到发言区的user 需要从举手列表中删除
				if _, ok := r.HandMemberList[userID]; ok {
					delete(r.HandMemberList, userID)
				}

			}
		} else if curdType == CrudType_DELETE { //从发言区移除
			//遍历要移除的举手发言列表成员
			for _, userID := range userIdList {
				//用户是否在当前列表中
				isExit, err := sliceutils.Containts(currentUserIdList, userID)
				if err != nil {
					logs.GetLogger().Error("sliceutils Containts err", err)
				}
				//在 移除
				if isExit {
					currentUserIdList = sliceutils.RemoveInt32(currentUserIdList, userID)
				}
			}
		}

		logs.GetLogger().Info("CurrentMemory SayingMemberList Is:", currentUserIdList)
		//根据新的举手列表更新内存
		r.SayingMemberList = currentUserIdList
	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("AddToSpeakArea In Memory Success")
	return errorcode.SUCCESS, "AddToSpeakArea SUCCESS"
}

//学生进入教室 维护在线列表
//Param classRoomId userId
func (self *ClassRoomMemoryManager) EntryClassRoom(classRoomId, userId int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("EntryClassRoom In Memory begins")
	if r, ok := self.ClassRooms[classRoomId]; ok {
		//获取当前教室的在线列表
		currentUserIdList := r.OnLineMemberList
		logs.GetLogger().Info("OriginMemory OnLineMemberList Is:", currentUserIdList)

		//用户是否在当前列表中
		isExit, err := sliceutils.Containts(currentUserIdList, userId)
		if err != nil {
			logs.GetLogger().Error("sliceutils Containts err", err)
		}
		//logs.GetLogger().Info("user is Exit", isExit)
		//不在 添加
		if !isExit {
			currentUserIdList = append(currentUserIdList, userId)
			//更新内存
			r.OnLineMemberList = currentUserIdList
		} else {
			logs.GetLogger().Error("用户已存在")
			return errorcode.USER_ID_IS_EXIST, errorcode.USER_ID_IS_EXIST_MSG
		}
		logs.GetLogger().Info("CurrentMemory OnLineMemberList Is:", currentUserIdList)
	} else {
		logs.GetLogger().Error("教室 %d 不存在", classRoomId)
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("EntryClassRoom In Memory Success")
	return errorcode.SUCCESS, "EntryClassRoom SUCCESS"
}

//退出教室 维护在线列表 发言列表 举手列表
//Param classRoomId userId
func (self *ClassRoomMemoryManager) ExitClassRoom(classRoomId, userId int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("ExitClassRoom In Memory begins")
	if r, ok := self.ClassRooms[classRoomId]; ok {

		//获取当前教室的在线列表
		onlineUserIdList := r.OnLineMemberList
		//获取当前教室的发言列表
		sayUserIdList := r.SayingMemberList

		//用户是否在当前在线列表中
		isExit, err := sliceutils.Containts(onlineUserIdList, userId)
		if err != nil {
			logs.GetLogger().Error("sliceutils Containts err", err)
		}
		//在 移除
		if isExit {
			onlineUserIdList = sliceutils.RemoveInt32(onlineUserIdList, userId)
			//更新内存在线列表
			r.OnLineMemberList = onlineUserIdList
		} else {
			logs.GetLogger().Info("user is exit")
			return errorcode.USER_ID_IS_EXIT, errorcode.USER_ID_IS_EXIT_MSG
		}

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

		}

		//用户在举手列表中 移除
		if _, ok := r.HandMemberList[userId]; ok {
			delete(r.HandMemberList, userId)
		}

	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("ExitClassRoom In Memory Success")
	return errorcode.SUCCESS, "ExitClassRoom SUCCESS"
}

//将学生移除 添加教室动作请求
//Param classRoomId userIdList
func (self *ClassRoomMemoryManager) DeleteAddStudents(classRoomId int32, userIdList []int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("DeleteStudents In Memory begins")
	if r, ok := self.ClassRooms[classRoomId]; ok {
		//获取当前教室的成员列表
		currentUserIdList := r.MemberList
		logs.GetLogger().Info("OriginMemory MemberList Is:", currentUserIdList)
		//更新内存
		r.MemberList = userIdList
		logs.GetLogger().Info("CurrentMemory MemberList Is:", r.MemberList)

	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("DeleteStudents In Memory Success")
	return errorcode.SUCCESS, "DeleteStudents SUCCESS"
}

//清空举手列表动作请求
//Param classRoomId
//return errorCode errMSG
func (self *ClassRoomMemoryManager) HandsListlear(classRoomId int32) (int32, string) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("HandsListlear In Memory begins")

	if r, ok := self.ClassRooms[classRoomId]; ok {
		//清空举手列表
		for _, handMember := range r.HandMemberList {
			delete(r.HandMemberList, handMember.UserId)
		}
	} else {
		logs.GetLogger().Error("教室不存在")
		return errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	logs.GetLogger().Info("HandsListlear In Memory Success")
	return errorcode.SUCCESS, "HandsListlear SUCCESS"
}
