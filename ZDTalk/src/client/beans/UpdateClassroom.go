package beans

type UpdateClassroomRequest struct {
	ClientBaseRequest //CMD=CREATE_CLASS_ROOM
	RoomId            int
	ClassRoomLogo     string
	Description       string //教室说明
}

type UpdateClassroomResponse struct {
	ClientBaseResponse //code=1标识成功 3：此名字的教室已经存在

}
