package doctor

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/migrate"
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

// CheckRouter validates router.md: valid state name, release/epic directories exist.
// epicFilter, if non-empty, skips directory checks when the router epic doesn't match.
func CheckRouter(root, epicFilter string, overrides ...DoctorDependencies) error {
	deps := doctorDependencies(overrides)
	routerPath := filepath.Join(root, "router.md")
	raw, err := os.ReadFile(routerPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("router.md not found: %w", data.ErrConfigNotFound)
	}
	if err != nil {
		return fmt.Errorf("router.md unreadable: %w", err)
	}

	state, err := deps.RouterReader.ReadState(string(raw))
	if err != nil {
		return fmt.Errorf("router.md invalid state block: %w", err)
	}

	if epicFilter != "" && state.Epic != epicFilter {
		return nil
	}

	if state.Release != "" && state.Release != "none" {
		releasePath := filepath.Join(root, "releases", state.Release)
		if _, err := os.Stat(releasePath); os.IsNotExist(err) {
			return fmt.Errorf("router.md release %q directory not found", state.Release)
		}
	}

	if state.Epic != "" && state.Epic != "none" {
		if state.Release == "" || state.Release == "none" {
			return fmt.Errorf("router.md has epic %q but no release", state.Epic)
		}
		epicPath := filepath.Join(root, "releases", state.Release, "epics", state.Epic)
		if _, err := os.Stat(epicPath); os.IsNotExist(err) {
			return fmt.Errorf("router.md epic %q directory not found", state.Epic)
		}
	}

	return nil
}

// migrationOperationRepair is the suggestion for an incomplete migration
// operation: doctor names the recovery command but never runs it, because an
// incomplete migration is finished by the migration command, not by a doctor
// repair.
const migrationOperationRepair = "Run `savepoint migrate --recover` to resume the named operation"

