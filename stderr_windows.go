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
