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

func SignUp(response http.ResponseWriter, request *bean.UserSignUpRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("SignUp begins")
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
	//用户昵称为空
	if request.UserName == "" {
		writeErrMsg(errorcode.USER_NAME_CAN_NOT_NULL, errorcode.USER_NAME_CAN_NOT_NULL_MSG, response)
		return
	}
	//角色不能为空
	if request.Role == Role_IS_NULL {
		writeErrMsg(errorcode.USER_ROLE_CAN_NOT_NULL, errorcode.USER_ROLE_CAN_NOT_NULL_MSG, response)

	}

	//deviceType
	if request.DeviceType != DeviceType_Android && request.DeviceType != DeviceType_IPhone && request.DeviceType != DeviceType_PC {
		writeErrMsg(errorcode.DEVICE_TYPE_IS_ERR, errorcode.DEVICE_TYPE_IS_ERR_MSG, response)
		return

	}
	//角色编码正确 身份: 1=VIP学员 2=普通学员 3=试听学员 4=游客 6老师
	if Role_NORMAL != request.Role && Role_Teacher != request.Role && Role_TYR != request.Role && Role_VIP != request.Role && Role_VISITOR != request.Role {
		writeErrMsg(errorcode.USER_ROLE_IS_ERR, errorcode.USER_ROLE_IS_ERR_MSG, response)
		return

	} else {
		//Process接口对象
		userProcessManager := new(process.UserProcess)
		//memory接口对象
		userMemoryManager := memory.GetUserInfoMemoryManager()
		//返回参数 结构体(Code ErrMsg)
		result := new(bean.UserSignUpResponse) //反馈到客户端

		//IM 创建用户
		signUpResponse_IM := UserSignUp_IM(request, response)
		if signUpResponse_IM.Result == imbean.SUCCESS { //IM 创建用户成功
			logs.GetLogger().Info("IM 创建用户成功")
			//数据库中创建用户
			//Param ChatId Role DeviceType LoginName UserName UserIcon Password
			result.Code, result.UserId, result.ErrMsg = userProcessManager.SignUp(signUpResponse_IM.UserId, request.Role, request.DeviceType, request.LoginName, request.UserName, request.UserIcon, request.Password)
			if result.Code == errorcode.SUCCESS {
				//内存中创建教室
				//Param classRoomID(数据库创建成功后的)
				result.Code, result.UserId, result.ErrMsg = userMemoryManager.SignUp(result.UserId, signUpResponse_IM.UserId, request.Role, request.DeviceType, request.LoginName, request.UserName, request.UserIcon, request.Password)
				if result.Code == errorcode.SUCCESS {
					//创建用户成功 读取用户信息
					if result.UserId != 0 {
						userInfo := userMemoryManager.GetUserInfo(result.UserId)
						result.ChatId = userInfo.ChatId
						result.Role = userInfo.Role
						result.UserIcon = userInfo.UserIcon
						result.UserName = userInfo.UserName
					}
				}
			}

			//返回参数转化成Json数据 返回到客户端
			datas, errThree := json.Marshal(result)
			if errThree != nil {
				logs.GetLogger().Error(errThree.Error())
			}
			logs.GetLogger().Info("SignUp result:" + string(datas))
			fmt.Fprintln(response, string(datas))
			logs.GetLogger().Info("SignUp end")
		} else {
			writeErrMsg(signUpResponse_IM.Result, signUpResponse_IM.ErrorMessage, response)
			return
		}
	}

	logs.GetLogger().Info("=============================================================\n")

}

//在IM服务器创建用户
func UserSignUp_IM(request *bean.UserSignUpRequest, response http.ResponseWriter) *imbean.SignUpResponse {
	logs.GetLogger().Info("UserSignUp_IM Is begin")

	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct
	userSignUpRequest_IM := imbean.GetSignUpRequest()
	userSignUpRequest_IM.RequestServerId = 1
	userSignUpRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	userSignUpRequest_IM.RequestTime = timeutils.GetTimeStamp()
	userSignUpRequest_IM.SkipDBOperat = false
	userSignUpRequest_IM.MarkId = imbean.PUSH

	userSignUpRequest_IM.LoginName = request.LoginName
	userSignUpRequest_IM.UserName = request.UserName
	userSignUpRequest_IM.Password = request.Password
	userSignUpRequest_IM.UserIcon = request.UserIcon
	userSignUpRequest_IM.IsUpdateUser = true //如果登录名已存在，是否要覆盖已注册用户
	logs.GetLogger().Info("userSignUpRequest_IM data is:", userSignUpRequest_IM)
	//json
	userRequest_Bytes, err := json.Marshal(userSignUpRequest_IM)
	if err != nil {
		logs.GetLogger().Info("userSignUpRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("userSignUpRequest_IM to json is:", string(userRequest_Bytes))
	//response
	userSignUpResponse_IM := new(imbean.SignUpResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, userRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("UserSignUp_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.USER_SIGNUP_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &userSignUpResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error("userSignUpResponse_IM to json is err:", errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("UserSignUp_IM Result IS：", userSignUpResponse_IM)
	logs.GetLogger().Info("UserSignUp_IM Result IS：", string(data))

	logs.GetLogger().Info("UserSignUp_IM Is end")

	return userSignUpResponse_IM
}
