package logger

import (
	"log"
)

// Logger is a logger object.
type Logger struct{}

// New is a logger constructor.
func New() *Logger {
	return &Logger{}
}

// Info prints information log.
func (l *Logger) Info(format string, args ...interface{}) {
	log.Printf("[INFO] "+format+"\n", args...)
}

// Error prints error log.
func (l *Logger) Error(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format+"\n", args...)
}

// Panic prints panic log and thrown a panic.
func (l *Logger) Panic(format string, args ...interface{}) {
	log.Panicf("[PANIC] "+format+"\n", args...)
}
