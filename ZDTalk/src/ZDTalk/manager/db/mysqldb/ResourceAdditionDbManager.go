package mysqldb

//对UploadResource表的数据库操作
import (
	"ZDTalk/db/database"
	dbInfo "ZDTalk/manager/db/info" //从数据库得到数据
	logs "ZDTalk/utils/log4go"
	"database/sql"
	//"fmt"
)

type ResourceAdditionDbManager struct {
}

//查询上传资源附属信息列表 即resourceaddition数据表的所有记录
//Param 空
func (self ResourceAdditionDbManager) LoadAllResourceAddition() (map[int32]*dbInfo.ResourceAdditionDbInfo, error) {
	logs.GetLogger().Info("开始读取数据库所有上传文件附属属性 LoadAllResourceAddition 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("数据库所有上传文件资源附属属性 错误 LoadAllResourceAddition: %s", r)
			return
		}
	}()

	sqls := "select * from resourceaddition"
	stmt, err := database.DBConn.Prepare(sqls)
	if err != nil {
		logs.GetLogger().Error("LoadAllResourceAddition Prepare Error:" + err.Error())
		return nil, err
	}
	rows, errTwo := stmt.Query()
	if errTwo != nil {
		logs.GetLogger().Error("LoadAllResourceAddition Query error:" + err.Error())
		return nil, err
	}
	defer rows.Close()
	//定义切片 存储uploadResourceAdditionDbInfo
	resourceAdditionDbMap := make(map[int32]*dbInfo.ResourceAdditionDbInfo)
	for rows.Next() {
		var fileIdSql sql.NullInt64
		var fileContentCoutSql sql.NullInt64

		rows.Scan(&fileIdSql, &fileContentCoutSql)

		resourceAdditionDbInfo := &dbInfo.ResourceAdditionDbInfo{}
		if fileIdSql.Valid {
			resourceAdditionDbInfo.FileId = int32(fileIdSql.Int64)
			//fmt.Println("fileIdSql ", uploadResourceDbInfo.FileId)
		}

		if fileContentCoutSql.Valid {
			resourceAdditionDbInfo.FileContentCount = int32(fileContentCoutSql.Int64)
			//fmt.Println("userIdSql ", uploadResourceDbInfo.UserId)
		}

		resourceAdditionDbMap[resourceAdditionDbInfo.FileId] = resourceAdditionDbInfo
	}
	return resourceAdditionDbMap, nil
}

//添加 fileid filecount
func (self ResourceAdditionDbManager) InsertResourceAddition(fileId int32, fileContentCount int32) (int64, error) {
	logs.GetLogger().Info("上传文件 InsertResourceAddition 数据库类型 mysql")
	defer func() {
		if r := recover(); r != nil {
			logs.GetLogger().Error("上传文件 错误：", r) //添加上传文件资源
			//Param userID fileName filePath fileTime fileThumbPath
			//return fileId("上传文件 错误 InsertResource: %s", r)
			return
		}
	}()
	sqlInsert := "insert into resourceaddition(FileId,FileContentCount)values(?,?)"

	stmt, err := database.DBConn.Prepare(sqlInsert)
	if err != nil {
		logs.GetLogger().Error("InsertResourceAddition Prepare error:" + err.Error())
		return ZERO, err
	}

	defer stmt.Close()

	result, errTwo := stmt.Exec(fileId, fileContentCount)
	if errTwo != nil {
		logs.GetLogger().Error("InsertResourceAddition Exec error:" + errTwo.Error())
		//fmt.Println("CreateClassRoom Exec error:" + err.Error())
		return ZERO, errTwo
	}

	id, errThree := result.LastInsertId()
	if errThree != nil {
		logs.GetLogger().Error("InsertResourceAddition LastInsertId error:" + errThree.Error())
		//fmt.Println("InsertResource LastInsertId error:" + err.Error())
		return ZERO, errThree
	}
	logs.GetLogger().Info(" InsertResourceAddition Id %d ", id)
	return id, nil

}
