package bean

type UserSignUpRequest struct {
	ClientBaseRequest        //CMD=SIGN_UP
	UserIcon          string `json:"userIcon"`
	UserName          string `json:"userName"`
	Role              int32  `json:"role"` //角色1学生 3老师 必须
	LoginName         string `json:"loginName"`
	Password          string `json:"password"`
	DeviceType        int32  `json:"deviceType"` //	1 Android 2 Iphone 3 PC
}

type UserSignUpResponse struct {
	ClientBaseResponse
	UserId   int32  `json:"userId"`
	ChatId   int32  `json:"chatId"`
	UserIcon string `json:"userIcon"`
	UserName string `json:"userName"`
	Role     int32  `json:"role"`
}
