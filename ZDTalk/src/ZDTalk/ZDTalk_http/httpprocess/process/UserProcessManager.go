package process

import (
	"ZDTalk/errorcode"
	"ZDTalk/manager/db"
	"ZDTalk/manager/db/info"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
)

type UserProcess struct {
}

//登录
//Param loginName Password
//return code errMsg UserInfo
func (userProcess *UserProcess) Login(loginName, password string) (int32, string, *info.ClassRoomMemberDbInfo) {
	logs.GetLogger().Info("Login begins")
	//userDBManager 接口对象
	userDBManager := &db.ClassRoomDbMemberManager{}

	userInfo, err := userDBManager.Login(loginName, password)
	if err != nil {
		logs.GetLogger().Error("Login dbm err", err)
		return errorcode.LOGIN_DATABASE_IS_ERROR, err.Error(), nil
	}
	if userInfo != nil {
		logs.GetLogger().Error("name or password is error")
		return errorcode.LOGINNAME_PSW_IS_ERROR, errorcode.LOGINNAME_PSW_IS_ERROR_MSG, nil
	}
	logs.GetLogger().Info("Login Success")
	return errorcode.SUCCESS, "", userInfo
}

//用户注册
//Param ChatId Role DeviceType LoginName UserName UserIcon Password
//return code userId errmsg
func (userProcess *UserProcess) SignUp(ChatId, Role, DeviceType int32, LoginName, UserName, UserIcon, Password string) (int32, int32, string) {
	logs.GetLogger().Info("SignUp begins")
	//userDBManager 接口对象
	userDBManager := &db.ClassRoomDbMemberManager{}

	userId, err := userDBManager.SignUp(ChatId, Role, DeviceType, LoginName, UserName, UserIcon, Password)
	if err != nil {
		logs.GetLogger().Error("SignUp dbm err", err)
		return errorcode.USER_SIGN_UP__IS_ERR_IN_DATABASE, ZERO, err.Error()
	}
	id := int32(userId)
	logs.GetLogger().Info("Login Success")
	return errorcode.SUCCESS, id, ""
}

//修改用户信息
//Param UserID Role DeviceType  UserName UserIcon 根据内存获取实时参数
//return code Success err
func (userProcess *UserProcess) UserInfoUpdate(UserId int32) (int32, string) {
	logs.GetLogger().Info("UserInfoUpdate begins")
	//userDBManager 接口对象
	userDBManager := &db.ClassRoomDbMemberManager{}
	//Memory接口对象
	userMemoryManager := memory.GetUserInfoMemoryManager()

	//内存操作已经完成 当前的userInfo即为处理客户端请求后的usernfo
	user := userMemoryManager.GetUserInfo(UserId)
	code, err := userDBManager.UserInfoUpdate(UserId, user.Role, user.DeviceType, user.UserName, user.UserIcon)
	if err != nil {
		logs.GetLogger().Error("UserInfoUpdate dbm err", err)
		return errorcode.USERINFO_UPDATE_IS_ERR_IN_DATABASE, err.Error()
	}
	if 0 == code {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("UserInfoUpdate Success")
	return errorcode.SUCCESS, ""
}

//修改用户密码
//Param UserID oldPassword newPassword
//return code Success err
func (userProcess *UserProcess) UserPwdUpdate(UserId int32, newPassword, oldPassword string) (int32, string) {
	logs.GetLogger().Info("UserPwdUpdate In DATABase begins")
	//userDBManager 接口对象
	userDBManager := &db.ClassRoomDbMemberManager{}

	code, err := userDBManager.UserPwdUpdate(UserId, newPassword, oldPassword)
	if err != nil {
		logs.GetLogger().Error("UserPwdUpdate dbm err", err)
		return errorcode.USERPWD_UPDATE_IS_ERR_IN_DATABASE, err.Error()
	}
	if 0 == code {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("UserPwdUpdate In Datase Success")
	return errorcode.SUCCESS, ""
}

//重置用户密码
//Param UserID  newPassword
//return code Success err
func (userProcess *UserProcess) UserPwdReset(UserId int32, password string) (int32, string) {
	logs.GetLogger().Info("UserPwdReset In DATABase begins")
	//userDBManager 接口对象
	userDBManager := &db.ClassRoomDbMemberManager{}

	code, err := userDBManager.UserPwdReset(UserId, password)
	if err != nil {
		logs.GetLogger().Error("UserPwdReset dbm err", err)
		return errorcode.USERPWD_UPDATE_IS_ERR_IN_DATABASE, err.Error()
	}
	if 0 == code {
		return errorcode.DATABASE_ROW_AFFECT_NULL, "数据库删除操作没有行受影响"
	}
	logs.GetLogger().Info("UserPwdReset In Datase Success")
	return errorcode.SUCCESS, ""
}

//模糊查询用户信息 用户昵称
//返回[]userInfo,error
//return []userInfo
func (userProcess *UserProcess) FuzzySearchUserInfo(userName string) []info.ClassRoomMemberDbInfo {
	logs.GetLogger().Info("FuzzySearchUserInfo In DATABase begins")
	//返回的userinfo 切片
	usersInfo := []info.ClassRoomMemberDbInfo{}
	//userDBManager 接口对象
	userDBManager := &db.ClassRoomDbMemberManager{}

	users, err := userDBManager.FuzzySearchUserInfo(userName)
	if err != nil {
		logs.GetLogger().Error("FuzzySearchUserInfo dbm err", err)
		return nil
	}
	for _, userInfo := range users {
		if userInfo != nil {
			usersInfo = append(usersInfo, *userInfo)
		}
	}

	logs.GetLogger().Info("FuzzySearchUserInfo In Datase Success")
	return usersInfo
}
