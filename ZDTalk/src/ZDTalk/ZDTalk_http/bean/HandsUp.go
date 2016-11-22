package bean

type HandsUpRequest struct {
	ClientBaseRequest       //CMD=HANDS_UP
	ClassRoomId       int32 `json:"classRoomId"` //教室ID
	HandsType         int32 `json:"handsType"`   //1=举手，2=取消举手
}

type HandsUpResponse struct {
	ClientBaseResponse
}
