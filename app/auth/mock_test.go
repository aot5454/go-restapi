package auth

import (
	"crypto/rsa"
	"go-restapi/app"
	"go-restapi/app/user"
	"go-restapi/utils"

	"github.com/stretchr/testify/mock"
)

type mockUserStorage struct {
	mock.Mock
	user.UserStorage
}

func (m *mockUserStorage) GetUserByUsername(username string) (*user.UserModel, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserModel), args.Error(1)
}

// ----------------------------

type mockRefreshTokenStorage struct {
	mock.Mock
	RefreshTokenStorage
}

func (m *mockRefreshTokenStorage) GetRefreshTokenByToken(token string) (*RefreshTokenModel, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RefreshTokenModel), args.Error(1)
}

func (m *mockRefreshTokenStorage) UpsertRefreshToken(refreshToken RefreshTokenModel) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

// ----------------------------

type mockAuthService struct {
	mock.Mock
	AuthService
}

func (m *mockAuthService) Login(req AuthRequest) (*AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AuthResponse), args.Error(1)
}

// ----------------------------

type mockUtils struct {
	mock.Mock
	utils.Utils
}

func (m *mockUtils) CheckPasswordHash(hashPwd, pwd string) bool {
	args := m.Called(hashPwd, pwd)
	return args.Bool(0)
}

func (m *mockUtils) GetPrivateKey() (*rsa.PrivateKey, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*rsa.PrivateKey), args.Error(1)
}

func (m *mockUtils) GetAccessToken(privateKey *rsa.PrivateKey, data app.TokenData, expireHour int) (string, int64, error) {
	args := m.Called(privateKey, data, expireHour)
	if args.Get(0) == nil {
		return "", 0, args.Error(2)
	}
	return args.String(0), args.Get(1).(int64), args.Error(2)
}

func (m *mockUtils) GetUUID() string {
	args := m.Called()
	return args.String(0)
}
