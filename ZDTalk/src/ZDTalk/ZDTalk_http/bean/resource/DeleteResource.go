package resource

import (
	"ZDTalk/ZDTalk_http/bean"
)

type DeleteResourceRequest struct {
	bean.ClientBaseRequest       //CMD=UPLOAD_RESOURCE
	FileId                 int32 `json:"fileId"`
}

type DeleteResourceResponse struct {
	bean.ClientBaseResponse
}
