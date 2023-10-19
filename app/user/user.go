package user

import "go-restapi/utils"

type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8,max=50"`
	FirstName string `json:"firstname" validate:"required,min=3,max=50"`
	LastName  string `json:"lastname" validate:"required,min=3,max=50"`
}

const UserTableName = "users"

type UserModel struct {
	ID        int64  `db:"id" gorm:"primaryKey" `
	Username  string `db:"username" gorm:"unique"`
	Password  string `db:"password"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	Status    int    `db:"status" gorm:"default:1"`
}

func New(userStorage UserStorage) UserHandler {
	service := NewUserService(userStorage, utils.NewUtils())
	handler := NewUserHandler(service)
	return handler
}
