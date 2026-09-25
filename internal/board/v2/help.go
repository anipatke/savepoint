package v2

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode/savepoint/internal/styles"
)

// renderHelp describes the focused surface's keys and derives gated owner
// actions from the same decision used by dispatch.
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
	}
	if model.SidebarFocused {
		lines = append(lines,
			helpRow("1–4", "set priority: Critical, High, Medium, Low"),
			helpRow("K / shift+↑", "move the Objective up within its priority group"),
			helpRow("J / shift+↓", "move the Objective down within its priority group"),
		)
	}
	lines = append(lines,
		helpRow("enter / v", "open the focused record"),
		helpRow("i / I", "open Issues"),
		helpRow("?", "close this help"),
	)
	if model.Issues != nil && model.Issues.Detail == nil {
		lines = append(lines,
			helpRow("space", "advance the selected Issue"),
			helpRow("backspace", "retreat the selected Issue"),
		)
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
	if key == objectiveCloseKey {
		key = "space"
	}
	return fmt.Sprintf("%s%s", styles.ColumnTitle.Render(key+":"), styles.CardMeta.Render(" "+label))
}
