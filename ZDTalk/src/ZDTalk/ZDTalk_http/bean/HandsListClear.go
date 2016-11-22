package bean

type HandsListClearRequest struct {
	ClientBaseRequest       //CMD=CLEAR_HAND_LIST
	ClassRoomId       int32 `json:"classRoomId"` //教室ID
}

type HandsListClearResponse struct {
	ClientBaseResponse
}
