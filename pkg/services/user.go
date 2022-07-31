package services

import (
	"porcupine/pkg/models"
	"porcupine/pkg/storage"

	"github.com/upper/db/v4"
	"gopkg.in/mgo.v2/bson"
)

var userCollection = storage.GetCollection("user")

type UserService struct{}

func (us *UserService) Create(username, email, hash string) (*models.User, error) {
	user := models.User{
		ID:       bson.NewObjectId(),
		Username: username,
		Email:    email,
		Hash:     hash,
	}
	_, err := userCollection.Insert(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) FindOne(condition db.Cond) (*models.User, error) {
	var user models.User
	res := userCollection.Find(condition)
	if err := res.One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) FindById(id interface{}) (*models.User, error) {
	var user models.User
	res := userCollection.Find(map[string]interface{}{"_id": id})
	if err := res.One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) FindIfExists(condition db.Cond) (bool, error) {
	res := userCollection.Find(condition)
	total, err := res.Count()
	if err != nil {
		return false, err
	}
	if total == 0 {
		return false, nil
	}
	return true, nil
}
