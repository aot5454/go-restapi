package utils

import (
	"crypto/rsa"
	"go-restapi/app"
)

type Utils interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GetPage(ctx app.Context) (int, error)
	GetPageSize(ctx app.Context) (int, error)
	GetTotalPage(total, pageSize int) int
	GetAccessToken(key *rsa.PrivateKey, data app.TokenData, expireHour int) (string, int64, error)
	GetPrivateKey() (*rsa.PrivateKey, error)
	GetUUID() string
}

type utils struct{}

func NewUtils() Utils {
	return &utils{}
}
