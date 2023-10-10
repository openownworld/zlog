# zlog

## Config

```go
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
```

## DefaultConfig

```go
    Config{
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
```

## 例子

```go
    zlog.Error("A", "B")
    zlog.WithField("log", "test").Info("A", "B")
    zlog.Info("A", "B")
    zlog.Println("A", "B")
```