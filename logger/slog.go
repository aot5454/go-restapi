package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var logLevel = &slog.LevelVar{}
var logger *slog.Logger
var pid = os.Getpid()
var goversion string = "unknown"

func init() {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		goversion = buildInfo.GoVersion
	}
}

func New() *slog.Logger {
	logLevel.Set(slog.LevelDebug)

	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
			source.Function = filepath.Base(source.Function)
		}
		return a
	}

	opts := &slog.HandlerOptions{
		Level:       logLevel,
		AddSource:   true,
		ReplaceAttr: replace,
	}

	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, opts)
	if strings.ToLower(os.Getenv("ENV")) == "local" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger = slog.New(handler)

	slog.SetDefault(logger)

	return logger
}

func Debugf(format string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelDebug) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(format, args...), pcs[0])
	r.AddAttrs([]slog.Attr{
		slog.Int("pid", pid),
		slog.String("go_version", goversion),
	}...)
	_ = logger.Handler().Handle(context.Background(), r)
}

func Warnf(format string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelWarn) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

func Infof(format string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

func Errorf(format string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

func AppErrorf(handler slog.Handler, format string, args ...any) {
	if !handler.Enabled(context.Background(), slog.LevelError) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), slog.LevelError, fmt.Sprintf(format, args...), pcs[0])
	r.AddAttrs([]slog.Attr{
		slog.Int("pid", pid),
		slog.String("go_version", goversion),
	}...)

	_ = handler.Handle(context.Background(), r)
}
