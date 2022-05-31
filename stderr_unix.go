//go:build !windows
// +build !windows

package go_logs

import (
	"fmt"
	"os"
	"os/exec"
)

func cuttingOff() {
	if ini {
		ini = false
		return
	}
	if infoLog != nil {
		_ = exec.Command("sh", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log"))).Run()
	}
	if errLog != nil {
		_ = exec.Command("sh", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "err.log"))).Run()
	}
	if debugLog != nil {
		_ = exec.Command("sh", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log"))).Run()
	}
	if warnLog != nil {
		_ = exec.Command("sh", "-c", fmt.Sprintf("cp /dev/null %s", fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log"))).Run()
	}
}
