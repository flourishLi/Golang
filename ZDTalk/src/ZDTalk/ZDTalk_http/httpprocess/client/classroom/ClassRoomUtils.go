// AuthorityCheckout
package classroom

import (
	"ZDTalk/ZDTalk_http/bean"
	//	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
	"net/http"
)

//检查用户是否有相应的权限
//Param   requestUserId roleId 1 2 3 4学生 6教师
//return true 有权限 false  没有权限 1 用户存在 0 用户不存在
func AuthorityCheckout(requestUserId, roleId int32) (bool, int32) {
	//Memory接口对象 用户表
	userInfoMemoryManager := memory.GetUserInfoMemoryManager()
	users := userInfoMemoryManager.Users
	if userInfo, ok := users[requestUserId]; ok { //用户存在
		if roleId == userInfo.Role {
			return true, User_IS_EXIT
		} else { //用户不具有权限
			return false, User_IS_EXIT
		}

	} else { //用户不存在
		return false, User_IS_NOT_EXIT
	}
}

//获取教室的classroomImid
func GetClassRoomIMId(classRoomId int32) int32 {
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()

	//获取教室id对应的ClassRoomIMid
	if r, ok := classRoomMemoryManager.ClassRooms[classRoomId]; ok {
		return r.ClassRoomIMId
	} else {
		return 0
	}
}

//获取用户的imid
func GetUserIMId(userId int32) int32 {
	//memory接口对象
	userInfoMemoryManager := memory.GetUserInfoMemoryManager()

	//获取教室id对应的ClassRoomIMid
	if r, ok := userInfoMemoryManager.Users[userId]; ok {
		return r.ChatId
	} else {
		return 0
	}
}

func writeErrMsg(code int32, errMsg string, response http.ResponseWriter) {
	result := new(bean.ClientBaseResponse)
	result.Code = code
	result.ErrMsg = errMsg
	son, _ := json.Marshal(result)
	logs.GetLogger().Error("Http接口 请求失败 result %s", string(son))
	fmt.Fprintf(response, string(son))
}
