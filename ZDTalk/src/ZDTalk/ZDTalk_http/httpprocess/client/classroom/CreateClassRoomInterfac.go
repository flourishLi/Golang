package classroom

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/ZDTalk_http/httpprocess/process"
	"ZDTalk/ZDTalk_http/imbean"
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

func CreateClassroom(response http.ResponseWriter, request *bean.CreateClassroomRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("create classroom begins")
	//教室名称为空
	if request.ClassRoomName == "" {
		writeErrMsg(errorcode.CLASS_ROOM_NAME_CAN_NOT_NULL, errorcode.CLASS_ROOM_NAME_CAN_NOT_NULL_MSG, response)
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
	//Process接口对象
	classRoomProcessManager := new(process.ClassRoomProcess)
	//Memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//系统时间
	creatTime := time.Now().Unix()
	//教室对应IM的Id 需要调用IM服务器接口获取 待定
	var classRoomIMId int32 = 0
	//返回参数 结构体(Code ErrMsg RoomId)
	result := new(bean.CreateClassroomResponse) //反馈到客户端

	//调用IM同步数据 创建群
	createGroupResponse_IM := CreateGroup_IM(request, response)
	//IM 创建群成功
	if createGroupResponse_IM.Result == imbean.SUCCESS {
		classRoomIMId = createGroupResponse_IM.ClassRoomIMId
		logs.GetLogger().Info("CreateGroup_IM Is Sucess")
		//数据库中创建教室
		//Param classRoomName教室名称 classRoomLogo教室头像 description教室说明 ClassRoomCourse课程信息 CreatorUserId创建者ID CreateTime创建时间 ClassRoomIMId 教室对应IM的Id
		result.Code, result.ErrMsg, result.ClassRoomId = classRoomProcessManager.CreateClassRoom(request.ClassRoomName, request.ClassRoomLogo, request.Description, request.ClassRoomCourse, request.RequestUserId, request.SettingStatus, creatTime, classRoomIMId)
		if result.Code == errorcode.SUCCESS {
			//内存中创建教室
			//Param classRoomID(数据库创建成功后的) classRoomName教室名称 classRoomLogo教室头像 description教室说明 ClassRoomCourse课程信息 CreatorUserId创建者ID CreateTime创建时间 ClassRoomIMId 教室对应IM的Id
			result.Code, result.ErrMsg, result.ClassRoomId = classRoomMemoryManager.CreateClassroom(result.ClassRoomId, request.ClassRoomLogo, request.ClassRoomName, request.Description, request.ClassRoomCourse, request.RequestUserId, creatTime, classRoomIMId, request.SettingStatus)
		}

		//返回参数转化成Json数据 返回到客户端
		datas, errThree := json.Marshal(result)
		if errThree != nil {
			logs.GetLogger().Error("Json Parse Error：", errThree.Error())
		}
		logs.GetLogger().Info("create classroom result:" + string(datas))
		fmt.Fprintln(response, string(datas))
		logs.GetLogger().Info("create classroom end")
	} else { //IM 创建群失败
		writeErrMsg(createGroupResponse_IM.Result, createGroupResponse_IM.ErrorMessage, response)
		return
	}
	logs.GetLogger().Info("=============================================================\n")
}

//在IM服务器创建群
//return CreateGroupResponse
func CreateGroup_IM(request *bean.CreateClassroomRequest, response http.ResponseWriter) *imbean.CreateGroupResponse {
	logs.GetLogger().Info("CreateGroup_IM Is begin")

	userIMid := GetUserIMId(request.RequestUserId)
	timeString := strconv.Itoa(int(timeutils.GetTimeStamp()))
	//request struct
	createGroupRequest_IM := imbean.GetCreateGroupRequest()
	createGroupRequest_IM.RequestServerId = 1
	createGroupRequest_IM.RequestServerToken = cryptoutils.Md5Encode("key" + timeString)
	createGroupRequest_IM.RequestTime = timeutils.GetTimeStamp()
	createGroupRequest_IM.SkipDBOperat = false
	createGroupRequest_IM.MarkId = imbean.PUSH

	createGroupRequest_IM.ClassRoomName = request.ClassRoomName
	createGroupRequest_IM.ClassRoomLogo = request.ClassRoomLogo
	createGroupRequest_IM.Description = request.Description
	createGroupRequest_IM.RequestUserId = userIMid
	createGroupRequest_IM.IsDisGroup = imbean.GROUPTYPE //1=讨论组，2=群 3聊天室
	logs.GetLogger().Info("createGroupRequest_IM data is:", createGroupRequest_IM)
	//json
	groupRequest_Bytes, err := json.Marshal(createGroupRequest_IM)
	if err != nil {
		logs.GetLogger().Info("createGroupRequest_IM to json is Wrong", err)
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("createGroup_IMRequest_IM to json is:", string(groupRequest_Bytes))
	//response
	createGroupResponse_IM := new(imbean.CreateGroupResponse)

	//向IM请求数据
	data, errOne := httputil.PostData(config.ConfigNodeInfo.IMHttpGroupUrl, groupRequest_Bytes)
	if errOne != nil {
		logs.GetLogger().Info("IM Handle iS error：", errOne)
		writeErrMsg(errorcode.CREATE_GROUP_IM_IS_ERROR, err.Error(), response)
		return nil
	}

	//返回参数转化成struct数据 返回到客户端
	errTwo := json.Unmarshal(data, &createGroupResponse_IM)
	if errTwo != nil {
		logs.GetLogger().Error("createGroupResponse_IM to json is err:", errTwo.Error())
		writeErrMsg(errorcode.JSON_PRASEM_ERROR, err.Error(), response)
		return nil
	}
	logs.GetLogger().Info("createGroup_IMRequest_IM Result IS：", createGroupResponse_IM)
	logs.GetLogger().Info("createGroup_IMRequest_IM Result IS：", string(data))

	logs.GetLogger().Info("createGroup_IMRequest_IM Is end")

	return createGroupResponse_IM
}
