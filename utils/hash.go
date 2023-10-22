package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *utils) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *utils) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *utils) GetUUID() string {
	return uuid.New().String()
}
