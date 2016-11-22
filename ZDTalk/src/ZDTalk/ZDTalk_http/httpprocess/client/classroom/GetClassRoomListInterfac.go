package classroom

import (
	"ZDTalk/ZDTalk_http/bean"
	"ZDTalk/errorcode"
	"ZDTalk/manager/memory"
	logs "ZDTalk/utils/log4go"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetClassRoomList(response http.ResponseWriter, request *bean.GetClassRoomListRequest) {
	logs.GetLogger().Info("=============================================================")
	logs.GetLogger().Info("GetClassRoomList begins")
	//memory接口对象
	classRoomMemoryManager := memory.GetClassRoomMemoryManager()
	//客户端反馈结果
	result := new(bean.GetClassRoomListResponse)

	d := classRoomMemoryManager.GetClassRoomList()

	for _, v := range d {
		i := bean.ClassRoomInfo{}
		i.ClassRoomLogo = v.ClassLogo
		i.ClassRoomName = v.ClassName
		i.Description = v.Description
		i.ClassRoomId = v.ClassId

		result.Rooms = append(result.Rooms, &i)
	}

	result.Code = errorcode.SUCCESS

	datas, err1 := json.Marshal(result)
	if err1 != nil {
		logs.GetLogger().Error(err1.Error())
	}
	logs.GetLogger().Info("GetClassRoomList result:" + string(datas))
	fmt.Fprintln(response, string(datas))

	logs.GetLogger().Info("GetClassRoomList end")
	logs.GetLogger().Info("=============================================================\n")
}
