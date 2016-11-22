package user

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/ZDTalk_http/httpprocess/process"
	"ZDTalk/errorcode"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUserInfoFuzzy(response http.ResponseWriter, request *bean.GetUserInfoFuzzyRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("Login begins")
	//用户昵称为空
	if request.UserName == "" {
		writeErrMsg(errorcode.USER_NAME_CAN_NOT_NULL, errorcode.USER_NAME_CAN_NOT_NULL_MSG, response)
		return
	}

	//Process接口对象
	userProcessManager := new(process.UserProcess)
	//返回参数 结构体(Code ErrMsg)
	result := new(bean.GetUserInfoFuzzyResponse) //反馈到客户端

	users := userProcessManager.FuzzySearchUserInfo(request.UserName)

	if users == nil {
		result.Code = errorcode.FUZZY_SEARCH_USERINFO_RESULT_IS_NULL
		result.ErrMsg = errorcode.FUZZY_SEARCH_USERINFO_RESULT_IS_NULL_MSG
	} else {
		result.Code = errorcode.SUCCESS
		for _, u := range users {
			userInfo := bean.UserInfo{}
			userInfo.ChatId = u.ChatId
			userInfo.Role = u.Role
			userInfo.UserIcon = u.UserIcon
			userInfo.UserId = u.UserId
			userInfo.UserName = u.UserName

			result.Data = append(result.Data, userInfo)
		}
	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("Login result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("Login end")
	logs.GetLogger().Info("=============================================================\n")

}
