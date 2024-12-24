package coze

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Logger ...
type Logger interface {
	Log(ctx context.Context, level LogLevel, message string, args ...interface{})
}

type LevelLogger interface {
	Logger
	SetLevel(level LogLevel)
}

type LogLevel int

// LogLevelTrace ...
const (
	LogLevelTrace LogLevel = iota + 1
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// String ...
func (r LogLevel) String() string {
	switch r {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return ""
	}
}

type stdLogger struct {
	log *log.Logger
}

// NewStdLogger ...
func NewStdLogger() Logger {
	return &stdLogger{
		log: log.New(os.Stderr, "", log.LstdFlags),
	}
}

// Log ...
func (l *stdLogger) Log(ctx context.Context, level LogLevel, message string, args ...interface{}) {
	if len(args) == 0 {
		_ = l.log.Output(2, "["+level.String()+"] "+message)
	} else {
		_ = l.log.Output(2, "["+level.String()+"] "+fmt.Sprintf(message, args...))
	}
}

type levelLogger struct {
	Logger
	level LogLevel
}

// NewLevelLogger ...
func NewLevelLogger(logger Logger, level LogLevel) LevelLogger {
	return &levelLogger{
		Logger: logger,
		level:  level,
	}
}

// SetLevel ...
func (l *levelLogger) SetLevel(level LogLevel) {
	l.level = level
}

// Log ...
func (l *levelLogger) Log(ctx context.Context, level LogLevel, message string, args ...interface{}) {
	if level >= l.level {
		l.Logger.Log(ctx, level, message, args...)
	}
}
