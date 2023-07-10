package core

import (
	"context"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
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
	CreateToken(tokenType token.TokenType, uid string, duration ...time.Duration) (string, *token.Payload, error)
	VerifyToken(token string) (*token.Payload, error)
}

type StorageComponent interface {
	UploadFile(ctx context.Context, data []byte, key string, contentType string) (url string, storageName string, err error)
	GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
	GetPresignedURLs(ctx context.Context, keys []string, expiration time.Duration) (map[string]string, error)
	DeleteFiles(ctx context.Context, keys []string) error
}

type GRPCServerComponent interface {
	WithAddress(address string)
	WithListener(lis net.Listener)
	WithServerOptions(serverOpts ...grpc.ServerOption)
	WithServeMuxOptions(muxOpts ...runtime.ServeMuxOption)
	WithUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)
	WithStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)
	GetLogger() appctx.Logger
	GetServer() *grpc.Server
	GetGateway() *runtime.ServeMux
	Start(ctx context.Context)
}

type GRPCClientComponent interface {
	WithPrefix(prefix string)
	WithAddress(address string)
	GetAddress() string
	GetLogger() appctx.Logger
	Dial(options ...grpc.DialOption) *grpc.ClientConn
	DialContext(ctx context.Context, options ...grpc.DialOption) *grpc.ClientConn
}
