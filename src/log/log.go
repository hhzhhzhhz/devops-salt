package log

import (
	"fmt"
	"github.com/devops-salt/src/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"io"
)

// Level log level
type Level int32

// Logger level
const (
	LvlDebug Level = iota
	LvlInfo
	LvlWarn
	LvlError
	LvlFatal
)

// Logger stores a logger
type Logger struct {
	log.Logger
	level Level
}

var stdLogger = New(os.Stdout)

// SetOutput sets the writer of standard logger
func SetOutput(w io.Writer) {
	stdLogger.SetOutput(w)
}

// SetLevel set log level
func SetLevel(l Level) {
	stdLogger.SetLevel(l)
}

// GetLevel returns current log level
func GetLevel() Level {
	return stdLogger.Level()
}

// Fatal fatal
func Fatal(format string, v ...interface{}) {
	stdLogger.Output(2, fmt.Sprintf("[FATAL] "+format+"\n", v...))
	os.Exit(1)
}

// Error error
func Error(format string, v ...interface{}) {
	if stdLogger.Level() <= LvlError {
		stdLogger.Output(2, fmt.Sprintf("[ERROR] "+format+"\n", v...))
	}
}

// Warn warn
func Warn(format string, v ...interface{}) {
	if stdLogger.Level() <= LvlWarn {
		stdLogger.Output(2, fmt.Sprintf("[WARN] "+format+"\n", v...))
	}
}

// Info info
func Info(format string, v ...interface{}) {
	if stdLogger.Level() <= LvlInfo {
		stdLogger.Output(2, fmt.Sprintf("[INFO] "+format+"\n", v...))
	}
}

// Debug debug
func Debug(format string, v ...interface{}) {
	if stdLogger.Level() <= LvlDebug {
		stdLogger.Output(2, fmt.Sprintf("[DEBUG] "+format+"\n", v...))
	}
}

// New creates a instance of Logger
func New(w io.Writer) *Logger {
	l := &Logger{level: LvlInfo}
	l.SetOutput(w)
	l.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return l
}

// Debug level
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LvlDebug {
		l.Output(2, fmt.Sprintf("[DEBUG] "+format+"\n", v...))
	}
}

// Info info
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LvlInfo {
		l.Output(2, fmt.Sprintf("[INFO] "+format+"\n", v...))
	}
}

// Warn warn
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= LvlWarn {
		l.Output(2, fmt.Sprintf("[WARN] "+format+"\n", v...))
	}
}

// Error error
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LvlError {
		l.Output(2, fmt.Sprintf("[ERROR] "+format+"\n", v...))
	}
}

// Fatal fatal
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf("[FATAL] "+format+"\n", v...))
	os.Exit(1)
}

// Level returns current logger level
func (l *Logger) Level() Level {
	return l.level
}

// SetLevel sets the logger level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func init() {
	stdLogger.SetOutput(&lumberjack.Logger{
		Filename:   config.GetLogFile(),
		MaxSize:    config.GetLogMaxSize(),
		MaxBackups: config.GetLogMaxBackups(),
		MaxAge:     config.GetLogMaxDuration(),
	})

	stdLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	switch config.GetLogLevel() {
	case "debug", "DEBUG":
		stdLogger.SetLevel(LvlDebug)
	case "warn", "WARN":
		stdLogger.SetLevel(LvlWarn)
	case "error", "ERROR":
		stdLogger.SetLevel(LvlError)
	default:
		stdLogger.SetLevel(LvlInfo)
	}
}
