package board

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
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

// renderTabIndicator renders the Epic detail subview tab strip, highlighting the
// active tab (0=Detail, 1=Audit).
func renderTabIndicator(tab int, width int) string {
	sep := styles.CardMeta.Render(" │ ")
	return tabLabel("DETAIL [1]", tab == 0) + sep +
		tabLabel("AUDIT [2]", tab == 1)
}

// tabLabel styles a single tab label, emphasized when it is the active tab.
func tabLabel(text string, active bool) string {
	if active {
		return styles.EpicItemFocused.Render(text)
	}
	return styles.CardMeta.Render(text)
}

// RenderReleaseDocs renders the top-level Release Docs overlay: a document
// selector (Release PRD, Overall PRD, Overall Design) above a scrollable,
// read-only body for the selected document. Missing or empty documents render an
// inline empty state rather than an error. offset scrolls the selected
// document's body. The selector is one fixed header row beyond the title, so its
// row is reserved from the body viewport like the Epic detail tab strip.
func RenderReleaseDocs(docs []data.ReleaseDoc, selected, overlayW, maxHeight, offset int) string {
	inner := overlayW - detailBorderPad
	if inner < 4 {
		inner = 4
	}

	lines := []string{
		styles.EpicTitleFocused.Render("RELEASE DOCS"),
		renderReleaseDocSelector(docs, selected),
	}

	body := releaseDocBody(docs, selected, inner)
	body = append(body, "", styles.CardMeta.Render("[/]:doc  esc:close"))
	lines = append(lines, visibleDetailLines(body, maxHeight-detailVerticalOverhead-1, offset)...)

	return styles.EpicDetailOverlay.Width(overlayW).Render(strings.Join(lines, "\n"))
}

// renderReleaseDocSelector renders the PRD/Design document selector, emphasizing
// the selected entry using the same active/inactive styling as the tab strip.
func renderReleaseDocSelector(docs []data.ReleaseDoc, selected int) string {
	if len(docs) == 0 {
		return styles.CardMeta.Render("(no documents)")
	}
	parts := make([]string, len(docs))
	for i, doc := range docs {
		parts[i] = tabLabel(doc.Label, i == selected)
	}
	return strings.Join(parts, styles.CardMeta.Render("  │  "))
}

// releaseDocBody returns the rendered body lines for the selected document,
// substituting a read-only empty state when the document is absent or blank.
func releaseDocBody(docs []data.ReleaseDoc, selected, width int) []string {
	if len(docs) == 0 {
		return []string{styles.CardMeta.Render("(no release documents available)")}
	}
	if selected < 0 || selected >= len(docs) {
		selected = 0
	}
	doc := docs[selected]
	switch {
	case !doc.Available:
		return []string{styles.CardMeta.Render("(" + doc.Label + " not found at " + doc.Path + ")")}
	case strings.TrimSpace(doc.Body) == "":
		return []string{styles.CardMeta.Render("(" + doc.Label + " is empty)")}
	default:
		return renderReleaseDocBody(doc.Body, width)
	}
}

// renderReleaseDocBody renders a raw Markdown document into wrapped, styled
// display lines. Unlike epicDetailBody it preserves the document's own line
// structure — blank lines and leading indentation survive — because supporting
// docs contain lists, tables, and code blocks whose layout carries meaning.
// Each source line is wrapped independently to the body width; WrapText is not
// applied across lines because it collapses whitespace and paragraph breaks.
// Frontmatter is stripped for a clean read-only view.
func renderReleaseDocBody(content string, width int) []string {
	if width < 4 {
		width = 4
	}
	var body []string
	for _, line := range stripFrontmatter(content) {
		body = append(body, wrapDocLine(line, width)...)
	}
	return body
}

// wrapDocLine renders a single source line: blank lines stay blank, Markdown
// headings take the board's heading styles, and every other line keeps its
// leading indentation while wrapping its content to the body width.
func wrapDocLine(line string, width int) []string {
	trimmed := strings.TrimSpace(line)
	switch {
	case trimmed == "":
		return []string{""}
	case strings.HasPrefix(trimmed, "# "):
		return styledWrap(strings.TrimPrefix(trimmed, "# "), "", width, styles.EpicTitleFocused)
	case strings.HasPrefix(trimmed, "## "):
		return styledWrap(strings.TrimPrefix(trimmed, "## "), "", width, styles.EpicItemFocused)
	case strings.HasPrefix(trimmed, "### "):
		return styledWrap(strings.TrimPrefix(trimmed, "### "), "", width, styles.EpicItemFocused)
	default:
		return styledWrap(trimmed, leadingWhitespace(line), width, styles.CardMeta)
	}
}

// styledWrap wraps text to width (reserving room for indent), prefixes each
// wrapped line with indent, and applies style. Release Docs content may contain
// code blocks and tables, so wrapping preserves interior whitespace via
// wrapDocText instead of routing through WrapText's word normalization. It
// always returns at least one line so blank-but-indented source lines still
// occupy a row. (Tabs in the text are normalized to spaces by the style's
// render layer, not here.)
func styledWrap(text, indent string, width int, style lipgloss.Style) []string {
	avail := width - len([]rune(indent))
	if avail < 4 {
		avail = 4
	}
	wrapped := wrapDocText(text, avail)
	if len(wrapped) == 0 {
		return []string{style.Render(indent)}
	}
	out := make([]string, len(wrapped))
	for i, w := range wrapped {
		out[i] = style.Render(indent + w)
	}
	return out
}

// wrapDocText wraps text to width while preserving interior whitespace within a
// kept line, breaking at the last space or tab that fits. Leading whitespace on
// continuation lines (the break point) is dropped so wrapped rows start flush.
func wrapDocText(text string, width int) []string {
	if width < 4 {
		width = 4
	}
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	var lines []string
	for len(runes) > width {
		breakAt := docWrapBreak(runes, width)
		lines = append(lines, strings.TrimRight(string(runes[:breakAt]), " \t"))
		runes = runes[breakAt:]
		for len(runes) > 0 && (runes[0] == ' ' || runes[0] == '\t') {
			runes = runes[1:]
		}
	}
	lines = append(lines, string(runes))
	return lines
}

// docWrapBreak returns the rune index to break runes at so the first line fits
// within width. It prefers the last space or tab within width; failing that it
// falls back to the first natural segment boundary, and finally a hard cut.
func docWrapBreak(runes []rune, width int) int {
	limit := min(width, len(runes))
	for i := limit; i > 0; i-- {
		if runes[i-1] == ' ' || runes[i-1] == '\t' {
			return i
		}
	}
	if segs := SplitLongWord(string(runes[:limit]), width); len(segs) > 0 {
		return len([]rune(segs[0]))
	}
	return limit
}

// leadingWhitespace returns the run of spaces and tabs at the start of line.
func leadingWhitespace(line string) string {
	return line[:len(line)-len(strings.TrimLeft(line, " \t"))]
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
