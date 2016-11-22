package user

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/ZDTalk_http/httpprocess/process"
	"ZDTalk/ZDTalk_http/imbean"
	"ZDTalk/config"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	"ZDTalk/utils/cryptoutils"
	"ZDTalk/utils/httputil"
	logs "ZDTalk/utils/log4go"
	"ZDTalk/utils/timeutils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ResetUserPwd(response http.ResponseWriter, request *bean.UserPwdResetRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("ResetUserPwd begins")
	//用户编号为空
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.USER_ID_CAN_NOT_NULL, errorcode.USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}
	//密码
	if request.NewPassword == "" {
		writeErrMsg(errorcode.NEW_PWD_CAN_NOT_NULL, errorcode.NEW_PWD_CAN_NOT_NULL_MSG, response)
		return
	}
	resetUserPwdResponse_IM := ResetUserPwd_IM(request, response)
	if resetUserPwdResponse_IM.Result == imbean.SUCCESS { //IM 成功
		//Process接口对象
		userProcessManager := new(process.UserProcess)
		//memory接口对象
		userMemoryManager := memory.GetUserInfoMemoryManager()
		//返回参数 结构体(Code ErrMsg)
		result := new(bean.UserPwdUpdateResponse) //反馈到客户端

		//内存中修改用户密码
		result.Code, result.ErrMsg = userMemoryManager.UserPwdReset(request.RequestUserId, request.NewPassword)
		if result.Code == errorcode.SUCCESS {
			//数据库中修改密码
			result.Code, result.ErrMsg = userProcessManager.UserPwdReset(request.RequestUserId, request.NewPassword)
		}
		datas, err1 := json.Marshal(result)
		if err1 != nil {
			logs.GetLogger().Error(err1.Error())
		}
		logs.GetLogger().Info("ResetUserPwd result:" + string(datas))
		fmt.Fprintln(response, string(datas))

		logs.GetLogger().Info("ResetUserPwd end")
	} else {
		writeErrMsg(resetUserPwdResponse_IM.Result, resetUserPwdResponse_IM.ErrorMessage, response)
		return
	}

	logs.GetLogger().Info("=============================================================\n")
}

//在IM服务器重置用户密码
func ResetUserPwd_IM(request *bean.UserPwdResetRequest, response http.ResponseWriter) *imbean.ResetUserPwdResponse {
	logs.GetLogger().Info("ResetUserPwd_IM Is begin")
	//获取用户id对应的IMid
	userIMId := GetUserIMId(request.RequestUserId)
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct
	userPwdResetRequest_IM := imbean.GetResetUserPwdRequest()
	userPwdResetRequest_IM.RequestServerId = 1
	userPwdResetRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	userPwdResetRequest_IM.RequestTime = timeutils.GetTimeStamp()
	userPwdResetRequest_IM.SkipDBOperat = false
	userPwdResetRequest_IM.MarkId = imbean.PUSH

	userPwdResetRequest_IM.UserId = userIMId
	userPwdResetRequest_IM.NewPassword = cryptoutils.Md5Encode(request.NewPassword)

	logs.GetLogger().Info("userPwdResetRequest_IM data is:", userPwdResetRequest_IM)
	//json
	userRequest_Bytes, err := json.Marshal(userPwdResetRequest_IM)
	if err != nil {
		logs.GetLogger().Info("userPwdResetRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("userPwdResetRequest_IM to json is:", string(userRequest_Bytes))
	//response
	userPwdResetResponse_IM := new(imbean.ResetUserPwdResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, userRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("ResetUserPwd_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.USER_SIGNUP_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &userPwdResetResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error("userPwdResetResponse_IM to json is err:", errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("ResetUserPwd_IM Result IS：", userPwdResetResponse_IM)
	logs.GetLogger().Info("ResetUserPwd_IM Result IS：", string(data))

	logs.GetLogger().Info("ResetUserPwd_IM Is end")

	return userPwdResetResponse_IM
}
