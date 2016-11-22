// DeleteGroupMember
package imbean

//群管理员删除群成员的请求
//cmd=DELETE_GROUP_MEMBERS_BY_GROUP_MANAGER
type DeleteGroupMemberRequest struct {
	RequestBase
	StudentIds    string `json:"users"` //要删除的群用户id集合,以","分隔拼接
	ClassRoomIMId int32  `json:"groupId"`
	RequestUserId int32  `json:"userId"`
}

func GetDeleteGroupMemberRequest() *DeleteGroupMemberRequest {
	request := new(DeleteGroupMemberRequest)
	request.Command = DELETE_GROUP_MEMBER
	return request
}

//群管理员删除群成员的请求的响应
type DeleteGroupMemberResponse struct {
	ResponseBase
}

func (self DeleteGroupMemberRequest) GetCmd() string {
	return self.Command
}
