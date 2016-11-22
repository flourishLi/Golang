package bean

type GetClassRoomListRequest struct {
	ClientBaseRequest //CMD=GET_CLASS_ROOM_LIST
}

type ClassRoomInfo struct {
	ClassRoomName string `json:"classRoomName"`
	ClassRoomLogo string `json:"classRoomLogo"`
	Description   string `json:"description"`
	ClassRoomId   int32  `json:"classRoomId"`
}

type GetClassRoomListResponse struct {
	ClientBaseResponse                  //code=1标识成功 3：此名字的教室已经存在
	Rooms              []*ClassRoomInfo `json:"rooms"`
}
