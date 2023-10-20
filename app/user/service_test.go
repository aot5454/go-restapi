package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var mockCreateUserRequest = CreateUserRequest{
	Username:  "test",
	Password:  "password",
	FirstName: "test",
	LastName:  "test",
}

var mockGetListUserResponse = []GetListUserResponse{
	{
		ID:        1,
		Username:  "test",
		FirstName: "test",
		LastName:  "test",
		Status:    "Active",
	},
	{
		ID:        2,
		Username:  "test2",
		FirstName: "test2",
		LastName:  "test2",
		Status:    "Inactive",
	},
}

var mockUserModel = []UserModel{
	{
		ID:        1,
		Username:  "test",
		Password:  "password",
		FirstName: "test",
		LastName:  "test",
		Status:    1,
	},
	{
		ID:        2,
		Username:  "test2",
		Password:  "password2",
		FirstName: "test2",
		LastName:  "test2",
		Status:    0,
	},
}

func TestCreateUserService(t *testing.T) {
	t.Run("Should return success", func(tc *testing.T) {
		storage := &mockUserStorage{}
		storage.On("GetUserByUsername", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
		storage.On("CreateUser", mock.Anything).Return(nil)
		utils := &mockUtils{}
		utils.On("HashPassword", mock.Anything).Return("password", nil)

		service := NewUserService(storage, utils)

		err := service.CreateUser(nil, mockCreateUserRequest)
		assert.NoError(tc, err)
	})

	t.Run("Should return error (Other error)", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("GetUserByUsername", mock.Anything).Return(nil, errors.New("error"))
		storage.On("CreateUser", mock.Anything).Return(nil)
		utils := &mockUtils{}

		service := NewUserService(storage, utils)

		err := service.CreateUser(nil, mockCreateUserRequest)
		assert.Error(t, err)
	})

	t.Run("Should return error (Dup data)", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("GetUserByUsername", mock.Anything).Return(&mockUserModel[0], nil)
		storage.On("CreateUser", mock.Anything).Return(nil)
		utils := &mockUtils{}

		service := NewUserService(storage, utils)

		err := service.CreateUser(nil, mockCreateUserRequest)
		assert.Error(t, err)
	})

	t.Run("Should return error (HashPassword)", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("GetUserByUsername", mock.Anything).Return(nil, nil)
		storage.On("CreateUser", mock.Anything).Return(nil)
		utils := &mockUtils{}
		utils.On("HashPassword", mock.Anything).Return("", errors.New("error"))

		service := NewUserService(storage, utils)

		err := service.CreateUser(nil, mockCreateUserRequest)
		assert.Error(t, err)
	})
}

func TestGetListUserService(t *testing.T) {
	t.Run("Should return success", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("GetListUser", mock.Anything, mock.Anything).Return(mockUserModel, nil)
		utils := &mockUtils{}

		service := NewUserService(storage, utils)

		got, err := service.GetListUser(nil, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, len(mockUserModel), len(got))
		assert.EqualValues(t, mockGetListUserResponse, got)
	})

	t.Run("Should return error", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("GetListUser", mock.Anything, mock.Anything).Return([]UserModel{}, errors.New("error"))
		utils := &mockUtils{}

		service := NewUserService(storage, utils)

		_, err := service.GetListUser(nil, 1, 10)
		assert.Error(t, err)
	})
}

func TestCountListUserService(t *testing.T) {
	t.Run("Should return success", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("CountListUser").Return(int64(1), nil)
		utils := &mockUtils{}

		service := NewUserService(storage, utils)

		got, err := service.CountListUser(nil)
		assert.NoError(t, err)
		assert.Equal(t, 1, got)
	})

	t.Run("Should return error", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("CountListUser").Return(int64(0), errors.New("error"))
		utils := &mockUtils{}

		service := NewUserService(storage, utils)

		_, err := service.CountListUser(nil)
		assert.Error(t, err)
	})
}