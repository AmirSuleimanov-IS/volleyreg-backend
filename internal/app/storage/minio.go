package storage

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Storage — загрузка фото и видео участников в бакет Minio
type Storage struct {
	client *minio.Client
	bucket string
}

func New() (*Storage, error) {
	client, err := minio.New(os.Getenv("MINIO_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("MINIO_ROOT_USER"), os.Getenv("MINIO_ROOT_PASSWORD"), ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &Storage{client: client, bucket: os.Getenv("MINIO_BUCKET")}, nil
}

// Upload сохраняет файл из формы в Minio под ключом key
func (s *Storage) Upload(ctx context.Context, key string, file *multipart.FileHeader) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = s.client.PutObject(ctx, s.bucket, key, f, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	return err
}
