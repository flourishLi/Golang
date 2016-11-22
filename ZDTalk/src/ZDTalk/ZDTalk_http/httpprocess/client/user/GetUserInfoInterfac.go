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

func GetUserInfo(response http.ResponseWriter, request *bean.GetUserInfoRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("GetUserInfo begains")
	//user为空
	if request.UserId == 0 {
		writeErrMsg(errorcode.USER_ID_CAN_NOT_NULL, errorcode.USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}
	//客户端反馈结果
	result := new(bean.GetUserInfoResponse)
	//memory接口对象
	userInfoMemoryManager := memory.GetUserInfoMemoryManager()

	//获取用户信息
	userInfo := userInfoMemoryManager.GetUserInfo(request.UserId)

	if userInfo == nil {
		result.Code = errorcode.USER_ID_IS_NOT_EXIT
		result.ErrMsg = errorcode.USER_ID_IS_NOT_EXIT_MSG
	} else {
		result.Data.UserId = userInfo.UserId
		result.Data.ChatId = userInfo.ChatId
		result.Data.Role = userInfo.Role
		result.Data.UserName = userInfo.UserName
		result.Data.UserIcon = userInfo.UserIcon

		result.Code = errorcode.SUCCESS
	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("GetUserInfo result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("GetUserInfo end")
	logs.GetLogger().Info("=============================================================\n")

}
