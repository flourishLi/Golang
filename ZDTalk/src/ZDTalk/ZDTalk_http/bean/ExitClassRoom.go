package bean

type ExitClassRoomRequest struct {
	ClientBaseRequest       //CMD=EXIT_CLASSROOM
	ClassRoomId       int32 `json:"classRoomId"` //教室ID
}

type ExitClassRoomResponse struct {
	ClientBaseResponse
}
