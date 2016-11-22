// memoryCode
package memory

const (
	ClassRoomID_IS_NULL = 0 //教室编号为空
	StudentIds_IS_NULL  = 0 //学生集合为空
	CrudType_ADD        = 1 //添加指令 添加到教室 添加到发言区 添加到禁言区
	CrudType_DELETE     = 2 //删除指令 从教室删除 从发言区移除 从禁言区删除

	HandsType_Up   = 1 //举手
	HandsType_Down = 2 //取消举手

	ClasRoom_IS_NOT_EXIT = 0 //教室不存在
	ClasRoom_IS_EXIT     = 1 //教室不存在

)