// CheckMigration reports an incomplete migration operation under root's
// project directory as a named diagnostic, naming the operation, the paths
// not yet verified, and the recovery command — the same read-only detector
// upgrade-assets and the board's write commands consult at their write
// boundary. root is the project's .savepoint directory, as every other check
// in this file expects; the project directory PendingOperation expects is one
// level up. A project with no operation directory, or one whose operation
// already completed and removed itself, reports no problem.
func CheckMigration(root string) []Problem {
	report, err := migrate.PendingOperation(filepath.Dir(root))
	if err != nil {
		return []Problem{{
			File:    root,
			Message: fmt.Sprintf("[migrate-multiple-operations] %v", err),
			Repair:  "Resolve or remove the extra directories under .savepoint/.migration/ so only one operation remains, then run `savepoint migrate --recover`",
		}}
	}
	if report == nil {
		return nil
	}
	return []Problem{{
		File:    root,
		Message: fmt.Sprintf("[migrate-operation-incomplete] %s", report.RecoveryGuidance()),
		Repair:  migrationOperationRepair,
	}}
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

// CheckProject validates schema, record, identity, path, and reference-graph
// diagnostics by loading root through the shared data-layer project loader
// rather than re-deriving V2 schema, record, or graph rules here. A V1
// project (no explicit schema_version, or an absent config.yml) reports no
// problems from this check; V1 structural diagnostics remain CheckStructure's
// job. A V2 project's structural diagnostic is reported with a stable name
// plus the record path and identity context the loader already carries.
// Doctor only reads: the project loader never writes to the project.
func CheckProject(root string, overrides ...DoctorDependencies) []Problem {
	_, problems := loadProjectChecks(root, doctorDependencies(overrides))
	return problems
}

func loadProjectChecks(root string, deps DoctorDependencies) (*data.Project, []Problem) {
	project, err := deps.ProjectLoader.Load(root)
	if err != nil {
		name := v2DiagnosticName(err)
		return nil, []Problem{{
			File:    root,
			Message: fmt.Sprintf("[%s] %v", name, err),
			Repair:  V2ProblemRepair(name),
		}}
	}
	if project == nil {
		return nil, []Problem{{
			File:    root,
			Message: "[v2-project-error] project loader returned no project",
			Repair:  "Review the V2 project diagnostic and fix the reported record",
		}}
	}
	if project.SchemaVersion == data.SchemaVersionV2 {
		return project, v2ConsistencyProblems(project.V2)
	}
	return project, nil
}

// CheckReleaseReadiness reports Release-level readiness findings from the
// same indexed records and canonical gate resolver used by other V2
// consumers. Structural load errors remain CheckProject's responsibility, so
// this check returns no duplicate finding when the index cannot be built.
func CheckReleaseReadiness(root string, overrides ...DoctorDependencies) []Problem {
	project, _ := loadProjectChecks(root, doctorDependencies(overrides))
	return releaseDiagnosticsForProject(project).Problems
}

type releaseDiagnostics struct {
	Problems []Problem
	Notes    []string
}

func releaseDiagnosticsForProject(project *data.Project) releaseDiagnostics {
	if project == nil || project.SchemaVersion != data.SchemaVersionV2 || project.V2 == nil || len(project.V2.Releases) == 0 {
		return releaseDiagnostics{}
	}
	return releaseDiagnosticsForIndex(project.Root, project.V2)
}

// releaseDiagnosticsForIndex is the V2-only form used by the live doctor
// runtime. It accepts the already loaded index rather than schema-dispatching
// through data.LoadProject, keeping legacy discovery out of ordinary health
// checks while the compatibility helper above remains available to historical
// doctor tests.
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
		repair = fmt.Sprintf("Record a Release-scoped Check and current freshness evidence for %s; doctor does not create evidence", release.ID)
	case data.GateBlockClearanceNeedsWork:
		name = "v2-release-clearance-needs-work"
		detail = fmt.Sprintf("Release evidence needs work: %s", blocker.Detail)
		repair = fmt.Sprintf("Resolve the findings from the Release Check for %s, then record a fresh CLEAR Check and freshness assessment", release.ID)
	case data.GateBlockClearanceStale:
		name = "v2-release-clearance-stale"
		detail = fmt.Sprintf("Release evidence is stale: %s", blocker.Detail)
		repair = fmt.Sprintf("Record a current freshness assessment for the latest Release Check on %s", release.ID)
	case data.GateBlockClearanceUnknown:
		name = "v2-release-clearance-unknown"
		detail = fmt.Sprintf("Release evidence is unknown: %s", blocker.Detail)
		repair = fmt.Sprintf("Record freshness evidence assessed by an independent checker for the latest Release Check on %s", release.ID)
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
	if kind == data.ObjectiveConsistencyDoneWithoutClearance {
		return HealthMissingEvidence
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

// v2DiagnosticName maps a data.LoadProject error to the stable diagnostic
// name doctor reports it under. Every V2 structural sentinel in
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

// IssuePostureReport computes a V2 project's Issue backlog counts directly
// from the loaded index at report time. It returns nil for a V1 project or
// one that fails to load — CheckProject already reports load failures
// separately, and Issue posture is advisory only: it never contributes to
// DiagnosticReport.HasProblems.
func IssuePostureReport(root string) *IssuePosture {
	project, err := data.LoadProject(root)
	if err != nil || project.SchemaVersion != data.SchemaVersionV2 {
		return nil
	}
	return issuePostureForIndex(project.V2)
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

// CheckStructure validates release/epic/task structure and YAML across the project.
// epicFilter, if non-empty, restricts checks to matching epics.
func CheckStructure(root string, epicFilter string, overrides ...DoctorDependencies) []Problem {
	deps := doctorDependencies(overrides)
	var problems []Problem

	releasesPath := filepath.Join(root, "releases")
	if _, err := os.Stat(releasesPath); os.IsNotExist(err) {
		problems = append(problems, Problem{File: releasesPath, Message: "releases directory not found"})
		return problems
	}

	releases, err := deps.Discoverer.ListReleases(root)
	if err != nil {
		problems = append(problems, Problem{File: releasesPath, Message: fmt.Sprintf("listing releases: %v", err)})
		return problems
	}

	if len(releases) == 0 {
		problems = append(problems, Problem{File: releasesPath, Message: "no release directories found"})
		return problems
	}

	for _, release := range releases {
		checkReleasePRD(release.Path, release.ID, deps.Parser, &problems)

		epics, err := deps.Discoverer.ListEpics(root, release.ID)
		if err != nil {
			problems = append(problems, Problem{
				File:    filepath.Join(release.Path, "epics"),
				Message: fmt.Sprintf("listing epics in release %q: %v", release.ID, err),
			})
			continue
		}

		for _, epic := range epics {
			if epicFilter != "" && epic.ID != epicFilter && !strings.HasPrefix(epic.ID, epicFilter) {
				continue
			}

			checkEpicDetail(epic.Path, epic.ID, deps.Parser, &problems)

			tasks, err := deps.Discoverer.ListTasks(root, release.ID, epic.ID)
			if err != nil {
				problems = append(problems, Problem{
					File:    filepath.Join(epic.Path, "tasks"),
					Message: fmt.Sprintf("listing tasks in epic %q: %v", epic.ID, err),
				})
				continue
			}

			for _, task := range tasks {
				checkTaskFile(task.Path, deps.Parser, &problems)
			}
		}
	}

	return problems
}

func checkReleasePRD(releasePath string, releaseID string, parser taskParser, problems *[]Problem) {
	prdPath := filepath.Join(releasePath, releaseID+"-PRD.md")
	raw, err := os.ReadFile(prdPath)
	if os.IsNotExist(err) {
		*problems = append(*problems, Problem{File: prdPath, Message: "release PRD file not found"})
		return
	}
	if err != nil {
		*problems = append(*problems, Problem{File: prdPath, Message: fmt.Sprintf("unreadable: %v", err)})
		return
	}
	validateFrontmatter(prdPath, string(raw), parser, problems)
}

func checkEpicDetail(epicPath string, epicID string, parser taskParser, problems *[]Problem) {
	prefix := extractPrefix(epicID)
	detailPath := filepath.Join(epicPath, prefix+"-Detail.md")
	raw, err := os.ReadFile(detailPath)
	if os.IsNotExist(err) {
		*problems = append(*problems, Problem{File: detailPath, Message: "epic detail file not found"})
		return
	}
	if err != nil {
		*problems = append(*problems, Problem{File: detailPath, Message: fmt.Sprintf("unreadable: %v", err)})
		return
	}
	content := string(raw)
	validateFrontmatter(detailPath, content, parser, problems)
	checkEpicStatus(content, detailPath, parser, problems)
}

// checkEpicStatus reports a non-canonical epic status that load-time
// normalization heals silently, reading the raw frontmatter value. A
// missing or empty status produces no problem.
func checkEpicStatus(content, path string, parser taskParser, problems *[]Problem) {
	fm, err := parser.ParseFrontmatter(content)
	if err != nil {
		return
	}
	status, _ := fm["status"].(string)
	for _, diagnostic := range data.DiagnoseEpicStatus(data.EpicStatus(status)) {
		*problems = append(*problems, Problem{File: path, Message: diagnostic.Message})
	}
}

func extractPrefix(epicID string) string {
	if idx := strings.IndexByte(epicID, '-'); idx != -1 {
		return epicID[:idx]
	}
	return epicID
}

func checkTaskFile(path string, parser taskParser, problems *[]Problem) {
	raw, err := os.ReadFile(path)
	if err != nil {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("unreadable: %v", err)})
		return
	}

	content := string(raw)
	fm, err := parser.ParseFrontmatter(content)
	if err != nil {
		line := extractYAMLLine(err)
		*problems = append(*problems, Problem{File: path, Line: line, Message: fmt.Sprintf("invalid frontmatter: %v", err)})
		return
	}

	checkRequiredString(fm, path, "id", problems)
	checkRequiredString(fm, path, "objective", problems)
	checkDependsOn(fm, path, problems)
	checkTaskLifecycle(fm, path, problems)
	checkComplexity(fm, path, problems)

	if !hasAcceptanceCriteria(content) {
		*problems = append(*problems, Problem{File: path, Message: "task missing ## Acceptance Criteria section"})
	}
}

func checkTaskLifecycle(fm map[string]any, path string, problems *[]Problem) {
	status, hasStatus := fm["status"].(string)
	stage, hasStage := fm["stage"].(string)
	phase, hasPhase := fm["phase"].(string)
	diagnostics := data.DiagnoseTaskLifecycle(data.TaskLifecycleDiagnosticInput{
		Metadata: data.TaskLifecycleMetadata{
			Status: data.ColumnType(status),
			Stage:  data.ProgressStage(stage),
			Phase:  data.ProgressStage(phase),
		},
		HasStatus: hasStatus,
		HasStage:  hasStage,
		HasPhase:  hasPhase,
	})
	for _, diagnostic := range diagnostics {
		*problems = append(*problems, Problem{File: path, Message: diagnostic.Message})
	}
}

func checkRequiredString(fm map[string]any, path, field string, problems *[]Problem) {
	val, ok := fm[field]
	if !ok {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("task missing required frontmatter field: %s", field)})
		return
	}
	s, ok := val.(string)
	if !ok || s == "" {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("task frontmatter field %q must be a non-empty string", field)})
	}
}

