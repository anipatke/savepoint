package v2

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/codehealth"
	"github.com/opencode/savepoint/internal/styles"
)

// The popover has a fixed size whatever the data: a title, a status line, five
// signal rows, a rule, up to five lines about the selected signal, the
// sign-off line and the re-run command.
const (
	healthPopoverMaxWidth = 72
	healthPopoverInner    = 15
	healthPopoverHeight   = healthPopoverInner + 2
	healthRows            = 5
	healthRowMarkWidth    = 4
	healthRowGap          = 2
	healthRowNameCap      = 16
	healthRowValueMin     = 8
	healthDetailLines     = 5
)

// Every phrase the Code Health popover adds to what codehealth already worded
// lives here (STYLE-09).
const (
	textHealthTitle          = "CODE HEALTH"
	textHealthLoading        = "Reading saved results…"
	textHealthLoadFailed     = "Could not read history, showing older results: %v"
	textHealthRefreshFailed  = "Refresh failed; nothing was saved: %v"
	textHealthCancelled      = "Refresh cancelled; nothing was saved."
	textHealthCancelling     = "Cancelling…"
	textHealthRefreshing     = "Refreshing %d of %d: %s"
	textHealthStarting       = "Refreshing: starting…"
	textHealthEsc            = "Esc to cancel"
	textHealthChecking       = "Checking whether your code has changed since…"
	textHealthMeasured       = "%s · %s · %s"
	textHealthRerun          = "Re-run: savepoint health check %s"
	textHealthOlderMark      = "*"
	textHealthOlderNote      = " (measured on other code)"
	textHealthHistoryTitle   = "  last checks, newest first"
	textHealthNoHistory      = "No checks yet."
	textHealthDetailJoin     = " · "
	textHealthValueJoin      = " · "
	textHealthBlocking       = "!"
	textHealthInstances      = "%s ×%d"
	textHealthWorstOf        = " (worst of %d: %s)"
	textHealthWorstOfUnnamed = " (worst of %d)"
	textHealthNoTrend        = "no trend yet"
	textHealthHistoryNote    = "History: %s"
	textHealthManualNote     = "Dimmed rows are manual refreshes; they do not feed the sparklines."
)

// healthRowLabel renames a signal where the number beside it counts something
// other than the signal's own name.
var healthRowLabel = map[codehealth.Capability]string{
	codehealth.CapabilityTests: "Tests failing",
}

// healthSignalName is the row's name: the signal's own, or its row label with
// any "×N" the grouping added.
func healthSignalName(row codehealth.DashboardRow) string {
	label, ok := healthRowLabel[row.Capability]
	if !ok {
		return row.CapabilityText
	}
	return label + strings.TrimPrefix(row.CapabilityText, codehealth.CapabilityText(row.Capability))
}

// healthShortName stands in for a signal name too long for a row.
var healthShortName = map[codehealth.Capability]string{
	codehealth.CapabilityDependencyVulnerability: "Dependencies",
}

var healthNotConfiguredLines = []string{
	"Code Health is not set up. It shows if your code stays healthy.",
	"Run `savepoint health setup` to turn it on.",
}

var healthFirstRunLines = []string{
	"Nothing has been measured yet.",
	"Press R to measure your code now; nothing in your project changes.",
}

