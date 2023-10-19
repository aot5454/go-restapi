package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockCreateUserRequest = CreateUserRequest{
	Username:  "test",
	Password:  "password",
	FirstName: "test",
	LastName:  "test",
}

func TestCreateUserService(t *testing.T) {

	t.Run("Should return success", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("CreateUser", mock.Anything).Return(nil)
		utils := &mockUtils{}
		utils.On("HashPassword", mock.Anything).Return("password", nil)

		service := NewUserService(storage, utils)

		err := service.CreateUser(nil, mockCreateUserRequest)
		assert.NoError(t, err)
	})

	t.Run("Should return error (HashPassword)", func(t *testing.T) {
		storage := &mockUserStorage{}
		storage.On("CreateUser", mock.Anything).Return(nil)
		utils := &mockUtils{}
		utils.On("HashPassword", mock.Anything).Return("", errors.New("error"))

		service := NewUserService(storage, utils)

		err := service.CreateUser(nil, mockCreateUserRequest)
		assert.Error(t, err)
	})

}
