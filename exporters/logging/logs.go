package logging

import (
	"os"

	"golang.org/x/exp/slog"
)

var (
	Logger *slog.Logger
)

// var Logs *zap.Logger
// var Buffer *zap.Logger

// func InitLogs() {
// 	logConfig := zap.NewProductionEncoderConfig()

// 	logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
// 	logConfig.LineEnding = ","
// 	fileEncoder := zapcore.NewJSONEncoder(logConfig)
// 	consoleEncoder := zapcore.NewConsoleEncoder(logConfig)

// 	bufferConfig := zapcore.EncoderConfig{
// 		MessageKey: "msg",
// 		LineEnding: ",",
// 	}

// 	bufferEncoder := zapcore.NewJSONEncoder(bufferConfig)

// 	_, err0 := os.Stat("./tmp")

// 	if err0 != nil {
// 		os.MkdirAll("tmp/logs", 0777)
// 	}

// 	logFile, err := os.OpenFile("tmp/logs/log.json", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	bufferFile, err := os.OpenFile("tmp/logs/buffer.json", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

// 	defaultLogLevel := zapcore.ErrorLevel
// 	core := zapcore.NewTee(
// 		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), defaultLogLevel),
// 		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
// 	)

// 	bufferCore := zapcore.NewTee(
// 		zapcore.NewCore(bufferEncoder, zapcore.AddSync(bufferFile), defaultLogLevel))

// 	Logs = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
// 	Buffer = zap.New(bufferCore, zap.AddStacktrace(zapcore.ErrorLevel))

// }

func InitLogger() {
	logOpts := slog.HandlerOptions{Level: slog.LevelInfo}
	textHandler := logOpts.NewTextHandler(os.Stdout)
	Logger = slog.New(textHandler)
}

func LogGeneric(lvl string, msg string, pkg string) {

	switch lvl {
	case "info":
		Logger.Info(msg, slog.String("package", pkg))
	case "warning":
		Logger.Warn(msg, slog.String("package", pkg))
	case "debub":
		Logger.Debug(msg, slog.String("package", pkg))
	}
}

func LogError(err error, msg string, pkg string) {
	Logger.Error(msg, err, slog.String("package", pkg))
}
