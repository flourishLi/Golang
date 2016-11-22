package imbean

const (
	SUCCESS              = 0 //与IM系统交互的Http接口返回值 返回成功
	GROUPTYPE            = 3 //isDisGroup	1=讨论组，2=群 3聊天室
	ENTERGROUP_BYMANAGER = 1 //1是管理员邀请用户入群；2是用户主动申请入群
	ENTERGROUP_BYUSER    = 2 //1是管理员邀请用户入群；2是用户主动申请入群
	PUSH                 = 0 //0=推送，1=不推送
	DONT_PUSH            = 1 //0=推送，1=不推送
	NEEDMANAGERPOWER     = 0 //	0=只有管理员可以修改，1=无须为管理员，均可修改 （Go版）
	NOMANAGERPOWER       = 1 //	0=只有管理员可以修改，1=无须为管理员，均可修改 （Go版）

)
