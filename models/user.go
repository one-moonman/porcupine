package models

type User struct {
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Hash     string `bson:"hash" json:"hash"`
}
