package book

type Book struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}

type BookModel struct {
	ID     int64  `db:"id" gorm:"primaryKey" `
	Title  string `db:"title"`
	Author string `db:"author"`
}

var BookTableName = "books"

func New(bookStorage BookStorage) BookHandler {
	return NewHandler(NewBookService(bookStorage))
}