func checkDependsOn(fm map[string]any, path string, problems *[]Problem) {
	val, ok := fm["depends_on"]
	if !ok {
		return
	}
	switch val.(type) {
	case []any, []string:
	default:
		*problems = append(*problems, Problem{File: path, Message: "task frontmatter field depends_on must be a list"})
	}
}

func checkComplexity(fm map[string]any, path string, problems *[]Problem) {
	var tier data.ComplexityTier
	if v, ok := fm["complexity_tier"].(string); ok {
		tier = data.ComplexityTier(v)
	}
	reason, _ := fm["complexity_reason"].(string)
	if err := data.ValidateComplexity(tier, reason); err != nil {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("task complexity invalid: %v", err)})
	}
}

func hasAcceptanceCriteria(content string) bool {
	normalized := strings.ReplaceAll(content, "\r\n", "\n")
	idx := strings.Index(normalized, "## Acceptance Criteria")
	if idx == -1 {
		return false
	}
	section := normalized[idx+len("## Acceptance Criteria"):]
	if next := strings.Index(section, "\n## "); next != -1 {
		section = section[:next]
	}
	section = strings.TrimSpace(section)
	return section != ""
}

func validateFrontmatter(path, content string, parser taskParser, problems *[]Problem) {
	_, err := parser.ParseFrontmatter(content)
	if err != nil {
		line := extractYAMLLine(err)
		*problems = append(*problems, Problem{File: path, Line: line, Message: fmt.Sprintf("invalid frontmatter: %v", err)})
	}
}

