package v2

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/styles"
)

// A framed surface — a card or a column — spends two cells on its border and
// two on its horizontal padding. Lip Gloss's Width sets the padded block, so a
// frame of outer width w carries a block of w-borderCells and text of
// w-borderCells-paddingCells. Keeping both names here is what stops a rule or a
// title from being measured against the wrong one and wrapping a cell early.
const (
	borderCells  = 2
	paddingCells = 2
	cardChrome   = borderCells + paddingCells
)

// TaskCard is one Task together with the resolved values its badges read. It
// exists so rendering is pure: renderCard formats a TaskCard and calls no
// resolver, so what is on screen is exactly what the resolvers returned.
type TaskCard struct {
	Task      *data.TaskV2
	Clearance data.Clearance
	// Decision is the one gate decision that governs this Task's next move.
	// It is the zero decision for a done Task, which has no next move: a
	// closed Task's state is its clearance and how it closed.
	Decision data.GateDecision
	// ByException is true when completion is, or was, allowed only by a
	// recorded exception: GateDecision.AllowedByException while the Task is
	// still open, and the Task's own recorded exception once it is done and
	// the completion decision no longer applies to it.
	ByException bool
	// ByWaiver mirrors ByException for an owner-recorded Task-check waiver:
	// GateDecision.AllowedByWaiver while open, the Task's own recorded
	// CheckWaiver once done.
	ByWaiver bool
}

// newTaskCard resolves everything a card shows about task, through the
// existing resolvers and nothing else. The decision it reads is chosen exactly
// as data.ResolveNext chooses one — start for a planned Task, advance for one
// mid build or test, completion for one at audit — so a card and the Next area
// never report two different judgements about the same Task.
func newTaskCard(index *data.V2Index, task *data.TaskV2) TaskCard {
	card := TaskCard{Task: task, Clearance: data.ResolveClearance(index, task.ID)}

	switch {
	case task.Status == data.ColumnDone:
		card.ByException = task.Evidence != nil && task.Evidence.Exception != nil
		card.ByWaiver = task.Evidence != nil && task.Evidence.CheckWaiver != nil
		return card
	case task.Status == data.ColumnPlanned:
		card.Decision = data.ResolveTaskStart(index, task.ID)
	case task.Stage != data.StageAudit:
		card.Decision = data.ResolveTaskAdvance(index, task.ID)
	default:
		card.Decision = data.ResolveTaskCompletion(index, task.ID)
	}

	card.ByException = card.Decision.AllowedByException
	card.ByWaiver = card.Decision.AllowedByWaiver
	return card
}

// groupTaskCards resolves every Task in index and groups the cards by recorded
// status, in ascending Task ID order so a project renders the same way twice.
// The three groups are the three columns: no status is promoted to a column of
// its own, and a Task's group comes from TaskV2.Status alone.
func groupTaskCards(index *data.V2Index) map[data.ColumnType][]TaskCard {
	return groupTaskCardsForRelease(index, "", "")
}

// groupTaskCardsFor groups the cards for the Tasks in view: the ones the
// selected Objective owns, or every Task when no Objective is selected. Which
// Tasks those are is taskIDsInView's answer, read from index.ObjectiveTasks;
// grouping them by status is this function's only other job.
func groupTaskCardsFor(index *data.V2Index, objectiveID string) map[data.ColumnType][]TaskCard {
	return groupTaskCardsForRelease(index, "", objectiveID)
}

// groupTaskCardsForRelease is the card projection for a Release context. The
// release filter is resolved before cards are built, so rendering still sees
// only already-resolved TaskCard values and never performs membership work.
func groupTaskCardsForRelease(index *data.V2Index, releaseID, objectiveID string) map[data.ColumnType][]TaskCard {
	grouped := map[data.ColumnType][]TaskCard{
		data.ColumnPlanned:    {},
		data.ColumnInProgress: {},
		data.ColumnDone:       {},
	}
	if index == nil {
		return grouped
	}

	for _, id := range taskIDsInReleaseView(index, releaseID, objectiveID) {
		task := index.Tasks[id]
		grouped[task.Status] = append(grouped[task.Status], newTaskCard(index, task))
	}
	return grouped
}

