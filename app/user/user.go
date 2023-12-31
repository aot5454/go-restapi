package user

import (
	"errors"
	"go-restapi/utils"
)

type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8,max=50"`
	FirstName string `json:"firstname" validate:"required,min=3,max=50"`
	LastName  string `json:"lastname" validate:"required,min=3,max=50"`
}

type GetListUserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Status    string `json:"status"`
}

type GetUserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Status    string `json:"status"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstname" validate:"required,min=3,max=50"`
	LastName  string `json:"lastname" validate:"required,min=3,max=50"`
	Status    string `json:"status" validate:"required,oneof=active inactive"`
}

const UserTableName = "users"

type UserModel struct {
	ID        int64  `db:"id" gorm:"primaryKey" `
	Username  string `db:"username" gorm:"unique"`
	Password  string `db:"password"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Status    int    `db:"status" gorm:"default:1"`
}

var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrUserNotFound = errors.New("user not found")

func New(userStorage UserStorage) UserHandler {
	service := NewUserService(userStorage, utils.NewUtils())
	handler := NewUserHandler(service, utils.NewUtils())
	return handler
}
