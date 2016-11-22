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

func Login(response http.ResponseWriter, request *bean.LoginRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("Login begins")
	//登录账号为空
	if request.LoginName == "" {
		writeErrMsg(errorcode.LOGIN_NAME_CAN_NOT_NULL, errorcode.LOGIN_NAME_CAN_NOT_NULL_MSG, response)
		return
	}
	//密码为空
	if request.Password == "" {
		writeErrMsg(errorcode.LOGIN_PASSWORD_CAN_NOT_NULL, errorcode.LOGIN_PASSWORD_CAN_NOT_NULL_MSG, response)
		return
	}

	//memory接口对象
	userMemoryManager := memory.GetUserInfoMemoryManager()

	//返回参数 结构体(Code ErrMsg)
	result := new(bean.LoginResponse) //反馈到客户端

	Code, ErrMsg, Data := userMemoryManager.Login(request.LoginName, request.Password)

	userInfo := &bean.UserInfo{}
	if Data != nil {

		userInfo.UserId = Data.UserId
		userInfo.ChatId = Data.ChatId
		userInfo.Role = Data.Role
		userInfo.UserIcon = Data.UserIcon
		userInfo.UserName = Data.UserName
	}
	result.Data = userInfo
	result.Code = Code
	result.ErrMsg = ErrMsg

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("Login result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("Login end")
	logs.GetLogger().Info("=============================================================\n")

}
