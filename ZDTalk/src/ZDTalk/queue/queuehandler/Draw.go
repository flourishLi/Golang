package queuehandler

import (
	"ZDTalk/queue/transmitter"
	//	"ZDTalk/config"
	"ZDTalk/manager/memory"
	msg "ZDTalk/queue/customMsg"
	//	"ZDTalk/queue/publisher"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
)

//画线命令
type CmdDrawLine struct {
	Index   int32  `json:"index"`   //命令索引号
	Color   string `json:"color"`   //颜色
	PointAx int16  `json:"pointAx"` //起点x坐标
	PointAy int16  `json:"pointAy"` //起点y坐标
	PointBx int16  `json:"pointBx"` //终点x坐标
	PointBy int16  `json:"pointBy"` //终点y坐标
}

// 绘制文本命令
type CmdDrawText struct {
	Index     int32  `json:"index"`     //命令索引号
	PointX    int16  `json:"pointX"`    //文本框左上点x坐标
	PointY    int16  `json:"pointY"`    //文本框左上点y坐标
	TextWidth int16  `json:"textWidth"` //文本框宽度
	Text      string `json:"text"`      //绘制文本
}

//转发的内容
type DrawContent struct {
	ClassroomId int32        `json:"classroomId"` //当前教室id
	FileId      int32        `json:"fileId"`      //当前ppt图片的url
	FileUrl     string       `json:"fileUrl"`     //当前ppt图片的url
	CmdSet      CmdSetStruct `json:"cmdSet"`      //绘制命令集
	Width       int16        `json:"width"`       //绘图区宽度
	Height      int16        `json:"height"`      //绘图区高度
	FontSize    int16        `json:"fontSize"`    //字号

}

type DrawCommand struct {
}

type CmdSetStruct struct {
	CmdDrawLine []CmdDrawLine `json:"cmdDrawLine"`
	CmdDrawText []CmdDrawText `json:"cmdDrawText"` //划线命令
}

//转发command消息
func (draw DrawCommand) TransportMessage(senderImUserId int32, clientMsg msg.SendCustomClientMessageRequest) {
	drawContent := DrawContent{}
	err := json.Unmarshal(clientMsg.MessageContent, &drawContent)

	if err != nil {
		logs.GetLogger().Info("json parse MessageContent is wrong", err)
		return
	}
	logs.GetLogger().Info("--------------- TransportMessage ----------------")
	if drawContent.ClassroomId == 0 {
		logs.GetLogger().Info(" classRoomId is 0 ")
		return
	}

	//获取当前教室的在线列表结合
	onLineUsers := GetOnLineUsers(drawContent.ClassroomId)
	logs.GetLogger().Info("转发消息时，当前在线用户集合 ", onLineUsers)
	//向每一个成员转发消息
	for _, onLineUserId := range onLineUsers {
		receiverImUserId := memory.GetUserInfoMemoryManager().GetUserIMId(onLineUserId) //targetid

		if receiverImUserId != senderImUserId {
			logs.GetLogger().Info("转发 绘制 receiveIMUserId", receiverImUserId)
			logs.GetLogger().Info("转发 绘制 content", string(clientMsg.MessageContent))
			transmitter.SendMessage2IMServer(clientMsg.ServerId, receiverImUserId, receiverImUserId, clientMsg.MessageFormat, clientMsg.MessageContent, clientMsg.MessageId, clientMsg.SendTime)
		}
	}
}
