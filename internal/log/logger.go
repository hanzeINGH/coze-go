package log

import (
	"fmt"
	"log"
	"os"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	// LogOff 关闭日志
	LogOff LogLevel = iota
	// LogError 只记录错误
	LogError
	// LogWarn 记录警告和错误
	LogWarn
	// LogInfo 记录信息、警告和错误
	LogInfo
	// LogDebug 记录所有内容
	LogDebug
)

// Logger 定义了日志接口
type Logger interface {
	// Debugf 打印调试日志
	Debugf(format string, v ...interface{})
	// Infof 打印信息日志
	Infof(format string, v ...interface{})
	// Warnf 打印警告日志
	Warnf(format string, v ...interface{})
	// Errorf 打印错误日志
	Errorf(format string, v ...interface{})
}

// defaultLogger 是默认的日志实现
type defaultLogger struct {
	level  LogLevel
	logger *log.Logger
}

var std Logger = NewDefaultLogger(LogInfo)

// NewDefaultLogger 创建一个新的默认日志实例
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

// 全局函数
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
