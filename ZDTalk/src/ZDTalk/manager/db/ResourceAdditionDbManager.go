package db

import (
	"ZDTalk/db/database"
	dbInfo "ZDTalk/manager/db/info"

	"ZDTalk/manager/db/mysqldb"
)

type ResourceAdditionDbManager struct {
}

//查询上传资源附属信息列表 即resourceaddition数据表的所有记录
//Param 空
func (self *ResourceAdditionDbManager) LoadAllResourceAddition() (map[int32]*dbInfo.ResourceAdditionDbInfo, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ResourceAdditionDbManager{}.LoadAllResourceAddition()
	}
	return nil, nil
}

//添加资源
//Param fileID filecontentcount
//return fileId
func (self *ResourceAdditionDbManager) InsertResourceAddition(fileId int32, fileContentCount int32) (int64, error) {
	if database.DBType == database.MYSQL {
		return mysqldb.ResourceAdditionDbManager{}.InsertResourceAddition(fileId, fileContentCount)
	}
	return ZERO, nil
}
