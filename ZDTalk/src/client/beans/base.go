package beans

//客户端请求基类
type ClientBaseRequest struct {
	Command string `json:"CMD"` //命令编号
}

//客户端响应基类
type ClientBaseResponse struct {
	Code   int    `json:"code"`   //响应编码
	ErrMsg string `json:"errMsg"` //错误信息
}
