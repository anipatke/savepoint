package v2

import (
	"maps"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

const (
	goalLabel           = "Goal"
	goalsLabel          = "Goals"
	goalSelectorKey     = "g"
	goalSelectorAlias   = "r"
	goalToggleKey       = "C"
	releaseActiveMarker = "►"
)

// orderedReleaseIDs is the selector's one source of ordering. The index is a
// map because identity lookup is the data contract; the board presents it in
// stable R-### order for repeatable navigation.
func orderedReleaseIDs(index *data.V2Index) []string {
	if index == nil {
		return nil
	}
	return slices.Sorted(maps.Keys(index.Releases))
}

func releaseIndex(releases []string, selected string) int {
	for i, release := range releases {
		if release == selected {
			return i
		}
	}
	if len(releases) == 0 {
		return 0
	}
	return 0
}

func releaseExists(index *data.V2Index, id string) bool {
	if index == nil || id == "" {
		return false
	}
	_, ok := index.Releases[id]
	return ok
}

func selectedRelease(state ProjectState) string {
	if state.Index == nil || state.Router == nil || state.Router.Release == "" {
		return ""
	}
	if diagnostic := state.Next.SelectionDiagnostic; diagnostic != nil {
		switch diagnostic.Kind {
		case data.SelectionReleaseMissing, data.SelectionReleaseNotFound, data.SelectionReleaseArchived:
			return ""
		}
	}
	if releaseExists(state.Index, state.Router.Release) {
		return state.Router.Release
	}
	return ""
}

func releaseLabel(index *data.V2Index, id string) string {
	if index == nil {
		return id
	}
	release, ok := index.Releases[id]
	if !ok || strings.TrimSpace(release.Title) == "" {
		return id
	}
	return id + " — " + release.Title
}

func renderReleaseSelector(index *data.V2Index, releases []string, cursor, width, height int) string {
	width = releaseOverlayWidth(width)
	inner := width - 4
	if inner < 2 {
		inner = 2
	}

	lines := []string{
		styles.ColumnTitleFocused.Render("SELECT " + strings.ToUpper(goalLabel)),
		strings.Repeat("─", inner),
	}

	if len(releases) == 0 {
		lines = append(lines, styles.TaskItem.Render("(no "+goalsLabel+" in this project)"))
	} else {
		start, end := releaseWindow(releases, cursor, height)
		if start > 0 {
			lines = append(lines, styles.CardMeta.Render("↑ more"))
		}
		for i := start; i < end; i++ {
			marker := releaseStatusMarker(index, releases[i])
			label := xansi.Truncate(releaseLabel(index, releases[i]), inner-2-lipgloss.Width(glyphCheckPending)-1, "…")
			if i == cursor {
				lines = append(lines, styles.TaskItemFocused.Render(releaseActiveMarker+" ")+marker+" "+styles.TaskItemFocused.Render(label))
			} else {
				lines = append(lines, styles.TaskItem.Render("  ")+marker+" "+styles.TaskItem.Render(label))
			}
		}
		if end < len(releases) {
			lines = append(lines, styles.CardMeta.Render("↓ more"))
		}
	}

	lines = append(lines, "",
		styles.CardMeta.Render("enter:select  "+goalToggleKey+":close/reopen  v:detail  esc:cancel"))
	return styles.DetailOverlay.Width(width).Render(strings.Join(lines, "\n"))
}

func releaseWindow(releases []string, cursor, height int) (int, int) {
	if len(releases) == 0 {
		return 0, 0
	}
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(releases) {
		cursor = len(releases) - 1
	}
	rows := len(releases)
	if height > 0 {
		rows = height - 8
		if rows < 1 {
			rows = 1
		}
		if rows > len(releases) {
			rows = len(releases)
		}
	}
	start := 0
	if cursor >= rows {
		start = cursor - rows + 1
	}
	return start, start + rows
}

func releaseOverlayWidth(termWidth int) int {
	if termWidth <= 0 {
		return 1
	}
	width := termWidth - 4
	if width > 56 {
		width = 56
	}
	if width < 8 {
		width = termWidth
	}
	if width < 1 {
		return 1
	}
	return width
}

func (m Model) renderReleaseOverlay(base string, width, height int) string {
	selector := renderReleaseSelector(m.State.Index, m.Releases, m.ReleaseCursor, releaseOverlayWidth(width), height)
	return overlayOnV2Base(dimV2Lines(base), selector, width, height)
}

// dimV2Lines keeps the loaded board visible behind the selector while making
// the focused list the active surface. It is local to V2 so the package never
// reaches into the V1 board's rendering helpers.
func dimV2Lines(value string) string {
	dim := lipgloss.NewStyle().Faint(true)
	lines := strings.Split(value, "\n")
	for i, line := range lines {
		lines[i] = dim.Render(line)
	}
	return strings.Join(lines, "\n")
}

func overlayOnV2Base(base, overlay string, termWidth, termHeight int) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")
	overlayHeight := len(overlayLines)
	overlayWidth := 0
	for _, line := range overlayLines {
		if lineWidth := lipgloss.Width(line); lineWidth > overlayWidth {
			overlayWidth = lineWidth
		}
	}

	startY := (termHeight - overlayHeight) / 2
	if startY < 0 {
		startY = 0
	}
	startX := (termWidth - overlayWidth) / 2
	if startX < 0 {
		startX = 0
	}

	for len(baseLines) < termHeight {
		baseLines = append(baseLines, "")
	}
	for i, line := range baseLines {
		overlayIndex := i - startY
		if overlayIndex < 0 || overlayIndex >= overlayHeight {
			continue
		}
		left := xansi.Truncate(line, startX, "")
		leftWidth := lipgloss.Width(left)
		if leftWidth < startX {
			left += strings.Repeat(" ", startX-leftWidth)
		}
		baseLines[i] = left + overlayLines[overlayIndex]
	}
	return strings.Join(baseLines, "\n")
}
