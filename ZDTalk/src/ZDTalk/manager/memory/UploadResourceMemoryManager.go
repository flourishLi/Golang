package memory

import (
	"ZDTalk/constant/db"
	"ZDTalk/errorcode"
	"ZDTalk/manager/db/mysqldb"
	logs "ZDTalk/utils/log4go"
	stringUtils "ZDTalk/utils/stringutils"
	"sort"
	"strconv"
	"sync"
)

//上传文件信息表
type UploadResourceMemoryInfo struct {
	FileId        int32
	UserId        int32
	FileName      string
	FilePath      string
	FileThumbPath string
	FileTime      int64
	IsDelete      int32
	FileType      int32
}

//教室信息内存管理结构体
type UploadResourceMemoryManager struct {
	UploadResourceMemoryInfos map[int32]map[int32]*UploadResourceMemoryInfo
	Lock                      sync.Mutex
}

//全局内存变量
// 对应数据库表uploadresource
var uploadResourceMemoryManager *UploadResourceMemoryManager

//初始化全局内存变量uploadResourceMemoryManager 即读取将数据表uploadresource到内存
func GetUploadResourceMemoryManager() *UploadResourceMemoryManager {
	if uploadResourceMemoryManager == nil {
		logs.Logs.Info("------------- ZDTalk uploadResource Memory Initial begin-------------")

		//初始化uploadResourceMemoryManager
		uploadResourceMemoryManager = &UploadResourceMemoryManager{}
		uploadResourceMemoryManager.UploadResourceMemoryInfos = make(map[int32]map[int32]*UploadResourceMemoryInfo)

		//加载数据库中uploadresource数据表到内存 客户端的请求操作直接去内存读取数据
		manager := new(mysqldb.ResourceDbManager)
		result, err := manager.LoadAllUpLoadResource()

		if err != nil {
			logs.GetLogger().Error("LoadAllUploadResource Error:" + err.Error())
			return nil
		}

		for _, uploadResourceDbInfo := range result {
			var resourceMap map[int32]*UploadResourceMemoryInfo
			if _, ok := uploadResourceMemoryManager.UploadResourceMemoryInfos[uploadResourceDbInfo.UserId]; !ok {
				resourceMap = make(map[int32]*UploadResourceMemoryInfo)
			} else {
				resourceMap = uploadResourceMemoryManager.UploadResourceMemoryInfos[uploadResourceDbInfo.UserId]
			}
			//根据uploadResourceDbInfo初始化uploadResourceMemoryInfo
			resource := &UploadResourceMemoryInfo{uploadResourceDbInfo.FileId, uploadResourceDbInfo.UserId, uploadResourceDbInfo.FileName, uploadResourceDbInfo.FilePath, uploadResourceDbInfo.FileThumbPath, uploadResourceDbInfo.FileTime, uploadResourceDbInfo.IsDelete, uploadResourceDbInfo.FileType}
			//设置uploadResourceMemoryManager的UploadResourceMemoryInfo属性
			resourceMap[resource.FileId] = resource
			uploadResourceMemoryManager.UploadResourceMemoryInfos[resource.UserId] = resourceMap

		}
		logs.Logs.Info("------------- ZDTalk Memory uploadResource Initial end-------------")
	}
	return uploadResourceMemoryManager
}

//获取用户的上传资源信息
//func (self *UploadResourceMemoryManager) GetUploadResource(userId int32) *UploadResourceMemoryInfo {
//	self.Lock.Lock()
//	defer self.Lock.Unlock()

//	if r, ok := self.UploadResourceMemoryInfos[userId]; ok {
//		return r
//	}

//	return nil
//}

