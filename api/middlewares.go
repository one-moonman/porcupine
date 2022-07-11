package api

import (
	"bug-free-octo-broccoli/configs"
	"bug-free-octo-broccoli/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/upper/db/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) Login() gin.HandlerFunc {

	type ReqBody struct {
		Username string `json:"username" binding:"required"`
		Email    string ` json:"email" binding:"required"`
		Password string ` json:"password" binding:"required"`
	}

	return func(ctx *gin.Context) {
		var body ReqBody
		ctx.BindJSON(&body)

		user, err := h.repository.FindOne(db.Cond{"email": body.Email})
		if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		if !utils.ComparePasswordHash(user.Hash, body.Password) {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Invalid password credential"})
			return
		}

		ctx.Set("user", user)
	}
}

func (h *Handler) VerifyAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(403, gin.H{
				"success": false,
				"message": "No Authorization header provided"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			ctx.AbortWithStatusJSON(403, gin.H{
				"success": false,
				"message": "Could not find bearer token in Authorization header"})
			return
		}

		claims, err := utils.DecodeToken(token, configs.ACCESS_TOKEN_SECRET)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		id, _ := primitive.ObjectIDFromHex(claims["sub"].(string))
		user, err := h.repository.FindById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		ismember, err := h.memory.SIsMember("BL_"+user.ID.Hex(), token)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}
		if ismember {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": "Token is blacklisted"})
			return
		}
		ctx.Set("claims", claims)
		ctx.Set("token", token)
		ctx.Set("user", user)
	}
}

func (h *Handler) VerifyRefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(403, gin.H{
				"success": false,
				"message": "No Authorization header provided"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			ctx.AbortWithStatusJSON(403, gin.H{
				"success": false,
				"message": "Could not find bearer token in Authorization header"})
			return
		}

		claims, err := utils.DecodeToken(token, configs.REFRESH_TOKEN_SECRET)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		id, _ := primitive.ObjectIDFromHex(claims["sub"].(string))
		user, err := h.repository.FindById(id)
		if err != nil {
			ctx.AbortWithStatusJSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		userId := user.ID.Hex()
		key := userId + "_" + claims["pair"].(string)
		_, err = h.memory.Get(key)
		if err == redis.Nil {
			ctx.AbortWithStatusJSON(404, gin.H{
				"success": false,
				"message": "Token not in store"})
			return
		} else if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}
		ctx.Set("user", user)
	}
}
