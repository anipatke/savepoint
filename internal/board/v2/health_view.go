package v2

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/codehealth"
	"github.com/opencode/savepoint/internal/styles"
)

// Every phrase the Code Health screen adds to what codehealth already worded
// lives here (STYLE-09).
const (
	textHealthTitle         = "CODE HEALTH"
	textHealthLoading       = "Reading saved results…"
	textHealthLoadFailed    = "Code Health history could not be read, so the screen may be out of date: %v"
	textHealthRefreshFailed = "Refresh failed; nothing was saved: %v"
	textHealthCancelled     = "Refresh cancelled; nothing was saved."
	textHealthCancelling    = "Cancelling…"
	textHealthRefreshing    = "Refreshing %d of %d: %s (%s)"
	textHealthStarting      = "Refreshing: starting…"
	textHealthEsc           = "Esc to cancel"
	textHealthChecking      = "Checking whether your code has changed since…"
	textHealthMeasured      = "Measured %s · %s"
	textHealthHistory       = "RECENT CHECKS"
	textHealthNoHistory     = "No earlier checks."
	textHealthDetailTitle   = "SIGNAL DETAIL"
)

var healthNotConfiguredLines = []string{
	"Code Health is not set up for this project.",
	"",
	"It answers one question: is your code staying healthy? It reads your",
	"test, coverage, complexity, duplication, and dependency reports and says",
	"in plain words what is good, what to watch, and what needs attention.",
	"",
	"Run `savepoint health setup` to turn it on.",
}

var healthFirstRunLines = []string{
	"Nothing has been measured yet.",
	"",
	"Press R to measure your code now. Nothing is changed in your project.",
}

// healthGlyph reinforces the label word so meaning never rests on color.
var healthGlyph = map[codehealth.Classification]string{
	codehealth.ClassificationGood:           "✓",
	codehealth.ClassificationWatch:          "◆",
	codehealth.ClassificationNeedsAttention: "✗",
	codehealth.ClassificationUnknown:        "?",
}

func healthLabelStyle(c codehealth.Classification) lipgloss.Style {
	switch c {
	case codehealth.ClassificationGood:
		return styles.HealthGood
	case codehealth.ClassificationWatch:
		return styles.HealthWatch
	case codehealth.ClassificationNeedsAttention:
		return styles.HealthNeedsAttention
	}
	return styles.HealthUnknown
}

func healthLabel(c codehealth.Classification, text string) string {
	return healthLabelStyle(c).Render(healthGlyph[c] + " " + text)
}

// healthPinned is the heading and the live status lines: they stay on screen
// while the body scrolls.
func (m Model) healthPinned(width int) []string {
	h := m.Health
	lines := []string{
		styles.HealthAccent.Render(textHealthTitle),
		styles.Divider.Render(strings.Repeat("─", width)),
	}
	switch {
	case h.Refresh != nil && h.Refresh.Cancelling:
		lines = append(lines, styles.HealthWatch.Render(fitLine(textHealthCancelling, width)))
	case h.Refresh != nil:
		lines = append(lines, styles.HealthWatch.Render(fitLine(refreshLine(h.Refresh), width)), styles.CardMeta.Render(textHealthEsc))
	case h.Notice != "":
		lines = append(lines, wrapDetailLine(h.Notice, width)...)
	}
	return lines
}

func refreshLine(r *healthRefresh) string {
	if !r.Started {
		return textHealthStarting
	}
	p := r.Progress
	signal := codehealth.CapabilityText(p.Capability)
	if p.Name != "" {
		signal += " " + p.Name
	}
	return fmt.Sprintf(textHealthRefreshing, p.Position, p.Total, signal, p.Provider)
}

// healthBody is the scrolling content and the line the cursor is on, or
// noCursorLine when the reader has scrolled freely.
func (m Model) healthBody(width int) (lines []string, cursorLine int) {
	h := m.Health
	if !h.Loaded {
		if h.Notice != "" {
			return nil, 0
		}
		return []string{styles.CardMeta.Render(textHealthLoading)}, 0
	}
	d := h.Dashboard
	switch d.State {
	case codehealth.DashboardNotConfigured:
		return plainLines(healthNotConfiguredLines, width), 0
	case codehealth.DashboardFirstRun:
		return plainLines(healthFirstRunLines, width), 0
	}
	if h.Detail && h.Cursor < len(d.Rows) {
		// A detail has no cursor, so the offset is the reader's own scroll.
		return healthDetailLines(*d, d.Rows[h.Cursor], width), noCursorLine
	}
	lines, cursorLine = healthOverviewLines(h, width)
	if h.Scrolled {
		cursorLine = noCursorLine
	}
	return lines, cursorLine
}

// noCursorLine tells healthWindow not to pull the window back to a cursor.
const noCursorLine = -1

func plainLines(text []string, width int) []string {
	var out []string
	for _, line := range text {
		out = append(out, wrapDetailLine(line, width)...)
	}
	return out
}

