package logging

import (
	"log/slog"
	"os"
)

var (
	Logger *slog.Logger
)

func InitLogger() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func LogGeneric(lvl string, msg string, pkg string) {

	switch lvl {
	case "info":
		Logger.Info(msg, slog.String("package", pkg))
	case "warning":
		Logger.Warn(msg, slog.String("package", pkg))
	case "debug":
		Logger.Debug(msg, slog.String("package", pkg))
	}
}

func LogError(err error, msg string, pkg string) {
	Logger.Error(msg, err, slog.String("package", pkg))
}
