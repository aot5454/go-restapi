package auth

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mockRefreshTokenStorageData = RefreshTokenModel{
	ID:        1,
	UserID:    1,
	Token:     "fcd277b6-562c-49f6-8146-051bb339fb8c",
	ExpiredAt: time.Date(2021, 8, 25, 15, 13, 7, 0, time.UTC),
	CreatedAt: time.Date(2021, 8, 24, 15, 13, 7, 0, time.UTC),
	CreatedBy: "admin",
	UpdatedAt: time.Date(2021, 8, 24, 15, 13, 7, 0, time.UTC),
	UpdatedBy: "admin",
}

type testRefreshTokenStorageSuite struct {
	suite.Suite
	sqlmockDB *sql.DB
	mock      sqlmock.Sqlmock
	gormDB    *gorm.DB
	data      RefreshTokenModel
}

func (s *testRefreshTokenStorageSuite) SetupTest() {
	sqlmockDB, mock, _ := sqlmock.New()

	mock.ExpectQuery(`SELECT VERSION()`).WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("7.2"))

	gormDB, _ := gorm.Open(mysql.New(mysql.Config{Conn: sqlmockDB}), &gorm.Config{})

	s.sqlmockDB = sqlmockDB
	s.mock = mock
	s.gormDB = gormDB
	s.data = mockRefreshTokenStorageData
}

func (s *testRefreshTokenStorageSuite) TearDownTest() {
	s.sqlmockDB.Close()
}

func (s *testRefreshTokenStorageSuite) TestGetRefreshTokenByToken() {
	s.Run("Should return nil", func() {
		s.mock.ExpectQuery(`SELECT`).WithArgs(s.data.Token).
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "user_id", "token", "expired_at", "created_at", "created_by", "updated_at", "updated_by"}).
				AddRow(s.data.ID, s.data.UserID, s.data.Token, s.data.ExpiredAt, s.data.CreatedAt, s.data.CreatedBy, s.data.UpdatedAt, s.data.UpdatedBy))
		storage := NewRefreshTokenStorage(s.gormDB)
		got, err := storage.GetRefreshTokenByToken(s.data.Token)
		s.NoError(err)
		s.Equal(&s.data, got)
	})

	s.Run("Should return error", func() {
		s.mock.ExpectQuery(`SELECT`).WithArgs(s.data.Token).WillReturnError(sql.ErrNoRows)
		storage := NewRefreshTokenStorage(s.gormDB)
		got, err := storage.GetRefreshTokenByToken(s.data.Token)
		s.Error(err)
		s.Nil(got)
	})
}

func (s *testRefreshTokenStorageSuite) TestUpsertRefreshToken() {
	defer func() {
		_ = s.mock.ExpectationsWereMet()
	}()

	s.Run("Should return nil", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		storage := NewRefreshTokenStorage(s.gormDB)
		err := storage.UpsertRefreshToken(s.data)
		s.NoError(err)
	})

	s.Run("Should return error when insert", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(`UPDATE`).WillReturnResult(sqlmock.NewResult(0, 0))
		s.mock.ExpectCommit()

		s.mock.ExpectBegin()
		s.mock.ExpectExec(`INSERT`).WillReturnError(sql.ErrConnDone)
		s.mock.ExpectCommit()

		storage := NewRefreshTokenStorage(s.gormDB)
		err := storage.UpsertRefreshToken(s.data)
		s.Error(err)
	})

	s.Run("Should return error when update", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(`UPDATE`).WillReturnError(sql.ErrConnDone)
		s.mock.ExpectCommit()

		storage := NewRefreshTokenStorage(s.gormDB)
		err := storage.UpsertRefreshToken(s.data)
		s.Error(err)
	})
}

func TestRefreshTokenStorage(t *testing.T) {
	suite.Run(t, new(testRefreshTokenStorageSuite))
}
