// Package logging provides a configured zerolog logger for the application.
package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is the global logger instance.
var Logger zerolog.Logger

// Init initializes the global logger with the specified log level.
// It configures zerolog to use human-readable console output with:
// - ISO-8601 timestamp format
// - Caller information (file:line)
func Init(level string) {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	// Set global log level
	zerolog.SetGlobalLevel(logLevel)

	// Configure time format to ISO-8601
	zerolog.TimeFieldFormat = time.RFC3339

	// Create console writer for human-readable output
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	// Create logger with caller info and timestamp
	Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()
}

// Debug logs a debug message.
func Debug() *zerolog.Event {
	return Logger.Debug()
}

// Info logs an info message.
func Info() *zerolog.Event {
	return Logger.Info()
}

// Warn logs a warning message.
func Warn() *zerolog.Event {
	return Logger.Warn()
}

// Error logs an error message.
func Error() *zerolog.Event {
	return Logger.Error()
}

// Fatal logs a fatal message and exits.
func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

// Panic logs a panic message and panics.
func Panic() *zerolog.Event {
	return Logger.Panic()
}
