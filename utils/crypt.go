package utils

import (
	"log"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("AllYourBase")

type MyCustomClaims struct {
	Pair string `json:"pair"`
	jwt.StandardClaims
}

func GenerateToken(pair, sub string) string {
	claims := MyCustomClaims{
		pair,
		jwt.StandardClaims{
			ExpiresAt: 15000,
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

func CryprtPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func ComparePasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
