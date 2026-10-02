//go:build windows

package codehealth

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// configureProcess ends the tool's whole process tree on cancel. Windows has no
// process groups in the standard library, and Process.Kill ends only the tool,
// so taskkill walks the tree first, while the parent link still exists.
// WaitDelay still bounds how long a survivor can hold the output pipes open.
func configureProcess(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		pid := strconv.Itoa(cmd.Process.Pid)
		exec.Command(taskkill(), "/T", "/F", "/PID", pid).Run()
		return cmd.Process.Kill()
	}
}

func taskkill() string {
	if root := os.Getenv("SystemRoot"); root != "" {
		return filepath.Join(root, "System32", "taskkill.exe")
	}
	return "taskkill"
}
