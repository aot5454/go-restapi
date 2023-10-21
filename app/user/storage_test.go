package user

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mockUserStorageData = UserModel{
	ID:        1,
	Username:  "test",
	Password:  "password",
	FirstName: "test",
	LastName:  "test",
	Status:    1,
}

var mockUserStorageDataList = []UserModel{
	{
		ID:        1,
		Username:  "test",
		Password:  "password",
		FirstName: "test",
		LastName:  "test",
		Status:    1,
	},
}

type testStorageSuite struct {
	suite.Suite
	sqlmockDB *sql.DB
	mock      sqlmock.Sqlmock
	gormDB    *gorm.DB
	data      UserModel
}

func (s *testStorageSuite) SetupTest() {
	sqlmockDB, mock, _ := sqlmock.New()

	mock.ExpectQuery(`SELECT VERSION()`).WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("7.2"))

	gormDB, _ := gorm.Open(mysql.New(mysql.Config{Conn: sqlmockDB}), &gorm.Config{})

	s.sqlmockDB = sqlmockDB
	s.mock = mock
	s.gormDB = gormDB
	s.data = mockUserStorageData
}

func (s *testStorageSuite) TearDownTest() {
	s.sqlmockDB.Close()
}

func (s *testStorageSuite) TestCreateUserStorage() {
	s.Run("Should return nil", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		storage := NewUserStorage(s.gormDB)
		err := storage.CreateUser(s.data)
		s.NoError(err)
	})

	s.Run("Should return error", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("INSERT").WillReturnError(sql.ErrConnDone)
		s.mock.ExpectCommit()

		storage := NewUserStorage(s.gormDB)
		err := storage.CreateUser(s.data)
		s.Error(err)
	})
}

func (s *testStorageSuite) TestGetListUserStorage() {
	s.Run("Should return nil", func() {
		s.mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "username", "password", "first_name", "last_name", "status"}).
				AddRow(1, "test", "password", "test", "test", 1))

		storage := NewUserStorage(s.gormDB)
		got, err := storage.GetListUser(1, 10)
		s.NoError(err)
		s.EqualValues(mockUserStorageDataList, got)
	})

	s.Run("Should return error", func() {
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)

		storage := NewUserStorage(s.gormDB)
		_, err := storage.GetListUser(1, 10)
		s.Error(err)
	})
}

func (s *testStorageSuite) TestCountListUserStorage() {
	s.Run("Should return nil", func() {
		s.mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.
				NewRows([]string{"count"}).
				AddRow(1))

		storage := NewUserStorage(s.gormDB)
		got, err := storage.CountListUser()
		s.NoError(err)
		s.EqualValues(1, got)
	})

	s.Run("Should return error", func() {
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)

		storage := NewUserStorage(s.gormDB)
		_, err := storage.CountListUser()
		s.Error(err)
	})
}

func (s *testStorageSuite) TestGetUserByUsername() {
	s.Run("Should return nil", func() {
		s.mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "username", "password", "first_name", "last_name", "status"}).
				AddRow(1, "test", "password", "test", "test", 1))

		storage := NewUserStorage(s.gormDB)
		got, err := storage.GetUserByUsername("test")
		s.NoError(err)
		s.EqualValues(mockUserStorageData, *got)
	})

	s.Run("Should return error", func() {
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)

		storage := NewUserStorage(s.gormDB)
		_, err := storage.GetUserByUsername("test")
		s.Error(err)
	})
}

func (s *testStorageSuite) TestGetUserByID() {
	s.Run("Should return nil", func() {
		s.mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.
				NewRows([]string{"id", "username", "password", "first_name", "last_name", "status"}).
				AddRow(1, "test", "password", "test", "test", 1))

		storage := NewUserStorage(s.gormDB)
		got, err := storage.GetUserByID(1)
		s.NoError(err)
		s.EqualValues(mockUserStorageData, *got)
	})

	s.Run("Should return error", func() {
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)

		storage := NewUserStorage(s.gormDB)
		_, err := storage.GetUserByID(1)
		s.Error(err)
	})
}

func (s *testStorageSuite) TestUpdateUser() {
	s.Run("Should return nil", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		storage := NewUserStorage(s.gormDB)
		err := storage.UpdateUser(s.data)
		s.NoError(err)
	})

	s.Run("Should return error", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec("UPDATE").WillReturnError(sql.ErrConnDone)
		s.mock.ExpectCommit()

		storage := NewUserStorage(s.gormDB)
		err := storage.UpdateUser(s.data)
		s.Error(err)
	})
}

func TestUserStorage(t *testing.T) {
	suite.Run(t, new(testStorageSuite))
}
