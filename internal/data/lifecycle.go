package data

import (
	"fmt"
)

const LegacyTaskStatusTodo ColumnType = "todo"
const LegacyTaskStatusComplete ColumnType = "complete"
const LegacyTaskStatusCompleted ColumnType = "completed"
const LegacyTaskStageImplementation ProgressStage = "implementation"

type TaskLifecycleMetadata struct {
	Status ColumnType
	Column ColumnType
	Stage  ProgressStage
	Phase  ProgressStage
}

type TaskLifecycleState struct {
	Status ColumnType
	Stage  ProgressStage
}

func CanonicalTaskStatuses() []ColumnType {
	return []ColumnType{ColumnPlanned, ColumnInProgress, ColumnDone}
}

// ParseTaskLifecycle heals recoverable lifecycle metadata at load time so a
// task file never blocks board load; doctor surfaces the same issues as
// notifications via DiagnoseTaskLifecycle.
func ParseTaskLifecycle(metadata TaskLifecycleMetadata) TaskLifecycleState {
	rawStatus := firstTaskStatus(metadata.Column, metadata.Status)
	state := TaskLifecycleState{
		Status: NormalizeTaskStatusForLoad(rawStatus),
	}
	if !IsCanonicalTaskStatus(state.Status) {
		state.Status = ColumnPlanned
	}

	if state.Status == ColumnInProgress {
		state.Stage = NormalizeTaskStageForLoad(firstProgressStage(metadata.Stage, metadata.Phase))
		if !IsCanonicalStage(state.Stage) {
			state.Stage = StageBuild
		}
	}

	return state
}

func NormalizeTaskStatusForLoad(value ColumnType) ColumnType {
	if value == "" {
		return ColumnPlanned
	}
	if status, ok := ResolveTaskStatusAlias(value); ok {
		return status
	}
	return value
}

func ResolveTaskStatusAlias(value ColumnType) (ColumnType, bool) {
	switch value {
	case LegacyTaskStatusTodo:
		return ColumnPlanned, true
	case LegacyTaskStatusComplete, LegacyTaskStatusCompleted:
		return ColumnDone, true
	default:
		return "", false
	}
}

func IsLegacyTaskStageAlias(value ProgressStage) bool {
	switch value {
	case LegacyTaskStageImplementation:
		return true
	default:
		return false
	}
}

func NormalizeTaskStageForLoad(value ProgressStage) ProgressStage {
	if value == LegacyTaskStageImplementation {
		return StageBuild
	}
	return value
}

func ValidateTaskLifecycleStateForWrite(state TaskLifecycleState) error {
	if !IsCanonicalTaskStatus(state.Status) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", state.Status)
	}

	return validateCanonicalTaskStage(state)
}

func AdvanceTaskLifecycleState(state TaskLifecycleState) (TaskLifecycleState, error) {
	state = normalizeTaskLifecycleStateForTransition(state)
	switch state.Status {
	case ColumnPlanned:
		return TaskLifecycleState{Status: ColumnInProgress, Stage: StageBuild}, nil
	case ColumnInProgress:
		switch state.Stage {
		case StageBuild:
			return TaskLifecycleState{Status: ColumnInProgress, Stage: StageTest}, nil
		case StageTest:
			return TaskLifecycleState{Status: ColumnInProgress, Stage: StageAudit}, nil
		case StageAudit:
			return TaskLifecycleState{Status: ColumnDone}, nil
		default:
			return TaskLifecycleState{}, fmt.Errorf("unknown stage %q", state.Stage)
		}
	case ColumnDone:
		return TaskLifecycleState{}, fmt.Errorf("task is already done")
	default:
		return TaskLifecycleState{}, fmt.Errorf("unknown status %q", state.Status)
	}
}

func RetreatTaskLifecycleState(state TaskLifecycleState) (TaskLifecycleState, error) {
	state = normalizeTaskLifecycleStateForTransition(state)
	switch state.Status {
	case ColumnDone:
		return TaskLifecycleState{Status: ColumnInProgress, Stage: StageAudit}, nil
	case ColumnInProgress:
		switch state.Stage {
		case StageAudit:
			return TaskLifecycleState{Status: ColumnInProgress, Stage: StageTest}, nil
		case StageTest:
			return TaskLifecycleState{Status: ColumnInProgress, Stage: StageBuild}, nil
		case StageBuild:
			return TaskLifecycleState{Status: ColumnPlanned}, nil
		default:
			return TaskLifecycleState{}, fmt.Errorf("unknown stage %q", state.Stage)
		}
	case ColumnPlanned:
		return TaskLifecycleState{Status: ColumnPlanned}, nil
	default:
		return TaskLifecycleState{}, fmt.Errorf("unknown status %q", state.Status)
	}
}

