package log

import (
	"io"
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(err error, msg string, kv ...any)
	Fatal(msg string, kv ...any)

	With(kv ...any) Logger

	SkipFrames(skip int) Logger

	Slog() *slog.Logger
}

type Options struct {
	// Writer defaults to os.Stderr
	Writer io.Writer

	// ShowTimestamp defaults to false
	ShowTimestamp bool

	// ShowCaller defaults to true
	ShowCaller bool

	// ShowDebugLogs defaults to false
	ShowDebugLogs bool

	// ShowLogLevel defaults to true
	ShowLogLevel bool

	// JSONFormat does JSON formatted logs
	JSONFormat bool
}

var defaultOptions = Options{
	Writer:        os.Stderr,
	ShowTimestamp: false,
	ShowCaller:    true,
	ShowDebugLogs: false,
	ShowLogLevel:  true,
	JSONFormat:    false,
}

var defaultLogger Logger

func SetDefaultLogger(l Logger) {
	defaultLogger = l
}

func DefaultLogger() Logger {
	if defaultLogger == nil {
		defaultLogger = New()
	}

	return defaultLogger
}

func New(options ...Options) Logger {
	return newZeroLogger(options...)
}
