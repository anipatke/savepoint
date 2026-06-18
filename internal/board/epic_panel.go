package board

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/styles"
)

// RenderEpicDetail renders an overlay showing the content of an E##-Detail.md file.
func RenderEpicDetail(epicSlug, content string, overlayW, maxHeight, offset int, tab int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	tabIndicator := renderTabIndicator(tab, inner)
	lines := []string{
		styles.EpicTitleFocused.Render("EPIC DETAIL"),
		tabIndicator,
	}

	body := epicDetailBody(content, inner)
	body = append(body, "", styles.CardMeta.Render("1:Detail 2:Audit  esc:close"))
	lines = append(lines, visibleDetailLines(body, maxHeight-detailVerticalOverhead-1, offset)...)

	return styles.EpicDetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

func renderTabIndicator(tab int, width int) string {
	var detail, audit string
	if tab == 0 {
		detail = styles.EpicItemFocused.Render("DETAIL [1]")
		audit = styles.CardMeta.Render("AUDIT [2]")
	} else {
		detail = styles.CardMeta.Render("DETAIL [1]")
		audit = styles.EpicItemFocused.Render("AUDIT [2]")
	}
	return detail + styles.CardMeta.Render(" │ ") + audit
}

// stripFrontmatter removes YAML frontmatter (between leading --- markers) from content.
func stripFrontmatter(content string) []string {
	lines := strings.Split(content, "\n")
	start := 0
	if len(lines) > 0 && strings.TrimSpace(lines[0]) == "---" {
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				start = i + 1
				break
			}
		}
	}
	return lines[start:]
}

// epicDetailBody parses markdown content into display lines, stripping frontmatter.
func epicDetailBody(content string, width int) []string {
	if strings.TrimSpace(content) == "" || content == "(no detail available)" {
		return []string{styles.CardMeta.Render("(no detail available)")}
	}

	lines := stripFrontmatter(content)

	var body []string
	para := newParagraphFlusher(&body, width)
	skip := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			heading := strings.ToLower(strings.TrimPrefix(trimmed, "## "))
			skip = strings.Contains(heading, "component") || strings.Contains(heading, "files")
		}
		if skip {
			para.flush()
			continue
		}
		switch {
		case strings.HasPrefix(trimmed, "# "):
			para.flush()
			body = append(body, styles.EpicTitleFocused.Render(strings.TrimPrefix(trimmed, "# ")))
		case strings.HasPrefix(trimmed, "## "):
			para.flush()
			body = append(body, "", styles.EpicItemFocused.Render(strings.TrimPrefix(trimmed, "## ")))
		case strings.HasPrefix(trimmed, "### "):
			para.flush()
			body = append(body, styles.EpicItemFocused.Render(strings.TrimPrefix(trimmed, "### ")))
		case strings.HasPrefix(trimmed, "|"):
			para.flush()
			body = append(body, styles.CardMeta.Render(trimmed))
		case isListItem(trimmed):
			para.startItem(trimmed)
		case trimmed == "":
			para.flush()
			body = append(body, "")
		default:
			para.add(trimmed)
		}
	}
	para.flush()
	return body
}

// isListItem reports whether a trimmed line begins a Markdown list item.
func isListItem(trimmed string) bool {
	return strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ")
}

// isOrderedItem reports whether a trimmed line begins an ordered list item, e.g.
// "1. step" or "2) step".
func isOrderedItem(trimmed string) bool {
	i := 0
	for i < len(trimmed) && trimmed[i] >= '0' && trimmed[i] <= '9' {
		i++
	}
	if i == 0 || i+1 >= len(trimmed) {
		return false
	}
	return (trimmed[i] == '.' || trimmed[i] == ')') && trimmed[i+1] == ' '
}

