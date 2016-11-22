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

func DeleteClassroom(response http.ResponseWriter, request *bean.DeleteClassroomRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("delete classroom begains")

	//ClassRoomId为空
	if request.ClassRoomId == ClassRoomID_IS_NULL {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}
	//请求用户的Id为空
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.REQUEST_USER_ID_CAN_NOT_NULL, errorcode.REQUEST_USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}

	//检查用户是否具有老师权限
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
	result := new(bean.DeleteClassroomResponse) //反馈到客户端

	//调用IM同步数据
	DissolutionGroupResponse_IM := DissolutionGroup_IM(request, response)
	//IM解散群成功
	if DissolutionGroupResponse_IM.Result == imbean.SUCCESS {
		logs.GetLogger().Info("DissolutionGroupResponse_IM Is Sucess")
		//内存中删除教室
		result.Code, result.ErrMsg = classRoomMemoryManager.DeleteClassroom(request.ClassRoomId)
		if result.Code == errorcode.SUCCESS {
			//数据库中删除教室
			result.Code, result.ErrMsg = classRoomProcessManager.DeleteClassRoom(request.ClassRoomId)
		}

		datas, err1 := json.Marshal(result)
		if err1 != nil {
			logs.GetLogger().Error(err1.Error())
		}
		logs.GetLogger().Info("delete classroom result:" + string(datas))
		fmt.Fprintln(response, string(datas))

		logs.GetLogger().Info("delete classroom end")
	} else {
		writeErrMsg(DissolutionGroupResponse_IM.Result, DissolutionGroupResponse_IM.ErrorMessage, response)
		return
	}
	logs.GetLogger().Info("=============================================================\n")
}

//在IM服务器解散群
func DissolutionGroup_IM(request *bean.DeleteClassroomRequest, response http.ResponseWriter) *imbean.DissolutionGroupResponse {
	logs.GetLogger().Info("DissolutionGroup_IM Is begin")
	//获取教室id对应的ClassRoomIMid
	classroomIMId := GetClassRoomIMId(request.ClassRoomId)

	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct base
	dissolutionGroupRequest_IM := imbean.GetDissolutionGroupRequest()
	dissolutionGroupRequest_IM.RequestServerId = 1
	dissolutionGroupRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	dissolutionGroupRequest_IM.RequestTime = timeutils.GetTimeStamp()
	dissolutionGroupRequest_IM.SkipDBOperat = false
	dissolutionGroupRequest_IM.MarkId = imbean.PUSH
	//client
	dissolutionGroupRequest_IM.UserId = request.RequestUserId
	dissolutionGroupRequest_IM.ClassRoomIMIds = []int32{classroomIMId}
	logs.GetLogger().Info("dissolutionGroupRequest_IM data is:", dissolutionGroupRequest_IM)
	//json
	groupRequest_Bytes, err := json.Marshal(dissolutionGroupRequest_IM)
	if err != nil {
		logs.GetLogger().Info("dissolutionGroupRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("dissolutionGroupRequest_IM to json is:", string(groupRequest_Bytes))
	//response
	dissolutionGroupResponse_IM := new(imbean.DissolutionGroupResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMHttpGroupUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("DissolutionGroup_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.DISSOLUTION_GROUP_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &dissolutionGroupResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("DissolutionGroup_IM Result IS：", dissolutionGroupResponse_IM)
	logs.GetLogger().Info("DissolutionGroup_IM Result IS：", string(data))
	logs.GetLogger().Info("DissolutionGroup_IM Is end")

	return dissolutionGroupResponse_IM
}
