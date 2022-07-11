package configs

import (
	"log"
	"os"

	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
	"upper.io/db.v3/mongo"
)

var (
	ACCESS_TOKEN_SECRET  = ""
	ACCESS_TOKEN_AGE     = ""
	REFRESH_TOKEN_SECRET = ""
	REFRESH_TOKEN_AGE    = ""

	REDIS_ADDRESS  = ""
	MONGO_HOST     = ""
	MONGO_DATABASE = ""
)

var MongoOptions = mongo.ConnectionURL{
	Database: MONGO_DATABASE,
	Host:     MONGO_HOST,
}

var RedisOptions = redis.Options{
	Addr:     REDIS_ADDRESS,
	Password: "",
	DB:       0,
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	ACCESS_TOKEN_SECRET = os.Getenv("ACCESS_TOKEN_SECRET")
	ACCESS_TOKEN_AGE = os.Getenv("ACCESS_TOKEN_AGE")
	REFRESH_TOKEN_SECRET = os.Getenv("REFRESH_TOKEN_SECRET")
	REFRESH_TOKEN_AGE = os.Getenv("REFRESH_TOKEN_AGE")

	REDIS_ADDRESS = os.Getenv("REDIS_ADDRESS")
	MONGO_HOST = os.Getenv("MONGO_HOST")
	MONGO_DATABASE = os.Getenv("MONGO_DATABASE")
}
