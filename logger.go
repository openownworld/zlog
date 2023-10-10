package zlog

// Level is the log level.
type Level int

// Enums log level constants.
const (
	LevelNil Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

// String turns the LogLevel to string.
func (lv *Level) String() string {
	return LevelStrings[*lv]
}

// LevelStrings is the map from log level to its string representation.
var LevelStrings = map[Level]string{
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelWarn:  "warn",
	LevelError: "error",
	LevelPanic: "panic",
	LevelFatal: "fatal",
}

// LevelNames is the map from string to log level.
var LevelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"panic": LevelPanic,
	"fatal": LevelFatal,
}

// Field is the user defined log field.
type Field struct {
	Key   string
	Value interface{}
}

// Logger is the underlying logging work for tRPC framework.
type Logger interface {
	// Println print console
	Println(args ...interface{})
	// Printf print console
	Printf(format string, args ...interface{})
	// Debug logs to DEBUG log. Arguments are handled in the manner of fmt.Print.
	Debug(args ...interface{})
	// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
	Debugf(format string, args ...interface{})
	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})
	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})
	// Warn logs to WARN log. Arguments are handled in the manner of fmt.Print.
	Warn(args ...interface{})
	// Warnf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warnf(format string, args ...interface{})
	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})
	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface{})
	// Panic logs to PANIC log. Arguments are handled in the manner of fmt.Print.
	Panic(args ...interface{})
	// Panicf logs to PANIC log. Arguments are handled in the manner of fmt.Printf.
	Panicf(format string, args ...interface{})
	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// All Fatal logs will exit by calling os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatal(args ...interface{})
	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Fatalf(format string, args ...interface{})
	// Sync calls the underlying Core's Sync method, flushing any buffered log entries.
	// Applications should take care to call Sync before exiting.
	// 在默认情况下，日志记录器是没有缓冲的。但是在进程退出之前调用 Sync() 方法是一个好习惯。
	Sync() error
	// SetLevel set the output log level.
	SetLevel(level Level)
	// WithFields set some user defined data to logs, such as uid, imei, etc.
	// Fields must be paired.
	WithFields(fields map[string]interface{}) Logger
	// WithField field
	WithField(key string, value interface{}) Logger
	// With add user defined fields to Logger. Fields support multiple values.
	With(fields ...Field) Logger
}
