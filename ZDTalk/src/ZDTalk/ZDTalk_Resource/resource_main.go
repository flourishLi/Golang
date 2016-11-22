package main

import (
	UpLoadProcess "ZDTalk/ZDTalk_Resource/process"
	uploadConfig "ZDTalk/ZDTalk_Resource/uploadconfig"
	logs "ZDTalk/utils/log4go"
)

func main() {

	logs.GetLogger().Info("ZDTalkHttpResource Server prepare start")
	err := uploadConfig.LoadConfig("resourceConfig.json")
	if err != nil {
		logs.GetLogger().Error("加载配置文件失败 停止资源服务 " + err.Error())
		return
	}
	logs.GetLogger().Info("ZDTalkHttpResource Server start Success")
	logs.GetLogger().Info("ZDTalkHttpResource Server location --> %s , port --> %d , ZDTalkResourceKey--> %s", uploadConfig.UploadNodeInfo.ZDTalkResourceLocation, uploadConfig.UploadNodeInfo.ZDTalkResourcePort, uploadConfig.UploadNodeInfo.ZDTalkResourceKey)
	UpLoadProcess.InitUpload(uploadConfig.UploadNodeInfo.ZDTalkResourceLocation, uploadConfig.UploadNodeInfo.ZDTalkResourcePort)
}
