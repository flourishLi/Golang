package protocol

//退出教室 IM协议
type ExitClassRoomProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type ExitClassRoomClientProtocal struct {
	UserId      int32 `json:"userId"`
	ClassRoomId int32 `json:"classRoomId"` //教室编号
}

func GetExitClassRoomProtocol() *ExitClassRoomProtocol {
	exitClassRoomProtocol := &ExitClassRoomProtocol{}
	exitClassRoomProtocol.Command = MESSAGE_CMD_NOGROUP

	//设置命令编号
	exitClassRoomProtocol.ProtocalCommand = EXIT_CLASS_ROOM
	exitClassRoomProtocol.NotPushOffLine = true

	return exitClassRoomProtocol
}
