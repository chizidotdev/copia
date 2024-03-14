package util

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"mime/multipart"
)

func GenerateRandString(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.RawURLEncoding.EncodeToString(b)

	return state, nil
}

type ParseImageResult struct {
	File        io.Reader
	Name        string
	ContentType string
}

func ParseImage(image *multipart.FileHeader) (ParseImageResult, error) {
	//file := imgString[strings.IndexByte(imgString, ',')+1:]
	var result ParseImageResult
	result.Name = image.Filename
	result.ContentType = image.Header.Get("Content-Type")

	file, err := image.Open()
	if err != nil {
		return result, err
	}

	result.File = file
	return result, nil
}
