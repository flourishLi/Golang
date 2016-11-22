package memory

import (
	"ZDTalk/ZDTalk_http/imbean"
	"ZDTalk/ZDTalk_http/protocol"
	"ZDTalk/config"
	"ZDTalk/utils/cryptoutils"
	"ZDTalk/utils/httputil"
	logs "ZDTalk/utils/log4go"
	"ZDTalk/utils/timeutils"
	"encoding/json"

	"strconv"
)

//掉线发通知
//向IM发送通知
func SendMessage_ExitClassRoom(userIdZD int32, classRoomId int32, receiveUserSlice []int32) {
	logs.GetLogger().Info("SendMessage_ExitClassRoom Is begin In timeOut")
	if userIdZD == 0 {
		logs.GetLogger().Error("用户Id 为 0 不发送用户退出教室的自定义推送")
		return
	}
	if GetUserInfoMemoryManager().GetUserIMId(userIdZD) == 0 {
		logs.GetLogger().Error("用户IMUserId 为 0 不发送用户退出教室的自定义推送")
		return
	}
	if classRoomId == 0 {
		logs.GetLogger().Error("教室Id 为0 不发送用户退出教室的自定义推送")
		return
	}
	if receiveUserSlice == nil || len(receiveUserSlice) == 0 {
		logs.GetLogger().Error("教室 %d 接收某人上线的用户集合为空 不发送用户退出教室的自定义推送", classRoomId)
		return
	}

	//获取用户id对应的IMid
	userIMId := GetUserInfoMemoryManager().GetUserIMId(userIdZD)
	//获取时间戳
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))

	receiveIMUserSlice := make([]int32, 0)
	for _, value := range receiveUserSlice {
		if value != userIdZD {
			imUserId := GetUserInfoMemoryManager().GetUserIMId(value)
			receiveIMUserSlice = append(receiveIMUserSlice, imUserId)
		}
	}

	//构建要发送给IM的客户端数据 作为最终IM协议的Content
	exitClientProtocal := protocol.ExitClassRoomClientProtocal{userIMId, classRoomId}
	data, err := json.Marshal(exitClientProtocal)
	if err != nil {
		logs.GetLogger().Info("exitClientProtocal Marshal err:", err)
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
	exitRoomRequest_protocol.UserIds = receiveIMUserSlice
	exitRoomRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("exitRoomRequest_protocol data is:", exitRoomRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(exitRoomRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("exitRoomRequest_protocol to json is Wrong", err)
	}
	logs.GetLogger().Info("exitRoomRequest_protocol to json is:", string(groupRequest_Bytes))

	//向IM请求数据
	_, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_ExitClassRoom Handle iS error：", errOne)
	}

}

