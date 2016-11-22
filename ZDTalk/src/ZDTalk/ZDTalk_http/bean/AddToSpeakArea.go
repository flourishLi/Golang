package bean

type AddToSpeakAreaRequest struct {
	ClientBaseRequest         //CMD=ADD_REMOVE_SPEAK_AREA
	ClassRoomId       int32   `json:"classRoomId"` //教室ID
	StudentIds        []int32 `json:"studentIds"`  //被添加到发言去的学生用户Id集合
	CrudType          int32   `json:"crudType"`    //1 添加，2 移除
}

type AddToSpeakAreaResponse struct {
	ClientBaseResponse
}
