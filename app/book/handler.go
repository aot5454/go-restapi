package book

import (
	"go-restapi/app"
)

type bookHandler struct {
	bookSvc BookService
}

type BookHandler interface {
	GetAllBook(ctx app.Context)
	CreateBook(ctx app.Context)
}

func NewHandler(bookSvc BookService) BookHandler {
	return &bookHandler{
		bookSvc: bookSvc,
	}
}

func (h *bookHandler) GetAllBook(ctx app.Context) {
	books, err := h.bookSvc.GetAllBook()
	if err != nil {
		ctx.StoreError(err)
		return
	}
	ctx.OK(books)
}

func (h *bookHandler) CreateBook(ctx app.Context) {
	var book BookRequest
	if err := ctx.Bind(&book); err != nil {
		ctx.BadRequest(err)
		return
	}

	if _, err := ctx.Validate(&book); err != nil {
		ctx.BadRequest(err)
		return
	}

	if err := h.bookSvc.CreateBook(book); err != nil {
		ctx.StoreError(err)
		return
	}

	ctx.OK(nil)
}
