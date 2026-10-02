//go:build windows

package codehealth

import (
	"os/exec"
	"syscall"
)

type execCmd = exec.Cmd

func newExecCmd(s ToolSpec) *exec.Cmd {
	return exec.Command(s.Executable, s.Args...)
}

const (
	processSynchronize = 0x00100000
	waitTimeout        = 258
)

// processAlive reports whether pid is still running: a process that cannot be
// opened is gone, and one whose handle is not yet signalled is running.
func processAlive(pid int) bool {
	h, err := syscall.OpenProcess(processSynchronize, false, uint32(pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(h)
	state, err := syscall.WaitForSingleObject(h, 0)
	return err == nil && state == waitTimeout
}
