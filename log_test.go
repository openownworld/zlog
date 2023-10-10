package zlog

import (
	"go.uber.org/zap"
	"testing"
)

func BenchmarkLogger(b *testing.B) {
	cfg := GetDefaultConfig()
	cfg.FileLogger = false
	InitLog(cfg)
	// 基准函数会运行目标代码b.N次。
	for i := 0; i < b.N; i++ {
		Error("A", "B")
	}
}

func BenchmarkZap(b *testing.B) {
	logger, _ := zap.NewProduction() // fast
	//logger, _ := zap.NewDevelopment()
	// 基准函数会运行目标代码b.N次。
	for i := 0; i < b.N; i++ {
		logger.Info("只是一个测试")
	}
}

func TestLogger(t *testing.T) {
	Error("A", "B")
	WithField("log", "test").Info("A", "B")
	Info("A", "B")
	Println("A", "B")
}

func TestCfg(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.ShortCaller = false
	cfg.ServiceName = "app"
	InitLog(cfg)
	Error("A", "B")
	With(Field{Key: "describe", Value: "instance"}).Error("describe2")
	WithField("log", "test").Info("A", "B")
}

func TestCfgJson(t *testing.T) {
	cfg := GetDefaultConfig()
	cfg.ConsoleLoggerJSON = true
	cfg.FileLoggerJSON = true
	cfg.ShortCaller = false
	cfg.ErrorFileEnable = true
	InitLog(cfg)
	Error("A", "B")
	WithField("log", "test").Info("A", "B")
}
