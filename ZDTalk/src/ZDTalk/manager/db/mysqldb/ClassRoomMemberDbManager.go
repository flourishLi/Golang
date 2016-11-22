package mysqldb

import (
	database "ZDTalk/db/database"
	dbInfo "ZDTalk/manager/db/info"
	logs "ZDTalk/utils/log4go"
	"database/sql"
	"fmt"
)

type ClassRoomMemberDbManager struct {
}

//查询用户列表 即userinfo数据表的所有记录
//Param 空
//return ClassRoomMemberDbInfo切片 error
func (manager ClassRoomMemberDbManager) LoadAllUser() (map[int32]*dbInfo.ClassRoomMemberDbInfo, error) {
	logs.GetLogger().Info("开始读取数据库所有用户 LoadAllUser 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("读取数据库所有用户 错误 LoadAllUser: %s", r)
			return
		}
	}()

	sqls := "select * from userinfo"

	stmt, err := database.DBConn.Prepare(sqls)
	if err != nil {
		logs.GetLogger().Error("LoadAllUser Prepare Error:" + err.Error())
		return nil, err
	}
	rows, errTwo := stmt.Query()
	if errTwo != nil {
		logs.GetLogger().Error("LoadAllUser Query error:" + err.Error())
		return nil, err
	}
	defer rows.Close()
	//定义切片 存储ClassRoomMemberDbInfo
	userInfoSlice := make(map[int32]*dbInfo.ClassRoomMemberDbInfo)

	for rows.Next() {
		var userIdSql sql.NullInt64
		var chatIdSql sql.NullInt64
		var roleSql sql.NullInt64
		var deviceTypeSql sql.NullInt64
		var loginNameSql sql.NullString
		var userNameSql sql.NullString
		var userIconSql sql.NullString
		var passwordSql sql.NullString
		var yyTokenSql sql.NullString

		ClassRoomListBuf := make([]byte, 0)

		rows.Scan(&userIdSql, &chatIdSql, &roleSql, &deviceTypeSql, &loginNameSql, &userNameSql, &userIconSql, &passwordSql, &yyTokenSql, &ClassRoomListBuf)

		memberInfo := &dbInfo.ClassRoomMemberDbInfo{}
		if userIdSql.Valid {
			memberInfo.UserId = int32(userIdSql.Int64)
			//fmt.Println("userIdSql ", memberInfo.UserId)
		}
		if chatIdSql.Valid {
			memberInfo.ChatId = int32(chatIdSql.Int64)
			//fmt.Println("chatIdSql ", memberInfo.ChatId)
		}
		if roleSql.Valid {
			memberInfo.Role = int32(roleSql.Int64)
			//fmt.Println("roleSql ", memberInfo.Role)
		}
		if deviceTypeSql.Valid {
			memberInfo.DeviceType = int32(deviceTypeSql.Int64)
			//fmt.Println("deviceTypeSql ", memberInfo.DeviceType)
		}
		if loginNameSql.Valid {
			memberInfo.LoginName = loginNameSql.String
			//fmt.Println("loginNameSql ", memberInfo.LoginName)
		}
		if userNameSql.Valid {
			memberInfo.UserName = userNameSql.String
			//fmt.Println("userNameSql ", memberInfo.UserName)
		}
		if userIconSql.Valid {
			memberInfo.UserIcon = userIconSql.String
			//fmt.Println("userIconSql ", memberInfo.UserIcon)
		}
		if passwordSql.Valid {
			memberInfo.Password = passwordSql.String
			//fmt.Println("passwordSql ", memberInfo.Password)
		}
		if yyTokenSql.Valid {
			memberInfo.YYToken = yyTokenSql.String
			//fmt.Println("yYTokenSql ", memberInfo.YYToken)
		}
		if ClassRoomListBuf != nil && len(ClassRoomListBuf) > 0 {
			memberInfo.ClassRoomList = memberInfo.FromBytes(ClassRoomListBuf)
			//fmt.Println("ClassRoomListBuf ", memberInfo.ClassRoomList)
		}
		userInfoSlice[memberInfo.UserId] = memberInfo
	}
	return userInfoSlice, nil
}

