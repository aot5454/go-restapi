package book

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestStorage(t *testing.T) {
	sqlmockDB, mock, _ := sqlmock.New()

	mock.ExpectQuery(`SELECT VERSION()`).WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("7.2"))

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlmockDB}), &gorm.Config{})
	if err != nil {
		t.Errorf("Error should be nil, got: %v", err)
	}

	t.Run("GetAllBook: Should return array", func(t *testing.T) {
		data := sqlmock.NewRows([]string{"id", "title", "author"}).AddRow(1, "test1", "test2").AddRow(2, "test1", "test2")
		mock.ExpectQuery("SELECT").WillReturnRows(data)

		storage := NewBookStorage(gormDB)
		got, err := storage.GetAllBook()
		if err != nil {
			t.Errorf("Error should be nil, got: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("Length of books should be 2, got: %d", len(got))
		}
	})

	t.Run("GetAllBook: Should return error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("Should return error"))

		storage := NewBookStorage(gormDB)
		_, err := storage.GetAllBook()
		if err == nil {
			t.Errorf("Error should be not nil")
		}
	})

	t.Run("CreateBook: Should return nil", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		storage := NewBookStorage(gormDB)
		err := storage.CreateBook(BookRequest{})
		if err != nil {
			t.Errorf("Error should be nil, got: %v", err)
		}
	})

	t.Run("CreateBook: Should return error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT").WillReturnError(errors.New("Should return error"))
		mock.ExpectRollback()

		storage := NewBookStorage(gormDB)
		err := storage.CreateBook(BookRequest{})
		if err == nil {
			t.Errorf("Error should be not nil")
		}
	})
}
