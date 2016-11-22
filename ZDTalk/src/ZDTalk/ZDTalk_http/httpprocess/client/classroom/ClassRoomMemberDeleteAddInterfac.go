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
	"ZDTalk/utils/sliceutils"
	"ZDTalk/utils/timeutils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ClassRoomMemberDeleteAdd(response http.ResponseWriter, request *bean.DeleteAddStudentsRequest) {
	logs.GetLogger().Info("=============================================================")

	logs.GetLogger().Info("ClassRoomMemberDeleteAdd begins")

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
	//crudType值错误 必须为1 添加 2删除 直接返回
	if request.CrudType != CrudType_ADD && request.CrudType != CrudType_DELETE {
		writeErrMsg(errorcode.CURD_TYPE_IS_WRONG, errorcode.CURD_TYPE_IS_WRONG_MSG, response)
		return
	}
	//检查用户是否具有老师权限
	hasAuthority, isExit := AuthorityCheckout(request.RequestUserId, Teacher)
	if User_IS_NOT_EXIT == isExit {
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
	//Memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()
	//返回参数 结构体(Code ErrMsg)
	result := new(bean.DeleteAddStudentsResponse) //反馈到客户端
	//根据crudType IM处理
	isSuccess := CrudTypeIMHandle(request, response, request.CrudType)
	//IM 处理成功
	if isSuccess {
		//根据crudType 重新计算需要更新的memberlist
		memberList := GetMemberList_ByCrudType(response, classRoomMemoryManager, request.CrudType, request.ClassRoomId, request.StudentIds)
		if memberList != nil {
			//内存中更新教室
			result.Code, result.ErrMsg = classRoomMemoryManager.DeleteAddStudents(request.ClassRoomId, memberList)
			if result.Code == errorcode.SUCCESS {
				//数据库中更新教室 利用内存的数据更新
				result.Code, result.ErrMsg = classRoomProcessManager.DeleteAddStudents(request.ClassRoomId)
				//数据库操作成功
				if result.Code == errorcode.SUCCESS {
					//IM发送消息
					stdResponse_protocol := SendMessage_DeleteAddStudents(request, response)
					//发送消息失败
					if stdResponse_protocol.Result != imbean.SUCCESS {
						writeErrMsg(stdResponse_protocol.Result, stdResponse_protocol.ErrorMessage, response)
						return
					}
				}
			}
			datas, err1 := json.Marshal(result)
			if err1 != nil {
				logs.GetLogger().Error(err1.Error())
			}
			logs.GetLogger().Info("ClassRoomMemberDeleteAdd  result:" + string(datas))
			fmt.Fprintln(response, string(datas))

			logs.GetLogger().Info("ClassRoomMemberDeleteAdd end")
		}
	}
	logs.GetLogger().Info("=============================================================\n")
}

//根据crudType IM处理成员
//Param ClassRoomMemoryManager crudType,
func CrudTypeIMHandle(request *bean.DeleteAddStudentsRequest, response http.ResponseWriter, crudType int32) bool {
	if CrudType_ADD == crudType { //添加成员
		//调用IM添加群成员
		addGroupMemberResponse_IM := AddGroupMember_IM(request, response)
		if addGroupMemberResponse_IM.Result == imbean.SUCCESS { //IM处理成功
			logs.GetLogger().Info("AddGroupMember_IM Success")
			return true
		} else { //IM处理失败
			writeErrMsg(addGroupMemberResponse_IM.Result, addGroupMemberResponse_IM.ErrorMessage, response)
			return false
		}
	} else if CrudType_DELETE == crudType { //删除成员
		//调用IM删除群成员
		deleteGroupMemberResponse_IM := DeleteGroupMember_IM(request, response)
		if deleteGroupMemberResponse_IM.Result == imbean.SUCCESS { //IM处理成功
			logs.GetLogger().Info("deleteGroupMemberResponse_IM Success")
			return true
		} else { //IM处理失败
			writeErrMsg(deleteGroupMemberResponse_IM.Result, deleteGroupMemberResponse_IM.ErrorMessage, response)
			return false
		}
	} else {
		return false
	}
}

