package resource

//获取上传文件的列表

import (
	"ZDTalk/ZDTalk_http/bean"
)

type ResourceListRequest struct {
	bean.ClientBaseRequest       //CMD=GET_RESOURCE_LIST
	StartSearchIndex       int32 `json:"startSearchIndex"` //获取文件列表的起始编号0 1 2 3 ...
	LimitCount             int32 `json:"limitCount"`       //列表数据的数量
}

type ResourceInfo struct {
	FileId           int32  `json:"fileId"`           //文件Id
	FileUrl          string `json:"fileUrl"`          //文件网络路径
	FileTime         int64  `json:"fileTime"`         //文件上传时间
	FileName         string `json:"fileName"`         //文件上传时间
	FileContentCount int32  `json:"fileContentCount"` //子文件数量

}

type ResourceListResponse struct {
	bean.ClientBaseResponse
	ResourceList []ResourceInfo `json:"resourceList"`
}
