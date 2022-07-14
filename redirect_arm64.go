//go:build !windows && !darwin && arm64
// +build !windows,!darwin,arm64

package go_logs

import (
	"os"
	"syscall"
)

func redirectStderr(f *os.File) error {
	if err := syscall.Dup3(int(f.Fd()), int(os.Stderr.Fd()), 0); err != nil {
		return err
	}
	return nil
}
func redirectStdout(f *os.File) error {
	if err := syscall.Dup3(int(f.Fd()), int(os.Stdout.Fd()), 0); err != nil {
		return err
	}
	return nil
}
