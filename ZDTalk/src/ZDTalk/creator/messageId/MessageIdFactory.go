package messageId

var meeeageIdCreator *MessageIdCreator

func InitMessageIdCreator(serverType int32) {
	if meeeageIdCreator == nil {
		meeeageIdCreator = CreateMessageIdCreator(serverType)
	}
}

func CreateId() int64 {
	return meeeageIdCreator.CreateId()
}

const (
	//	GATEWAY_SERVER_ID = 0
	//	IM_SERVER_ID      = 9
	//	GROUP_SERVER_ID   = 8
	//	SUPPORT_SERVER_ID = 7
	ZDTALK_MESSAGE_ID_START_NUM = 5 //早道中的MessageId 首位数字
)
