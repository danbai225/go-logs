//+build !windows

package go_logs

import (
	"os"
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
