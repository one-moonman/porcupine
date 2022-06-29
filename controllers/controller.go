package controllers

import (
	"bug-free-octo-broccoli/configs"
	"bug-free-octo-broccoli/models"
	"bug-free-octo-broccoli/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/upper/db/v4"
	"gopkg.in/mgo.v2/bson"
)

var userCollection = configs.GetCollection(configs.Connection, "user")
var Expiration int64 = 15000

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody map[string]string

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		res := userCollection.Find(db.Cond{
			"username": reqBody["username"],
			"email":    reqBody["email"]})

		if total, _ := res.Count(); total != 0 {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "Email or username already taken"})
			return
		}

		hash, _ := utils.CryprtPassword(reqBody["password"])

		user := models.User{
			Username: reqBody["username"],
			Email:    reqBody["email"],
			Hash:     hash}

		_, err := userCollection.Insert(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    user})
	}
}

type User struct {
	ID       interface{} `bson:"_id" json:"_id"`
	Username string      `bson:"username" json:"username"`
	Email    string      `bson:"email" json:"email"`
	Hash     string      `bson:"hash" json:"hash"`
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody map[string]string

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		res := userCollection.Find(db.Cond{"email": reqBody["email"]})

		var user User
		if err := res.One(&user); err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		if !utils.ComparePasswordHash(user.Hash, reqBody["password"]) {
			ctx.JSON(401, gin.H{
				"success": false,
				"message": "Invalid password credential"})
			return
		}

		userId := user.ID.(bson.ObjectId).Hex()
		pairId := uuid.New().String()

		accessToken := utils.GenerateToken(pairId, string(userId))
		refreshToken := utils.GenerateToken(pairId, string(userId))

		key := userId + "_" + pairId

		err := configs.RDB.Set(configs.Ctx, key, map[string]interface{}{
			"refreshToken": refreshToken,
			"expiresAt":    Expiration}, 0).Err()
		if err != nil {
			fmt.Println(err.Error())
		}

		ctx.JSON(200, gin.H{
			"success": true,
			"data": map[string]string{
				"accessToken":  accessToken,
				"refreshToken": refreshToken}})
	}
}

func Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// extract from header
	}
}

func RefreshTokens() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// extract from header
	}
}

func DeleteAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// extract from header
	}
}
