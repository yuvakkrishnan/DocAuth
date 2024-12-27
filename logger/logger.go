package logger

import (
    "fmt"
    "log"
    "os"
)

type Logger struct {
    infoLogger  *log.Logger
    errorLogger *log.Logger
}

func NewLogger() *Logger {
    return &Logger{
        infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
        errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}

func (l *Logger) Println(message string) {
    l.infoLogger.Println(message)
}

func (l *Logger) Info(message string) {
    l.infoLogger.Println(message)
}

func (l *Logger) Infof(format string, args ...interface{}) {
    l.infoLogger.Println(fmt.Sprintf(format, args...))
}

func (l *Logger) Error(message string) {
    l.errorLogger.Println(message)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
    l.errorLogger.Println(fmt.Sprintf(format, args...))
}