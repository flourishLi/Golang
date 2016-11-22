package process

import (
	"ZDTalk/errorcode"
	"ZDTalk/manager/db"
	"ZDTalk/manager/db/info"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
)

type ClassRoomProcess struct {
}

const (
	ZERO = 0
)

//创建教室
//Param classRoomName教室名称 classRoomLogo教室头像 description教室说明 ClassRoomCourse课程信息 CreatorUserId创建者ID CreateTime创建时间 ClassRoomIMId 教室对应IM的Id
//return errcode errMsg id 教室ID
func (roomProcess *ClassRoomProcess) CreateClassRoom(classRoomName, classRoomLogo, description, classRoomCourse string, creatorUserId int32, settingStatus []int32, createTime int64, classRoomIMId int32) (int32, string, int32) {
	logs.GetLogger().Info("create classroom In DataBase begins")
	//ClassRoomDbManager 接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}

	roomID, err := classRoomDBManager.CreateClassRoom(classRoomName, classRoomLogo, description, classRoomCourse, creatorUserId, settingStatus, createTime, classRoomIMId)
	if err != nil {
		logs.GetLogger().Error("CreateClassRoom dbm err", err)
		return errorcode.Create_CLASSROOM_DATEBASE_ERR, err.Error(), ZERO
	}
	id := int32(roomID)
	logs.GetLogger().Info("Create ClassRoom In DataBase Success")
	return errorcode.SUCCESS, "", id
}

//删除教室
//Param roomId教室ID
//return errcode errMsg
func (roomProcess *ClassRoomProcess) DeleteClassRoom(roomID int32) (int32, string) {
	logs.GetLogger().Info("Delete classroom In Database begins")

	//ClassRoomDbManager 接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}

	affectedRow, err := classRoomDBManager.DeleteClassRoom(roomID)
	if err != nil {
		logs.GetLogger().Error("DeleteClassRoom dbm err", err)
		return errorcode.DELETE_CLASSROOM_DATEBASE_ERR, err.Error()
	}
	deleteFlag := int32(affectedRow)
	//数据库操作没有任何影响
	if deleteFlag == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, errorcode.DATABASE_ROW_AFFECT_NULL_MSG
	}
	logs.GetLogger().Info("DeleteClassRoom  In DataBase Success")
	return errorcode.SUCCESS, ""
}

//更新教室
//Param roomId教室ID 参数利用内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) UpdateClassRoom(roomID int32) (int32, string) {
	logs.GetLogger().Info("Update classroom In Database begins")

	//ClassRoomDbManager 接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory  接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//内存操作已经完成 当前的classRoomInfo即为处理客户端请求后的classRoomInfo
	//通过Memory接口对象 获取该教室已存在的信息 非必须参数为空时更新原值
	classRoomInfo := classRoomMemoryManager.GetClassRoom(roomID)

	logs.GetLogger().Info("The Update classRoomInfo Is:", classRoomInfo)

	affectedRow, err := classRoomDBManager.UpdateClassRoom(roomID, classRoomInfo.ClassName, classRoomInfo.ClassLogo, classRoomInfo.Description, classRoomInfo.ClassRoomCourse)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom dbm err", err)
		return errorcode.UPDATE_CLASSROOM_DATEBASE_ERR, err.Error()
	}
	updateFlag := int32(affectedRow)
	if updateFlag == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("UpdateClassRoom  In DataBase Success")
	return errorcode.SUCCESS, ""
}

//教室设置 SettingStatus更新
//Param roomId教室ID teacherUserID 参数利用内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) UpdateClassRoomSettingStatus(roomID, teacherUserID int32) (int32, string) {
	//ClassRoomDbManager 接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory  接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//内存操作已经完成 当前的classRoomInfo即为处理客户端请求后的classRoomInfo
	//通过Memory接口对象 获取该教室已存在的信息 非必须参数为空时更新原值
	classRoomInfo := classRoomMemoryManager.GetClassRoom(roomID)

	logs.GetLogger().Info("The Update classRoomInfo Is:", classRoomInfo)

	affectedRow, err := classRoomDBManager.UpdateClassRoomSettingStatus(roomID, teacherUserID, classRoomInfo.SettingStatus)
	if err != nil {
		logs.GetLogger().Error("ClassRoomSettingStatus dbm err", err)
		return errorcode.UPDATE_CLASSROOM_DATEBASE_ERR, err.Error()
	}
	updateFlag := int32(affectedRow)
	if updateFlag == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("ClassRoomSettingStatus  In DataBase Success")
	return errorcode.SUCCESS, ""
}

