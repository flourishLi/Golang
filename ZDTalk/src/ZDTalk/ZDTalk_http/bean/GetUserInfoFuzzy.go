package bean

type GetUserInfoFuzzyRequest struct {
	ClientBaseRequest        //CMD=USERINFO_UPDATE 含有UserId
	UserName          string `json:"userName"` //用户昵称
}

type GetUserInfoFuzzyResponse struct {
	ClientBaseResponse
	Data []UserInfo `json:"data"`
}
