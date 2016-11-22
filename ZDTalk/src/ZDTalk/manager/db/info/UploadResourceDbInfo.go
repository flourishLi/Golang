package info

//上传资源Db 结构体
type UploadResourceDbInfo struct {
	FileId        int32
	UserId        int32
	FileName      string
	FilePath      string
	FileThumbPath string
	FileTime      int64
	IsDelete      int32
	FileType      int32
}
type SubResourceDbInfo struct {
	SubFileId   int32 //子文件名称
	FileId      int32 //父文件名称
	UserId      int32
	SubFileName string
	SubFilePath string
	SubFileTime int64
	SubFileType int32

	IsDelete int32 //0未被删除1已删除
}

type ResourceAdditionDbInfo struct {
	FileId           int32
	FileContentCount int32
}
