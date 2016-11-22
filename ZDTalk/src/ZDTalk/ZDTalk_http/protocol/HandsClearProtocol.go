package protocol

//清空举手列表 IM协议
type HandsClearProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type HandsClearClientProtocal struct {
	TeacherUserId int32 `json:"teacherUserId"`
	ClassRoomId   int32 `json:"classRoomId"` //教室编号
}

func GetHandsClearProtocol() *HandsClearProtocol {
	handsClearProtocol := &HandsClearProtocol{}
	handsClearProtocol.Command = MESSAGE_CMD_NOGROUP

	//设置命令编号
	handsClearProtocol.ProtocalCommand = HANDS_CLEAR
	handsClearProtocol.NotPushOffLine = true

	return handsClearProtocol
}
