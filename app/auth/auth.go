package auth

import (
	"errors"
	"go-restapi/app/user"
	"go-restapi/utils"
	"time"
)

const (
	RefreshTokenTableName  = "refresh_tokens"
	AccessTokenExpireHour  = 6
	RefreshTokenExpireHour = 24
	FormatDateTime         = "2006-01-02 15:04:05"
)

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken          string `json:"accessToken"`
	AccessTokenExpireAt  string `json:"accessTokenExpireAt"`
	RefreshToken         string `json:"refreshToken"`
	RefreshTokenExpireAt string `json:"refreshTokenExpireAt"`
}

var ErrUserNotFound = errors.New("user not found")
var ErrPasswordNotMatch = errors.New("password not match")

type RefreshTokenModel struct {
	ID        int       `db:"id" gorm:"primaryKey" `
	UserID    int       `db:"user_id"`
	Token     string    `db:"token" gorm:"unique"`
	ExpiredAt time.Time `db:"expired_at"`
	CreatedAt time.Time `db:"created_at" gorm:"autoCreateTime"`
	CreatedBy string    `db:"created_by" gorm:"default:'SYSTEM'"`
	UpdatedAt time.Time `db:"updated_at" gorm:"autoUpdateTime"`
	UpdatedBy string    `db:"updated_by"`
}

func New(userStorage user.UserStorage, refreshTokenStorage RefreshTokenStorage) AuthHandler {
	return NewAuthHandler(NewAuthService(userStorage, refreshTokenStorage, utils.NewUtils()))
}
