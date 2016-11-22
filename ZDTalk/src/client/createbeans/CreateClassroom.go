package createbeans

import (
	"client/beans"
)

func CreateClassroom() beans.CreateClassroomRequest {
	requestInfo := beans.CreateClassroomRequest{}
	//requestInfo.ClassRoomLogo = "weasyr"
	requestInfo.ClassRoomName = "ffff"
	//	requestInfo.Description = "afgsdfgd"
	requestInfo.Command = "CREATE_CLASS_ROOM"
	return requestInfo
}
