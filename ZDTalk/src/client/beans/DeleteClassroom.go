package beans

type DeleteClassroomRequest struct {
	ClientBaseRequest //CMD=CREATE_CLASS_ROOM
	RoomId            int
}

type DeleteClassroomResponse struct {
	ClientBaseResponse //code=1标识成功 3：此名字的教室已经存在

}