//查询用户信息
//Param userID
//return ClassRoomMemberDbInfo error
func (manager ClassRoomMemberDbManager) GetUserInfoInfo(userId int32) (memberInfo *dbInfo.ClassRoomMemberDbInfo, err error) {
	logs.GetLogger().Info("开始读取数据库用户信息 GetUserInfo 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("读取数据库用户信息 错误 GetUserInfo: %s", r)
			return
		}
	}()

	sqls := "select * from userinfo where UserId=?"

	stmt, err := database.DBConn.Prepare(sqls)
	if err != nil {
		logs.GetLogger().Error("GetUserInfo Prepare Error:" + err.Error())
		return nil, err
	}
	rows, errTwo := stmt.Query(userId)
	if errTwo != nil {
		logs.GetLogger().Error("GetUserInfo Query error:" + err.Error())
		//fmt.Println("LoadAllGroups Query error:" + err.Error())
		return nil, err
	}
	defer rows.Close()

	memberInfo = &dbInfo.ClassRoomMemberDbInfo{}

	for rows.Next() {
		var userIdSql sql.NullInt64
		var chatIdSql sql.NullInt64
		var roleSql sql.NullInt64
		var deviceTypeSql sql.NullInt64
		var loginNameSql sql.NullString
		var userNameSql sql.NullString
		var userIconSql sql.NullString
		var passwordSql sql.NullString
		var yyTokenSql sql.NullString

		ClassRoomListBuf := make([]byte, 0)

		rows.Scan(&userIdSql, &chatIdSql, &roleSql, &deviceTypeSql, &loginNameSql, &userNameSql, &userIconSql, &passwordSql, &yyTokenSql, &ClassRoomListBuf)

		if userIdSql.Valid {
			memberInfo.UserId = int32(userIdSql.Int64)
			fmt.Println("userIdSql ", memberInfo.UserId)
		}
		if chatIdSql.Valid {
			memberInfo.ChatId = int32(chatIdSql.Int64)
			fmt.Println("chatIdSql ", memberInfo.ChatId)
		}
		if roleSql.Valid {
			memberInfo.Role = int32(roleSql.Int64)
			fmt.Println("roleSql ", memberInfo.Role)
		}
		if deviceTypeSql.Valid {
			memberInfo.DeviceType = int32(deviceTypeSql.Int64)
			fmt.Println("deviceTypeSql ", memberInfo.DeviceType)
		}
		if loginNameSql.Valid {
			memberInfo.LoginName = loginNameSql.String
			fmt.Println("loginNameSql ", memberInfo.LoginName)
		}
		if userNameSql.Valid {
			memberInfo.UserName = userNameSql.String
			fmt.Println("userNameSql ", memberInfo.UserName)
		}
		if userIconSql.Valid {
			memberInfo.UserIcon = userIconSql.String
			fmt.Println("userIconSql ", memberInfo.UserIcon)
		}
		if passwordSql.Valid {
			memberInfo.Password = passwordSql.String
			fmt.Println("passwordSql ", memberInfo.Password)
		}
		if yyTokenSql.Valid {
			memberInfo.YYToken = yyTokenSql.String
			fmt.Println("yyTokenSql ", memberInfo.YYToken)
		}
		if ClassRoomListBuf != nil && len(ClassRoomListBuf) > 0 {
			memberInfo.ClassRoomList = memberInfo.FromBytes(ClassRoomListBuf)
			fmt.Println("ClassRoomListBuf ", memberInfo.ClassRoomList)
		}
		logs.GetLogger().Info("GetUserInfo memberInfo", memberInfo)
	}
	return memberInfo, nil
}

//模糊查询用户信息 用户昵称
//返回[] *userInfo,error

