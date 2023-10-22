package router

import (
	"go-restapi/app"
	"go-restapi/app/auth"
	"go-restapi/app/book"
	"go-restapi/app/user"

	"gorm.io/gorm"
)

func Router(r *app.Router, db *gorm.DB) *app.Router {
	bookStorege := book.NewBookStorage(db)
	userStorage := user.NewUserStorage(db)
	refreshTokenStorage := auth.NewRefreshTokenStorage(db)

	authHandler := auth.New(userStorage, refreshTokenStorage)
	bookHandler := book.New(bookStorege)
	userHandler := user.New(userStorage)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", authHandler.Login)

		v1.GET("/books", bookHandler.GetAllBook)
		v1.POST("/books", bookHandler.CreateBook)

		v1.POST("/users", userHandler.CreateUser)
		v1.GET("/users", userHandler.GetListUser)
		v1.GET("/users/:id", userHandler.GetUserByID)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
	}

	r.NoRoute()
	return r
}
