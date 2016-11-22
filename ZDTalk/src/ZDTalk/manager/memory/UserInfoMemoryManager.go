package memory

import (
	"ZDTalk/errorcode"
	"ZDTalk/manager/db"
	dbInfo "ZDTalk/manager/db/info"
	logs "ZDTalk/utils/log4go"
	"sync"
)

//Userinfo 信息
type UserInfo struct {
	UserId        int32
	ChatId        int32 //用户对应的IM中的userId
	Role          int32
	DeviceType    int32 //设备登录类型
	LoginName     string
	UserName      string
	UserIcon      string
	Password      string
	YYToken       string //YY提供的Token
	ClassRoomList []int32
}

//HandUserinfo 信息
type HandUserInfo struct {
	*UserInfo
	HandTime int64 //举手时间
}

//用户信息内存管理结构体
type UserInfoMemoryManager struct {
	Users map[int32]*UserInfo
	Lock  sync.Mutex
}

//全局内存变量
//对应数据库表 userInfo
var userInfoMemoryManager *UserInfoMemoryManager

//初始化全局内存变量UserInfoMemoryManager 即读取将数据表userInfo到内存
func GetUserInfoMemoryManager() *UserInfoMemoryManager {
	if userInfoMemoryManager == nil {
		logs.Logs.Info("------------- ZDTalk userInfo Memory Initial begin-------------")

		//初始化userInfoManager
		userInfoMemoryManager = &UserInfoMemoryManager{}
		userInfoMemoryManager.Users = make(map[int32]*UserInfo)
		//加载数据库中userInfo数据表到内存 客户端的请求操作直接去内存读取数据 即声明的全局变量userInfoManager
		manager := new(db.ClassRoomDbMemberManager)
		result, err := manager.LoadAllUser()

		if err != nil {
			logs.GetLogger().Error("LoadAllUser Error:" + err.Error())
			return nil
		}

		//初始化map ClassRoomMemberDbInfo 赋值数据库返回的结果
		dbUserInfo := make(map[int32]*dbInfo.ClassRoomMemberDbInfo)
		dbUserInfo = result

		for _, userInfo := range dbUserInfo {
			//根据dbUserInfo初始化UserInfo
			user := &UserInfo{userInfo.UserId, userInfo.ChatId, userInfo.Role, userInfo.DeviceType, userInfo.LoginName, userInfo.UserName, userInfo.UserIcon, userInfo.Password, userInfo.YYToken, userInfo.ClassRoomList}
			//设置userInfoManager的ClassRooms属性
			userInfoMemoryManager.Users[user.UserId] = user
		}
		logs.Logs.Info("ZDTalk userInfo Memory 用户数量 %d ", len(dbUserInfo))
		logs.Logs.Info("------------- ZDTalk userInfo Memory Initial end-------------")
	}
	return userInfoMemoryManager
}

//获取用户的imid
func (self *UserInfoMemoryManager) GetUserIMId(userId int32) int32 {

	//获取教室id对应的ClassRoomIMid
	if r, ok := self.Users[userId]; ok {
		return r.ChatId
	} else {
		return 0
	}
}

//通过用户的IMUserId获取UserId(ZD)
func (self *UserInfoMemoryManager) GetUserId(userImId int32) int32 {
	//memory接口对象
	var userId int32 = 0
	for _, userInfo := range self.Users {
		if userInfo.ChatId == userImId {
			userId = userInfo.UserId
			break
		} else {
			continue
		}
	}
	return userId
}

//获取全部用户
func (self *UserInfoMemoryManager) GetAllUsers() []*UserInfo {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	data := []*UserInfo{}
	for _, v := range self.Users {
		data = append(data, v)
	}
	return data
}

//获取用户信息
func (self *UserInfoMemoryManager) GetUserInfo(userId int32) *UserInfo {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if u, ok := self.Users[userId]; ok {
		return u
	}
	return nil
}

//获取举手列表 用户信息列表
func (self *UserInfoMemoryManager) GetHandsUserList(classRoomId int32) []HandUserInfo {
	logs.GetLogger().Info("GetHandsUserList In Memory Beginning")

	self.Lock.Lock()
	classRoomMemoryManager.Lock.Lock()
	defer self.Lock.Unlock()
	defer classRoomMemoryManager.Lock.Unlock()

	//初始化用户信息列表 切片
	userInfoSlice := []HandUserInfo{}
	//根据教室编号 获取该教室的举手列表
	if r, ok := classRoomMemoryManager.ClassRooms[classRoomId]; ok {
		handMemberList := r.HandMemberList
		logs.GetLogger().Info("HandMemberList Is:", handMemberList)
		//根据举手列表 获取各个成员的个人信息
		for _, handMember := range handMemberList {
			//根据userID从全局内存变量UserInfoMemoryManager获取对应的用户信息
			if u, ok := self.Users[handMember.UserId]; ok {
				handUserInfo := HandUserInfo{u, handMember.HandsTime}
				userInfoSlice = append(userInfoSlice, handUserInfo)
			} else {
				logs.GetLogger().Info("UserID: not exit", handMember.UserId)
			}
		}
	} else {
		return nil
	}
	logs.GetLogger().Info("GetHandsUserList In Memory Success")

	return userInfoSlice
}

