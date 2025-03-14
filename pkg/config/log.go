package config

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

// CustomFormatter is a logrus formatter that includes file name and line number
type CustomFormatter struct {
	// You can add any fields here if you need to customize it further
}

// Format will format the log entries with the file name and line number
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	_, file, line, ok := runtime.Caller(8) // Adjust the 8 as needed to get the caller's info from where the log is generated
	if !ok {
		file = "unknown"
		line = 0
	}

	// Set up mild log color based on log level
	var levelColor string
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = "\033[34m" // Blue for Debug
	case logrus.InfoLevel:
		levelColor = "\033[32m" // Green for Info
	case logrus.WarnLevel:
		levelColor = "\033[33m" // Yellow for Warning
	case logrus.ErrorLevel:
		levelColor = "\033[31m" // Red for Error
	case logrus.FatalLevel:
		levelColor = "\033[35m" // Magenta for Fatal
	case logrus.PanicLevel:
		levelColor = "\033[36m" // Cyan for Panic
	default:
		levelColor = "\033[37m" // Default gray color
	}

	// Mild color for file and line
	fileLineColor := "\033[90m" // Light Gray

	// Reset color at the end of the message
	resetColor := "\033[0m"

	// Custom log format: timestamp, level with color, file, line, and message
	return []byte(fmt.Sprintf("%s%s [%s%s:%d%s] %s%s: %s%s\n",
		levelColor,                               // Color for log level
		entry.Time.Format("2006-01-02 15:04:05"), // Timestamp
		fileLineColor,                            // Color for file and line
		file, line,                               // File and Line (correct formatting)
		resetColor,           // Reset color after file/line
		entry.Level.String(), // Log Level as a string
		resetColor,           // Reset color after log level
		entry.Message,        // Log Message
		resetColor,           // Reset color after message
	)), nil
}

var log = logrus.New()

// GetLogger returns the configured logger
func GetLogger() *logrus.Logger {
	log.SetFormatter(&CustomFormatter{}) // Set custom formatter
	log.SetLevel(logrus.DebugLevel)      // Set log level to Debug or as per your need
	return log
}
