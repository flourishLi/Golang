// GetAllUsers
package bean

type GetAllUserRequest struct {
	ClientBaseRequest //CMD=GET_ALLUSERS 含有UserId
}

type GetAllUserResponse struct {
	ClientBaseResponse
	Data []UserInfo `json:"data"`
}
