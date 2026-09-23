package v2

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
	"github.com/opencode/savepoint/internal/styles"
)

// This file draws an open detail and resolves nothing: it is handed the
// RecordDetail detail.go already settled, and it reaches no index, no record
// map, and no resolver. What is on screen is exactly what was resolved.
//
// The evidence wording is internal/resume's, called rather than copied — the
// same decision the Next area records. A detail that spelled out its own
// clearance, dependency, exception, or replan sentences would be a second copy
// of the vocabulary with nothing keeping the two in step (STYLE-07, STYLE-09).

const (
	// detailIndent is the gutter a section's lines sit in, so a heading and its
	// content read as one block.
	detailIndent = "  "
	// detailTimeFormat renders a recorded timestamp as it was recorded, rather
	// than as a friendlier approximation of it: the exact instant is the part a
	// reader is checking.
	detailTimeFormat = time.RFC3339
	// noCheckRecorded is the Check section for a record with no Check. A record
	// with no evidence says so; it never renders an empty section, because
	// "nothing recorded" and "nothing shown" are different statements.
	noCheckRecorded = "(no Check recorded)"
	// dependencySatisfied is the absence of a block. Every reason a dependency
	// is *not* satisfied is internal/resume's wording; this is the one case
	// that has no reason to state.
	dependencySatisfied = "Satisfied."
	// notRecorded is what an optional reference the record did not declare
	// reads as.
	notRecorded = "(none recorded)"
)

// renderDetail draws the overlay at the board's body size, in the columns'
// own frame so opening one moves no geometry: the heading, a rule, and as much
// of the record as the height budget fits, with an indicator for whatever is
// scrolled out of view.
func renderDetail(detail RecordDetail, width, height, offset int) string {
	textW := columnTextWidth(width)
	bodyH := columnBodyHeight(height)

	lines := []string{
		styles.ColumnTitleFocused.Render(detailHeading(detail.Kind)),
		styles.Divider.Render(strings.Repeat("─", textW)),
	}
	lines = append(lines, detailWindow(detailLines(detail, textW), detailBudget(bodyH), offset)...)

	return frameColumn(lines, textW, bodyH, true)
}

func detailHeading(kind DetailKind) string {
	return detailKindLabel(kind) + " DETAIL"
}

// detailKindLabel is the public name of a detail kind: the Release-backed
// kind is presented as a Goal.
func detailKindLabel(kind DetailKind) string {
	if kind == DetailRelease {
		return strings.ToUpper(goalLabel)
	}
	return string(kind)
}

// detailLines is the overlay's whole content, already folded to width, in the
// order a reader needs it: who this record is, what it waits on, what evidence
// stands behind its badge, the Check chain, the follow-ups, and last the body.
func detailLines(detail RecordDetail, width int) []string {
	lines := identityRows(detail, width)

	switch detail.Kind {
	case DetailRelease:
		lines = append(lines, detailSection(strings.ToUpper(goalLabel)+" PROMISE", releasePromiseLines(detail), width)...)
		lines = append(lines, detailSection("MEMBER OBJECTIVES", releaseObjectiveLines(detail.MemberObjectives), width)...)
	case DetailObjective:
		lines = append(lines, detailSection("OWNED TASKS", refLines(detail.OwnedTasks), width)...)
		lines = append(lines, detailSection("OBJECTIVE DEPENDENCIES", objectiveDependencyLines(detail.ObjectiveDependencies), width)...)
	default:
		lines = append(lines, detailSection("DEPENDENCIES", dependencyLines(detail.Dependencies), width)...)
	}

	waived := detail.Evidence != nil && detail.Evidence.CheckWaiver != nil
	lines = append(lines, detailSection("CLEARANCE", clearanceLines(detail.Kind, detail.Clearance, waived, detail.ByException), width)...)
	if detail.Kind == DetailRelease {
		lines = append(lines, detailSection(strings.ToUpper(goalLabel)+" READINESS", releaseReadinessLines(detail), width)...)
		lines = append(lines, detailSection("OWNER VALIDATION", releaseOwnerValidationLines(detail), width)...)
	} else {
		lines = append(lines, detailSection("OWNER VALIDATION", ownerValidationLines(detail.Evidence), width)...)
	}
	lines = append(lines, detailSection("EXCEPTION", exceptionLines(detail.Evidence), width)...)
	lines = append(lines, detailSection("REPLAN", replanLines(detail.Evidence), width)...)
	lines = append(lines, detailSection("CHECKS", checkHistoryLines(detail.Checks), width)...)
	lines = append(lines, detailSection("ISSUES", issueLines(detail.Issues), width)...)
	if detail.Kind == DetailRelease && detail.LegacyCompletion != nil {
		lines = append(lines, detailSection("HISTORICAL COMPLETION", historicalCompletionLines(detail), width)...)
	}
	lines = append(lines, detailSection("BODY", bodyLines(detail.Body), width)...)

	return lines
}

