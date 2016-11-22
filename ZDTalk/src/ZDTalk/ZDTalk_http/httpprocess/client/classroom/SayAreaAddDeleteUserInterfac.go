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

func AddDeleteToSpeakArea(response http.ResponseWriter, request *bean.AddToSpeakAreaRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("AddToSpeakArea begins")

	//ClassRoomId为空 直接返回
	if request.ClassRoomId == ClasRoom_IS_NOT_EXIT {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}
	//UserId为空 直接返回
	if len(request.StudentIds) == StudentIds_IS_NULL {
		writeErrMsg(errorcode.REQUEST_USER_ID_CAN_NOT_NULL, errorcode.REQUEST_USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}
	//crudType值错误 必须为1 添加 2删除 直接返回
	if request.CrudType != CrudType_ADD && request.CrudType != CrudType_DELETE {
		writeErrMsg(errorcode.CURD_TYPE_IS_WRONG, errorcode.CURD_TYPE_IS_WRONG_MSG, response)
		return
	}
	//process接口对象
	classRoomProcessManager := new(process.ClassRoomProcess)
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()
	//返回参数 结构体(Code ErrMsg)
	result := new(bean.AddToSpeakAreaResponse) //反馈到客户端

	//内存中更新教室
	result.Code, result.ErrMsg = classRoomMemoryManager.DeleteAddToSpeakArea(request.ClassRoomId, request.CrudType, request.StudentIds)
	if result.Code == errorcode.SUCCESS {
		//数据库中更新教室 利用内存的数据更新
		result.Code, result.ErrMsg = classRoomProcessManager.DeleteAddToSpeakArea(request.ClassRoomId)
		//数据库操作成功
		if result.Code == errorcode.SUCCESS {
			speakerResponse_protocol := SendMessage_DeleteAddSpeaker(request, response)
			//发送消息失败
			if speakerResponse_protocol.Result != imbean.SUCCESS {
				writeErrMsg(speakerResponse_protocol.Result, speakerResponse_protocol.ErrorMessage, response)
				return
			}
		}
	}
	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("AddToSpeakArea  result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("AddToSpeakArea end")
	logs.GetLogger().Info("=============================================================\n")
}

//向IM发送通知
func SendMessage_DeleteAddSpeaker(request *bean.AddToSpeakAreaRequest, response http.ResponseWriter) *protocol.ResponseBase {
	logs.GetLogger().Info("SendMessage_DeleteAddSpeaker Is begin")

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
	deleteAddSpeakerClientProtocal := protocol.DeleteAddSpeakerToAreaClientProtocal{request.RequestUserId, request.StudentIds, request.CrudType, request.ClassRoomId}
	data, err := json.Marshal(deleteAddSpeakerClientProtocal)
	if err != nil {
		logs.GetLogger().Info("deleteAddSpeakerClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
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
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("DeleteAddSpeakerRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	DeleteAddSpeakerResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_DeleteAddSpeaker Handle iS error：", errOne)
		writeErrMsg(errorcode.HANDSFORBID_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &DeleteAddSpeakerResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("SendMessage_DeleteAddSpeaker Result IS：", DeleteAddSpeakerResponse_protocol)
	logs.GetLogger().Info("SendMessage_DeleteAddSpeaker Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_DeleteAddSpeaker Is end")

	return DeleteAddSpeakerResponse_protocol
}
