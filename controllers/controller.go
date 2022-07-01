package controllers

import (
	"bug-free-octo-broccoli/configs"
	"bug-free-octo-broccoli/models"
	"bug-free-octo-broccoli/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/upper/db/v4"
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

		var user models.UserDocument
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

		userId := user.ID.Hex()
		pairId := uuid.New().String()

		accessToken := utils.GenerateToken(pairId, string(userId))
		refreshToken := utils.GenerateToken(pairId, string(userId))

		key := userId + "_" + pairId
		expiration := time.Now().Add(5 * time.Minute).Unix()
		value, _ := json.Marshal(map[string]interface{}{
			"refreshToken": refreshToken,
			"expiresAt":    expiration})
		err := configs.RDB.Set(configs.Ctx, key, value, 0).Err()
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

func Me() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := ctx.Get("user")
		if !err {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "User does not exist oin context"})
			return
		}
		ctx.JSON(200, gin.H{
			"success": true,
			"data":    user})
	}
}

func Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := ctx.Get("user")
		if !err {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "User does not exist oin context"})
			return
		}

		claims, err := ctx.Get("claims")
		if !err {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "Claims not in context. Login first"})
			return
		}

		token, err := ctx.Get("token")
		if !err {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "Token not in context. Login first"})
			return
		}
		key := user.(models.UserDocument).ID.Hex() + "_" + claims.(jwt.MapClaims)["pair"].(string)
		configs.RDB.Del(configs.Ctx, key)
		configs.RDB.SAdd(configs.Ctx, "BL_"+user.(models.UserDocument).ID.Hex(), token)

		ctx.JSON(200, gin.H{
			"success": true,
			"message": "Loged out"})
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