// taskDep describes a parsed task's dependency information.
type taskDep struct {
	File      string
	ID        string
	Release   string
	Epic      string
	DependsOn []string
}

// CheckDependencies validates task dependency integrity:
// missing deps, duplicate IDs, and dependency cycles.
// epicFilter restricts checks to matching epics if non-empty.
func CheckDependencies(root string, epicFilter string, overrides ...DoctorDependencies) []Problem {
	deps := doctorDependencies(overrides)
	var problems []Problem

	releases, err := deps.Discoverer.ListReleases(root)
	if err != nil {
		problems = append(problems, Problem{Message: fmt.Sprintf("listing releases: %v", err)})
		return problems
	}

	var allTasks []taskDep
	idSet := make(map[string]string) // id -> first file seen
	epicsByRelease := make(map[string]map[string]string)

	for _, release := range releases {
		epics, err := deps.Discoverer.ListEpics(root, release.ID)
		if err != nil {
			continue
		}
		epicsByRelease[release.ID] = make(map[string]string, len(epics))
		for _, epic := range epics {
			epicsByRelease[release.ID][epic.ID] = readEpicStatus(epic.Path, epic.ID, deps.Parser)
			if epicFilter != "" && epic.ID != epicFilter && !strings.HasPrefix(epic.ID, epicFilter) {
				continue
			}
			tasks, err := deps.Discoverer.ListTasks(root, release.ID, epic.ID)
			if err != nil {
				continue
			}
			for _, t := range tasks {
				td := parseTaskDep(t.Path, deps.Parser)
				if td == nil {
					continue
				}
				td.Release = release.ID
				td.Epic = epic.ID
				allTasks = append(allTasks, *td)
				if existing, ok := idSet[td.ID]; ok {
					problems = append(problems, Problem{
						File:    td.File,
						Message: fmt.Sprintf("duplicate task ID %q (first seen in %s)", td.ID, existing),
					})
				} else {
					idSet[td.ID] = td.File
				}
			}
		}
	}

	// Check for missing dependencies and cycles
	graph := make(map[string][]string) // id -> list of dependencies
	idToFile := make(map[string]string)
	resolverTasks := make([]data.Task, 0, len(allTasks))

	for _, td := range allTasks {
		idToFile[td.ID] = td.File
		resolverTasks = append(resolverTasks, data.Task{
			ID:      td.ID,
			Release: td.Release,
			Epic:    td.Epic,
			Column:  data.ColumnPlanned,
		})
		graph[td.ID] = nil
	}

	for _, td := range allTasks {
		dependent := data.Task{ID: td.ID, Release: td.Release, Epic: td.Epic}
		for _, dep := range td.DependsOn {
			resolved := data.ResolveDependency(dep, dependent, resolverTasks, epicsByRelease[td.Release])
			switch resolved.Kind {
			case data.DependencyTask:
				graph[td.ID] = append(graph[td.ID], resolved.ID)
			case data.DependencyEpic:
			default:
				problems = append(problems, Problem{
					File:    td.File,
					Message: fmt.Sprintf("depends_on references non-existent task %q", dep),
				})
			}
		}
	}

	// Cycle detection using DFS
	cycleProblems := detectCycles(graph, idToFile)
	problems = append(problems, cycleProblems...)

	return problems
}

