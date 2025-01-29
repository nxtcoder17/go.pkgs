package main

import (
	"fmt"
	"log/slog"

	"github.com/nxtcoder17/go.pkgs/log"
)

func x(l *slog.Logger) {
	l.Debug("Hello")
	l.Info("Hello")
	l.Warn("Hello")
	l.Error("hello", "err", fmt.Errorf("this is an error"))
}

func main() {
	logger := log.New(log.Options{
		ShowTimestamp: false,
		ShowCaller:    true,
		ShowLogLevel:  true,
		ShowDebugLogs: false,
	})
	err := fmt.Errorf("this is an error")

	logger.Debug("hello", "msg", "greeting")
	logger.Info("hello", "msg", "greeting")
	logger.Warn("hello", "msg", "greeting")
	logger.Error(err, "hello", "msg", fmt.Errorf("this is an error"))

	l := logger.With("msg", "with/greeting")

	l.Debug("Hello")
	l.Info("Hello")
	l.Warn("Hello")
	l.Error(err, "msg", fmt.Errorf("this is an error"))

	// sl := l.Slog()
	x(l.Slog())
}
