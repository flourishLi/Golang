package imbean

//群管理员添加群成员的请求
//cmd=ADD_GROUP_MEMBER
type AddGroupMemberRequest struct {
	RequestBase
	StudentIds    string `json:"users"` //要添加的群用户id集合,以","分隔拼接
	ClassRoomIMId int32  `json:"groupId"`
	RequestType   int32  `json:"requestType"` //	1是管理员邀请用户入群；2是用户主动申请入群
	RequestUserId int32  `json:"userId"`      //requestType为1时，表示邀请者(管理员)的id；requestType为2时，忽略该字段
	Content       string `json:"content"`
	IsDisGroup    int32  `json:"isDisGroup"` //isDisGroup=1为讨论组，isDisGroup=2为群
}

func GetAddGroupMemberRequest() *AddGroupMemberRequest {
	request := new(AddGroupMemberRequest)
	request.Command = ADD_GROUP_MEMBER
	return request
}

//群管理员添加群成员的请求的响应
type AddGroupMemberResponse struct {
	ResponseBase
}

func (self AddGroupMemberRequest) GetCmd() string {
	return self.Command
}