// NormalizeDefectLifecycleForLoad heals recoverable defect lifecycle metadata
// at load time so a defect file never blocks board load; doctor surfaces the
// same issues as notifications via DiagnoseDefectLifecycle.
func NormalizeDefectLifecycleForLoad(d *Defect) {
	d.Status = NormalizeDefectStatusForLoad(d.Status)
	if d.Status != DefectInProgress {
		d.Stage = ""
		return
	}
	d.Stage = NormalizeTaskStageForLoad(d.Stage)
	if !IsCanonicalStage(d.Stage) {
		d.Stage = StageBuild
	}
}

func NormalizeDefectStatusForLoad(value DefectStatus) DefectStatus {
	if status, ok := ResolveDefectStatusAlias(value); ok {
		return status
	}
	if !IsCanonicalDefectStatus(value) {
		return DefectOpen
	}
	return value
}

// ResolveDefectStatusAlias maps task-style statuses agents sometimes write
// into defect frontmatter onto the defect lifecycle.
func ResolveDefectStatusAlias(value DefectStatus) (DefectStatus, bool) {
	switch ColumnType(value) {
	case ColumnPlanned, LegacyTaskStatusTodo:
		return DefectOpen, true
	case ColumnDone, LegacyTaskStatusComplete, LegacyTaskStatusCompleted:
		return DefectResolved, true
	default:
		return "", false
	}
}

// EpicStatusAudited is the only canonical epic status that is not shared with
// the task status vocabulary; planned/in_progress/done reuse the Column*
// constants so the epic and task lifecycles never drift apart.
const EpicStatusAudited EpicStatus = "audited"

type EpicStatus string

// CanonicalEpicStatuses is the single source of the epic-status vocabulary.
func CanonicalEpicStatuses() []EpicStatus {
	return []EpicStatus{
		EpicStatus(ColumnPlanned),
		EpicStatus(ColumnInProgress),
		EpicStatus(ColumnDone),
		EpicStatusAudited,
	}
}

// ResolveEpicStatusAlias maps non-canonical values agents sometimes leak into
// epic frontmatter — task-style completions and stray router states — onto the
// canonical epic vocabulary. Router states are not valid epic statuses; this
// only heals them, mirroring ResolveDefectStatusAlias.
func ResolveEpicStatusAlias(value EpicStatus) (EpicStatus, bool) {
	switch ColumnType(value) {
	case LegacyTaskStatusComplete, LegacyTaskStatusCompleted:
		return EpicStatus(ColumnDone), true
	case LegacyTaskStatusTodo:
		return EpicStatus(ColumnPlanned), true
	}
	switch value {
	case "epic-design", "epic-task-breakdown", "task-building":
		return EpicStatus(ColumnPlanned), true
	default:
		return "", false
	}
}

func IsCanonicalDefectStatus(value DefectStatus) bool {
	switch value {
	case DefectOpen, DefectInProgress, DefectResolved:
		return true
	default:
		return false
	}
}

func IsCanonicalTaskStatus(value ColumnType) bool {
	switch value {
	case ColumnPlanned, ColumnInProgress, ColumnDone:
		return true
	default:
		return false
	}
}

func IsCanonicalStage(value ProgressStage) bool {
	switch value {
	case StageBuild, StageTest, StageAudit:
		return true
	default:
		return false
	}
}

func validateCanonicalTaskStage(state TaskLifecycleState) error {
	if state.Status == ColumnInProgress {
		if state.Stage == "" {
			return fmt.Errorf("stage is required when task status is in_progress. Add 'stage: build' to task frontmatter")
		}
		if !IsCanonicalStage(state.Stage) {
			return fmt.Errorf("invalid stage %q: use build, test, or audit. Add 'stage: build' to task frontmatter", state.Stage)
		}
		return nil
	}

	if state.Stage != "" {
		return fmt.Errorf("stage field %q is only valid when status is in_progress. Remove 'stage' or change status to in_progress", state.Stage)
	}

	return nil
}

func normalizeTaskLifecycleStateForTransition(state TaskLifecycleState) TaskLifecycleState {
	state.Status = NormalizeTaskStatusForLoad(state.Status)
	if state.Status != ColumnInProgress {
		state.Stage = ""
		return state
	}
	if state.Stage == "" {
		state.Stage = StageBuild
	}
	return state
}

func firstTaskStatus(values ...ColumnType) ColumnType {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func firstProgressStage(values ...ProgressStage) ProgressStage {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
