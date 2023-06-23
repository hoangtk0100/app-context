package redisdb

import (
	"context"
	"fmt"
	"strings"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/pflag"
)

const (
	defaultPoolSize     = 0 // 0 is unlimited number of socket connections
	defaultMinIdleConns = 10
)

// <user>:<password>@<host>:<port>/<db_number>
type redisDBOpt struct {
	prefix       string
	url          string
	poolSize     int
	minIdleConns int
}

type redisDB struct {
	id     string
	client *redis.Client
	logger appctx.Logger
	*redisDBOpt
}

func NewRedisDB(id, prefix string) *redisDB {
	return &redisDB{
		id: id,
		redisDBOpt: &redisDBOpt{
			prefix:       strings.TrimSpace(prefix),
			poolSize:     defaultPoolSize,
			minIdleConns: defaultMinIdleConns,
		},
	}
}

func (r *redisDB) ID() string {
	return r.id
}

func (r *redisDB) InitFlags() {
	prefix := r.prefix
	if prefix != "" {
		prefix += "-"
	}

	pflag.StringVar(&r.url,
		fmt.Sprintf("%surl", prefix),
		"redis://localhost:6379",
		"Redis connection-string - Ex: redis:<user>:<password>@<host>:<port>/<db_name>",
	)

	pflag.IntVar(&r.poolSize,
		fmt.Sprintf("%spool-size", prefix),
		defaultPoolSize,
		"Redis pool size",
	)

	pflag.IntVar(&r.minIdleConns,
		fmt.Sprintf("%spool-min-idle", prefix),
		defaultMinIdleConns,
		"Redis min idle connections",
	)
}

func (r *redisDB) isDisabled() bool {
	return r.url == ""
}

func (r *redisDB) Run(ac appctx.AppContext) error {
	if r.isDisabled() {
		return nil
	}

	r.logger = ac.Logger(r.id)

	opt, err := redis.ParseURL(r.url)
	if err != nil {
		r.logger.Error(err, "Cannot parse Redis URL")
		return err
	}

	opt.PoolSize = r.poolSize
	opt.MinIdleConns = r.minIdleConns

	client := redis.NewClient(opt)

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		r.logger.Error(err, "Cannot connect Redis ")
		return err
	}

	r.client = client

	r.logger.Infof("Connect Redis on %s", opt.Addr)

	return nil
}

func (r *redisDB) Stop() error {
	return nil
}

func (r *redisDB) GetDB() *redis.Client {
	return r.client
}
