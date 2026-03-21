package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type zerologLogger struct {
	logger zerolog.Logger
}

func New(level string) Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	var logLevel zerolog.Level
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &zerologLogger{logger: logger}
}

func (l *zerologLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug().Msgf(msg, args...)
}

func (l *zerologLogger) Info(msg string, args ...interface{}) {
	l.logger.Info().Msgf(msg, args...)
}

func (l *zerologLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn().Msgf(msg, args...)
}

func (l *zerologLogger) Error(msg string, args ...interface{}) {
	l.logger.Error().Msgf(msg, args...)
}

func (l *zerologLogger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatal().Msgf(msg, args...)
}