//老师上下课 ClassRoomStatus更新
//Param roomId教室ID teacherUserID 参数利用内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) UpdateClassRoomStatus(roomID, teacherUserID int32) (int32, string) {
	//ClassRoomDbManager 接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory  接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//内存操作已经完成 当前的classRoomInfo即为处理客户端请求后的classRoomInfo
	//通过Memory接口对象 获取该教室已存在的信息 非必须参数为空时更新原值
	classRoomInfo := classRoomMemoryManager.GetClassRoom(roomID)

	logs.GetLogger().Info("The Update classRoomInfo Is:", classRoomInfo)

	affectedRow, err := classRoomDBManager.UpdateClassRoomStatus(roomID, teacherUserID, classRoomInfo.ClassRoomStatus)
	if err != nil {
		logs.GetLogger().Error("ClassRoomStatus dbm err", err)
		return errorcode.UPDATE_CLASSROOM_DATEBASE_ERR, err.Error()
	}
	updateFlag := int32(affectedRow)
	if updateFlag == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("ClassRoomStatus  In DataBase Success")
	return errorcode.SUCCESS, ""
}

//教室的禁止举手状态更新 ForbidHandStatus
//Param roomId教室ID forbidHandStatus 参数利用内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) UpdateForbidHandStatus(roomID int32) (int32, string) {
	//ClassRoomDbManager 接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory  接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//内存操作已经完成 当前的classRoomInfo即为处理客户端请求后的classRoomInfo
	//通过Memory接口对象 获取该教室已存在的信息 非必须参数为空时更新原值
	classRoomInfo := classRoomMemoryManager.GetClassRoom(roomID)

	logs.GetLogger().Info("The Update classRoomInfo Is:", classRoomInfo)

	affectedRow, err := classRoomDBManager.UpdateHandForbidStatus(roomID, classRoomInfo.ForbidHandStatus)
	if err != nil {
		logs.GetLogger().Error("UpdateForbidHandStatus dbm err", err)
		return errorcode.UPDATE_FORBIDHAND_STATUS_DATEBASE_ERR, err.Error()
	}
	updateFlag := int32(affectedRow)
	if updateFlag == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, errorcode.DATABASE_ROW_AFFECT_NULL_MSG
	}
	logs.GetLogger().Info("UpdateForbidHandStatus  In DataBase Success")
	return errorcode.SUCCESS, ""
}

//举手(取消)动作请求
//Param classRoomId 参数通过内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) HandsUp(classRoomID int32) (int32, string) {
	logs.GetLogger().Info("HandsUp In Database begins")

	//ClassRoomDbManager接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory接口对象
	classRoomManager := memory.GetClassRoomMemoryManager()

	//根据教室编号从内存获取该教室的举手列表
	classRoomDbInfo := classRoomManager.GetClassRoom(classRoomID)
	//map对象
	currentUserIdList := classRoomDbInfo.HandMemberList

	logs.GetLogger().Info("The HandsUp UserIDList Is:", currentUserIdList)

	//slice对象 数据库存储的是列表
	userIdSlice := []info.HandsMember{}
	for _, handMember := range currentUserIdList {
		userIdSlice = append(userIdSlice, handMember)
	}

	logs.GetLogger().Info("userIdSlice IS:", userIdSlice)
	//内存操作已经完成 当前的HandMemberList即为处理客户端请求后的HandMemberList
	//根据计算的userIDlist 更新classroominfo的HandMemberList
	code, err := classRoomDBManager.HandsUp(userIdSlice, classRoomID)
	if err != nil {
		logs.GetLogger().Error("HandsUp dbm err", err)
		return errorcode.HANDS_UP_DATEBASE_ERROR, err.Error()
	}
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("HandsUp  In DataBase Success")
	return errorcode.SUCCESS, ""
}

//禁止 解禁请求 禁言区
//Param classRoomId 参数通过内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) ForbidArea(classRoomID int32) (int32, string) {
	logs.GetLogger().Info("ForbidArea In DataBase begins")

	//ClassRoomDbManager接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//根据教室编号从内存获取该教室当前禁止举手列表
	classRoomDbInfo := classRoomMemoryManager.GetClassRoom(classRoomID)
	currentUserIdList := classRoomDbInfo.ForbidSayMemberList

	logs.GetLogger().Info("The ForbidArea UserIDList Is:", currentUserIdList)

	//内存操作已经完成 当前的ForbidSayMemberList即为删除studentids后的ForbidSayMemberList
	//根据计算的userIDlist 更新classroominfo的HandMemberList
	code, err := classRoomDBManager.ForbidArea(currentUserIdList, classRoomID)
	if err != nil {
		logs.GetLogger().Error("ForbidArea dbm err", err)
		return errorcode.HANDS_FORBID_DATEBASE_ERROR, err.Error()
	}
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("ForbidArea  In DataBase Success")
	return errorcode.SUCCESS, ""
}

