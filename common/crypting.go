package common

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("AllYourBase")

const (
	ACCESS_TOKEN_SECRET  = "3234"
	ACCESS_TOKEN_AGE     = "sadas"
	REFRESH_TOKEN_SECRET = " asda"
	REFRESH_TOKEN_AGE    = "ASD"
)

type MyCustomClaims struct {
	Pair string `json:"pair"`
	jwt.StandardClaims
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

func (util *Utilities) GenerateAccessToken(pair, sub string) string {
	claims := MyCustomClaims{
		pair,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			Issuer:    "test",
			Subject:   sub,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err.Error())
	}
	return ss
}

func (util *Utilities) GenerateRefreshToken(pair, sub string) string {
	claims := MyCustomClaims{
		pair,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			Issuer:    "test",
			Subject:   sub,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err.Error())
	}
	return ss
}

func (util *Utilities) CryprtPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func (util *Utilities) ComparePasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
