package bean

//	"ZDTalk/manager/db/info"

type GetClassRoomInfoRequest struct {
	ClientBaseRequest       //CMD=GET_CLASS_ROOM_INFO
	ClassRoomId       int32 `json:"classRoomId"`
}

type GetClassRoomInfoResponse struct {
	ClientBaseResponse             //code=1标识成功 3：此名字的教室已经存在
	ClassRoomId         int32      `json:"classRoomId"`
	ClassRoomIMId       int32      `json:"classRoomIMId"`
	CreatorUserId       int32      `json:"creatorUserId"`
	ClassRoomStatus     int32      `json:"classRoomStatus"`
	SettingStatus       []int32    `json:"settingStatus"`
	CreateTime          int64      `json:"CceateTime"`
	ClassRoomName       string     `json:"classRoomName"`
	ClassRoomLogo       string     `json:"classRoomLogo"`
	Description         string     `json:"introduction"`
	ClassRoomCourse     string     `json:"classRoomCourse"`
	MemberList          []UserInfo `json:"memberList"`          //成员Id集合
	OnLineMemberList    []UserInfo `json:"onLineMemberList"`    //在线成员Id集合
	HandMemberList      []UserInfo `json:"handMemberList"`      //举手成员Id集合
	ForbidSayMemberList []UserInfo `json:"forbidSayMemberList"` //被禁言的成员Id集合
	ForbidHandStatus    int32      `json:"forbidHandStatus"`    //教室的禁止举手状态 0可举手 1禁止举手

	SayingMemberList []UserInfo `json:"sayingMemberList"` //正在发言的成员Id集合
}
