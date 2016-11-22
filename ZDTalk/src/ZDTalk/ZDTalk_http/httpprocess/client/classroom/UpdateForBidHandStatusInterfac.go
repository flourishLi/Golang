// ClassRoomSettingStatus
package classroom

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/ZDTalk_http/httpprocess/process"
	"ZDTalk/ZDTalk_http/imbean"
	"ZDTalk/ZDTalk_http/protocol"
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

func UpdateHandForbidStatus(response http.ResponseWriter, request *bean.HandsForbidRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("UpdateHandForbidStatus begins")
	//roomID为空 直接返回
	if request.ClassRoomId == ClassRoomID_IS_NULL {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}
	//请求用户的Id为空
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.REQUEST_USER_ID_CAN_NOT_NULL, errorcode.REQUEST_USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}
	// 禁止举手状态0 不禁止举手, 1 禁止举手 必须参数
	if ForbidHandStatus_ForBid != request.ForbidHandStatus && ForbidHandStatus_NO != request.ForbidHandStatus {
		writeErrMsg(errorcode.FORBID_HAND_STATUS_IS_WRONG, errorcode.FORBID_HAND_STATUS_IS_WRONG_MSG, response)
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
	result := new(bean.HandsForbidResponse) //反馈到客户端

	//内存中更新教室
	result.Code, result.ErrMsg = classRoomMemoryManager.UpdateHandForbidStatus(request.ClassRoomId, request.ForbidHandStatus)
	if result.Code == errorcode.SUCCESS {
		//数据库中更新教室 利用内存的数据更新
		result.Code, result.ErrMsg = classRoomProcessManager.UpdateForbidHandStatus(request.ClassRoomId)
		//数据库操作成功
		if result.Code == errorcode.SUCCESS {
			fbStatusResponse_protocol := SendMessage_UpdateHandForbidStatus(request, response)
			//发送消息失败
			if fbStatusResponse_protocol.Result != imbean.SUCCESS {
				writeErrMsg(fbStatusResponse_protocol.Result, fbStatusResponse_protocol.ErrorMessage, response)
				return
			}
		}
	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("UpdateHandForbidStatus classroom result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("UpdateHandForbidStatus end")
	logs.GetLogger().Info("=============================================================\n")
}

func SendMessage_UpdateHandForbidStatus(request *bean.HandsForbidRequest, response http.ResponseWriter) *protocol.ResponseBase {
	logs.GetLogger().Info("SendMessage_UpdateHandForbidStatus Is begin")
	//获取用户id对应的IMid
	userIMId := GetUserIMId(request.RequestUserId)
	//获取时间戳
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))

	//获取该教室的在线成员列表
	userMemoryManager := memory.GetUserInfoMemoryManager()
	onLineUsers := userMemoryManager.GetOnlineUserList(request.ClassRoomId)

	//获取对应的Imuserid集合
	imOnlineUsers := make([]int32, 0)
	for _, userInfo := range onLineUsers {
		imOnlineUsers = append(imOnlineUsers, userInfo.ChatId)
	}

	//构建要发送给IM的客户端数据 作为最终IM协议的Content
	updateHandForbidStatusClientProtocal := protocol.UpdateHandForbidStatusClientProtocal{request.RequestUserId, request.ClassRoomId, request.ForbidHandStatus}
	data, err := json.Marshal(updateHandForbidStatusClientProtocal)
	if err != nil {
		logs.GetLogger().Info("updateHandForbidStatusClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	//构建最终发送给IM的协议
	//requestBase设置
	updateFBStatusRequest_protocol := protocol.GetUpdateHandForbidStatusProtocol()
	updateFBStatusRequest_protocol.RequestServerId = 1
	updateFBStatusRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	updateFBStatusRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	updateFBStatusRequest_protocol.SkipDBOperat = false
	updateFBStatusRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	updateFBStatusRequest_protocol.SenderUserId = userIMId
	updateFBStatusRequest_protocol.UserIds = imOnlineUsers
	updateFBStatusRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("updateFBStatusRequest_protocol data is:", updateFBStatusRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(updateFBStatusRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("updateFBStatusRequest_protocol to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("updateRoomStatusRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	updateFBStatusResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_UpdateHandForbidStatus Handle iS error：", errOne)
		writeErrMsg(errorcode.ROOMSTATUS_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &updateFBStatusResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("SendMessage_UpdateHandForbidStatus Result IS：", updateFBStatusResponse_protocol)
	logs.GetLogger().Info("SendMessage_UpdateHandForbidStatus Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_UpdateHandForbidStatus Is end")

	return updateFBStatusResponse_protocol
}
