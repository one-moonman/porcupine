package storage

import (
	"bug-free-octo-broccoli/model"
	"fmt"
	"log"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mongo"
	"gopkg.in/mgo.v2/bson"
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

type UserRepository struct{}

func (repo *UserRepository) Create(username, email, hash string) (*model.User, error) {
	user := model.User{
		ID:       bson.NewObjectId(),
		Username: username,
		Email:    email,
		Hash:     hash}
	_, err := Connection.Collection("user").Insert(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindOne(condition db.Cond) (*model.User, error) {
	var user model.User
	res := Connection.Collection("user").Find(condition)
	if err := res.One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindById(id interface{}) (*model.User, error) {
	var user model.User
	res := Connection.Collection("user").Find(map[string]interface{}{"_id": id})
	if err := res.One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) FindIfExists(condition db.Cond) (bool, error) {
	res := Connection.Collection("user").Find(condition)
	total, err := res.Count()
	if err != nil {
		return false, err
	}
	if total == 0 {
		return false, nil
	}
	return true, nil
}
