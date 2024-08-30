package config

import (
	"github.com/ecommerce-api/pkg/helper"
)

var (
	RedisHost string
	RedisPort string
)

func RedisLoad() {
	RedisHost = helper.GetEnv("REDIS_HOST", "127.0.0.1")
	RedisPort = helper.GetEnv("REDIS_PORT", "6379")
}