func readEpicStatus(epicPath, epicID string, parser taskParser) string {
	prefix := extractPrefix(epicID)
	raw, err := os.ReadFile(filepath.Join(epicPath, prefix+"-Detail.md"))
	if err != nil {
		return ""
	}
	fm, err := parser.ParseFrontmatter(string(raw))
	if err != nil {
		return ""
	}
	status, _ := fm["status"].(string)
	return status
}

func parseTaskDep(path string, parser taskParser) *taskDep {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	fm, err := parser.ParseFrontmatter(string(raw))
	if err != nil {
		return nil
	}
	id, _ := fm["id"].(string)
	if id == "" {
		return nil
	}
	var deps []string
	switch v := fm["depends_on"].(type) {
	case []any:
		for _, d := range v {
			if s, ok := d.(string); ok {
				deps = append(deps, s)
			}
		}
	case []string:
		deps = v
	}
	return &taskDep{
		File:      path,
		ID:        id,
		DependsOn: deps,
	}
}

// detectCycles runs DFS on the dependency graph and returns cycle problems.
// Uses a path stack to accurately reconstruct cycle paths (avoids parent-map
// overwrite issues that produced inaccurate paths).
func detectCycles(graph map[string][]string, idToFile map[string]string) []Problem {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make(map[string]int)
	path := make([]string, 0)

	for id := range graph {
		color[id] = white
	}

	var problems []Problem

	var dfs func(id string)
	dfs = func(id string) {
		color[id] = gray
		path = append(path, id)
		for _, dep := range graph[id] {
			switch color[dep] {
			case white:
				dfs(dep)
			case gray:
				cycleStart := -1
				for i, n := range path {
					if n == dep {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := path[cycleStart:]
					cyclePath := make([]string, 0, len(cycle))
					for _, cid := range cycle {
						if f, ok := idToFile[cid]; ok {
							cyclePath = append(cyclePath, f)
						} else {
							cyclePath = append(cyclePath, cid)
						}
					}
					problems = append(problems, Problem{
						Message: fmt.Sprintf("dependency cycle detected: %s", strings.Join(cyclePath, " → ")),
					})
				}
			}
		}
		path = path[:len(path)-1]
		color[id] = black
	}

	for id := range graph {
		if color[id] == white {
			dfs(id)
		}
	}
	return problems
}

// CheckAuditState finds audit proposal files without matching audit-pending state in the router.
func CheckAuditState(root string, overrides ...DoctorDependencies) []Problem {
	deps := doctorDependencies(overrides)
	var problems []Problem

	routerPath := filepath.Join(root, "router.md")
	raw, err := os.ReadFile(routerPath)
	if err != nil {
		return problems
	}

	state, err := deps.RouterReader.ReadState(string(raw))
	if err != nil {
		return problems
	}

	releases, err := deps.Discoverer.ListReleases(root)
	if err != nil {
		return problems
	}

	for _, release := range releases {
		epics, err := deps.Discoverer.ListEpics(root, release.ID)
		if err != nil {
			continue
		}
		for _, epic := range epics {
			prefix := extractPrefix(epic.ID)
			auditPath := filepath.Join(epic.Path, prefix+"-Audit.md")
			if _, err := os.Stat(auditPath); os.IsNotExist(err) {
				continue
			}
			if state.State != "audit-pending" || state.Epic != epic.ID {
				problems = append(problems, Problem{
					File:    auditPath,
					Message: fmt.Sprintf("audit proposal exists but router state is %q (epic: %q) — expected audit-pending for %q", state.State, state.Epic, epic.ID),
				})
			}
		}
	}

	return problems
}

// CheckOrphans finds tasks whose epic prefix in their ID does not match any existing epic directory.
func CheckOrphans(root string, overrides ...DoctorDependencies) []Problem {
	deps := doctorDependencies(overrides)
	var problems []Problem

	existingEpics := make(map[string]bool)
	releasesPath := filepath.Join(root, "releases")
	releaseDirs, err := deps.Discoverer.ListRootDirs(releasesPath)
	if err != nil {
		problems = append(problems, Problem{File: releasesPath, Message: fmt.Sprintf("listing releases: %v", err)})
		return problems
	}

	for _, release := range releaseDirs {
		epicsPath := filepath.Join(releasesPath, release, "epics")
		epics, err := deps.Discoverer.ListRootDirs(epicsPath)
		if err != nil {
			continue
		}
		for _, epic := range epics {
			existingEpics[epic] = true
		}
	}

	// Collect all tasks and check their epic references
	allReleases, err := deps.Discoverer.ListReleases(root)
	if err != nil {
		return problems
	}

	for _, release := range allReleases {
		epics, err := deps.Discoverer.ListEpics(root, release.ID)
		if err != nil {
			continue
		}
		for _, epic := range epics {
			tasks, err := deps.Discoverer.ListTasks(root, release.ID, epic.ID)
			if err != nil {
				continue
			}
			for _, t := range tasks {
				raw, err := os.ReadFile(t.Path)
				if err != nil {
					continue
				}
				fm, err := deps.Parser.ParseFrontmatter(string(raw))
				if err != nil {
					continue
				}
				id, _ := fm["id"].(string)
				if id == "" {
					continue
				}
				idx := strings.IndexByte(id, '/')
				if idx == -1 {
					continue
				}
				taskEpic := id[:idx]
				if !existingEpics[taskEpic] {
					problems = append(problems, Problem{
						File:    t.Path,
						Message: fmt.Sprintf("orphaned task: epic %q does not exist in any release — consider moving to .savepoint/orphans/", taskEpic),
					})
				}
			}
		}
	}

	return problems
}

// CheckDefects validates all defect files across the project:
// frontmatter validity, status/stage lifecycle, and reference format.
func CheckDefects(root string, overrides ...DoctorDependencies) []Problem {
	deps := doctorDependencies(overrides)
	var problems []Problem

	releases, err := deps.Discoverer.ListReleases(root)
	if err != nil {
		return problems
	}

	taskIDs := collectTaskIDs(root, deps)

	for _, release := range releases {
		defects, err := deps.Discoverer.ListDefects(root, release.ID)
		if err != nil {
			problems = append(problems, Problem{
				File:    filepath.Join(root, "releases", release.ID, "defects"),
				Message: fmt.Sprintf("listing defects in release %q: %v", release.ID, err),
			})
			continue
		}
		for _, d := range defects {
			checkDefectFile(d.Path, deps.Parser, taskIDs, &problems)
		}
	}

	return problems
}

func checkDefectFile(path string, parser taskParser, taskIDs map[string]bool, problems *[]Problem) {
	raw, err := os.ReadFile(path)
	if err != nil {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("unreadable: %v", err)})
		return
	}

	defect, err := parser.ParseDefectFile(path, string(raw))
	if err != nil {
		line := extractYAMLLine(err)
		*problems = append(*problems, Problem{File: path, Line: line, Message: fmt.Sprintf("defect parse error: %v", err)})
		return
	}

	checkDefectLifecycle(string(raw), path, parser, problems)

	if defect.ID == "" {
		*problems = append(*problems, Problem{File: path, Message: "defect missing required frontmatter field: id"})
	}
	if defect.Severity == "" {
		*problems = append(*problems, Problem{File: path, Message: "defect missing required frontmatter field: severity"})
	}

	checkDefectReference(path, defect.Reference, taskIDs, problems)
	checkDefectReference(path, defect.Introduced, taskIDs, problems)
}

