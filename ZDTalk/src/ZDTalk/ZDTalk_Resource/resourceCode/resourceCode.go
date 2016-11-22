package process

const (
	SUCCESS int32 = 1    //上传文件成功
	FAIL    int32 = 1039 //上传文件失败

	FAIL_MSG string = "Upload FIle fail" //上传文件失败

	UserId_IS_null   int32 = 2 //userid为空
	FileName_Is_null int32 = 3

	UserId_IS_null_Msg   string = "userId 不能为空" //userid为空
	FileName_Is_null_Msg string = "fileName 不能为空"
)
