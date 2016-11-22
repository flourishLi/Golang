package imbean

//重置用户密码的请求
//cmd=RESET_PASSWORD
type ResetUserPwdRequest struct {
	RequestBase
	NewPassword string `json:"password"` //新密码
	UserId      int32  `json:"userId"`   //
}

func GetResetUserPwdRequest() *ResetUserPwdRequest {
	request := new(ResetUserPwdRequest)
	request.Command = RESET_USERPWD
	return request
}

//重置用户密码的请求的响应
type ResetUserPwdResponse struct {
	ResponseBase
}

func (self ResetUserPwdRequest) GetCmd() string {
	return self.Command
}
