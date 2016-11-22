package user

import (
	"ZDTalk/ZDTalk_http/bean"
	"encoding/json"
	"fmt"
	"net/http"
)

func writeErrMsg(code int32, errMsg string, response http.ResponseWriter) {
	result := new(bean.ClientBaseResponse)
	result.Code = code
	result.ErrMsg = errMsg
	son, _ := json.Marshal(result)
	fmt.Fprintf(response, string(son))
}
