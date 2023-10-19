package user

import "gorm.io/gorm"

type UserStorage interface {
	CreateUser(UserModel) error
}

type userStorage struct {
	db *gorm.DB
}

func NewUserStorage(db *gorm.DB) UserStorage {
	return &userStorage{db: db}
}

func (s *userStorage) CreateUser(user UserModel) error {
	q := s.db.Table(UserTableName).Create(&user)
	if q.Error != nil {
		return q.Error
	}
	return nil
}