//将学生添加 移除教室动作请求
//Param classRoomId 参数通过内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) DeleteAddStudents(classRoomID int32) (int32, string) {

	logs.GetLogger().Info("DeleteStudents In Database begins")

	//ClassRoomDbManager接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//根据教室编号从内存获取该教室当前教室成员列表
	classRoomDbInfo := classRoomMemoryManager.GetClassRoom(classRoomID)
	currentUserIdList := classRoomDbInfo.MemberList

	//根据计算的userIDlist 更新classroominfo的MemberList
	code, err := classRoomDBManager.DeleteAddStudents(currentUserIdList, classRoomID)
	//数据库操作发送错误
	if err != nil {
		logs.GetLogger().Error("DeleteStudents dbm err", err)
		return errorcode.DELETESTUDENTS_DATEBASE_ERROR, err.Error()
	}
	//数据库操作没有影响
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, errorcode.DATABASE_ROW_AFFECT_NULL_MSG
	}

	logs.GetLogger().Info("DeleteStudents  In DataBase Success")
	return errorcode.SUCCESS, ""
}

//清空举手列表动作请求
//Param classRoomId
//return errcode errMsg
func (roomProcess *ClassRoomProcess) HandsListClear(classRoomID int32) (int32, string) {
	logs.GetLogger().Info("HandsListClear In DataBase begins")

	//ClassRoomDbManager接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}

	affectedRow, err := classRoomDBManager.HandsListClear(classRoomID)
	if err != nil {
		logs.GetLogger().Error("HandsListClear dbm err", err)
		return errorcode.HANDS_CLEAR_DATEBASE_ERROR, err.Error()
	}
	updateFlag := int32(affectedRow)
	if updateFlag == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("HandsListClear  In DataBase Success")

	return errorcode.SUCCESS, ""
}

//添加 移除举手到 发言区 动作请求
//Param classRoomId 参数通过内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) DeleteAddToSpeakArea(classRoomID int32) (int32, string) {
	logs.GetLogger().Info("AddToSpeakArea In DataBase begins")

	//ClassRoomDbManager接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory接口对象
	classRoomManager := memory.GetClassRoomMemoryManager()

	//根据教室编号从内存获取该教室的发言列表
	classRoomDbInfo := classRoomManager.GetClassRoom(classRoomID)
	currentUserIdList := classRoomDbInfo.SayingMemberList

	logs.GetLogger().Info("The AddToSpeakArea UserIDList Is:", currentUserIdList)
	//内存操作已经完成 当前的SayingMemberList即为添加studentids后的SayingMemberList
	//根据计算的userIDlist 更新classroominfo的HandMemberList
	code, err := classRoomDBManager.DeleteAddToSpeakArea(currentUserIdList, classRoomID)
	if err != nil {
		logs.GetLogger().Error("AddToSpeakArea dbm err", err)
		return errorcode.ADDSPEAKTOAREA_DATEBASE_ERROR, err.Error()
	}
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("AddToSpeakArea  In DataBase Success")

	return errorcode.SUCCESS, ""
}

//学生进入教室动作请求
//Param classRoomId onlineuserIdList参数通过内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) EntryClassRoom(classRoomID int32) (int32, string) {
	logs.GetLogger().Info("EntryClassRoom In DataBase begins")

	//ClassRoomDbManager接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory接口对象
	classRoomManager := memory.GetClassRoomMemoryManager()

	//根据教室编号从内存获取该教室的在线列表
	classRoomDbInfo := classRoomManager.GetClassRoom(classRoomID)
	currentUserIdList := classRoomDbInfo.OnLineMemberList

	logs.GetLogger().Info("The EntryClassRoom UserIDList Is:", currentUserIdList)
	//内存操作已经完成
	//根据计算的userIDlist 更新classroominfo的OnLineMemberList
	code, err := classRoomDBManager.EntryClassroom(currentUserIdList, classRoomID)
	if err != nil {
		logs.GetLogger().Error("EntryClassRoom dbm err", err)
		return errorcode.TOFORBID_DATEBASE_ERROR, err.Error()
	}
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("EntryClassRoom  In DataBase Success")

	return errorcode.SUCCESS, ""
}

