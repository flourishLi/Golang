package protocol

//禁止 解禁 禁言区 IM协议
type ForbidAreaProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type ForbidAreaClientProtocal struct {
	TeacherUserId int32   `json:"teacherUserId"`
	StudentIds    []int32 `json:"studentIds"`  //学生ID集合
	CrudType      int32   `json:"crudType"`    //1禁止  2 解禁
	ClassRoomId   int32   `json:"classRoomId"` //教室编号
}

func GetForbidAreaProtocol() *ForbidAreaProtocol {
	forbidAreaProtocol := &ForbidAreaProtocol{}
	forbidAreaProtocol.Command = MESSAGE_CMD_NOGROUP
	//设置命令编号
	forbidAreaProtocol.ProtocalCommand = FORBID_AREA
	forbidAreaProtocol.NotPushOffLine = true

	return forbidAreaProtocol
}
