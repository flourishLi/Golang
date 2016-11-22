package classroom

import (
	"ZDTalk/ZDTalk_http/bean/resource"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetResourceList(response http.ResponseWriter, request *resource.ResourceListRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("GetResourceList begains")

	//LimitCount为空 直接返回
	if request.LimitCount == File_LimitCount_Is_NULL {
		writeErrMsg(errorcode.FILE_SEARCH_LIMIT_IS_NULL, errorcode.FILE_SEARCH_LIMIT_IS_NULL_MSG, response)
		return
	}

	//memory接口对象
	resourceMemoryManager := memory.GetUploadResourceMemoryManager()
	additionMemoryManager := memory.GetResourceAdditionMemoryManager()

	//客户端反馈结果
	result := new(resource.ResourceListResponse)
	result.ResourceList = make([]resource.ResourceInfo, 0)
	d := resourceMemoryManager.GetUploadResourceList(request.RequestUserId)
	//文件数量
	fileCount := int32(len(d))
	if request.StartSearchIndex < fileCount {
		//获取指定数量的文件
		var i int32 = 0
		var count int32 = 0 //计数器
		for i = fileCount - 1 - request.StartSearchIndex; i < fileCount; i-- {
			value := d[i]
			fileInfo := resource.ResourceInfo{}
			fileInfo.FileId = value.FileId
			fileInfo.FileTime = value.FileTime
			fileInfo.FileUrl = value.FilePath
			fileInfo.FileName = value.FileName
			//获取文件的子文件数量
			subFile := additionMemoryManager.GetResourceAddition(fileInfo.FileId)
			if subFile == nil {
				logs.GetLogger().Error("获取子文件数量时 父文件id出错")
			} else {
				fileInfo.FileContentCount = subFile.FileContentCount
			}
			result.ResourceList = append(result.ResourceList, fileInfo)
			//达到查询数量 跳出循环
			count++
			if count == request.LimitCount || count == fileCount-request.StartSearchIndex {
				break
			}
		}
	}
	result.Code = errorcode.SUCCESS

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("GetResourceList result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("GetResourceList end")
	logs.GetLogger().Info("=============================================================\n")
}
