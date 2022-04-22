//go:build windows
// +build windows

package go_logs

import (
	"fmt"
	"os"
	"syscall"
)

func cuttingOff() {
	if ini {
		ini = false
		return
	}
	if infoLog != nil {
		infoLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log")
		infoLog = empty(path)
	}
	if errLog != nil {
		errLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "err.log")
		errLog = empty(path)
	}

	if debugLog != nil {
		debugLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log")
		debugLog = empty(path)
	}
	if warnLog != nil {
		warnLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log")
		warnLog = empty(path)
	}

}
func empty(path string) *os.File {
	os.Remove(path)
	create, err := os.Create(path)
	if err == nil {
		return create
	}
	return nil
}

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func setStdHandle(stdhandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) error {
	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		return err
	}
	os.Stderr = f
	return err
}

// redirectStdout to the file passed in
func redirectStdout(f *os.File) error {
	err := setStdHandle(syscall.STD_OUTPUT_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		return err
	}
	os.Stdout = f
	return err
}
