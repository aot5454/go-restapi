package user

import (
	"go-restapi/app"

	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) CreateUser(ctx app.Context, req CreateUserRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// ----------------------------

type mockUserStorage struct {
	mock.Mock
}

func (m *mockUserStorage) CreateUser(model UserModel) error {
	args := m.Called(model)
	return args.Error(0)
}

// ----------------------------

type mockUtils struct {
	mock.Mock
}

func (m *mockUtils) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockUtils) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}