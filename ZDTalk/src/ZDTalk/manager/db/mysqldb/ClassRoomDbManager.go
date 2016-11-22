package mysqldb

import (
	database "ZDTalk/db/database"
	dbInfo "ZDTalk/manager/db/info"
	logs "ZDTalk/utils/log4go"
	"database/sql"
	"fmt"
)

const (
	ZERO = 0
	One  = 1
)

type ClassRoomDbManager struct {
}

//查询教室列表 即classroominfo数据表的所有记录
//Param 空
//return ClassRoomDbInfo切片 error
func (manager ClassRoomDbManager) LoadAllClassRoom() (map[int32]*dbInfo.ClassRoomDbInfo, error) {
	logs.GetLogger().Info("开始读取数据库所有群 loadAllClassRoom 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("加载数据库所有群 错误 LoadGroups: %s", r)
			return
		}
	}()

	sqls := "select * from classroominfo"
	//logs.GetLogger().Info("sqls " + sqls)
	stmt, err := database.DBConn.Prepare(sqls)
	if err != nil {
		logs.GetLogger().Error("LoadAllGroups Prepare Error:" + err.Error())
		return nil, err
	}
	rows, errTwo := stmt.Query()
	if errTwo != nil {
		logs.GetLogger().Error("LoadAllGroups Query error:" + err.Error())
		//fmt.Println("LoadAllGroups Query error:" + err.Error())
		return nil, err
	}
	defer rows.Close()
	//定义切片 存储ClassRoomDbInfo
	classRoomSlice := make(map[int32]*dbInfo.ClassRoomDbInfo)
	//	logs.GetLogger().Info(rows.Next())
	for rows.Next() {
		var classRoomIdSql sql.NullInt64
		var classRoomIMIdSql sql.NullInt64
		var creatorUserIdSql sql.NullInt64
		var classRoomStatusSql sql.NullInt64
		var createTimeSql sql.NullInt64
		var classRoomNameSql sql.NullString
		var classRoomLogoSql sql.NullString
		var descriptionSql sql.NullString
		var classRoomCourseSql sql.NullString
		var forbidHandStatusSql sql.NullInt64

		settingStatusBuf := make([]byte, 0)
		memberBuf := make([]byte, 0)
		onLineMemberBuf := make([]byte, 0)
		handMemberBuf := make([]byte, 0)
		forbidSayMemberBuf := make([]byte, 0)
		sayingMemberBuf := make([]byte, 0)

		rows.Scan(&classRoomIdSql, &classRoomIMIdSql, &creatorUserIdSql, &classRoomStatusSql, &settingStatusBuf, &createTimeSql, &classRoomNameSql, &classRoomLogoSql, &descriptionSql, &classRoomCourseSql, &memberBuf, &onLineMemberBuf, &handMemberBuf, &forbidSayMemberBuf, &sayingMemberBuf, &forbidHandStatusSql)

		roomInfo := &dbInfo.ClassRoomDbInfo{}
		if classRoomIdSql.Valid {
			roomInfo.ClassRoomId = int32(classRoomIdSql.Int64)
			//fmt.Println("classRoomIdSql ", roomInfo.ClassRoomId)
		}
		if classRoomIMIdSql.Valid {
			roomInfo.ClassRoomIMId = int32(classRoomIMIdSql.Int64)
			//fmt.Println("classRoomIMIdSql ", roomInfo.ClassRoomId)
		}
		if creatorUserIdSql.Valid {
			roomInfo.CreatorUserId = int32(creatorUserIdSql.Int64)
			//fmt.Println("creatorUserIdSql ", roomInfo.CreatorUserId)
		}
		if classRoomStatusSql.Valid {
			roomInfo.ClassRoomStatus = int32(classRoomStatusSql.Int64)
			//fmt.Println("classRoomStatusSql ", roomInfo.ClassRoomStatus)
		}
		if createTimeSql.Valid {
			roomInfo.CreateTime = createTimeSql.Int64
			//fmt.Println("createTimeSql ", roomInfo.CreateTime)
		}
		if classRoomNameSql.Valid {
			roomInfo.ClassRoomName = classRoomNameSql.String
			//fmt.Println("classRoomIMIdSql ", roomInfo.ClassRoomName)
		}
		if classRoomLogoSql.Valid {
			roomInfo.ClassRoomLogo = classRoomLogoSql.String
			//fmt.Println("classRoomLogoSql ", roomInfo.ClassRoomLogo)
		}
		if descriptionSql.Valid {
			roomInfo.Description = descriptionSql.String
			//fmt.Println("descriptionSql ", roomInfo.Description)
		}
		if classRoomCourseSql.Valid {
			roomInfo.ClassRoomCourse = classRoomCourseSql.String
			//fmt.Println("classRoomCourseSql ", roomInfo.ClassRoomCourse)
		}
		if settingStatusBuf != nil && len(settingStatusBuf) > 0 {
			roomInfo.SettingStatus = roomInfo.FromBytes(settingStatusBuf)
			//fmt.Println("settingStatusBuf ", roomInfo.SettingStatus)
		}
		if memberBuf != nil && len(memberBuf) > 0 {
			roomInfo.MemberList = roomInfo.FromBytes(memberBuf)
			//fmt.Println("memberBuf ", roomInfo.MemberList)
		}
		if onLineMemberBuf != nil && len(onLineMemberBuf) > 0 {
			roomInfo.OnLineMemberList = roomInfo.FromBytes(onLineMemberBuf)
			//fmt.Println("onLineMemberBuf ", roomInfo.OnLineMemberList)
		}
		if handMemberBuf != nil && len(handMemberBuf) > 0 {
			roomInfo.HandMemberList = roomInfo.StructFromBytes(handMemberBuf)
			//fmt.Println("handMemberBuf ", roomInfo.HandMemberList)
		}
		if forbidSayMemberBuf != nil && len(forbidSayMemberBuf) > 0 {
			roomInfo.ForbidSayMemberList = roomInfo.FromBytes(forbidSayMemberBuf)
			//fmt.Println("forbidSayMemberBuf ", roomInfo.ForbidSayMemberList)
		}
		if forbidHandStatusSql.Valid {
			roomInfo.ForbidHandStatus = int32(forbidHandStatusSql.Int64)
			//fmt.Println("forbidSayMemberBuf ", roomInfo.ForbidSayMemberList)
		}
		if sayingMemberBuf != nil && len(sayingMemberBuf) > 0 {
			roomInfo.SayingMemberList = roomInfo.FromBytes(sayingMemberBuf)
			//fmt.Println("sayingMemberBuf ", roomInfo.SayingMemberList)
		}
		classRoomSlice[roomInfo.ClassRoomId] = roomInfo
	}
	return classRoomSlice, nil
}

