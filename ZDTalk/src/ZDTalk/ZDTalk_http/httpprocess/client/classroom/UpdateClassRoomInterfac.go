package classroom

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

func UpdateClassroom(response http.ResponseWriter, request *bean.UpdateClassroomRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("update classroom begins")

	//roomID为空 直接返回
	if request.ClassRoomId == ClasRoom_IS_NOT_EXIT {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}
	//请求用户的Id为空
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.REQUEST_USER_ID_CAN_NOT_NULL, errorcode.REQUEST_USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}
	//检查用户是否具有权限
	hasAuthority, isExit := AuthorityCheckout(request.RequestUserId, Teacher)
	if isExit == User_IS_NOT_EXIT {
		writeErrMsg(errorcode.USER_ID_IS_NOT_EXIT, errorcode.USER_ID_IS_NOT_EXIT_MSG, response)
		return
	}

	if !hasAuthority {
		//没有权限
		writeErrMsg(errorcode.REQUEST_USER_ID_HAS_NO_AHTHORITY, errorcode.REQUEST_USER_ID_HAS_NO_AHTHORITY_MSG, response)
		return
	}
	//process接口对象
	classRoomProcessManager := new(process.ClassRoomProcess)
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//返回参数 结构体(Code ErrMsg)
	result := new(bean.UpdateClassroomResponse) //反馈到客户端

	//调用IM同步数据
	updateGroupResponse_IM := UpdateGroup_IM(request, response)
	//IM更新群成功
	if updateGroupResponse_IM.Result == imbean.SUCCESS {
		logs.GetLogger().Info("updateGroupResponse_IM Is Sucess")
		//内存中更新教室
		result.Code, result.ErrMsg = classRoomMemoryManager.UpdateClassroom(request.ClassRoomId, request.ClassRoomName, request.ClassRoomLogo, request.Description, request.ClassRoomCourse)
		if result.Code == errorcode.SUCCESS {
			//数据库中更新教室 利用内存的数据更新
			result.Code, result.ErrMsg = classRoomProcessManager.UpdateClassRoom(request.ClassRoomId)
		}

		datas, err1 := json.Marshal(result)
		if err1 != nil {
			logs.GetLogger().Error(err1.Error())
		}
		logs.GetLogger().Info("UpdateClassroom classroom result:" + string(datas))
		fmt.Fprintln(response, string(datas))

		logs.GetLogger().Info("update classroom end")
	} else {
		writeErrMsg(updateGroupResponse_IM.Result, updateGroupResponse_IM.ErrorMessage, response)
		return
	}
	logs.GetLogger().Info("=============================================================\n")
}

//在IM服务器修改群资料
func UpdateGroup_IM(request *bean.UpdateClassroomRequest, response http.ResponseWriter) *imbean.UpdateGroupInfoResponse {
	logs.GetLogger().Info("UpdateGroup_IM Is begin")

	//获取教室id对应的ClassRoomIMid
	classroomIMId := GetClassRoomIMId(request.ClassRoomId)

	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct base
	updateGroupRequest_IM := imbean.GetUpdateGroupInfoRequest()
	updateGroupRequest_IM.RequestServerId = 1
	updateGroupRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	updateGroupRequest_IM.RequestTime = timeutils.GetTimeStamp()
	updateGroupRequest_IM.SkipDBOperat = false
	updateGroupRequest_IM.MarkId = imbean.PUSH
	//client
	updateGroupRequest_IM.ClassRoomName = request.ClassRoomName
	updateGroupRequest_IM.ClassRoomLogo = request.ClassRoomLogo
	updateGroupRequest_IM.ClassRoomIMId = classroomIMId
	updateGroupRequest_IM.Description = request.Description
	updateGroupRequest_IM.RequestUserId = request.RequestUserId
	updateGroupRequest_IM.NeedManagerPower = imbean.NOMANAGERPOWER
	logs.GetLogger().Info("updateGroupRequest_IM data is:", updateGroupRequest_IM)
	//tojson
	groupRequest_Bytes, err := json.Marshal(updateGroupRequest_IM)
	if err != nil {
		logs.GetLogger().Info("updateGroupRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("updateGroupRequest_IM to json is:", string(groupRequest_Bytes))
	//response
	updateGroupResponse_IM := new(imbean.UpdateGroupInfoResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMHttpGroupUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("UpdateGroup_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.UPDATE_GROUP_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &updateGroupResponse_IM)
	if errTwo != nil {
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		logs.GetLogger().Error("updateGroupResponse_IM unmarshal is error:", errTwo.Error())
		return nil
	}
	logs.GetLogger().Info("UpdateGroup_IM IS：", updateGroupResponse_IM)
	logs.GetLogger().Info("UpdateGroup_IM IS：", string(data))

	logs.GetLogger().Info("UpdateGroup_IM Is end")

	return updateGroupResponse_IM
}
