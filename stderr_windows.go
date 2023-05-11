//go:build windows
// +build windows

// Package go_logs 提供日志文件切割和重定向标准输出和标准错误流的功能
package go_logs

import (
	"fmt"
	"os"
	"syscall"
)

var (
	// kernel32.dll 中的 SetStdHandle 函数
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

// setStdHandle 将指定的标准输出或标准错误流重定向到给定的文件句柄
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

// redirectStderr 将标准错误流重定向到指定的文件
func redirectStderr(f *os.File) error {
	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		return err
	}
	os.Stderr = f
	return err
}

// redirectStdout 将标准输出流重定向到指定的文件
func redirectStdout(f *os.File) error {
	err := setStdHandle(syscall.STD_OUTPUT_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		return err
	}
	os.Stdout = f
	return err
}

// cuttingOff 根据一定的条件切割日志文件
func cuttingOff() {
	// ini 是全局变量，用于记录初始状态
	if ini {
		ini = false
		return
	}
	// 如果 infoLog 已经被初始化，则关闭当前文件并创建一个新的文件
	if infoLog != nil {
		infoLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log")
		infoLog = empty(path)
	}
	// 如果 errLog 已经被初始化，则关闭当前文件并创建一个新的文件
	if errLog != nil {
		errLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "err.log")
		errLog = empty(path)
	}

	// 如果 debugLog 已经被初始化，则关闭当前文件并创建一个新的文件
	if debugLog != nil {
		debugLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log")
		debugLog = empty(path)
	}
	// 如果 warnLog 已经被初始化，则关闭当前文件并创建一个新的文件
	if warnLog != nil {
		warnLog.Close()
		path := fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log")
		warnLog = empty(path)
	}
}

// empty 删除指定的文件并创建一个新的文件
func empty(path string) *os.File {
	os.Remove(path)
	create, err := os.Create(path)
	if err == nil {
		return create
	}
	return nil
}
