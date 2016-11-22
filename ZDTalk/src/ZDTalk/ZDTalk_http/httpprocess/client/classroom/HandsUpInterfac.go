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
	"time"
)

func HandsUp(response http.ResponseWriter, request *bean.HandsUpRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("HandsUp begins")

	//ClassRoomId为空 直接返回
	if request.ClassRoomId == ClasRoom_IS_NOT_EXIT {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}

	//请求用户的Id为空
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.REQUEST_USER_ID_CAN_NOT_NULL, errorcode.REQUEST_USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}

	//handsType值错误 直接返回
	if request.HandsType != HandsType_Up && request.HandsType != HandsType_Down {
		writeErrMsg(errorcode.HANDS_TYPE_IS_WRONG, errorcode.HAND_TYPE_IS_WRONG_MSG, response)
		return
	}

	//检查用户是否具有权限
	hasAuthority, isExit := AuthorityCheckout(request.RequestUserId, Teacher)
	if isExit == User_IS_NOT_EXIT {
		writeErrMsg(errorcode.USER_ID_IS_NOT_EXIT, errorcode.USER_ID_IS_NOT_EXIT_MSG, response)
		return
	}

	if hasAuthority {
		//没有权限 只有学生举手
		writeErrMsg(errorcode.REQUEST_USER_ID_HAS_NO_AHTHORITY, errorcode.REQUEST_USER_ID_HAS_NO_AHTHORITY_MSG, response)
		return
	}
	//process接口对象
	classRoomProcessManager := new(process.ClassRoomProcess)
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()
	//返回参数 结构体(Code ErrMsg)
	result := new(bean.HandsUpResponse) //反馈到客户端
	//举手时间
	handsTime := time.Now().Unix()
	//内存中更新教室
	result.Code, result.ErrMsg = classRoomMemoryManager.HandsUp(request.ClassRoomId, request.RequestUserId, request.HandsType, handsTime)
	if result.Code == errorcode.SUCCESS {
		//数据库中更新教室 利用内存的数据更新
		result.Code, result.ErrMsg = classRoomProcessManager.HandsUp(request.ClassRoomId)
		//数据库操作成功
		if result.Code == errorcode.SUCCESS {
			handsUpResponse_protocol := SendMessage_HandsUp(request, response)
			//发送消息失败
			if handsUpResponse_protocol.Result != imbean.SUCCESS {
				writeErrMsg(handsUpResponse_protocol.Result, handsUpResponse_protocol.ErrorMessage, response)
				return
			}
		}
	}
	//最后向客户端返回数据
	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("HandsUp  result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("HandsUp end")
	logs.GetLogger().Info("=============================================================\n")

}

func SendMessage_HandsUp(request *bean.HandsUpRequest, response http.ResponseWriter) *protocol.ResponseBase {
	logs.GetLogger().Info("SendMessage_HandsUp_IM Is begin")
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
	handsUpClientProtocal := protocol.HandsUpClientProtocal{request.ClassRoomId, request.RequestUserId, request.HandsType}
	data, err := json.Marshal(handsUpClientProtocal)
	if err != nil {
		logs.GetLogger().Info("handsUpClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
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
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("handsUpRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	handsUpResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("handsUpRequest_protocol Handle iS error：", errOne)
		writeErrMsg(errorcode.HANDSUP_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &handsUpResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("SendMessage_HandsUp Result IS：", handsUpResponse_protocol)
	logs.GetLogger().Info("SendMessage_HandsUp Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_HandsUp Is end")

	return handsUpResponse_protocol
}