//获取全部资源列表
func (self *UploadResourceMemoryManager) GetUploadResourceList(userId int32) []*UploadResourceMemoryInfo {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	logs.GetLogger().Info("GetUploadResourceList In Memory Success: ", self.UploadResourceMemoryInfos)
	data := []*UploadResourceMemoryInfo{}
	ok := self.IsExistUser(userId)
	if !ok {
		logs.GetLogger().Info("用户 %d 没有文件", userId)
		return data
	}
	//key 排序
	sorted_keys := make([]string, 0)
	resourceMap := self.UploadResourceMemoryInfos[userId]
	for k, _ := range resourceMap {

		sorted_keys = append(sorted_keys, stringUtils.GetFormatString(k))
	}
	sort.Strings(sorted_keys)

	for _, b := range sorted_keys {
		logs.GetLogger().Info("sortKey", b)
	}

	//排序后的切片
	for _, v := range sorted_keys {

		vInt, _ := strconv.ParseInt(v, 10, 32)
		if resourceInfo, ok := resourceMap[int32(vInt)]; ok {
			if resourceInfo.IsDelete == db.IS_DELETE_INSTERT {
				data = append(data, resourceInfo)
			}
		}
	}

	return data
}

//添加上传文件资源
//Param userID fileName url fileTime fileThumbPath
//return code errMsg fileId
func (self *UploadResourceMemoryManager) UploadResource(fileId int32, userId int32, fileName string, filePath string, fileTime int64) (int32, string, int32) {
	logs.GetLogger().Info("UploadResource In Memory begins")
	self.Lock.Lock()
	defer self.Lock.Unlock()

	r := &UploadResourceMemoryInfo{}
	r.FileId = fileId
	r.UserId = userId
	r.FileName = fileName
	r.FilePath = filePath
	//r.FileThumbPath = fileThumbPath
	r.FileTime = fileTime
	r.IsDelete = 0
	var resourceMap map[int32]*UploadResourceMemoryInfo
	if _, ok := self.UploadResourceMemoryInfos[userId]; !ok {
		resourceMap = make(map[int32]*UploadResourceMemoryInfo)
		resourceMap[fileId] = r
	} else {
		resourceMap = self.UploadResourceMemoryInfos[userId]
		resourceMap[fileId] = r
	}
	self.UploadResourceMemoryInfos[userId] = resourceMap

	logs.GetLogger().Info("UploadResource In Memory Success: ", self.UploadResourceMemoryInfos)
	logs.GetLogger().Info("UploadResource In Memory Success:", r)
	return errorcode.SUCCESS, "", r.FileId
}

//删除上传文件资源
//Param fileId
//return code errMsg
func (self *UploadResourceMemoryManager) DeleteResource(userId int32, fileId int32) (int32, string) {
	logs.GetLogger().Info("DeleteResource In Memory begins")
	self.Lock.Lock()
	defer self.Lock.Unlock()
	ok := self.IsExistUser(userId)
	logs.GetLogger().Info("DeleteResource In Memory ", self.UploadResourceMemoryInfos)
	if !ok {
		logs.GetLogger().Info("FileID IS not Exit:", fileId)
		return errorcode.RESOURCE_IS_NOT_EXIST, errorcode.RESOURCE_IS_NOT_EXIST_MSG
	}
	resourceMap := self.UploadResourceMemoryInfos[userId]
	if resourceInfo, ok := resourceMap[fileId]; !ok {
		logs.GetLogger().Info("FileID IS not Exit:", fileId)
		return errorcode.RESOURCE_IS_NOT_EXIST, errorcode.RESOURCE_IS_NOT_EXIST_MSG
	} else {
		if db.IS_DELETE_DELETE == resourceInfo.IsDelete { //已经删除
			return errorcode.RESOURCE_IS_DELETE, errorcode.RESOURCE_IS_DELETE_MSG
		} else { //删除
			resourceInfo.IsDelete = db.IS_DELETE_DELETE
		}
	}

	logs.GetLogger().Info("DeleteResource In Memory Success:")
	return errorcode.SUCCESS, "DeleteResource SUCCESS IN MEMORY"
}

//用户是否存在资源
func (self *UploadResourceMemoryManager) IsExistUser(userId int32) bool {
	_, ok := self.UploadResourceMemoryInfos[userId]
	return ok
}