// checkDefectLifecycle reports lifecycle metadata that load-time
// normalization heals silently, reading raw frontmatter because
// ParseDefectFile returns already-healed values.
func checkDefectLifecycle(content, path string, parser taskParser, problems *[]Problem) {
	fm, err := parser.ParseFrontmatter(content)
	if err != nil {
		return
	}
	status, _ := fm["status"].(string)
	stage, _ := fm["stage"].(string)
	diagnostics := data.DiagnoseDefectLifecycle(data.DefectStatus(status), data.ProgressStage(stage))
	for _, diagnostic := range diagnostics {
		*problems = append(*problems, Problem{File: path, Message: diagnostic.Message})
	}
}

// checkDefectReference validates a reference field that looks like a task ref (contains /).
func checkDefectReference(path, ref string, taskIDs map[string]bool, problems *[]Problem) {
	if ref == "" {
		return
	}
	idx := strings.IndexByte(ref, '/')
	if idx == -1 {
		return
	}
	if idx == 0 || idx == len(ref)-1 {
		*problems = append(*problems, Problem{
			File:    path,
			Message: fmt.Sprintf("defect reference %q has empty epic or task component", ref),
		})
		return
	}
	if len(taskIDs) > 0 && !taskIDs[ref] {
		*problems = append(*problems, Problem{
			File:    path,
			Message: fmt.Sprintf("defect reference %q does not match any known task ID", ref),
		})
	}
}

