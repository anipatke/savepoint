package v2

import (
	"path/filepath"
	"strings"
	"testing"

	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/testutil"
)

func TestStaleSelectionWithholdsLaunchAdviceOnEverySurface(t *testing.T) {
	for name, router := range map[string]string{
		"missing issue": "state: task\nrelease: R-001\nobjective: O-001\ntask: T-002\nissue: I-999\n",
		"missing task":  "state: task\nrelease: R-001\nobjective: O-001\ntask: T-999\n",
	} {
		t.Run(name, func(t *testing.T) {
			root := writeAdviceProject(t, true)
			testutil.WriteFile(t, filepath.Join(root, "router.md"), "# Router\n\n## Current state\n\n```yaml\n"+router+"```\n")
			resumeText, objectiveText, taskText, plainText := adviceSurfaces(t, root)
			for surface, text := range map[string]string{"resume": resumeText, "objective": objectiveText, "task": taskText, "plain": plainText} {
				if strings.Contains(text, "May start together") || strings.Contains(text, "begin T-") {
					t.Errorf("%s still advertises a launch under a stale selection:\n%s", surface, text)
				}
			}
			for surface, text := range map[string]string{"objective": objectiveText, "task": taskText} {
				if !strings.Contains(text, "not current") {
					t.Errorf("%s does not explain the withheld advice:\n%s", surface, text)
				}
			}
		})
	}
}

func TestLaneHeadingControlsAreStripped(t *testing.T) {
	root := writeLaneProject(t, true)
	lanes := "lanes:\n  - {key: core, title: \"Core\\x1b[2J\\x07\"}\n  - {key: board, title: Board}\n"
	writeFixtureObjective(t, root, "O-001", "First objective", "in_progress", lanes)
	view := openSizedBoard(t, root, 120, 48).View()
	if strings.Contains(view, "\x1b[2J") || strings.Contains(view, "\x07") {
		t.Errorf("board view keeps authored terminal controls: %q", view)
	}
	if !strings.Contains(xansi.Strip(view), "Core") {
		t.Errorf("board lost the heading text:\n%s", view)
	}
}
