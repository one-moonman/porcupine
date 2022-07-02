package common

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Utilities struct{}

func (util *Utilities) BindBody(ctx *gin.Context) map[string]string {
	var reqBody map[string]string

	if err := ctx.BindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		ctx.Abort()
	}
	return reqBody
}

func (util *Utilities) ExtractAuthBearerToken(ctx *gin.Context) string {
	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" {
		ctx.String(http.StatusForbidden, "No Authorization header provided")
		ctx.Abort()
		return ""
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == auth {
		ctx.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
		ctx.Abort()
		return ""
	}
	return token
}

func (util *Utilities) DecodeToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	decodedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		return nil, err
	}
	if !decodedToken.Valid {
		err := errors.New("token not valid")
		return nil, err
	}
	return claims, nil
}
