//go:build windows

package run

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// configureProcess ends the agent's whole process tree on cancel. Windows has no
// process groups in the standard library, so taskkill walks the tree first,
// while the parent link still exists. It mirrors the health runner's setup.
func configureProcess(cmd *exec.Cmd) {
	cmd.Cancel = func() error { return killTree(cmd) }
}

// interruptProcess reports false: Windows has no signal that reaches a console
// child, so a stop goes straight to the tree kill.
func interruptProcess(*exec.Cmd) bool { return false }

// killTree ends the agent and everything it started. It must run while the
// agent is alive, since taskkill finds children through the parent link.
func killTree(cmd *exec.Cmd) error {
	pid := strconv.Itoa(cmd.Process.Pid)
	exec.Command(taskkill(), "/T", "/F", "/PID", pid).Run()
	return cmd.Process.Kill()
}

func taskkill() string {
	if root := os.Getenv("SystemRoot"); root != "" {
		return filepath.Join(root, "System32", "taskkill.exe")
	}
	return "taskkill"
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
