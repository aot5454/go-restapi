package auth

import (
	"errors"
	"go-restapi/app"
	"go-restapi/app/user"
	"go-restapi/utils"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Login(req AuthRequest) (*AuthResponse, error)
}

type authService struct {
	userStroage         user.UserStorage
	refreshTokenStorage RefreshTokenStorage
	utils               utils.Utils
}

func NewAuthService(userStroage user.UserStorage, refreshTokenStorage RefreshTokenStorage, utils utils.Utils) AuthService {
	return &authService{
		userStroage:         userStroage,
		refreshTokenStorage: refreshTokenStorage,
		utils:               utils,
	}
}

func (s *authService) Login(req AuthRequest) (*AuthResponse, error) {
	u, err := s.userStroage.GetUserByUsername(req.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if !s.utils.CheckPasswordHash(req.Password, u.Password) {
		return nil, ErrPasswordNotMatch
	}

	privateKey, err := s.utils.GetPrivateKey()
	if err != nil {
		return nil, err
	}

	tokenData := app.TokenData{
		UserID:    int(u.ID),
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      "admin",
	}

	accessToken, accessTokenExpire, err := s.utils.GetAccessToken(privateKey, tokenData, AccessTokenExpireHour)
	if err != nil {
		return nil, err
	}

	refreshToken := s.utils.GetUUID()
	refreshTokenExpire := time.Now().Add(time.Hour * RefreshTokenExpireHour).Unix()
	refreshTokenModel := RefreshTokenModel{
		UserID:    int(u.ID),
		Token:     refreshToken,
		ExpiredAt: time.Unix(refreshTokenExpire, 0),
		CreatedBy: u.Username,
		UpdatedBy: u.Username,
	}
	if err := s.refreshTokenStorage.UpsertRefreshToken(refreshTokenModel); err != nil {
		return nil, err
	}

	accessTokenExpireAt := time.Unix(accessTokenExpire, 0).Format(FormatDateTime)
	refreshTokenExpireAt := time.Unix(refreshTokenExpire, 0).Format(FormatDateTime)
	res := &AuthResponse{
		AccessToken:          accessToken,
		AccessTokenExpireAt:  accessTokenExpireAt,
		RefreshToken:         refreshToken,
		RefreshTokenExpireAt: refreshTokenExpireAt,
	}

	return res, nil
}
