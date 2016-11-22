package httphandler

import (
	"ZDTalk/ZDTalk_http/httpprocess/client"
	logs "ZDTalk/utils/log4go"
	StringUtils "ZDTalk/utils/stringutils"
	"fmt"
	"net/http"
	"strings"
)

type MyHandler struct{}

var mux map[string]func(http.ResponseWriter, *http.Request)

func InitServlet(location string, httpPort int32) {
	logs.GetLogger().Info("------------- ZDTalk Http Init -------------")
	defer func() {
		if err := recover(); err != nil {
			logs.GetLogger().Error(err)
		}
	}()
	port := ":" + StringUtils.GetFormatString(httpPort)
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux[location] = client.Filter
	server := http.Server{
		Addr:    port,
		Handler: &MyHandler{},
	}
	err := server.ListenAndServe()
	if err != nil {
		logs.GetLogger().Info("listen ZDTalk http server error" + err.Error())
		panic(err.Error())
	}
	logs.GetLogger().Info("------------- ZDTalk Http Init 成功-------------")
}

func (*MyHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Access-Control-Allow-Origin", "*")
	params := strings.Split(request.URL.String(), "?")
	if fu, ok := mux[params[0]]; ok {
		fu(response, request)
	} else {
		fmt.Fprintf(response, "%v", "CMD不对: cmd="+params[0])
	}
}
