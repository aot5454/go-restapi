package user

import (
	"go-restapi/app"
	"go-restapi/utils"

	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) CreateUser(ctx app.Context, req CreateUserRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockUserService) GetListUser(ctx app.Context, page, pageSize int) ([]GetListUserResponse, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]GetListUserResponse), args.Error(1)
}

func (m *mockUserService) CountListUser(ctx app.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *mockUserService) GetUserByID(ctx app.Context, id int) (*GetUserResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*GetUserResponse), args.Error(1)
}

func (m *mockUserService) UpdateUser(ctx app.Context, id int, req UpdateUserRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *mockUserService) DeleteUser(ctx app.Context, id int) error {
	args := m.Called(ctx, id)
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

func (m *mockUserStorage) GetListUser(page, pageSize int) ([]UserModel, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]UserModel), args.Error(1)
}

func (m *mockUserStorage) GetUserByUsername(username string) (*UserModel, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserModel), args.Error(1)
}

func (m *mockUserStorage) GetUserByID(id int) (*UserModel, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserModel), args.Error(1)
}

func (m *mockUserStorage) CountListUser() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserStorage) UpdateUser(model UserModel) error {
	args := m.Called(model)
	return args.Error(0)
}

func (m *mockUserStorage) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// ----------------------------

type mockUtils struct {
	mock.Mock
	utils.Utils
}

func (m *mockUtils) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockUtils) CheckPasswordHash(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

func (m *mockUtils) GetPage(ctx app.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *mockUtils) GetPageSize(ctx app.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *mockUtils) GetTotalPage(total, pageSize int) int {
	args := m.Called(total, pageSize)
	return args.Int(0)
}
