package bean

type DeleteAddStudentsRequest struct {
	ClientBaseRequest         //CMD=ADD_TO_SPEAK_AREA
	ClassRoomId       int32   `json:"classRoomId"` //教室ID
	StudentIds        []int32 `json:"studentIds"`  //被移除教室的学生用户Id集合
	CrudType          int32   `json:"crudType"`    //1 添加，2 移除
}

type DeleteAddStudentsResponse struct {
	ClientBaseResponse
}
