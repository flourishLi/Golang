package bean

type GetUserInfoRequest struct {
	ClientBaseRequest       //CMD=USERINFO_UPDATE 含有UserId
	UserId            int32 `json:"userId"` //待查用户的id
}

type GetUserInfoResponse struct {
	ClientBaseResponse
	Data UserInfo `json:"data"`
}
