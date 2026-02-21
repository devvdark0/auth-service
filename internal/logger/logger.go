package logger

import (
	"log/slog"
	"os"
)

const (
	devEnv  = "dev"
	prodEnv = "prod"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(env string) *Logger {
	var log *slog.Logger

	switch env {
	case devEnv:
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case prodEnv:
		log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return &Logger{log}
}