// identityRows name the record: its ID, title, and recorded lifecycle, plus the
// Objective a Task belongs to. A Task's stage row is always present, because
// "no stage" is a fact about a planned or closed Task rather than a missing
// field.
func identityRows(detail RecordDetail, width int) []string {
	rows := []string{
		fieldRow("ID", detail.ID),
		fieldRow("Title", detail.Title),
		fieldRow("Status", string(detail.Status)),
	}
	if detail.Kind == DetailTask {
		rows = append(rows, fieldRow("Stage", orNone(stageLabel(detail.Stage))))
	}
	if detail.Owner != nil {
		rows = append(rows, fieldRow("Objective", refLabel(*detail.Owner)))
	}

	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		for _, line := range wrapDetailLine(row, width) {
			lines = append(lines, styles.CardMeta.Render(line))
		}
	}
	return lines
}

// releasePromiseLines exposes the authored Release promise as structured
// fields, while BODY below still carries the source markdown verbatim.
func releasePromiseLines(detail RecordDetail) []string {
	return []string{
		fieldRow("Outcome", orNone(detail.Outcome)),
		fieldRow("Why", orNone(detail.Why)),
		fieldRow("Success Conditions", orNone(detail.SuccessConditions)),
		fieldRow("Boundaries", orNone(detail.Boundaries)),
	}
}

// releaseObjectiveLines reports derived membership without implying that a
// Release owns the Objective or its Tasks. The counts come from the resolved
// detail value; the renderer reads no index and runs no completion resolver.
func releaseObjectiveLines(entries []ReleaseObjectiveProgress) []string {
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		line := fmt.Sprintf("%s — Tasks %d/%d done; %s", refLabel(entry.Objective), entry.TasksDone, entry.TasksTotal, objectiveCompletionPhrase(entry.Decision))
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return []string{"(no Objectives in this Goal)"}
	}
	return lines
}

func objectiveCompletionPhrase(decision data.GateDecision) string {
	if decision.Allowed {
		if decision.AllowedByException {
			return "completion allowed by exception"
		}
		return "completion allowed"
	}
	if len(decision.Blockers) == 0 {
		return "completion not allowed"
	}
	if decision.Blockers[0].Detail != "" {
		return "completion blocked: " + decision.Blockers[0].Detail
	}
	return "completion blocked"
}

// releaseReadinessLines reports the canonical Release completion decision.
// Historical proof and exception permission stay distinct from current V2
// technical clearance, even when the gate allows the Release to be done.
func releaseReadinessLines(detail RecordDetail) []string {
	decision := detail.ReleaseDecision
	if decision == nil {
		return []string{"No Goal completion decision was resolved."}
	}
	if decision.AllowedByLegacyCompletion {
		return []string{resume.HistoricalCompletionPhrase(detail.ID, decision.LegacyCompletion)}
	}
	if decision.AllowedByException {
		return []string{"Completion: " + resume.ExceptionPhrase(decision.Exception)}
	}
	if decision.Allowed {
		return []string{resume.ReleaseAcceptancePhraseForEvidence(detail.ID, detail.Evidence, &detail.Clearance)}
	}

	lines := make([]string, 0, len(decision.Blockers))
	for _, blocker := range decision.Blockers {
		lines = append(lines, resume.ReleaseBlockerPhrase(blocker))
	}
	if len(lines) == 0 {
		return []string{"Goal completion is not currently allowed."}
	}
	return lines
}

