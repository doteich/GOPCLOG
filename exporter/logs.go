package exporter

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logs *zap.Logger

func InitLogs() {
	logConfig := zap.NewProductionEncoderConfig()
	logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(logConfig)

	_, err0 := os.Stat("./tmp")

	if err0 != nil {
		os.MkdirAll("tmp/logs", 0777)
	}

	logFile, err := os.OpenFile("tmp/logs/log.json", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		fmt.Println(err)
	}

	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	Logs = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

}

func WriteLogFile() {

}
