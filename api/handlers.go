package api

import (
	"bug-free-octo-broccoli/configs"
	"bug-free-octo-broccoli/model"
	"bug-free-octo-broccoli/storage"
	"bug-free-octo-broccoli/utils"

	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/upper/db/v4"
)

type Handler struct {
	// primary database, repository for users
	repository storage.UserRepository
	// secondary database, in-memory storage for tokens
	memory storage.MemoryStorage
}

func (h *Handler) Register() gin.HandlerFunc {

	type ReqBody struct {
		Username string `json:"username" binding:"required"`
		Email    string ` json:"email" binding:"required"`
		Password string ` json:"password" binding:"required"`
	}

	return func(ctx *gin.Context) {
		var body ReqBody
		ctx.BindJSON(&body)

		exist, _ := h.repository.FindIfExists(db.Cond{
			"username": body.Username,
			"email":    body.Email,
		})
		if exist {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": "Email or username already taken",
			})
			return
		}

		hash, _ := utils.CryprtPassword(body.Password)
		user, err := h.repository.Create(body.Username, body.Email, hash)
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
		user := ctx.MustGet("user").(*model.User)

		userId := user.ID.Hex()
		pairId := uuid.New().String()

		duration, _ := time.ParseDuration(configs.ACCESS_TOKEN_AGE)

		accessToken := utils.GenerateAccessToken(pairId, userId, configs.ACCESS_TOKEN_SECRET, time.Now().Add(duration).Unix())

		duration, _ = time.ParseDuration(configs.REFRESH_TOKEN_AGE)
		expiration := time.Now().Add(duration).Unix()
		refreshToken := utils.GenerateAccessToken(pairId, userId, configs.REFRESH_TOKEN_SECRET, expiration)

		key := userId + "_" + pairId
		value, _ := json.Marshal(map[string]interface{}{
			"refreshToken": refreshToken,
			"expiresAt":    expiration,
		})
		h.memory.Set(key, value, 0)

		ctx.JSON(200, gin.H{
			"success": true,
			"data": map[string]string{
				"accessToken":  accessToken,
				"refreshToken": refreshToken,
			},
		})
	}
}

func (h *Handler) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*model.User)
		claims := ctx.MustGet("user").(jwt.MapClaims)
		token := ctx.MustGet("token")

		key := user.ID.Hex() + "_" + claims["pair"].(string)
		h.memory.Del(key)
		h.memory.SAdd("BL_"+user.ID.Hex(), token)

		ctx.JSON(200, gin.H{
			"success": true,
			"message": "User:" + user.Username + " logged out",
		})
	}
}

// when using in production, the hash field from the user struct should not be returned for security
func (h *Handler) Me() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"success": true,
			"data":    ctx.MustGet("user"),
		})
	}
}
