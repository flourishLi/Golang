package database

import (
	config "ZDTalk/config"
	log4go "ZDTalk/utils/log4go"
	"database/sql"
	"fmt"
	//_ "go-mssqldb"
	_ "mysql"
	"os"
	"sync"
	"time"
	//	_ "github.com/mattn/go-sqlite3"
)

var logs = log4go.GetLogger()

const (
	MYSQL  = 1
	MSSQL  = 2
	SQLITE = 3
)

var lock sync.Mutex

var DBType int
var DBConn *sql.DB

func OpenDb() {
	var drivername string = ""
	var confs string = ""

	if !config.LoadDBConf("db.json") {
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}
	dbConf := config.DBInfo

	if dbConf.DBType == config.DB_MYSQL {
		DBType = MYSQL
		drivername = dbConf.DriverName
		//	confs := "root:root@tcp(localhost:3306)/test"
		confs = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConf.UserName,
			dbConf.UserPassword, dbConf.HostAddress, dbConf.Port, dbConf.DBName)

	} else if dbConf.DBType == config.DB_MSSQL {
		DBType = MSSQL
		drivername = dbConf.DriverName
		//	confs := "server=127.0.0.1;user id=sa;password=root123123;database=gosmc;port=1433;datasource=SQLEXPRESS"
		confs = fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;port=%d",
			dbConf.HostAddress, dbConf.UserName, dbConf.UserPassword, dbConf.DBName, dbConf.Port)
	} else if dbConf.DBType == "sqlite3" {
		DBType = SQLITE
		drivername = dbConf.DriverName
		confs = fmt.Sprintf("%s/%s", dbConf.HostAddress, dbConf.DBName)
	}
	//	fmt.Println("DB--drivername:", drivername)
	//	fmt.Println("DB--conf:", confs)

	logs.Info("DB--drivername:", drivername)
	logs.Info("DB--conf:", confs)
	db, err := sql.Open(drivername, confs)

	if err != nil {
		logs.Error("加载数据库信息错误：", err.Error())

		time.Sleep(1 * time.Second)

		os.Exit(11)

	} else {
		logs.Info("加载数据库信息成功:")
	}
	DBConn = db
	DBConn.SetMaxOpenConns(200)
	DBConn.SetMaxIdleConns(50)
}
