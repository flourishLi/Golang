package bean

//客户端请求基类
type ClientBaseRequest struct {
	Command       string `json:"cmd"`           //命令编号
	RequestUserId int32  `json:"requestUserId"` //请求者用户Id
	Token         string `json:"token"`         //Token
	RequestTime   int64  `json:"requestTime"`   //请求时间
}

//客户端响应基类
type ClientBaseResponse struct {
	Code   int32  `json:"code"`   //响应编码
	ErrMsg string `json:"errMsg"` //错误信息
}
