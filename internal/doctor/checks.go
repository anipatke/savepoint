package doctor

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"gopkg.in/yaml.v3"
)

// CheckConfig validates config.yml: exists, valid YAML, required fields present.
func CheckConfig(root string) error {
	configPath := filepath.Join(root, "config.yml")
	raw, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("config.yml not found: %w", data.ErrConfigNotFound)
	}
	if err != nil {
		return fmt.Errorf("config.yml unreadable: %w", err)
	}

	var fields map[string]any
	if err := yaml.Unmarshal(raw, &fields); err != nil {
		return fmt.Errorf("config.yml invalid YAML: %w", err)
	}

	if _, ok := fields["quality_gates"]; !ok {
		return fmt.Errorf("config.yml missing required field: quality_gates")
	}
	if _, ok := fields["theme"]; !ok {
		return fmt.Errorf("config.yml missing required field: theme")
	}

	return nil
}

// Problem describes a single issue found during a structure check. Repair, when
// set, is the typed repair suggestion for the problem; report formatting falls
// back to SuggestRepair message matching when it is empty. Category assigns the
// problem to one of the report's four health categories; a zero value defaults
// to HealthMalformedData, the category every structural check below already
// belongs to.
type Problem struct {
	File     string
	Line     int
	Message  string
	Repair   string
	Category HealthCategory
}

type releaseDiagnostics struct {
	Problems []Problem
	Notes    []string
}

// releaseDiagnosticsForIndex is the V2-only form used by the live doctor
// runtime. It accepts the already loaded index, keeping project loading and
// legacy discovery outside ordinary health checks.
func releaseDiagnosticsForIndex(root string, index *data.V2Index) releaseDiagnostics {
	if index == nil || len(index.Releases) == 0 {
		return releaseDiagnostics{}
	}

	var diagnostics releaseDiagnostics
	for _, releaseID := range slices.Sorted(maps.Keys(index.Releases)) {
		release := index.Releases[releaseID]
		if release.LegacyCompletion != nil {
			if problem := legacyCompletionProblem(root, release); problem != nil {
				diagnostics.Problems = append(diagnostics.Problems, *problem)
			} else if release.Status == data.ColumnDone {
				diagnostics.Notes = append(diagnostics.Notes, fmt.Sprintf(
					"release %s: historical completion is preserved in the legacy archive; it is historical evidence, not a current CLEAR Check",
					releaseID,
				))
			}
		}

		if release.Status != data.ColumnInProgress && release.Status != data.ColumnDone {
			continue
		}

		decision := data.ResolveReleaseCompletion(index, releaseID)
		if decision.Allowed {
			continue
		}
		clearance := data.ResolveClearance(index, releaseID)
		for _, blocker := range decision.Blockers {
			// A Goal may be active before its first Objective is planned. The
			// completion gate still blocks closing it without members.
			if release.Status == data.ColumnInProgress && blocker.Kind == data.GateBlockReleaseNoObjectives {
				continue
			}
			diagnostics.Problems = append(diagnostics.Problems, releaseBlockerProblem(release, blocker, clearance.Check))
		}
	}

	return diagnostics
}

func legacyCompletionProblem(root string, release *data.ReleaseV2) *Problem {
	archivePath, ok := resolveLegacyArchivePath(root, release.LegacyCompletion.ArchivePath)
	if !ok {
		problem := Problem{
			File:     release.Source.Path,
			Message:  fmt.Sprintf("[v2-release-legacy-dangling] release %s legacy completion archive path %q is not confined to the project", release.ID, release.LegacyCompletion.ArchivePath),
			Repair:   "Restore or correct the legacy_completion archive_path in the Release record; doctor does not rewrite historical evidence",
			Category: HealthMalformedData,
		}
		return &problem
	}
	info, err := os.Stat(archivePath)
	if err != nil || info.IsDir() {
		problem := Problem{
			File:     release.Source.Path,
			Message:  fmt.Sprintf("[v2-release-legacy-dangling] release %s legacy completion archive %q is missing", release.ID, release.LegacyCompletion.ArchivePath),
			Repair:   "Restore the archived legacy source or correct the legacy_completion archive_path in the Release record; doctor does not create evidence",
			Category: HealthMalformedData,
		}
		return &problem
	}
	return nil
}

func resolveLegacyArchivePath(root, reference string) (string, bool) {
	clean := filepath.Clean(filepath.FromSlash(reference))
	if reference == "" || filepath.IsAbs(reference) || clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", false
	}

	base := root
	rootName := filepath.Base(filepath.Clean(root))
	if clean == rootName || strings.HasPrefix(clean, rootName+string(filepath.Separator)) {
		base = filepath.Dir(root)
	}
	return filepath.Join(base, clean), true
}

