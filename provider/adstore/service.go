package adstore

import (
	"log"

	"github.com/go-redis/redis/v8"
)

func NewStorage(logg *log.Logger, client redis.Cmdable) *Storage {
	return &Storage{
		Log:    logg,
		Client: client,
	}
}
