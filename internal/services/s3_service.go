package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	appConfig "ecommerce-gin/internal/config" // rename this OR using Aliasing

	awsConfig "github.com/aws/aws-sdk-go-v2/config" // rename this

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

func InitS3() {

	cfg := appConfig.Cfg

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: cfg.S3Endpoint,
			}, nil
		},
	)

	awsCfg, _ := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(cfg.S3Region),
		awsConfig.WithEndpointResolverWithOptions(customResolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3Key,
			cfg.S3Secret,
			"",
		)),
	)

	s3Client = s3.NewFromConfig(awsCfg)
}

func UploadToS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	newName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

	// Determine content type based on extension
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".webp":
		contentType = "image/webp"
	case ".gif":
		contentType = "image/gif"
	}

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(appConfig.Cfg.S3Bucket),
		Key:         aws.String(newName),
		Body:        file,
		ACL:         "public-read",
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", err
	}

	return newName, nil
}
