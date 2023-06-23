package common

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GormDBComponent interface {
	GetDB() *gorm.DB
}

type RedisDBComponent interface {
	GetDB() *redis.Client
}
