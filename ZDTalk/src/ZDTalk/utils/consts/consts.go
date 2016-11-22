package consts

const (
	NSQLOOKUPD_TCP_ADDRESS  = "127.0.0.1:4160"
	NSQLOOKUPD_HTTP_ADDRESS = "127.0.0.1:4161"
	//	NSQ_TCP_ADDRESS  = "118.126.142.86:4150" //cisco所用服务器地址
	//	NSQ_HTTP_ADDRESS = "118.126.142.86:4151" //cisco所用服务器地址
	NSQ_TCP_ADDRESS         = "127.0.0.1:4150"
	NSQ_HTTP_ADDRESS        = "127.0.0.1:4151"
	TOPIC_GROUP_TO_IMSERVER = "im_msg"        //群发消息到NSQ的topic
	TOPIC_IMSERVER_TO_GROUP = "gp_msg"        //ImServer发消息到NSQ的topic
	NSQ_CHANNEL             = "clcong"        //发消息到NSQ的channel
	FZONE                   = "fzone"         //朋友圈转发通知
	FZONE_CHANNEL           = "fzone_channel" //朋友圈转发通知
	APNS_GATEWAY_SANDBOX    = "gateway.sandbox.push.apple.com:2195"
	APNS_GATEWAY            = "gateway.push.apple.com:2195"

	APNS_CERTIFICATEFILE_PREFIX = "pem/"
	APNS_KEYFILE_PREFIX         = "pem/"
	APNS_CERTIFICATEFILE_SUFFIX = "/cert.pem"
	APNS_KEYFILE_SUFFIX         = "/key.pem"

	APNS_CERTIFICATEFILE_TEST  = "pem/test/cert.pem"
	APNS_KEYFILE_TEST          = "pem/test/key.pem"
	NEED_FRIENDSHIP            = false
	TOPIC_GROUP_TO_PROXY       = "group_proxy"
	TOPIC_GROUP_TO_PROXY_OTHER = "group_proxy_other"
	TOPIC_IM_TO_PROXY          = "im_proxy"
	TOPIC_IM_TO_PROXY_OTHER    = "im_proxy_other"
	PROXY_CHANNEL              = "proxy"
)

var Serverstatus = true
