//+build windows

package go_logs

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
