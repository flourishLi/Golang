package protocol

//老师上下课IM协议
type UpdateHandForbidStatusProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type UpdateHandForbidStatusClientProtocal struct {
	TeacherUserId    int32 `json:"teacherUserId"`
	ClassRoomId      int32 `json:"classRoomId"` //教室编号
	ForbidHandStatus int32 `json:"forbidHandStatus"`
}

func GetUpdateHandForbidStatusProtocol() *UpdateHandForbidStatusProtocol {
	updateHandForbidStatusProtocol := &UpdateHandForbidStatusProtocol{}
	updateHandForbidStatusProtocol.Command = MESSAGE_CMD_NOGROUP

	//设置命令编号
	updateHandForbidStatusProtocol.ProtocalCommand = HAND_FORBID_STATUS
	updateHandForbidStatusProtocol.NotPushOffLine = true

	return updateHandForbidStatusProtocol
}
