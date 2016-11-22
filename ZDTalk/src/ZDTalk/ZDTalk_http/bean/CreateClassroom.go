package bean

type CreateClassroomRequest struct {
	ClientBaseRequest         //CMD=CREATE_CLASS_ROOM
	ClassRoomName     string  `json:"classRoomName"`
	ClassRoomLogo     string  `json:"classRoomLogo"`
	ClassRoomCourse   string  `json:"classRoomCourse"`
	Description       string  `json:"description"`   //教室说明
	SettingStatus     []int32 `json:"settingStatus"` //教室设置状态 1 允许所有人打字，2 允许学生录音，3 禁止游客打字，4 禁止游客举手
}

type CreateClassroomResponse struct {
	ClientBaseResponse       //code=1标识成功 3：此名字的教室已经存在
	ClassRoomId        int32 `json:"classRoomId"` //教室的ID
}
