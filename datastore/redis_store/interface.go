package redis_store

import "github.com/go-redis/redis/v8"

type Cmdable interface {
	redis.Cmdable
}
