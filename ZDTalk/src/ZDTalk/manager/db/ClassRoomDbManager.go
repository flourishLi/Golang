package db

import (
	"ZDTalk/db/database"
	dbInfo "ZDTalk/manager/db/info"
	"ZDTalk/manager/db/mysqldb"
)

const (
	ZERO = 0
)

//函数的接收者和函数的定义必须在一个包内
type ClassRoomDbManager struct {
}

//加载所有教室 教室列表
//Param 空
//return ChassRoomDbInfo切片 error
func (manager ClassRoomDbManager) LoadAllClassRoom() (map[int32]*dbInfo.ClassRoomDbInfo, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.LoadAllClassRoom()
	}
	return nil, nil
}

//查询教室信息
//Param roomID 教室ID
//return *dbInfo.ClassRoomDbInfo error
func (manager ClassRoomDbManager) GetClassRoomInfo(roomID int32) (*dbInfo.ClassRoomDbInfo, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.GetClassRoomInfo(roomID)
	}
	return nil, nil
}

//创建教室
//Param classRoomName教室名称 classRoomLogo教室头像 description教室说明 ClassRoomCourse课程信息 CreatorUserId创建者ID CreateTime创建时间 ClassRoomIMId 教室对应IM的Id
//return roomId err
func (manager ClassRoomDbManager) CreateClassRoom(classRoomName, classRoomLogo, description, classRoomCourse string, creatorUserId int32, settingStatus []int32, createTime int64, classRoomIMId int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.CreateClassRoom(classRoomName, classRoomLogo, description, classRoomCourse, creatorUserId, settingStatus, createTime, classRoomIMId)
	}
	return ZERO, nil
}

//更新教室状态信息
//Param roomId教室ID classRoomName教室名称 classRoomLogo教室头像 description教室说明 course
//return int=1 更新成功 err
func (manager ClassRoomDbManager) UpdateClassRoom(classRoomID int32, classRoomName, classRoomLogo, description string, course string) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.UpdateClassRoom(classRoomID, classRoomName, classRoomLogo, description, course)
	}
	return ZERO, nil
}

//教室设置
//Param roomId教室ID teacherUserId教师ID settingStatus
//return int=1 更新成功 err
func (manager ClassRoomDbManager) UpdateClassRoomSettingStatus(classRoomID, teacherUserID int32, settingStatus []int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.UpdateClassRoomSettingStatus(classRoomID, teacherUserID, settingStatus)
	}
	return ZERO, nil
}

//老师上下课 roomStatus设置
//Param roomId教室ID teacherUserId教师ID classRoomStatus
//return int=1 更新成功 err
func (manager ClassRoomDbManager) UpdateClassRoomStatus(classRoomID, teacherUserID, classRoomStatus int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.UpdateClassRoomStatus(classRoomID, teacherUserID, classRoomStatus)
	}
	return ZERO, nil
}

//教室的禁止举手状态设置
//Param roomId教室ID  ForbidHandsStatus
//return int=1 更新成功 err
func (manager ClassRoomDbManager) UpdateHandForbidStatus(classRoomID, forbidHandsStatus int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.UpdateHandForbidStatus(classRoomID, forbidHandsStatus)
	}
	return ZERO, nil
}

//删除教室
//Param roomId教室ID
//Return int=1 删除成功 err
func (manager ClassRoomDbManager) DeleteClassRoom(classRoomID int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.DeleteClassRoom(classRoomID)
	}
	return ZERO, nil
}

//举手(取消)动作请求
//Param classRoomId userIdList
//return int=1 更新成功 err
func (manager ClassRoomDbManager) HandsUp(userIdList []dbInfo.HandsMember, classRoomID int32) (int64, error) {

	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.HandsUp(userIdList, classRoomID)
	}
	return ZERO, nil
}

////禁止 解禁请求 禁言区
//Param classRoomId userIdList
//return int=1 更新成功 err
func (manager ClassRoomDbManager) ForbidArea(userIdList []int32, classRoomID int32) (int64, error) {

	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.ForbidArea(userIdList, classRoomID)
	}
	return ZERO, nil
}

//清空举手列表动作请求
//Param classRoomId
//return int=1 更新成功 err
func (manager ClassRoomDbManager) HandsListClear(classRoomID int32) (int64, error) {

	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.HandsListClear(classRoomID)
	}
	return ZERO, nil
}

//添加举手到发言区动作请求
//Param classRoomId userIdList
//return int=1 更新成功 err
func (manager ClassRoomDbManager) DeleteAddToSpeakArea(userIdList []int32, classRoomID int32) (int64, error) {

	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.DeleteAddToSpeakArea(userIdList, classRoomID)
	}
	return ZERO, nil
}

//将学生添加 移除教室动作请求
//Param classRoomId studentIds集合 userIdList=原有的studentid集合-studentsId
//return int=1 更新成功 err
func (manager ClassRoomDbManager) DeleteAddStudents(userIdList []int32, classRoomID int32) (int64, error) {

	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.DeleteAddStudents(userIdList, classRoomID)
	}
	return ZERO, nil
}

//学生进入教室请求 及实时更新在线成员列表
//Param classRoomId userId userIdList=原有的studentid集合+-studentsId
//return int=1 更新成功 err
//维护数据库的OnLineMemberList字段 在线成员列表
func (manager ClassRoomDbManager) EntryClassroom(userIdList []int32, classRoomID int32) (int64, error) {

	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.EntryClassroom(userIdList, classRoomID)
	}
	return ZERO, nil
}

//退出教室请求
//Param classRoomId userId userIdList=原有的studentid集合+-studentsId
//return int=1 更新成功 err
//维护数据库 在线成员列表 举手列表 发言列表
func (manager ClassRoomDbManager) ExitClassroom(onlineUserList, handUserList, sayUserlist []int32, classRoomID int32) (int64, error) {

	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomDbManager{}.ExitClassroom(onlineUserList, handUserList, sayUserlist, classRoomID)
	}
	return ZERO, nil
}
