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
	Debug("Debug")
	Println("Println")
	Info("Info")
	Warn("Warn")
	Err("Err")
	SetLogsDir("logs-2")
	time.Sleep(2 * time.Second)
	Debug("Debug")
	Println("Println")
	Info("Info")
	Warn("Warn")
	Err("Err")
}
func TestJson(t *testing.T) {
	PrintJson(struct {
		Name string
		Age  int64
	}{Name: "test", Age: 18})
}
