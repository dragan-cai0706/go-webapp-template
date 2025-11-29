package shared

import (
	"github.com/redis/go-redis/v9"
)

// RedisClient Redis客户端包装
type RedisClient struct {
	*redis.Client
	keyPrefix string
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(config RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisClient{
		Client:    rdb,
		keyPrefix: config.KeyPrefix,
	}
}

// FormatKey 格式化Redis键，添加前缀
func (r *RedisClient) FormatKey(key string) string {
	if r.keyPrefix == "" {
		return key
	}
	return r.keyPrefix + ":" + key
}

