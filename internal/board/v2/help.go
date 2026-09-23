package v2

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/styles"
)

// renderHelp is deliberately derived from the same action set as key
// dispatch. A capability absent from the focused gate decision cannot appear
// here as a tempting key that only refuses when pressed.
func renderHelp(model Model, width, height int) string {
	textWidth := width - 4
	if textWidth < 1 {
		textWidth = 1
	}

	lines := []string{
		styles.ColumnTitleFocused.Render("KEYBOARD HELP"),
		fitLine(strings.Repeat("─", textWidth), textWidth),
		helpRow("↑↓ / j k", "move within the focused surface"),
		helpRow("←→ / h l", "move between columns, and into/out of Objectives at the edge"),
		helpRow(goalSelectorKey, "open the Goal selector"),
		helpRow("enter / v", "open the focused record"),
		helpRow("i / I", "open Issues"),
		helpRow("?", "close this help"),
	}

	if target, ok := model.focusedActionTarget(); ok {
		gateActions := actionsForRecord(model.State.Index, target)
		actions := append([]BoardAction(nil), gateActions...)
		if selection, ok := selectionAction(model.State.Index, target); ok {
			actions = append(actions, selection)
		}
		if len(actions) > 0 {
			lines = append(lines, "", styles.ColumnTitle.Render("OWNER ACTIONS"))
			for _, action := range actions {
				lines = append(lines, helpRow(action.Key, action.Label))
			}
		}
		if len(gateActions) == 0 {
			if decision, _, resolved := recordDecision(model.State.Index, target); resolved {
				lines = append(lines, "", styles.ColumnTitle.Render("CURRENT GATE"), helpRow("", decisionRefusal(decision)))
			}
		}
	}

	lines = append(lines, "", styles.CardMeta.Render("Executor and checker work belongs to their sessions."), styles.CardMeta.Render("esc/q:close"))
	return lipgloss.NewStyle().Width(width).MaxHeight(height).Render(strings.Join(lines, "\n"))
}

func helpRow(key, label string) string {
	if key == "" {
		return styles.CardMeta.Render(label)
	}
	return fmt.Sprintf("%s%s", styles.ColumnTitle.Render(key+":"), styles.CardMeta.Render(" "+label))
}
