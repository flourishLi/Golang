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

func UpdateUserPwd(response http.ResponseWriter, request *bean.UserPwdUpdateRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("UpdateUserPwd begins")
	//用户编号为空
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.USER_ID_CAN_NOT_NULL, errorcode.USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}
	//oldPassword
	if request.OldPassword == "" {
		writeErrMsg(errorcode.OLD_PWD_CAN_NOT_NULL, errorcode.OLD_PWD_CAN_NOT_NULL_MSG, response)
		return
	}
	//newpassword
	if request.NewPassword == "" {
		writeErrMsg(errorcode.NEW_PWD_CAN_NOT_NULL, errorcode.NEW_PWD_CAN_NOT_NULL_MSG, response)
		return
	}
	//IM 修改用户信息
	updateUserPwdResponse_IM := UpdateUserPwd_IM(request, response)
	if updateUserPwdResponse_IM.Result == imbean.SUCCESS { //IM 成功
		//内存中修改用户
		//Process接口对象
		userProcessManager := new(process.UserProcess)
		//memory接口对象
		userMemoryManager := memory.GetUserInfoMemoryManager()
		//返回参数 结构体(Code ErrMsg)
		result := new(bean.UserPwdUpdateResponse) //反馈到客户端

		//内存中修改用户密码
		result.Code, result.ErrMsg = userMemoryManager.UserPwdUpdate(request.RequestUserId, request.OldPassword, request.NewPassword)
		if result.Code == errorcode.SUCCESS {
			//数据库中修改密码
			result.Code, result.ErrMsg = userProcessManager.UserPwdUpdate(request.RequestUserId, request.NewPassword, request.OldPassword)
		}
		datas, err1 := json.Marshal(result)
		if err1 != nil {
			logs.GetLogger().Error(err1.Error())
		}
		logs.GetLogger().Info("UpdateUserPwd result:" + string(datas))
		fmt.Fprintln(response, string(datas))

		logs.GetLogger().Info("UpdateUserPwd end")
	} else { //IM 失败
		writeErrMsg(updateUserPwdResponse_IM.Result, updateUserPwdResponse_IM.ErrorMessage, response)
		return
	}

	logs.GetLogger().Info("=============================================================\n")

}

//在IM服务器修改用户密码
func UpdateUserPwd_IM(request *bean.UserPwdUpdateRequest, response http.ResponseWriter) *imbean.UpdateUserPwdResponse {
	logs.GetLogger().Info("UpdateUserPwd_IM Is begin")
	//获取用户id对应的IMid
	userIMId := GetUserIMId(request.RequestUserId)
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct
	userPwdUpdateRequest_IM := imbean.GetUpdateUserPwdRequest()
	userPwdUpdateRequest_IM.RequestServerId = 1
	userPwdUpdateRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	userPwdUpdateRequest_IM.RequestTime = timeutils.GetTimeStamp()
	userPwdUpdateRequest_IM.SkipDBOperat = false
	userPwdUpdateRequest_IM.MarkId = imbean.PUSH

	userPwdUpdateRequest_IM.UserId = userIMId
	userPwdUpdateRequest_IM.NewPassword = cryptoutils.Md5Encode(request.NewPassword)
	userPwdUpdateRequest_IM.OldPassword = cryptoutils.Md5Encode(request.OldPassword)

	logs.GetLogger().Info("userPwdUpdateRequest_IM data is:", userPwdUpdateRequest_IM)
	//json
	userRequest_Bytes, err := json.Marshal(userPwdUpdateRequest_IM)
	if err != nil {
		logs.GetLogger().Info("userPwdUpdateRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("userPwdUpdateRequest_IM to json is:", string(userRequest_Bytes))
	//response
	userPwdUpdateResponse_IM := new(imbean.UpdateUserPwdResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, userRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("UpdateUserPwd_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.USER_SIGNUP_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &userPwdUpdateResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error("userPwdUpdateResponse_IM to json is err:", errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("UpdateUserPwd_IM Result IS：", userPwdUpdateResponse_IM)
	logs.GetLogger().Info("UpdateUserPwd_IM Result IS：", string(data))

	logs.GetLogger().Info("UpdateUserPwd_IM Is end")

	return userPwdUpdateResponse_IM
}
