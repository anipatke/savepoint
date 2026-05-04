package board

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// RunTUI launches the Bubble Tea board with optional release/epic filters.
func RunTUI(release, epic string) error {
	model, err := newProjectModel(".", release, epic)
	if err != nil {
		return err
	}

	cfg, err := model.Dependencies.ConfigReader.Read(filepath.Join(model.Root, "config.yml"))
	if err != nil {
		return err
	}
	model.Theme = cfg.Theme

	applyColorProfile()

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}
