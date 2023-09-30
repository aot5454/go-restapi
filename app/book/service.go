package book

type bookService struct {
	bookStorage BookStorage
}

type BookService interface {
	GetAllBook() ([]Book, error)
	CreateBook(BookRequest) error
}

func NewService(bookStorage BookStorage) BookService {
	return &bookService{bookStorage: bookStorage}
}

func (s *bookService) GetAllBook() ([]Book, error) {
	books, err := s.bookStorage.GetAllBook()
	if err != nil {
		return nil, err
	}

	var result []Book
	for _, book := range books {
		result = append(result, Book{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
		})
	}

	return result, nil
}

func (s *bookService) CreateBook(book BookRequest) error {
	return s.bookStorage.CreateBook(book)
}