// healthGlyph reinforces the label word so meaning never rests on color.
var healthGlyph = map[codehealth.Classification]string{
	codehealth.ClassificationGood:           "✓",
	codehealth.ClassificationWatch:          "~",
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

func refreshLine(r *healthRefresh) string {
	if !r.Started {
		return textHealthStarting
	}
	p := r.Progress
	signal := codehealth.CapabilityText(p.Capability)
	if p.Name != "" {
		signal += " " + p.Name
	}
	return fmt.Sprintf(textHealthRefreshing, p.Position, p.Total, signal)
}

// healthPopoverWidth is the popover's outer width: 72 columns, or what the
// terminal leaves.
func healthPopoverWidth(termWidth int) int {
	return max(min(healthPopoverMaxWidth, termWidth-4), 8)
}

func (m Model) renderHealthOverlay(base string, width, height int) string {
	return overlayOnV2Base(dimV2Lines(base), m.renderHealthPopover(width), width, height)
}

func (m Model) renderHealthPopover(termWidth int) string {
	outer := healthPopoverWidth(termWidth)
	inner := outer - 4
	lines := m.healthPopoverLines(inner)
	for i, line := range lines {
		lines[i] = xansi.Truncate(line, inner, "…")
	}
	for len(lines) < healthPopoverInner {
		lines = append(lines, "")
	}
	lines = lines[:healthPopoverInner]
	return styles.DetailOverlay.Width(outer - 2).Render(strings.Join(lines, "\n"))
}

func (m Model) healthPopoverLines(width int) []string {
	h := m.Health
	title := styles.HealthAccent.Render(textHealthTitle)
	if !h.Loaded {
		if h.Notice != "" {
			return append([]string{title, ""}, plainLines([]string{h.Notice}, width)...)
		}
		return []string{title, "", styles.CardMeta.Render(textHealthLoading)}
	}
	d := h.Dashboard
	if h.History {
		return m.healthHistoryLines(width)
	}
	switch d.State {
	case codehealth.DashboardNotConfigured:
		return append([]string{title, ""}, plainLines(healthNotConfiguredLines, width)...)
	case codehealth.DashboardFirstRun:
		lines := append([]string{title, ""}, plainLines(healthFirstRunLines, width)...)
		lines = append(lines, "", m.healthStatusLine(width))
		for len(lines) < healthPopoverInner-1 {
			lines = append(lines, "")
		}
		return append(lines[:healthPopoverInner-1], m.healthFooterLine())
	}

	lines := []string{
		title + styles.CardMeta.Render("  "+fmt.Sprintf(textHealthMeasured, d.MeasuredText, d.OriginText, d.Headline)),
		m.healthStatusLine(width),
	}
	lines = append(lines, healthSignalRows(h, width)...)
	lines = append(lines, styles.Divider.Render(strings.Repeat("─", width)))
	lines = append(lines, healthSelectedLines(h)...)
	for len(lines) < healthPopoverInner-2 {
		lines = append(lines, "")
	}
	lines = append(lines[:healthPopoverInner-2], styles.CardMeta.Render(d.SignOff), m.healthFooterLine())
	return lines
}

// healthHistoryLines is the last ten checks in the same frame as the signals:
// the footer stays on the last line and the note sits just above it.
func (m Model) healthHistoryLines(width int) []string {
	h := m.Health
	lines := []string{styles.HealthAccent.Render(textHealthTitle) + styles.CardMeta.Render(textHealthHistoryTitle), m.healthStatusLine(width)}
	entries := h.Dashboard.History
	if len(entries) == 0 {
		lines = append(lines, styles.CardMeta.Render(textHealthNoHistory))
	}
	var whenW, originW int
	for _, e := range entries[:min(len(entries), codehealth.MaxDashboardHistory)] {
		whenW = max(whenW, lipgloss.Width(e.WhenText))
		originW = max(originW, lipgloss.Width(e.OriginText))
	}
	for i, e := range entries {
		if i == codehealth.MaxDashboardHistory {
			break
		}
		head := padCells(e.WhenText, whenW) + "  " + padCells(e.OriginText, originW) + "  "
		label := healthGlyph[e.Overall] + " " + e.OverallText
		if e.Origin == codehealth.OriginManual {
			lines = append(lines, fitLine(styles.CardMeta.Render(head+label), width))
			continue
		}
		lines = append(lines, fitLine(head+healthLabelStyle(e.Overall).Render(label), width))
	}
	for len(lines) < healthPopoverInner-2 {
		lines = append(lines, "")
	}
	return append(lines[:healthPopoverInner-2], styles.CardMeta.Render(textHealthManualNote), m.healthFooterLine())
}

// healthStatusLine is the one live line under the title: progress, a notice,
// or how current the measured code is.
func (m Model) healthStatusLine(width int) string {
	h := m.Health
	switch {
	case h.Refresh != nil && h.Refresh.Cancelling:
		return styles.HealthWatch.Render(textHealthCancelling)
	case h.Refresh != nil:
		return styles.HealthWatch.Render(fitLine(refreshLine(h.Refresh), width-len(textHealthEsc)-2)) + styles.CardMeta.Render("  "+textHealthEsc)
	case h.Notice != "":
		return styles.HealthWatch.Render(h.Notice)
	case h.Freshness != nil:
		return styles.CardMeta.Render(h.Freshness.Text)
	case h.Dashboard != nil && h.Dashboard.State == codehealth.DashboardMeasured:
		return styles.CardMeta.Render(textHealthChecking)
	}
	return ""
}

// healthFooterLine is the command to re-run the official check, or the manual
// refresh note while a refresh runs.
func (m Model) healthFooterLine() string {
	if m.Health.Refresh != nil {
		return ""
	}
	return styles.CardMeta.Render(m.Health.Rerun)
}

// healthSignalRows is one line per signal. The scale (the aim) and the
// direction are what a glance needs, so they keep their full width; the name
// and then the value give way when the line is tight.
func healthSignalRows(h *HealthOverlay, width int) []string {
	rows := h.rows()
	shown := rows[:min(len(rows), healthRows)]
	older := healthOlderMark(h)
	names := make([]string, len(shown))
	values := make([]string, len(shown))
	sparks := make([]string, len(shown))
	var nameW, valueW, sparkW, aimW int
	for i, row := range shown {
		names[i] = healthRowName(row)
		values[i] = healthFigure(row) + older
		sparks[i] = healthSpark(row, healthLabelStyle(row.Label))
		nameW = max(nameW, lipgloss.Width(names[i]))
		valueW = max(valueW, lipgloss.Width(values[i]))
		sparkW = max(sparkW, lipgloss.Width(sparks[i]))
		aimW = max(aimW, lipgloss.Width(row.Aim))
	}
	room := width - healthRowMarkWidth - sparkW - aimW - 3*healthRowGap
	if nameW+valueW > room {
		nameW = min(nameW, healthRowNameCap)
		valueW = max(min(valueW, room-nameW), healthRowValueMin)
	}
	gap := strings.Repeat(" ", healthRowGap)
	out := make([]string, 0, len(shown))
	for i, row := range shown {
		marker := "  "
		if i == h.Cursor {
			marker = "▸ "
		}
		style := healthLabelStyle(row.Label)
		line := marker + style.Render(healthGlyph[row.Label]) + " " + padCells(cutCells(names[i], nameW), nameW) +
			gap + padCells(cutValue(values[i], valueW), valueW) + gap + padCells(sparks[i], sparkW)
		if row.Aim != "" {
			line += styles.CardMeta.Render(gap + row.Aim)
		}
		if i == h.Cursor {
			line = healthSelectedStyle.Render(line)
		}
		out = append(out, fitLine(line, width))
	}
	return out
}

// healthFigure is the value column: the number alone, with ! when a
// vulnerability finding blocks sign-off. Rows without a figure show their value.
func healthFigure(row codehealth.DashboardRow) string {
	figure := row.Figure
	if figure == "" {
		figure = row.Value
	}
	if row.Capability == codehealth.CapabilityDependencyVulnerability && row.BlocksSignOff() {
		figure += textHealthBlocking
	}
	return figure
}

// healthRowName is the signal's name, with the instance name when one signal
// has several.
func healthRowName(row codehealth.DashboardRow) string {
	name := healthSignalName(row)
	if short, ok := healthShortName[row.Capability]; ok && row.Name == "" && row.CapabilityText == codehealth.CapabilityText(row.Capability) && lipgloss.Width(name) > healthRowNameCap {
		name = short
	}
	if row.Name != "" {
		name += " " + row.Name
	}
	return name
}

// cutCells shortens text to width cells, ending in an ellipsis.
func cutCells(text string, width int) string {
	return xansi.Truncate(text, width, "…")
}

// cutValue shortens a value to width cells. A value that ends in its number
// ("hardest function scores 46") keeps the number and drops words before it.
func cutValue(text string, width int) string {
	if lipgloss.Width(text) <= width {
		return text
	}
	words := strings.Fields(text)
	last := words[len(words)-1]
	if len(words) > 2 && strings.ContainsAny(last, "0123456789") {
		for keep := len(words) - 2; keep > 0; keep-- {
			short := strings.Join(words[:keep], " ") + "… " + last
			if lipgloss.Width(short) <= width {
				return short
			}
		}
	}
	return cutCells(text, width)
}

// healthOlderMark flags values measured on code other than the code now: moved
// on, or from a different branch.
func healthOlderMark(h *HealthOverlay) string {
	if h.Freshness != nil && (h.Freshness.State == codehealth.CodeMovedOn || h.Freshness.State == codehealth.CodeOtherBranch) {
		return textHealthOlderMark
	}
	return ""
}

// healthSpark is the sparkline with its last block coloured by the signal's
// label, then the word for its direction.
func healthSpark(row codehealth.DashboardRow, style lipgloss.Style) string {
	if row.Spark == "" {
		return styles.CardMeta.Render(textHealthNoTrend)
	}
	blocks := []rune(row.Spark)
	last := len(blocks) - 1
	out := string(blocks[:last]) + style.Render(string(blocks[last]))
	if row.SparkWord != "" {
		out += " " + row.SparkWord
	}
	return out
}

var healthSelectedStyle = lipgloss.NewStyle().Bold(true)

func padCells(text string, width int) string {
	if gap := width - lipgloss.Width(text); gap > 0 {
		return text + strings.Repeat(" ", gap)
	}
	return text
}

// healthSelectedLines explains the selected signal in at most five lines.
func healthSelectedLines(h *HealthOverlay) []string {
	rows := h.rows()
	if h.Cursor >= len(rows) {
		return nil
	}
	row := rows[h.Cursor]
	question := row.Question
	if healthOlderMark(h) != "" {
		question += textHealthOlderNote
	}
	if row.Value != "" {
		question += textHealthValueJoin + row.Value
	}
	lines := []string{styles.HealthAccent.Render(question), row.Meaning}
	if row.SparkNote != "" {
		lines = append(lines, fmt.Sprintf(textHealthHistoryNote, row.SparkNote))
	}
	signOff, where := "", ""
	if row.SignOff != "" {
		signOff = "Sign-off: " + row.SignOff
	}
	if row.Where != "" {
		where = "Where: " + row.Where
	}
	next := "Next: " + row.NextStep
	// The five lines always end with Next; sign-off and where share a line
	// when the history note needs the room.
	if len(lines)+len(nonEmpty(signOff, where))+1 > healthDetailLines {
		return append(lines, strings.Join(nonEmpty(signOff, where), textHealthDetailJoin), next)
	}
	return append(append(lines, nonEmpty(signOff, where)...), next)
}

func nonEmpty(texts ...string) []string {
	var out []string
	for _, text := range texts {
		if text != "" {
			out = append(out, text)
		}
	}
	return out
}

func plainLines(text []string, width int) []string {
	var out []string
	for _, line := range text {
		out = append(out, wrapDetailLine(line, width)...)
	}
	return out
}
