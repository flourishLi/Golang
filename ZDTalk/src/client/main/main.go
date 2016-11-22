package main

import (
	"client/createbeans"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	IP   = "test2.clcong.com"
	PORT = "8012"
	PATH = "interface"
)

func main() {
	SendRequest(jsonFormat())
}

//生成请求需要的json数据
func jsonFormat() string {

	info := createbeans.DeleteClassroom()
	b, _ := json.Marshal(info)
	return string(b)
}

//发送HTTP请求并打印响应结果
func SendRequest(parmas string) {

	//处理参数
	url := "http://" + IP + ":" + PORT + "/" + PATH
	fmt.Println("请求地址路径：", url)
	body := strings.NewReader(parmas)
	fmt.Println(body)

	//发送请求
	client := new(http.Client)
	req, _ := http.NewRequest("POST", url, body)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求返回响应中的错误：", err.Error())
		return
	}
	defer resp.Body.Close()

	//解析响应
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("解析返回响应结果错误：", err.Error())
		return
	}
	fmt.Println(string(bytes))
}
