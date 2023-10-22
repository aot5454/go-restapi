package auth

import (
	"errors"
	"go-restapi/app/user"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type testServiceSuite struct {
	suite.Suite
}

func (s *testServiceSuite) TestLogin() {

	var now = time.Now()
	var mockAccessTokenExpireAt = now.Add(time.Hour * AccessTokenExpireHour).Unix()
	var mockRefreshTokenExpireAt = now.Add(time.Hour * RefreshTokenExpireHour).Unix()

	var mockReq = AuthRequest{
		Username: "admin",
		Password: "password",
	}

	var mockUserModel = []user.UserModel{
		{
			ID:        1,
			Username:  "admin",
			Password:  "password",
			FirstName: "admin",
			LastName:  "admin",
			Status:    1,
		},
	}

	var mockAuthResponseData = &AuthResponse{
		AccessToken:          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc5Nzk2ODcsImZpcnN0bmFtZSI6IkFtaXlhIiwiaWF0IjoxNjk3OTU4MDg3LCJpc3MiOiJnby1yZXN0YXBpIiwibGFzdG5hbWUiOiJBcm1zdHJvbmciLCJyb2xlIjoiYWRtaW4iLCJzdWIiOjEzLCJ1c2VySUQiOjEzLCJ1c2VybmFtZSI6ImFkbWluIn0.RETaz4RWH1VHloJXv-NVrDY1VdgRTXbPY6dxaEXVlF2kiqsbHMdPkY8KT-wjV0k06Jwc3KtYgcSqrT4iWEa5Ej8JnoF3ag56OXkbnTH0cvdA9oTgltSgSUd1OlafUK3IPS-8XbFFSHbV-3oCN8tiUQMLAmF78DMBPB3H2B2JsegL5No0P-WUb3ZOzVEVyOXRs5sh57EXrFWQb-t-ZPRha6P-gskE0uyt2BQ3FmTNuk17yk2ssA_iGW21wdwsYAOmAbeEjaHgcyIQPOVkkMBNBvY6qlOKzjJYzYC3sScF0o-3EvTqd6UPpTPBhqY7R5mpNGD6NylE4sZI9K840aKs6w",
		AccessTokenExpireAt:  time.Unix(mockAccessTokenExpireAt, 0).Format(FormatDateTime),
		RefreshToken:         "fcd277b6-562c-49f6-8146-051bb339fb8c",
		RefreshTokenExpireAt: time.Unix(mockRefreshTokenExpireAt, 0).Format(FormatDateTime),
	}

	s.Run("Should return error when user not found", func() {
		errWant := ErrUserNotFound
		userStroage := &mockUserStorage{}
		userStroage.On("GetUserByUsername", mockReq.Username).Return(nil, gorm.ErrRecordNotFound)
		refreshTokenStorage := &mockRefreshTokenStorage{}
		utils := &mockUtils{}

		service := NewAuthService(userStroage, refreshTokenStorage, utils)
		_, err := service.Login(mockReq)
		s.ErrorIs(err, errWant)
	})

	s.Run("Should return error when other error", func() {
		errWant := errors.New("error")
		userStroage := &mockUserStorage{}
		userStroage.On("GetUserByUsername", mockReq.Username).Return(nil, errWant)
		refreshTokenStorage := &mockRefreshTokenStorage{}
		utils := &mockUtils{}

		service := NewAuthService(userStroage, refreshTokenStorage, utils)
		_, err := service.Login(mockReq)
		s.ErrorIs(err, errWant)
	})

	s.Run("Should return error when password not match", func() {
		userStroage := &mockUserStorage{}
		userStroage.On("GetUserByUsername", mockReq.Username).Return(&mockUserModel[0], nil)
		refreshTokenStorage := &mockRefreshTokenStorage{}
		utils := &mockUtils{}
		utils.On("CheckPasswordHash", mockReq.Password, mockUserModel[0].Password).Return(false)

		service := NewAuthService(userStroage, refreshTokenStorage, utils)
		_, err := service.Login(mockReq)
		s.ErrorIs(err, ErrPasswordNotMatch)
	})

	s.Run("Should return error when get private key", func() {
		errWant := errors.New("error")
		userStroage := &mockUserStorage{}
		userStroage.On("GetUserByUsername", mockReq.Username).Return(&mockUserModel[0], nil)
		refreshTokenStorage := &mockRefreshTokenStorage{}
		utils := &mockUtils{}
		utils.On("CheckPasswordHash", mockReq.Password, mockUserModel[0].Password).Return(true)
		utils.On("GetPrivateKey").Return(nil, errWant)

		service := NewAuthService(userStroage, refreshTokenStorage, utils)
		_, err := service.Login(mockReq)
		s.ErrorIs(err, errWant)
	})

	s.Run("Should return error when get access token", func() {
		errWant := errors.New("error")
		userStroage := &mockUserStorage{}
		userStroage.On("GetUserByUsername", mockReq.Username).Return(&mockUserModel[0], nil)
		refreshTokenStorage := &mockRefreshTokenStorage{}
		utils := &mockUtils{}
		utils.On("CheckPasswordHash", mockReq.Password, mockUserModel[0].Password).Return(true)
		utils.On("GetPrivateKey").Return(nil, nil)
		utils.On("GetAccessToken", mock.Anything, mock.Anything, mock.Anything).Return("xxx", int64(1), errWant)

		service := NewAuthService(userStroage, refreshTokenStorage, utils)
		_, err := service.Login(mockReq)
		s.ErrorIs(err, errWant)
	})

	s.Run("Should return error when upsert data", func() {
		errWant := errors.New("error")

		userStroage := &mockUserStorage{}
		userStroage.On("GetUserByUsername", mockReq.Username).Return(&mockUserModel[0], nil)

		refreshTokenStorage := &mockRefreshTokenStorage{}
		refreshTokenStorage.On("UpsertRefreshToken", mock.Anything).Return(errWant)

		utils := &mockUtils{}
		utils.On("CheckPasswordHash", mockReq.Password, mockUserModel[0].Password).Return(true)
		utils.On("GetPrivateKey").Return(nil, nil)
		utils.On("GetAccessToken", mock.Anything, mock.Anything, mock.Anything).Return("xxx", int64(1), nil)
		utils.On("GetUUID").Return("xxx")

		service := NewAuthService(userStroage, refreshTokenStorage, utils)
		_, err := service.Login(mockReq)
		s.ErrorIs(err, errWant)
	})

	s.Run("Should return success", func() {
		userStroage := &mockUserStorage{}
		userStroage.On("GetUserByUsername", mockReq.Username).Return(&mockUserModel[0], nil)

		refreshTokenStorage := &mockRefreshTokenStorage{}
		refreshTokenStorage.On("UpsertRefreshToken", mock.Anything).Return(nil)

		utils := &mockUtils{}
		utils.On("CheckPasswordHash", mockReq.Password, mockUserModel[0].Password).Return(true)
		utils.On("GetPrivateKey").Return(nil, nil)
		utils.On("GetAccessToken", mock.Anything, mock.Anything, mock.Anything).Return(mockAuthResponseData.AccessToken, mockAccessTokenExpireAt, nil)
		utils.On("GetUUID").Return(mockAuthResponseData.RefreshToken)

		service := NewAuthService(userStroage, refreshTokenStorage, utils)
		got, err := service.Login(mockReq)
		s.NoError(err)
		s.Equal(mockAuthResponseData, got)
	})
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(testServiceSuite))
}