func releaseBlockerProblem(release *data.ReleaseV2, blocker data.GateBlocker, latestCheck string) Problem {
	name := "v2-release-readiness"
	detail := blocker.Detail
	repair := "Review the Release's recorded readiness requirement and make the authoritative change manually; doctor never creates Release evidence"
	category := HealthMissingEvidence

	switch blocker.Kind {
	case data.GateBlockReleaseNoObjectives:
		name = "v2-release-no-objectives"
		detail = fmt.Sprintf("release %s has no member Objectives", release.ID)
		repair = fmt.Sprintf("Assign at least one Objective by adding release: %s to that Objective; doctor does not create membership", release.ID)
	case data.GateBlockReleaseObjectiveIncomplete:
		name = "v2-release-objective-incomplete"
		if blocker.Objective != "" {
			detail = fmt.Sprintf("member Objective %s is incomplete: %s", blocker.Objective, blocker.Detail)
		}
		repair = "Complete the named member Objective through its existing completion gate before treating the Release as ready; doctor does not change Objective records"
	case data.GateBlockClearanceMissing:
		name = "v2-release-clearance-missing"
		detail = fmt.Sprintf("Release evidence is missing: %s", blocker.Detail)
		repair = fmt.Sprintf("Record a Release-scoped Check for %s; doctor does not create evidence", release.ID)
	case data.GateBlockClearanceNeedsWork:
		name = "v2-release-clearance-needs-work"
		detail = fmt.Sprintf("Release evidence needs work: %s", blocker.Detail)
		repair = fmt.Sprintf("Resolve the findings from the Release Check for %s, then record a fresh CLEAR Check", release.ID)
	case data.GateBlockClearanceStale:
		name = "v2-release-clearance-stale"
		detail = fmt.Sprintf("Release evidence is stale: %s", blocker.Detail)
		repair = fmt.Sprintf("Run a new Release Check on %s; a freshness assessment marks the latest one stale", release.ID)
	case data.GateBlockClearanceUnknown:
		name = "v2-release-clearance-unknown"
		detail = fmt.Sprintf("Release evidence is unknown: %s", blocker.Detail)
		repair = fmt.Sprintf("Run a new Release Check on %s; a freshness assessment marks the latest one unknown", release.ID)
	case data.GateBlockCheckerAuthority:
		name = "v2-release-checker-authority"
		detail = fmt.Sprintf("Release evidence is unknown: %s", blocker.Detail)
		repair = fmt.Sprintf("Record the latest Release Check and freshness assessment with independent checker provenance for %s", release.ID)
	case data.GateBlockReleaseIssueUnresolved:
		name = "v2-release-issue-unresolved"
		detail = fmt.Sprintf("material blocker remains unresolved: %s", blocker.Detail)
		repair = "Resolve or explicitly except the named material Issue, then re-check the Release"
	case data.GateBlockOwnerAcceptance:
		if release.Evidence != nil && release.Evidence.OwnerValidation != nil && release.Evidence.OwnerValidation.AcceptedCheck != "" && release.Evidence.OwnerValidation.AcceptedCheck != latestCheck {
			name = "v2-release-owner-acceptance-stale"
			detail = fmt.Sprintf("owner acceptance is stale: it names %s, but the latest Release Check is %s", release.Evidence.OwnerValidation.AcceptedCheck, latestCheck)
		} else {
			name = "v2-release-owner-acceptance-missing"
			detail = fmt.Sprintf("owner acceptance is missing for current Release Check %s", latestCheck)
		}
		repair = fmt.Sprintf("Record owner acceptance for current Release Check %s in the Release record; doctor does not create acceptance", latestCheck)
	}

	return Problem{
		File:     release.Source.Path,
		Message:  fmt.Sprintf("[%s] release %s: %s", name, release.ID, detail),
		Repair:   repair,
		Category: category,
	}
}

