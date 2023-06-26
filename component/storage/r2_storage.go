package storage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type storageOpt struct {
	bucketName string
	region     string
	accessKey  string
	secretKey  string
	endPoint   string
	domain     string
}

type r2Storage struct {
	id            string
	name          string
	client        *s3.Client
	presignClient *s3.PresignClient
	logger        appctx.Logger
	*storageOpt
}

func NewR2Storage(id string) *r2Storage {
	return &r2Storage{
		id:         id,
		name:       "R2",
		storageOpt: new(storageOpt),
	}
}

func (storage *r2Storage) ID() string {
	return storage.id
}

func (storage *r2Storage) InitFlags() {
	pflag.StringVar(&storage.accessKey, "storage-access-key", "", "Cloud storage access key")
	pflag.StringVar(&storage.secretKey, "storage-secret-key", "", "Cloud storage secret key")
	pflag.StringVar(&storage.region, "storage-region", "", "Cloud storage region")
	pflag.StringVar(&storage.bucketName, "storage-bucket", "", "Cloud storage bucket name")
	pflag.StringVar(&storage.endPoint, "storage-end-point", "", "Cloud storage end point")
	pflag.StringVar(&storage.domain, "storage-domain", "", "Cloud storage domain")
}

func (storage *r2Storage) Run(ac appctx.AppContext) error {
	storage.logger = ac.Logger(storage.id)
	r2Resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: storage.endPoint,
			}, nil
		})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				storage.accessKey,
				storage.secretKey,
				"",
			),
		),
	)

	if err != nil {
		storage.logger.Fatal(err, "Cannot setup storage storage")
	}

	storage.client = s3.NewFromConfig(cfg)
	storage.presignClient = s3.NewPresignClient(storage.client)

	storage.logger.Info("Setup storage storage : ", storage.id)
	return nil
}

func (storage *r2Storage) Stop() error {
	return nil
}

func (storage *r2Storage) UploadFile(ctx context.Context, data []byte, key string, contentType string) (url string, storageName string, err error) {
	fileBytes := bytes.NewReader(data)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(storage.bucketName),
		Key:         aws.String(key),
		Body:        fileBytes,
		ContentType: aws.String(contentType),
	}

	_, err = storage.client.PutObject(ctx, params)
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	url = fmt.Sprintf("%s/%s", storage.domain, key)
	return url, storage.name, nil
}

func (storage *r2Storage) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return getPresignedURL(ctx, storage.presignClient, storage.bucketName, key, expiration)
}

func (storage *r2Storage) GetPresignedURLs(ctx context.Context, keys []string, expiration time.Duration) (map[string]string, error) {
	return getPresignedURLs(ctx, storage.presignClient, storage.bucketName, keys, expiration)
}

func (storage *r2Storage) DeleteFiles(ctx context.Context, keys []string) error {
	return deleteFiles(ctx, storage.client, storage.bucketName, keys)
}

func deleteFiles(ctx context.Context, client *s3.Client, bucketName string, keys []string) error {
	params := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{
			Objects: toOIDs(keys),
		},
	}

	_, err := client.DeleteObjects(ctx, params)
	return errors.WithStack(err)
}

func toOIDs(keys []string) []types.ObjectIdentifier {
	oids := make([]types.ObjectIdentifier, len(keys))
	for index := 0; index < len(oids); index++ {
		oids[index] = types.ObjectIdentifier{
			Key: &(keys[index]),
		}
	}

	return oids
}

func getPresignedURL(ctx context.Context, client *s3.PresignClient, bucketName string, key string, expiration time.Duration) (string, error) {
	req, err := client.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", errors.WithStack(err)
	}

	return req.URL, nil
}

func getPresignedURLs(ctx context.Context, client *s3.PresignClient, bucketName string, keys []string, expiration time.Duration) (map[string]string, error) {
	urls := make(map[string]string)

	for _, key := range keys {
		req, err := client.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = expiration
		})

		if err != nil {
			return nil, errors.WithStack(err)
		}

		urls[key] = req.URL
	}

	return urls, nil
}
