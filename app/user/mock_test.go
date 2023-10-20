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

func (m *mockUserService) GetListUser(ctx app.Context) ([]GetListUserResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]GetListUserResponse), args.Error(1)
}

// ----------------------------

type mockUserStorage struct {
	mock.Mock
}

func (m *mockUserStorage) CreateUser(model UserModel) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *mockUserStorage) GetListUser() ([]UserModel, error) {
	args := m.Called()
	return args.Get(0).([]UserModel), args.Error(1)
}

func (m *mockUserStorage) GetUserByUsername(username string) (*UserModel, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserModel), args.Error(1)
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
