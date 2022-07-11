package storage

import (
	"bug-free-octo-broccoli/configs"
	"bug-free-octo-broccoli/model"
	"fmt"
	"log"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
	"gopkg.in/mgo.v2/bson"
)

var connection = func() db.Session {
	sess, err := mongo.Open(configs.MongoOptions)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	fmt.Println("Database Connected")
	return sess
}()

func CloseJsonStorage() {
	err := connection.Close()
	if err != nil {
		log.Fatalf("db.Close(): %q\n", err)
	}
	fmt.Println("Database Disconnected")
}

type UserRepository struct{}

func (repo *UserRepository) Create(username, email, hash string) (*model.User, error) {
	user := model.User{
		ID:       bson.NewObjectId(),
		Username: username,
		Email:    email,
		Hash:     hash}
	_, err := connection.Collection("user").Insert(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindOne(condition db.Cond) (*model.User, error) {
	var user model.User
	res := connection.Collection("user").Find(condition)
	if err := res.One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindById(id interface{}) (*model.User, error) {
	var user model.User
	res := connection.Collection("user").Find(map[string]interface{}{"_id": id})
	if err := res.One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindIfExists(condition db.Cond) (bool, error) {
	res := connection.Collection("user").Find(condition)
	total, err := res.Count()
	if err != nil {
		return false, err
	}
	if total == 0 {
		return false, nil
	}
	return true, nil
}
