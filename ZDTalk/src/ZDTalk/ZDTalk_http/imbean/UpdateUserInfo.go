// CreatGroupRequest
package imbean

//修改用户信息的请求
//cmd=UPDATE_USER_INFO
type UpdateUserInfoRequest struct {
	RequestBase
	UserName string `json:"userName"` //用户昵称 必须
	UserIcon string `json:"userIcon"` //头像
	UserId   int32  `json:"userId"`   //指定用户的userId(可选，如果此字段为空，系统会返回一个userId)
	OtherId  int32  `json:"otherId"`  //在业务服务器中的id(可选)
}

func GetUpdateUserInfoRequest() *UpdateUserInfoRequest {
	request := new(UpdateUserInfoRequest)
	request.Command = UPDATE_USERINFO
	return request
}

//修改用户信息的请求的响应
type UpdateUserInfoResponse struct {
	ResponseBase
}

func (self UpdateUserInfoRequest) GetCmd() string {
	return self.Command
}
