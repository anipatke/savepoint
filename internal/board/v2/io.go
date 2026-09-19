package v2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

// actionMsg is the typed result of every V2 board write. Update handles it by
// showing the result and, for a successful write, starting the same load
// command used at startup. Filesystem access never occurs in the reducer.
type actionMsg struct {
	message string
	err     error
	reload  bool
}

func actionFailure(err error, subject string) tea.Msg {
	if err == nil {
		return actionMsg{err: fmt.Errorf("%s failed", subject)}
	}
	if errors.Is(err, data.ErrV2SourceConflict) || errors.Is(err, data.ErrMtimeConflict) {
		return actionMsg{err: fmt.Errorf("%s changed on disk: refresh before retrying", subject)}
	}
	return actionMsg{err: err}
}

func writeGuard(root string) (string, bool) {
	migration, guidance, err := pendingMigration(root)
	if err != nil {
		return err.Error(), true
	}
	if migration.Pending {
		return guidance, true
	}
	return "", false
}

func freshV2Index(root string) (*data.V2Index, error) {
	project, err := data.LoadProject(root)
	if err != nil {
		return nil, err
	}
	if project.SchemaVersion != data.SchemaVersionV2 || project.V2 == nil {
		return nil, fmt.Errorf("board write requires a V2 project")
	}
	return project.V2, nil
}

// writeOwnerAcceptanceCmd re-reads the project and resolves the completion
// gate immediately before preparing the evidence write. A stale in-memory
// card can therefore never turn an owner wait into an acceptance for a newer
// or different Check.
func writeOwnerAcceptanceCmd(root string, target actionTarget) tea.Cmd {
	return func() tea.Msg {
		if message, refused := writeGuard(root); refused {
			return actionMsg{err: fmt.Errorf("write refused: %s", message)}
		}

		index, err := freshV2Index(root)
		if err != nil {
			return actionFailure(err, "owner acceptance")
		}
		decision, clearance, ok := recordDecision(index, target)
		if !ok {
			return actionMsg{err: fmt.Errorf("owner acceptance target %s %s is no longer present", target.Kind, target.ID)}
		}
		if !clearanceIsCurrent(clearance) || !hasOwnerAcceptanceBlock(decision) {
			return actionMsg{err: fmt.Errorf("owner acceptance refused for %s: %s", target.ID, decisionRefusal(decision))}
		}

		checkID := clearance.Check
		switch target.Kind {
		case DetailTask:
			task := index.Tasks[target.ID]
			setOwnerAcceptance(&task.Evidence, checkID)
			if err := data.WriteTaskEvidenceV2(task); err != nil {
				return actionFailure(err, "task evidence")
			}
		case DetailObjective:
			objective := index.Objectives[target.ID]
			setOwnerAcceptance(&objective.Evidence, checkID)
			if err := data.WriteObjectiveEvidenceV2(objective); err != nil {
				return actionFailure(err, "objective evidence")
			}
		default:
			return actionMsg{err: fmt.Errorf("owner acceptance is not supported for %s", target.Kind)}
		}

		return actionMsg{
			message: fmt.Sprintf("Owner accepted Check %s for %s.", checkID, target.ID),
			reload:  true,
		}
	}
}

func setOwnerAcceptance(evidence **data.Evidence, checkID string) {
	if *evidence == nil {
		*evidence = &data.Evidence{}
	}
	if (*evidence).OwnerValidation == nil {
		(*evidence).OwnerValidation = &data.OwnerValidation{Required: true}
	}
	(*evidence).OwnerValidation.AcceptedCheck = checkID
	(*evidence).OwnerValidation.AcceptedBy = data.Actor{Role: data.ActorRoleOwner, Session: ownerBoardSession}
}

// writeExceptionCompletionCmd is the only lifecycle write exposed by this
// package: the gate has already granted owner authority through a recorded
// exception. It does not make executor or checker transitions available.
func writeExceptionCompletionCmd(root string, target actionTarget) tea.Cmd {
	return func() tea.Msg {
		if message, refused := writeGuard(root); refused {
			return actionMsg{err: fmt.Errorf("write refused: %s", message)}
		}

		index, err := freshV2Index(root)
		if err != nil {
			return actionFailure(err, "exception completion")
		}
		decision, _, ok := recordDecision(index, target)
		if !ok {
			return actionMsg{err: fmt.Errorf("exception completion target %s %s is no longer present", target.Kind, target.ID)}
		}
		if !decision.Allowed || decision.Actor != data.ActorRoleOwner || !decision.AllowedByException {
			return actionMsg{err: fmt.Errorf("exception completion refused for %s: %s", target.ID, decisionRefusal(decision))}
		}

		switch target.Kind {
		case DetailTask:
			task := index.Tasks[target.ID]
			task.Status = data.ColumnDone
			task.Stage = ""
			if err := data.WriteTaskV2(task); err != nil {
				return actionFailure(err, "task completion")
			}
		case DetailObjective:
			objective := index.Objectives[target.ID]
			objective.Status = data.ColumnDone
			if err := data.WriteObjectiveV2(objective); err != nil {
				return actionFailure(err, "objective completion")
			}
		default:
			return actionMsg{err: fmt.Errorf("exception completion is not supported for %s", target.Kind)}
		}

		return actionMsg{
			message: fmt.Sprintf("%s completed by recorded exception.", target.ID),
			reload:  true,
		}
	}
}

// writeSelectionCmd re-reads both the V2 index and router before passing only
// the requested selection to data.WriteRouterStateV2. The data writer owns
// byte preservation and its final freshness check; this command owns the
// board's pending-migration and target-validity boundary.
func writeSelectionCmd(root string, selection data.RouterSelectionV2) tea.Cmd {
	return func() tea.Msg {
		if message, refused := writeGuard(root); refused {
			return actionMsg{err: fmt.Errorf("write refused: %s", message)}
		}

		index, err := freshV2Index(root)
		if err != nil {
			return actionFailure(err, "selection")
		}
		if err := validateSelectionAgainstIndex(index, selection); err != nil {
			return actionMsg{err: err}
		}

		path := filepath.Join(root, "router.md")
		info, err := os.Stat(path)
		if err != nil {
			return actionFailure(err, "router selection")
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return actionFailure(err, "router selection")
		}
		if _, err := data.NewRouterReader().ReadStateV2(string(content)); err != nil {
			return actionFailure(err, "router selection")
		}
		if err := data.WriteRouterStateV2(root, selection, info.ModTime()); err != nil {
			return actionFailure(err, "router.md")
		}

		return actionMsg{
			message: selectionMessage(selection),
			reload:  true,
		}
	}
}

func validateSelectionAgainstIndex(index *data.V2Index, selection data.RouterSelectionV2) error {
	if selection.Objective != "" {
		if _, ok := index.Objectives[selection.Objective]; !ok {
			return fmt.Errorf("selection target %s is no longer present", selection.Objective)
		}
	}
	if selection.Task != "" {
		task, ok := index.Tasks[selection.Task]
		if !ok {
			return fmt.Errorf("selection target %s is no longer present", selection.Task)
		}
		if task.Objective != selection.Objective {
			return fmt.Errorf("selection target %s belongs to %s, not %s", selection.Task, task.Objective, selection.Objective)
		}
	}
	return nil
}

func selectionMessage(selection data.RouterSelectionV2) string {
	if selection.Task == "" {
		return fmt.Sprintf("Selection recorded: %s.", selection.Objective)
	}
	return fmt.Sprintf("Selection recorded: %s / %s.", selection.Objective, selection.Task)
}
