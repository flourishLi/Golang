package protocol

//添加 移除学生到发言区 IM协议
type DeleteAddSpeakerToAreaProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type DeleteAddSpeakerToAreaClientProtocal struct {
	TeacherUserId int32   `json:"teacherUserId"`
	StudentIds    []int32 `json:"studentIds"`  //学生ID集合
	CrudType      int32   `json:"crudType"`    //1禁止  2 解禁
	ClassRoomId   int32   `json:"classRoomId"` //教室编号
}

func GetDeleteAddSpeakerToAreaProtocol() *DeleteAddSpeakerToAreaProtocol {
	addSpeakerToAreaProtocol := &DeleteAddSpeakerToAreaProtocol{}
	addSpeakerToAreaProtocol.Command = MESSAGE_CMD_NOGROUP

	//设置命令编号
	addSpeakerToAreaProtocol.ProtocalCommand = SAY_AREA_ADD_USER
	addSpeakerToAreaProtocol.NotPushOffLine = true

	return addSpeakerToAreaProtocol
}
