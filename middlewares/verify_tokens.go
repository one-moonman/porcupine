package middlewares

import (
	"bug-free-octo-broccoli/configs"
	"bug-free-octo-broccoli/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = configs.GetCollection(configs.Connection, "user")

func VerifyAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.String(http.StatusForbidden, "No Authorization header provided")
			ctx.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			ctx.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
			ctx.Abort()
			return
		}

		claims := jwt.MapClaims{}
		decodedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("AllYourBase"), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.String(http.StatusUnauthorized, err.Error())
				ctx.Abort()

				return
			}
			ctx.String(http.StatusBadRequest, err.Error())
			ctx.Abort()

			return
		}
		if !decodedToken.Valid {
			ctx.String(http.StatusBadRequest, "token not valid")
			ctx.Abort()
			return
		}

		id, _ := primitive.ObjectIDFromHex(claims["sub"].(string))
		var user models.UserDocument
		res := userCollection.Find(map[string]interface{}{"_id": id})
		if err := res.One(&user); err != nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}
		ismember, err := configs.RDB.SIsMember(configs.Ctx, "BL_"+user.ID.Hex(), decodedToken.Raw).Result()
		if err != nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}
		fmt.Println(ismember)
		if ismember {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": "Token is blacklisted"})
			ctx.Abort()
			return
		}
		ctx.Set("claims", claims)
		ctx.Set("token", token)

		ctx.Set("user", user)
	}
}

func VerifyRefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.String(http.StatusForbidden, "No Authorization header provided")
			ctx.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			ctx.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
			ctx.Abort()
			return
		}

		claims := jwt.MapClaims{}
		decodedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("AllYourBase"), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.String(http.StatusUnauthorized, err.Error())
				ctx.Abort()

				return
			}
			ctx.String(http.StatusBadRequest, err.Error())
			ctx.Abort()

			return
		}
		if !decodedToken.Valid {
			ctx.String(http.StatusBadRequest, "token not valid")
			ctx.Abort()
			return
		}

		id, _ := primitive.ObjectIDFromHex(claims["sub"].(string))
		var user models.UserDocument
		res := userCollection.Find(map[string]interface{}{"_id": id})
		if err := res.One(&user); err != nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}

		userId := user.ID.Hex()
		key := userId + "_" + claims["pair"].(string)
		_, err = configs.RDB.Get(configs.Ctx, key).Result()
		if err == redis.Nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": "token not in tokenstore"})
			ctx.Abort()
			return
		} else if err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
	}
}
