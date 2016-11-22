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

func ForbidArea(response http.ResponseWriter, request *bean.ForbidAreaRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("ForbidArea begins")

	//ClassRoomId为空 直接返回
	if request.ClassRoomId == ClassRoomID_IS_NULL {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}
	//UserId为空 直接返回
	if len(request.StudentIds) == StudentIds_IS_NULL {
		writeErrMsg(errorcode.REQUEST_USER_ID_CAN_NOT_NULL, errorcode.REQUEST_USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}
	//crudType值错误 必须为1 禁止 2解禁直接返回
	if request.CrudType != CrudType_ADD && request.CrudType != CrudType_DELETE {
		writeErrMsg(errorcode.CURD_TYPE_IS_WRONG, errorcode.CURD_TYPE_IS_WRONG_MSG, response)
		return
	}
	//UserId为空 直接返回
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
	classRoomProcess := new(process.ClassRoomProcess)
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()
	//返回参数 结构体(Code ErrMsg)
	result := new(bean.ForbidAreaResponse) //反馈到客户端

	//内存中更新教室
	result.Code, result.ErrMsg = classRoomMemoryManager.ForbidArea(request.ClassRoomId, request.CrudType, request.StudentIds)
	if result.Code == errorcode.SUCCESS {
		//数据库中更新教室 利用内存的数据更新
		result.Code, result.ErrMsg = classRoomProcess.ForbidArea(request.ClassRoomId)
		//数据库操作成功
		if result.Code == errorcode.SUCCESS {
			handsForbidResponse_protocol := SendMessage_ForbidArea(request, response)
			//发送消息失败
			if handsForbidResponse_protocol.Result != imbean.SUCCESS {
				writeErrMsg(handsForbidResponse_protocol.Result, handsForbidResponse_protocol.ErrorMessage, response)
				return
			}
		}
	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("ForbidArea  result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("ForbidArea end")
	logs.GetLogger().Info("=============================================================\n")
}

func SendMessage_ForbidArea(request *bean.ForbidAreaRequest, response http.ResponseWriter) *protocol.ResponseBase {
	logs.GetLogger().Info("SendMessage_HandsForbid Is begin")
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
	handsForbidClientProtocal := protocol.ForbidAreaClientProtocal{request.RequestUserId, request.StudentIds, request.CrudType, request.ClassRoomId}
	data, err := json.Marshal(handsForbidClientProtocal)
	if err != nil {
		logs.GetLogger().Info("handsForbidClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	//构建最终发送给IM的协议
	//requestBase设置
	handsForbidRequest_protocol := protocol.GetForbidAreaProtocol()
	handsForbidRequest_protocol.RequestServerId = 1
	handsForbidRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	handsForbidRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	handsForbidRequest_protocol.SkipDBOperat = false
	handsForbidRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	handsForbidRequest_protocol.SenderUserId = userIMId
	handsForbidRequest_protocol.UserIds = imOnlineUsers
	handsForbidRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("handsForbidRequest_protocol data is:", handsForbidRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(handsForbidRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("handsForbidRequest_protocol to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("handsForbidRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	handsForbidResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("handsForbid_protocol Handle iS error：", errOne)
		writeErrMsg(errorcode.HANDSFORBID_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &handsForbidResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("SendMessage_HandsForbid Result IS：", handsForbidResponse_protocol)
	logs.GetLogger().Info("SendMessage_HandsForbid Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_HandsForbid Is end")

	return handsForbidResponse_protocol
}
