package protocol

//举手 取消请求 IM协议
type HandsUpProtocol struct {
	ProtocolShell
}

//来自客户端的数据
type HandsUpClientProtocal struct {
	ClassRoomId int32 `json:"classRoomId"` //教室编号
	HandsUserId int32 `json:"handsUserId"` //举手学生的Id
	HandsType   int32 `json:"handsType"`   //	1=举手，2=取消举手 必须参数
}

func GetHandsUpProtocol() *HandsUpProtocol {
	handsUpProtocol := &HandsUpProtocol{}
	handsUpProtocol.Command = MESSAGE_CMD_NOGROUP
	//设置命令编号
	handsUpProtocol.ProtocalCommand = HANDS_UP
	handsUpProtocol.NotPushOffLine = true

	return handsUpProtocol
}
