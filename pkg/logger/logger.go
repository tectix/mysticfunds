package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}

type zapLogger struct {
	logger *zap.SugaredLogger
}

// NewLogger creates a new Logger
func NewLogger(logLevel string) Logger {
	config := zap.NewProductionConfig()

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}
	config.Level.SetLevel(level)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &zapLogger{
		logger: logger.Sugar(),
	}
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) With(keysAndValues ...interface{}) Logger {
	return &zapLogger{
		logger: l.logger.With(keysAndValues...),
	}
}
