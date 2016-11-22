package classroom

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetHandsUserList(response http.ResponseWriter, request *bean.GetHandsUserListRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("GetHandsUserList begins")

	//ClassRoomId为空 直接返回
	if request.ClassRoomId == ClassRoomID_IS_NULL {
		writeErrMsg(errorcode.CLASS_ROOM_ID_ERROR, errorcode.CLASS_ROOM_ID_ERROR_MSG, response)
		return
	}

	//memory接口对象
	userInfoMemoryManager := memory.GetUserInfoMemoryManager()
	//返回参数 结构体(Code ErrMsg)
	result := new(bean.GetHandsUserListResponse) //反馈到客户端

	//获取举手列表
	userinfoLsit := userInfoMemoryManager.GetHandsUserList(request.ClassRoomId)

	if userinfoLsit != nil {
		result.Code = errorcode.SUCCESS
		result.ErrMsg = ""
	} else {
		result.Code = errorcode.CLASS_ROOM_IS_NOT_EXIST
		result.ErrMsg = errorcode.CLASS_ROOM_IS_NOT_EXIST_MSG
	}
	for _, v := range userinfoLsit {
		i := bean.HandUserInfo{}
		i.UserId = v.UserId
		i.ChatId = v.ChatId
		i.Role = v.Role
		i.UserName = v.UserName
		i.UserIcon = v.UserIcon
		i.HandTime = v.HandTime
		result.UserList = append(result.UserList, i)
	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("GetHandsUserList  result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("GetHandsUserList end")
	logs.GetLogger().Info("=============================================================\n")
}
