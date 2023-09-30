package book

func New(bookStorage BookStorage) BookHandler {
	return NewHandler(NewService(bookStorage))
}