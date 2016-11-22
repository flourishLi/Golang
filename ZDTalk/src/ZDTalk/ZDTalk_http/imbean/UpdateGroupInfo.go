// UpdateGroupInfo.go
package imbean

//修改群资料的请求
type UpdateGroupInfoRequest struct {
	RequestBase
	ClassRoomName    string `json:"groupName"`        //群名称（不修改时可为空，但不允许和群头像URL、群简介同时为空）
	ClassRoomLogo    string `json:"groupIcon"`        //群头像URL（不修改时可为空，但不允许和群名称、群简介同时为空）
	Description      string `json:"groupIntroduce"`   //群简介 （不修改时可为空，但不允许和群名称、群Icon同时为空）（Go版）
	NeedManagerPower int32  `json:"needManagerPower"` //被修改群资料的群ID
	RequestUserId    int32  `json:"userId"`           //修改群资料的用户ID
	ClassRoomIMId    int32  `json:"groupId"`          //被修改群资料的群ID
}

func GetUpdateGroupInfoRequest() *UpdateGroupInfoRequest {
	request := new(UpdateGroupInfoRequest)
	request.Command = UPDATE_GROUP_INFO
	return request
}

//修改群资料的请求的响应
type UpdateGroupInfoResponse struct {
	ResponseBase
}

func (self UpdateGroupInfoRequest) GetCmd() string {
	return self.Command
}
