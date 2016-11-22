package protocol

//IM协议 命令编号
const (
	MESSAGE_CMD              string = "PUSH_GROUP_NOTIFICATION"
	MESSAGE_CMD_NOGROUP      string = "PUSH_NOTIFICATION"
	START_FINISH_CLASS       int16  = 0x0621 //老师 上下课
	HANDS_UP                 int16  = 0x0622 //举手(取消)动作
	FORBID_AREA              int16  = 0x0623 //禁止(解除) 禁言区
	HANDS_CLEAR              int16  = 0x0624 //清空举手列表
	SAY_AREA_ADD_USER        int16  = 0x0625 //添加(移除)发言区的成员
	CLASS_ROOM_MEMBER_DELETE int16  = 0x0626 //将学生移除教室
	STUDENT_ENTRY_CLASSROOM  int16  = 0x0627 //学生进入教室
	HAND_FORBID_STATUS       int16  = 0x0628 //教室的禁止举手状态
	ROOM_STATUS              int16  = 0x0629 //教室设置
	EXIT_CLASS_ROOM          int16  = 0x0630 //退出教室

)
