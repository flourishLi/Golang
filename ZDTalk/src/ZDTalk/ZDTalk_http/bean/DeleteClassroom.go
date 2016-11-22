package bean

type DeleteClassroomRequest struct {
	ClientBaseRequest       //CMD=DELETE_CLASS_ROOM
	ClassRoomId       int32 `json:"classRoomId"`
}

type DeleteClassroomResponse struct {
	ClientBaseResponse //code=1标识成功3：此教室已删除
}