// releaseOwnerValidationLines makes the Release's mandatory owner boundary
// explicit. Unlike Task and Objective records, a Release does not need an
// owner_validation.required flag to require acceptance.
func releaseOwnerValidationLines(detail RecordDetail) []string {
	lines := []string{"Required: yes (Goal completion)"}
	checkID := detail.Clearance.Check
	if checkID == "" {
		return append(lines, "Accepted: (not recorded; a current Goal Check is required first)")
	}

	if detail.Evidence != nil && detail.Evidence.OwnerValidation != nil {
		accepted := detail.Evidence.OwnerValidation
		if accepted.AcceptedCheck == checkID {
			return append(lines, fmt.Sprintf("Accepted: Check %s, by %s", accepted.AcceptedCheck, resume.ActorLabel(accepted.AcceptedBy)))
		}
		if accepted.AcceptedCheck != "" {
			lines = append(lines, fmt.Sprintf("Accepted: Check %s, by %s", accepted.AcceptedCheck, resume.ActorLabel(accepted.AcceptedBy)))
			lines = append(lines, fmt.Sprintf("Current Check: %s (acceptance is not current)", checkID))
			return lines
		}
	}
	return append(lines, "Accepted: "+notRecorded)
}

func historicalCompletionLines(detail RecordDetail) []string {
	reference := detail.LegacyCompletion
	if reference == nil {
		return []string{resume.HistoricalCompletionPhrase(detail.ID, nil)}
	}
	return []string{
		resume.HistoricalCompletionPhrase(detail.ID, reference),
		"Source: " + reference.SourcePath,
		"Archive: " + reference.ArchivePath,
		"SHA-256: " + reference.SHA256,
	}
}

// detailSection renders one titled block, or nothing at all when the section
// has no content. A leading blank line separates it from what came before.
func detailSection(heading string, body []string, width int) []string {
	if len(body) == 0 {
		return nil
	}
	lines := []string{"", styles.ColumnTitle.Render(heading)}
	for _, entry := range body {
		for _, line := range wrapDetailLine(detailIndent+entry, width) {
			lines = append(lines, styles.CardMeta.Render(line))
		}
	}
	return lines
}

// clearanceLines state the resolved clearance twice over: the badge a card or
// sidebar row would show, so the overlay agrees with them at a glance, and
// the sentence that says which of the states it is and on what recorded
// basis. All three read the one value ResolveClearance returned; none
// inspects a Check. waived is meaningful only for a Task detail — an
// Objective's Check is never waivable — and reads whether the record's own
// Evidence carries a CheckWaiver, the same fact the card badge reads off
// GateDecision.AllowedByWaiver while the Task is still open. byException is
// meaningful only for an Objective detail, and only changes the badge
// (objectiveCheckBadge folds it to a plain green tick, matching the sidebar
// row); the sentence below the badge still states the true clearance state,
// and the EXCEPTION section further down still names the recorded owner,
// Check, and reason in full — nothing here hides that an exception was
// used, only the compact badge collapses it into the same signal a passed
// Check gives.
func clearanceLines(kind DetailKind, clearance data.Clearance, waived, byException bool) []string {
	badge := objectiveCheckBadge(clearance.State, byException)
	if kind == DetailTask {
		badge = taskCheckBadge(clearance.State, waived)
	}
	return []string{
		badge.Text(),
		resume.ClearancePhrase(&clearance),
	}
}

// dependencyLines name each declared Task dependency, the level it requires,
// and whether the resolver found it satisfied. Nothing here reads the
// dependency Task's own status or evidence.
func dependencyLines(entries []DependencyEntry) []string {
	lines := make([]string, 0, len(entries)*2)
	for _, entry := range entries {
		lines = append(lines, fmt.Sprintf("%s — requires %s", refLabel(entry.Target), entry.Requires))
		if entry.Decision.Satisfied {
			lines = append(lines, detailIndent+dependencySatisfied)
			continue
		}
		lines = append(lines, detailIndent+resume.DependencyPhrase(entry.Decision.Block))
	}
	return lines
}

