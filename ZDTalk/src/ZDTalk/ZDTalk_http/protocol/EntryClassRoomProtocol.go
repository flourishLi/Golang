package protocol

//学生进入教室 IM协议
type EntryClassRoomProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type EntryClassRoomClientProtocal struct {
	UserId      int32 `json:"userId"`
	ClassRoomId int32 `json:"classRoomId"` //教室编号
}

func GetEntryClassRoomProtocol() *EntryClassRoomProtocol {
	entryClassRoomProtocol := &EntryClassRoomProtocol{}
	entryClassRoomProtocol.Command = MESSAGE_CMD_NOGROUP

	//设置命令编号
	entryClassRoomProtocol.ProtocalCommand = STUDENT_ENTRY_CLASSROOM
	entryClassRoomProtocol.NotPushOffLine = true

	return entryClassRoomProtocol
}
