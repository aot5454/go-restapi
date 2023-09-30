package book

type bookStorage struct {

}

type BookStorage interface {
	GetAllBook() ([]BookModel, error)
	CreateBook(BookRequest) error
}

func NewStorage() BookStorage {
	return &bookStorage{}
}

func (s *bookStorage) GetAllBook() ([]BookModel, error) {
	return []BookModel{}, nil
}

func (s *bookStorage) CreateBook(book BookRequest) error {
	return nil
}