// v2ConsistencyProblems reports every evaluation-level inconsistency
// data.InspectTaskConsistency finds in a successfully loaded V2 index —
// mismatches between a Task's recorded status and its recorded evidence that
// do not fail the load itself. Each Problem's File names the Task's source
// record so the reported record path and identity context point straight at
// the file to review; doctor only reports these, it never repairs or
// rewrites a Check, an evidence field, or a record status.
func v2ConsistencyProblems(index *data.V2Index) []Problem {
	var problems []Problem
	for _, diagnostic := range data.InspectTaskConsistency(index) {
		name := v2ConsistencyDiagnosticName(diagnostic.Kind)
		file := v2TaskSourcePath(index, diagnostic.Task)
		problems = append(problems, Problem{
			File:     file,
			Message:  fmt.Sprintf("[%s] task %s: %s", name, diagnostic.Task, diagnostic.Detail),
			Repair:   V2ConsistencyRepair(name),
			Category: v2ConsistencyCategory(diagnostic.Kind),
		})
	}
	for _, diagnostic := range data.InspectObjectiveConsistency(index) {
		name := v2ObjectiveConsistencyDiagnosticName(diagnostic.Kind)
		file := v2ObjectiveSourcePath(index, diagnostic.Objective)
		problems = append(problems, Problem{
			File:     file,
			Message:  fmt.Sprintf("[%s] objective %s: %s", name, diagnostic.Objective, diagnostic.Detail),
			Repair:   V2ConsistencyRepair(name),
			Category: v2ObjectiveConsistencyCategory(diagnostic.Kind),
		})
	}
	for _, diagnostic := range data.InspectIssueConsistency(index) {
		name := v2IssueConsistencyDiagnosticName(diagnostic.Kind)
		file := v2IssueSourcePath(index, diagnostic.Issue)
		problems = append(problems, Problem{
			File:     file,
			Message:  fmt.Sprintf("[%s] issue %s: %s", name, diagnostic.Issue, diagnostic.Detail),
			Repair:   V2ConsistencyRepair(name),
			Category: v2IssueConsistencyCategory(diagnostic.Kind),
		})
	}
	return problems
}

// v2ConsistencyCategory maps a data.ConsistencyDiagnosticKind to the health
// category doctor reports it under. A done Task without current clearance
// means work is claimed complete but its evidence has not caught up yet —
// nothing is broken, so it is missing evidence, not malformed data. The
// other kinds mean the record's own fields disagree with each other after a
// hand edit, which is malformed data.
func v2ConsistencyCategory(kind data.ConsistencyDiagnosticKind) HealthCategory {
	if kind == data.ConsistencyDoneWithoutClearance {
		return HealthMissingEvidence
	}
	return HealthMalformedData
}

// v2ObjectiveConsistencyCategory mirrors v2ConsistencyCategory for Objective
// diagnostics: a done Objective without current clearance is missing
// evidence, not malformed data.
func v2ObjectiveConsistencyCategory(kind data.ObjectiveConsistencyDiagnosticKind) HealthCategory {
	switch kind {
	case data.ObjectiveConsistencyDoneWithoutClearance:
		return HealthMissingEvidence
	case data.ObjectiveConsistencyPlannedWithStartedTask:
		// A status left behind by work that started without the board is
		// a warning to review, not broken data.
		return HealthPendingReview
	}
	return HealthMalformedData
}

// v2IssueConsistencyCategory mirrors v2ConsistencyCategory for Issue
// diagnostics: a verified Issue whose proof Check has been superseded is
// missing evidence — its proof is no longer current — not malformed data.
func v2IssueConsistencyCategory(kind data.IssueConsistencyDiagnosticKind) HealthCategory {
	if kind == data.IssueConsistencyProofSuperseded {
		return HealthMissingEvidence
	}
	return HealthMalformedData
}

// v2TaskSourcePath looks up taskID's source record path in index, falling
// back to the ID itself if the task is somehow absent — InspectTaskConsistency
// only ever names Tasks already present in the same index.
func v2TaskSourcePath(index *data.V2Index, taskID string) string {
	if task, ok := index.Tasks[taskID]; ok {
		return task.Source.Path
	}
	return taskID
}

// v2ObjectiveSourcePath looks up objectiveID's source record path in index,
// falling back to the ID itself if the objective is somehow absent —
// InspectObjectiveConsistency only ever names Objectives already present in
// the same index.
func v2ObjectiveSourcePath(index *data.V2Index, objectiveID string) string {
	if objective, ok := index.Objectives[objectiveID]; ok {
		return objective.Source.Path
	}
	return objectiveID
}

// v2IssueSourcePath looks up issueID's source record path in index, falling
// back to the ID itself if the issue is somehow absent — InspectIssueConsistency
// only ever names Issues already present in the same index.
func v2IssueSourcePath(index *data.V2Index, issueID string) string {
	if issue, ok := index.Issues[issueID]; ok {
		return issue.Source.Path
	}
	return issueID
}

