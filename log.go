// @Author: openownworld
// @Email:  openownworld@163.com
// @File:   wrapper.go
// @Description:

package zlog

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

const (
	maxStack  = 20
	separator = "capture panic---------------------------------------"
)

// PrintPanicLog recover 打印堆栈信息
func PrintPanicLog() {
	if err := recover(); err != nil {
		t := time.Now().Format("2006-01-02 15:04:05.000") + " "
		str := fmt.Sprintf("\n%s%s start\nruntime error: %v\ntraceback:\n", t, separator, err)
		i := 2
		for {
			pc, file, line, ok := runtime.Caller(i)
			if !ok || i > maxStack {
				break
			}
			str += fmt.Sprintf("\tstack: %d %v [ file: %s:%d ] func: %s\n", i-1, ok, file, line,
				runtime.FuncForPC(pc).Name())
			i++
		}
		str += t + separator + " end\n" + string(debug.Stack())
		Error(str)
		//debug.PrintStack()
	}
}

// Sync calls the underlying Core's Sync method, flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func Sync() error {
	return GetDefaultLogger().Sync()
}

// SetLevel set the output log level.
func SetLevel(level Level) {
	GetDefaultLogger().SetLevel(level)
}

// Println 打印日志到终端
func Println(args ...interface{}) {
	GetDefaultLogger().Println(args...)
}

// Printf 打印日志到终端
func Printf(format string, args ...interface{}) {
	GetDefaultLogger().Printf(format, args...)
}

// Debug logs a message at level Debug on the compatibleLogger.
func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

// Debugf logs a message at level Debug on the compatibleLogger.
func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

// Info logs a message at level Info on the compatibleLogger.
func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

// Infof logs a message at level Info on the compatibleLogger.
func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

// Warn logs a message at level Warn on the compatibleLogger.
func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

// Warnf logs a message at level Warn on the compatibleLogger.
func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warnf(format, args...)
}

// Error logs a message at level Error on the compatibleLogger.
func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

// Errorf logs a message at level Error on the compatibleLogger.
func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

// Panic logs a message at level Painc on the compatibleLogger.  followed by a call to panic().
func Panic(args ...interface{}) {
	GetDefaultLogger().Panic(args...)
}

// Panicf logs a message at level Painc on the compatibleLogger.
func Panicf(format string, args ...interface{}) {
	GetDefaultLogger().Panicf(format, args...)
}

// Fatal logs a message at level Fatal on the compatibleLogger.
func Fatal(args ...interface{}) {
	GetDefaultLogger().Fatal(args...)
}

// Fatalf logs a message at level Fatal on the compatibleLogger. followed by a call to os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	GetDefaultLogger().Fatalf(format, args...)
}

// With return a logger with an extra field.
func With(fields ...Field) Logger {
	return GetDefaultLogger().With(fields...)
}

// WithField return a logger with an extra field.
func WithField(key string, value interface{}) Logger {
	return GetDefaultLogger().WithField(key, value)
}

// WithFields return a logger with extra fields.
func WithFields(fields map[string]interface{}) Logger {
	return GetDefaultLogger().WithFields(fields)
}
