/*
 * GoAstra CLI - Logger Template
 *
 * Generates structured logging using Uber's Zap library.
 * Provides environment-aware configuration (dev vs prod).
 */
package common

// LoggerGo returns the logger.go template with Zap configuration.
func LoggerGo() string {
	return `package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
 * Logger wraps zap.SugaredLogger with convenience methods.
 * Provides structured logging with context support.
 */
type Logger struct {
	*zap.SugaredLogger
}

/*
 * New creates a new Logger instance configured for the given environment.
 * Development mode uses colored console output.
 * Production mode uses JSON format for log aggregation.
 */
func New(env string) *Logger {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	}

	// Set log level from environment
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		var zapLevel zapcore.Level
		if err := zapLevel.UnmarshalText([]byte(level)); err == nil {
			config.Level = zap.NewAtomicLevelAt(zapLevel)
		}
	}

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		// Fallback to no-op logger
		return &Logger{zap.NewNop().Sugar()}
	}

	return &Logger{logger.Sugar()}
}

/*
 * WithFields returns a new logger with the given fields attached.
 */
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{l.SugaredLogger.With(args...)}
}

/*
 * WithRequestID returns a logger with request ID context.
 */
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{l.SugaredLogger.With("request_id", requestID)}
}

/*
 * WithUserID returns a logger with user ID context.
 */
func (l *Logger) WithUserID(userID uint) *Logger {
	return &Logger{l.SugaredLogger.With("user_id", userID)}
}

/*
 * Sync flushes any buffered log entries.
 * Should be called before application exit.
 */
func (l *Logger) Sync() {
	_ = l.SugaredLogger.Sync()
}
`
}
