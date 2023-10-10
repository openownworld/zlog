package main

import (
	"github.com/openownworld/zlog"
)

func init() {
	zlog.Warn("default logger init")
}

func f() {
	a := 0
	b := 10 / a
	_ = b
}

func main() {
	defer zlog.PrintPanicLog()
	zlog.Error("A", "B")
	f()
}
