// @Author: openownworld
// @Email:  openownworld@163.com
// @File:   wrapper.go
// @Description:

// Package log zap的特性
//高性能：zap 对日志输出进行了多项优化以提高它的性能
//日志分级：有 Debug，Info，Warn，Error，DPanic，Panic，Fatal 等
//日志记录结构化：日志内容记录是结构化的，比如 json 格式输出
//自定义格式：用户可以自定义输出的日志格式
//自定义公共字段：用户可以自定义公共字段，输出的日志内容就共同拥有了这些字段
//调试：可以打印文件名、函数名、行号、日志时间等，便于调试程序
//自定义调用栈级别：可以根据日志级别输出它的调用栈信息
//Namespace：日志命名空间。定义命名空间后，所有日志内容就在这个命名空间下。命名空间相当于一个文件夹
//支持 hook 操作

package zlog

import (
	"fmt"
	"github.com/go-ini/ini"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Levels is the map from string to zapcore.Level.
var Levels = map[string]zapcore.Level{
	"":      zapcore.DebugLevel,
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"panic": zapcore.PanicLevel,
	"fatal": zapcore.FatalLevel,
}

var levelToZapLevel = map[Level]zapcore.Level{
	LevelDebug: zapcore.DebugLevel,
	LevelInfo:  zapcore.InfoLevel,
	LevelWarn:  zapcore.WarnLevel,
	LevelError: zapcore.ErrorLevel,
	LevelPanic: zapcore.PanicLevel,
	LevelFatal: zapcore.FatalLevel,
}

var (
	atomLevel      zap.AtomicLevel
	callerSkipNum  = 2
	consoleSkipNum = 3
	defaultLogger  Logger
	mtx            sync.RWMutex
)

// Config 封装高性能日志库zap
// 支持按最大天数保存，最大文件限制，最大文件数限制
// 支持错误日志分级 error level 能复制提取保存文件，支持动态设置日志分级，支持日志压缩
type Config struct {
	ServiceKey         string `ini:"serviceKey"`         // json service key
	ServiceName        string `ini:"serviceName"`        // json service name
	CustomTimeEnable   bool   `ini:"customTimeEnable"`   // custom time 2006-01-02 15:04:05.000
	LogFileName        string `ini:"logFileName"`        // all日志输出路径文件名
	ErrorFileName      string `ini:"errorFileName"`      // 错误日志分级复制输出路径文件名
	MaxSize            int    `ini:"maxSize"`            // Mb 最大文件限制，最大文件数限制
	MaxBackups         int    `ini:"maxBackups"`         // 最大文件数限制
	MaxDays            int    `ini:"maxDays"`            // 最大天数保存
	Compress           bool   `ini:"compress"`           // 启用日志压缩
	Level              string `ini:"level"`              // 日志级别
	StacktraceLevel    string `ini:"stacktraceLevel"`    // 输出调用堆栈 级别
	ErrorFileLevel     string `ini:"errorFileLevel"`     // 错误日志分级 级别
	ShortCaller        bool   `ini:"shortCaller"`        // 文件名行号 log/log.go:127 or 全路径
	FunctionEnable     bool   `ini:"functionEnable"`     // 函数路径 go-demo/libary/log.TestLogger
	SocketType         string `ini:"socketType"`         // socket type UDP
	SocketIP           string `ini:"socketIP"`           // server dst ip
	SocketPort         string `ini:"socketPort"`         // server dst port
	SocketLoggerEnable bool   `ini:"socketLoggerEnable"` // 启用 socket Logger
	SocketLoggerJSON   bool   `ini:"socketLoggerJSON"`   // 启用 socket LoggerJSON
	ErrorFileEnable    bool   `ini:"errorFileEnable"`    // 启用 错误日志分级复制输出
	FileLogger         bool   `ini:"fileLogger"`         // 启用 file Logger
	FileLoggerJSON     bool   `ini:"fileLoggerJSON"`     // 启用 file LoggerJSON
	ConsoleLogger      bool   `ini:"consoleLogger"`      // 启用 console Logger
	ConsoleLoggerJSON  bool   `ini:"consoleLoggerJSON"`  // 启用 console LoggerJSON
}

// InitLogByFile 确保日志最先初始化 log.ini
func InitLogByFile(filename string) error {
	runDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := path.Join(runDir, filename)
	var logConfig Config
	p, err := ini.Load(dir)
	if err != nil {
		return fmt.Errorf("open config file dir [%s] failed: %v", dir, err)
	}
	if err := p.MapTo(logConfig); err != nil {
		return err
	}
	return InitLog(logConfig)
}

// InitLogByReader 确保日志最先初始化
func InitLogByReader(reader io.Reader) error {
	var logConfig Config
	p, err := ini.Load(reader)
	if err != nil {
		return fmt.Errorf("load config stream failed: %v", err)
	}
	if err := p.MapTo(logConfig); err != nil {
		return err
	}
	return InitLog(logConfig)
}

// InitLog 确保日志最先初始化
func InitLog(config Config) error {
	t := NewZLogger(config)
	mtx.Lock()
	defaultLogger = t
	mtx.Unlock()
	return nil
}

// GetDefaultConfig  默认 Config
func GetDefaultConfig() Config {
	return Config{
		LogFileName:        "./logs/log.log",
		ErrorFileName:      "./logs/error.log",
		MaxSize:            20, // Mb
		MaxBackups:         15,
		MaxDays:            15,
		Compress:           true,
		Level:              "debug",
		StacktraceLevel:    "panic",
		ErrorFileLevel:     "error",
		ShortCaller:        true,
		FunctionEnable:     false,
		FileLogger:         true,
		ErrorFileEnable:    false,
		ConsoleLogger:      true,
		FileLoggerJSON:     false,
		ConsoleLoggerJSON:  false,
		SocketLoggerEnable: false,
		SocketLoggerJSON:   false,
		SocketType:         "udp",
		SocketIP:           "127.0.0.1",
		SocketPort:         "9990",
	}
}

// SetDefaultLogger implements
func SetDefaultLogger(logger Logger) {
	mtx.Lock()
	defaultLogger = logger
	mtx.Unlock()
}

// GetDefaultLogger defaults logger
func GetDefaultLogger() Logger {
	mtx.RLock()
	l := defaultLogger
	mtx.RUnlock()
	if l != nil {
		return l
	}
	t := NewZLogger(GetDefaultConfig())
	mtx.Lock()
	defaultLogger = t
	mtx.Unlock()
	return t
}

type zLogger struct {
	logger *zap.Logger
}

// NewZLogger creates a new logger
func NewZLogger(logConfig Config) Logger {
	return &zLogger{
		logger: getLogger(logConfig),
	}
}

// Sync 在默认情况下，日志记录器是没有缓冲的。但是在进程退出之前调用 Sync() 方法是一个好习惯。
func (z *zLogger) Sync() error {
	return z.logger.Sync()
}

// SetLevel set the output log level.
func (z *zLogger) SetLevel(level Level) {
	v := levelToZapLevel[level]
	atomLevel.SetLevel(v)
}

// Println 打印日志到终端
func (z *zLogger) Println(args ...interface{}) {
	fmt.Printf("%s %s %s %s", getNowTimeMs(), "console", getCallerInfo(consoleSkipNum), fmt.Sprintln(args...))
}

// Printfln 打印日志到终端 conslone
func (z *zLogger) Printfln(format string, args ...interface{}) {
	fmt.Printf("%s %s %s %s\n", getNowTimeMs(), "console", getCallerInfo(consoleSkipNum), fmt.Sprintf(format, args...))
}

// Printf 打印日志到终端 默认加换行
func (z *zLogger) Printf(format string, args ...interface{}) {
	fmt.Printf("%s %s %s %s\n", getNowTimeMs(), "console", getCallerInfo(consoleSkipNum), fmt.Sprintf(format, args...))
}

// Debug logs a message at level Debug on the compatibleLogger.
func (z *zLogger) Debug(args ...interface{}) {
	z.logger.Debug(getLogMsg(args...))
}

// Debugf logs a message at level Debug on the compatibleLogger.
func (z *zLogger) Debugf(format string, args ...interface{}) {
	z.logger.Debug(fmt.Sprintf(format, args...))
}

// Info logs a message at level Info on the compatibleLogger.
func (z *zLogger) Info(args ...interface{}) {
	z.logger.Info(getLogMsg(args...))
}

// Infof logs a message at level Info on the compatibleLogger.
func (z *zLogger) Infof(format string, args ...interface{}) {
	z.logger.Info(fmt.Sprintf(format, args...))
}

// Warn logs a message at level Warn on the compatibleLogger.
func (z *zLogger) Warn(args ...interface{}) {
	z.logger.Warn(getLogMsg(args...))
}

// Warnf logs a message at level Warn on the compatibleLogger.
func (z *zLogger) Warnf(format string, args ...interface{}) {
	z.logger.Warn(fmt.Sprintf(format, args...))
}

// Error logs a message at level Error on the compatibleLogger.
func (z *zLogger) Error(args ...interface{}) {
	z.logger.Error(getLogMsg(args...))
}

// Errorf logs a message at level Error on the compatibleLogger.
func (z *zLogger) Errorf(format string, args ...interface{}) {
	z.logger.Error(fmt.Sprintf(format, args...))
}

// Panic logs a message at level Painc on the compatibleLogger.  followed by a call to panic().
func (z *zLogger) Panic(args ...interface{}) {
	z.logger.Panic(getLogMsg(args...))
}

// Panicf logs a message at level Painc on the compatibleLogger.
func (z *zLogger) Panicf(format string, args ...interface{}) {
	z.logger.Panic(fmt.Sprintf(format, args...))
}

// Fatal logs a message at level Fatal on the compatibleLogger.
func (z *zLogger) Fatal(args ...interface{}) {
	z.logger.Fatal(getLogMsg(args...))
}

// Fatalf logs a message at level Fatal on the compatibleLogger. followed by a call to os.Exit(1).
func (z *zLogger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal(fmt.Sprintf(format, args...))
}

// With return a logger with an extra field.
func (z *zLogger) With(fields ...Field) Logger {
	f := make([]zap.Field, len(fields))
	i := 0
	for _, v := range fields {
		f[i] = zap.Any(v.Key, v.Value)
		i++
	}
	n := *z
	n.logger = n.logger.With(f...)
	return &zLogWrapper{logger: &n}
}

// WithField return a logger with an extra field.
func (z *zLogger) WithField(key string, value interface{}) Logger {
	n := *z
	n.logger = n.logger.With(zap.Any(key, value))
	return &zLogWrapper{logger: &n}
}

// WithFields return a logger with extra fields.
func (z *zLogger) WithFields(fields map[string]interface{}) Logger {
	f := make([]zap.Field, len(fields))
	i := 0
	for k, v := range fields {
		f[i] = zap.Any(k, v)
		i++
	}
	n := *z
	n.logger = n.logger.With(f...)
	return &zLogWrapper{logger: &n}
}

func getLogMsg(args ...interface{}) string {
	// strings.TrimRight(fmt.Sprintln(args...), "\n")
	s := fmt.Sprintln(args...)
	if len(s) == 0 {
		return s
	}
	s = s[:len(s)-1]
	return s
}

func getLogMsgf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// EncoderOption option
type EncoderOption struct {
	formatter, timeFmt                string
	colorLevel, shortCaller, function bool
}

func newEncoder(op EncoderOption) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "name",
		CallerKey:     "caller",
		FunctionKey:   "func",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		LineEnding:    zapcore.DefaultLineEnding,
		// zapcore.LowercaseColorLevelEncoder
		// zapcore.LowercaseLevelEncoder // 小写编码器
		// zapcore.CapitalLevelEncoder
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     newTimeEncoder(op.timeFmt),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if op.colorLevel {
		encoderCfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}
	if op.shortCaller {
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	}
	if !op.function {
		encoderCfg.FunctionKey = ""
	}
	switch op.formatter {
	case "console":
		return zapcore.NewConsoleEncoder(encoderCfg)
	case "json":
		return zapcore.NewJSONEncoder(encoderCfg)
	default:
		return zapcore.NewConsoleEncoder(encoderCfg)
	}
}

