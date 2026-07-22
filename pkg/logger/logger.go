package logger

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func New() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "INFO  ", log.LstdFlags),
		warn:  log.New(os.Stdout, "WARN  ", log.LstdFlags),
		error: log.New(os.Stderr, "ERROR ", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.info.Printf(msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.warn.Printf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.error.Printf(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.error.Fatalf(msg, args...)
}

