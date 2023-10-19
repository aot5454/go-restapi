package utils

type Utils interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type utils struct{}

func NewUtils() Utils {
	return &utils{}
}
