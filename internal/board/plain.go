package board

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

const plainNonTTYWarning = "[non-interactive mode — run in a TTY to launch the board UI]"
const plainAuditSignal = "[◆ audit proposals pending]"

// RenderPlainTable renders a plain text three-column task table for non-TTY output.
func RenderPlainTable(model Model) string {
	var b strings.Builder

	fmt.Fprintln(&b, plainNonTTYWarning)
	if hasAuditProposals(model.Root) {
		fmt.Fprintln(&b, plainAuditSignal)
	}
	fmt.Fprintln(&b)

	cols := []struct {
		label string
		col   data.ColumnType
	}{
		{"PLANNED", data.ColumnPlanned},
		{"IN PROGRESS", data.ColumnInProgress},
		{"DONE", data.ColumnDone},
	}

	for _, c := range cols {
		tasks := model.Tasks[c.col]
		fmt.Fprintln(&b, c.label)
		if len(tasks) == 0 {
			fmt.Fprintln(&b, "  (none)")
		}
		for _, t := range tasks {
			title := t.Title
			if title == "" {
				title = "(no title)"
			}
			fmt.Fprintf(&b, "  %-52s  %s\n", t.ID, title)
		}
		fmt.Fprintln(&b)
	}

	return b.String()
}

// hasAuditProposals reports whether any audit file under root contains a Proposed Changes section.
func hasAuditProposals(root string) bool {
	releasesDir := filepath.Join(root, "releases")
	releases, err := os.ReadDir(releasesDir)
	if err != nil {
		return false
	}
	for _, r := range releases {
		if !r.IsDir() {
			continue
		}
		epicsDir := filepath.Join(releasesDir, r.Name(), "epics")
		epics, err := os.ReadDir(epicsDir)
		if err != nil {
			continue
		}
		for _, e := range epics {
			if !e.IsDir() {
				continue
			}
			short := e.Name()
			if idx := strings.Index(short, "-"); idx >= 0 {
				short = short[:idx]
			}
			auditPath := filepath.Join(epicsDir, e.Name(), short+"-Audit.md")
			raw, err := os.ReadFile(auditPath)
			if err != nil {
				continue
			}
			if strings.Contains(string(raw), "## Proposed Changes") {
				return true
			}
		}
	}
	return false
}
