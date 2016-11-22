package memory

import (
	"ZDTalk/errorcode"
	"ZDTalk/manager/db/mysqldb"
	logs "ZDTalk/utils/log4go"
	"sync"
)

//资源附属信息表
type ResourceAdditionMemoryInfo struct {
	FileId           int32
	FileContentCount int32
}

//资源附属信息内存管理结构体
type ResourceAdditionMemoryManager struct {
	ResourceAdditionMemoryInfos map[int32]*ResourceAdditionMemoryInfo
	Lock                        sync.Mutex
}

//全局内存变量
// 对应数据库表resourceaddition
var resourceAdditionMemoryManager *ResourceAdditionMemoryManager

//初始化全局内存变量uploadResourceMemoryManager 即读取将数据表uploadresource到内存
func GetResourceAdditionMemoryManager() *ResourceAdditionMemoryManager {
	if resourceAdditionMemoryManager == nil {
		logs.Logs.Info("------------- ZDTalk ResourceAddition Memory Initial begin-------------")
		//初始化resourceAdditionMemoryManager
		resourceAdditionMemoryManager = &ResourceAdditionMemoryManager{}
		resourceAdditionMemoryManager.ResourceAdditionMemoryInfos = make(map[int32]*ResourceAdditionMemoryInfo)

		//加载数据库中resourceaddition数据表到内存 客户端的请求操作直接去内存读取数据
		manager := new(mysqldb.ResourceAdditionDbManager)
		result, err := manager.LoadAllResourceAddition()

		if err != nil {
			logs.GetLogger().Error("LoadResourceAddition Error:" + err.Error())
			return nil
		}

		for _, resourceAdditionDbInfo := range result {
			//根据uploadResourceDbInfo初始化uploadResourceMemoryInfo
			memoryInfo := &ResourceAdditionMemoryInfo{resourceAdditionDbInfo.FileId, resourceAdditionDbInfo.FileContentCount}
			//设置uploadResourceMemoryManager的UploadResourceMemoryInfo属性
			resourceAdditionMemoryManager.ResourceAdditionMemoryInfos[memoryInfo.FileId] = memoryInfo

		}
		logs.Logs.Info("------------- ZDTalk Memory ResourceAddition Initial end-------------")
	}
	return resourceAdditionMemoryManager
}

//获取某一文件资源附属信息
func (resource *ResourceAdditionMemoryManager) GetResourceAddition(fileId int32) *ResourceAdditionMemoryInfo {
	resource.Lock.Lock()
	defer resource.Lock.Unlock()

	if r, ok := resource.ResourceAdditionMemoryInfos[fileId]; ok {
		return r
	}
	return nil
}

//获取全部资源列表
func (resource *ResourceAdditionMemoryManager) GetResourceAdditionList() []*ResourceAdditionMemoryInfo {
	resource.Lock.Lock()
	defer resource.Lock.Unlock()

	data := []*ResourceAdditionMemoryInfo{}

	for _, v := range resource.ResourceAdditionMemoryInfos {
		data = append(data, v)
	}

	return data
}

//添加上传文件附属信息
//Param fileid fileCOntentcount
//return code errMsg fileId
func (resource *ResourceAdditionMemoryManager) InsertResourceAddition(fileId, fileContentCount int32) (int32, string, int32) {
	logs.GetLogger().Info("InsertResourceAddition In Memory begins")
	resource.Lock.Lock()
	defer resource.Lock.Unlock()

	r := &ResourceAdditionMemoryInfo{}
	r.FileId = fileId
	r.FileContentCount = fileContentCount
	resource.ResourceAdditionMemoryInfos[r.FileId] = r
	logs.GetLogger().Info("InsertResourceAddition In Memory Success:", r)
	return errorcode.SUCCESS, "", r.FileId
}
