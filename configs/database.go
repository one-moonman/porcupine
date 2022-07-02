package configs

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v9"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
)

var settings = mongo.ConnectionURL{
	Database: `test`,
	Host:     `127.0.0.1`,
}

var Connection = func() db.Session {
	sess, err := mongo.Open(settings)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	fmt.Println("Database Connected")
	return sess
}()

func GetCollection(connection db.Session, collectionName string) db.Collection {
	return connection.Collection(collectionName)
}

var UserCollection = Connection.Collection("user")

var Ctx = context.Background()
var RDB = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