// badges is the card's state line: stage, how a done Task closed, clearance,
// and every blocker that is not already stated by the clearance badge. The
// order is fixed so a reader learns one shape.
func (c TaskCard) badges() []Badge {
	var badges []Badge

	if badge, ok := stageBadge(c.Task.Status, c.Task.Stage); ok {
		badges = append(badges, badge)
	}
	if c.Task.Status == data.ColumnDone {
		badges = append(badges, completionBadge(c.Clearance.State, c.ByException, c.ByWaiver))
	}
	if c.Task.Status != data.ColumnPlanned {
		badges = append(badges, taskCheckBadge(c.Clearance.State, c.ByWaiver))
	}
	for _, blocker := range c.Decision.Blockers {
		if badge, ok := blockerBadge(blocker); ok {
			badges = append(badges, badge)
		}
	}
	if c.ByException && c.Task.Status != data.ColumnDone {
		badges = append(badges, exceptionBadge())
	}

	return badges
}

// renderCard draws one card at the given outer width: its T### identity, the
// human title its author wrote, and its badges. It reads only the resolved
// values on the card — no index, no resolver, no Check ID comparison, no
// freshness inspection, and no dependency evaluation.
//
// The title is the card's label. A V2 Task's objective field is an O###
// identity reference, never display language, so there is no fallback to it
// and no other source for the label: a titleless Task cannot reach here,
// because DecodeTaskV2 refuses to load one.
func renderCard(card TaskCard, width int, focused bool) string {
	textW := width - cardChrome
	if textW < 4 {
		textW = 4
	}

	identity := styles.CardMeta.Render(xansi.Truncate(card.Task.ID, textW, "…"))
	lines := []string{identity}
	tStyle := titleStyle(card.Task.Status, focused)
	for _, titleLine := range wrapTitleLines(card.Task.Title, textW, 2) {
		lines = append(lines, tStyle.Render(titleLine))
	}
	if badges := renderBadgeLines(card.badges(), textW); len(badges) > 0 {
		lines = append(lines, badges...)
	}

	return cardStyle(card.Task.Status, focused).Width(textW + paddingCells).Render(strings.Join(lines, "\n"))
}

// renderBadgeLines packs badges onto as few lines as fit within width, so a
// narrow card wraps between badges instead of through one.
func renderBadgeLines(badges []Badge, width int) []string {
	var lines []string
	current := ""
	for _, badge := range badges {
		rendered := badge.Render()
		switch {
		case current == "":
			current = rendered
		case lipgloss.Width(current)+2+lipgloss.Width(rendered) <= width:
			current += "  " + rendered
		default:
			lines = append(lines, current)
			current = rendered
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// cardStyle and titleStyle are the whole of focus presentation: the accent
// changes, the geometry does not. Both styles carry identical padding and
// border sides, so a focused card occupies exactly the cells an unfocused one
// does (the visual identity treats a layout that moves under focus as a bug).
func cardStyle(status data.ColumnType, focused bool) lipgloss.Style {
	if focused {
		if status == data.ColumnPlanned {
			return styles.CardBoxFocusedPlanned
		}
		if status == data.ColumnDone {
			return styles.CardBoxFocusedDone
		}
		return styles.CardBoxFocused
	}
	return styles.CardBox
}

func titleStyle(status data.ColumnType, focused bool) lipgloss.Style {
	if focused {
		if status == data.ColumnPlanned {
			return styles.TaskItem
		}
		if status == data.ColumnDone {
			return styles.TaskItemFocusedDone
		}
		return styles.TaskItemFocused
	}
	return styles.TaskItem
}
