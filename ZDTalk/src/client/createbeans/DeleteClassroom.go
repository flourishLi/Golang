package createbeans

import (
	"client/beans"
)

func DeleteClassroom() beans.DeleteClassroomRequest {
	requestInfo := beans.DeleteClassroomRequest{}
	//requestInfo.RoomId = 12
	requestInfo.Command = "DELETE_CLASS_ROOM"
	return requestInfo
}
