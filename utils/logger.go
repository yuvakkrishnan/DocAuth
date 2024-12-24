package utils

import (
	"log"
	"os"
)

// Logger structure
type Logger struct {
	logger *log.Logger
}

// NewLogger creates and returns a new logger instance
func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an informational message
func (l *Logger) Info(message string) {
	l.logger.Println("INFO: " + message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.logger.Println("ERROR: " + message)
}
