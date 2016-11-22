package protocol

//老师上下课IM协议
type UpdateClassRoomStatusProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type UpdateRoomStatusClientProtocal struct {
	TeacherUserId int32 `json:"teacherUserId"`
	ClassRoomId   int32 `json:"classRoomId"` //教室编号
	RoomStatus    int32 `json:"roomStatus"`
}

func GetUpdateClassRoomStatusProtocol() *UpdateClassRoomStatusProtocol {
	updateClassRoomStatusProtocol := &UpdateClassRoomStatusProtocol{}
	updateClassRoomStatusProtocol.Command = MESSAGE_CMD_NOGROUP

	//设置命令编号
	updateClassRoomStatusProtocol.ProtocalCommand = ROOM_STATUS
	updateClassRoomStatusProtocol.NotPushOffLine = true

	return updateClassRoomStatusProtocol
}
