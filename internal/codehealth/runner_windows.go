//go:build windows

package codehealth

import "os/exec"

// configureProcess kills the tool on cancel. Windows has no process groups in
// the standard library; WaitDelay still bounds how long a surviving child can
// hold the output pipes open.
func configureProcess(cmd *exec.Cmd) {
	cmd.Cancel = func() error { return cmd.Process.Kill() }
}
