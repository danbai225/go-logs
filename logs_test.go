package go_logs

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	time.Sleep(2 * time.Second)
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
}
