package utils

import "go-restapi/app"

type Utils interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
	GetPage(ctx app.Context) (int, error)
	GetPageSize(ctx app.Context) (int, error)
	GetTotalPage(total, pageSize int) int
}

type utils struct{}

func NewUtils() Utils {
	return &utils{}
}