//退出教室动作请求
//Param classRoomId onlineuserIdList参数通过内存获取
//return errcode errMsg
func (roomProcess *ClassRoomProcess) ExitClassRoom(classRoomID int32) (int32, string) {
	logs.GetLogger().Info("ExitClassRoom In DataBase begins")

	//ClassRoomDbManager接口对象
	classRoomDBManager := &db.ClassRoomDbManager{}
	//Memory接口对象
	classRoomManager := memory.GetClassRoomMemoryManager()

	//根据教室编号从内存获取该教室的在线列表 举手列表 发言列表
	classRoomDbInfo := classRoomManager.GetClassRoom(classRoomID)
	onlineUserIdList := classRoomDbInfo.OnLineMemberList
	sayUserIdList := classRoomDbInfo.SayingMemberList
	handUserMap := classRoomDbInfo.HandMemberList
	//slice对象 数据库存储的是列表
	handUserlist := make([]int32, 0)
	for _, handMember := range handUserMap {
		handUserlist = append(handUserlist, handMember.UserId)
	}

	logs.GetLogger().Info("The ExitClassRoom UserIDList Is:", onlineUserIdList)
	logs.GetLogger().Info("The ExitClassRoom UserIDList Is:", handUserlist)
	logs.GetLogger().Info("The ExitClassRoom UserIDList Is:", sayUserIdList)

	//内存操作已经完成
	//根据计算的userIDlist 更新classroominfo
	code, err := classRoomDBManager.ExitClassroom(onlineUserIdList, handUserlist, sayUserIdList, classRoomID)
	if err != nil {
		logs.GetLogger().Error("ExitClassRoom dbm err", err)
		return errorcode.TOFORBID_DATEBASE_ERROR, err.Error()
	}
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, errorcode.DATABASE_ROW_AFFECT_NULL_MSG
	}
	logs.GetLogger().Info("ExitClassRoom  In DataBase Success")

	return errorcode.SUCCESS, ""
}

//添加上传文件资源
//Param userID fileName filePath fileThumbPath fileTime
//return code errMsg fileId
func (roomProcess *ClassRoomProcess) UploadResource(userId int32, fileType int32, fileName, filePath string, fileTime int64) (int32, string, int32) {
	logs.GetLogger().Info("UploadResource In DataBase begins")

	//ResourceDbManager接口对象
	resourceDBManager := &db.ResourceDbManager{}

	code, err := resourceDBManager.InsertResource(userId, fileType, fileName, filePath, fileTime)
	if err != nil {
		logs.GetLogger().Error("UploadResource dbm err", err)
		return errorcode.UPLOAD_RESOURCE_DATEBASE_ERROR, err.Error(), ZERO
	}
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响", ZERO
	}
	logs.GetLogger().Info("UploadResource  In DataBase Success")
	fileId := int32(code)
	return errorcode.SUCCESS, "", fileId
}

//添加上传文件资源附属信息 filecount
//Param fileid fileContentCount
//return code errMsg fileId
func (roomProcess *ClassRoomProcess) InsertResourceAddition(fileId int32, fileContentCount int32) (int32, string, int32) {
	logs.GetLogger().Info("InsertResourceAddition In DataBase begins")

	//resourceAdditionManager接口对象
	resourceAdditionManager := &db.ResourceAdditionDbManager{}

	code, err := resourceAdditionManager.InsertResourceAddition(fileId, fileContentCount)
	if err != nil {
		logs.GetLogger().Error("InsertResourceAddition dbm err", err)
		return errorcode.UPLOAD_RESOURCE_DATEBASE_ERROR, err.Error(), ZERO
	}
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响", ZERO
	}
	logs.GetLogger().Info("InsertResourceAddition  In DataBase Success")
	id := int32(code)
	return errorcode.SUCCESS, "", id
}

//删除上传文件资源
//Param fileId
//return code errMsg
func (roomProcess *ClassRoomProcess) DeleteResource(fileId int32) (int32, string) {
	logs.GetLogger().Info("DeleteResource In DataBase begins")

	//ResourceDbManager接口对象
	resourceDBManager := &db.ResourceDbManager{}

	code, err := resourceDBManager.DeleteResource(fileId)
	if err != nil {
		logs.GetLogger().Error("DeleteResource dbm err", err)
		return errorcode.DELETE_RESOURCE_DATEBASE_ERROR, err.Error()
	}
	if code == 0 {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("DeleteResource  In DataBase Success")
	return errorcode.SUCCESS, ""
}
