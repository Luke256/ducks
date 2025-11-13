package s3

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Luke256/ducks/utils/compressor"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Storage struct {
	client     *s3.Client
	bucketName string
}

func NewS3Storage(bucketName string) (*S3Storage, error) {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	accessKey := os.Getenv("STORAGE_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_SECRET_KEY")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &S3Storage{
		client:     s3Client,
		bucketName: bucketName,
	}, nil
}

func (s *S3Storage) UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	rawImage, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer rawImage.Close()

	compressedImage, format, err := compressor.CompressImage(rawImage)
	if err != nil {
		return "", err
	}
	defer compressedImage.Close()

	fileID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s.%s", fileID.String(), strings.ToLower(format))

	uploader := manager.NewUploader(s.client)

	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   compressedImage,
		ContentType: aws.String("image/" + format),
	})
	if err != nil {
		return "", err
	}

	// ファイルのキャッシュ
	cachedImage, err := os.Create(filepath.Join(os.TempDir(), fileName))
	if err != nil {
		return "", err
	}
	defer cachedImage.Close()

	_, err = compressedImage.Seek(0, 0)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(cachedImage, compressedImage); err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *S3Storage) DeleteFile(fileName string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return err
	}

	cachedFilePath := filepath.Join(os.TempDir(), fileName)
	if _, err := os.Stat(cachedFilePath); err == nil {
		if err := os.Remove(cachedFilePath); err != nil {
			return err
		}
	}
	return nil
}

func (s *S3Storage) GetFileURL(fileName string) string {
	endpoint := os.Getenv("STORAGE_ENDPOINT")
	return fmt.Sprintf("%s/%s/%s", endpoint, s.bucketName, fileName)
}