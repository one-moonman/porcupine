package common

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("AllYourBase")

type MyCustomClaims struct {
	Pair string `json:"pair"`
	jwt.StandardClaims
}

func (util *Utilities) GenerateToken(pair, sub string) string {
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
