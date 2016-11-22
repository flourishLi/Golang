package classroom

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetClassRoomInfo(response http.ResponseWriter, request *bean.GetClassRoomInfoRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("GetClassRoomInfo begains")
	//ClassRoomId为空
	if request.ClassRoomId == ClassRoomID_IS_NULL {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}
	//客户端反馈结果
	result := new(bean.GetClassRoomInfoResponse)
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//获取指定的教室信息
	d := classRoomMemoryManager.GetClassRoom(request.ClassRoomId)

	if d == nil {
		result.Code = errorcode.CLASS_ROOM_IS_NOT_EXIST
		result.ErrMsg = errorcode.CLASS_ROOM_IS_EXIST_MSG
	} else {
		result.Code = errorcode.SUCCESS
		result.ClassRoomLogo = d.ClassLogo
		result.ClassRoomName = d.ClassName
		result.Description = d.Description
		result.ClassRoomId = d.ClassId
		result.ClassRoomIMId = d.ClassRoomIMId
		result.ClassRoomCourse = d.ClassRoomCourse
		result.ClassRoomStatus = d.ClassRoomStatus
		result.SettingStatus = d.SettingStatus
		result.CreateTime = d.CreateTime
		result.ForbidHandStatus = d.ForbidHandStatus
		result.CreatorUserId = d.CreatorUserId
		users := make([]bean.UserInfo, 1)
		//memberList
		users = GetUsersInfo(d.MemberList)
		result.MemberList = users
		//OnLineMemberList
		users = GetUsersInfo(d.OnLineMemberList)
		result.OnLineMemberList = users
		//ForbidSayMemberList
		users = GetUsersInfo(d.ForbidSayMemberList)
		result.ForbidSayMemberList = users

		//SayingMemberList
		users = GetUsersInfo(d.SayingMemberList)
		result.SayingMemberList = users
		//HandMemberList
		handList := []int32{}
		for userId, _ := range d.HandMemberList {
			user := GetUserInfo(userId)
			if user.UserId != 0 {
				handList = append(handList, userId)
			}
		}
		users = GetUsersInfo(handList)
		result.HandMemberList = users
	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("GetClassRoomInfo result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("GetClassRoomInfo end")
	logs.GetLogger().Info("=============================================================\n")
}

func GetUserInfo(userId int32) bean.UserInfo {
	//memory接口对象
	userInfoMemoryManager := memory.GetUserInfoMemoryManager()
	userInfo := bean.UserInfo{}
	userinfoMemory := userInfoMemoryManager.GetUserInfo(userId)
	if userinfoMemory != nil { //存在
		userInfo.UserId = userinfoMemory.UserId
		userInfo.ChatId = userinfoMemory.ChatId
		userInfo.UserName = userinfoMemory.UserName
		userInfo.UserIcon = userinfoMemory.UserIcon
		userInfo.Role = userinfoMemory.Role
	} else { //用户表中不存在该用户
		logs.GetLogger().Info("User ID is Not exit:", userId)
	}
	return userInfo
}

func GetUsersInfo(userIdList []int32) []bean.UserInfo {
	//memory接口对象
	userInfoMemoryManager := memory.GetUserInfoMemoryManager()
	users := make([]bean.UserInfo, 0)
	for _, userId := range userIdList {
		userInfoMemory := userInfoMemoryManager.GetUserInfo(userId)
		if userInfoMemory != nil { //存在
			userInfo := bean.UserInfo{}
			userInfo.UserId = userInfoMemory.UserId
			userInfo.ChatId = userInfoMemory.ChatId
			userInfo.UserName = userInfoMemory.UserName
			userInfo.UserIcon = userInfoMemory.UserIcon
			userInfo.Role = userInfoMemory.Role
			users = append(users, userInfo)
		} else { //用户表中不存在该用户
			logs.GetLogger().Info("User ID is Not exit:", userId)
		}
	}
	return users
}