func (manager ClassRoomMemberDbManager) FuzzyQueryUserInfo(key string) (map[int32]*dbInfo.ClassRoomMemberDbInfo, error) {

	logs.GetLogger().Info("用户昵称模糊查询 LoadAllUser 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("模糊查询 错误 LoadAllUser: %s", r)
			return
		}
	}()

	strsql := "select * from userinfo where UserName like ?"
	stmt, err := database.DBConn.Prepare(strsql)
	if err != nil {
		logs.GetLogger().Error("FuzzyQueryUserInfo Prepare error:" + err.Error())
	}
	rows, errTwo := stmt.Query("%" + key + "%")
	if errTwo != nil {
		logs.GetLogger().Error("FuzzyQueryUserInfo Query error:" + errTwo.Error())
	}
	defer rows.Close()

	//定义切片 存储ClassRoomMemberDbInfo
	userInfoSlice := make(map[int32]*dbInfo.ClassRoomMemberDbInfo)

	for rows.Next() {
		var userIdSql sql.NullInt64
		var chatIdSql sql.NullInt64
		var roleSql sql.NullInt64
		var deviceTypeSql sql.NullInt64
		var loginNameSql sql.NullString
		var userNameSql sql.NullString
		var userIconSql sql.NullString
		var passwordSql sql.NullString
		var yyTokenSql sql.NullString

		ClassRoomListBuf := make([]byte, 0)

		rows.Scan(&userIdSql, &chatIdSql, &roleSql, &deviceTypeSql, &loginNameSql, &userNameSql, &userIconSql, &passwordSql, &yyTokenSql, &ClassRoomListBuf)

		memberInfo := &dbInfo.ClassRoomMemberDbInfo{}
		if userIdSql.Valid {
			memberInfo.UserId = int32(userIdSql.Int64)
			//fmt.Println("userIdSql ", memberInfo.UserId)
		}
		if chatIdSql.Valid {
			memberInfo.ChatId = int32(chatIdSql.Int64)
			//fmt.Println("chatIdSql ", memberInfo.ChatId)
		}
		if roleSql.Valid {
			memberInfo.Role = int32(roleSql.Int64)
			//fmt.Println("roleSql ", memberInfo.Role)
		}
		if deviceTypeSql.Valid {
			memberInfo.DeviceType = int32(deviceTypeSql.Int64)
			//fmt.Println("deviceTypeSql ", memberInfo.DeviceType)
		}
		if loginNameSql.Valid {
			memberInfo.LoginName = loginNameSql.String
			//fmt.Println("loginNameSql ", memberInfo.LoginName)
		}
		if userNameSql.Valid {
			memberInfo.UserName = userNameSql.String
			//fmt.Println("userNameSql ", memberInfo.UserName)
		}
		if userIconSql.Valid {
			memberInfo.UserIcon = userIconSql.String
			//fmt.Println("userIconSql ", memberInfo.UserIcon)
		}
		if passwordSql.Valid {
			memberInfo.Password = passwordSql.String
			//fmt.Println("passwordSql ", memberInfo.Password)
		}
		if yyTokenSql.Valid {
			memberInfo.YYToken = yyTokenSql.String
			//fmt.Println("yYTokenSql ", memberInfo.YYToken)
		}
		if ClassRoomListBuf != nil && len(ClassRoomListBuf) > 0 {
			memberInfo.ClassRoomList = memberInfo.FromBytes(ClassRoomListBuf)
			//fmt.Println("ClassRoomListBuf ", memberInfo.ClassRoomList)
		}
		userInfoSlice[memberInfo.UserId] = memberInfo
	}
	return userInfoSlice, nil

}

//登录
//Param loginName Password
//return 1 成功 error
func (manager ClassRoomMemberDbManager) Login(loginName, password string) (*dbInfo.ClassRoomMemberDbInfo, error) {
	logs.GetLogger().Info("登录开始 Login 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("登录 错误 Login: %s", r)
			return
		}
	}()

	sqls := "select * from userinfo where LoginName=? and Password=?"

	stmt, err := database.DBConn.Prepare(sqls)
	if err != nil {
		logs.GetLogger().Error("Login Prepare Error:" + err.Error())
		return nil, err
	}
	rows, errTwo := stmt.Query(loginName, password)
	if errTwo != nil {
		logs.GetLogger().Error("Login Query error:" + err.Error())
		return nil, err
	}
	defer rows.Close()
	memberInfo := &dbInfo.ClassRoomMemberDbInfo{}

	for rows.Next() {
		var userIdSql sql.NullInt64
		var chatIdSql sql.NullInt64
		var roleSql sql.NullInt64
		var deviceTypeSql sql.NullInt64
		var loginNameSql sql.NullString
		var userNameSql sql.NullString
		var userIconSql sql.NullString
		var passwordSql sql.NullString
		var yYTokenSql sql.NullString

		ClassRoomListBuf := make([]byte, 0)

		rows.Scan(&userIdSql, &chatIdSql, &roleSql, &deviceTypeSql, &loginNameSql, &userNameSql, &userIconSql, &passwordSql, &yYTokenSql, &ClassRoomListBuf)

		if userIdSql.Valid {
			memberInfo.UserId = int32(userIdSql.Int64)
			fmt.Println("userIdSql ", memberInfo.UserId)
		}
		if chatIdSql.Valid {
			memberInfo.ChatId = int32(chatIdSql.Int64)
			fmt.Println("chatIdSql ", memberInfo.ChatId)
		}
		if roleSql.Valid {
			memberInfo.Role = int32(roleSql.Int64)
			fmt.Println("roleSql ", memberInfo.Role)
		}
		if deviceTypeSql.Valid {
			memberInfo.DeviceType = int32(deviceTypeSql.Int64)
			fmt.Println("deviceTypeSql ", memberInfo.DeviceType)
		}
		if loginNameSql.Valid {
			memberInfo.LoginName = loginNameSql.String
			fmt.Println("loginNameSql ", memberInfo.LoginName)
		}
		if userNameSql.Valid {
			memberInfo.UserName = userNameSql.String
			fmt.Println("userNameSql ", memberInfo.UserName)
		}
		if userIconSql.Valid {
			memberInfo.UserIcon = userIconSql.String
			fmt.Println("userIconSql ", memberInfo.UserIcon)
		}
		if passwordSql.Valid {
			memberInfo.Password = passwordSql.String
			fmt.Println("passwordSql ", memberInfo.Password)
		}
		if yYTokenSql.Valid {
			memberInfo.YYToken = yYTokenSql.String
			fmt.Println("yYTokenSql ", memberInfo.YYToken)
		}
		if ClassRoomListBuf != nil && len(ClassRoomListBuf) > 0 {
			memberInfo.ClassRoomList = memberInfo.FromBytes(ClassRoomListBuf)
			fmt.Println("ClassRoomListBuf ", memberInfo.ClassRoomList)
		}
		logs.GetLogger().Info("GetUserInfo memberInfo", memberInfo)
	}
	return memberInfo, nil
}

