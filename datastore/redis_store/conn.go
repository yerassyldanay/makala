package redis_store

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisConnectionParams struct {
	Host          string
	Port          int32
	LogicDatabase int32
}

func NewRedisConnection(params RedisConnectionParams) (*redis.Client, error) {
	redisFeedUrl, err := redis.ParseURL(fmt.Sprintf("redis://%s:%d/%d", params.Host, params.Port, params.LogicDatabase))
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url. err: %w", err)
	}

	inMemoryStorageFeed := redis.NewClient(redisFeedUrl)
	return inMemoryStorageFeed, err
}
