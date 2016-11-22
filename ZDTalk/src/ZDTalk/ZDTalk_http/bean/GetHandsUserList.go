package bean

type GetHandsUserListRequest struct {
	ClientBaseRequest       //CMD=GET_HANDS_USER_LIST
	ClassRoomId       int32 `json:"classRoomId"` //创建人
}

type GetHandsUserListResponse struct {
	ClientBaseResponse
	//userinfo
	UserList []HandUserInfo `json:"userList"`
}

type HandUserInfo struct {
	UserId   int32  `json:"userId"`
	ChatId   int32  `json:"chatId"` //用户对应的IM中的userId
	Role     int32  `json:"role"`
	UserName string `json:"userName"`
	UserIcon string `json:"userIcon"`
	HandTime int64  `json:"handTime"` //举手时间
}
