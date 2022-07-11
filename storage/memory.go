package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v9"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func CloseMemoryStorage() {
	err := rdb.Close()
	if err != nil {
		log.Fatalf("mem.Close(): %q\n", err)
	}
	fmt.Println("Memory Disconnected")
}

type MemoryStorage struct{}

func (mem *MemoryStorage) Set(key string, value interface{}, expiration time.Duration) error {
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (mem *MemoryStorage) Get(key string) (string, error) {
	value, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (mem *MemoryStorage) Del(key string) *redis.IntCmd {
	return rdb.Del(ctx, key)
}

func (mem *MemoryStorage) SAdd(key string, value interface{}) *redis.IntCmd {
	return rdb.SAdd(ctx, key, value)
}

func (mem *MemoryStorage) SIsMember(key string, value interface{}) (bool, error) {
	ismember, err := rdb.SIsMember(ctx, key, value).Result()
	if err != nil {
		return false, err
	}
	return ismember, nil
}
