//go:build !windows

package run

import (
	"os"
	"os/exec"
	"syscall"
)

// configureProcess puts the agent in its own process group and, on cancel, kills
// the whole group so children do not outlive the run. It mirrors the health
// runner's setup; the two packages stay independent (ARCH-04).
func configureProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return killTree(cmd) }
}

// interruptProcess asks the agent to stop with SIGINT, the way a terminal's
// Ctrl-C does, and reports whether the signal was delivered.
func interruptProcess(cmd *exec.Cmd) bool {
	return cmd.Process.Signal(os.Interrupt) == nil
}

// killTree ends every process in the agent's group. It is safe after the agent
// itself has exited: the group lives on while a grandchild does.
func killTree(cmd *exec.Cmd) error {
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		return cmd.Process.Kill()
	}
	return nil
}

// processAlive reports whether pid still exists. A zombie child of a dead
// parent is reparented and reaped, so existence is the right signal.
func processAlive(pid int) bool { return syscall.Kill(pid, 0) == nil }
