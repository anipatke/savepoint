//go:build !windows

package codehealth

import (
	"os/exec"
	"syscall"
)

// configureProcess puts the tool in its own process group and, on cancel, kills
// the whole group so children do not outlive the run.
func configureProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
}
