package bean

type UpdateClassRoomStatusRequest struct {
	ClientBaseRequest       //CMD=UPDATE_CLASSROOM_STATUS
	ClassRoomId       int32 `json:"classRoomId"`
	ClassRoomStatus   int32 `json:"classRoomStatus"` //教室设置状态 教室当前状态0 刚创建, 1 初始状态, 2 上课状态, 3 下课状态
}

type UpdateClassRoomStatusResponse struct {
	ClientBaseResponse //code=1标识成功 3：room id不存在
}
