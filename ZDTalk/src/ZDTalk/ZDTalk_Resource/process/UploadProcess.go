package process

import (
	resourceCode "ZDTalk/ZDTalk_Resource/resourceCode"
	uploadConfig "ZDTalk/ZDTalk_Resource/uploadconfig"
	log4go "ZDTalk/utils/log4go"
	stringUtils "ZDTalk/utils/stringutils"
	"archive/zip"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	//	"time"

	//	"imaging"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	//	"strconv"
	"strings"
	//	"time"
)

const (
	Upload_Dir    = "uploadFiles/"
	RESIZE_WIDTH  = 0
	RESIZE_HEIGHT = 100
)

var logs = log4go.GetLogger()
var reqres map[string]func(http.ResponseWriter, *http.Request)

type Myhandler struct{}

func InitUpload(location string, httpPort int32) {
	port := ":" + stringUtils.GetFormatString(httpPort)
	server := http.Server{
		Addr:    port,
		Handler: &Myhandler{},
		/*ReadTimeout: 50 * time.Second*/}

	reqres = make(map[string]func(http.ResponseWriter, *http.Request))
	reqres[location] = uploadFile
	reqres["/uploadFiles"] = staticFile
	server.ListenAndServe()
}

func (mh *Myhandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Access-Control-Allow-Origin", "*")

	fullUrl := strings.Split(request.URL.String(), "/")                //含参数
	locationParam := strings.Split(strings.Join(fullUrl, "/"), "?")[0] //Location
	logs.Info("完整请求地址 fullUrl ", fullUrl)
	logs.Info("地址参数 locationParam ", locationParam)  //接口后面的所有参数
	location := strings.Split(locationParam, "/")[1] //接口后面的地址 location
	logs.Info("地址参数 locationParam ", location)
	if fu, ok := reqres["/"+location]; ok {
		//if fu, ok := reqres["/"+params[1]]; ok {
		fu(response, request)
	} else {
		fmt.Println("上传文件CMD不对: cmd=" + request.URL.String())
		fmt.Fprintf(response, "%v", "上传文件CMD不对: cmd="+request.URL.String())
	}
}

func fileFilter(response http.ResponseWriter, request *http.Request) bool {
	url := request.RequestURI

	if len(url) == 0 {
		return false
	}

	logs.Debug("url is", url)
	if url[len(url)-1] == '/' {
		response.Write([]byte("404 page not found"))
		return false
	}

	return true
}

//查看目录 格式http://192.168.254.33:8017/uploadFiles/
//直接输入网络路径可下载 http://192.168.254.33:8017/uploadFiles/10001/2016-09-02%2016_09_11%20a/a.zip
func staticFile(response http.ResponseWriter, request *http.Request) {
	staticHandler := http.StripPrefix("/uploadFiles/", http.FileServer(http.Dir("uploadFiles")))

	if !fileFilter(response, request) {
		return
	}
	staticHandler.ServeHTTP(response, request)
}

func uploadFile(response http.ResponseWriter, request *http.Request) {
	logs.Info("上传文件 starting")
	defer request.Body.Close()
	if request.Method != "POST" {
		userIdparams := request.FormValue("userId")
		logs.Info("userId %d", userIdparams)
		logs.Info("上传文件 Mehtod %s", request.Method)
		fmt.Fprintf(response, "%v", "上传文件 请求方式错误 仅支持Post")
		return
	}

	urlParams := request.FormValue("url")
	userIdParams := request.FormValue("userId")
	logs.Info("userId %s", userIdParams)
	logs.Info("获取web端的url参数" + urlParams)

	//返回结构体
	s := UrlNode{}
	//参数检查
	if userIdParams == "" {
		logs.Error("userId 为空：")
		s.Code = resourceCode.UserId_IS_null
		s.ErrMsg = resourceCode.UserId_IS_null_Msg
		b, _ := json.Marshal(s)
		fmt.Fprintln(response, string(b))
		return
	}

	upFile, handler, err := request.FormFile(uploadConfig.UploadNodeInfo.ZDTalkResourceKey)
	if err != nil {
		logs.Error("接收上传文件错误：" + err.Error())
		s.Code = resourceCode.FAIL
		s.ErrMsg = err.Error()
		b, _ := json.Marshal(s)
		logs.Info("上传文件成功 result " + string(b))
		fmt.Fprintln(response, string(b))
		return
	}

	logs.Info("资源服务器得到mulFile", upFile)
	defer upFile.Close()
	logs.Info("handler.Filename " + handler.Filename)

	fileExt := filepath.Ext(handler.Filename)

	if bo, filedir := checkExt(fileExt, userIdParams, handler.Filename, upFile); bo {
		logs.Debug("上传文件成功，返回成功给客户端")
		if len(urlParams) != 0 {
			v := "<script>window.location.href='" + urlParams + filedir + "'</script>"
			s.Code = resourceCode.SUCCESS
			s.Url = v
			b, _ := json.Marshal(s)
			logs.Info("上传文件成功 result " + string(b))
			fmt.Fprintln(response, string(b))
		} else {
			s.Url = filedir
			s.Code = resourceCode.SUCCESS
			b, _ := json.Marshal(s)

			logs.Info("上传文件成功 result " + string(b))
			fmt.Fprintln(response, string(b))
		}
	} else {
		s.Url = filedir
		s.Code = resourceCode.SUCCESS
		b, _ := json.Marshal(s)
		logs.Info("上传文件成功 result " + string(b))
		fmt.Fprintln(response, string(b))
	}

}

