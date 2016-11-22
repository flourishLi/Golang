package bean

type UserInfoUpdateRequest struct {
	ClientBaseRequest        //CMD=USERINFO_UPDATE 含有UserId
	UserIcon          string `json:"userIcon"`
	UserName          string `json:"userName"`
	Role              int32  `json:"role"`       //角色1学生 3老师 必须
	DeviceType        int32  `json:"deviceType"` //	1 Android 2 Iphone 3 PC
}

type UserInfoUpdateResponse struct {
	ClientBaseResponse
}
