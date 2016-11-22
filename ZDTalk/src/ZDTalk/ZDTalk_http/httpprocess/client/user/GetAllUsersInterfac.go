package user

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAllUsers(response http.ResponseWriter, request *bean.GetAllUserRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("GetAllUsers begains")
	//客户端反馈结果
	result := new(bean.GetAllUserResponse)
	//memory接口对象
	userInfoMemoryManager := memory.GetUserInfoMemoryManager()

	//获取用户信息
	users := userInfoMemoryManager.GetAllUsers()
	for _, v := range users {

		userInfo := bean.UserInfo{}
		userInfo.UserId = v.UserId
		userInfo.ChatId = v.ChatId
		userInfo.Role = v.Role
		userInfo.UserName = v.UserName
		userInfo.UserIcon = v.UserIcon
		result.Data = append(result.Data, userInfo)
	}

	result.Code = errorcode.SUCCESS

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("GetAllUsers result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("GetAllUsers end")
	logs.GetLogger().Info("=============================================================\n")

}
