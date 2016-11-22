package uploadconfig

import (
	FileUtils "ZDTalk/utils/fileutils"
	"encoding/json"
	"fmt"
)

type UploadConfigNode struct {
	ZDTalkInterfaceUrl     string //ZDTalk Server InterfaceUrl 早道服务器完整地址
	ZDTalkResourceLocation string //ZDTalk Resource Url Location(***：8080/location)
	ZDTalkResourcePort     int32  //ZDTalk Resource Url	端口号 (***：port/location)
	ZDTalkResourceKey      string //ZDTalk 上传文件Request中的资源Key值
	ZDTalkResourcePath     string //文件存储路径
}

var UploadNodeInfo UploadConfigNode

func LoadConfig(configPath string) error {
	configBuf, err := FileUtils.ReadAll(configPath)
	if err != nil {
		return err
	}
	tmpNode := new(UploadConfigNode)
	err = json.Unmarshal(configBuf, tmpNode)
	if err != nil {
		fmt.Println("uploadConfig.json内容失败 " + err.Error())
		return err
	}
	UploadNodeInfo = *tmpNode
	return nil
}
