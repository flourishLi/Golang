package db

import (
	"ZDTalk/db/database"
	dbInfo "ZDTalk/manager/db/info"
	"ZDTalk/manager/db/mysqldb"
)

//函数的接收者和函数的定义必须在一个包内
type ClassRoomDbMemberManager struct {
}

//加载所有用户 用户列表
//Param 空
//return ClassRoomMemberDbInfo切片 error
func (manager ClassRoomDbMemberManager) LoadAllUser() (map[int32]*dbInfo.ClassRoomMemberDbInfo, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomMemberDbManager{}.LoadAllUser()
	}
	return nil, nil
}

//模糊查询用户信息 用户昵称
//返回[] *userInfo,error

func (manager ClassRoomDbMemberManager) FuzzySearchUserInfo(userName string) (map[int32]*dbInfo.ClassRoomMemberDbInfo, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomMemberDbManager{}.FuzzyQueryUserInfo(userName)
	}
	return nil, nil
}

//登录
////Param loginName Password
//return userinfo error
func (manager ClassRoomDbMemberManager) Login(loginName, password string) (*dbInfo.ClassRoomMemberDbInfo, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomMemberDbManager{}.Login(loginName, password)
	}
	return nil, nil
}

//用户注册
//Param ChatId Role DeviceType LoginName UserName UserIcon Password
//return userId err
func (manager ClassRoomDbMemberManager) SignUp(ChatId, Role, DeviceType int32, LoginName, UserName, UserIcon, Password string) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomMemberDbManager{}.SignUp(ChatId, Role, DeviceType, LoginName, UserName, UserIcon, Password)
	}
	return ZERO, nil
}

//修改用户信息
//Param UserID Role DeviceType UserName UserIcon
//return int 1 Success err
func (manager ClassRoomDbMemberManager) UserInfoUpdate(UserId, Role, DeviceType int32, UserName, UserIcon string) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomMemberDbManager{}.UserInfoUpdate(UserId, Role, DeviceType, UserName, UserIcon)
	}
	return ZERO, nil
}

//修改用户密码
//Param UserID newPassword oldPassword
//return int 1 Success err
func (manager ClassRoomDbMemberManager) UserPwdUpdate(UserId int32, newPassword, oldPassword string) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomMemberDbManager{}.UserPwdUpdate(UserId, newPassword, oldPassword)
	}
	return ZERO, nil
}

//重置用户密码
//Param UserID newPassword
//return int 1 Success err
func (manager ClassRoomDbMemberManager) UserPwdReset(UserId int32, newPassword string) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ClassRoomMemberDbManager{}.UserPwdReset(UserId, newPassword)
	}
	return ZERO, nil
}