// objectiveDependencyLines are the same for the Objectives an Objective waits
// on, resolved through ResolveObjectiveDependency rather than by reading the
// dependency Objective's status.
func objectiveDependencyLines(entries []ObjectiveDependencyEntry) []string {
	lines := make([]string, 0, len(entries)*2)
	for _, entry := range entries {
		lines = append(lines, refLabel(entry.Target))
		if entry.Decision.Satisfied {
			lines = append(lines, detailIndent+dependencySatisfied)
			continue
		}
		lines = append(lines, detailIndent+resume.ObjectiveDependencyPhrase(entry.Decision.Block))
	}
	return lines
}

// ownerValidationLines report whether the record declares that the owner must
// accept its outcome, and which Check the owner accepted. An absent block is
// the record saying acceptance is not required, so the section is always
// rendered rather than silently omitted.
func ownerValidationLines(evidence *data.Evidence) []string {
	validation := (*data.OwnerValidation)(nil)
	if evidence != nil {
		validation = evidence.OwnerValidation
	}
	if validation == nil || !validation.Required {
		return []string{"Required: no"}
	}

	lines := []string{"Required: yes"}
	if validation.AcceptedCheck == "" {
		return append(lines, "Accepted: "+notRecorded)
	}
	return append(lines, fmt.Sprintf("Accepted: Check %s, by %s",
		validation.AcceptedCheck, resume.ActorLabel(validation.AcceptedBy)))
}

// exceptionLines report a recorded exception as an exception — never as a
// clearance result — together with the requirement IDs it waives, the owner
// who recorded it, when, and the Check it applies to.
func exceptionLines(evidence *data.Evidence) []string {
	if evidence == nil || evidence.Exception == nil {
		return nil
	}
	exception := evidence.Exception
	return []string{
		resume.ExceptionPhrase(exception),
		"Requirements: " + strings.Join(exception.Requirements, ", "),
		"Applies to: Check " + exception.Check,
		fmt.Sprintf("Recorded by owner %s on %s", exception.Owner, exception.RecordedAt.Format(detailTimeFormat)),
	}
}

// replanLines report a recorded replan by its own reason and provenance.
func replanLines(evidence *data.Evidence) []string {
	if evidence == nil || evidence.Replan == nil {
		return nil
	}
	replan := evidence.Replan
	return []string{
		resume.ReplanPhrase(replan.Reason),
		fmt.Sprintf("Recorded by %s on %s", resume.ActorLabel(replan.RecordedBy), replan.RecordedAt.Format(detailTimeFormat)),
	}
}

// checkHistoryLines render the whole chain in recorded order, marking the one
// that counts and the ones a later run replaced. No entry is dropped, and a
// record with no Check says so rather than showing an empty section.
func checkHistoryLines(entries []CheckEntry) []string {
	if len(entries) == 0 {
		return []string{noCheckRecorded}
	}
	lines := make([]string, 0, len(entries)*2)
	for _, entry := range entries {
		lines = append(lines, fmt.Sprintf("%s  %s%s", entry.Check.ID, entry.Check.Result, checkMarker(entry)))
		lines = append(lines, detailIndent+fmt.Sprintf("recorded by %s on %s",
			resume.ActorLabel(entry.Check.CheckedBy), entry.Check.CheckedAt.Format(detailTimeFormat)))
	}
	return lines
}

// checkMarker says where one Check sits in the chain. Both markers can be
// printed together: the pair is unreachable in a project the loader accepts,
// and reporting it honestly beats choosing one of two contradictory facts.
func checkMarker(entry CheckEntry) string {
	switch {
	case entry.Latest && entry.Superseded:
		return "  [latest, superseded]"
	case entry.Latest:
		return "  [latest]"
	case entry.Superseded:
		return "  [superseded]"
	default:
		return ""
	}
}