//用户注册
//Param ChatId Role DeviceType LoginName UserName UserIcon Password
//return userId err
func (manager ClassRoomMemberDbManager) SignUp(ChatId, Role, DeviceType int32, LoginName, UserName, UserIcon, Password string) (int64, error) {
	logs.GetLogger().Info("开始创建用户 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("创建用户 错误: %s", r)
			return
		}
	}()

	sqlInsert := "insert into userinfo(ChatId,Role,DeviceType,LoginName,UserName,UserIcon,Password)values(?,?,?,?,?,?,?)"

	stmt, err := database.DBConn.Prepare(sqlInsert)
	if err != nil {
		logs.GetLogger().Error("CreateUser Prepare error:" + err.Error())
		//fmt.Println("CreateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(ChatId, Role, DeviceType, LoginName, UserName, UserIcon, Password)
	if errTwo != nil {
		logs.GetLogger().Error("CreateUser Exec error:" + errTwo.Error())
		//fmt.Println("CreateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	id, errThree := result.LastInsertId()
	if errThree != nil {
		logs.GetLogger().Error("CreateUser LastInsertId error:" + errThree.Error())
		//fmt.Println("CreateClassRoom LastInsertId error:" + err.Error())
		return ZERO, errThree
	}
	return id, nil
}

//修改用户信息
//Param UserID Role DeviceType UserName UserIcon
//return int 1 Success err
func (manager ClassRoomMemberDbManager) UserInfoUpdate(UserID, Role, DeviceType int32, UserName, UserIcon string) (int64, error) {
	logs.GetLogger().Info("开始修改用户信息 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("修改用户信息 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update userinfo set Role =?,DeviceType=?,UserName=?,UserIcon=? where UserId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UserInfoUpdate Prepare error:" + err.Error())
		//fmt.Println("CreateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(Role, DeviceType, UserName, UserIcon, UserID)
	if errTwo != nil {
		logs.GetLogger().Error("UserInfoUpdate Exec error:" + errTwo.Error())
		//fmt.Println("CreateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affectRow, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UserInfoUpdate affectRow error:" + errThree.Error())
		//fmt.Println("CreateClassRoom LastInsertId error:" + err.Error())
		return ZERO, errThree
	}
	return affectRow, nil
}

//修改用户密码
//Param UserID oldPassword newPassword
//return int 1 Success err
func (manager ClassRoomMemberDbManager) UserPwdUpdate(UserID int32, newPassword, oldPassword string) (int64, error) {
	logs.GetLogger().Info("开始修改用户密码 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("修改用户密码 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update userinfo set Password =? where UserId = ? and Password = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UserPwdUpdate Prepare error:" + err.Error())
		//fmt.Println("CreateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(newPassword, UserID, oldPassword)
	if errTwo != nil {
		logs.GetLogger().Error("UserPwdUpdate Exec error:" + errTwo.Error())
		//fmt.Println("CreateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affectRow, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UserPwdUpdate affectRow error:" + errThree.Error())
		//fmt.Println("CreateClassRoom LastInsertId error:" + err.Error())
		return ZERO, errThree
	}
	return affectRow, nil
}

//重置用户密码
//Param UserID newPassword
//return int 1 Success err
func (manager ClassRoomMemberDbManager) UserPwdReset(UserID int32, newPassword string) (int64, error) {
	logs.GetLogger().Info("开始重置用户密码 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("重置用户密码 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update userinfo set Password =? where UserId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("UserPwdReset Prepare error:" + err.Error())
		//fmt.Println("CreateClassRoom Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(newPassword, UserID)
	if errTwo != nil {
		logs.GetLogger().Error("UserPwdReset Exec error:" + errTwo.Error())
		//fmt.Println("CreateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affectRow, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("UserPwdReset affectRow error:" + errThree.Error())
		//fmt.Println("CreateClassRoom LastInsertId error:" + err.Error())
		return ZERO, errThree
	}
	return affectRow, nil
}
