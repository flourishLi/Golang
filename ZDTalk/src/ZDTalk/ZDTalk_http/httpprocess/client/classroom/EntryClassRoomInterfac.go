package classroom

//进入教室接口
import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/ZDTalk_http/httpprocess/process"
	"ZDTalk/ZDTalk_http/imbean"
	"ZDTalk/ZDTalk_http/protocol"
	"ZDTalk/config"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	"ZDTalk/queue/transmitter"
	"ZDTalk/utils/cryptoutils"
	"ZDTalk/utils/httputil"
	logs "ZDTalk/utils/log4go"
	"ZDTalk/utils/timeutils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func EntryClassRoom(response http.ResponseWriter, request *bean.EntryClassRoomRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("EntryClassRoom begins")

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
	result := new(bean.EntryClassRoomResponse) //反馈到客户端

	onLineMemoryManager := memory.GetQueueOnlineMemoryManager()
	onLineMember := &memory.QueueOnlineMember{}
	onLineMember.UserId = request.RequestUserId
	onLineMember.ClassRoomId = request.ClassRoomId
	onLineMember.SetKeepAliving(true)
	onLineMember.StartTime = timeutils.GetUnix13NowTime()
	//添加到在线用户列表中
	onLineMemoryManager.AddMemoryLockOnlineUser(onLineMember)

	//内存中更新教室
	result.Code, result.ErrMsg = classRoomMemoryManager.EntryClassRoom(request.ClassRoomId, request.RequestUserId)
	if result.Code == errorcode.SUCCESS {
		//数据库中更新教室 利用内存的数据更新
		result.Code, result.ErrMsg = classRoomProcessManager.EntryClassRoom(request.ClassRoomId)
		//数据库操作成功
		if result.Code == errorcode.SUCCESS {
			code, errMsg := SendMessage_EntryClassRoom(request, response)
			//发送消息失败
			if code != imbean.SUCCESS {
				writeErrMsg(code, errMsg, response)
				return
			}
			onLineMemoryManager.StartDetectOnLineHeartBeat(onLineMember)
			//给进入教室的人发送最新的绘制命令
			ReceiveDrawCommand(request.RequestUserId, request.ClassRoomId)
		} else {
			classRoomMemoryManager.RemoveOnLineUser(request.ClassRoomId, request.RequestUserId)
			writeErrMsg(result.Code, result.ErrMsg, response)
			return
		}

	} else {
		classRoomMemoryManager.RemoveOnLineUser(request.ClassRoomId, request.RequestUserId)
		writeErrMsg(result.Code, result.ErrMsg, response)
		return
	}
	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("EntryClassRoom  result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("EntryClassRoom end")
	logs.GetLogger().Info("=============================================================\n")
}

//向IM发送通知
func SendMessage_EntryClassRoom(request *bean.EntryClassRoomRequest, response http.ResponseWriter) (int32, string) {
	logs.GetLogger().Info("SendMessage_EntryClassRoom Is begin")
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
		if userInfo.ChatId != userIMId {
			imOnlineUsers = append(imOnlineUsers, userInfo.ChatId)
		}
	}
	//
	var userIdStr string
	length := len(imOnlineUsers)
	logs.GetLogger().Info("当前教室 %d 的在线数量为 %d 正在发送用户%d 进入教室的自定义通知", request.ClassRoomId, length, request.RequestUserId)

	if length == 0 {
		logs.GetLogger().Info("当前教室 %d 的在线数量 为 0 不发送进入教室的自定义通知", request.ClassRoomId)
		return imbean.SUCCESS, ""
	}
	for _, v := range imOnlineUsers {
		userIdStr += "["
		userIdStr += fmt.Sprintf("%d", v)
		userIdStr += "] "
	}

	logs.GetLogger().Info("当前教室 %d 的在线数量 %d 用户Id集合为 %s", request.ClassRoomId, length, userIdStr)

	//构建要发送给IM的客户端数据 作为最终IM协议的Content
	entryClientProtocal := protocol.EntryClassRoomClientProtocal{request.RequestUserId, request.ClassRoomId}
	data, err := json.Marshal(entryClientProtocal)
	if err != nil {
		logs.GetLogger().Info("entryClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return errorcode.JSON_PRASEM_ERROR, err.Error()
	}
	//构建最终发送给IM的协议
	//requestBase设置
	entryRoomRequest_protocol := protocol.GetEntryClassRoomProtocol()
	entryRoomRequest_protocol.RequestServerId = 1
	entryRoomRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	entryRoomRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	entryRoomRequest_protocol.SkipDBOperat = false
	entryRoomRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	entryRoomRequest_protocol.SenderUserId = userIMId
	entryRoomRequest_protocol.UserIds = imOnlineUsers
	entryRoomRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("entryRoomRequest_protocol data is:", entryRoomRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(entryRoomRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("entryRoomRequest_protocol to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return errorcode.JSON_PRASEM_ERROR, err.Error()
	}
	logs.GetLogger().Info("entryRoomRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	entryRoomResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_EntryClassRoom Handle iS error：", errOne)
		writeErrMsg(errorcode.STD_ENTRY_CLASSROOM_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return errorcode.STD_ENTRY_CLASSROOM_PROTOCOL_IM_IS_ERROR, err.Error()
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &entryRoomResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return errorcode.JSON_PRASEM_ERROR, err.Error()
	}
	logs.GetLogger().Info("SendMessage_EntryClassRoom Result IS：", entryRoomResponse_protocol)
	logs.GetLogger().Info("SendMessage_EntryClassRoom Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_EntryClassRoom Is end")

	return imbean.SUCCESS, ""
}

//进入教室接收最新的绘制命令
func ReceiveDrawCommand(userId int32, classRoomId int32) {
	receiveImUserId := GetUserIMId(userId)
	drawMemoryManager := memory.GetDrawCommandMemoryManager()
	drawContent := drawMemoryManager.DrawContent[classRoomId]

	if drawContent == nil {
		logs.GetLogger().Info("当前教室 %d 绘制命令集合为空 不发送绘制命令", classRoomId)
	} else {
		transmitter.SendMessage2IMServer(drawContent.ServerId, receiveImUserId, receiveImUserId, drawContent.MessageFormat, drawContent.MessageContent, drawContent.MessageId, drawContent.SendTime)
	}
}
