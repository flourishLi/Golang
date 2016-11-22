package log4go

import (
	"ZDTalk/utils/stringutils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	l4g "code.google.com/p/log4go"
)

var Logs = GetLogger()
var logger *l4g.Logger

var separator string = string(filepath.Separator)

func GetLogger() *l4g.Logger {

	if logger == nil {
		a := make(l4g.Logger, 0)
		logger = &a
		initLog(logger)
	}
	return logger
}

func initLog(logger *l4g.Logger) {
	fileName := stringutils.GetCurFilename()
	path := "logfile" + separator + fileName + separator
	_, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0777) //创建所有文件夹
		if err != nil {
			fmt.Println("failed to create folder")
			fmt.Println(err.Error())
		}
		fmt.Println("success to create folder")
	} else {
		fmt.Println("folder exists")
	}

	logFilename := path + fileName + ".log"
	flw := l4g.NewFileLogWriter(logFilename, false)
	flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetRotate(true)
	//flw.SetRotateSize(0)
	//flw.SetRotateLines(0)
	flw.SetRotateDaily(true)
	logger.AddFilter("file", l4g.FINEST, flw)
	logger.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	//	logger.AddFilter("file", l4g.INFO, flw)
	//	logger.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
	logger.Info("create log file:" + logFilename)
	logger.Info("The Time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
}
