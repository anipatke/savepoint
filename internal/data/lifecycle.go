package data

import (
	"fmt"
	"strings"
)

const LegacyTaskStatusTodo ColumnType = "todo"
const LegacyTaskStatusComplete ColumnType = "complete"
const LegacyTaskStatusCompleted ColumnType = "completed"
const LegacyTaskStageImplementation ProgressStage = "implementation"
const LegacyComplexityTierSmall ComplexityTier = "small"
const LegacyComplexityTierMed ComplexityTier = "med"
const LegacyComplexityTierNormal ComplexityTier = "normal"
const LegacyComplexityTierModerate ComplexityTier = "moderate"
const LegacyComplexityTierLarge ComplexityTier = "large"

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
	TaskLifecycleStatusAlias   TaskLifecycleDiagnosticCode = "status_alias"
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

func IsLegacyTaskStatusAlias(value ColumnType) bool {
	_, ok := ResolveTaskStatusAlias(value)
	return ok
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

	if status, ok := ResolveTaskStatusAlias(metadata.Status); ok {
		diagnostics = append(diagnostics, TaskLifecycleDiagnostic{
			Code:    TaskLifecycleStatusAlias,
			Message: fmt.Sprintf("task uses non-canonical status %q; replace with %q", metadata.Status, status),
		})
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
	wordCount := ComplexityReasonWordCount(reason)
	if wordCount > MaxComplexityReasonWords {
		return fmt.Errorf("complexity_reason has %d words; maximum is %d words", wordCount, MaxComplexityReasonWords)
	}
	return nil
}

func HealTaskMetadataForProgress(task *Task) bool {
	changed := false
	if tier, ok := ResolveComplexityTierAlias(task.ComplexityTier); ok {
		task.ComplexityTier = tier
		changed = true
	}
	words := strings.Fields(task.ComplexityReason)
	if len(words) > MaxComplexityReasonWords {
		task.ComplexityReason = strings.Join(words[:MaxComplexityReasonWords], " ")
		changed = true
	}
	return changed
}

func ResolveComplexityTierAlias(tier ComplexityTier) (ComplexityTier, bool) {
	switch tier {
	case LegacyComplexityTierSmall:
		return ComplexityLow, true
	case LegacyComplexityTierMed, LegacyComplexityTierNormal, LegacyComplexityTierModerate:
		return ComplexityMedium, true
	case LegacyComplexityTierLarge:
		return ComplexityHigh, true
	default:
		return "", false
	}
}

func ComplexityReasonWordCount(reason string) int {
	return len(strings.Fields(reason))
}

type DefectLifecycleDiagnosticCode string

const (
	DefectLifecycleStatusAlias   DefectLifecycleDiagnosticCode = "status_alias"
	DefectLifecycleInvalidStatus DefectLifecycleDiagnosticCode = "invalid_status"
	DefectLifecycleMissingStage  DefectLifecycleDiagnosticCode = "missing_stage"
	DefectLifecycleInvalidStage  DefectLifecycleDiagnosticCode = "invalid_stage"
	DefectLifecycleStaleStage    DefectLifecycleDiagnosticCode = "stale_stage"
)

type DefectLifecycleDiagnostic struct {
	Code    DefectLifecycleDiagnosticCode
	Message string
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

// DiagnoseDefectLifecycle reports every condition that
// NormalizeDefectLifecycleForLoad heals silently, from raw frontmatter values.
func DiagnoseDefectLifecycle(status DefectStatus, stage ProgressStage) []DefectLifecycleDiagnostic {
	var diagnostics []DefectLifecycleDiagnostic

	if alias, ok := ResolveDefectStatusAlias(status); ok {
		diagnostics = append(diagnostics, DefectLifecycleDiagnostic{
			Code:    DefectLifecycleStatusAlias,
			Message: fmt.Sprintf("defect uses non-canonical status %q; replace with %q", status, alias),
		})
	} else if status != "" && !IsCanonicalDefectStatus(status) {
		diagnostics = append(diagnostics, DefectLifecycleDiagnostic{
			Code:    DefectLifecycleInvalidStatus,
			Message: fmt.Sprintf("defect status invalid %q; use open, in_progress, or resolved (loads as open)", status),
		})
		return diagnostics
	}

	if NormalizeDefectStatusForLoad(status) == DefectInProgress {
		if stage == "" {
			diagnostics = append(diagnostics, DefectLifecycleDiagnostic{
				Code:    DefectLifecycleMissingStage,
				Message: "defect stage is required when status is in_progress (loads as build)",
			})
			return diagnostics
		}
		if !IsCanonicalStage(stage) {
			diagnostics = append(diagnostics, DefectLifecycleDiagnostic{
				Code:    DefectLifecycleInvalidStage,
				Message: fmt.Sprintf("defect stage invalid %q; use build, test, or audit (loads as build)", stage),
			})
		}
		return diagnostics
	}

	if stage != "" {
		diagnostics = append(diagnostics, DefectLifecycleDiagnostic{
			Code:    DefectLifecycleStaleStage,
			Message: fmt.Sprintf("defect stage %q is only valid when status is in_progress (ignored on load)", stage),
		})
	}

	return diagnostics
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
