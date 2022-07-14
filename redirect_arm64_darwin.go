//go:build darwin && arm64
// +build darwin,arm64

package go_logs

import (
	"os"
	"syscall"
)

func redirectStderr(f *os.File) error {
	if err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd())); err != nil {
		return err
	}
	return nil
}
func redirectStdout(f *os.File) error {
	if err := syscall.Dup2(int(f.Fd()), int(os.Stdout.Fd())); err != nil {
		return err
	}
	return nil
}