// issueLines list the follow-ups linked to this record, each by its ID, type,
// status, and title.
func issueLines(issues []*data.IssueV2) []string {
	lines := make([]string, 0, len(issues))
	for _, issue := range issues {
		lines = append(lines, resume.IssueLine(issue))
	}
	return lines
}

// bodyLines carry the record's own markdown, split on its line breaks and
// nothing else. The body is author-owned prose: it is displayed, not parsed,
// so no heading, checkbox, list marker, or fenced block here means anything to
// the board.
func bodyLines(body string) []string {
	trimmed := strings.Trim(body, "\n")
	if strings.TrimSpace(trimmed) == "" {
		return nil
	}
	return strings.Split(trimmed, "\n")
}

// refLines name a list of referenced records.
func refLines(refs []RecordRef) []string {
	lines := make([]string, 0, len(refs))
	for _, ref := range refs {
		lines = append(lines, refLabel(ref))
	}
	return lines
}

// refLabel names one referenced record by identity, title, and recorded status.
func refLabel(ref RecordRef) string {
	label := ref.ID
	if ref.Title != "" {
		label += " — " + ref.Title
	}
	if ref.Status != "" {
		label += " (" + string(ref.Status) + ")"
	}
	return label
}

func fieldRow(label, value string) string {
	return label + ": " + value
}

func orNone(value string) string {
	if value == "" {
		return notRecorded
	}
	return value
}

// wrapDetailLine folds one plain line into width, hard-breaking a token too
// long to fit rather than letting it push the frame sideways: long and wide
// content costs extra lines inside the overlay and never a wider board. A
// line's own leading indent is carried onto its continuations, so a wrapped
// sentence stays inside the block it belongs to instead of starting again at
// the frame.
func wrapDetailLine(text string, width int) []string {
	if width < 1 {
		width = 1
	}

	indent := text[:len(text)-len(strings.TrimLeft(text, " "))]
	body := width - len(indent)
	if body < 1 {
		indent, body = "", width
	}

	wrapped := strings.Split(lipgloss.NewStyle().Width(body).Render(strings.TrimLeft(text, " ")), "\n")
	for i, line := range wrapped {
		wrapped[i] = indent + line
	}
	return wrapped
}

// detailBudget is how many content lines the overlay has: its frame's body
// height, less the heading and the rule above the content.
func detailBudget(bodyH int) int {
	budget := bodyH - columnHeaderLines
	if budget < 1 {
		return 1
	}
	return budget
}

// detailWindow is the run of content lines the overlay shows, reserving a line
// for each scroll indicator it needs so an indicator never pushes a line out of
// the frame it was measured for.
func detailWindow(lines []string, budget, offset int) []string {
	if len(lines) <= budget {
		return lines
	}
	offset = clampDetailOffset(offset, len(lines), budget)

	available := budget
	if offset > 0 {
		available--
	}
	if offset+available < len(lines) {
		available--
	}
	if available < 1 {
		available = 1
	}
	end := min(offset+available, len(lines))

	window := make([]string, 0, budget)
	if offset > 0 {
		window = append(window, scrollIndicator("↑", offset, "above"))
	}
	window = append(window, lines[offset:end]...)
	if end < len(lines) {
		window = append(window, scrollIndicator("↓", len(lines)-end, "more"))
	}
	return window
}

// clampDetailOffset keeps a scroll offset inside the content: zero at the top,
// and at the bottom the first line of the last full window, which spends one of
// its lines on the indicator naming what is above it.
func clampDetailOffset(offset, total, budget int) int {
	maxOffset := total - budget + 1
	if maxOffset < 0 {
		maxOffset = 0
	}
	if offset > maxOffset {
		return maxOffset
	}
	if offset < 0 {
		return 0
	}
	return offset
}

// detailScrollLimit is the furthest an overlay of this size can scroll. The
// reducer clamps against it so a key press that would move past the end changes
// nothing, rather than moving state the window then ignores.
func detailScrollLimit(detail RecordDetail, width, height int) int {
	total := len(detailLines(detail, columnTextWidth(width)))
	return clampDetailOffset(total, total, detailBudget(columnBodyHeight(height)))
}
