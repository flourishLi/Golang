package bean

type HandsForbidRequest struct {
	ClientBaseRequest       //CMD=FORBID_HAND
	ClassRoomId       int32 `json:"classRoomId"`      //教室ID
	ForbidHandStatus  int32 `json:"forbidHandStatus"` //教室的禁止举手状态 0可举手 1禁止举手
}

type HandsForbidResponse struct {
	ClientBaseResponse
}
