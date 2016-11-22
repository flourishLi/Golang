package bean

type ForbidAreaRequest struct {
	ClientBaseRequest         //CMD=FORBID_HAND
	ClassRoomId       int32   `json:"classRoomId"` //教室ID
	StudentIds        []int32 `json:"studentIds"`  //学生ID集合
	CrudType          int32   `json:"crudType"`    //1 禁止，2 解禁
}

type ForbidAreaResponse struct {
	ClientBaseResponse
}
