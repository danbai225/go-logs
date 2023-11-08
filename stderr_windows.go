//go:build windows
// +build windows

package go_logs

import (
	"fmt"
	"os"
)

// cuttingOff 根据一定的条件切割日志文件
func cuttingOff() {
	// ini 是全局变量，用于记录初始状态
	if ini {
		ini = false
		return
	}
	// 如果 infoLog 已经被初始化，则关闭当前文件并创建一个新的文件
	if infoLog != nil {
		_ = infoLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log")
		infoLog = empty(path)
	}
	// 如果 errLog 已经被初始化，则关闭当前文件并创建一个新的文件
	if errLog != nil {
		_ = errLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "err.log")
		errLog = empty(path)
	}

	// 如果 debugLog 已经被初始化，则关闭当前文件并创建一个新的文件
	if debugLog != nil {
		_ = debugLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log")
		debugLog = empty(path)
	}
	// 如果 warnLog 已经被初始化，则关闭当前文件并创建一个新的文件
	if warnLog != nil {
		_ = warnLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log")
		warnLog = empty(path)
	}
}

// empty 删除指定的文件并创建一个新的文件
func empty(path string) *os.File {
	_ = os.Remove(path)
	create, err := os.Create(path)
	if err == nil {
		return create
	}
	return nil
}
