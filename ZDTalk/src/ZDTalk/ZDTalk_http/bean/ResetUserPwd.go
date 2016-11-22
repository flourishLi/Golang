package bean

type UserPwdResetRequest struct {
	ClientBaseRequest        //CMD=USERINFO_UPDATE 含有UserId
	NewPassword       string `json:"newPassword"`
}

type UserPwdResetResponse struct {
	ClientBaseResponse
}
