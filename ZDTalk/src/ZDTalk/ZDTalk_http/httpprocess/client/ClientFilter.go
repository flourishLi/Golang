package client

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/ZDTalk_http/bean/resource"

	"ZDTalk/ZDTalk_http/httpprocess/client/classroom"
	"ZDTalk/ZDTalk_http/httpprocess/client/user"
	//	"ZDTalk/ZDTalk_http/httpprocess/client/student"
	"ZDTalk/errorcode"
	logs "ZDTalk/utils/log4go"
	"ZDTalk/utils/stringutils"
	"encoding/json"
	"fmt"
	//	"fmt"
	//	studentbean "ZDTalk/ZDTalk_http/bean/App_Student"
	"io/ioutil"
	"net/http"
)

//拦截器
func Filter(response http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logs.GetLogger().Error(err)
			writeErrMsg(errorcode.REQUEST_FILTER_ERROR, errorcode.REQUEST_FILTER_ERROR_MSG, response)
		}
	}()

	if request.Method != "POST" {
		logs.GetLogger().Info(errorcode.REQUEST_METHOD_ERROR_MSG + ", Current Method is " + request.Method)
		writeErrMsg(errorcode.REQUEST_METHOD_ERROR, errorcode.REQUEST_METHOD_ERROR_MSG+", Current Method is "+request.Method, response)
		return
	}

	data, err := ioutil.ReadAll(request.Body)
	logs.GetLogger().Info("需要解析的对象json:", string(data))
	defer request.Body.Close()
	if err != nil {
		panic(err)
	}
	base := new(bean.ClientBaseRequest)
	err = json.Unmarshal(data, &base)
	if err != nil {
		logs.GetLogger().Error(err.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return
	}
	logs.GetLogger().Info("CMD IS:", base.Command)

	logs.GetLogger().Info("Unmarshal json to ClientBaseRequest is:", base)
	if stringutils.StringIsEmpty(base.Command) {
		logs.GetLogger().Error(errorcode.CMD_IS_NULL_ERROR_MSG)
		writeErrMsg(errorcode.CMD_IS_NULL, errorcode.CMD_IS_NULL_ERROR_MSG, response)
		return
	}
	switch base.Command {
	// 创建教室
	case "CREATE_CLASS_ROOM":
		request := bean.CreateClassroomRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to CreateClassroomRequest is:", request)
		classroom.CreateClassroom(response, &request)
		// 删除教室
	case "DELETE_CLASS_ROOM":
		request := bean.DeleteClassroomRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to DeleteClassroomRequest is:", request)
		classroom.DeleteClassroom(response, &request)
		// 修改教室
	case "UPDATE_CLASS_ROOM":
		request := bean.UpdateClassroomRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		classroom.UpdateClassroom(response, &request)
		//老师上下课
	case "UPDATE_CLASSROOM_STATUS":
		request := bean.UpdateClassRoomStatusRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to DeleteStudentsRequest is:", request)
		classroom.UpdateClassRoomStatus(response, &request)
		//教室设置
	case "UPDATE_CLASSROOM_SETTING_STATUS":
		request := bean.UpdateClassRoomSettingStatusRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to DeleteStudentsRequest is:", request)
		classroom.UpdateClassRoomSettingStatus(response, &request)
	//教室禁止举手状态设置
	case "FORBID_HAND":
		request := bean.HandsForbidRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to HandsForbidRequest is:", request)
		classroom.UpdateHandForbidStatus(response, &request)

	// 获取教室信息
	case "GET_CLASS_ROOM":
		request := bean.GetClassRoomInfoRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to GetClassRoomInfoRequest is:", request)
		classroom.GetClassRoomInfo(response, &request)
		// 获取教室列表
	case "GET_CLASS_ROOM_LIST":
		request := bean.GetClassRoomListRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to GetClassRoomListRequest is:", request)
		classroom.GetClassRoomList(response, &request)
		// 获取举手列表
	case "GET_HANDS_USER_LIST":
		request := bean.GetHandsUserListRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to HandsUserListRequest is:", request)
		classroom.GetHandsUserList(response, &request)
		//举手 取消 动作请求
	case "HANDS_UP":
		request := bean.HandsUpRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to HandsUpRequest is:", request)
		classroom.HandsUp(response, &request)
		//禁止 解禁 禁言区
	case "FORBID_AREA":
		request := bean.ForbidAreaRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to ForbidAreaRequest is:", request)
		classroom.ForbidArea(response, &request)
		//清空举手列表
	case "CLEAR_HAND_LIST":
		request := bean.HandsListClearRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to HandsListClearRequest is:", request)
		classroom.HandsListClear(response, &request)
		//获取在线用户列表
	case "GET_ONLINE_USER_LIST":
		request := bean.GetOnLineUserListRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to GetOnLineUserListRequest is:", request)
		classroom.GetOnLineUserList(response, &request)
		//添加 移除举手到发言区
	case "ADD_REMOVE_SPEAK_AREA":
		request := bean.AddToSpeakAreaRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to AddToSpeakAreaRequest is:", request)
		classroom.AddDeleteToSpeakArea(response, &request)

		//进入教室
	case "ENTRY_CLASSROOM":
		request := bean.EntryClassRoomRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to EntryClassRoomRequest is:", request)
		classroom.EntryClassRoom(response, &request)
	//退出教室
	case "EXIT_CLASSROOM":
		request := bean.ExitClassRoomRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to ExitClassRoomRequest is:", request)
		classroom.ExitClassRoom(response, &request)

		//将学生移除 添加教室
	case "DELADD_STUDENTS":
		request := bean.DeleteAddStudentsRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to DeleteStudentsRequest is:", request)
		classroom.ClassRoomMemberDeleteAdd(response, &request)
		//上传文件
	case "UPLOAD_RESOURCE":
		upLoadRequest := bean.UploadResourceRequest{}
		err = json.Unmarshal(data, &upLoadRequest)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to request is:", upLoadRequest)
		classroom.UpLoadResource(response, &upLoadRequest, request)
		//删除文件
	case "RESOURCE_DELETE":
		request := resource.DeleteResourceRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to DeleteResourceRequest is:", request)
		classroom.DeleteResource(response, &request)
		//获取文件列表
	case "GET_RESOURCE_LIST":
		request := resource.ResourceListRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to ResourceListRequest is:", request)
		classroom.GetResourceList(response, &request)
		//登录
	case "LOGIN":
		request := bean.LoginRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to LoginRequest is:", request)
		user.Login(response, &request)
		//注册
	case "SIGN_UP":
		request := bean.UserSignUpRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to UserSignUpRequest is:", request)
		user.SignUp(response, &request)
		//修改用户信息
	case "USERINFO_UPDATE":
		request := bean.UserInfoUpdateRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to UserInfoUpdateRequest is:", request)
		user.UpdateUserInfo(response, &request)
		//修改密码
	case "USER_PASSWORD_UPDATE":
		request := bean.UserPwdUpdateRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to UserPwdUpdateRequest is:", request)
		user.UpdateUserPwd(response, &request)
		//重置密码
	case "RESET_USER_PASSWORD":
		request := bean.UserPwdResetRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to UserPwdResetRequest is:", request)
		user.ResetUserPwd(response, &request)
		//用户查询
	case "GET_USERINFO":
		request := bean.GetUserInfoRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to GetUserInfoRequest is:", request)
		user.GetUserInfo(response, &request)
		//用户昵称模糊查询
	case "GET_USERINFO_FUZZY":
		request := bean.GetUserInfoFuzzyRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to GetUserInfoFuzzyRequest is:", request)
		user.GetUserInfoFuzzy(response, &request)
		//所有用户查询
	case "GET_ALLUSERS":
		request := bean.GetAllUserRequest{}
		err = json.Unmarshal(data, &request)
		if err != nil {
			logs.GetLogger().Error(err.Error())
			writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
			return
		}
		logs.GetLogger().Info("Unmarshal json to GetAllUserRequest is:", request)
		user.GetAllUsers(response, &request)
		//测试在线用户
	case "TestOnline":
		logs.GetLogger().Info("Test Online:")
		//		OnLineTest.UpdateOnLineTest(10033, 10, 2)
	//	OnLineTest.OnLineTest(10, 2)
	default:
		writeErrMsg(errorcode.NO_SUCH_CMD, "no such CMD", response)
	}

}

func writeErrMsg(code int32, errMsg string, response http.ResponseWriter) {
	result := new(bean.ClientBaseResponse)
	result.Code = code
	result.ErrMsg = errMsg
	son, _ := json.Marshal(result)
	fmt.Fprintf(response, string(son))
}
