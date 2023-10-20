package user

import "gorm.io/gorm"

type UserStorage interface {
	CreateUser(UserModel) error
	GetListUser(int, int) ([]UserModel, error)
	GetUserByUsername(string) (*UserModel, error)
	CountListUser() (int64, error)
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

func (s *userStorage) GetListUser(limit, offset int) ([]UserModel, error) {
	var users []UserModel
	q := s.db.Debug().Table(UserTableName).Limit(limit).Offset(offset).Find(&users)
	if q.Error != nil {
		return nil, q.Error
	}
	return users, nil
}

func (s *userStorage) CountListUser() (int64, error) {
	var count int64
	q := s.db.Debug().Table(UserTableName).Count(&count)
	if q.Error != nil {
		return 0, q.Error
	}
	return count, nil
}

func (s *userStorage) GetUserByUsername(username string) (*UserModel, error) {
	var user UserModel
	q := s.db.Debug().Table(UserTableName).Where("username = ?", username).First(&user)
	if q.Error != nil {
		return nil, q.Error
	}
	return &user, nil
}
