package imbean

type UserInfo struct {
	UserId   int32  `json:"userId"`
	ChatId   int32  `json:"chatId"`
	UserIcon string `json:"userIcon"`
	UserName string `json:"userName"`
	Role     int32  `json:"role"`
}
