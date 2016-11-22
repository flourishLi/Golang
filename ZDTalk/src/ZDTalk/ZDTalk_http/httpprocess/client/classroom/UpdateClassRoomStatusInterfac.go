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

func UpdateClassRoomStatus(response http.ResponseWriter, request *bean.UpdateClassRoomStatusRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("UpdateClassRoomStatus begins")
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
	//roomStatus 教室状态合法	教室当前状态0 刚创建, 1 初始状态, 2 上课状态, 3 下课状态 必须参数
	if RoomStatus_One != request.ClassRoomStatus && RoomStatus_Two != request.ClassRoomStatus && RoomStatus_Three != request.ClassRoomStatus && RoomStatus_ZERO != request.ClassRoomStatus {
		writeErrMsg(errorcode.CLASS_ROOM_STATUS_IS_WRONG, errorcode.CLASS_ROOM_STATUS_IS_WRONG_MSG, response)
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
	result := new(bean.UpdateClassRoomStatusResponse) //反馈到客户端

	//内存中更新教室
	result.Code, result.ErrMsg = classRoomMemoryManager.UpdateClassRoomStatus(request.ClassRoomId, request.RequestUserId, request.ClassRoomStatus)
	if result.Code == errorcode.SUCCESS {
		//数据库中更新教室 利用内存的数据更新
		result.Code, result.ErrMsg = classRoomProcessManager.UpdateClassRoomStatus(request.ClassRoomId, request.RequestUserId)
		//数据库操作成功
		if result.Code == errorcode.SUCCESS {
			roomStatusResponse_protocol := SendMessage_UpdateRoomStatus(request, response)
			//发送消息失败
			if roomStatusResponse_protocol.Result != imbean.SUCCESS {
				writeErrMsg(roomStatusResponse_protocol.Result, roomStatusResponse_protocol.ErrorMessage, response)
				return
			}
		}
	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("UpdateClassRoomStatus classroom result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("UpdateClassRoomStatus end")
	logs.GetLogger().Info("=============================================================\n")
}

func SendMessage_UpdateRoomStatus(request *bean.UpdateClassRoomStatusRequest, response http.ResponseWriter) *protocol.ResponseBase {
	logs.GetLogger().Info("SendMessage_UpdateRoomStatus Is begin")
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
	updateRoomStatusClientProtocal := protocol.UpdateRoomStatusClientProtocal{request.RequestUserId, request.ClassRoomId, request.ClassRoomStatus}
	data, err := json.Marshal(updateRoomStatusClientProtocal)
	if err != nil {
		logs.GetLogger().Info("updateRoomStatusClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	//构建最终发送给IM的协议
	//requestBase设置
	updateRoomStatusRequest_protocol := protocol.GetUpdateClassRoomStatusProtocol()
	updateRoomStatusRequest_protocol.RequestServerId = 1
	updateRoomStatusRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	updateRoomStatusRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	updateRoomStatusRequest_protocol.SkipDBOperat = false
	updateRoomStatusRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	updateRoomStatusRequest_protocol.SenderUserId = userIMId
	updateRoomStatusRequest_protocol.UserIds = imOnlineUsers
	updateRoomStatusRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("updateRoomStatusRequest_protocol data is:", updateRoomStatusRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(updateRoomStatusRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("updateRoomStatusRequest_protocol to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("updateRoomStatusRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	updateRoomStatusResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_UpdateRoomStatus Handle iS error：", errOne)
		writeErrMsg(errorcode.ROOMSTATUS_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &updateRoomStatusResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("SendMessage_UpdateRoomStatus Result IS：", updateRoomStatusResponse_protocol)
	logs.GetLogger().Info("SendMessage_UpdateRoomStatus Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_UpdateRoomStatus Is end")

	return updateRoomStatusResponse_protocol
}
