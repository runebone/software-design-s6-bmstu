package logger

import (
	"context"
	"os"
	"path/filepath"
	"user/internal/common/logger"
	"user/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(config config.LogConfig) *ZapLogger {
	var zapLogLevel zapcore.Level
	if err := zapLogLevel.UnmarshalText([]byte(config.Level)); err != nil {
		zapLogLevel = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	consoleWriteSyncer := zapcore.AddSync(os.Stdout)

	dir := filepath.Dir(config.Path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}

	var fileWriteSyncer zapcore.WriteSyncer
	if config.Path != "" {
		logFile, err := os.OpenFile(config.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		fileWriteSyncer = zapcore.AddSync(logFile)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriteSyncer, zapLogLevel),
		zapcore.NewCore(fileEncoder, fileWriteSyncer, zapLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &ZapLogger{logger: logger}
}

func (l *ZapLogger) zapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

func (l *ZapLogger) WithFields(fields map[string]interface{}) logger.Logger {
	return &ZapLogger{logger: l.logger.With(l.zapFields(fields)...)}
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Debug(msg, zap.Any("context", fields))
}

func (l *ZapLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Info(msg, zap.Any("context", fields))
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Warn(msg, zap.Any("context", fields))
}

func (l *ZapLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Error(msg, zap.Any("context", fields))
}

func (l *ZapLogger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.logger.Fatal(msg, zap.Any("context", fields))
}
