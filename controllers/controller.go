package controllers

import (
	"bug-free-octo-broccoli/configs"
	"bug-free-octo-broccoli/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = configs.GetCollection(configs.Connection, "user")

func Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody map[string]string

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		res := userCollection.Find(db.Cond{
			"username": reqBody["username"],
			"email":    reqBody["email"]})

		if total, _ := res.Count(); total != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Email or username already taken"})
			return
		}

		hash, _ := cryprtPassword(reqBody["password"])

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
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		res := userCollection.Find(db.Cond{"email": reqBody["email"]})

		var user models.User
		if err := res.One(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		if !comparePasswordHash(user.Hash, reqBody["password"]) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid password credential"})
			return
		}

		// generate token

	}
}

func cryprtPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func comparePasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateAccessToken() {

}
