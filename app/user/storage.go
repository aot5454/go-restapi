package user

import "gorm.io/gorm"

type UserStorage interface {
	CreateUser(UserModel) error
	GetListUser() ([]UserModel, error)
	GetUserByUsername(string) (*UserModel, error)
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

func (s *userStorage) GetListUser() ([]UserModel, error) {
	var users []UserModel
	q := s.db.Debug().Table(UserTableName).Find(&users)
	if q.Error != nil {
		return nil, q.Error
	}
	return users, nil
}

func (s *userStorage) GetUserByUsername(username string) (*UserModel, error) {
	var user UserModel
	q := s.db.Debug().Table(UserTableName).Where("username = ?", username).First(&user)
	if q.Error != nil {
		return nil, q.Error
	}
	return &user, nil
}
