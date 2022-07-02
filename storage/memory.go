package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type MemoryStorage struct{}

var Ctx = context.Background()
var RDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var ctx = context.Background()

func (mem *MemoryStorage) Set(key string, value interface{}, expiration time.Duration) error {
	err := RDB.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}
