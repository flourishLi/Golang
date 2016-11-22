package config

import (
	"ZDTalk/utils/fileutils"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
)

type DBConfig struct {
	DriverName   string
	UserName     string
	UserPassword string
	HostAddress  string
	Port         int32
	DBName       string
	DBType       string
}

const (
	DB_MYSQL = "mysql"
	DB_MSSQL = "mssql"
)

var DBInfo DBConfig

func LoadDBConf(config string) bool {

	d, err := fileutils.ReadAll(config)

	if err != nil {
		logs.GetLogger().Error("load db config file:" + config + " error " + err.Error())
		return false
	}

	dbConfig := DBConfig{}

	err = json.Unmarshal(d, &dbConfig)

	if err != nil {
		logs.GetLogger().Error("load db config file:" + config + " error " + err.Error())
		return false
	}
	DBInfo = dbConfig
	return true
}
