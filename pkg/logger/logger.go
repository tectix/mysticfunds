package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(logLevel string) {
	config := zap.NewProductionConfig()

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}
	config.Level.SetLevel(level)

	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	log.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	log.Fatal(message, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}