func healthOverviewLines(h *HealthOverlay, width int) (lines []string, cursorLine int) {
	d := h.Dashboard
	lines = append(lines, healthLabel(d.Overall, d.OverallText)+styles.CardMeta.Render("  overall"))
	lines = append(lines, wrapDetailLine(fmt.Sprintf(textHealthMeasured, d.MeasuredAt, d.OriginText), width)...)
	switch {
	case h.Freshness != nil:
		lines = append(lines, wrapDetailLine(h.Freshness.Text, width)...)
	default:
		lines = append(lines, styles.CardMeta.Render(fitLine(textHealthChecking, width)))
	}
	lines = append(lines, "")
	for i, row := range d.Rows {
		if i == h.Cursor {
			cursorLine = len(lines)
		}
		lines = append(lines, healthRowLines(row, i == h.Cursor, width)...)
	}
	lines = append(lines, "", styles.HealthAccent.Render(textHealthHistory))
	if len(d.History) == 0 {
		lines = append(lines, styles.CardMeta.Render(textHealthNoHistory))
	}
	for _, e := range d.History {
		lines = append(lines, fitLine(fmt.Sprintf("  %s  %s  %s", e.CreatedAt, e.OverallText, e.OriginText), width))
	}
	return lines, cursorLine
}

// healthRowLines is a signal's label, name, and trend on one line, and its
// reason beneath. A row that is not Good is never drawn like one.
func healthRowLines(row codehealth.DashboardRow, selected bool, width int) []string {
	marker := "  "
	if selected {
		marker = "▸ "
	}
	name := row.CapabilityText
	if row.Name != "" {
		name += " · " + row.Name
	}
	head := fmt.Sprintf("%s%s  %s", marker, healthLabel(row.Label, row.LabelText), name)
	if row.Trend != "" {
		head += styles.CardMeta.Render("  " + row.Trend)
	}
	lines := []string{fitLine(head, width)}
	for _, line := range wrapDetailLine("    "+row.Explanation, width) {
		lines = append(lines, styles.CardMeta.Render(line))
	}
	return lines
}

func healthDetailLines(d codehealth.Dashboard, row codehealth.DashboardRow, width int) []string {
	name := row.CapabilityText
	if row.Name != "" {
		name += " · " + row.Name
	}
	lines := []string{styles.HealthAccent.Render(textHealthDetailTitle), name, healthLabel(row.Label, row.LabelText), ""}
	field := func(label, value string) {
		if value != "" {
			lines = append(lines, wrapDetailLine(label+": "+value, width)...)
		}
	}
	lines = append(lines, wrapDetailLine(row.Explanation, width)...)
	lines = append(lines, "")
	field("Trend", row.Trend)
	field("Compared with", row.Basis)
	field("Outcome", row.OutcomeText)
	field("Required", row.RequiredText)
	field("Freshness", row.FreshnessText)
	field("Provider", strings.TrimSpace(string(row.Provider)+" "+row.ProviderVersion))
	field("Scope", strings.Join(row.Scope, ", "))
	field("Measured", row.CollectedAt)
	field("Snapshot", d.OriginText+" "+d.SnapshotID)
	lines = append(lines, "", styles.HealthAccent.Render("AFFECTED AREAS"))
	if len(row.Evidence) == 0 {
		lines = append(lines, styles.CardMeta.Render("  None recorded."))
	}
	for _, ref := range row.Evidence {
		lines = append(lines, wrapDetailLine("  "+evidenceText(ref), width)...)
	}
	return lines
}

func evidenceText(ref codehealth.EvidenceRef) string {
	out := ref.Path
	if ref.Line > 0 {
		out += fmt.Sprintf(":%d", ref.Line)
	}
	if ref.Note != "" {
		out += " — " + ref.Note
	}
	return out
}

// healthWindow is the first body line shown: the stored offset, moved just far
// enough to keep the cursor in view (when there is one) and clamped inside the
// content.
func healthWindow(total, cursorLine, offset, height int) int {
	height = max(height, 1)
	if cursorLine >= 0 {
		if cursorLine < offset {
			offset = cursorLine
		}
		if cursorLine >= offset+height {
			offset = cursorLine - height + 1
		}
	}
	return min(max(offset, 0), max(total-height, 0))
}

// healthBodyHeight is what is left for the body under the pinned lines.
func (m Model) healthBodyHeight(width, height int) int {
	return max(columnBodyHeight(height)-len(m.healthPinned(width)), 1)
}

// healthGeometry is the body's line count, the lines that fit, and the line
// the selected signal starts on, at the current size. The signal's line is
// reported even while the reader has scrolled past it.
func (m Model) healthGeometry() (total, room, cursorLine int) {
	w, height := m.detailViewport()
	width := columnTextWidth(w)
	lines, cursorLine := m.healthBody(width)
	if cursorLine == noCursorLine && !m.Health.Detail {
		_, cursorLine = healthOverviewLines(m.Health, width)
	}
	return len(lines), m.healthBodyHeight(width, height), cursorLine
}

func (m *Model) syncHealthScroll() {
	if m.Health == nil {
		return
	}
	w, height := m.detailViewport()
	width := columnTextWidth(w)
	lines, cursorLine := m.healthBody(width)
	m.Health.Offset = healthWindow(len(lines), cursorLine, m.Health.Offset, m.healthBodyHeight(width, height))
}

func renderHealth(m Model, w, height int) string {
	width := columnTextWidth(w)
	pinned := m.healthPinned(width)
	body, cursorLine := m.healthBody(width)
	room := m.healthBodyHeight(width, height)
	start := healthWindow(len(body), cursorLine, m.Health.Offset, room)
	end := min(start+room, len(body))

	lines := append(pinned, body[start:end]...)
	for len(lines) < len(pinned)+room {
		lines = append(lines, "")
	}
	return frameColumn(lines, width, columnBodyHeight(height), true)
}
