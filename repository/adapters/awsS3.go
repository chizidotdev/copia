package adapters

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/chizidotdev/shop/config"
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

	fmt.Println("Uploading file to s3", *client)
	return &S3Store{
		s3Client: client,
	}
}

func (s *S3Store) UploadFile(key string, file io.Reader) (string, error) {
	region := config.EnvVars.AWSRegion
	accessKey := config.EnvVars.AWSAccessKey
	secretKey := config.EnvVars.AWSSecretAccessKey

	cfg := aws.Config{
		Region: region,
		Credentials: credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		),
	}
	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

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
