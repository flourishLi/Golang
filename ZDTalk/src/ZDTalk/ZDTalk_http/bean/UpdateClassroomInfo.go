package bean

type UpdateClassroomRequest struct {
	ClientBaseRequest        //CMD=UPDATE_CLASS_ROOM
	ClassRoomId       int32  `json:"classRoomId"`
	ClassRoomLogo     string `json:"classRoomLogo"`
	Description       string `json:"description"` //教室说明
	ClassRoomName     string `json:"classRoomName"`
	ClassRoomStatus   int32  `json:"classRoomStatus"` //教室当前状态0 刚创建, 1 初始状态, 2 上课状态, 3 下课状态
	SettingStatus     int32  `json:"settingStatus"`   //教室设置状态 1 允许所有人打字，2 允许学生录音，3 禁止游客打字，4 禁止游客举手
	ClassRoomCourse   string `json:"classRoomCourse"`
}

type UpdateClassroomResponse struct {
	ClientBaseResponse //code=1标识成功 3：room id不存在
}