//心跳时 内存的在线用户列表中不存在 也要发送登录教室的通知 比如断网
func SendMessage_EntryClassRoom(userIdZD int32, classRoomId int32, receiveUserSlice []int32) {
	logs.GetLogger().Info("SendMessage_EntryClassRoom Is begin")
	if userIdZD == 0 {
		logs.GetLogger().Error("用户Id 为 0 不发送用户进入教室的自定义推送")
		return
	}
	if GetUserInfoMemoryManager().GetUserIMId(userIdZD) == 0 {
		logs.GetLogger().Error("用户IMUserId 为 0 不发送用户进入教室的自定义推送")
		return
	}
	if classRoomId == 0 {
		logs.GetLogger().Error("教室Id 为0 不发送用户进入教室的自定义推送")
		return
	}
	if receiveUserSlice == nil || len(receiveUserSlice) == 0 {
		logs.GetLogger().Error("接收某人上线的用户集合为空 不发送用户进入教室的自定义推送")
		return
	}

	//获取时间戳
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	userIMId := GetUserInfoMemoryManager().GetUserIMId(userIdZD)
	receiveIMUserSlice := []int32{}
	for _, value := range receiveUserSlice {
		if value != userIdZD {
			imUserId := GetUserInfoMemoryManager().GetUserIMId(value)
			receiveIMUserSlice = append(receiveIMUserSlice, imUserId)
		}
	}
	logs.GetLogger().Info("接收用户进入教室的自定义推送的用户集合 %d", receiveIMUserSlice)

	//构建要发送给IM的客户端数据 作为最终IM协议的Content
	entryClientProtocal := protocol.EntryClassRoomClientProtocal{userIMId, classRoomId}
	data, err := json.Marshal(entryClientProtocal)
	if err != nil {
		logs.GetLogger().Info("entryClientProtocal Marshal err:", err)
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
	entryRoomRequest_protocol.UserIds = receiveIMUserSlice
	entryRoomRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("entryRoomRequest_protocol data is:", entryRoomRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(entryRoomRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("entryRoomRequest_protocol to json is Wrong", err)
	}
	logs.GetLogger().Info("entryRoomRequest_protocol to json is:", string(groupRequest_Bytes))

	//向IM请求数据
	_, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_EntryClassRoom Handle iS error：", errOne)
	}
}

//退出发言列表通知
func SendMessage_DeleteAddSpeaker(userIdZD int32, classRoomId int32, receiveUserSlice []int32) {
	logs.GetLogger().Info("SendMessage_DeleteAddSpeaker Is begin")

	if userIdZD == 0 {
		logs.GetLogger().Error("用户Id 为 0 不发送用户退出教室的自定义推送")
		return
	}
	if GetUserInfoMemoryManager().GetUserIMId(userIdZD) == 0 {
		logs.GetLogger().Error("用户IMUserId 为 0 不发送用户退出教室的自定义推送")
		return
	}
	if classRoomId == 0 {
		logs.GetLogger().Error("教室Id 为0 不发送用户退出教室的自定义推送")
		return
	}
	if receiveUserSlice == nil || len(receiveUserSlice) == 0 {
		logs.GetLogger().Error("教室 %d 接收某人上线的用户集合为空 不发送用户退出教室的自定义推送", classRoomId)
		return
	}

	//获取用户id对应的IMid
	userIMId := GetUserInfoMemoryManager().GetUserIMId(userIdZD)
	//获取时间戳
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))

	//获取对应的Imuserid集合
	imOnlineUsers := make([]int32, 0)
	for _, userid := range receiveUserSlice {
		if userid != userIdZD {
			id := GetUserInfoMemoryManager().GetUserIMId(userid)
			imOnlineUsers = append(imOnlineUsers, id)
		}
	}

	//构建要发送给IM的客户端数据 作为最终IM协议的Content
	deleteAddSpeakerClientProtocal := protocol.DeleteAddSpeakerToAreaClientProtocal{userIdZD, imOnlineUsers, 2, classRoomId}
	data, err := json.Marshal(deleteAddSpeakerClientProtocal)
	if err != nil {
		logs.GetLogger().Info("deleteAddSpeakerClientProtocal Marshal err:", err)
		return
	}
	logs.GetLogger().Info("the content is：", string(data))
	//构建最终发送给IM的协议
	//requestBase设置
	DeleteAddSpeakerRequest_protocol := protocol.GetDeleteAddSpeakerToAreaProtocol()
	DeleteAddSpeakerRequest_protocol.RequestServerId = 1
	DeleteAddSpeakerRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	DeleteAddSpeakerRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	DeleteAddSpeakerRequest_protocol.SkipDBOperat = false
	DeleteAddSpeakerRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	DeleteAddSpeakerRequest_protocol.SenderUserId = userIMId
	DeleteAddSpeakerRequest_protocol.UserIds = imOnlineUsers
	DeleteAddSpeakerRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("DeleteAddSpeakerRequest_protocol data is:", DeleteAddSpeakerRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(DeleteAddSpeakerRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("DeleteAddSpeakerRequest_protocol to json is Wrong", err)
		return
	}
	logs.GetLogger().Info("DeleteAddSpeakerRequest_protocol to json is:", string(groupRequest_Bytes))

	//向IM请求数据
	_, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_DeleteAddSpeaker Handle iS error：", errOne)
		return
	}

	return
}

//退出举手列表通知
func SendMessage_HandsUp(userIdZD int32, classRoomId int32, receiveUserSlice []int32) {
	logs.GetLogger().Info("SendMessage_HandsUp_IM Is begin")

	if userIdZD == 0 {
		logs.GetLogger().Error("用户Id 为 0 不发送用户退出教室的自定义推送")
		return
	}
	if GetUserInfoMemoryManager().GetUserIMId(userIdZD) == 0 {
		logs.GetLogger().Error("用户IMUserId 为 0 不发送用户退出教室的自定义推送")
		return
	}
	if classRoomId == 0 {
		logs.GetLogger().Error("教室Id 为0 不发送用户退出教室的自定义推送")
		return
	}
	if receiveUserSlice == nil || len(receiveUserSlice) == 0 {
		logs.GetLogger().Error("教室 %d 接收某人上线的用户集合为空 不发送用户退出教室的自定义推送", classRoomId)
		return
	}

	//获取用户id对应的IMid
	userIMId := GetUserInfoMemoryManager().GetUserIMId(userIdZD)
	//获取时间戳
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))

	//获取对应的Imuserid集合
	imOnlineUsers := make([]int32, 0)
	for _, userId := range receiveUserSlice {
		if userId != userIdZD {
			id := GetUserInfoMemoryManager().GetUserIMId(userId)
			imOnlineUsers = append(imOnlineUsers, id)
		}
	}

	//构建要发送给IM的客户端数据 作为最终IM协议的Content
	handsUpClientProtocal := protocol.HandsUpClientProtocal{classRoomId, userIdZD, 2}
	data, err := json.Marshal(handsUpClientProtocal)
	if err != nil {
		logs.GetLogger().Info("handsUpClientProtocal Marshal err:", err)
		return
	}
	//构建最终发送给IM的协议
	//requestBase设置
	handsUpRequest_protocol := protocol.GetHandsUpProtocol()
	handsUpRequest_protocol.RequestServerId = 1
	handsUpRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	handsUpRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	handsUpRequest_protocol.SkipDBOperat = false
	handsUpRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	handsUpRequest_protocol.SenderUserId = userIMId
	handsUpRequest_protocol.UserIds = imOnlineUsers
	handsUpRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("handsUpRequest_protocol data is:", handsUpRequest_protocol)

	//handsUpRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(handsUpRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("handsUpRequest_protocol to json is Wrong", err)
		return
	}
	logs.GetLogger().Info("handsUpRequest_protocol to json is:", string(groupRequest_Bytes))

	//向IM请求数据
	_, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("handsUpRequest_protocol Handle iS error：", errOne)
		return
	}
	return
}
