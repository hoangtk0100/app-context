package core

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormDBComponent interface {
	GetDB() *gorm.DB
}

type RedisDBComponent interface {
	GetDB() *redis.Client
}

type CacheComponent interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
}

type EmailComponent interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachments []string,
	) error
}

type GinComponent interface {
	GetAddress() string
	GetRouter() *gin.Engine
	Start()
}

type PubSubComponent interface {
	Publish(ctx context.Context, topic pubsub.Topic, msg *pubsub.Message) error
	Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, unsubscribeFunc func())
}

type TokenMakerComponent interface {
	CreateToken(uid string, tokenType token.TokenType, duration ...time.Duration) (string, *token.Payload, error)
	VerifyToken(token string) (*token.Payload, error)
}

type StorageComponent interface {
	UploadFile(ctx context.Context, data []byte, key string, contentType string) (string, string, error)
	GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
	GetPresignedURLs(ctx context.Context, keys []string, expiration time.Duration) (map[string]string, error)
	DeleteFiles(ctx context.Context, keys []string) error
}
