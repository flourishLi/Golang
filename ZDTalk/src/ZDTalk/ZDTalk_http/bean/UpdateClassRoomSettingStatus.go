package bean

type UpdateClassRoomSettingStatusRequest struct {
	ClientBaseRequest         //CMD=UPDATE_CLASSROOM_SETTING_STATUS
	ClassRoomId       int32   `json:"classRoomId"`
	SettingStatus     []int32 `json:"settingStatus"` //教室设置状态 1 允许所有人打字，2 允许学生录音，3 禁止游客打字，4 禁止游客举手
}

type UpdateClassRoomSettingStatusResponse struct {
	ClientBaseResponse //code=1标识成功 3：room id不存在
}
