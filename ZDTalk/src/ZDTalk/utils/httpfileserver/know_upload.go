package main

import (
	co "ZDTalk/utils/cryptoutils"
	su "ZDTalk/utils/stringutils"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	//	Template_Dir = "./view/"
	Upload_Dir = "files/"

//	Upload_Dir = "G:/KuGouu/"
)

var reqres map[string]func(http.ResponseWriter, *http.Request)

type Myhandler struct{}

func InitUpload() {
	logs.Debug("Upload files server starting")
	server := http.Server{
		Addr: ":8006",
		//		Addr:    ":12345",
		Handler: &Myhandler{},
		/*ReadTimeout: 50 * time.Second*/}

	reqres = make(map[string]func(http.ResponseWriter, *http.Request))
	reqres["/uploadFile"] = uploadFile
	reqres["/files"] = staticFile
	server.ListenAndServe()
}
func (mh *Myhandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	params := strings.Split(request.URL.String(), "/")
	fmt.Println(request.URL.String())
	fmt.Println(params)
	//	logs.Debug(params)
	//	logs.Debug("/" + params[1])
	if fu, ok := reqres["/"+params[1]]; ok {
		fu(response, request)
	} else {
		fmt.Println("上传文件CMD不对: cmd=" + request.URL.String())
		fmt.Fprintf(response, "%v", "上传文件CMD不对: cmd="+request.URL.String())
	}
}

func staticFile(response http.ResponseWriter, request *http.Request) {
	staticHandler := http.StripPrefix("/files/", http.FileServer(http.Dir("files")))

	staticHandler.ServeHTTP(response, request)
}

func uploadFile(response http.ResponseWriter, request *http.Request) {
	logs.Debug("接收文件starting")
	defer request.Body.Close()
	if request.Method == "GET" {
		logs.Debug("接收文件——get")
		fmt.Fprintf(response, "%v", "上传文件方法不对 不能是get方法")
	} else if request.Method == "POST" {
		logs.Debug("接收文件_post")
		urlparams := request.FormValue("url")
		logs.Info("获取web端的url参数" + urlparams)
		upFile, handler, err := request.FormFile("uploadedfile")
		if err != nil {
			logs.Error("接收上传文件错误：" + err.Error())
			fmt.Fprintf(response, "%v", "上传的异常: "+err.Error())
			return
		}
		defer upFile.Close()
		fileExt := filepath.Ext(handler.Filename)
		if bo, filedir := checkExt(fileExt, upFile); bo {
			logs.Debug("上传文件成功，返回成功给客户端")

			if len(urlparams) != 0 {
				v := "<script>window.location.href='" + urlparams + filedir + "'</script>"
				fmt.Fprintln(response, v)
			} else {
				var s UrlNode
				s.Url = filedir
				b, _ := json.Marshal(s)
				fmt.Fprintln(response, string(b))
			}
			//			fmt.Fprintf(response, "%v", "上传完成,服务器地址:"+filedir)
		} else {
			logs.Debug("上传文件失败")
			fmt.Fprintf(response, "%v", "上传文件失败")
		}
	}

}

type UrlNode struct {
	Url string
}

//判断文件格式
func checkExt(fileExt string, upFile multipart.File) (bool, string) {
	switch fileExt {
	case ".txt":
		return createFile(fileExt, upFile)
	default:
		//目前所有文件类型都可上传，不处理文件类型
		return createFile(fileExt, upFile)
		//		return false, ""
	}
}

//保存上传的文件
func createFile(fileExt string, upFile multipart.File) (bool, string) {
	fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExt
	path := createpath(Upload_Dir, fileExt)
	_, err := os.Stat(path)
	if err != nil {
		logs.Debug("文件夹不存在")
		err = os.MkdirAll(path, 0777) //创建所有文件夹
		//		err = os.MkdirAll(path, 0644) //创建所有文件夹
		//		err = os.MkdirAll(path, os.ModeDir) //创建所有文件夹
		//		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			logs.Debug("创建新文件夹失败")
			logs.Debug(err.Error())
		}
		logs.Debug("已创建新文件夹")
	} else {
		logs.Debug("文件夹已经存在")
	}
	f, _ := os.OpenFile(path+fileName, os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()
	if _, err := io.Copy(f, upFile); err != nil {
		return false, ""
	} else {
		//		fileDir, _ := filepath.Abs(path + fileName)//绝对路径
		fileDir := path + fileName
		return true, fileDir
	}
}
func createpath(path, fileExt string) string {
	ext := strings.ToLower(su.SubString(fileExt, 1, len(fileExt)-1))
	path += co.Md5Encode(strconv.Itoa(time.Now().Year())) + string(filepath.Separator)
	path += co.Md5Encode(time.Now().Month().String()) + string(filepath.Separator)
	path += co.Md5Encode(strconv.Itoa(time.Now().Day())) + string(filepath.Separator)
	path += strconv.Itoa(time.Now().Hour()) + string(filepath.Separator)
	path += ext + string(filepath.Separator)
	//	fmt.Println(path)
	return path
}
