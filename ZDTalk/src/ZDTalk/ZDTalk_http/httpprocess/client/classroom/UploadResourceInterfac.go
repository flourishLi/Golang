// UploadResource
package classroom

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/ZDTalk_http/httpprocess/process"
	//	"ZDTalk/config"
	dbConstant "ZDTalk/constant/db"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	//	"os/exec"
	//	httpUtils "ZDTalk/utils/httputil"
	logs "ZDTalk/utils/log4go"
	timeUtils "ZDTalk/utils/timeutils"
	"path/filepath"
)

const (
	TEMP_FILE_PATH string = "tempFile/temp"
	FORM_FILE_NAME string = "uploadedfile" //必须和服务器端的formfile参数一致
)

func UpLoadResource(response http.ResponseWriter, request *bean.UploadResourceRequest, httpRequest *http.Request) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("UpLoadResource begins")

	//fileurl为空 直接返回
	if request.FileUrl == "" {
		writeErrMsg(errorcode.RESOURCE_NAME_IS_NOT_NULL, errorcode.RESOURCE_NAME_IS_NOT_NULL_MSG, response)
		return
	}
	//UserId为空 直接返回
	if request.RequestUserId == Request_User_IS_NULL {
		writeErrMsg(errorcode.REQUEST_USER_ID_CAN_NOT_NULL, errorcode.REQUEST_USER_ID_CAN_NOT_NULL_MSG, response)
		return
	}
	//FileName为空 直接返回
	if request.FileName == "" {
		writeErrMsg(errorcode.FILE_NAME_IS_NULL, errorcode.FILE_NAME_IS_NULL_MSG, response)
		return
	}
	//FileContentCount为空 直接返回
	if request.FileContentCount == Request_FileCount_IS_NULL {
		writeErrMsg(errorcode.FILE_COUNT_IS_NULL, errorcode.FILE_COUNT_IS_NULL_MSG, response)
		return
	}

	//检查用户是否具有权限
	hasAuthority, isExit := AuthorityCheckout(request.RequestUserId, Teacher)
	if isExit == User_IS_NOT_EXIT {
		writeErrMsg(errorcode.USER_ID_IS_NOT_EXIT, errorcode.USER_ID_IS_NOT_EXIT_MSG, response)
		return
	}

	if !hasAuthority {
		//没有权限
		writeErrMsg(errorcode.REQUEST_USER_ID_HAS_NO_AHTHORITY, errorcode.REQUEST_USER_ID_HAS_NO_AHTHORITY_MSG, response)
		return
	}
	//	serverREsponse := UpLoadFile(request, response, httpRequest)
	//	logs.GetLogger().Info("UpLoadFile Server Result IS:", serverREsponse)
	//上传文件成功
	//	if serverREsponse.Code == errorcode.UPLOADFILE_SUCESS {
	//process接口对象
	classRoomProcessManager := new(process.ClassRoomProcess)
	//memory接口对象
	resourceMemoryManager := memory.GetUploadResourceMemoryManager()
	resourceAdditionManager := memory.GetResourceAdditionMemoryManager()

	//返回参数 结构体(Code ErrMsg)
	result := new(bean.UploadResourceResponse) //反馈到客户端

	//系统时间
	creatTime := timeUtils.GetUnix13NowTime()
	fileName := request.FileName

	//数据库中添加上传文件
	result.Code, result.ErrMsg, result.FileId = classRoomProcessManager.UploadResource(request.RequestUserId, dbConstant.FILE_TYPE_ZIP, fileName, request.FileUrl, creatTime)
	if result.Code == errorcode.SUCCESS {
		//内存中添加资源
		result.Code, result.ErrMsg, result.FileId = resourceMemoryManager.UploadResource(result.FileId, request.RequestUserId, fileName, request.FileUrl, creatTime)
		if result.Code == errorcode.SUCCESS {
			//上传资源的同时 添加资源附属信息 fileContentCount
			result.Code, result.ErrMsg, result.FileId = classRoomProcessManager.InsertResourceAddition(result.FileId, request.FileContentCount)
			if result.Code == errorcode.SUCCESS {
				//添加到内存
				result.Code, result.ErrMsg, result.FileId = resourceAdditionManager.InsertResourceAddition(result.FileId, request.FileContentCount)
			}
		}
	}
	//	result.Url = serverREsponse.Url
	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("UpLoadResource result:" + string(datas))

	logs.GetLogger().Info("UpLoadResource end")
	fmt.Fprintf(response, string(datas))
	logs.GetLogger().Info("=============================================================\n")

}

