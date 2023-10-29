package core

import "io"

type FileUploadRepository interface {
	UploadFile(key string, file io.Reader) (string, error)
	DeleteFile(key string) error
}
