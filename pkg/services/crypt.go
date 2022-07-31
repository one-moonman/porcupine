package services

import "golang.org/x/crypto/bcrypt"

type CryptService struct{}

func (cs *CryptService) CryprtPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func (cs *CryptService) ComparePasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