//上传文件 Http响应 结构体
type UrlNode struct {
	Code   int32  `json:"code"`   //返回码
	ErrMsg string `json:"errMsg"` //错误信息
	Url    string `json:"url"`    //文件网络路径
}

//判断文件格式创建文件
func checkExt(fileExt, userId, fileName string, upFile multipart.File) (bool, string) {
	switch fileExt {
	case ".txt":
		return createFile(fileExt, userId, fileName, upFile)
	default:
		//目前所有文件类型都可上传，不处理文件类型
		return createFile(fileExt, userId, fileName, upFile)
		//		return false, ""
	}
}

//保存上传的文件  返回文件路径
func createFile(fileExt, userId, filename string, upFile multipart.File) (bool, string) {
	//建立文件目录
	fileNames := strings.Split(filename, ".")
	fileNameNoext := fileNames[0]
	for _, b := range fileNames {
		logs.Debug("names", b)
	}
	path := createPath(Upload_Dir, fileNameNoext, userId, fileExt)
	//读取目录
	_, err := os.Stat(path)
	if err != nil {
		logs.Info("文件夹不存在")
		err = os.MkdirAll(path, 0777) //创建所有文件夹
		if err != nil {
			logs.Info("创建新文件夹失败")
			logs.Info(err.Error())
		}
		logs.Info("已创建新文件夹")
	} else {
		logs.Info("文件夹已经存在")
	}
	//创建文件
	f, _ := os.OpenFile(path+filename, os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()
	_, errOne := io.Copy(f, upFile)
	if errOne != nil {
		return false, "创建文件失败"
	} else {
		fileDir := path + filename
		//TODO 判断是否为图片类型文件
		if strings.EqualFold(fileExt, ".jpg") || strings.EqualFold(fileExt, ".jpeg") || strings.EqualFold(fileExt, ".bmp") ||
			strings.EqualFold(fileExt, ".png") || strings.EqualFold(fileExt, ".gif") || strings.EqualFold(fileExt, ".tiff") || strings.EqualFold(fileExt, ".tif") {
			logs.Info("上传的文件为图片")
			return true, fileDir
		} else {
			if strings.EqualFold(fileExt, ".zip") {
				return true, fileDir
			} else {
				logs.Info("上传的文件为非压缩文件")
				return true, fileDir
			}
		}

	}
}

//func main() {
//	createPath("1", "a.zip", "10009", "zip")
//}

//创建文件保存目录
func createPath(path, fileName, userId, fileExt string) string {
	//ext := strings.ToLower(stringUtils.SubString(fileExt, 1, len(fileExt)-1))
	separator := "/"
	path += userId + separator
	cuTimeDir := time.Now().Format("2006-01-02 15_04_05") + "_" + fileName

	path += cuTimeDir + separator
	path = strings.Replace(path, " ", "_", -1)
	logs.Info("folder path" + path)
	return path
}

//解压文件
//srcFile 要解压的文件
//destPath 解压后存储的位置
func DeCompressZip(srcFile string, destPath string) (bool, error) {
	// 打开一个zip格式文件
	closer, err := zip.OpenReader(srcFile)
	if err != nil {
		return false, err
	}
	defer closer.Close()

	// 迭代压缩文件中的文件，打印出文件中的内容
	for _, file := range closer.File {
		logs.Info("文件名 %s:\n", file.Name)
		fileInfo := file.FileInfo()
		logs.Info(fileInfo.IsDir())
		readCloser, err := file.Open()
		if err != nil {
			logs.Error("DeCompressZip " + err.Error())
			return false, err
		}

		//获取当前目录状态
		_, err = os.Stat(destPath)
		if err != nil {
			logs.Info("文件夹不存在 创建文件夹 %s", destPath)
			os.MkdirAll(destPath, 0777)
		} else {
			logs.Info("文件夹已存在")
		}

		fileName := destPath + "/" + fileInfo.Name()

		logs.Info("文件夹已存在")
		//先创建文件
		w, err := os.Create(fileName)

		//再写入到文件中
		_, err = io.Copy(w, readCloser)
		if err != nil {
			logs.Error("DeCompressZip " + err.Error())
			return false, err
		}
		logs.Info("------------------------------------------")
		w.Close()
		readCloser.Close()
	}
	return true, nil
}

//文件重命名
//path 文件目录  pFileName文件名称
func FileListReName(path, pFileName string) {
	count := 0 //文件计数器 连接文件名称
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileExt := filepath.Ext(path)
		if strings.EqualFold(fileExt, ".jpg") || strings.EqualFold(fileExt, ".jpeg") || strings.EqualFold(fileExt, ".bmp") ||
			strings.EqualFold(fileExt, ".png") || strings.EqualFold(fileExt, ".gif") || strings.EqualFold(fileExt, ".tiff") || strings.EqualFold(fileExt, ".tif") {
			//获取当前文件的目录
			dir := filepath.Dir(path)
			//重命名文件
			fileName := pFileName + "_" + strconv.Itoa(count) + fileExt
			newPath := filepath.Join(dir, fileName)
			err := os.Rename(path, newPath)
			if err != nil {
				fmt.Println("重命名成功")
			}
			count++
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
