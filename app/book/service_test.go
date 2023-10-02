package book

import (
	"errors"
	"testing"
)

var bookModelMock = []BookModel{
	{
		ID:     1,
		Title:  "test",
		Author: "test",
	},
	{
		ID:     2,
		Title:  "test",
		Author: "test",
	},
}

type bookStorageMockSuccess struct {
	BookStorage
}

func (m *bookStorageMockSuccess) GetAllBook() ([]BookModel, error) {
	return bookModelMock, nil
}

func (m *bookStorageMockSuccess) CreateBook(book BookRequest) error {
	return nil
}

func TestBookServiceSuccessCase(t *testing.T) {
	t.Run("GetAllBook: Should return array", func(t *testing.T) {
		storage := &bookStorageMockSuccess{}
		svc := NewService(storage)

		books, err := svc.GetAllBook()
		if err != nil {
			t.Errorf("Error should be nil, got: %v", err)
		}

		if len(books) != 2 {
			t.Errorf("Length of books should be 2, got: %d", len(books))
		}
	})

	t.Run("CreateBook: Should return nil", func(t *testing.T) {
		storage := &bookStorageMockSuccess{}
		svc := NewService(storage)

		err := svc.CreateBook(BookRequest{})
		if err != nil {
			t.Errorf("Error should be nil, got: %v", err)
		}
	})
}

// --------------------

type bookStorageMockError struct {
	BookStorage
}

func (m *bookStorageMockError) GetAllBook() ([]BookModel, error) {
	return []BookModel{}, errors.New("error")
}

func (m *bookStorageMockError) CreateBook(book BookRequest) error {
	return errors.New("error")
}

func TestBookServiceErrorsCase(t *testing.T) {
	t.Run("GetAllBook: Should return error", func(t *testing.T) {
		storage := &bookStorageMockError{}
		svc := NewService(storage)

		_, err := svc.GetAllBook()
		if err == nil {
			t.Errorf("Error should be not nil, got: %v", err)
		}
	})

	t.Run("CreateBook: Should return error", func(t *testing.T) {
		storage := &bookStorageMockError{}
		svc := NewService(storage)

		err := svc.CreateBook(BookRequest{})
		if err == nil {
			t.Errorf("Error should be not nil, got: %v", err)
		}
	})
}
