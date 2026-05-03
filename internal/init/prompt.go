package init

import (
	"fmt"
	"io/fs"
)

func RenderMagicPrompt(templates fs.FS, projectName string) (string, error) {
	content, err := fs.ReadFile(templates, "magic-prompt.prompt.md")
	if err != nil {
		return "", fmt.Errorf("read magic prompt template: %w", err)
	}
	return interpolate(string(content), projectName), nil
}
