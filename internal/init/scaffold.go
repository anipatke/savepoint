package init

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const ReleaseNumber = "1"

func Scaffold(templates fs.FS, targetDir, projectName string, force bool) error {
	return fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", path, err)
		}

		targetPath := filepath.Join(targetDir, path)

		if d.IsDir() {
			if path == "." {
				return nil
			}
			return os.MkdirAll(targetPath, 0755)
		}

		if path == "." {
			return nil
		}

		content, err := fs.ReadFile(templates, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		interpolated := interpolate(string(content), projectName)

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("create parent dirs for %s: %w", targetPath, err)
		}

		return AtomicWrite(targetPath, []byte(interpolated))
	})
}

func ProjectNameFromDir(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "my-project"
	}
	return filepath.Base(abs)
}

func interpolate(content, projectName string) string {
	result := strings.ReplaceAll(content, "{{PROJECT_NAME}}", projectName)
	result = strings.ReplaceAll(result, "{{RELEASE_NUMBER}}", ReleaseNumber)
	return result
}
