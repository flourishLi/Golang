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

func UpdateClassRoomSettingStatus(response http.ResponseWriter, request *bean.UpdateClassRoomSettingStatusRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("UpdateClassRoomSettingStatus begins")
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
	//settingStatus 教室设置不合法 教室设置状态 1 允许所有人打字，2 允许学生录音，3 禁止游客打字，4 禁止游客举手
	for _, status := range request.SettingStatus {
		if SettingStatus_One == status || SettingStatus_Two == status || SettingStatus_Three == status || SettingStatus_Four == status {
			continue
		} else {
			writeErrMsg(errorcode.SETTING_STATUS_IS_WRONG, errorcode.SETTING_STATUS_IS_WRONG_MSG, response)
			return
		}
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
	result := new(bean.UpdateClassRoomSettingStatusResponse) //反馈到客户端

	//内存中更新教室
	result.Code, result.ErrMsg = classRoomMemoryManager.UpdateClassRoomSettingStatus(request.ClassRoomId, request.RequestUserId, request.SettingStatus)
	if result.Code == errorcode.SUCCESS {
		//数据库中更新教室 利用内存的数据更新
		result.Code, result.ErrMsg = classRoomProcessManager.UpdateClassRoomSettingStatus(request.ClassRoomId, request.RequestUserId)
		//数据库操作成功
		if result.Code == errorcode.SUCCESS {
			settingStatusResponse_protocol := SendMessage_UpdateClassRoomSettingStatus(request, response)
			//发送消息失败
			if settingStatusResponse_protocol.Result != imbean.SUCCESS {
				writeErrMsg(settingStatusResponse_protocol.Result, settingStatusResponse_protocol.ErrorMessage, response)
				return
			}
		}

	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("UpdateClassRoomSettingStatus classroom result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("UpdateClassRoomSettingStatus end")
	logs.GetLogger().Info("=============================================================\n")
}

func SendMessage_UpdateClassRoomSettingStatus(request *bean.UpdateClassRoomSettingStatusRequest, response http.ResponseWriter) *protocol.ResponseBase {
	logs.GetLogger().Info("SendMessage_UpdateClassRoomSettingStatus Is begin")
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
	updateSettingStatusClientProtocal := protocol.UpdateClassRoomSettingStatusClientProtocal{request.RequestUserId, request.ClassRoomId, request.SettingStatus}
	data, err := json.Marshal(updateSettingStatusClientProtocal)
	if err != nil {
		logs.GetLogger().Info("updateSettingStatusClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	//构建最终发送给IM的协议
	//requestBase设置
	updateSettingStatusRequest_protocol := protocol.GetUpdateClassRoomSettingStatusProtocol()
	updateSettingStatusRequest_protocol.RequestServerId = 1
	updateSettingStatusRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	updateSettingStatusRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	updateSettingStatusRequest_protocol.SkipDBOperat = false
	updateSettingStatusRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	updateSettingStatusRequest_protocol.SenderUserId = userIMId
	updateSettingStatusRequest_protocol.UserIds = imOnlineUsers
	updateSettingStatusRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("updateSettingStatusRequest_protocol data is:", updateSettingStatusRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(updateSettingStatusRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("updateSettingStatusRequest_protocol to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("handsForbidRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	updateSettingStatusResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_UpdateClassRoomSettingStatus Handle iS error：", errOne)
		writeErrMsg(errorcode.SETTINGSTATUS_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &updateSettingStatusResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("SendMessage_UpdateClassRoomSettingStatus Result IS：", updateSettingStatusResponse_protocol)
	logs.GetLogger().Info("SendMessage_UpdateClassRoomSettingStatus Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_UpdateClassRoomSettingStatus Is end")

	return updateSettingStatusResponse_protocol
}
