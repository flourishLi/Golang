package imbean

//修改用户密码的请求
//cmd=MODIFY_PASSWORD
type UpdateUserPwdRequest struct {
	RequestBase
	NewPassword string `json:"password"`    //新密码
	OldPassword string `json:"oldPassword"` //老密码
	UserId      int32  `json:"userId"`      //
}

func GetUpdateUserPwdRequest() *UpdateUserPwdRequest {
	request := new(UpdateUserPwdRequest)
	request.Command = UPDATE_USERPWD
	return request
}

//修改用户密码的请求的响应
type UpdateUserPwdResponse struct {
	ResponseBase
}

func (self UpdateUserPwdRequest) GetCmd() string {
	return self.Command
}
