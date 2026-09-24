package v2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/opencode/savepoint/internal/data"
)

// actionMsg is the typed result of every V2 board write. Update handles it by
// showing the result and, for a successful write, starting the same load
// command used at startup. Filesystem access never occurs in the reducer.
type actionMsg struct {
	message         string
	err             error
	reload          bool
	releaseRollback bool
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

func freshV2Index(root string) (*data.V2Index, error) {
	index, err := data.LoadV2Index(root)
	if err != nil {
		return nil, err
	}
	return index, nil
}

// writeOwnerAcceptanceCmd re-reads the project and resolves the completion
// gate immediately before preparing the evidence write. A stale in-memory
// card can therefore never turn an owner wait into an acceptance for a newer
// or different Check.
func writeOwnerAcceptanceCmd(root string, target actionTarget) tea.Cmd {
	return func() tea.Msg {
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

// writeExceptionCompletionCmd completes a record through a recorded exception.
func writeExceptionCompletionCmd(root string, target actionTarget) tea.Cmd {
	return func() tea.Msg {
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

		closed := data.RouterSelectionV2{}
		switch target.Kind {
		case DetailTask:
			task := index.Tasks[target.ID]
			task.Status = data.ColumnDone
			task.Stage = ""
			if err := data.WriteTaskV2(task); err != nil {
				return actionFailure(err, "task completion")
			}
			closed = data.RouterSelectionV2{Objective: task.Objective, Task: task.ID}
		case DetailObjective:
			objective := index.Objectives[target.ID]
			objective.Status = data.ColumnDone
			if err := data.WriteObjectiveV2(objective); err != nil {
				return actionFailure(err, "objective completion")
			}
			closed = data.RouterSelectionV2{Objective: objective.ID}
		default:
			return actionMsg{err: fmt.Errorf("exception completion is not supported for %s", target.Kind)}
		}

		return completedRecordAction(root, closed, fmt.Sprintf("%s completed by recorded exception.", target.ID), data.WriteRouterStateV2)
	}
}

// writeObjectiveCompletionCmd closes a cleared Objective after resolving its
// gate against the latest project state. Exception-only closure keeps its x key.
func writeObjectiveCompletionCmd(root, objectiveID string) tea.Cmd {
	return func() tea.Msg {
		index, err := freshV2Index(root)
		if err != nil {
			return actionFailure(err, "objective completion")
		}
		objective, ok := index.Objectives[objectiveID]
		if !ok {
			return actionMsg{err: fmt.Errorf("objective %s is no longer present", objectiveID)}
		}
		if objective.Status == data.ColumnDone {
			return actionMsg{err: fmt.Errorf("%s is already done", objectiveID)}
		}
		decision := data.ResolveObjectiveCompletion(index, objectiveID)
		if !decision.Allowed || decision.AllowedByException {
			return actionMsg{err: fmt.Errorf("cannot close %s: %s", objectiveID, decisionRefusal(decision))}
		}
		objective.Status = data.ColumnDone
		if err := data.WriteObjectiveV2(objective); err != nil {
			return actionFailure(err, "objective completion")
		}
		return completedRecordAction(root, data.RouterSelectionV2{Objective: objectiveID}, fmt.Sprintf("%s completed.", objectiveID), data.WriteRouterStateV2)
	}
}

type routerStateWriter func(string, data.RouterSelectionV2, time.Time) error

// completedRecordAction updates the router only after the completion write
// succeeds. If selection persistence fails, the completed record remains and
// the reload surfaces the stale-selection diagnostic alongside this message.
func completedRecordAction(root string, closed data.RouterSelectionV2, message string, writer routerStateWriter) tea.Msg {
	if err := advanceRouterAfterClosure(root, closed, writer); err != nil {
		closedID := closed.Task
		if closedID == "" {
			closedID = closed.Objective
		}
		return actionMsg{
			err:    fmt.Errorf("%s completed, but the router still selects %s; the selection could not be advanced: %w", closedID, closedID, err),
			reload: true,
		}
	}
	return actionMsg{message: message, reload: true}
}

func advanceRouterAfterClosure(root string, closed data.RouterSelectionV2, writer routerStateWriter) error {
	index, err := freshV2Index(root)
	if err != nil {
		return fmt.Errorf("reload project after completion: %w", err)
	}

	path := filepath.Join(root, "router.md")
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat router.md: %w", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read router.md: %w", err)
	}
	router, err := data.NewRouterReader().ReadStateV2(string(content))
	if err != nil {
		return fmt.Errorf("read router selection: %w", err)
	}
	current := data.RouterSelectionV2{
		Release:   router.Release,
		Objective: router.Objective,
		Task:      router.Task,
		Issue:     router.Issue,
	}
	selection, changed := data.RouterSelectionAfterClosureV2(index, current, closed)
	if !changed {
		return nil
	}
	if writer == nil {
		writer = data.WriteRouterStateV2
	}
	if err := writer(root, selection, info.ModTime()); err != nil {
		return fmt.Errorf("write router selection: %w", err)
	}
	return nil
}

// writeTaskAdvanceCmd is Space on a focused Task: start a planned Task, move
// an in-progress one to its next stage, or complete one at stage check —
// whichever ResolveTaskStart/ResolveTaskAdvance/ResolveTaskCompletion already
// govern for its current lifecycle position, re-resolved fresh so a stale
// in-memory card can never write a transition its own gate would refuse.
// Completion is never forced: an in-progress Task at stage check only
// reaches done when its own gate decision is already Allowed — through
// current clearance, a recorded exception, or an owner Task-check waiver —
// exactly the same authority the board's other lifecycle writes require.
func writeTaskAdvanceCmd(root, taskID string) tea.Cmd {
	return func() tea.Msg {
		index, err := freshV2Index(root)
		if err != nil {
			return actionFailure(err, "task advance")
		}
		task, ok := index.Tasks[taskID]
		if !ok {
			return actionMsg{err: fmt.Errorf("task %s is no longer present", taskID)}
		}

		waivedNow, started := false, false
		switch {
		case task.Status == data.ColumnDone:
			return actionMsg{err: fmt.Errorf("%s is already done", taskID)}
		case task.Status == data.ColumnPlanned:
			decision := data.ResolveTaskStart(index, taskID)
			if !decision.Allowed {
				return actionMsg{err: fmt.Errorf("cannot start %s: %s", taskID, decisionRefusal(decision))}
			}
			task.Status, task.Stage = data.ColumnInProgress, data.StageBuild
			started = true
		case task.Stage == data.StageAudit:
			decision := data.ResolveTaskCompletion(index, taskID)
			if !decision.Allowed {
				clearance := data.ResolveClearance(index, taskID)
				if !clearanceIsMissing(clearance) {
					// A Check was actually run and found a problem — needs_work,
					// stale, or unverified. That result stands; Space never
					// silently overrides an independent checker's finding.
					return actionMsg{err: fmt.Errorf("cannot complete %s: %s", taskID, decisionRefusal(decision))}
				}
				// No Check was ever requested for this Task. Completing it here
				// is an interactive owner action — only a human at the keyboard
				// reaches this key — so that action itself is the explicit
				// Task-check waiver TEST-09 requires. Record it, then complete
				// under it, rather than refusing and making the owner write the
				// same fact into the file by hand first.
				if task.Evidence == nil {
					task.Evidence = &data.Evidence{}
				}
				task.Evidence.CheckWaiver = &data.CheckWaiver{
					Task:       taskID,
					Reason:     "Owner completed this Task via the board without requesting a Task Check.",
					Actor:      data.Actor{Role: data.ActorRoleOwner, Session: ownerBoardSession},
					RecordedAt: time.Now().UTC(),
				}
				if err := data.WriteTaskEvidenceV2(task); err != nil {
					return actionFailure(err, "task check waiver")
				}
				waivedNow = true
			}
			task.Status, task.Stage = data.ColumnDone, ""
		default:
			decision := data.ResolveTaskAdvance(index, taskID)
			if !decision.Allowed {
				return actionMsg{err: fmt.Errorf("cannot advance %s: %s", taskID, decisionRefusal(decision))}
			}
			next, err := data.AdvanceTaskLifecycleState(data.TaskLifecycleState{Status: task.Status, Stage: task.Stage})
			if err != nil {
				return actionFailure(err, "task advance")
			}
			task.Status, task.Stage = next.Status, next.Stage
		}

		if err := data.WriteTaskV2(task); err != nil {
			return actionFailure(err, "task advance")
		}
		message := taskLifecycleMessage(taskID, task.Status, task.Stage)
		if waivedNow {
			message = fmt.Sprintf("%s completed; no Task Check was requested, so an owner waiver was recorded.", taskID)
		}
		if task.Status == data.ColumnDone {
			closed := data.RouterSelectionV2{Objective: task.Objective, Task: task.ID}
			return completedRecordAction(root, closed, message, data.WriteRouterStateV2)
		}
		if started {
			return startedTaskAction(index, task, message, data.WriteObjectiveV2)
		}
		return actionMsg{message: message, reload: true}
	}
}

type objectiveWriter func(*data.ObjectiveV2) error

// startedTaskAction moves the started Task's owning Objective from planned to
// in_progress, so every surface that shows the Objective's recorded status
// agrees that its work has begun (I-044). It runs only after the Task write
// succeeded: a failed Objective write leaves the Task started and reports the
// error, and doctor's planned-with-started-Task warning names the Objective
// until it is set. An Objective already in_progress or done is never rewritten.
func startedTaskAction(index *data.V2Index, task *data.TaskV2, message string, writer objectiveWriter) tea.Msg {
	objective, ok := index.Objectives[task.Objective]
	if !ok || objective.Status != data.ColumnPlanned {
		return actionMsg{message: message, reload: true}
	}
	objective.Status = data.ColumnInProgress
	if err := writer(objective); err != nil {
		return actionMsg{
			err:    fmt.Errorf("%s started, but Objective %s is still planned; its status could not be set to in_progress: %w", task.ID, objective.ID, err),
			reload: true,
		}
	}
	return actionMsg{message: fmt.Sprintf("%s; Objective %s is now in progress.", strings.TrimSuffix(message, "."), objective.ID), reload: true}
}

// writeTaskRetreatCmd is Backspace on a focused Task: move it one lifecycle
// step backward. Retreat carries no Check gate — only the owner's own
// keypress reaches it, matching AGENTS.md's "only the user may retreat a
// Task to an earlier status".
func writeTaskRetreatCmd(root, taskID string) tea.Cmd {
	return func() tea.Msg {
		index, err := freshV2Index(root)
		if err != nil {
			return actionFailure(err, "task retreat")
		}
		task, ok := index.Tasks[taskID]
		if !ok {
			return actionMsg{err: fmt.Errorf("task %s is no longer present", taskID)}
		}

		prev, err := data.RetreatTaskLifecycleState(data.TaskLifecycleState{Status: task.Status, Stage: task.Stage})
		if err != nil {
			return actionMsg{err: fmt.Errorf("cannot retreat %s: %v", taskID, err)}
		}
		task.Status, task.Stage = prev.Status, prev.Stage

		if err := data.WriteTaskV2(task); err != nil {
			return actionFailure(err, "task retreat")
		}
		return actionMsg{message: taskLifecycleMessage(taskID, task.Status, task.Stage), reload: true}
	}
}

// taskLifecycleMessage is the one status-bar phrase both the advance and the
// retreat write share, so a card's next status line reads the same lifecycle
// vocabulary whichever direction moved it.
func taskLifecycleMessage(taskID string, status data.ColumnType, stage data.ProgressStage) string {
	if stage == "" {
		return fmt.Sprintf("%s moved to %s.", taskID, status)
	}
	return fmt.Sprintf("%s moved to %s, stage %s.", taskID, status, stage)
}

// writeSelectionCmd re-reads both the V2 index and router before passing only
// the requested selection to data.WriteRouterStateV2. The data writer owns
// byte preservation and its final freshness check; this command owns the
// target-validity boundary.
func writeSelectionCmd(root string, selection data.RouterSelectionV2, expectedMtime ...time.Time) tea.Cmd {
	return func() tea.Msg {
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
		expected := info.ModTime()
		if len(expectedMtime) > 0 && !expectedMtime[0].IsZero() {
			if !info.ModTime().Equal(expectedMtime[0]) {
				return actionFailure(data.ErrMtimeConflict, "router selection")
			}
			expected = expectedMtime[0]
		}
		if err := data.WriteRouterStateV2(root, selection, expected); err != nil {
			return actionFailure(err, "router.md")
		}

		return actionMsg{
			message: selectionMessage(selection),
			reload:  true,
		}
	}
}

func releaseSelectionFailure(err error, subject string) actionMsg {
	msg, ok := actionFailure(err, subject).(actionMsg)
	if !ok {
		msg = actionMsg{err: err}
	}
	msg.reload = true
	msg.releaseRollback = true
	return msg
}

// writeReleaseSelectionCmd records only the optional Release context. The
// current Objective and Task are read from the router immediately before the
// canonical writer runs, so switching context never invents ownership or
// rewrites next_action prose.
func writeReleaseSelectionCmd(root, release string, expectedMtime ...time.Time) tea.Cmd {
	return func() tea.Msg {
		index, err := freshV2Index(root)
		if err != nil {
			return releaseSelectionFailure(err, "Goal selection")
		}
		if release == "" {
			return releaseSelectionFailure(fmt.Errorf("Goal selection requires a live Goal"), "Goal selection")
		}
		if _, ok := index.Releases[release]; !ok {
			return releaseSelectionFailure(fmt.Errorf("selection target %s is no longer present", release), "Goal selection")
		}

		path := filepath.Join(root, "router.md")
		info, err := os.Stat(path)
		if err != nil {
			return releaseSelectionFailure(err, "Goal selection")
		}
		expected := info.ModTime()
		if len(expectedMtime) > 0 && !expectedMtime[0].IsZero() {
			if !info.ModTime().Equal(expectedMtime[0]) {
				return releaseSelectionFailure(data.ErrMtimeConflict, "Goal selection")
			}
			expected = expectedMtime[0]
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return releaseSelectionFailure(err, "Goal selection")
		}
		router, err := data.NewRouterReader().ReadStateV2(string(content))
		if err != nil {
			return releaseSelectionFailure(err, "Goal selection")
		}
		selection := data.RouterSelectionV2{
			Release:   release,
			Objective: router.Objective,
			Task:      router.Task,
			Issue:     router.Issue,
		}
		if objective := index.Objectives[selection.Objective]; objective == nil || string(objective.Release) != release {
			selection.Objective = ""
			selection.Task = ""
		}
		if err := validateSelectionAgainstIndex(index, selection); err != nil {
			return releaseSelectionFailure(err, "Goal selection")
		}
		if err := data.WriteRouterStateV2(root, selection, expected); err != nil {
			return releaseSelectionFailure(err, "router.md")
		}

		return actionMsg{
			message: fmt.Sprintf("Goal selected: %s.", release),
			reload:  true,
		}
	}
}

func validateSelectionAgainstIndex(index *data.V2Index, selection data.RouterSelectionV2) error {
	if selection.Release != "" {
		if _, ok := index.Releases[selection.Release]; !ok {
			return fmt.Errorf("selection target %s is no longer present", selection.Release)
		}
	}
	if selection.Objective != "" {
		objective, ok := index.Objectives[selection.Objective]
		if !ok {
			return fmt.Errorf("selection target %s is no longer present", selection.Objective)
		}
		if selection.Release != "" && string(objective.Release) != selection.Release {
			objectiveRelease := string(objective.Release)
			if objectiveRelease == "" {
				objectiveRelease = "(none)"
			}
			return fmt.Errorf("selection target %s belongs to Goal %s, not %s", selection.Objective, objectiveRelease, selection.Release)
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
