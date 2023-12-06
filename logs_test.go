package go_logs

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	Debug("Debug")
	Println("Println")
	Info("Info")
	Warn("Warn")
	Err("Err")
	time.Sleep(2 * time.Second)
	Debug("Debug")
	Println("Println")
	Info("Info")
	Warn("Warn")
	Err("Err")
	println("123")
	SetWriteLogs(INFO | ERR | DEBUG)
	time.Sleep(2 * time.Second)
	Debug("Debug")
}
func TestDir(t *testing.T) {
	Debug("Debug-logs")
	Println("Println-logs")
	Info("Info-logs")
	Warn("Warn-logs")
	Err("Err-logs")
	SetLogsDir("logs-2")
	Debug("Debug-logs2")
	Println("Println-logs2")
	Info("Info-logs2")
	Warn("Warn-logs2")
	Err("Err-logs2")
}
func TestJson(t *testing.T) {
	PrintJson(struct {
		Name string
		Age  int64
	}{Name: "test", Age: 18})
}

func TestFlag(t *testing.T) {
	flag := byte(INFO)
	println(flag & INFO)
	flag = flag | ERR
	println(flag & ERR)
	println(flag & INFO)
	SetLevel(INFO)
	Debug("Debug")
	Info("Info")
	Warn("Warn")
	Err("Err")
}
func TestLogF(t *testing.T) {
	DebugF("test %s", "DebugF")
	InfoF("test %s", "InfoF")
	WarnF("test %s", "WarnF")
	ErrF("test %s", "ErrF")
}
