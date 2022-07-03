package model

import "gopkg.in/mgo.v2/bson"

type User struct {
	ID       bson.ObjectId `bson:"_id" json:"_id"`
	Username string        `bson:"username" json:"username" binding:"required"`
	Email    string        `bson:"email" json:"email" binding:"required"`
	Hash     string        `bson:"hash" json:"hash" binding:"required"`
}
