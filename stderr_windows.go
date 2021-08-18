//+build windows

package go_logs

import "syscall"

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
	exec.Command(fmt.Sprintf("@echo.>%s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "gin.log"))).Run()
	exec.Command(fmt.Sprintf("@echo.>%s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log"))).Run()
	exec.Command(fmt.Sprintf("@echo.>%s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "error.log"))).Run()
	exec.Command(fmt.Sprintf("@echo.>%s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log"))).Run()
	exec.Command(fmt.Sprintf("@echo.>%s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log"))).Run()
	exec.Command(fmt.Sprintf("@echo.>%s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "stdErr.log"))).Run()
}