// newTimeEncoder creates a time format encoder.
func newTimeEncoder(format string) zapcore.TimeEncoder {
	switch format {
	case "seconds":
		return zapcore.EpochTimeEncoder
	case "milliseconds":
		return zapcore.EpochMillisTimeEncoder
	case "nanoseconds":
		return zapcore.EpochNanosTimeEncoder
	case "utc":
		return zapcore.ISO8601TimeEncoder // ISO8601 UTC 时间格式
	default:
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			//enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000000") + "]")
			//enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
	}
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func getLevel(loglevel string) zapcore.Level {
	var level zapcore.Level
	switch loglevel {
	case "debug":
		// DebugLevel logs are typically voluminous, and are usually disabled in production.
		level = zap.DebugLevel
	case "info":
		// InfoLevel is the default logging priority.
		level = zap.InfoLevel
	case "warn":
		// WarnLevel logs are more important than Info, but don't need individual human review.
		level = zap.WarnLevel
	case "error":
		// ErrorLevel logs are high-priority. If an application is running smoothly,
		// it shouldn't generate any error-level logs.
		level = zap.ErrorLevel
	case "panic":
		// PanicLevel logs a message, then panics.
		level = zap.PanicLevel
	case "fatal":
		// FatalLevel logs a message, then calls os.Exit(1).
		level = zap.FatalLevel
	default:
		level = zap.DebugLevel
	}
	return level
}

func getNowTimeMs() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// getCallerInfo 0为当前栈 1上层调用栈
func getCallerInfo(stackNum int) string {
	//0为当前栈 1上层调用栈
	//pc, file, line, ok := runtime.Caller(stackNum)
	_, file, line, ok := runtime.Caller(stackNum)
	if !ok {
		file = "runtime call err"
	}
	//pcName := runtime.FuncForPC(pc).Name() //获取函数名，听说很耗时
	//where:= fmt.Sprintf("%v %s %d %t %s",pc,file,line,ok,pcName)
	//where := fmt.Sprintf("%s %d %s", file, line, pcName)
	where := fmt.Sprintf("%s:%d", file, line) //冒号拼接，goland可以直接点开到文件行数
	return where
}

func getLogger(logConfig Config) *zap.Logger {
	op := EncoderOption{timeFmt: "json", colorLevel: false, shortCaller: logConfig.ShortCaller,
		function: logConfig.FunctionEnable}
	//[1]文件log hook MaxBackups和MaxAge 任意达到限制，对应的文件就会被清理
	hookAll := lumberjack.Logger{
		Filename:   logConfig.LogFileName, // 日志文件路径 ./log/log.log
		MaxSize:    logConfig.MaxSize,     // 最大文件大小 M字节
		MaxBackups: logConfig.MaxBackups,  // 最多保留3个备份
		MaxAge:     logConfig.MaxDays,     // 文件最多保存多少天
		Compress:   logConfig.Compress,    // 是否压缩 disabled by default
	}
	hookError := lumberjack.Logger{
		Filename:   logConfig.ErrorFileName, // 日志文件路径 ./log/error.log
		MaxSize:    logConfig.MaxSize,       // 最大文件大小 M字节
		MaxBackups: logConfig.MaxBackups,    // 最多保留3个备份
		MaxAge:     logConfig.MaxDays,       // 文件最多保存多少天
		Compress:   logConfig.Compress,      // 是否压缩 disabled by default
	}
	//[2]设置level 动态level
	atomLevel = zap.NewAtomicLevel()
	atomLevel.SetLevel(getLevel(logConfig.Level))
	var errorLevel zapcore.Level
	if logConfig.ErrorFileLevel == "error" {
		errorLevel = zapcore.ErrorLevel
	} else {
		errorLevel = zapcore.WarnLevel
	}
	//[3]配置多个输出方式
	//TCP 没解决 当服务器断开，重连的问题
	//UDP 不存在这样的问题
	var socketCore zapcore.Core
	if logConfig.SocketLoggerEnable {
		addr := fmt.Sprintf("%s:%s", logConfig.SocketIP, logConfig.SocketPort)
		conn, err := net.DialTimeout("udp", addr, 3*time.Second)
		if err != nil {
			fmt.Println("err", logConfig.SocketType, addr, err.Error())
		} else {
			// read or write on conn
			//defer conn.Close()
			wSocket := zapcore.AddSync(conn)
			if logConfig.SocketLoggerJSON {
				socketEncoder := newEncoder(op)
				socketCore = zapcore.NewCore(socketEncoder, wSocket, atomLevel)
			} else {
				op.formatter = ""
				socketEncoder := newEncoder(op)
				socketCore = zapcore.NewCore(socketEncoder, wSocket, atomLevel)
			}
		}
	}
	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	var consoleWriter zapcore.WriteSyncer
	var allWriter zapcore.WriteSyncer
	var errorWriter zapcore.WriteSyncer
	if logConfig.ConsoleLogger {
		consoleWriter = zapcore.Lock(os.Stdout)
	} else {
		consoleWriter = zapcore.AddSync(ioutil.Discard)
	}
	if logConfig.FileLogger {
		allWriter = zapcore.AddSync(&hookAll)
		if logConfig.ErrorFileEnable {
			errorWriter = zapcore.AddSync(&hookError)
		} else {
			errorWriter = zapcore.AddSync(ioutil.Discard)
		}
	} else {
		allWriter = zapcore.AddSync(ioutil.Discard)
		errorWriter = zapcore.AddSync(ioutil.Discard)
	}
	// Optimize the Kafka output for machine consumption and the console output for human operators.
	op.formatter = ""
	fileEncoder := newEncoder(op)
	op.formatter = "json"
	jsonEncoder := newEncoder(op)
	if logConfig.FileLoggerJSON {
		fileEncoder = jsonEncoder
	}
	//终端支持彩色打印
	op.formatter = ""
	op.colorLevel = true
	consoleEncoder := newEncoder(op)
	if logConfig.ConsoleLoggerJSON {
		consoleEncoder = jsonEncoder
	}
	// Join the outputs, encoders, and level-handling functions into zapcore.Cores, then tee the four cores together.
	var core zapcore.Core
	if socketCore != nil {
		core = zapcore.NewTee(
			socketCore,
			zapcore.NewCore(fileEncoder, allWriter, atomLevel),
			zapcore.NewCore(fileEncoder, errorWriter, errorLevel),
			zapcore.NewCore(consoleEncoder, consoleWriter, atomLevel),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, allWriter, atomLevel),
			zapcore.NewCore(fileEncoder, errorWriter, errorLevel),
			zapcore.NewCore(consoleEncoder, consoleWriter, atomLevel),
		)
	}
	//[4]创建日志logger 设置初始化字段 service key
	filed := zap.Fields()
	if logConfig.ServiceKey == "" {
		logConfig.ServiceKey = "service"
	}
	if len(logConfig.ServiceName) != 0 {
		filed = zap.Fields(zap.String("service", logConfig.ServiceName))
	}
	// zap.Logger.Info("") 为 0 层
	// With 调用链使用的 Info 接口 ，比直接 Info 少一层 , With需要 we can add a layer to the debug
	//series function calls, so that the caller information can be set correctly.
	logger := zap.New(core, zap.Development(), zap.AddCaller(), zap.AddCallerSkip(callerSkipNum), zap.AddStacktrace(getLevel(logConfig.StacktraceLevel)), filed)
	//输出调用堆栈 主要是调用函数 zap.AddStacktrace()
	//defer logger.Sync()
	//logger.Info("default logger init " + logConfig.Level + " success")
	//logger.Error("default logger init " + logConfig.Level + " success")
	return logger
}
