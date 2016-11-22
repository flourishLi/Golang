package httputil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func PostData(url string, data []byte) ([]byte, error) {

	body := bytes.NewBuffer(data)
	res, err := http.Post(url, "application/json;charset=utf-8", body)

	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func PostForm(url string, data url.Values) ([]byte, error) {

	body := ioutil.NopCloser(strings.NewReader(data.Encode())) //把form数据编下码
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//这个一定要加，不加form的值post不过去，被坑了两小时
	//看下发送的结构

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	result, _ := ioutil.ReadAll(resp.Body)

	return result, nil
}

func PostMultiForm(url string, data url.Values) ([]byte, error) {

	body := ioutil.NopCloser(strings.NewReader(data.Encode())) //把form数据编下码
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "multipart/form-data; param=value")
	//这个一定要加，不加form的值post不过去，被坑了两小时
	//看下发送的结构

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	result, _ := ioutil.ReadAll(resp.Body)

	return result, nil
}

func UrlRequest(url, method string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return bytes, nil
}

//生成URL请求路径
//map为请求的参数的key-value键值对
func ParseUrl(server string, m map[string]string) string {
	u, _ := url.Parse(server)
	q := u.Query()
	for k, v := range m {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
