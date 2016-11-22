package bean

type UserPwdUpdateRequest struct {
	ClientBaseRequest        //CMD=USERINFO_UPDATE 含有UserId
	OldPassword       string `json:"oldPassword"`
	NewPassword       string `json:"newPassword"`
}

type UserPwdUpdateResponse struct {
	ClientBaseResponse
}
