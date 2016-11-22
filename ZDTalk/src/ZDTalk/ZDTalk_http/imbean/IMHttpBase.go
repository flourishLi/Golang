package imbean

type RequestBase struct {
	Command            string `json:"cmd"`
	SkipDBOperat       bool   `json:"skipDBOperat"`    //是否是忽略数据库操作，在业务服务器处理数据库操作，IM服务器仅操作内存数据的时候使用，默认为false
	RequestTime        int32  `json:"requestTime"`     //以秒为单位的时间戳
	RequestServerId    int32  `json:"requestServerId"` //请求服务器的ID
	RequestServerToken string `json:"requestServerToken"`
	MarkId             int32  `json:"markId"` // 0 =推送 1=不推送
}

func (self RequestBase) GetCmd() string {
	return self.Command
}

type BaseInterface interface {
	GetCmd() string
}

type ResponseBase struct {
	Result       int32  `json:"code"`
	ErrorMessage string `json:"errMsg"`
	ResponseTime int64  `json:"responseTime"`
}

func CreateResponseBase() *ResponseBase {
	return &ResponseBase{}
}
