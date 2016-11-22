package beans

type CreateClassroomRequest struct {
	ClientBaseRequest //CMD=CREATE_CLASS_ROOM
	ClassRoomName     string
	ClassRoomLogo     string
	Description       string //教室说明
}

type CreateClassroomResponse struct {
	ClientBaseResponse     //code=1标识成功 3：此名字的教室已经存在
	RoomId             int //教室的ID
}
