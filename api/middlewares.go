package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/upper/db/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqBody := h.utils.BindBody(ctx)

		user, err := h.repository.FindOne(db.Cond{"email": reqBody["email"]})
		if err != nil {
			ctx.JSON(400, gin.H{
				"success": false,
				"message": err.Error()})
			return
		}

		if !h.utils.ComparePasswordHash(user.Hash, reqBody["password"]) {
			ctx.JSON(401, gin.H{
				"success": false,
				"message": "Invalid password credential"})
			return
		}

		ctx.Set("user", user)
	}
}

func (h *Handler) VerifyAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := h.utils.ExtractAuthBearerToken(ctx)

		claims, err := h.utils.DecodeToken(token)
		if err != nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}

		id, _ := primitive.ObjectIDFromHex(claims["sub"].(string))
		user, err := h.repository.FindById(id)
		if err != nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}

		ismember, err := h.memory.SIsMember("BL_"+user.ID.Hex(), token)
		if err != nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}
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

func (h *Handler) VerifyRefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := h.utils.ExtractAuthBearerToken(ctx)

		claims, err := h.utils.DecodeToken(token)
		if err != nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}

		id, _ := primitive.ObjectIDFromHex(claims["sub"].(string))
		user, err := h.repository.FindById(id)
		if err != nil {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": err.Error()})
			ctx.Abort()
			return
		}

		userId := user.ID.Hex()
		key := userId + "_" + claims["pair"].(string)
		_, err = h.memory.Get(key)
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
