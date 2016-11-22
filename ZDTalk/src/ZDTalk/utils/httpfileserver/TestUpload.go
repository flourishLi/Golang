package main

import (
	"code.google.com/p/log4go"
	//	"fmt"
	//	co "clcong/utils/cryptoutils"
	l4g "ZDTalk/utils/log4go"

	//	"net/http"
	//	"os"
	//	"path/filepath"
	//	"strconv"
	//	"time"
)

var logs *log4go.Logger = l4g.GetLogger()

func main() {
	logs.Debug("server starting")

	InitUpload()
	//	up.MainServer()
	//	path := "G:/KuGouu/"
	//		createpath(path)
	//	filePath(path)
}

//func staticServer() {
//	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("files"))))
//	http.ListenAndServe(":2235", nil)
//}
//func filePath(path string) {

//	_, err := os.Stat(path)
//	if err == nil {
//		fmt.Println("文件夹不存在")
//		//		path1 := path + strconv.Itoa(time.Now().Year()) + string(filepath.Separator) + time.Now().Month().String() + string(filepath.Separator) + strconv.Itoa(time.Now().Day())
//		path1 := createPathForDownLoad(path)
//		fmt.Println(path)
//		fmt.Println(path1)
//		os.MkdirAll(path1, os.ModeDir)

//		//		filePath(path)
//		//				filePath("G:/KuGouu")
//	} else {
//		fmt.Println("文件夹存在")
//		fmt.Println(path)
//	}
//}
//func createPathForDownLoad(path string) string {
//	fileName := strconv.FormatInt(time.Now().UnixNano(), 10)
//	path += co.Md5Encode(strconv.Itoa(time.Now().Year())) + string(filepath.Separator)
//	path += co.Md5Encode(time.Now().Month().String()) + string(filepath.Separator)
//	path += co.Md5Encode(strconv.Itoa(time.Now().Day()+1)) + string(filepath.Separator)
//	fmt.Println(path + fileName)
//	return path
//}
