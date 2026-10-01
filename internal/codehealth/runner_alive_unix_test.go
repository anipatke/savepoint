//go:build !windows

package codehealth

import (
	"os/exec"
	"syscall"
)

type execCmd = exec.Cmd

func newExecCmd(s ToolSpec) *exec.Cmd {
	return exec.Command(s.Executable, s.Args...)
}

// processAlive reports whether pid still exists. A zombie child of a dead
// parent is reparented and reaped, so existence is the right signal.
func processAlive(pid int) bool { return syscall.Kill(pid, 0) == nil }
