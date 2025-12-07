/*
 * GoAstra Backend - Structured Logger
 *
 * Production-grade logging using zap for high performance.
 * Supports structured logging with configurable levels.
 */
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
 * Logger wraps zap.SugaredLogger for convenient structured logging.
 */
type Logger struct {
	*zap.SugaredLogger
	level zapcore.Level
}

/*
 * New creates a new logger instance with the specified level.
 * Configures JSON output for production, console for development.
 */
func New(level string) *Logger {
	zapLevel := parseLevel(level)

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	/* Use console encoder for development */
	if os.Getenv("APP_ENV") == "development" {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zapLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		/* Fallback to default logger */
		logger = zap.NewNop()
	}

	return &Logger{
		SugaredLogger: logger.Sugar(),
		level:         zapLevel,
	}
}

/*
 * WithFields returns a new logger with additional context fields.
 */
func (l *Logger) WithFields(fields ...interface{}) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(fields...),
		level:         l.level,
	}
}

/*
 * WithRequestID returns a logger with request ID context.
 */
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithFields("request_id", requestID)
}

/*
 * WithUserID returns a logger with user ID context.
 */
func (l *Logger) WithUserID(userID uint) *Logger {
	return l.WithFields("user_id", userID)
}

/*
 * Sync flushes any buffered log entries.
 * Should be called before application exit.
 */
func (l *Logger) Sync() {
	_ = l.SugaredLogger.Sync()
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
