package storage

import (
	"fmt"
	"log"
	"porcupine/pkg/config"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
)

// var Connection db.Session

// func ConnectJsonStorage(settings db.ConnectionURL) error {
// 	var err error
// 	Connection, err = mongo.Open(config.MongoOptions)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Database Connected")
// 	return nil
// }

// func CloseJsonStorage() error {
// 	err := Connection.Close()
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Database DisConnected")
// 	return nil
// }

var connection = func() db.Session {
	sess, err := mongo.Open(config.MongoOptions)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	fmt.Println("Database Connected")
	return sess
}()

func GetCollection(collectionName string) db.Collection {
	return connection.Collection(collectionName)
}
