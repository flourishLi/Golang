package bean

type UploadResourceRequest struct {
	ClientBaseRequest        //CMD=UPLOAD_RESOURCE
	FileUrl           string `json:"fileUrl"`          //文件网络路径
	FileName          string `json:"fileName"`         //被指定的文件名称
	FileContentCount  int32  `json:"fileContentCount"` //被指定的文件名称
}

type UploadResourceResponse struct {
	ClientBaseResponse
	FileId int32 `json:"fileId"`
	//	Url    string `json:"url"`
}
