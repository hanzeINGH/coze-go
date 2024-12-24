package log

import (
	"fmt"
	"log"
	"os"
)

// LogLevel defines the logging level
type LogLevel int

const (
	// LogOff disables all logging
	LogOff LogLevel = iota
	// LogError logs only errors
	LogError
	// LogWarn logs warnings and errors
	LogWarn
	// LogInfo logs information, warnings and errors
	LogInfo
	// LogDebug logs everything
	LogDebug
)

// Logger defines the logging interface
type Logger interface {
	// Debugf prints debug logs
	Debugf(format string, v ...interface{})
	// Infof prints information logs
	Infof(format string, v ...interface{})
	// Warnf prints warning logs
	Warnf(format string, v ...interface{})
	// Errorf prints error logs
	Errorf(format string, v ...interface{})
}

// defaultLogger is the default implementation of Logger
type defaultLogger struct {
	level  LogLevel
	logger *log.Logger
}

var std Logger = NewDefaultLogger(LogInfo)

// NewDefaultLogger creates a new default logger instance
func NewDefaultLogger(level LogLevel) Logger {
	return &defaultLogger{
		level:  level,
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	if l.level >= LogDebug {
		l.logger.Output(2, fmt.Sprintf("[DEBUG] "+format, v...))
	}
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	if l.level >= LogInfo {
		l.logger.Output(2, fmt.Sprintf("[INFO] "+format, v...))
	}
}

func (l *defaultLogger) Warnf(format string, v ...interface{}) {
	if l.level >= LogWarn {
		l.logger.Output(2, fmt.Sprintf("[WARN] "+format, v...))
	}
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	if l.level >= LogError {
		l.logger.Output(2, fmt.Sprintf("[ERROR] "+format, v...))
	}
}

// Global functions
func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	std.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

func SetLogger(logger Logger) {
	if logger != nil {
		std = logger
	}
}

func SetLevel(level LogLevel) {
	if instance, ok := std.(*defaultLogger); ok {
		instance.level = level
		return
	}
}
