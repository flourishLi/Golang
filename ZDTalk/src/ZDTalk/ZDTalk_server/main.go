package main

import (
	"ZDTalk/ZDTalk_http/httphandler"
	"ZDTalk/config"
	MessageIdFactory "ZDTalk/creator/messageId"
	"ZDTalk/db/database"
	"ZDTalk/manager/memory"
	NSQPublisher "ZDTalk/queue/publisher"
	queueHandler "ZDTalk/queue/queuehandler"
	log4go "ZDTalk/utils/log4go"
	//	"ZDTalk/utils/timeutils"
	//	"ZDTalk/manager/onLine"
	"os"
	"time"

	"github.com/nsqio/go-nsq"
)

var logs = log4go.GetLogger()

func main() {
	logs.Info("------------- ZDTalk beginning server -------------")
	err := config.LoadConfig("config.json")
	if err != nil {
		logs.Error("加载配置文件失败 停止服务 " + err.Error())
		return
	}

	//初始化消息Id创建工厂
	MessageIdFactory.InitMessageIdCreator(MessageIdFactory.ZDTALK_MESSAGE_ID_START_NUM)

	//加载数据库
	database.OpenDb()
	NSQPublisher.InitPublisherPool()

	logs.Info("------------- ZDTalk Initial Memory Begin -------------")
	//加载classroominfo数据表到roomMemoryManager
	memory.GetClassRoomMemoryManager()
	//加载userinfo数据表到userInfoMemoryManager
	memory.GetUserInfoMemoryManager()
	//加载uploadresource数据表到uploadResourceMemoryManager
	memory.GetUploadResourceMemoryManager()
	//加载drawCommandMemoryManager
	memory.GetDrawCommandMemoryManager()

	//初始化在线用户列表
	memory.GetQueueOnlineMemoryManager()

	logs.Info("------------- ZDTalk Initial Memory End -------------")

	logs.Info("ZDTalk Http interface location--> %s , port--> %d", config.ConfigNodeInfo.ZDTalkInterfaceLocation, config.ConfigNodeInfo.ZDTalkInterfacePort)

	InitIMServerQueue(config.ConfigNodeInfo.IMServerRecvCustomMsgNsqUrl, config.GetNsqRecvTopic())

	//	go onLine.EntryClassRoomAndSendHeartBeat()

	httphandler.InitServlet(config.ConfigNodeInfo.ZDTalkInterfaceLocation, config.ConfigNodeInfo.ZDTalkInterfacePort)

}

//初始化与接收IMServer消息的NSQ
func InitIMServerQueue(imserverQueueIp string, imserverQueueTopic string) {
	logs.Info("ZDTalk 开始连接 IMServer NSQ Address IP %s , Topic %s", imserverQueueIp, imserverQueueTopic)
	config := nsq.NewConfig()

	imserverQueueConsumer, err := nsq.NewConsumer(imserverQueueTopic,
		imserverQueueTopic+"_channel", config)
	if err != nil {
		logs.Error("初始化IMServer 消费者失败:", err)
		time.Sleep(1 * time.Second)
		os.Exit(0)
		return
	}

	handler := &queueHandler.IMServerCustomMessageReceiveHandler{NsqConsumer: imserverQueueConsumer}
	imserverQueueConsumer.AddHandler(handler)

	err = imserverQueueConsumer.ConnectToNSQD(imserverQueueIp)
	if err != nil {
		logs.Error("连接IMServer nsqd失败 error %s  Url %s , Topic %s", err.Error(), imserverQueueIp, imserverQueueTopic)
		time.Sleep(1 * time.Second)
		os.Exit(0)
		return
	}
	logs.Info("ZDTalk 连接 ArrowIMServer nsqd成功: Url %s , Topic %s", imserverQueueIp, imserverQueueTopic)
}
