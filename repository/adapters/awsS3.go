package adapters

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3BucketName = "copia-server"
)

type S3Store struct {
	s3Client *s3.Client
}

type S3StoreArgs struct {
	Region    string
	AccessKey string
	SecretKey string
}

func NewS3Store(args S3StoreArgs) *S3Store {
	cfg := aws.Config{
		Region: args.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			args.AccessKey,
			args.SecretKey,
			"",
		),
	}
	client := s3.NewFromConfig(cfg)

	return &S3Store{
		s3Client: client,
	}
}

func (s *S3Store) UploadFile(key string, file io.Reader) (string, error) {

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

func (s *S3Store) DeleteFile(key string) error {
	_, err := s.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s3BucketName),
		Key:    aws.String(key),
	})

	return err
}

// DownloadFile TODO: Not implemented
func (s *S3Store) DownloadFile(key string) ([]byte, error) {
	return nil, nil
}
