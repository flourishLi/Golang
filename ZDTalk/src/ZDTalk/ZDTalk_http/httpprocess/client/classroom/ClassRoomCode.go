// ClassRoomCode
package classroom

const (
	ClassRoomID_IS_NULL       = 0 //教室编号为空
	StudentIds_IS_NULL        = 0 //学生集合为空
	Request_User_IS_NULL      = 0 //请求的用户不存在
	Request_FileCount_IS_NULL = 0 //文件数量为空
	File_LimitCount_Is_NULL   = 0 //查找文件数量为空

	CrudType_ADD    = 1 //添加指令 添加到教室 添加到发言区 添加到禁言区
	CrudType_DELETE = 2 //删除指令 从教室删除 从发言区移除 从禁言区删除
	HandsType_Up    = 1 //举手
	HandsType_Down  = 2 //取消举手

	ClasRoom_IS_NOT_EXIT = 0 //教室不存在
	ClasRoom_IS_EXIT     = 1 //教室存在

	User_IS_NOT_EXIT = 0 //用户不存在
	User_IS_EXIT     = 1 //学生存在

	SettingStatus_One   = 1 //1 允许所有人打字，2 允许学生录音，3 禁止游客打字，4 禁止游客举手
	SettingStatus_Two   = 2
	SettingStatus_Three = 3
	SettingStatus_Four  = 4

	RoomStatus_One   = 1 //0 刚创建, 1 初始状态, 2 上课状态, 3 下课状态 必须参数
	RoomStatus_Two   = 2
	RoomStatus_Three = 3
	RoomStatus_ZERO  = 0

	ForbidHandStatus_ForBid = 1 //0 不禁止, 1 禁止
	ForbidHandStatus_NO     = 0

	STUDENT_VIP     = 1 //vip学员
	STUDENT_NORMAL  = 2 //普通学员
	STUDENT_TRY     = 3 //试听学员
	STUDENT_VISITOR = 4 //游客
	Teacher         = 6
)
