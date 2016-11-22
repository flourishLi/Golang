package beans

type GetClassRoomInfoRequest struct {
	ClientBaseRequest //CMD=GET_CLASS_ROOM_INFO
	ClassRoomName     string
}

type GetClassRoomInfoResponse struct {
	ClientBaseResponse     //code=1标识成功 3：此名字的教室已经存在
	RoomId             int //教室的ID
	ClassRoomName      string
	ClassRoomLogo      string
	Description        string
}
