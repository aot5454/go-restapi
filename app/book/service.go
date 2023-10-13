package book

type bookService struct {
	bookStorage BookStorage
}

type BookService interface {
	GetAllBook() ([]Book, error)
	CreateBook(BookRequest) error
}

func NewBookService(bookStorage BookStorage) BookService {
	return &bookService{bookStorage: bookStorage}
}

func (s *bookService) GetAllBook() ([]Book, error) {
	books, err := s.bookStorage.GetAllBook()
	if err != nil {
		return nil, err
	}
	result := s.mappingBookModelToBook(books)
	return result, nil
}

func (s *bookService) CreateBook(book BookRequest) error {
	return s.bookStorage.CreateBook(book)
}

func (s *bookService) mappingBookModelToBook(books []BookModel) []Book {
	var result []Book
	for _, book := range books {
		result = append(result, Book(book))
	}
	return result
}
