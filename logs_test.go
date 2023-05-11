package go_logs

import (
	"os"
	"sync"
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

func TestPanic(t *testing.T) {
	var std, _ = os.Create("_tempStderr.log")
	redirectStderr(std)
	redirectStdout(std)
	println("1")
	os.Stderr.WriteString("123")
	os.Stdout.WriteString("321")
	panic(2)
}
func TestFlag(t *testing.T) {
	flag := byte(INFO | STD)
	println(flag & STD)
	println(flag & INFO)
	flag = flag | ERR
	println(flag & ERR)
	println(flag & STD)
	println(flag & INFO)
}
func TestSTD(t *testing.T) {
	SetRedirectStdLog()
	Println("123")
	group := &sync.WaitGroup{}
	time.Sleep(time.Second)
	go func() {
		group.Add(1)
		defer func() {
			if r := recover(); r != nil {
				group.Done()
			}
		}()
		panic("panic")
	}()
	os.Stderr.WriteString("err")
	group.Wait()
	println("test")
	panic("panic")
}