// renderSectionBody reflows a plain multi-line Markdown section body into wrapped,
// styled display lines. Consecutive prose lines join into one paragraph; blank
// lines, list items (-, *, or "N."), and sub-headings break the paragraph so the
// source file's own hard wrapping neither orphans trailing words nor flattens
// list and paragraph structure. Used by the task and defect detail overlays,
// which render whole section bodies rather than headed epic content.
func renderSectionBody(content string, width int) []string {
	var body []string
	para := newParagraphFlusher(&body, width)
	for _, line := range strings.Split(strings.TrimSpace(content), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			para.flush()
			body = append(body, "")
		case strings.HasPrefix(trimmed, "### "):
			para.flush()
			body = append(body, styles.EpicItemFocused.Render(strings.TrimPrefix(trimmed, "### ")))
		case isListItem(trimmed) || isOrderedItem(trimmed):
			para.startItem(trimmed)
		default:
			para.add(trimmed)
		}
	}
	para.flush()
	return body
}

// paragraphFlusher reflows consecutive Markdown source lines into one paragraph
// before wrapping, so the source file's own hard line breaks do not leave ragged
// orphan words when re-wrapped at the panel width. Each list item is its own
// reflowed block; blank lines and block elements flush the buffer.
type paragraphFlusher struct {
	body  *[]string
	width int
	buf   []string
}

func newParagraphFlusher(body *[]string, width int) *paragraphFlusher {
	return &paragraphFlusher{body: body, width: width}
}

// add appends a continuation line to the current paragraph.
func (p *paragraphFlusher) add(line string) {
	p.buf = append(p.buf, line)
}

// startItem flushes the current paragraph and begins a new one for a list item,
// so adjacent items with no blank line between them stay separate.
func (p *paragraphFlusher) startItem(line string) {
	p.flush()
	p.buf = append(p.buf, line)
}

// flush wraps the buffered paragraph at the panel width and clears the buffer.
func (p *paragraphFlusher) flush() {
	if len(p.buf) == 0 {
		return
	}
	joined := strings.Join(p.buf, " ")
	for _, wrapped := range WrapText(joined, p.width) {
		*p.body = append(*p.body, styles.CardMeta.Render(wrapped))
	}
	p.buf = p.buf[:0]
}

// RenderEpicAuditTab renders an overlay showing audit findings from an E##-Audit.md file.
func RenderEpicAuditTab(epicSlug, content string, overlayW, maxHeight, offset int, tab int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	tabIndicator := renderTabIndicator(tab, inner)
	lines := []string{
		styles.GlyphAudit.Render("EPIC AUDIT"),
		tabIndicator,
	}

	body := epicAuditBody(content, inner)
	body = append(body, "", styles.CardMeta.Render("1:Detail 2:Audit  a:mark audited  esc:close"))
	lines = append(lines, visibleDetailLines(body, maxHeight-detailVerticalOverhead-1, offset)...)

	return styles.EpicDetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

// epicAuditHiddenSectionHeadings lists markdown section headings suppressed in the audit tab overlay.
// Sections that are implementation details or planning artifacts clutter the summary view.
var epicAuditHiddenSectionHeadings = map[string]struct{}{
	"12. Distribution & build": {},
	"Acceptance Criteria":      {},
	"Architectural notes":      {},
	"Boundaries":               {},
	"Context Files":            {},
	"Implemented As":           {},
	"Implemented as":           {},
	"Implementation Plan":      {},
	"Manual audit override":    {},
	"Proposed Changes":         {},
	"Quality Review":           {},
	"With":                     {},
}

func epicAuditBody(content string, width int) []string {
	if strings.TrimSpace(content) == "" || content == "(no audit available)" {
		return []string{styles.CardMeta.Render("(no audit available)")}
	}

	lines := stripFrontmatter(content)

	var body []string
	para := newParagraphFlusher(&body, width)
	inHiddenSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "## "):
			para.flush()
			sectionName := strings.TrimPrefix(trimmed, "## ")
			_, inHiddenSection = epicAuditHiddenSectionHeadings[sectionName]
			if !inHiddenSection {
				body = append(body, "", styles.EpicItemFocused.Render(sectionName))
			}
		case inHiddenSection:
			para.flush()
		case strings.HasPrefix(trimmed, "### "):
			para.flush()
			body = append(body, styles.EpicItemFocused.Render(strings.TrimPrefix(trimmed, "### ")))
		case strings.HasPrefix(trimmed, "- [x] ") || strings.HasPrefix(trimmed, "- [X] "):
			para.flush()
			text := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- [x] "), "- [X] ")
			body = append(body, renderChecklistSentences(text, "[x] ", width, styles.TagDone)...)
		case strings.HasPrefix(trimmed, "- [ ] "):
			para.flush()
			text := strings.TrimPrefix(trimmed, "- [ ] ")
			body = append(body, renderChecklistSentences(text, "[ ] ", width, styles.CardMeta)...)
		case strings.HasPrefix(trimmed, "- "):
			para.flush()
			body = append(body, styles.CardMeta.Render("• "+strings.TrimPrefix(trimmed, "- ")))
		case trimmed == "":
			para.flush()
			body = append(body, "")
		default:
			para.add(trimmed)
		}
	}
	para.flush()
	return body
}