//创建新的httprequest
//param url formFileName filePath resource
func newfileUploadRequest(uri string, paramName, path string, resource []byte) (*http.Request, error) {
	//创建临时文件
	file, err := os.Create(path)
	//将文件流写入临时文件
	ioutil.WriteFile(path, resource, 0666)
	//重新打开文件
	file, err = os.Open(path)
	//写入磁盘
	file.Sync()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", uri, body)

	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

//利用文件上传服务器 上传文件
//return 	UpLoadFile
//func UpLoadFile(request *bean.UploadResourceRequest, response http.ResponseWriter, httpRequest *http.Request) UrlNode {
//	logs.GetLogger().Info("UpLoadFile to Server Is begin")
//	//创建带有文件格式的文件名称
//	//	fileName := TEMP_FILE_PATH + filepath.Ext(request.FileName)
//	//	httpUtils.PostData()

//	//	req, errTwo := newfileUploadRequest(config.ConfigNodeInfo.UploadResourceUrl, FORM_FILE_NAME, fileName, request.FileBuff)
//	//	logs.GetLogger().Info("request.FileBuff ", string(request.FileBuff))
//	//	//创建request
//	//	if errTwo != nil {
//	//		logs.GetLogger().Info("NewRequest is error：", errTwo)
//	//	}

//	//	resp, errThree := http.DefaultClient.Do(req)
//	//	if errThree != nil {
//	//		logs.GetLogger().Info("Post is error：", errThree)
//	//	}
//	//	buf := new(bytes.Buffer)
//	//	w := multipart.NewWriter(buf)

//	logs.GetLogger().Error("转发到Resource服务 httpRequest ", httpRequest)
//	mulFile, _, err := httpRequest.FormFile("fileBuff")
//	logs.GetLogger().Error("转发到Resource服务 mulFile ", mulFile)
//	resp, err := PostData(config.ConfigNodeInfo.UploadResourceUrl, mulFile)
//	if err != nil {
//		logs.GetLogger().Error("转发到Resource服务出错 " + err.Error())
//	}
//	//获取结果
//	data, _ := ioutil.ReadAll(resp.Body)
//	resp.Body.Close()
//	logs.GetLogger().Info("UpLoadFile to Server Result IS：", string(data))

//	urlNode := UrlNode{}
//	errFour := json.Unmarshal(data, &urlNode)

//	if errFour != nil {
//		logs.GetLogger().Error("UpLoadFileResponse To Json is err:", errFour.Error())
//		writeErrMsg(errorcode.JSON_PRASEM_ERROR, errFour.Error(), response)
//	}

//	logs.GetLogger().Info("UpLoadFile to Server Result IS：", urlNode)

//	logs.GetLogger().Info("UpLoadFile to Server Is end")

//	return urlNode

//}

//func PostData(url string, reader io.Reader) (*http.Response, error) {

//	//	body := bytes.NewBuffer(data)

//	res, err := http.Post(url, "multipart/form-data", reader)

//	if err != nil {
//		return nil, err
//	}
//	return res, err
//	//	result, err := ioutil.ReadAll(res.Body)
//	//	defer res.Body.Close()
//	//	if err != nil {
//	//		return nil, err
//	//	}

//	//	return result, nil
//}

type UrlNode struct {
	Code   int32  `json:"code"`
	ErrMsg string `json:"errMsg"`
	Url    string `json:"url"`
}
