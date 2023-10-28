package core


type FileUploadRepository interface {
	UploadFile(key string, file []byte) (string, error)
	DeleteFile(key string) error
}
