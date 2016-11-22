package config

import (
	FileUtils "ZDTalk/utils/fileutils"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
)

type ConfigNode struct {
	IMMessageUrl string //消息推送
	//UploadResourceUrl             string //资源上传 Http接口地址
	IMHttpGroupUrl string //IM 群 Http接口地址
	IMHttpUrl      string //IM 个人 Http接口地址

	ZDTalkInterfaceLocation       string //ZDTalk Url Location(***：8080/location)
	ZDTalkInterfacePort           int32  //ZDTalk Url	端口号 (***：port/location)
	IMServerRecvCustomMsgNsqUrl   string //接收IMServer 自定义消息NSQ Url地址
	IMServerRecvCustomMsgNsqTopic string //接收IMServer 自定义消息NSQ Topic
	IMServerRecvCustomMsgNsqId    int32  //接收IMServer 自定义消息NSQ Id

	IMServerSendCustomMsgNsqUrl   string //发送IMServer 自定义消息NSQ Url地址
	IMServerSendCustomMsgNsqTopic string //发送IMServer 自定义消息NSQ Topic
	IMServerSendCustomMsgNsqId    int32  //发送IMServer 自定义消息NSQ Id
	OnLineTimeSpan                int64  //刷新在线用户的时间间隔
	AppId                         int32  //IM的appId
}

var ConfigNodeInfo ConfigNode

func LoadConfig(configPath string) error {
	configBuf, err := FileUtils.ReadAll(configPath)
	if err != nil {
		return err
	}
	tmpNode := new(ConfigNode)
	err = json.Unmarshal(configBuf, tmpNode)
	if err != nil {
		logs.GetLogger().Error("解析config.json内容失败 " + err.Error())
		return err
	}
	ConfigNodeInfo = *tmpNode
	return nil
}

//获取接收IM服务消息的NSQ 的Topic完整地址
func GetNsqRecvTopic() string {
	//	if ConfigNodeInfo == nil {
	//		logs.GetLogger().Error("初始化config.json文件时有问题，配置文件加载失败，请查看")
	//		return ""
	//	}
	return ConfigNodeInfo.IMServerRecvCustomMsgNsqTopic + fmt.Sprintf("%d", ConfigNodeInfo.IMServerRecvCustomMsgNsqId)
}

//获取发送IM服务消息的NSQ 的Topic完整地址
func GetNsqSendTopic() string {
	//	if ConfigNodeInfo == nil {
	//		logs.GetLogger().Error("初始化config.json文件时有问题，配置文件加载失败，请查看")
	//		return ""
	//	}
	return ConfigNodeInfo.IMServerSendCustomMsgNsqTopic + "_" + fmt.Sprintf("%d", ConfigNodeInfo.IMServerSendCustomMsgNsqId)
}
