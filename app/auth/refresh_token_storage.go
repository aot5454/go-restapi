package auth

import "gorm.io/gorm"

type RefreshTokenStorage interface {
	GetRefreshTokenByToken(token string) (*RefreshTokenModel, error)
	UpsertRefreshToken(refreshToken RefreshTokenModel) error
}

type refreshTokenStorage struct {
	db *gorm.DB
}

func NewRefreshTokenStorage(db *gorm.DB) RefreshTokenStorage {
	return &refreshTokenStorage{
		db: db,
	}
}

func (s *refreshTokenStorage) GetRefreshTokenByToken(token string) (*RefreshTokenModel, error) {
	var refreshToken RefreshTokenModel
	if err := s.db.Debug().Table(RefreshTokenTableName).Where("token = ?", token).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (s *refreshTokenStorage) UpsertRefreshToken(refreshToken RefreshTokenModel) error {
	q := s.db.Debug().Table(RefreshTokenTableName).Where("user_id = ?", refreshToken.UserID).Updates(refreshToken)
	if q.Error != nil {
		return q.Error
	}
	if q.RowsAffected == 0 {
		q := s.db.Debug().Table(RefreshTokenTableName).Create(&refreshToken)
		if q.Error != nil {
			return q.Error
		}
	}
	return nil
}