//查询教室信息
//Param roomID 教室ID
//return roomInfo *dbInfo.ClassRoomDbInfo description教室说明 error
func (manager ClassRoomDbManager) GetClassRoomInfo(roomID int32) (roomInfo *dbInfo.ClassRoomDbInfo, err error) {
	logs.GetLogger().Info("开始读取教室信息 GetClassRoomInfo 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("读取教室信息 错误 GetClassRoomInfo: %s", r)
			return
		}
	}()
	roomInfo = &dbInfo.ClassRoomDbInfo{}
	sqlQuery := "select ClassRoomId,ClassRoomIMId,CreatorUserId,ClassRoomStatus,SettingStatus,CreateTime,ClassRoomName,ClassRoomLogo,Description,ClassRoomCourse,MemberList,OnLineMemberList,HandMemberList,ForbidSayMemberList,SayingMemberList from classroominfo where ClassRoomId=? "
	stmt, err := database.DBConn.Prepare(sqlQuery)
	if err != nil {
		logs.GetLogger().Error("GetClassRoomInfo Prepare error:" + err.Error())
		return nil, err
	}
	rows, errTwo := stmt.Query(roomID)
	if errTwo != nil {
		logs.GetLogger().Error("GetClassRoomInfo Query error:" + errTwo.Error())
		return nil, errTwo
	}
	defer rows.Close()

	for rows.Next() {
		var classRoomIdSql sql.NullInt64
		var classRoomIMIdSql sql.NullInt64
		var creatorUserIdSql sql.NullInt64
		var classRoomStatusSql sql.NullInt64
		var createTimeSql sql.NullInt64
		var classRoomNameSql sql.NullString
		var classRoomLogoSql sql.NullString
		var descriptionSql sql.NullString
		var classRoomCourseSql sql.NullString
		var forbidHandStatusSql sql.NullInt64

		settingStatusBuf := make([]byte, 0)
		memberBuf := make([]byte, 0)
		onLineMemberBuf := make([]byte, 0)
		handMemberBuf := make([]byte, 0)
		forbidSayMemberBuf := make([]byte, 0)

		sayingMemberBuf := make([]byte, 0)

		rows.Scan(&classRoomIdSql, &classRoomIMIdSql, &creatorUserIdSql, &classRoomStatusSql, &settingStatusBuf, &createTimeSql, &classRoomNameSql, &classRoomLogoSql, &descriptionSql, &classRoomCourseSql, &memberBuf, &onLineMemberBuf, &handMemberBuf, &forbidSayMemberBuf, &sayingMemberBuf, &forbidHandStatusSql)

		if classRoomIdSql.Valid {
			roomInfo.ClassRoomId = int32(classRoomIdSql.Int64)
			fmt.Println("classRoomIdSql ", roomInfo.ClassRoomId)
		}
		if classRoomIMIdSql.Valid {
			roomInfo.ClassRoomIMId = int32(classRoomIMIdSql.Int64)
			fmt.Println("classRoomIMIdSql ", roomInfo.ClassRoomId)
		}
		if creatorUserIdSql.Valid {
			roomInfo.CreatorUserId = int32(creatorUserIdSql.Int64)
			fmt.Println("creatorUserIdSql ", roomInfo.CreatorUserId)
		}
		if classRoomStatusSql.Valid {
			roomInfo.ClassRoomStatus = int32(classRoomStatusSql.Int64)
			fmt.Println("classRoomStatusSql ", roomInfo.ClassRoomStatus)
		}
		if createTimeSql.Valid {
			roomInfo.CreateTime = createTimeSql.Int64
			fmt.Println("createTimeSql ", roomInfo.CreateTime)
		}
		if classRoomNameSql.Valid {
			roomInfo.ClassRoomName = classRoomNameSql.String
			fmt.Println("classRoomIMIdSql ", roomInfo.ClassRoomName)
		}
		if classRoomLogoSql.Valid {
			roomInfo.ClassRoomLogo = classRoomLogoSql.String
			fmt.Println("classRoomLogoSql ", roomInfo.ClassRoomLogo)
		}
		if descriptionSql.Valid {
			roomInfo.Description = descriptionSql.String
			fmt.Println("descriptionSql ", roomInfo.Description)
		}
		if classRoomCourseSql.Valid {
			roomInfo.ClassRoomCourse = classRoomCourseSql.String
			fmt.Println("classRoomCourseSql ", roomInfo.ClassRoomCourse)
		}
		if settingStatusBuf != nil && len(settingStatusBuf) > 0 {
			roomInfo.SettingStatus = roomInfo.FromBytes(settingStatusBuf)
			fmt.Println("settingStatusBuf ", roomInfo.SettingStatus)
		}
		if memberBuf != nil && len(memberBuf) > 0 {
			roomInfo.MemberList = roomInfo.FromBytes(memberBuf)
			fmt.Println("memberBuf ", roomInfo.MemberList)
		}
		if onLineMemberBuf != nil && len(onLineMemberBuf) > 0 {
			roomInfo.OnLineMemberList = roomInfo.FromBytes(onLineMemberBuf)
			fmt.Println("onLineMemberBuf ", roomInfo.OnLineMemberList)
		}
		if handMemberBuf != nil && len(handMemberBuf) > 0 {
			roomInfo.HandMemberList = roomInfo.StructFromBytes(handMemberBuf)
			fmt.Println("handMemberBuf ", roomInfo.HandMemberList)
		}
		if forbidSayMemberBuf != nil && len(forbidSayMemberBuf) > 0 {
			roomInfo.ForbidSayMemberList = roomInfo.FromBytes(forbidSayMemberBuf)
			fmt.Println("forbidSayMemberBuf ", roomInfo.ForbidSayMemberList)
		}
		if forbidHandStatusSql.Valid {
			roomInfo.ForbidHandStatus = int32(forbidHandStatusSql.Int64)
			//fmt.Println("forbidSayMemberBuf ", roomInfo.ForbidSayMemberList)
		}
		if sayingMemberBuf != nil && len(sayingMemberBuf) > 0 {
			roomInfo.SayingMemberList = roomInfo.FromBytes(sayingMemberBuf)
			fmt.Println("sayingMemberBuf ", roomInfo.SayingMemberList)
		}
		logs.GetLogger().Info("ClassRoomInfo:", roomInfo)
	}
	return roomInfo, nil
}

