package images

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioServiceInterface interface {
	UploadFile(file []byte, contentType, id string) (string, error)
	DeleteFile(path string) error
}

type MinioService struct {
	Client     *minio.Client
	BucketName string
}

func NewMinioService(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (MinioServiceInterface, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("Bucket %s successfully created", bucketName)
	}

	return &MinioService{Client: client, BucketName: bucketName}, nil
}

func (m *MinioService) UploadFile(file []byte, contentType, id string) (string, error) {
	imageUUID := uuid.New().String()
	filePath := fmt.Sprintf("%s/%s", id, imageUUID)

	// Преобразуем []byte в io.Reader для передачи в MinIO
	reader := bytes.NewReader(file)

	// Загрузка файла в MinIO
	_, err := m.Client.PutObject(
		context.Background(),
		m.BucketName,
		filePath,
		reader,
		int64(len(file)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", err
	}

	log.Printf("File successfully uploaded to %s", filePath)
	return filePath, nil
}

func (m *MinioService) DeleteFile(filePath string) error {
	filePath = strings.TrimPrefix(filePath, "/images/")
	err := m.Client.RemoveObject(context.Background(), m.BucketName, filePath, minio.RemoveObjectOptions{})
	if err != nil {
		log.Printf("Error deleting file %s: %v", filePath, err)
		return err
	}

	_, err = m.Client.StatObject(context.Background(), m.BucketName, filePath, minio.StatObjectOptions{})
	if err == nil {
		log.Printf("File %s still exists after deletion attempt", filePath)
		return fmt.Errorf("file %s still exists after deletion attempt", filePath)
	}
	log.Printf("File %s successfully deleted", filePath)
	return nil
}
