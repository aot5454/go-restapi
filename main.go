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
	"go-restapi/config"
	"go-restapi/database"
	"go-restapi/logger"
	"go-restapi/router"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.NewMysqlDB(conf)
	if err != nil {
		panic(err)
	}

	logger := logger.New()
	r := app.NewRouter(logger, conf)
	r = router.Router(r, db)

	srv := http.Server{
		Addr:              ":" + conf.Server.Port,
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
		sqlDB, _ := db.DB()
		if err := sqlDB.Close(); err != nil {
			logger.Info("Database close: " + err.Error())
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
