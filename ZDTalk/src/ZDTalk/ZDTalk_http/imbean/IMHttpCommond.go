package imbean

const (
	CREATE_GROUP        string = "CREATE_GROUP"                          //创建群
	DISSOLUTION_GROUP   string = "BACKSTAGE_DISSOLUTION_GROUP"           //后台解散群
	DELETE_GROUP_MEMBER string = "DELETE_GROUP_MEMBERS_BY_GROUP_MANAGER" //删除群成员
	ADD_GROUP_MEMBER    string = "ADD_GROUP_MEMBER"                      //添加群成员
	UPDATE_GROUP_INFO   string = "UPDATE_GROUP_INFO"                     //修改群资料
	QUIT_GROUP          string = "QUIT_GROUP_BY_OWN"                     //群成员主动退群
	SIGNUP              string = "REGISTER"                              //用户注册
	UPDATE_USERINFO     string = "UPDATE_USER_INFO"                      //更改用户信息
	UPDATE_USERPWD      string = "MODIFY_PASSWORD"                       //更改用户密码

	RESET_USERPWD string = "RESET_PASSWORD" //重置密码

)