// v2ObjectiveConsistencyDiagnosticName maps a
// data.ObjectiveConsistencyDiagnosticKind (InspectObjectiveConsistency's
// result) to the stable diagnostic name doctor reports it under. Every kind
// InspectObjectiveConsistency defines has a name here so a project's
// diagnostic name never changes between doctor runs.
func v2ObjectiveConsistencyDiagnosticName(kind data.ObjectiveConsistencyDiagnosticKind) string {
	switch kind {
	case data.ObjectiveConsistencyDoneWithoutClearance:
		return "v2-objective-done-without-clearance"
	case data.ObjectiveConsistencyIncompleteTask:
		return "v2-objective-done-with-incomplete-task"
	case data.ObjectiveConsistencyPlannedWithStartedTask:
		return "v2-objective-planned-with-started-task"
	default:
		return "v2-objective-evidence-inconsistency"
	}
}

// v2IssueConsistencyDiagnosticName maps a data.IssueConsistencyDiagnosticKind
// (InspectIssueConsistency's result) to the stable diagnostic name doctor
// reports it under. Every kind InspectIssueConsistency defines has a name
// here so a project's diagnostic name never changes between doctor runs.
func v2IssueConsistencyDiagnosticName(kind data.IssueConsistencyDiagnosticKind) string {
	switch kind {
	case data.IssueConsistencyProofSuperseded:
		return "v2-issue-verified-proof-superseded"
	default:
		return "v2-issue-evidence-inconsistency"
	}
}

// v2ConsistencyDiagnosticName maps a data.ConsistencyDiagnosticKind
// (InspectTaskConsistency's result) to the stable diagnostic name doctor
// reports it under. Every kind InspectTaskConsistency defines has a name
// here so a project's diagnostic name never changes between doctor runs.
func v2ConsistencyDiagnosticName(kind data.ConsistencyDiagnosticKind) string {
	switch kind {
	case data.ConsistencyDoneWithoutClearance:
		return "v2-done-without-clearance"
	case data.ConsistencyAcceptanceSuperseded:
		return "v2-acceptance-superseded"
	case data.ConsistencyEvidenceContradictsStatus:
		return "v2-evidence-contradicts-status"
	default:
		return "v2-evidence-inconsistency"
	}
}

