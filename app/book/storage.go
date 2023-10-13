package book

import "gorm.io/gorm"

type bookStorage struct {
	db *gorm.DB
}

type BookStorage interface {
	GetAllBook() ([]BookModel, error)
	CreateBook(BookRequest) error
}

func NewBookStorage(db *gorm.DB) BookStorage {
	return &bookStorage{db: db}
}

func (s *bookStorage) GetAllBook() ([]BookModel, error) {
	books := []BookModel{}
	q := s.db.Table(BookTableName).Find(&books)
	if q.Error != nil {
		return nil, q.Error
	}
	return books, nil
}

func (s *bookStorage) CreateBook(book BookRequest) error {
	bookModel := BookModel{
		Title:  book.Title,
		Author: book.Author,
	}
	q := s.db.Table(BookTableName).Create(&bookModel)
	if q.Error != nil {
		return q.Error
	}
	return nil
}
