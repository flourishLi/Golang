// QuitGroup
package imbean

//群成员主动退出群的请求

//cmd=QUIT_GROUP_BY_OWN
type QuitGroupRequest struct {
	RequestBase
	NewManagerId int32 `json:"newManagerId"`
}

func GetQuitGroupRequest() *QuitGroupRequest {
	request := new(QuitGroupRequest)
	request.Command = QUIT_GROUP
	return request
}

//群成员自动退出群的请求的响应
type QuitGroupResponse struct {
	ResponseBase
}

func (self QuitGroupRequest) GetCmd() string {
	return self.Command
}
