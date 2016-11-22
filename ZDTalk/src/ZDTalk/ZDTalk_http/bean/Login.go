package bean

type LoginRequest struct {
	ClientBaseRequest        //CMD=LOGIN
	LoginName         string `json:"loginName"`
	Password          string `json:"password"`
}

type LoginResponse struct {
	ClientBaseResponse
	Data *UserInfo `json:"data"`
}