//在IM服务器 管理员删除群成员
func DeleteGroupMember_IM(request *bean.DeleteAddStudentsRequest, response http.ResponseWriter) *imbean.DeleteGroupMemberResponse {
	logs.GetLogger().Info("DeleteGroupMember_IM Is begin")
	//im的userid
	userIMid := GetUserIMId(request.RequestUserId)
	//获取教室id对应的ClassRoomIMid
	classroomIMId := GetClassRoomIMId(request.ClassRoomId)
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct base
	deleteGroupMemberRequest_IM := imbean.GetDeleteGroupMemberRequest()
	deleteGroupMemberRequest_IM.RequestServerId = 1
	deleteGroupMemberRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	deleteGroupMemberRequest_IM.RequestTime = timeutils.GetTimeStamp()
	deleteGroupMemberRequest_IM.SkipDBOperat = false
	deleteGroupMemberRequest_IM.MarkId = imbean.PUSH
	//client
	deleteGroupMemberRequest_IM.RequestUserId = userIMid
	deleteGroupMemberRequest_IM.ClassRoomIMId = classroomIMId
	for _, id := range request.StudentIds {
		imid := GetUserIMId(id)
		deleteGroupMemberRequest_IM.StudentIds += strconv.Itoa(int(imid)) + ","
	}
	logs.GetLogger().Info("deleteGroupMemberRequest_IM data is:", deleteGroupMemberRequest_IM)
	//json
	groupRequest_Bytes, err := json.Marshal(deleteGroupMemberRequest_IM)
	if err != nil {
		logs.GetLogger().Info("deleteGroupMemberRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("deleteGroupMemberRequest_IM to json is:", string(groupRequest_Bytes))
	//response
	deleteGroupMemberResponse_IM := new(imbean.DeleteGroupMemberResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMHttpGroupUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("deleteGroupMember_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.DELETE_GROUP_MEMBER_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &deleteGroupMemberResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("DeleteGroupMember_IM Result IS：", string(data))
	logs.GetLogger().Info("DeleteGroupMember_IM Is end")
	return deleteGroupMemberResponse_IM
}

//在IM服务器 成员自动退群
func QuitGroup_IM(request *bean.DeleteAddStudentsRequest, response http.ResponseWriter) *imbean.QuitGroupResponse {
	logs.GetLogger().Info("QuitGroup_IM Is begin")

	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct base
	quiteGroupRequest_IM := imbean.GetQuitGroupRequest()
	quiteGroupRequest_IM.RequestServerId = 1
	quiteGroupRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	quiteGroupRequest_IM.RequestTime = timeutils.GetTimeStamp()
	quiteGroupRequest_IM.SkipDBOperat = false
	quiteGroupRequest_IM.MarkId = imbean.PUSH
	//client
	quiteGroupRequest_IM.NewManagerId = 1 //被改为管理员的群成员ID

	logs.GetLogger().Info("quiteGroupRequest_IM data is:", quiteGroupRequest_IM)
	//json
	groupRequest_Bytes, err := json.Marshal(quiteGroupRequest_IM)
	if err != nil {
		logs.GetLogger().Info("quiteGroupRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("quiteGroupRequest_IM to json is:", string(groupRequest_Bytes))
	//response
	quiteGroupMemberResponse_IM := new(imbean.QuitGroupResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMHttpGroupUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("quiteGroup_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.QUITE_GROUP_MEMBER_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &quiteGroupMemberResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("QuitGroup_IM Result IS：", quiteGroupMemberResponse_IM)
	logs.GetLogger().Info("QuitGroup_IM Result IS：", string(data))

	logs.GetLogger().Info("QuitGroup_IM Is end")

	return quiteGroupMemberResponse_IM
}

//在IM服务器 添加群成员
func AddGroupMember_IM(request *bean.DeleteAddStudentsRequest, response http.ResponseWriter) *imbean.AddGroupMemberResponse {
	logs.GetLogger().Info("AddGroupMember_IM Is begin")
	//获取教室id对应的ClassRoomIMid
	classroomIMId := GetClassRoomIMId(request.ClassRoomId)
	userIMid := GetUserIMId(request.RequestUserId)
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct base
	addGroupMemberRequest_IM := imbean.GetAddGroupMemberRequest()
	addGroupMemberRequest_IM.RequestServerId = 1 //请求服务器的ID
	addGroupMemberRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	addGroupMemberRequest_IM.RequestTime = timeutils.GetTimeStamp()
	addGroupMemberRequest_IM.SkipDBOperat = false
	addGroupMemberRequest_IM.MarkId = imbean.PUSH //0=推送，1=不推送
	//client
	addGroupMemberRequest_IM.IsDisGroup = imbean.GROUPTYPE
	addGroupMemberRequest_IM.ClassRoomIMId = classroomIMId
	addGroupMemberRequest_IM.RequestType = imbean.ENTERGROUP_BYUSER //1是管理员邀请用户入群；2是用户主动申请入群
	addGroupMemberRequest_IM.RequestUserId = userIMid               //requestType为1时，表示邀请者(管理员)的id；requestType为2时，忽略该字段
	for _, id := range request.StudentIds {
		imid := GetUserIMId(id)
		addGroupMemberRequest_IM.StudentIds += strconv.Itoa(int(imid)) + ","
	}
	addGroupMemberRequest_IM.IsDisGroup = imbean.GROUPTYPE

	logs.GetLogger().Info("addGroupMemberRequest_IM data is:", addGroupMemberRequest_IM)
	//json
	groupRequest_Bytes, err := json.Marshal(addGroupMemberRequest_IM)
	if err != nil {
		logs.GetLogger().Info("addGroupMemberRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("addGroupMemberRequest_IM to json is:", string(groupRequest_Bytes))
	//response
	addGroupMemberResponse_IM := new(imbean.AddGroupMemberResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMHttpGroupUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("quiteGroup_IM Handle iS error：", errOne)
		writeErrMsg(errorcode.QUITE_GROUP_MEMBER_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &addGroupMemberResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("AddGroupMember_IM Result IS：", addGroupMemberResponse_IM)
	logs.GetLogger().Info("AddGroupMember_IM Result IS：", string(data))
	logs.GetLogger().Info("AddGroupMember_IM Is end")
	return addGroupMemberResponse_IM
}

//根据crudType创建需要更新的用户列表
//Param ClassRoomMemoryManager crudType, classRoomID userList
//return memberList
func GetMemberList_ByCrudType(response http.ResponseWriter, classRoomMemoryManager *memory.ClassRoomMemoryManager, crudType, classRoomID int32, userIdList []int32) []int32 {
	memberList := make([]int32, 0)
	if r, ok := classRoomMemoryManager.ClassRooms[classRoomID]; ok {
		//获取当前教室的成员列表
		memberList = r.MemberList
		logs.GetLogger().Info("OriginMemory MemberList Is:", memberList)
		//移除教室
		if CrudType_DELETE == crudType {
			//遍历要移除的学生列表成员
			for _, userID := range userIdList {
				//用户是否在当前列表中
				isExit, err := sliceutils.Containts(memberList, userID)
				if err != nil {
					logs.GetLogger().Error("sliceutils Containts err", err)
				}
				//在 移除
				if isExit {
					memberList = sliceutils.RemoveInt32(memberList, userID)
				} else {
					logs.GetLogger().Info("UserID Is Not Exit:", userID)
				}
			}
		} else if CrudType_ADD == crudType { //添加到教室
			//遍历要添加的学生列表成员
			for _, userID := range userIdList {
				//用户是否在当前列表中
				isExit, err := sliceutils.Containts(memberList, userID)
				if err != nil {
					logs.GetLogger().Error("sliceutils Containts err", err)
				}
				//不在 添加
				if !isExit {
					memberList = append(memberList, userID)
				} else {
					logs.GetLogger().Info("UserID is already Exit:", userID)
				}
			}
		}
		return memberList
	} else { //教室不存在
		writeErrMsg(errorcode.CLASS_ROOM_IS_NOT_EXIST, errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG, response)
		return nil
	}
}

//向IM发送通知
func SendMessage_DeleteAddStudents(request *bean.DeleteAddStudentsRequest, response http.ResponseWriter) *protocol.ResponseBase {
	logs.GetLogger().Info("SendMessage_DeleteAddStudents Is begin")
	//获取用户id对应的IMid
	userIMId := GetUserIMId(request.RequestUserId)
	//获取时间戳
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))

	//获取该教室的在线成员列表
	userMemoryManager := memory.GetUserInfoMemoryManager()
	onLineUsers := userMemoryManager.GetOnlineUserList(request.ClassRoomId)

	//获取对应的Imuserid集合
	imOnlineUsers := []int32{}
	for _, userInfo := range onLineUsers {
		imOnlineUsers = append(imOnlineUsers, userInfo.ChatId)
	}

	//构建要发送给IM的客户端数据 作为最终IM协议的Content
	deleteAddStudendtClientProtocal := protocol.DeleteAddClientProtocal{request.RequestUserId, request.StudentIds, request.CrudType, request.ClassRoomId}
	data, err := json.Marshal(deleteAddStudendtClientProtocal)
	if err != nil {
		logs.GetLogger().Info("deleteAddStudendtClientProtocal Marshal err:", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	//构建最终发送给IM的协议
	//requestBase设置
	DeleteAddStdRequest_protocol := protocol.GetDeleteAddStudentProtocol()
	DeleteAddStdRequest_protocol.RequestServerId = 1
	DeleteAddStdRequest_protocol.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	DeleteAddStdRequest_protocol.RequestTime = timeutils.GetTimeStamp()
	DeleteAddStdRequest_protocol.SkipDBOperat = false
	DeleteAddStdRequest_protocol.MarkId = imbean.PUSH
	//request其他设置
	DeleteAddStdRequest_protocol.SenderUserId = userIMId
	DeleteAddStdRequest_protocol.UserIds = imOnlineUsers
	DeleteAddStdRequest_protocol.SendContent = string(cryptoutils.Base64Encode(data))

	logs.GetLogger().Info("DeleteAddStdRequest_protocol data is:", DeleteAddStdRequest_protocol)

	//handsForbidRequest_protocol转成byte
	groupRequest_Bytes, err := json.Marshal(DeleteAddStdRequest_protocol)
	if err != nil {
		logs.GetLogger().Info("DeleteAddStdRequest_protocol to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("DeleteAddStdRequest_protocol to json is:", string(groupRequest_Bytes))
	//response
	DeleteAddStdResponse_protocol := new(protocol.ResponseBase)

	//向IM请求数据
	responseData, errOne := httputil.PostData(config.ConfigNodeInfo.IMMessageUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("SendMessage_DeleteAddStudents Handle iS error：", errOne)
		writeErrMsg(errorcode.HANDSFORBID_PROTOCOL_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(responseData, &DeleteAddStdResponse_protocol)
	if errTwo != nil {
		logs.GetLogger().Error(errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)

		return nil
	}
	logs.GetLogger().Info("SendMessage_DeleteAddStudents Result IS：", DeleteAddStdResponse_protocol)
	logs.GetLogger().Info("SendMessage_DeleteAddStudents Result IS：", string(responseData))
	logs.GetLogger().Info("SendMessage_DeleteAddStudents Is end")

	return DeleteAddStdResponse_protocol
}
