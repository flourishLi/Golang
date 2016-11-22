package classroom

import (
	"ZDTalk/ZDTalk_http/bean/resource"
	"ZDTalk/ZDTalk_http/httpprocess/process"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
	"net/http"
)

func DeleteResource(response http.ResponseWriter, request *resource.DeleteResourceRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("DeleteResource begins")

	//fileId为空 直接返回
	if request.FileId == 0 {
		writeErrMsg(errorcode.FILE_ID_IS_NOT_EXIST, errorcode.FILE_ID_IS_NOT_EXIST_MSG, response)
		return
	}
	//请求用户的Id为空
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
	classRoomProcessManager := new(process.ClassRoomProcess)
	//memory接口对象
	resourceMemoryManager := memory.GetUploadResourceMemoryManager()

	//返回参数 结构体(Code ErrMsg)
	result := new(resource.DeleteResourceResponse) //反馈到客户端

	//内存中删除资源
	result.Code, result.ErrMsg = resourceMemoryManager.DeleteResource(request.RequestUserId, request.FileId)
	if result.Code == errorcode.SUCCESS {
		//数据库中删除上传文件
		result.Code, result.ErrMsg = classRoomProcessManager.DeleteResource(request.FileId)
	}

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("DeleteResource result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("DeleteResource end")
	logs.GetLogger().Info("=============================================================\n")
}
