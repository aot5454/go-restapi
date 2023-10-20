package router

import (
	"go-restapi/app"
	"go-restapi/app/book"
	"go-restapi/app/user"

	"gorm.io/gorm"
)

func Router(r *app.Router, db *gorm.DB) *app.Router {
	bookStorege := book.NewBookStorage(db)
	userStorage := user.NewUserStorage(db)

	bookHandler := book.New(bookStorege)
	userHandler := user.New(userStorage)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/books", bookHandler.GetAllBook)
		v1.POST("/books", bookHandler.CreateBook)

		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users", userHandler.GetListUser)
	}

	r.NoRoute()
	return r
}
