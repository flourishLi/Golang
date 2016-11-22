package errorcode

const (
	SUCCESS int32 = 1 //成功
	//request参数错误 100开头
	CMD_IS_NULL           int32  = 1000
	CMD_IS_NULL_ERROR_MSG string = "cmd is null is Illegal"

	JSON_PRASEM_ERROR       int32  = 1001                          //JSON解析错误
	CLASS_ROOM_ID_ERROR     int32  = 1002                          //ClassRoomId 错误
	CLASS_ROOM_ID_ERROR_MSG string = "chassRoomId == 0 is Illegal" //ClassRoomId 错误

	CLASS_ROOM_IS_EXIST     int32  = 1003                     //教室已经存在
	CLASS_ROOM_IS_EXIST_MSG string = "the chassRoom is exist" //教室已经存在

	CLASS_ROOM_IS_NOT_EXIST     int32  = 1004                        //教室不存在
	CLASS_ROOM_IS_NOT_EXIST_MSG string = "the classRoom isn't exist" //教室不存在

	CLASS_ROOM_NAME_CAN_NOT_NULL     int32  = 1005                          //创建教室时，教室名字为空
	CLASS_ROOM_NAME_CAN_NOT_NULL_MSG string = "CLASSROOM_NAME_CAN_NOT_NULL" //创建教室时，教室名字为空

	REQUEST_USER_ID_CAN_NOT_NULL     int32  = 1006                   //请求的id不能为空
	REQUEST_USER_ID_CAN_NOT_NULL_MSG string = "USER_ID_CAN_NOT_NULL" //

	REQUEST_USER_ID_IS_NOT_EXIT     int32  = 1007                              //请求的id数据库中不存在
	REQUEST_USER_ID_IS_NOT_EXIT_MSG string = "USER_ID_IS_NOT_EXIT_IN_DATABASE" //

	REQUEST_USER_ID_HAS_NO_AHTHORITY     int32  = 1008                       //请求的id没有权限
	REQUEST_USER_ID_HAS_NO_AHTHORITY_MSG string = "USER_ID_HAS_No_Authority" //

	GET_CLASSROOMLIST_ERR     int32  = 1007 //获取数据库教室列表失败
	GET_CLASSROOMLIST_ERR_MSG string = "初始化教室列表失败"

	GET_CLASSROOMINFO_DATEBASE_ERR int32 = 1008 //获取数据库教室信息失败

	Create_CLASSROOM_DATEBASE_ERR int32 = 1009 //创建教室失败 数据库中

	DELETE_CLASSROOM_DATEBASE_ERR int32 = 1010 //删除教室失败 数据库中

	UPDATE_CLASSROOM_DATEBASE_ERR int32 = 1011 //更新教室失败 数据库中

	DATABASE_ROW_AFFECT_NULL     int32  = 1012             //数据库操作时  没有行受影响
	DATABASE_ROW_AFFECT_NULL_MSG string = "操作重复，数据库没有任何影响" //数据库操作时  没有行受影响

	HANDS_UP_DATEBASE_ERROR     int32 = 1013 //举手(取消)动作请求 更新失败
	HANDS_FORBID_DATEBASE_ERROR int32 = 1014 //举手(取消)动作请求 更新失败

	HANDS_CLEAR_DATEBASE_ERROR    int32 = 1015 //清空举手动作请求 更新失败
	ADDSPEAKTOAREA_DATEBASE_ERROR int32 = 1016 // 添加举手到发言列表
	DELETESTUDENTS_DATEBASE_ERROR int32 = 1017 // 将学生移除教室

	HANDS_TYPE_IS_WRONG    int32  = 1018
	HAND_TYPE_IS_WRONG_MSG string = "HandsType must be 1 or 2"

	CURD_TYPE_IS_WRONG     int32  = 1018
	CURD_TYPE_IS_WRONG_MSG string = "CURDType must be 1 or 2"

	SETTING_STATUS_IS_WRONG     int32  = 1019
	SETTING_STATUS_IS_WRONG_MSG string = "settingStatus must be 1 or 2 or 3 or 4"

	CLASS_ROOM_STATUS_IS_WRONG     int32  = 1020
	CLASS_ROOM_STATUS_IS_WRONG_MSG string = "RoomStatusStatus must be 0 or 1 or 2 or 3"

	RESOURCE_IS_NOT_EXIST     int32  = 1021                       //资源不存在
	RESOURCE_IS_NOT_EXIST_MSG string = "the Resource isn't exist" //资源不存在

	RESOURCE_NAME_IS_NOT_NULL     int32  = 1022                          //资源不存在
	RESOURCE_NAME_IS_NOT_NULL_MSG string = "the ResourceName isn't null" //资源名称不能为空

	RESOURCE_PATH_IS_NOT_NULL     int32  = 1023                          //资源不存在
	RESOURCE_PATH_IS_NOT_NULL_MSG string = "the ResourcePath isn't null" //资源路径不能为空

	UPLOAD_RESOURCE_DATEBASE_ERROR int32 = 1024 //上传资源错误
	DELETE_RESOURCE_DATEBASE_ERROR int32 = 1025 // 删除资源错误

	FILE_ID_IS_NOT_EXIST     int32  = 1026                     //资源不存在
	FILE_ID_IS_NOT_EXIST_MSG string = "the fileId isn't exist" //资源不存在

	RESOURCE_IS_DELETE     int32  = 1027                      //资源已删除
	RESOURCE_IS_DELETE_MSG string = "the Resource is deleted" //资源已删除

	LOGIN_NAME_CAN_NOT_NULL     int32  = 1028                         //登录账号不能为空
	LOGIN_NAME_CAN_NOT_NULL_MSG string = "the loginName can not null" //登录账号不能为空

	LOGIN_PASSWORD_CAN_NOT_NULL     int32  = 1029                        //登录密码不能为空
	LOGIN_PASSWORD_CAN_NOT_NULL_MSG string = "the Password can not null" //登录密码不能为空

	LOGIN_DATABASE_IS_ERROR     int32  = 1030                      //登录校验数据库发送错误
	LOGIN_DATABASE_IS_ERROR_MSG string = "Login error in database" //登录校验数据库发送错误

	LOGINNAME_PSW_IS_ERROR     int32  = 1031                             //用户名或者密码错误
	LOGINNAME_PSW_IS_ERROR_MSG string = "loginName or Password is error" //用户名或者密码错误

	LOGINNAME_IS_NOT_EXIST     int32  = 1032                     //用户名不存在
	LOGINNAME_IS_NOT_EXIST_MSG string = "loginName is not exist" //用户名不存在

	CREATE_GROUP_IM_IS_ERROR        int32 = 1032 //创建群
	DISSOLUTION_GROUP_IM_IS_ERROR   int32 = 1033 //解散群
	UPDATE_GROUP_IM_IS_ERROR        int32 = 1034 //更新群
	DELETE_GROUP_MEMBER_IM_IS_ERROR int32 = 1035 //删除群成员
	ADD_GROUP_MEMBER_IM_IS_ERROR    int32 = 1036 //添加群成员

	QUITE_GROUP_MEMBER_IM_IS_ERROR int32 = 1037

	UPLOADFILE_SUCESS int32 = 1
	UPLOADFILE_FAIL   int32 = 1039

	HANDSUP_PROTOCOL_IM_IS_ERROR       int32 = 1040 //举手消息错误
	HANDSFORBID_PROTOCOL_IM_IS_ERROR   int32 = 1041 //禁止举手消息错误
	HANDSCLEAR_PROTOCOL_IM_IS_ERROR    int32 = 1042 //清除举手消息错误
	SETTINGSTATUS_PROTOCOL_IM_IS_ERROR int32 = 1043 //settingStatus消息错误
	ROOMSTATUS_PROTOCOL_IM_IS_ERROR    int32 = 1044 //roomStatus消息错误

	USER_SIGNUP_IS_SUCCESS     int32  = 1045                   //注册成功
	USER_SIGNUP_IS_SUCCESS_MSG string = "User Sign up Success" //注册成功
	USER_SIGNUP_IS_ERROR       int32  = 1046                   //注册失败

	USER_ID_IS_NOT_EXIT     int32  = 1047 //用户存在
	USER_ID_IS_NOT_EXIT_MSG string = "userID is not exit"

	USER_PWD_IS_ERROR     int32  = 1048 //用户密码错误
	USER_PWD_IS_ERROR_MSG string = "User Password is wrong"

	USER_NAME_CAN_NOT_NULL     int32  = 1049                        //用户昵称不能为空
	USER_NAME_CAN_NOT_NULL_MSG string = "the Username can not null" //登录账号不能为空

	USER_ROLE_CAN_NOT_NULL     int32  = 1050                    //用户角色不能为空
	USER_ROLE_CAN_NOT_NULL_MSG string = "the role can not null" //登录账号不能为空

	USER_ROLE_IS_ERR     int32  = 1051                         //用户角色编号错误
	USER_ROLE_IS_ERR_MSG string = "the role must be 1 2 3 4 6" //

	USER_SIGN_UP__IS_ERR_IN_DATABASE   int32 = 1052
	USERINFO_UPDATE_IS_ERR_IN_DATABASE int32 = 1053
	USERPWD_UPDATE_IS_ERR_IN_DATABASE  int32 = 1054 //

	USER_ID_CAN_NOT_NULL     int32  = 1055                      //用户id不能为空
	USER_ID_CAN_NOT_NULL_MSG string = "the UserId can not null" //不能为空

	DEVICE_TYPE_IS_ERR     int32  = 1056                           //设备类型错误
	DEVICE_TYPE_IS_ERR_MSG string = "the deviceType must be 1 2 3" //不能为空

	OLD_PWD_CAN_NOT_NULL     int32  = 1057                           //设备类型错误
	OLD_PWD_CAN_NOT_NULL_MSG string = "the oldpassword can not null" //不能为空

	NEW_PWD_CAN_NOT_NULL     int32  = 1058                            //设备类型错误
	NEW_PWD_CAN_NOT_NULL_MSG string = "the new password can not null" //不能为空

	FUZZY_SEARCH_USERINFO_IS_ERR             int32  = 1059 //模糊查询错误
	FUZZY_SEARCH_USERINFO_RESULT_IS_NULL     int32  = 1060 //模糊查询结果为空
	FUZZY_SEARCH_USERINFO_RESULT_IS_NULL_MSG string = "用户昵称查询结果为空"

	USER_SIGNUP_IM_IS_ERROR      int32 = 1061 //创建用户
	USER_PWD_UPDATE_IM_IS_ERROR  int32 = 1062 //修改密码
	USER_PWD_RESET_IM_IS_ERROR   int32 = 1063 //重置密码
	USER_INFO_UPDATE_IM_IS_ERROR int32 = 1064 //修改信息

	TOFORBID_DATEBASE_ERROR int32 = 1065 // 添加举手到发言列表

	TOFORBID_PROTOCOL_IM_IS_ERROR            int32 = 1066 //学生添加移除到禁烟区消息错误
	STD_ENTRY_CLASSROOM_PROTOCOL_IM_IS_ERROR int32 = 1067 //学生进入教室消息错误

	UPDATE_FORBIDHAND_STATUS_DATEBASE_ERR int32 = 1068 //更新教室失败 禁止举手状态

	FORBID_HAND_STATUS_IS_WRONG     int32  = 1069
	FORBID_HAND_STATUS_IS_WRONG_MSG string = "forbidHandStatus must be 0 or 1"

	EXIT_CLASSROOM_PROTOCOL_IM_IS_ERROR int32 = 1070 //退出教室消息错误

	USER_ID_IS_EXIST     int32  = 1047 //用户存在
	USER_ID_IS_EXIST_MSG string = "userID is exist"

	USER_ID_IS_EXIT     int32  = 1071 //用户已退出
	USER_ID_IS_EXIT_MSG string = "userID is exit"

	FILE_NAME_IS_NULL     int32  = 1072 //资源名称为空
	FILE_NAME_IS_NULL_MSG string = "filename can not be null"

	FILE_COUNT_IS_NULL     int32  = 1073 //子资源数量为空
	FILE_COUNT_IS_NULL_MSG string = "subfile count not be null"

	FILE_SEARCH_INDEX_IS_NULL     int32  = 1074 //索引为空
	FILE_SEARCH_INDEX_IS_NULL_MSG string = "searchIndex is null"

	FILE_SEARCH_LIMIT_IS_NULL     int32  = 1075 //limit为空
	FILE_SEARCH_LIMIT_IS_NULL_MSG string = "Limit is null"

	//自动退出
	NO_SUCH_CMD              int32  = 2000                             //指令 错误
	REQUEST_FILTER_ERROR     int32  = 3000                             //服务器内部Http过滤器错误
	REQUEST_FILTER_ERROR_MSG string = "Server Inner Http Filter Error" //服务器内部Http过滤器错误

	REQUEST_METHOD_ERROR     int32  = 4000                                             //请求方式错误,仅支持Post请求
	REQUEST_METHOD_ERROR_MSG string = "Request Method Error, Only Support Method POST" //请求方式错误,仅支持Post请求
)
