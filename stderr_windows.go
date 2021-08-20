//+build windows

package go_logs

import (
	"fmt"
	"os"
	"syscall"
)

const (
	kernel32dll = "kernel32.dll"
)

func rewriteStderrFile() error {
	if !StderrFile {
		return nil
	}
	kernel32 := syscall.NewLazyDLL(kernel32dll)
	setStdHandle := kernel32.NewProc("SetStdHandle")
	sh := syscall.STD_ERROR_HANDLE
	v, _, err := setStdHandle.Call(uintptr(sh), uintptr(stdErrLog.Fd()))
	if v == 0 {
		return err
	}
	return nil
}
func cuttingOff() {
	if ini {
		ini = false
		return
	}
	infoLog.Close()
	path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log")
	infoLog = empty(path)

	errLog.Close()
	path = fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "error.log")
	errLog = empty(path)

	debugLog.Close()
	path = fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log")
	debugLog = empty(path)

	warnLog.Close()
	path = fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log")
	warnLog = empty(path)
	warnLog, _ = os.Open(path)
}
func empty(path string) *os.File {
	os.Remove(path)
	create, err := os.Create(path)
	if err == nil {
		return create
	}
	return nil
}