// CheckAuditRegister validates audit-register findings under the audit/ tree:
// structural parse failures, field problems that load-time normalization heals
// silently, and cross-record lifecycle and link validation against the
// discovered project. A root without an audit/ tree (or one without findings)
// produces no problems. Doctor only reports; it never edits audit files.
func CheckAuditRegister(root string, overrides ...DoctorDependencies) []Problem {
	deps := doctorDependencies(overrides)
	var problems []Problem

	findings, err := data.LoadAuditFindings(root)
	if err != nil {
		problems = append(problems, Problem{
			File:    filepath.Join(root, "audit", "findings"),
			Message: fmt.Sprintf("audit finding unloadable: %v", err),
			Repair:  auditFrontmatterRepair,
		})
		return problems
	}

	if _, err := data.LoadAuditRuns(root); err != nil {
		problems = append(problems, Problem{
			File:    filepath.Join(root, "audit", "runs"),
			Message: fmt.Sprintf("audit run unloadable: %v", err),
			Repair:  auditFrontmatterRepair,
		})
	}

	if len(findings) == 0 {
		return problems
	}

	for _, finding := range findings {
		checkAuditFindingFields(finding.Path, deps.Parser, &problems)
	}

	pathByID := make(map[string]string, len(findings))
	for _, finding := range findings {
		pathByID[finding.ID] = finding.Path
	}
	items := collectAuditWorkItems(root, deps)
	for _, v := range data.ValidateAuditFindings(findings, items) {
		problems = append(problems, Problem{
			File:    pathByID[v.FindingID],
			Message: v.Message,
			Repair:  AuditValidationRepair(v.Code),
		})
	}

	return problems
}

