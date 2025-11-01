// @Author: openownworld
// @Email:  openownworld@163.com
// @File:   wrapper.go
// @Description:

package zlog

type zLogWrapper struct {
	logger Logger
}

// Sync 在默认情况下，日志记录器是没有缓冲的。但是在进程退出之前调用 Sync() 方法是一个好习惯。
func (z *zLogWrapper) Sync() error {
	return z.logger.Sync()
}

// SetLevel set the output log level.
func (z *zLogWrapper) SetLevel(level Level) {
	z.logger.SetLevel(level)
}

// Println 打印日志到终端
// var buf strings.Builder
// buf.WriteString(getNowTimeMs())
// buf.WriteString(" console ")
// buf.WriteString(getCallerInfo(consoleSkipNum))
// buf.WriteString(" ")
// fmt.Fprintln(&buf, args...)
// fmt.Print(buf.String())
func (z *zLogWrapper) Println(args ...interface{}) {
	z.Println(args...)
}

// Printfln 打印日志到终端 conslone
func (z *zLogWrapper) Printfln(format string, args ...interface{}) {
	z.Printfln(format, args...)
}

// Printf 打印日志到终端 默认加换行
func (z *zLogWrapper) Printf(format string, args ...interface{}) {
	z.Printf(format, args...)
}

// Debug logs a message at level Debug on the compatibleLogger.
func (z *zLogWrapper) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}

// Debugf logs a message at level Debug on the compatibleLogger.
func (z *zLogWrapper) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args...)
}

// Info logs a message at level Info on the compatibleLogger.
func (z *zLogWrapper) Info(args ...interface{}) {
	z.logger.Info(args...)
}

// Infof logs a message at level Info on the compatibleLogger.
func (z *zLogWrapper) Infof(format string, args ...interface{}) {
	z.logger.Infof(format, args...)
}

// Warn logs a message at level Warn on the compatibleLogger.
func (z *zLogWrapper) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}

// Warnf logs a message at level Warn on the compatibleLogger.
func (z *zLogWrapper) Warnf(format string, args ...interface{}) {
	z.logger.Warnf(format, args...)
}

// Error logs a message at level Error on the compatibleLogger.
func (z *zLogWrapper) Error(args ...interface{}) {
	z.logger.Error(args...)
}

// Errorf logs a message at level Error on the compatibleLogger.
func (z *zLogWrapper) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args...)
}

// Panic logs a message at level Painc on the compatibleLogger.  followed by a call to panic().
func (z *zLogWrapper) Panic(args ...interface{}) {
	z.logger.Panic(args...)
}

// Panicf logs a message at level Painc on the compatibleLogger.
func (z *zLogWrapper) Panicf(format string, args ...interface{}) {
	z.logger.Panicf(format, args...)
}

// Fatal logs a message at level Fatal on the compatibleLogger.
func (z *zLogWrapper) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}

// Fatalf logs a message at level Fatal on the compatibleLogger. followed by a call to os.Exit(1).
func (z *zLogWrapper) Fatalf(format string, args ...interface{}) {
	z.logger.Fatalf(format, args...)
}

// With return a logger with an extra field.
func (z *zLogWrapper) With(fields ...Field) Logger {
	return z.logger.With(fields...)
}

// WithField return a logger with an extra field.
func (z *zLogWrapper) WithField(key string, value interface{}) Logger {
	return z.logger.WithField(key, value)
}

// WithFields return a logger with extra fields.
func (z *zLogWrapper) WithFields(fields map[string]interface{}) Logger {
	return z.logger.WithFields(fields)
}