// v2DiagnosticName maps V2 runtime errors to the stable diagnostic name
// doctor reports them under. Every V2 structural sentinel in
// internal/data/errors.go has a name here so a project's diagnostic name
// never changes between doctor runs.
func v2DiagnosticName(err error) string {
	switch {
	case errors.Is(err, data.ErrMalformedSchemaVersion):
		return "schema-version-malformed"
	case errors.Is(err, data.ErrUnsupportedSchemaVersion):
		return "schema-version-unsupported"
	case errors.Is(err, data.ErrV2MissingField):
		return "v2-missing-field"
	case errors.Is(err, data.ErrV2InvalidID):
		if strings.Contains(err.Error(), "release id") {
			return "v2-release-invalid-id"
		}
		return "v2-invalid-id"
	case errors.Is(err, data.ErrV2InvalidOwnership):
		return "v2-invalid-ownership"
	case errors.Is(err, data.ErrV2InvalidLifecycle):
		return "v2-invalid-lifecycle"
	case errors.Is(err, data.ErrV2InvalidDependency):
		return "v2-invalid-dependency"
	case errors.Is(err, data.ErrV2InvalidReleaseReference):
		return "v2-invalid-release-reference"
	case errors.Is(err, data.ErrV2DuplicateID):
		return "v2-duplicate-id"
	case errors.Is(err, data.ErrV2PathMismatch):
		return "v2-path-mismatch"
	case errors.Is(err, data.ErrV2UnsafePath):
		return "v2-unsafe-path"
	case errors.Is(err, data.ErrV2MissingOwner):
		return "v2-missing-owner"
	case errors.Is(err, data.ErrV2MissingRelease):
		return "v2-missing-release"
	case errors.Is(err, data.ErrV2ReleaseMissingSection):
		return "v2-release-missing-section"
	case errors.Is(err, data.ErrV2ReleaseLegacyMalformed):
		return "v2-release-legacy-malformed"
	case errors.Is(err, data.ErrV2MissingDependencyTarget):
		return "v2-missing-dependency-target"
	case errors.Is(err, data.ErrV2SelfDependency):
		return "v2-self-dependency"
	case errors.Is(err, data.ErrV2DependencyCycle):
		return "v2-dependency-cycle"
	case errors.Is(err, data.ErrV2Malformed):
		return "v2-record-malformed"
	case errors.Is(err, data.ErrV2CheckMalformed):
		return "v2-check-malformed"
	case errors.Is(err, data.ErrV2CheckMissingScopeTarget):
		if strings.Contains(err.Error(), "scope names missing release") {
			return "v2-check-missing-release-scope-target"
		}
		return "v2-check-missing-scope-target"
	case errors.Is(err, data.ErrV2CheckMissingReference):
		return "v2-check-missing-reference"
	case errors.Is(err, data.ErrV2CheckSupersedesConflict):
		return "v2-check-supersedes-conflict"
	case errors.Is(err, data.ErrV2EvidenceMalformed):
		return "v2-evidence-malformed"
	case errors.Is(err, data.ErrV2EvidenceMissingReference):
		return "v2-evidence-missing-reference"
	case errors.Is(err, data.ErrV2CheckImmutable):
		return "v2-check-immutable"
	case errors.Is(err, data.ErrV2IssueMissingDuplicateTarget):
		return "v2-issue-missing-duplicate-target"
	case errors.Is(err, data.ErrV2IssueMissingEscalationTarget):
		return "v2-issue-missing-escalation-target"
	case errors.Is(err, data.ErrV2IssueSelfDuplicate):
		return "v2-issue-self-duplicate"
	case errors.Is(err, data.ErrV2IssueDuplicateCycle):
		return "v2-issue-duplicate-cycle"
	case errors.Is(err, data.ErrV2IssueMissingLinkTarget):
		return "v2-issue-missing-link-target"
	case errors.Is(err, data.ErrV2IssueUnpairedCheckLink):
		return "v2-issue-unpaired-check-link"
	case errors.Is(err, data.ErrV2IssueResolutionRequired):
		return "v2-issue-resolution-required"
	case errors.Is(err, data.ErrV2IssueResolutionNotAllowed):
		return "v2-issue-resolution-not-allowed"
	case errors.Is(err, data.ErrV2IssueResolutionMissingProof):
		return "v2-issue-resolution-missing-proof"
	case errors.Is(err, data.ErrV2IssueResolutionUnusableProof):
		return "v2-issue-resolution-unusable-proof"
	case errors.Is(err, data.ErrV2IssueResolutionFieldMismatch):
		return "v2-issue-resolution-field-mismatch"
	case errors.Is(err, data.ErrV2IssueAlreadyExists):
		return "v2-issue-already-exists"
	case errors.Is(err, data.ErrV2IssueHistoryNotAppendOnly):
		return "v2-issue-history-not-append-only"
	case errors.Is(err, data.ErrV2IssueMalformed):
		return "v2-issue-malformed"
	default:
		return "v2-project-error"
	}
}

func (p Problem) Error() string {
	if p.Line > 0 {
		return fmt.Sprintf("%s:%d: %s", p.File, p.Line, p.Message)
	}
	if p.File != "" {
		return fmt.Sprintf("%s: %s", p.File, p.Message)
	}
	return p.Message
}

// IssuePosture summarizes a V2 project's Issue backlog by status and type,
// computed at report time. Doctor stores no separate summary, register, or
// cached total anywhere in the project: every count here is derived fresh
// from the loaded index.
type IssuePosture struct {
	StatusCounts map[data.IssueStatus]int
	TypeCounts   map[data.IssueType]int
	// Pending lists every open or in_progress Issue, in sorted ID order, so
	// the pending-review health category can name each one rather than only
	// its aggregate counts.
	Pending []IssueEntry
}

// IssueEntry names one Issue in the pending-review health category: its
// identity, type, and status.
type IssueEntry struct {
	ID     string
	Type   data.IssueType
	Status data.IssueStatus
}

// issuePostureForIndex derives the advisory Issue summary from an already
// loaded V2 index. The live doctor uses this form so it cannot invoke the V1
// schema dispatcher merely to count Issues.
func issuePostureForIndex(index *data.V2Index) *IssuePosture {
	if index == nil {
		return nil
	}
	posture := &IssuePosture{
		StatusCounts: index.IssueStatusCounts(),
		TypeCounts:   index.IssueTypeCounts(),
	}
	for _, id := range slices.Sorted(maps.Keys(index.Issues)) {
		issue := index.Issues[id]
		if issue.Status != data.IssueStatusOpen && issue.Status != data.IssueStatusInProgress {
			continue
		}
		posture.Pending = append(posture.Pending, IssueEntry{ID: issue.ID, Type: issue.Type, Status: issue.Status})
	}
	return posture
}
