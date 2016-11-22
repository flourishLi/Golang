package bean

type EntryClassRoomRequest struct {
	ClientBaseRequest       //CMD=ENTRY_CLASSROOM
	ClassRoomId       int32 `json:"classRoomId"` //教室ID
}

type EntryClassRoomResponse struct {
	ClientBaseResponse
}
