package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"porcupine/pkg/storage"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
)

var ctx = context.Background()

type MyCustomClaims struct {
	Pair string `json:"pair"`
	jwt.StandardClaims
}

type TokenService struct{}

func (ts *TokenService) DecodeJwt(token, secret string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	decodedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
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

func (ts *TokenService) SignJwt(pair, sub, secret string, expiration int64) string {
	claims := MyCustomClaims{
		pair,
		jwt.StandardClaims{
			ExpiresAt: expiration,
			Issuer:    "test",
			Subject:   sub,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println(err.Error())
	}
	return ss
}

func (ts *TokenService) IsInStore(key string) (bool, error) {
	_, err := storage.RDB.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, errors.New("Token not in store")
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (ts *TokenService) IsBlacklisted(userId, token string) (bool, error) {
	ismember, err := storage.RDB.SIsMember(ctx, "BL_"+userId, token).Result()
	if err != nil {
		return false, err
	}
	if ismember {
		return true, nil
	}
	return false, nil
}

func (ts *TokenService) PushToStore(key, token string, expiration time.Duration) error {
	value, _ := json.Marshal(map[string]interface{}{
		"refreshToken": token,
		"expiresAt":    expiration,
	})
	err := storage.RDB.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ts *TokenService) PushToBlacklist(userId, pairId, token string) error {
	// delete token from the store
	err := storage.RDB.Del(ctx, userId+"_"+pairId).Err()
	if err != nil {
		return err
	}
	// push token to the blacklist with key
	err = storage.RDB.SAdd(ctx, "BL_"+userId, token).Err()
	if err != nil {
		return err
	}
	return nil
}
