package storage

import (
	"context"

	"github.com/go-redis/redis/v9"
)

var Ctx = context.Background()
var RDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
