package protocol

//移除 添加学生到教室 IM协议
type DeleteAddStudentProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type DeleteAddClientProtocal struct {
	TeacherUserId int32   `json:"teacherUserId"`
	StudentIds    []int32 `json:"studentIds"`  //学生ID集合
	CrudType      int32   `json:"crudType"`    //1禁止  2 解禁
	ClassRoomId   int32   `json:"classRoomId"` //教室编号
}

func GetDeleteAddStudentProtocol() *DeleteAddStudentProtocol {
	deleteAddStudentProtocol := &DeleteAddStudentProtocol{}
	deleteAddStudentProtocol.Command = MESSAGE_CMD_NOGROUP

	//设置命令编号
	deleteAddStudentProtocol.ProtocalCommand = CLASS_ROOM_MEMBER_DELETE
	deleteAddStudentProtocol.NotPushOffLine = true

	return deleteAddStudentProtocol
}
