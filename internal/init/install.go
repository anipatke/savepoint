package init

import (
	"fmt"
	"os/exec"
)

func InstallDependencies(dir string) error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install failed in %s: %w\noutput: %s", dir, err, string(out))
	}
	return nil
}
