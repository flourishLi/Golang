package customMsg

/***
各种命令编号 全在这里
*/
const (
	SEND_CUSTOM_SERVER_MESSAGE_REQUEST        = 0x405 //发送自定义消息到服务器
	SEND_CUSTOM_SERVER_MESSAGE_RESPONSE       = 0x406 //发送自定义消息到服务器的响应
	SEND_CUSTOM_CLIENT_MESSAGE_REQUEST        = 0x407 //服务器发送到客户端的消息
	SEND_CUSTOM_CLIENT_MESSAGE_RESPONSE       = 0x408 //服务器发送到客户端消息的响应
	ONLINE_HEARTBEAT                    int16 = 1793  //	心跳协议 不转发
	ONLINE_TRANSPOND                    int16 = 1794  //	转发消息
	DRAW_TRANSPOND                      int16 = 1795  //	绘制命令

)
