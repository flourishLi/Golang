package protocol

import (
	"ZDTalk/ZDTalk_http/imbean"
)

//最终要发送给IM的协议
type ProtocolShell struct {
	imbean.RequestBase
	SenderUserId int32 `json:"userId"` //发送通知者Id
	//ClassRoomsId    []int32 `json:"groupIds"`       //接收通知的群Id
	UserIds         []int32 `json:"userIds"`        //	接收通知的用户ID集合，最多一次不要超过1000个
	ProtocalCommand int16   `json:"notifyType"`     //发送通知类型
	SendContent     string  `json:"content"`        //发送内容 base64编码
	ShowContent     string  `json:"showContent"`    //在通知栏上显示的内容
	ShowName        string  `json:"showName"`       //在通知栏上显示的名称
	NotPushOffLine  bool    `json:"notPushOffLine"` //不推送离线消息(为true时跳过推送离线消息，false时推送离线消息)
}

func GetProtocolShell() *ProtocolShell {
	protocolShell := &ProtocolShell{}
	//设置命令编号
	protocolShell.Command = MESSAGE_CMD_NOGROUP
	return protocolShell
}
