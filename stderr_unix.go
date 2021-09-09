//go:build !windows
// +build !windows

package go_logs

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func rewriteStderrFile() error {
	if !StderrFile {
		return nil
	}
	if err := syscall.Dup2(int(stdErrLog.Fd()), int(os.Stderr.Fd())); err != nil {
		return err
	}
	return nil
}
func cuttingOff() {
	if ini {
		ini = false
		return
	}
	exec.Command("bash", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "gin.log"))).Run()
	exec.Command("bash", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log"))).Run()
	exec.Command("bash", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "err.log"))).Run()
	exec.Command("bash", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log"))).Run()
	exec.Command("bash", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log"))).Run()
}
