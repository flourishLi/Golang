package mysqldb

//子资源 数据库管理 (暂未使用，属于冗余代码)

import (
	"ZDTalk/db/database"
	logs "ZDTalk/utils/log4go"

	//"fmt"
)

type ResourceSubDbManager struct {
}

//添加子资源
func (self ResourceSubDbManager) InsertSubResource(userId int32, fileId int32, subFileName string, subFilePath string, subFileType int32, subFileTime int64, isDelete int32) (int64, error) {
	logs.GetLogger().Info("上传文件 Subresource 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("上传文件 错误：", r) //添加上传文件资源
			//Param userID fileName filePath fileTime fileThumbPath
			//return fileId("上传文件 错误 InsertResource: %s", r)
			return
		}
	}()
	sqlInsert := "insert into subresource(UserId,FileId,SubFileName,SubFilePath,SubFileType,SubFileTime,IsDelete)values(?,?,?,?,?,?,?)"

	stmt, err := database.DBConn.Prepare(sqlInsert)
	if err != nil {
		logs.GetLogger().Error("Subresource Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(userId, fileId, subFileName, subFilePath, subFileTime, subFileType, isDelete)
	if errTwo != nil {
		logs.GetLogger().Error("Subresource Exec error:" + errTwo.Error())
		//fmt.Println("CreateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	id, errThree := result.LastInsertId()
	if errThree != nil {
		logs.GetLogger().Error("Subresource LastInsertId error:" + errThree.Error())
		//fmt.Println("InsertResource LastInsertId error:" + err.Error())
		return ZERO, errThree
	}
	logs.GetLogger().Info("添加子资源结果 InsertResource Id %d ", id)
	return id, nil
}

////添加子资源
//Param fileId
//Return int=1 删除成功 err
func (self ResourceSubDbManager) DeleteSubResource(fileId int32, isDelete int32) (int64, error) {
	logs.GetLogger().Info("开始删除文件 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("更新教室状态 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update subResource set IsDelete = ? where FileId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("DeleteResource Prepare error:" + err.Error())
		//fmt.Println("DeleteResource Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(isDelete, fileId)
	if errTwo != nil {
		logs.GetLogger().Error("DeleteSubResource Exec error:" + errTwo.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("DeleteResource RowsAffected error:" + errThree.Error())
		//fmt.Println("DeleteResource RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	logs.GetLogger().Info("DB 删除文件结果 DeleteSubResource 影响行数 %d ", affect)
	return affect, nil
}
