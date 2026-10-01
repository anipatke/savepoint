//go:build windows

package codehealth

import "os/exec"

type execCmd = exec.Cmd

func newExecCmd(s ToolSpec) *exec.Cmd {
	return exec.Command(s.Executable, s.Args...)
}

func processAlive(pid int) bool { return false }
