package createbeans

import (
	"client/beans"
)

func UpdateClassroom() beans.UpdateClassroomRequest {
	requestInfo := beans.UpdateClassroomRequest{}
	requestInfo.RoomId = 12
	//	requestInfo.ClassRoomLogo = "sdfsdd"
	//	requestInfo.Description = "afgsdfgd"
	requestInfo.Command = "UPDATE_CLASS_ROOM"
	return requestInfo
}
