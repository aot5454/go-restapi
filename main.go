package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-restapi/app"
	"go-restapi/app/book"
	"go-restapi/app/user"
	"go-restapi/database"
	"go-restapi/logger"
)

func main() {

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")
	db, err := database.NewMysqlDB(username, password, host, port, databaseName)
	if err != nil {
		panic(err)
	}

	logger := logger.New()
	r := app.NewRouter(logger)

	bookStorege := book.NewBookStorage(db)
	userStorage := user.NewUserStorage(db)

	bookHandler := book.New(bookStorege)
	userHandler := user.New(userStorage)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/books", bookHandler.GetAllBook)
		v1.POST("/books", bookHandler.CreateBook)

		v1.POST("/users", userHandler.CreateUser)
	}

	r.NoRoute()

	srv := http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		d := time.Duration(5 * time.Second)
		fmt.Printf("shutting down int %s ...", d)
		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Info("HTTP server Shutdown: " + err.Error())
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("HTTP server ListenAndServe: " + err.Error())
		return
	}

	<-idleConnsClosed
	fmt.Println("gracefully")
}
