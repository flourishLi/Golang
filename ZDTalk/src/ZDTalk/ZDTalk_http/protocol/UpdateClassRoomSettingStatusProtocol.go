package protocol

//教室设置 IM协议
type UpdateClassRoomSettingStatusProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type UpdateClassRoomSettingStatusClientProtocal struct {
	TeacherUserId int32   `json:"teacherUserId"`
	ClassRoomId   int32   `json:"classRoomId"`   //教室编号
	SettingStatus []int32 `json:"settingStatus"` //1 允许所有人打字，2 允许学生录音，3 禁止游客打字，4 禁止游客举手

}

func GetUpdateClassRoomSettingStatusProtocol() *UpdateClassRoomSettingStatusProtocol {
	updateClassRoomSettingStatusProtocol := &UpdateClassRoomSettingStatusProtocol{}
	updateClassRoomSettingStatusProtocol.Command = MESSAGE_CMD_NOGROUP
	//设置命令编号
	updateClassRoomSettingStatusProtocol.ProtocalCommand = START_FINISH_CLASS
	updateClassRoomSettingStatusProtocol.NotPushOffLine = true

	return updateClassRoomSettingStatusProtocol
}
