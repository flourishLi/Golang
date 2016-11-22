package bean

type GetOnLineUserListRequest struct {
	ClientBaseRequest       //CMD=GET_ONLINE_USER_LIST
	ClassRoomId       int32 `json:"classRoomId"` //教室ID
}

type GetOnLineUserListResponse struct {
	ClientBaseResponse
	//userinfo
	UserList []UserInfo `json:"userList"`
}
