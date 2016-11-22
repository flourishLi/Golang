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

func ExitClassRoom(response http.ResponseWriter, request *bean.ExitClassRoomRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("ExitClassRoom begins")

	//ClassRoomId为空 直接返回
	if request.ClassRoomId == ClassRoomID_IS_NULL {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}
	//UserId为空 直接返回
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.REQUEST_USER_ID_CAN_NOT_NULL, errorcode.REQUEST_USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}

	//process接口对象
	classRoomProcessManager := new(process.ClassRoomProcess)
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()
	//返回参数 结构体(Code ErrMsg)
	result := new(bean.ExitClassRoomResponse) //反馈到客户端
	onLineMemoryManager := memory.GetQueueOnlineMemoryManager()
	//有异常时，将用户重新添加到在线用户列表中
	onLineUser := onLineMemoryManager.GetMemoryOnlineUser(request.ClassRoomId, request.RequestUserId)
	//内存中更新教室
	result.Code, result.ErrMsg = classRoomMemoryManager.ExitClassRoom(request.ClassRoomId, request.RequestUserId)
	if result.Code != errorcode.SUCCESS {
		onLineMemoryManager.AddMemoryLockOnlineUser(onLineUser)
		logs.GetLogger().Error("Http接口 用户 %d 退出教室 失败", request.RequestUserId, request.ClassRoomId)
		writeErrMsg(result.Code, result.ErrMsg, response)
		return
	} else {

		//数据库中更新教室 利用内存的数据更新
		result.Code, result.ErrMsg = classRoomProcessManager.ExitClassRoom(request.ClassRoomId)
		//数据库操作成功
		if result.Code != errorcode.SUCCESS {
			onLineMemoryManager.AddMemoryLockOnlineUser(onLineUser)
			writeErrMsg(result.Code, result.ErrMsg, response)
			return
		} else {
			onLineMemoryManager.RemoveMemoryLockOnlineUser(onLineUser.ClassRoomId, onLineUser.UserId)
			code, errMsg := SendMessage_ExitClassRoom(request, response)
			//发送消息失败
			if code != imbean.SUCCESS {
				writeErrMsg(code, errMsg, response)
				return
			}
		}
	}
	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("ExitClassRoom  result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("ExitClassRoom end")
	logs.GetLogger().Info("=============================================================\n")
}

//向IM发送通知
func SendMessage_ExitClassRoom(request *bean.ExitClassRoomRequest, response http.ResponseWriter) (int32, string) {
	logs.GetLogger().Info("SendMessage_ExitClassRoom Is begin")
	//获取用户id对应的IMid
	userIMId := GetUserIMId(request.RequestUserId)
	//获取时间戳
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))

	//获取该教室的在线成员列表
	userMemoryManager := memory.GetUserInfoMemoryManager()
	onLineUsers := userMemoryManager.GetOnlineUserList(request.ClassRoomId)
	length := len(onLineUsers)
	if length == 0 {

		logs.GetLogger().Info("当前教室 %d 的在线数量 为 0 不发送用户 %d 退出教室的自定义通知", request.ClassRoomId, request.RequestUserId)
		return imbean.SUCCESS, ""
	}

	//获取对应的Imuserid集合
	imOnlineUsers := make([]int32, 0)
	for _, userInfo := range onLineUsers {
		if userInfo.ChatId != userIMId {
			imOnlineUsers = append(imOnlineUsers, userInfo.ChatId)
		}
	}

	//构建要发送给IM的客户端数据 作为最终IM协议的Content
	exitClientProtocal := protocol.ExitClassRoomClientProtocal{request.RequestUserId, request.ClassRoomId}
	data, err := json.Marshal(exitClientProtocal)
	if err != nil {
		logs.GetLogger().Info("exitClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return errorcode.JSON_PRASEM_ERROR, err.Error()
	}
	//构建最终发送给IM的协议
	//requestBase设置
	exitRoomRequest_protocol := protocol.GetExitClassRoomProtocol()
	exitRoomRequest_protocol.RequestServerId = 1
	exitRoomRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	exitRoomRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	exitRoomRequest_protocol.SkipDBOperat = false
	exitRoomRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	exitRoomRequest_protocol.SenderUserId = userIMId
	exitRoomRequest_protocol.UserIds = imOnlineUsers
	exitRoomRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("exitRoomRequest_protocol data is:", exitRoomRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(exitRoomRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("exitRoomRequest_protocol to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return errorcode.JSON_PRASEM_ERROR, err.Error()
	}
	logs.GetLogger().Info("exitRoomRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	exitRoomResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_ExitClassRoom Handle iS error：", errOne)
		writeErrMsg(errorcode.EXIT_CLASSROOM_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return errorcode.EXIT_CLASSROOM_PROTOCOL_IM_IS_ERROR, err.Error()
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &exitRoomResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return errorcode.JSON_PRASEM_ERROR, err.Error()
	}
	logs.GetLogger().Info("SendMessage_ExitClassRoom Result IS：", exitRoomResponse_protocol)
	logs.GetLogger().Info("SendMessage_ExitClassRoom Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_ExitClassRoom Is end")

	return imbean.SUCCESS, ""
}
