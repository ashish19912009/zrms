package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

/**
Example usage
logger.Info("Starting authentication service", map[string]interface{}{
	"env": "local",
})
logger.Debug("Debugging mode active", nil)
logger.Warn("High CPU usage detected", nil)
logger.Error("Database connection failed", nil, nil)
*/

var Logger zerolog.Logger

// InitLogger initializes the global logger with console & file output
func InitLogger(serviceName string, level zerolog.Level, logFile string) {
	// Console output
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	// File output
	fileWriter, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	multiWriter := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

	Logger = zerolog.New(multiWriter).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger().
		Level(level)
}

// Info logs an informational message
func Info(msg string, fields map[string]interface{}) {
	logEvent := Logger.Info()
	for key, value := range fields {
		logEvent = logEvent.Interface(key, value)
	}
	logEvent.Msg(msg)
}

// Debug logs a debug-level message
func Debug(msg string, fields map[string]interface{}) {
	logEvent := Logger.Debug()
	for key, value := range fields {
		logEvent = logEvent.Interface(key, value)
	}
	logEvent.Msg(msg)
}

// Warn logs a warning message
func Warn(msg string, fields map[string]interface{}) {
	logEvent := Logger.Warn()
	for key, value := range fields {
		logEvent = logEvent.Interface(key, value)
	}
	logEvent.Msg(msg)
}

// Error logs an error message
func Error(msg string, err error, fields map[string]interface{}) {
	logEvent := Logger.Error().Err(err)
	for key, value := range fields {
		logEvent = logEvent.Interface(key, value)
	}
	logEvent.Msg(msg)
}

// Fatal logs a fatal error and exits
func Fatal(msg string, err error, fields map[string]interface{}) {
	logEvent := Logger.Fatal().Err(err)
	for key, value := range fields {
		logEvent = logEvent.Interface(key, value)
	}
	logEvent.Msg(msg)
	os.Exit(1) // Ensure exit on fatal error
}

// Panic logs a panic message and panics
func Panic(msg string, err error, fields map[string]interface{}) {
	logEvent := Logger.Panic().Err(err)
	for key, value := range fields {
		logEvent = logEvent.Interface(key, value)
	}
	logEvent.Msg(msg)
	panic(err)
}
