package db

import (
	"ZDTalk/db/database"
	"ZDTalk/manager/db/mysqldb"
)

type ResourceSubDbManager struct {
}

//添加上传文件资源
//Param userID fileName filePath fileTime fileThumbPath
//return fileId
func (self *ResourceSubDbManager) InsertResource(userId int32, fileId int32, subFileName string, subFilePath string, subFileType int32, subFileTime int64, isDelete int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ResourceSubDbManager{}.InsertSubResource(userId, fileId, subFileName, subFilePath, subFileType, subFileTime, isDelete)
	}
	return ZERO, nil
}

//删除上传文件资源
//Param fileID
//return int=1 删除成功
func (self *ResourceSubDbManager) DeleteResource(fileId int32, isDelete int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ResourceSubDbManager{}.DeleteSubResource(fileId, isDelete)
	}
	return ZERO, nil
}
