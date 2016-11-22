package db

import (
	"ZDTalk/db/database"
	"ZDTalk/manager/db/mysqldb"
)

type ResourceDbManager struct {
}

//添加上传文件资源
//Param userID fileName filePath fileTime fileThumbPath
//return fileId
func (self *ResourceDbManager) InsertResource(userId int32, fileType int32, fileName string, filePath string, fileTime int64) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ResourceDbManager{}.InsertResource(userId, fileType, fileName, filePath, fileTime)
	}
	return ZERO, nil
}

//删除上传文件资源
//Param fileID
//return int=1 删除成功
func (self *ResourceDbManager) DeleteResource(fileId int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ResourceDbManager{}.DeleteResource(fileId)
	}
	return ZERO, nil
}
