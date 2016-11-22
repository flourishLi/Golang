package mysqldb

//对UploadResource表的数据库操作
import (
	"ZDTalk/db/database"
	dbInfo "ZDTalk/manager/db/info" //从数据库得到数据
	logs "ZDTalk/utils/log4go"
	"database/sql"
	//"fmt"
)

type ResourceDbManager struct {
}

//查询上传资源列表列表 即uploadresource数据表的所有记录
//Param 空
//return UploadResourceDbInfo切片 error
func (self ResourceDbManager) LoadAllUpLoadResource() (map[int32]*dbInfo.UploadResourceDbInfo, error) {
	logs.GetLogger().Info("开始读取数据库所有上传文件资源 LoadAllUpLoadResource 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("数据库所有上传文件资源 错误 LoadGroups: %s", r)
			return
		}
	}()

	sqls := "select * from uploadresource where IsDelete=0"
	stmt, err := database.DBConn.Prepare(sqls)
	if err != nil {
		logs.GetLogger().Error("LoadAllUpLoadResource Prepare Error:" + err.Error())
		return nil, err
	}
	rows, errTwo := stmt.Query()
	if errTwo != nil {
		logs.GetLogger().Error("LoadAllUpLoadResource Query error:" + err.Error())
		return nil, err
	}
	defer rows.Close()
	//定义切片 存储uploadResourceDbInfo
	uploadResourceDbSlice := make(map[int32]*dbInfo.UploadResourceDbInfo)
	for rows.Next() {
		var userIdSql sql.NullInt64
		var fileIdSql sql.NullInt64
		var fileNameSql sql.NullString
		var filePathSql sql.NullString
		var fileThumbPathSql sql.NullString
		var fileTimeSql sql.NullInt64
		var isDeleteSql sql.NullInt64
		var fileTypeSql sql.NullInt64

		rows.Scan(&fileIdSql, &userIdSql, &fileNameSql, &filePathSql, &fileThumbPathSql, &fileTimeSql, &isDeleteSql, &fileTypeSql)

		uploadResourceDbInfo := &dbInfo.UploadResourceDbInfo{}
		if fileIdSql.Valid {
			uploadResourceDbInfo.FileId = int32(fileIdSql.Int64)
			//fmt.Println("fileIdSql ", uploadResourceDbInfo.FileId)
		}

		if userIdSql.Valid {
			uploadResourceDbInfo.UserId = int32(userIdSql.Int64)
			//fmt.Println("userIdSql ", uploadResourceDbInfo.UserId)
		}
		if fileNameSql.Valid {
			uploadResourceDbInfo.FileName = fileNameSql.String
			//fmt.Println("fileNameSql ", uploadResourceDbInfo.FileName)
		}
		if filePathSql.Valid {
			uploadResourceDbInfo.FilePath = filePathSql.String
			//fmt.Println("filePathSql ", uploadResourceDbInfo.FilePath)
		}
		if fileThumbPathSql.Valid {
			uploadResourceDbInfo.FileThumbPath = fileThumbPathSql.String
			//fmt.Println("fileThumbPathSql ", uploadResourceDbInfo.FileThumbPath)
		}
		if fileTimeSql.Valid {
			uploadResourceDbInfo.FileTime = fileTimeSql.Int64
			//fmt.Println("fileTimeSql ", uploadResourceDbInfo.FileTime)
		}
		if isDeleteSql.Valid {
			uploadResourceDbInfo.IsDelete = int32(isDeleteSql.Int64)
			//fmt.Println("isDeleteSql ", uploadResourceDbInfo.IsDelete)
		}
		if fileTypeSql.Valid {
			uploadResourceDbInfo.FileType = int32(fileTypeSql.Int64)
			//fmt.Println("isDeleteSql ", uploadResourceDbInfo.IsDelete)
		}

		uploadResourceDbSlice[uploadResourceDbInfo.FileId] = uploadResourceDbInfo
	}
	return uploadResourceDbSlice, nil
}

func (self ResourceDbManager) InsertResource(userId int32, fileType int32, fileName string, filePath string, fileTime int64) (int64, error) {
	logs.GetLogger().Info("上传文件 InsertResource 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("上传文件 错误：", r) //添加上传文件资源
			//Param userID fileName filePath fileTime fileThumbPath
			//return fileId("上传文件 错误 InsertResource: %s", r)
			return
		}
	}()
	sqlInsert := "insert into uploadresource(UserId,FileType,FileName,FilePath,FileTime,IsDelete)values(?,?,?,?,?,?)"

	stmt, err := database.DBConn.Prepare(sqlInsert)
	if err != nil {
		logs.GetLogger().Error("InsertResource Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(userId, fileType, fileName, filePath, fileTime, ZERO)
	if errTwo != nil {
		logs.GetLogger().Error("InsertResource Exec error:" + errTwo.Error())
		//fmt.Println("CreateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	id, errThree := result.LastInsertId()
	if errThree != nil {
		logs.GetLogger().Error("InsertResource LastInsertId error:" + errThree.Error())
		//fmt.Println("InsertResource LastInsertId error:" + err.Error())
		return ZERO, errThree
	}
	logs.GetLogger().Info("上传文件 InsertResource Id %d ", id)
	return id, nil

}

//删除资源
//Param fileId
//Return int=1 删除成功 err
func (self ResourceDbManager) DeleteResource(fileId int32) (int64, error) {
	logs.GetLogger().Info("开始删除文件 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("更新教室状态 错误: %s", r)
			return
		}
	}()

	sqlUpdate := "update uploadresource set IsDelete = ? where FileId = ?"

	stmt, err := database.DBConn.Prepare(sqlUpdate)
	if err != nil {
		logs.GetLogger().Error("DeleteResource Prepare error:" + err.Error())
		//fmt.Println("DeleteResource Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(One, fileId)
	if errTwo != nil {
		logs.GetLogger().Error("DeleteResource Exec error:" + errTwo.Error())
		//fmt.Println("UpdateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	affect, errThree := result.RowsAffected()
	if errThree != nil {
		logs.GetLogger().Error("DeleteResource RowsAffected error:" + errThree.Error())
		//fmt.Println("DeleteResource RowsAffected error:" + err.Error())
		return ZERO, errThree
	}
	logs.GetLogger().Info("删除文件结果 DeleteResource 影响行数 %d ", affect)
	return affect, nil
}
