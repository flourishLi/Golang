// CreatGroupRequest
package imbean

import (
	"ZDTalk/config"
)

//创建群的请求
//cmd=CREATE_GROUP
type CreateGroupRequest struct {
	RequestBase
	ClassRoomName string `json:"groupName"`
	ClassRoomLogo string `json:"groupIcon"`
	AppId         int32  `json:"appId"` //7133
	StudentIds    string `json:"users"` //需要添加的用户集合,以","分隔拼接
	Content       string `json:"content"`
	IsDisGroup    int32  `json:"isDisGroup"`     //isDisGroup=1为讨论组，isDisGroup=2为群
	Description   string `json:"groupIntroduce"` //群简介
	RequestUserId int32  `json:"userId"`
}

func GetCreateGroupRequest() *CreateGroupRequest {
	request := new(CreateGroupRequest)
	request.Command = CREATE_GROUP
	request.AppId = config.ConfigNodeInfo.AppId
	return request
}

//创建群的请求的响应
type CreateGroupResponse struct {
	ResponseBase
	ClassRoomIMId int32 `json:"groupId"`
}

func (self CreateGroupRequest) GetCmd() string {
	return self.Command
}
