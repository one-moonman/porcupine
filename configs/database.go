package configs

import (
	"fmt"
	"log"

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
	fmt.Println("connected")
	return sess
}()

func GetCollection(connection db.Session, collectionName string) db.Collection {
	return connection.Collection(collectionName)
}
