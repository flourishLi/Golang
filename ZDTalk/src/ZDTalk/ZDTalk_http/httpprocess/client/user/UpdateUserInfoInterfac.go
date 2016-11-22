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

func UpdateUserInfo(response http.ResponseWriter, request *bean.UserInfoUpdateRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("Login begins")
	//用户编号为空
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.USER_ID_CAN_NOT_NULL, errorcode.USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}

	//role
	if request.Role != Role_NORMAL && request.Role != Role_Teacher && request.Role != Role_TYR && request.Role != Role_VIP && request.Role != Role_VISITOR {
		writeErrMsg(errorcode.USER_ROLE_IS_ERR, errorcode.USER_ROLE_IS_ERR_MSG, response)
		return
	}

	//deviceType
	if request.DeviceType != DeviceType_Android && request.DeviceType != DeviceType_IPhone && request.DeviceType != DeviceType_PC {
		writeErrMsg(errorcode.DEVICE_TYPE_IS_ERR, errorcode.DEVICE_TYPE_IS_ERR_MSG, response)
		return

	}
	//Process接口对象
	userProcessManager := new(process.UserProcess)
	//memory接口对象
	userMemoryManager := memory.GetUserInfoMemoryManager()
	//返回参数 结构体(Code ErrMsg)
	result := new(bean.UserInfoUpdateResponse) //反馈到客户端

	//IM 修改用户信息
	updateUserinfoResponse_IM := UpdateUserInfo_IM(request, response)
	if updateUserinfoResponse_IM.Result == imbean.SUCCESS { //IM 成功
		//内存中修改用户
		//Param UserID Role DeviceType LoginName UserName UserIcon
		result.Code, result.ErrMsg = userMemoryManager.UserInfoUpdate(request.RequestUserId, request.Role, request.DeviceType, request.UserName, request.UserIcon)
		if result.Code == errorcode.SUCCESS {
			//数据库中修改用户
			result.Code, result.ErrMsg = userProcessManager.UserInfoUpdate(request.RequestUserId)
		}
		datas, err1 := json.Marshal(result)
		if err1 != nil {
			logs.GetLogger().Error(err1.Error())
		}
		logs.GetLogger().Info("Login result:" + string(datas))
		fmt.Fprintln(response, string(datas))

		logs.GetLogger().Info("Login end")
	} else {
		writeErrMsg(updateUserinfoResponse_IM.Result, updateUserinfoResponse_IM.ErrorMessage, response)
		return
	}
	logs.GetLogger().Info("=============================================================\n")

}

//在IM服务器修改用户信息
func UpdateUserInfo_IM(request *bean.UserInfoUpdateRequest, response http.ResponseWriter) *imbean.UpdateUserInfoResponse {
	logs.GetLogger().Info("UpdateUserInfo_IM Is begin")
	//获取用户id对应的IMid
	userIMId := GetUserIMId(request.RequestUserId)
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct
	userInfoUpdateRequest_IM := imbean.GetUpdateUserInfoRequest()
	userInfoUpdateRequest_IM.RequestServerId = 1
	userInfoUpdateRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	userInfoUpdateRequest_IM.RequestTime = timeutils.GetTimeStamp()
	userInfoUpdateRequest_IM.SkipDBOperat = false
	userInfoUpdateRequest_IM.MarkId = imbean.PUSH

	userInfoUpdateRequest_IM.UserId = userIMId
	userInfoUpdateRequest_IM.UserName = request.UserName
	userInfoUpdateRequest_IM.UserIcon = request.UserIcon

	logs.GetLogger().Info("userInfoUpdateRequest_IM data is:", userInfoUpdateRequest_IM)
	//json
	userRequest_Bytes, err := json.Marshal(userInfoUpdateRequest_IM)
	if err != nil {
		logs.GetLogger().Info("userInfoUpdateRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("userInfoUpdateRequest_IM to json is:", string(userRequest_Bytes))
	//response
	userInfoUpdateResponse_IM := new(imbean.UpdateUserInfoResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, userRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("UpdateUserInfo_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.USER_SIGNUP_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &userInfoUpdateResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error("userInfoUpdateResponse_IM to json is err:", errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("UpdateUserInfo_IM Result IS：", userInfoUpdateResponse_IM)
	logs.GetLogger().Info("UpdateUserInfo_IM Result IS：", string(data))

	logs.GetLogger().Info("UpdateUserInfo_IM Is end")

	return userInfoUpdateResponse_IM
}

//获取用户的imid
func GetUserIMId(userId int32) int32 {
	//memory接口对象
	userInfoMemoryManager := memory.GetUserInfoMemoryManager()

	//获取教室id对应的ClassRoomIMid
	if r, ok := userInfoMemoryManager.Users[userId]; ok {
		return r.ChatId
	} else {
		return 0
	}
}
