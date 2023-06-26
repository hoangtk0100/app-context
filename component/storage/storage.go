package storage

import (
	"context"
	"time"
)

type Storage interface {
	UploadFile(ctx context.Context, data []byte, key string, contentType string) (url string, storageName string, err error)
	GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
	GetPresignedURLs(ctx context.Context, keys []string, expiration time.Duration) (map[string]string, error)
	DeleteFiles(ctx context.Context, keys []string) error
}
