package adapters

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chizidotdev/copia/internal/app/core"
	"io"
	"log"
)

var (
	s3BucketName = "copia-server"
)

type S3Repository struct {
	s3Client *s3.Client
}

func NewS3Repository() core.FileUploadRepository {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error loading config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &S3Repository{
		s3Client: client,
	}
}

func (s *S3Repository) UploadFile(key string, file io.Reader) (string, error) {
	uploader := manager.NewUploader(s.s3Client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(key),
		Body:   file,
		ACL:    "public-read",
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func (s *S3Repository) DeleteFile(key string) error {
	_, err := s.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(key),
	})

	return err
}

// DownloadFile TODO: Not implemented
func (s *S3Repository) DownloadFile(key string) ([]byte, error) {
	return nil, nil
}