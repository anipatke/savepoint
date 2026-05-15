package data

import "fmt"

const LegacyTaskStatusTodo ColumnType = "todo"
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

type TaskLifecycleTransition struct {
	From TaskLifecycleState
	To   TaskLifecycleState
}

type TaskLifecycleDiagnosticCode string

const (
	TaskLifecycleMissingStatus TaskLifecycleDiagnosticCode = "missing_status"
	TaskLifecycleInvalidStatus TaskLifecycleDiagnosticCode = "invalid_status"
	TaskLifecycleLegacyPhase   TaskLifecycleDiagnosticCode = "legacy_phase"
	TaskLifecycleMissingStage  TaskLifecycleDiagnosticCode = "missing_stage"
	TaskLifecycleInvalidStage  TaskLifecycleDiagnosticCode = "invalid_stage"
	TaskLifecycleStaleStage    TaskLifecycleDiagnosticCode = "stale_stage"
)

type TaskLifecycleDiagnostic struct {
	Code    TaskLifecycleDiagnosticCode
	Message string
}

type TaskLifecycleDiagnosticInput struct {
	Metadata  TaskLifecycleMetadata
	HasStatus bool
	HasStage  bool
	HasPhase  bool
}

func CanonicalTaskStatuses() []ColumnType {
	return []ColumnType{ColumnPlanned, ColumnInProgress, ColumnDone}
}

func CanonicalTaskStages() []ProgressStage {
	return []ProgressStage{StageBuild, StageTest, StageAudit}
}

func ParseTaskLifecycle(metadata TaskLifecycleMetadata) (TaskLifecycleState, error) {
	rawStatus := firstTaskStatus(metadata.Column, metadata.Status)
	state := TaskLifecycleState{
		Status: NormalizeTaskStatusForLoad(rawStatus),
	}

	if state.Status == ColumnInProgress {
		state.Stage = firstProgressStage(metadata.Stage, metadata.Phase)
	}

	if err := validateLoadedTaskLifecycle(rawStatus, state, metadata.Phase); err != nil {
		return TaskLifecycleState{}, err
	}

	return state, nil
}

func NormalizeTaskStatusForLoad(value ColumnType) ColumnType {
	switch value {
	case "", LegacyTaskStatusTodo:
		return ColumnPlanned
	case ColumnPlanned, ColumnInProgress, ColumnDone:
		return value
	default:
		return value
	}
}

func IsLegacyTaskStatusAlias(value ColumnType) bool {
	_, ok := ResolveTaskStatusAlias(value)
	return ok
}

