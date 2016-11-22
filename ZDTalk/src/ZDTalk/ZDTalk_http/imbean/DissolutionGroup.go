package imbean

//后台解散群请求
type DissolutionGroupRequest struct {
	RequestBase
	UserId         int32   `json:"userId"`   //解散群的用户ID,如果为0,默认是后台管理员,不是实际用户
	ClassRoomIMIds []int32 `json:"groupIds"` //要解散的群ID集合
}

func GetDissolutionGroupRequest() *DissolutionGroupRequest {
	request := new(DissolutionGroupRequest)
	request.Command = DISSOLUTION_GROUP
	return request
}

//后台解散群请求的响应
type DissolutionGroupResponse struct {
	ResponseBase
}
