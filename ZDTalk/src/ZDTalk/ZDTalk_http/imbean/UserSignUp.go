// CreatGroupRequest
package imbean

import (
	"ZDTalk/config"
)

//创建用户的请求
//cmd=SIGN_UP
type SignUpRequest struct {
	RequestBase
	LoginName    string `json:"loginName"`    //登录名 必须
	UserName     string `json:"userName"`     //用户昵称 必须
	Password     string `json:"password"`     //MD5加密 必须
	UserIcon     string `json:"userIcon"`     //头像
	UserId       int32  `json:"userId"`       //指定用户的userId(可选，如果此字段为空，系统会返回一个userId)
	OtherId      int32  `json:"otherId"`      //在业务服务器中的id(可选)
	AppId        int32  `json:"appId"`        //APPID与服务器保持一致即可
	IsUpdateUser bool   `json:"isUpdateUser"` //如果登录名已存在，是否要覆盖已注册用户
}

func GetSignUpRequest() *SignUpRequest {
	request := new(SignUpRequest)
	request.Command = SIGNUP
	request.AppId = config.ConfigNodeInfo.AppId
	return request
}

//创建用户的请求的响应
type SignUpResponse struct {
	ResponseBase
	UserId int32 `json:"userId"`
}

func (self SignUpRequest) GetCmd() string {
	return self.Command
}
