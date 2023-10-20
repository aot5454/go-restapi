package user

import (
	"errors"
	"go-restapi/app"
	"go-restapi/utils"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(app.Context, CreateUserRequest) error
	GetListUser(app.Context, int, int) ([]GetListUserResponse, error)
	GetUserByID(app.Context, int) (*GetUserResponse, error)
	CountListUser(app.Context) (int, error)
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

	checkDup, err := s.userStorage.GetUserByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if checkDup != nil {
		return ErrUsernameAlreadyExists
	}

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

func (s *userService) GetListUser(ctx app.Context, page int, pageSize int) ([]GetListUserResponse, error) {
	limit := pageSize
	offset := (page - 1) * pageSize
	users, err := s.userStorage.GetListUser(limit, offset)
	if err != nil {
		return nil, err
	}

	var res []GetListUserResponse
	for _, user := range users {
		status := "Active"
		if user.Status == 0 {
			status = "Inactive"
		}
		res = append(res, GetListUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Status:    status,
		})
	}
	return res, nil
}

func (s *userService) GetUserByID(ctx app.Context, id int) (*GetUserResponse, error) {
	user, err := s.userStorage.GetUserByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	status := "Active"
	if user.Status == 0 {
		status = "Inactive"
	}
	res := GetUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Status:    status,
	}
	return &res, nil
}

func (s *userService) CountListUser(ctx app.Context) (int, error) {
	count, err := s.userStorage.CountListUser()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
