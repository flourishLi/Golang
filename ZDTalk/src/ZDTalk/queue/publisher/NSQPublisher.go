package publisher

//NSQ 生产者(发送消息到NSQ)
import (
	//	"ZDTalk/queue/customMsg"
	logs "ZDTalk/utils/log4go"
	"sync"

	"github.com/nsqio/go-nsq"
)

type QueuePublisher interface {
	SendMessage(topic string, data []byte) error
	//	Close()
}

type NSQPublisher struct {
	Publisher *nsq.Producer
	Server    string
}

type QueuePublisherPool struct {
	mutex                 sync.Mutex
	queues                map[string]QueuePublisher
	MaxConnections        int
	CurrentUseConnections int
}

var queuePublisherPool *QueuePublisherPool

func GetQueuePool() *QueuePublisherPool {
	return queuePublisherPool
}

func InitPublisherPool() {
	logs.GetLogger().Info("----------------------- queue init -----------------------")
	queuePublisherPool = &QueuePublisherPool{}
	queuePublisherPool.MaxConnections = 200
}

func (self NSQPublisher) SendMessage(topic string, data []byte) error {
	logs.GetLogger().Info("-------------------------- NSQ SendMessage topic %s--------------------------", topic)
	return self.Publisher.Publish(topic, data)
}

func CreateQueuePublisher(server string) (QueuePublisher, error) {
	config := nsq.NewConfig()
	publisher, err := nsq.NewProducer(server, config)
	queue := &NSQPublisher{}
	queue.Publisher = publisher
	queue.Server = server
	if err != nil {
		logs.GetLogger().Error("Create Producer(Failed):%s", err)
		return nil, err
	}
	return queue, nil
}

func (self *QueuePublisherPool) GetQueue(server string) (QueuePublisher, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.queues == nil {
		self.queues = make(map[string]QueuePublisher)
	}

	if q, ok := self.queues[server]; ok {
		return q, nil
	} else {
		q, err := CreateQueuePublisher(server)

		if err == nil {
			self.queues[server] = q
		}

		return q, err
	}

	return nil, nil
}

//func SendMessageToIMServer(msg customMsg.ArrowMessage) {
//	logs.GetLogger().Info("----------------------- SendMessageToIMServer -----------------------")
//	q := GetQueuePool()

//	p, err := q.GetQueue("IMServer_from_gateway_1")
//	//	p, err := q.GetQueue("IMServer_from_server_1")

//	if err != nil {
//		logs.GetLogger().Error("get queue error:" + err.Error())
//	}

//	err = p.SendMessage("IMServer_from_server_1", msg.ToBytes(customMsg.SERVER_MESSAGE))
//	if err != nil {
//		logs.GetLogger().Error("send message error: " + err.Error())
//	}
//}