//创建教室
//Param classRoomName教室名称 classRoomLogo教室头像 description教室说明 ClassRoomCourse课程信息 CreatorUserId创建者ID CreateTime创建时间 ClassRoomIMId 教室对应IM的Id ClassRoomStatus=0 刚创建 SettingStatus
//return roomId err
func (manager ClassRoomDbManager) CreateClassRoom(classRoomName, classRoomLogo, description, classRoomCourse string, creatorUserId int32, settingStatus []int32, createTime int64, classRoomIMId int32) (int64, error) {
	logs.GetLogger().Info("开始创建教室 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("创建教室 错误: %s", r)
			return
		}
	}()

	sqlInsert := "insert into classroominfo(ClassRoomName,ClassRoomLogo,Description,ClassRoomCourse,CreatorUserId,CreateTime,ClassRoomIMId,ClassRoomStatus,SettingStatus)values(?,?,?,?,?,?,?,?,?)"

	stmt, err := database.DBConn.Prepare(sqlInsert)
	if err != nil {
		logs.GetLogger().Error("CreateClassRoom Prepare error:" + err.Error())
		//fmt.Println("CreateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	//数据库类型为blob
	//定义接口对象
	memberinfo := &dbInfo.ClassRoomDbInfo{}
	//	memberlist := []int32{creatorUserId}
	//	buffList := memberinfo.ToBytes(memberlist)
	//举手列表
	//	handMember := dbInfo.HandsMember{creatorUserId, createTime}
	//	handList := []dbInfo.HandsMember{handMember}
	//	handBuf := memberinfo.StuctToBytes(handList)
	//SettingStatus
	settingStatusBuf := memberinfo.ToBytes(settingStatus)

	result, errTwo := stmt.Exec(classRoomName, classRoomLogo, description, classRoomCourse, creatorUserId, createTime, classRoomIMId, ZERO, settingStatusBuf)
	if errTwo != nil {
		logs.GetLogger().Error("CreateClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("CreateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	id, errThree := result.LastInsertId()
	if errThree != nil {
		logs.GetLogger().Error("CreateClassRoom LastInsertId error:" + errThree.Error())
		//fmt.Println("CreateClassRoom LastInsertId error:" + err.Error())
		return ZERO, errThree
	}
	return id, nil
}

//更新教室状态信息
//Param roomId教室ID classRoomName教室名称 classRoomLogo教室头像 description教室说明 course
//return int=1 更新成功 err
func (manager ClassRoomDbManager) UpdateClassRoom(classRoomID int32, classRoomName, classRoomLogo, description string, course string) (int64, error) {
	logs.GetLogger().Info("开始更新教室状态 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("更新教室状态 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set ClassRoomName =?,ClassRoomLogo=?,Description=?,ClassRoomCourse=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		//fmt.Println("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(classRoomName, classRoomLogo, description, course, classRoomID)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//教室设置 SettingStatus 教室设置状态 1 允许所有人打字，2 允许学生录音，3 禁止游客打字，4 禁止游客举手
//Param roomId教室ID teacherUserId教师ID settingStatus
//return int=1 更新成功 err
func (manager ClassRoomDbManager) UpdateClassRoomSettingStatus(classRoomID, teacherUserID int32, settingStatus []int32) (int64, error) {
	logs.GetLogger().Info("开始更新教室状态 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("更新教室状态 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set SettingStatus=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		//fmt.Println("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()
	//数据库类型为blob
	//定义接口对象
	memberinfo := &dbInfo.ClassRoomDbInfo{}
	//SettingStatus
	settingStatusBuf := memberinfo.ToBytes(settingStatus)
	result, errTwo := stmt.Exec(settingStatusBuf, classRoomID)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//老师上下课 ClassRoomStatus 教室当前状态0 刚创建, 1 初始状态, 2 上课状态, 3 下课状态
//Param roomId教室ID teacherUserId教师ID classroomStatus
//return int=1 更新成功 err
func (manager ClassRoomDbManager) UpdateClassRoomStatus(classRoomID, teacherUserID, classRoomStatus int32) (int64, error) {
	logs.GetLogger().Info("开始老师上下课设置 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("开始老师上下课设置 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set ClassRoomStatus=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		//fmt.Println("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(classRoomStatus, classRoomID)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//删除教室
//Param roomId教室ID
//Return int=1 删除成功 err
func (manager ClassRoomDbManager) DeleteClassRoom(classRoomID int32) (int64, error) {
	logs.GetLogger().Info("开始删除教室 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("删除教室状态 错误: %s", r)
			return
		}
	}()

	sqlDelete := "delete from classroominfo where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlDelete)
	if err != nil {
		logs.GetLogger().Error("DeleteClassRoom Prepare error:" + err.Error())
		//fmt.Println("DeleteClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(classRoomID)
	if errTwo != nil {
		logs.GetLogger().Error("DeleteClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("DeleteClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("DeleteClassRoom RowsAffected error:" + errThree.Error())
		//fmt.Println("DeleteClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//举手(取消)动作请求 前端只传userID 需要重新计算useIDlist
//Param classRoomId userIdList
//return int=1 更新成功 err
//维护数据库的handmemberList字段 举手列表
func (manager ClassRoomDbManager) HandsUp(userIdList []dbInfo.HandsMember, classRoomId int32) (int64, error) {
	//声明接口对象
	classRoomDbInfo := &dbInfo.ClassRoomDbInfo{}

	logs.GetLogger().Info("开始处理举手(取消)动作请求 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("举手(取消)动作请求 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set HandMemberList=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		//fmt.Println("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()
	//将struct转化成[]byte
	userIdListBuf := classRoomDbInfo.StuctToBytes(userIdList)

	result, errTwo := stmt.Exec(userIdListBuf, classRoomId)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//禁止 解禁请求 禁言区
//Param classRoomId studentIds集合 userIdList=studentIds+原有的studentid
//return int=1 更新成功 err
//维护数据库的ForbidSayMemberList字段 被禁言的列表
func (manager ClassRoomDbManager) ForbidArea(userIdList []int32, classRoomId int32) (int64, error) {
	//声明接口对象
	classRoomDbInfo := &dbInfo.ClassRoomDbInfo{}

	logs.GetLogger().Info("开始处理禁止解禁 禁言区请求 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("禁止解禁 禁言区请求 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set ForbidSayMemberList=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	//将[]int32转化成[]byte
	userIdListBuf := classRoomDbInfo.ToBytes(userIdList)
	result, errTwo := stmt.Exec(userIdListBuf, classRoomId)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//更新教室的禁止举手状态 教室的禁止举手状态 0可举手 1禁止举手
//Param classRoomId ForbidHandsStatus
//return int=1 更新成功 err
//维护数据库的ForbidHandStatus字段
func (manager ClassRoomDbManager) UpdateHandForbidStatus(classRoomId, ForbidHandsStatus int32) (int64, error) {

	logs.GetLogger().Info("开始处理教室禁止举手状态请求 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("教室禁止举手状态 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set ForbidHandStatus=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateHandForbidStatus Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(ForbidHandsStatus, classRoomId)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateHandForbidStatus Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateHandForbidStatus RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//添加 移除举手到发言区请求
//Param classRoomId studentIds集合 userIdList=studentIds+-原有的studentid
//return int=1 更新成功 err
//维护数据库的SayingMemberList字段 发言列表
func (manager ClassRoomDbManager) DeleteAddToSpeakArea(userIdList []int32, classRoomId int32) (int64, error) {
	//声明接口对象
	classRoomDbInfo := &dbInfo.ClassRoomDbInfo{}
	logs.GetLogger().Info("添加举手到发言区请求 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("添加举手到发言区 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set SayingMemberList=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	userIdListBuf := classRoomDbInfo.ToBytes(userIdList)
	result, errTwo := stmt.Exec(userIdListBuf, classRoomId)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//清空举手列表请求
//Param classRoomId
//return int=1 更新成功 err
//维护数据库的handmemberList字段 举手列表
func (manager ClassRoomDbManager) HandsListClear(classRoomId int32) (int64, error) {

	logs.GetLogger().Info("开始清空举手列表请求 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("清空举手列表 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set HandMemberList=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()
	handsListBuf := []byte{}
	result, errTwo := stmt.Exec(handsListBuf, classRoomId)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//将学生移除 添加教室举手请求
//Param classRoomId studentIds集合 userIdList=原有的studentid集合+-studentsId
//return int=1 更新成功 err
//维护数据库的MemberList字段 发言列表
func (manager ClassRoomDbManager) DeleteAddStudents(userIdList []int32, classRoomId int32) (int64, error) {
	//声明接口对象
	classRoomDbInfo := &dbInfo.ClassRoomDbInfo{}
	logs.GetLogger().Info("开始处理将学生移除教室举手请求 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("将学生移除教室举手请求 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set MemberList=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()
	userIdListBuf := classRoomDbInfo.ToBytes(userIdList)
	result, errTwo := stmt.Exec(userIdListBuf, classRoomId)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}

//学生进入教室请求
//Param classRoomId userId userIdList=原有的studentid集合+-studentsId
//return int=1 更新成功 err
//维护数据库的OnLineMemberList字段 在线成员列表
func (manager ClassRoomDbManager) EntryClassroom(userIdList []int32, classRoomId int32) (int64, error) {
	//声明接口对象
	classRoomDbInfo := &dbInfo.ClassRoomDbInfo{}
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("进入教室 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set OnLineMemberList=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UpdateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()
	userIdListBuf := classRoomDbInfo.ToBytes(userIdList)
	result, errTwo := stmt.Exec(userIdListBuf, classRoomId)
	if errTwo != nil {
		logs.GetLogger().Error("UpdateClassRoom Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UpdateClassRoom RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}

	return affect, nil
}

//退出教室请求
//Param classRoomId userId userIdList=原有的studentid集合+-studentsId
//return int=1 更新成功 err
//维护数据库 在线成员列表 举手列表 发言列表
func (manager ClassRoomDbManager) ExitClassroom(onlineUserIdList, handUserList, sayUserList []int32, classRoomId int32) (int64, error) {
	//声明接口对象
	classRoomDbInfo := &dbInfo.ClassRoomDbInfo{}
	logs.GetLogger().Info("开始处理将学生退出教室请求 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("退出教室 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update classroominfo set OnLineMemberList=?,HandMemberList=?,SayingMemberList=? where ClassRoomId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("ExitClassroom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	//数据转换
	onlineUserListBuf := classRoomDbInfo.ToBytes(onlineUserIdList)
	handUserListBuf := classRoomDbInfo.ToBytes(handUserList)
	sayUserListBuf := classRoomDbInfo.ToBytes(sayUserList)

	result, errTwo := stmt.Exec(onlineUserListBuf, handUserListBuf, sayUserListBuf, classRoomId)
	if errTwo != nil {
		logs.GetLogger().Error("ExitClassroom Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("ExitClassroom RowsAffected error:" + errThree.Error())
		//fmt.Println("UpdateClassRoom RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	return affect, nil
}
