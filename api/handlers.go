package api

import (
	"bug-free-octo-broccoli/common"
	"bug-free-octo-broccoli/model"
	"bug-free-octo-broccoli/storage"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/upper/db/v4"
)

type Handler struct {
	repository storage.UserRepository
	utils      common.Utilities
}

func (h *Handler) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqBody := h.utils.BindBody(ctx)
		exist, _ := h.repository.FindIfExists(db.Cond{
			"username": reqBody["username"],
			"email":    reqBody["email"],
		})
		if exist {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "Email or username already taken",
			})
			return
		}

		hash, _ := h.utils.CryprtPassword(reqBody["password"])

		user, err := h.repository.Create(reqBody["username"], reqBody["email"], hash)
		if err != nil {
			ctx.JSON(500, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"success": true,
			"data":    user,
		})
	}
}

func (h *Handler) GenerateTokens() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, onContext := ctx.Get("user")
		if !onContext {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "User does not exist oin context"})
			return
		}
		userId := user.(model.User).ID.Hex()
		pairId := uuid.New().String()

		accessToken := h.utils.GenerateToken(pairId, string(userId))
		refreshToken := h.utils.GenerateToken(pairId, string(userId))

		key := userId + "_" + pairId
		expiration := time.Now().Add(5 * time.Minute).Unix()
		value, _ := json.Marshal(map[string]interface{}{
			"refreshToken": refreshToken,
			"expiresAt":    expiration})
		err := storage.RDB.Set(storage.Ctx, key, value, 0).Err()
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

func (h *Handler) Logout() gin.HandlerFunc {
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
		key := user.(model.User).ID.Hex() + "_" + claims.(jwt.MapClaims)["pair"].(string)
		storage.RDB.Del(storage.Ctx, key)
		storage.RDB.SAdd(storage.Ctx, "BL_"+user.(model.User).ID.Hex(), token)

		ctx.JSON(200, gin.H{
			"success": true,
			"message": "Loged out"})
	}
}

func (h *Handler) Me() gin.HandlerFunc {
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
