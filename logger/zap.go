package logger

import (
	"log"

	"go.uber.org/zap"
)

func NewZap() (*zap.Logger, func()) {
	config := zap.NewProductionConfig()
	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}
	undo := zap.ReplaceGlobals(logger)

	return logger, func() {
		undo()
		_ = logger.Sync()
	}
}
