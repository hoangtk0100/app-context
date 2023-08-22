package storage

import (
	"bytes"
	"context"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type s3Storage struct {
	id            string
	name          string
	client        *s3.Client
	presignClient *s3.PresignClient
	logger        appctx.Logger
	*storageOpt
}

func NewS3Storage(id string) *s3Storage {
	return &s3Storage{
		id:         id,
		name:       "S3",
		storageOpt: new(storageOpt),
	}
}

func (storage *s3Storage) ID() string {
	return storage.id
}

func (storage *s3Storage) InitFlags() {
	pflag.StringVar(&storage.accessKey, "storage-access-key", "", "Cloud storage access key")
	pflag.StringVar(&storage.secretKey, "storage-secret-key", "", "Cloud storage secret key")
	pflag.StringVar(&storage.region, "storage-region", "", "Cloud storage region")
	pflag.StringVar(&storage.bucketName, "storage-bucket", "", "Cloud storage bucket name")
	pflag.StringVar(&storage.endPoint, "storage-end-point", "", "Cloud storage end point")
	pflag.StringVar(&storage.domain, "storage-domain", "", "Cloud storage domain")
}

func (storage *s3Storage) Run(ac appctx.AppContext) error {
	storage.logger = ac.Logger(storage.id)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(storage.region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				storage.accessKey,
				storage.secretKey,
				"",
			),
		),
	)

	if err != nil {
		storage.logger.Fatal(err, ErrCannotSetupStorage.Error())
	}

	storage.client = s3.NewFromConfig(cfg)
	storage.presignClient = s3.NewPresignClient(storage.client)

	storage.logger.Info("Setup storage: ", storage.id)
	return nil
}

func (storage *s3Storage) Stop() error {
	return nil
}

func (storage *s3Storage) UploadFile(ctx context.Context, data []byte, key string, contentType string) (url string, storageName string, err error) {
	fileBytes := bytes.NewReader(data)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(storage.bucketName),
		Key:         aws.String(key),
		Body:        fileBytes,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACL(types.BucketCannedACLPrivate),
	}

	_, err = storage.client.PutObject(ctx, params)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	url = fmt.Sprintf("%s/%s", storage.domain, key)
	return url, storage.name, nil
}

func (storage *s3Storage) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return getPresignedURL(ctx, storage.presignClient, storage.bucketName, key, expiration)
}

func (storage *s3Storage) GetPresignedURLs(ctx context.Context, keys []string, expiration time.Duration) (map[string]string, error) {
	return getPresignedURLs(ctx, storage.presignClient, storage.bucketName, keys, expiration)
}

func (storage *s3Storage) DeleteFiles(ctx context.Context, keys []string) error {
	return deleteFiles(ctx, storage.client, storage.bucketName, keys)
}
