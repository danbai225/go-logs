//go:build windows
// +build windows

package go_logs

import (
	"fmt"
	"os"
)

func cuttingOff() {
	if ini {
		ini = false
		return
	}
	infoLog.Close()
	path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log")
	infoLog = empty(path)

	errLog.Close()
	path = fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "err.log")
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
