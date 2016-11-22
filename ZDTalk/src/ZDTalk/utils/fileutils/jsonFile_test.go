package fileutils

import (
	"fmt"
	"testing"
)

func Test_OpenFile(t *testing.T) {
	//	OpenFile()
	result := OpenJsonFileWithParam("conf.json", new(ServerList))
	if info, ok := result.(*ServerList); ok {
		for _, v := range info.ServerList {
			fmt.Println(v.ProxyIndex)
			fmt.Println(v.ProxyMainCoreIp)
			fmt.Println(v.ProxyMainGroupIp)
			fmt.Println(v.ProxyStauts)
			fmt.Println(v.ProxyWeight)
		}
	} else {
		fmt.Println("error")
	}
}
