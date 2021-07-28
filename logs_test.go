package go_logs

import (
	"testing"
)

func Test(t *testing.T) {
	Debug("Debug")
	Println("Println")
	Info("Info")
	Warn("Warn")
	_ = Err("Err")
}
