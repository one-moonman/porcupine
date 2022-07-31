package storage

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"
)

//var ctx = context.Background()
var RDB *redis.Client

func ConnectCacheStorage(opt *redis.Options) {
	RDB = redis.NewClient(opt)
	fmt.Println("Memory connected")
}

func CloseCacheStorage() error {
	err := RDB.Close()
	if err != nil {
		return err
	}
	fmt.Println("Memory Disconnected")
	return nil
}

var ctx = context.Background()

// type CacheStorage struct{}

// func (cache *CacheStorage) Set(key string, value interface{}, expiration time.Duration) error {
// 	err := rdb.Set(ctx, key, value, expiration).Err()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (cache *CacheStorage) Get(key string) (string, error) {
// 	value, err := rdb.Get(ctx, key).Result()
// 	if err != nil {
// 		return "", err
// 	}
// 	return value, nil
// }

// // better error handling
// func (cache *CacheStorage) Del(key string) *redis.IntCmd {
// 	return rdb.Del(ctx, key)
// }

// // better error handling
// func (cache *CacheStorage) SAdd(key string, value interface{}) *redis.IntCmd {
// 	return rdb.SAdd(ctx, key, value)
// }

// func (cache *CacheStorage) SIsMember(key string, value interface{}) (bool, error) {
// 	ismember, err := rdb.SIsMember(ctx, key, value).Result()
// 	if err != nil {
// 		return false, err
// 	}
// 	return ismember, nil
// }