const epicActiveMarker = "►"

// RenderEpicSidebar renders the fixed left sidebar listing epics with active indicator.
// If epics is empty and selected is non-empty, selected is shown as the sole entry.
func RenderEpicSidebar(epics []string, selected string, width int, focus bool, cursor int, status map[string]string, maxHeight int) string {
	inner := width - epicPanelOverhead
	if inner < 2 {
		inner = 2
	}
	list := epics
	if len(list) == 0 && selected != "" {
		list = []string{selected}
	}

	title := styles.ColumnTitle.Render("EPICS")
	if focus {
		title = styles.EpicTitleFocused.Render("EPICS")
	}
	lines := []string{title, strings.Repeat("─", inner)}
	for i, e := range list {
		g := epicSidebarGlyph(status, e)
		gw := lipgloss.Width(g)
		if gw < 1 {
			gw = 1
		}
		label := truncate(e, inner-2-gw)
		if focus && len(epics) > 0 && i == cursor {
			lines = append(lines, styles.EpicItemFocused.Render(epicActiveMarker+" "+g+" "+label))
		} else if !focus && e == selected {
			lines = append(lines, styles.EpicItemFocused.Render(epicActiveMarker+" "+g+" "+label))
		} else {
			lines = append(lines, styles.TaskItem.Render("  "+g+" "+label))
		}
	}
	if len(list) == 0 {
		lines = append(lines, styles.TaskItem.Render("(none)"))
	}
	if maxHeight > 0 && len(lines) > maxHeight {
		items := lines[2:]
		available := maxHeight - 3
		if available < 1 {
			available = 1
		}
		clipped := make([]string, 0, maxHeight)
		clipped = append(clipped, lines[0], lines[1])
		clipped = append(clipped, items[:min(available, len(items))]...)
		if len(items) > available {
			clipped = append(clipped, renderScrollIndicator("↓", len(items)-available, "more"))
		}
		lines = clipped
	}
	style := styles.EpicPanel.Width(width)
	if focus && len(epics) > 0 {
		style = styles.EpicPanelFocused.Width(width)
	}
	return style.Render(strings.Join(lines, "\n"))
}

func epicSidebarGlyph(status map[string]string, epicID string) string {
	if status == nil {
		return statusGlyphDefault
	}
	s, ok := status[epicID]
	if !ok {
		return statusGlyphDefault
	}
	return statusGlyph(s)
}

// RenderEpicDropdown renders the epic selection dropdown overlay.
func RenderEpicDropdown(epics []string, cursor int, width int) string {
	inner := width - epicPanelOverhead
	if inner < 2 {
		inner = 2
	}

	lines := []string{styles.ColumnTitleFocused.Render("SELECT EPIC"), strings.Repeat("─", inner)}
	for i, e := range epics {
		label := truncate(e, inner-2)
		if i == cursor {
			lines = append(lines, styles.TaskItemFocused.Render(epicActiveMarker+" "+label))
		} else {
			lines = append(lines, styles.TaskItem.Render("  "+label))
		}
	}
	if len(epics) == 0 {
		lines = append(lines, styles.TaskItem.Render("(none)"))
	}
	lines = append(lines, "", styles.CardMeta.Render("↑↓:nav  enter:select  esc:cancel"))
	return styles.EpicPanel.Width(width).Render(strings.Join(lines, "\n"))
}