//获取在线用户列表 用户信息列表
func (self *UserInfoMemoryManager) GetOnlineUserList(classRoomId int32) []*UserInfo {
	logs.GetLogger().Info("GetOnlineUserList In Memory begins")

	self.Lock.Lock()
	defer self.Lock.Unlock()

	//初始化用户信息列表 切片
	userInfoSlice := []*UserInfo{}
	//根据教室编号 获取该教室的在线用户列表
	if r, ok := classRoomMemoryManager.ClassRooms[classRoomId]; ok {
		onLineMemberList := r.OnLineMemberList
		//根据在线用户列表 获取各个成员的个人信息
		for _, userID := range onLineMemberList {
			//根据userID从全局内存变量UserInfoMemoryManager获取对应的用户信息
			if u, ok := self.Users[userID]; ok {
				userInfoSlice = append(userInfoSlice, u)
			} else {
				logs.GetLogger().Info("UserID: not exit", userID)
			}

		}
	} else {
		return nil
	}
	logs.GetLogger().Info("GetOnlineUserList In Memory Success")

	return userInfoSlice
}

//登录
func (self *UserInfoMemoryManager) Login(loginName, password string) (int32, string, *UserInfo) {
	logs.GetLogger().Info("login In Memory begins")

	self.Lock.Lock()
	defer self.Lock.Unlock()

	userList := self.Users
	//logs.GetLogger().Info("users is:", userList)
	for _, user := range userList {
		if user.LoginName == loginName { //查询到用户名
			if user.Password == password { //密码相等 登录成功
				logs.GetLogger().Info("login In Memory Success")
				return errorcode.SUCCESS, "", user
			} else { //登录失败
				return errorcode.LOGINNAME_PSW_IS_ERROR, errorcode.LOGINNAME_PSW_IS_ERROR_MSG, nil

			}
		}
	}
	return errorcode.LOGINNAME_IS_NOT_EXIST, errorcode.LOGINNAME_IS_NOT_EXIST_MSG, nil
}

//用户注册
//Param UserId ChatId Role DeviceType LoginName UserName UserIcon Password
//return code userId errMsg
func (self *UserInfoMemoryManager) SignUp(UserId, ChatId, Role, DeviceType int32, LoginName, UserName, UserIcon, Password string) (int32, int32, string) {
	logs.GetLogger().Info("SignUp In Memory begins")
	self.Lock.Lock()
	defer self.Lock.Unlock()

	userInfo := new(UserInfo)
	userInfo.UserId = UserId
	userInfo.ChatId = ChatId
	userInfo.Role = Role
	userInfo.UserIcon = UserIcon
	userInfo.LoginName = LoginName
	userInfo.UserName = UserName
	userInfo.DeviceType = DeviceType
	userInfo.Password = Password

	self.Users[UserId] = userInfo

	logs.GetLogger().Info("SignUp In Memory Success")

	return errorcode.SUCCESS, UserId, ""
}

//修改用户信息
//Param UserID Role DeviceType UserName UserIcon
//return  code errMsg
func (self *UserInfoMemoryManager) UserInfoUpdate(UserId, Role, DeviceType int32, UserName, UserIcon string) (int32, string) {
	logs.GetLogger().Info("Update userInfo In Memory begins")

	self.Lock.Lock()
	defer self.Lock.Unlock()
	//传引用已修改原值
	if r, ok := self.Users[UserId]; ok {
		if Role != 0 {
			r.Role = Role
		}
		if DeviceType != 0 {
			r.DeviceType = DeviceType
		}
		if UserIcon != "" {
			r.UserIcon = UserIcon
		}
		if UserName != "" {
			r.UserName = UserName
		}
	} else {
		logs.GetLogger().Error("用户不存在")
		return errorcode.USER_ID_IS_NOT_EXIT, errorcode.USER_ID_IS_NOT_EXIT_MSG
	}
	logs.GetLogger().Info("Update userInfo  In Memory Success:")
	return errorcode.SUCCESS, ""
}

//修改用户密码
//Param UserID newPassword oldPassword
//return code Success errMsg
func (self *UserInfoMemoryManager) UserPwdUpdate(UserId int32, oldPassword, newPassword string) (int32, string) {
	logs.GetLogger().Info("Update userPwd In Memory begins")

	self.Lock.Lock()
	defer self.Lock.Unlock()
	//传引用已修改原值
	if r, ok := self.Users[UserId]; ok {
		if r.Password == oldPassword {
			r.Password = newPassword
		} else {
			logs.GetLogger().Error("用户密码错误")
			return errorcode.USER_PWD_IS_ERROR, errorcode.USER_PWD_IS_ERROR_MSG

		}
	} else {
		logs.GetLogger().Error("用户不存在")
		return errorcode.USER_ID_IS_NOT_EXIT, errorcode.USER_ID_IS_NOT_EXIT_MSG
	}
	logs.GetLogger().Info("Update userPwd  In Memory Success:")
	return errorcode.SUCCESS, ""
}

//重置用户密码
//Param UserID newPassword
//return code Success errMsg
func (self *UserInfoMemoryManager) UserPwdReset(UserId int32, newPassword string) (int32, string) {
	logs.GetLogger().Info("Reset userPwd In Memory begins")

	self.Lock.Lock()
	defer self.Lock.Unlock()
	//传引用已修改原值
	if r, ok := self.Users[UserId]; ok {
		r.Password = newPassword

	} else {
		logs.GetLogger().Error("用户不存在")
		return errorcode.USER_ID_IS_NOT_EXIT, errorcode.USER_ID_IS_NOT_EXIT_MSG
	}
	logs.GetLogger().Info("Update userPwd  In Memory Success:")
	return errorcode.SUCCESS, ""
}