// checkAuditFindingFields re-reads a finding's raw frontmatter and reports every
// field problem that NormalizeFindingForLoad heals silently, mirroring
// checkDefectLifecycle: LoadAuditFindings returns already-healed values, so the
// raw parse recovers the original problems.
func checkAuditFindingFields(path string, parser taskParser, problems *[]Problem) {
	raw, err := os.ReadFile(path)
	if err != nil {
		*problems = append(*problems, Problem{File: path, Message: fmt.Sprintf("unreadable: %v", err)})
		return
	}
	finding, err := parser.ParseRawFindingFile(path, string(raw))
	if err != nil {
		*problems = append(*problems, Problem{
			File:    path,
			Message: fmt.Sprintf("audit finding unloadable: %v", err),
			Repair:  auditFrontmatterRepair,
		})
		return
	}
	for _, diagnostic := range data.DiagnoseFinding(finding, path) {
		*problems = append(*problems, Problem{
			File:    path,
			Message: diagnostic.Message,
			Repair:  AuditFindingRepair(diagnostic.Code),
		})
	}
}

// collectAuditWorkItems flattens discovery into the work-item ID sets that audit
// finding links resolve against.
func collectAuditWorkItems(root string, deps DoctorDependencies) data.AuditWorkItems {
	var items data.AuditWorkItems

	releases, err := deps.Discoverer.ListReleases(root)
	if err != nil {
		return items
	}
	for _, release := range releases {
		items.Releases = append(items.Releases, release.ID)

		if defects, err := deps.Discoverer.ListDefects(root, release.ID); err == nil {
			for _, defect := range defects {
				items.Defects = append(items.Defects, defect.ID)
			}
		}

		epics, err := deps.Discoverer.ListEpics(root, release.ID)
		if err != nil {
			continue
		}
		for _, epic := range epics {
			items.Epics = append(items.Epics, epic.ID)
			tasks, err := deps.Discoverer.ListTasks(root, release.ID, epic.ID)
			if err != nil {
				continue
			}
			for _, t := range tasks {
				if id := readTaskID(t.Path, deps.Parser); id != "" {
					items.Tasks = append(items.Tasks, id)
				}
			}
		}
	}

	return items
}

// collectTaskIDs reads all task frontmatter IDs across the project.
func collectTaskIDs(root string, deps DoctorDependencies) map[string]bool {
	ids := make(map[string]bool)

	releases, err := deps.Discoverer.ListReleases(root)
	if err != nil {
		return ids
	}

	for _, release := range releases {
		epics, err := deps.Discoverer.ListEpics(root, release.ID)
		if err != nil {
			continue
		}
		for _, epic := range epics {
			tasks, err := deps.Discoverer.ListTasks(root, release.ID, epic.ID)
			if err != nil {
				continue
			}
			for _, t := range tasks {
				if id := readTaskID(t.Path, deps.Parser); id != "" {
					ids[id] = true
				}
			}
		}
	}

	return ids
}

// readTaskID reads a task file's frontmatter id, returning "" for any file that
// cannot be read or parsed; structural problems are CheckStructure's job.
func readTaskID(path string, parser taskParser) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	fm, err := parser.ParseFrontmatter(string(raw))
	if err != nil {
		return ""
	}
	id, _ := fm["id"].(string)
	return id
}

func extractYAMLLine(err error) int {
	s := err.Error()
	const prefix = "yaml: line "
	if idx := strings.Index(s, prefix); idx != -1 {
		rest := s[idx+len(prefix):]
		if end := strings.IndexByte(rest, ':'); end != -1 {
			if line, err := strconv.Atoi(rest[:end]); err == nil {
				return line
			}
		}
	}
	return 0
}
