package user

import (
	"go-restapi/app"
	"go-restapi/utils"
)

type UserService interface {
	CreateUser(app.Context, CreateUserRequest) error
}

type userService struct {
	userStorage UserStorage
	utils       utils.Utils
}

func NewUserService(userStorage UserStorage, utils utils.Utils) UserService {
	return &userService{
		userStorage: userStorage,
		utils:       utils,
	}
}

func (s *userService) CreateUser(ctx app.Context, req CreateUserRequest) error {
	hashPassword, err := s.utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := UserModel{
		Username:  req.Username,
		Password:  hashPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    1,
	}
	return s.userStorage.CreateUser(user)
}
