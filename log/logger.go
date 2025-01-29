package log

import (
	"io"
	"log/slog"
)

type Logger interface {
	Debug(msg string, kv ...any)
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(err error, msg string, kv ...any)
	Fatal(msg string, kv ...any)

	With(kv ...any) Logger

	Slog() *slog.Logger
}

type Options struct {
	Writer io.Writer

	ShowTimestamp bool
	ShowCaller    bool
	ShowDebugLogs bool
	ShowLogLevel  bool
}