func ResolveTaskStatusAlias(value ColumnType) (ColumnType, bool) {
	switch value {
	case LegacyTaskStatusTodo:
		return ColumnPlanned, true
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

func ValidateTaskLifecycle(task *Task) error {
	if err := ValidateComplexity(task.ComplexityTier, task.ComplexityReason); err != nil {
		return err
	}

	return ValidateTaskLifecycleStateForWrite(TaskLifecycleState{Status: task.Column, Stage: task.Stage})
}

func ValidateTaskLifecycleStateForWrite(state TaskLifecycleState) error {
	if !IsCanonicalTaskStatus(state.Status) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", state.Status)
	}

	return validateCanonicalTaskStage(state)
}

func ValidateTaskLifecycleTransition(transition TaskLifecycleTransition) error {
	if transition.From.Status != "" || transition.From.Stage != "" {
		if err := ValidateTaskLifecycleStateForWrite(transition.From); err != nil {
			return fmt.Errorf("invalid source lifecycle: %w", err)
		}
	}
	return ValidateTaskLifecycleStateForWrite(transition.To)
}

func TaskLifecycleStateFromTask(task Task) TaskLifecycleState {
	return TaskLifecycleState{Status: task.Column, Stage: task.Stage}
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

func AdvanceTaskLifecycle(task *Task) (TaskLifecycleTransition, error) {
	from := normalizeTaskLifecycleStateForTransition(TaskLifecycleStateFromTask(*task))
	to, err := AdvanceTaskLifecycleState(from)
	if err != nil {
		return TaskLifecycleTransition{}, err
	}
	transition := TaskLifecycleTransition{From: from, To: to}
	if err := ValidateTaskLifecycleTransition(transition); err != nil {
		return TaskLifecycleTransition{}, err
	}
	applyTaskLifecycleState(task, to)
	return transition, nil
}

func RetreatTaskLifecycle(task *Task) (TaskLifecycleTransition, error) {
	from := normalizeTaskLifecycleStateForTransition(TaskLifecycleStateFromTask(*task))
	to, err := RetreatTaskLifecycleState(from)
	if err != nil {
		return TaskLifecycleTransition{}, err
	}
	transition := TaskLifecycleTransition{From: from, To: to}
	if err := ValidateTaskLifecycleTransition(transition); err != nil {
		return TaskLifecycleTransition{}, err
	}
	applyTaskLifecycleState(task, to)
	return transition, nil
}

func SameTaskLifecycleForTransition(a, b Task) bool {
	return normalizeTaskLifecycleStateForTransition(TaskLifecycleStateFromTask(a)) ==
		normalizeTaskLifecycleStateForTransition(TaskLifecycleStateFromTask(b))
}

func DiagnoseTaskLifecycle(input TaskLifecycleDiagnosticInput) []TaskLifecycleDiagnostic {
	metadata := input.Metadata
	var diagnostics []TaskLifecycleDiagnostic

	if !input.HasStatus || metadata.Status == "" {
		diagnostics = append(diagnostics, TaskLifecycleDiagnostic{
			Code:    TaskLifecycleMissingStatus,
			Message: "task missing required frontmatter field: status",
		})
		return diagnostics
	}

	if !IsCanonicalTaskStatus(metadata.Status) && !IsLegacyTaskStatusAlias(metadata.Status) {
		diagnostics = append(diagnostics, TaskLifecycleDiagnostic{
			Code:    TaskLifecycleInvalidStatus,
			Message: fmt.Sprintf("task status invalid: invalid task status %q; use planned, in_progress, or done", metadata.Status),
		})
		return diagnostics
	}

	if input.HasPhase {
		diagnostics = append(diagnostics, TaskLifecycleDiagnostic{
			Code:    TaskLifecycleLegacyPhase,
			Message: "task uses legacy frontmatter field phase; rename phase to stage",
		})
	}

	effectiveStage := firstProgressStage(metadata.Stage, metadata.Phase)
	if metadata.Status == ColumnInProgress {
		if effectiveStage == "" {
			diagnostics = append(diagnostics, TaskLifecycleDiagnostic{
				Code:    TaskLifecycleMissingStage,
				Message: "task stage is required when status is in_progress",
			})
			return diagnostics
		}
		if !IsCanonicalStage(effectiveStage) {
			diagnostics = append(diagnostics, TaskLifecycleDiagnostic{
				Code:    TaskLifecycleInvalidStage,
				Message: fmt.Sprintf("task stage invalid: invalid stage %q; use build, test, or audit", effectiveStage),
			})
		}
		return diagnostics
	}

	if input.HasStage && metadata.Stage != "" {
		diagnostics = append(diagnostics, TaskLifecycleDiagnostic{
			Code:    TaskLifecycleStaleStage,
			Message: fmt.Sprintf("task stage field %q is only valid when status is in_progress", metadata.Stage),
		})
	}

	return diagnostics
}

// ValidateComplexity checks that complexity_tier and complexity_reason are
// mutually consistent and within allowed bounds.
func ValidateComplexity(tier ComplexityTier, reason string) error {
	if tier == "" && reason == "" {
		return nil
	}
	if tier != "" && !IsCanonicalComplexityTier(tier) {
		return fmt.Errorf("invalid complexity_tier %q: use low, medium, high, or spike", tier)
	}
	if tier != "" && reason == "" {
		return fmt.Errorf("complexity_reason is required when complexity_tier is set")
	}
	if reason != "" && tier == "" {
		return fmt.Errorf("complexity_tier is required when complexity_reason is set")
	}
	if len(reason) > MaxComplexityReasonLen {
		return fmt.Errorf("complexity_reason exceeds %d characters", MaxComplexityReasonLen)
	}
	return nil
}

func validateDefectLifecycle(d *Defect) error {
	if d.Status == "" {
		d.Status = DefectOpen
	}
	if !IsCanonicalDefectStatus(d.Status) {
		return fmt.Errorf("invalid defect status %q: use open, in_progress, or resolved", d.Status)
	}
	if d.Status == DefectInProgress {
		if d.Stage == "" {
			return fmt.Errorf("stage is required when defect status is in_progress")
		}
		if !IsCanonicalStage(d.Stage) {
			return fmt.Errorf("invalid stage %q: use build, test, or audit", d.Stage)
		}
		return nil
	}
	if d.Stage != "" {
		return fmt.Errorf("stage field %q is only valid when defect status is in_progress", d.Stage)
	}
	return nil
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

func IsCanonicalColumn(value ColumnType) bool {
	return IsCanonicalTaskStatus(value)
}

func IsCanonicalStage(value ProgressStage) bool {
	switch value {
	case StageBuild, StageTest, StageAudit:
		return true
	default:
		return false
	}
}

func validateLoadedTaskLifecycle(rawStatus ColumnType, state TaskLifecycleState, phase ProgressStage) error {
	if rawStatus != "" && !IsLegacyTaskStatusAlias(rawStatus) && !IsCanonicalTaskStatus(rawStatus) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", rawStatus)
	}
	if state.Status == ColumnInProgress {
		return validateCanonicalTaskStage(state)
	}
	if phase != "" && !IsCanonicalStage(phase) && !IsLegacyTaskStageAlias(phase) {
		return fmt.Errorf("invalid legacy phase %q: use build, test, or audit, or remove 'phase'", phase)
	}
	return nil
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

func applyTaskLifecycleState(task *Task, state TaskLifecycleState) {
	task.Column = state.Status
	task.Stage = state.Stage
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

func IsCanonicalComplexityTier(tier ComplexityTier) bool {
	switch tier {
	case ComplexityLow, ComplexityMedium, ComplexityHigh, ComplexitySpike:
		return true
	default:
		return false
	}
}
