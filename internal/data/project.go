package data

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

// Project is the schema-dispatched entry point for loading a .savepoint
// project. LoadProject detects schema_version exactly once, before any
// record discovery runs, then hands off to an isolated schema-specific
// loader.
type Project struct {
	Root          string
	SchemaVersion SchemaVersion
	V1            *Discover
	V2            *V2Index
}

// LoadProject detects the project schema from root's config.yml and
// dispatches to the schema-specific loader. root is a .savepoint directory
// path. Schema selection depends only on config.yml's schema_version field;
// it never consults package, release, or .upgrade-manifest.yml versions.
func LoadProject(root string) (*Project, error) {
	configPath := filepath.Join(root, "config.yml")
	version, err := ReadSchemaVersion(configPath)
	if err != nil {
		return nil, err
	}

	switch version {
	case SchemaVersionV2:
		return loadProjectV2(root)
	default:
		return loadProjectV1(root)
	}
}

// loadProjectV1 isolates transitional V1 discovery behind the schema
// dispatch boundary without changing Discover's existing behavior.
func loadProjectV1(root string) (*Project, error) {
	return &Project{Root: root, SchemaVersion: SchemaVersionV1, V1: NewDiscover()}, nil
}

// loadProjectV2 dispatches to the V2 identity-keyed project index. It
// deliberately fails closed on any structural diagnostic from LoadV2Index
// rather than falling back to V1 discovery.
func loadProjectV2(root string) (*Project, error) {
	index, err := LoadV2Index(root)
	if err != nil {
		return nil, err
	}
	return &Project{Root: root, SchemaVersion: SchemaVersionV2, V2: index}, nil
}

// V2Index is the identity-keyed project index for a V2 project's Objective
// and Task records, built by LoadV2Index. Objectives and Tasks are looked up
// by their global ID, never by directory-derived numbering or path. Task
// ownership in ObjectiveTasks is explicit: it comes from each Task's
// objective field, not from the directory a Task file happens to live in, so
// moving a Task file retains both its identity and its ownership.
//
// V2Index construction validates structure only: duplicate identity, path
// safety, ownership, and reference-graph diagnostics. It does not infer
// Check clearance, owner acceptance, or dependency satisfaction; those gates
// belong to E43.
type V2Index struct {
	Objectives map[string]*ObjectiveV2
	Tasks      map[string]*TaskV2
	Checks     map[string]*CheckV2
	// Issues holds every durable follow-up record keyed by its global I###
	// identity, the identity that survives repair, recheck, and reopening.
	Issues map[string]*IssueV2
	// ObjectiveTasks maps an Objective ID to the sorted IDs of the Tasks
	// that declare it as their owner.
	ObjectiveTasks map[string][]string
	// TaskIssues maps a Task ID to the sorted IDs of the Issues that name it
	// as carrying their repair.
	TaskIssues map[string][]string
	// CheckIssues maps a Check ID to the sorted IDs of the Issues that name
	// it. Both maps are built from the Issue side: an Issue may legally name
	// a recheck or proof Check that never recorded it, while the pairing rule
	// guarantees every Check-side link appears there too.
	CheckIssues map[string][]string
	// ScopeChecks maps a Task or Objective ID to the IDs of the Checks that
	// name it as their scope target, in recorded (ascending C### ID) order.
	ScopeChecks map[string][]string
	// LatestCheck maps a Task or Objective ID to the most recently recorded
	// Check ID for that target — the last entry of ScopeChecks[id].
	LatestCheck map[string]string
}

// LoadV2Index discovers and validates the V2 Objective and Task records
// under root (a .savepoint directory), returning one authoritative,
// identity-keyed index. It fails closed on the first structural diagnostic
// encountered — duplicate/malformed IDs, unsafe paths, missing ownership,
// missing dependency targets, self-dependencies, or Task/Objective cycles —
// in deterministic ID order, so repeated runs against the same project
// report the same diagnostic first.
func LoadV2Index(root string) (*V2Index, error) {
	objectives, tasks, err := DiscoverV2Records(root)
	if err != nil {
		return nil, err
	}

	checks, err := DiscoverV2Checks(root)
	if err != nil {
		return nil, err
	}

	issues, err := DiscoverV2Issues(root)
	if err != nil {
		return nil, err
	}

	index := &V2Index{
		Objectives:     objectives,
		Tasks:          tasks,
		Checks:         checks,
		Issues:         issues,
		ObjectiveTasks: map[string][]string{},
		TaskIssues:     map[string][]string{},
		CheckIssues:    map[string][]string{},
		ScopeChecks:    map[string][]string{},
		LatestCheck:    map[string]string{},
	}

	for _, id := range slices.Sorted(maps.Keys(tasks)) {
		task := tasks[id]
		if _, ok := objectives[task.Objective]; !ok {
			return nil, fmt.Errorf("%w: %s: task %s references missing objective %s", ErrV2MissingOwner, task.Source.Path, task.ID, task.Objective)
		}
		index.ObjectiveTasks[task.Objective] = append(index.ObjectiveTasks[task.Objective], task.ID)
	}

	if err := ValidateV2ReferenceGraphs(index); err != nil {
		return nil, err
	}

	if err := indexChecks(index); err != nil {
		return nil, err
	}

	if err := validateEvidenceReferences(index); err != nil {
		return nil, err
	}

	if err := validateIssueDuplicateGraph(index); err != nil {
		return nil, err
	}

	if err := validateIssueLinkTargets(index); err != nil {
		return nil, err
	}

	if err := validateCheckIssuePairing(index); err != nil {
		return nil, err
	}

	if err := validateIssueResolutionObligations(index); err != nil {
		return nil, err
	}

	indexIssueLinks(index)

	return index, nil
}

// validateIssueLinkTargets resolves every Issue reference that crosses record
// families — a Check's issues, and an Issue's tasks and checks — against the
// records actually discovered, so a link that names nothing fails the load
// closed instead of being carried as an unresolvable string. Checks are
// walked in recorded order and Issues in sorted ID order, so a project with
// several dangling references reports the same one first on every run.
func validateIssueLinkTargets(index *V2Index) error {
	for _, id := range sortedV2CheckIDs(index.Checks) {
		check := index.Checks[id]
		for _, issueID := range check.Issues {
			if _, ok := index.Issues[issueID]; !ok {
				return fmt.Errorf("%w: %s: check %s names missing issue %s", ErrV2IssueMissingLinkTarget, check.Source.Path, id, issueID)
			}
		}
	}

	for _, id := range slices.Sorted(maps.Keys(index.Issues)) {
		issue := index.Issues[id]
		for _, taskID := range issue.Tasks {
			if _, ok := index.Tasks[taskID]; !ok {
				return fmt.Errorf("%w: %s: issue %s names missing task %s", ErrV2IssueMissingLinkTarget, issue.Source.Path, id, taskID)
			}
		}
		for _, checkID := range issue.Checks {
			if _, ok := index.Checks[checkID]; !ok {
				return fmt.Errorf("%w: %s: issue %s names missing check %s", ErrV2IssueMissingLinkTarget, issue.Source.Path, id, checkID)
			}
		}
	}

	return nil
}

// validateCheckIssuePairing enforces the one authoritative direction of the
// Check-to-Issue relation. A Check is immutable, so its issues list is the
// record of which Issues that evaluation opened, and the Issue's mutable
// checks list must mirror it. A Check naming an Issue that does not name it
// back is a named diagnostic over both records rather than a reconciliation:
// the two sides disagree about what was observed, and only their authors know
// which one is wrong. The reverse asymmetry is legal — an Issue may name a
// recheck or proof Check that observed it without recording it. It runs after
// validateIssueLinkTargets, where every named Issue is proven to exist.
func validateCheckIssuePairing(index *V2Index) error {
	for _, id := range sortedV2CheckIDs(index.Checks) {
		check := index.Checks[id]
		for _, issueID := range check.Issues {
			if !slices.Contains(index.Issues[issueID].Checks, id) {
				return fmt.Errorf("%w: %s: check %s names issue %s, which does not name check %s in its checks", ErrV2IssueUnpairedCheckLink, check.Source.Path, id, issueID, id)
			}
		}
	}
	return nil
}

// indexIssueLinks builds the Task-to-Issue and Check-to-Issue maps once at
// load, so no consumer has to walk every Issue to answer what follow-up hangs
// off one record. It runs after link validation, where every reference has
// been proven to resolve.
func indexIssueLinks(index *V2Index) {
	for _, id := range slices.Sorted(maps.Keys(index.Issues)) {
		issue := index.Issues[id]
		for _, taskID := range issue.Tasks {
			index.TaskIssues[taskID] = appendUniqueIssueID(index.TaskIssues[taskID], id)
		}
		for _, checkID := range issue.Checks {
			index.CheckIssues[checkID] = appendUniqueIssueID(index.CheckIssues[checkID], id)
		}
	}
}

// appendUniqueIssueID appends id unless it is already the last entry. A record
// may legally repeat a reference, and a repeat lands next to its twin because
// one Issue's references are walked together and Issues are walked in
// ascending ID order — so the maps stay sorted and list each Issue once.
func appendUniqueIssueID(ids []string, id string) []string {
	if len(ids) > 0 && ids[len(ids)-1] == id {
		return ids
	}
	return append(ids, id)
}

// validateIssueDuplicateGraph validates every duplicate_of reference by
// global ID: the canonical Issue must exist and must be a different Issue,
// and the resulting graph must not close into a cycle where every Issue
// points at another and none is canonical. It walks IDs in sorted order so a
// project reports the same diagnostic on every run, and reuses the shared
// cycle walk so a duplicate chain is described exactly as a dependency cycle
// is.
func validateIssueDuplicateGraph(index *V2Index) error {
	ids := slices.Sorted(maps.Keys(index.Issues))
	edges := make(map[string][]string, len(ids))

	for _, id := range ids {
		issue := index.Issues[id]
		if issue.DuplicateOf == "" {
			continue
		}
		if issue.DuplicateOf == id {
			return fmt.Errorf("%w: %s: issue %s duplicates itself", ErrV2IssueSelfDuplicate, issue.Source.Path, id)
		}
		if _, ok := index.Issues[issue.DuplicateOf]; !ok {
			return fmt.Errorf("%w: %s: issue %s duplicates missing issue %s", ErrV2IssueMissingDuplicateTarget, issue.Source.Path, id, issue.DuplicateOf)
		}
		edges[id] = append(edges[id], issue.DuplicateOf)
	}

	if cycle := findV2Cycle(ids, edges); cycle != nil {
		return fmt.Errorf("%w: issue duplicate_of cycle %s", ErrV2IssueDuplicateCycle, describeV2Cycle(cycle, func(id string) string { return index.Issues[id].Source.Path }))
	}
	return nil
}

// validateEvidenceReferences resolves every Check reference named in Task
// and Objective evidence blocks — last_check, freshness.check,
// owner_validation.accepted_check, and exception.check — against
// index.Checks, so a load fails closed on a dangling reference exactly as
// Check supersedes references do. It walks Task IDs then Objective IDs in
// sorted order so diagnostics stay deterministic across runs.
func validateEvidenceReferences(index *V2Index) error {
	for _, id := range slices.Sorted(maps.Keys(index.Tasks)) {
		task := index.Tasks[id]
		if err := checkEvidenceReferences(index, task.Source.Path, "task", task.ID, task.Evidence); err != nil {
			return err
		}
	}
	for _, id := range slices.Sorted(maps.Keys(index.Objectives)) {
		objective := index.Objectives[id]
		if err := checkEvidenceReferences(index, objective.Source.Path, "objective", objective.ID, objective.Evidence); err != nil {
			return err
		}
	}
	return nil
}

// checkEvidenceReferences validates evidence's Check references in field
// declaration order, so a record with multiple dangling references always
// reports the same one first.
func checkEvidenceReferences(index *V2Index, path, recordKind, id string, evidence *Evidence) error {
	if evidence == nil {
		return nil
	}

	type reference struct {
		field string
		check string
	}
	refs := []reference{{"last_check", evidence.LastCheck}}
	if evidence.Freshness != nil {
		refs = append(refs, reference{"freshness.check", evidence.Freshness.Check})
	}
	if evidence.OwnerValidation != nil {
		refs = append(refs, reference{"owner_validation.accepted_check", evidence.OwnerValidation.AcceptedCheck})
	}
	if evidence.Exception != nil {
		refs = append(refs, reference{"exception.check", evidence.Exception.Check})
	}

	for _, ref := range refs {
		if ref.check == "" {
			continue
		}
		if _, ok := index.Checks[ref.check]; !ok {
			return fmt.Errorf("%w: %s: %s %s %s names missing check %s", ErrV2EvidenceMissingReference, path, recordKind, id, ref.field, ref.check)
		}
	}
	return nil
}

// indexChecks resolves each Check's scope target against the rest of index,
// groups Checks per scope target in recorded (ID) order, validates
// supersedes chains, and records each target's latest Check. It walks Check
// IDs in sorted order so a given project reports the same diagnostic first
// on every run, exactly as ValidateV2ReferenceGraphs does for Task/Objective
// dependencies.
func indexChecks(index *V2Index) error {
	ids := sortedV2CheckIDs(index.Checks)

	for _, id := range ids {
		check := index.Checks[id]
		if !checkScopeTargetExists(index, check.Scope) {
			return fmt.Errorf("%w: %s: check %s scope names missing %s %s", ErrV2CheckMissingScopeTarget, check.Source.Path, id, check.Scope.Kind, check.Scope.ID)
		}
		index.ScopeChecks[check.Scope.ID] = append(index.ScopeChecks[check.Scope.ID], id)
	}

	if err := validateCheckSupersedesChains(index, ids); err != nil {
		return err
	}

	for target, checkIDs := range index.ScopeChecks {
		index.LatestCheck[target] = checkIDs[len(checkIDs)-1]
	}

	return nil
}

// sortedV2CheckIDs orders global Check identities by their numeric suffix.
// Check IDs are allocated numerically, so lexical ordering would put C1000
// before C999 and make a newly recorded Check appear older than its
// predecessor. Comparing normalized decimal strings avoids an integer-width
// limit while retaining a lexical tie-breaker for unusual zero-padded IDs.
func sortedV2CheckIDs(checks map[string]*CheckV2) []string {
	return slices.SortedFunc(maps.Keys(checks), func(a, b string) int {
		return compareV2CheckIDs(a, b)
	})
}

func compareV2CheckIDs(a, b string) int {
	an := normalizedV2CheckNumber(a)
	bn := normalizedV2CheckNumber(b)
	if len(an) != len(bn) {
		if len(an) < len(bn) {
			return -1
		}
		return 1
	}
	if an < bn {
		return -1
	}
	if an > bn {
		return 1
	}
	return strings.Compare(a, b)
}

func normalizedV2CheckNumber(id string) string {
	digits := strings.TrimPrefix(id, "C")
	digits = strings.TrimLeft(digits, "0")
	if digits == "" {
		return "0"
	}
	return digits
}

func checkScopeTargetExists(index *V2Index, scope CheckScope) bool {
	switch scope.Kind {
	case CheckScopeTask:
		_, ok := index.Tasks[scope.ID]
		return ok
	case CheckScopeObjective:
		_, ok := index.Objectives[scope.ID]
		return ok
	default:
		return false
	}
}

// validateCheckSupersedesChains validates, for every Check with a
// supersedes reference: the target Check exists, it shares the same scope,
// no two Checks supersede the same target, and the resulting graph has no
// cycle. ids is the sorted Check ID order indexChecks already computed, so
// diagnostics and cycle discovery stay deterministic across runs.
func validateCheckSupersedesChains(index *V2Index, ids []string) error {
	edges := make(map[string][]string, len(ids))
	supersededBy := make(map[string]string, len(ids))

	for _, id := range ids {
		check := index.Checks[id]
		if check.Supersedes == "" {
			continue
		}

		target, ok := index.Checks[check.Supersedes]
		if !ok {
			return fmt.Errorf("%w: %s: check %s supersedes missing check %s", ErrV2CheckMissingReference, check.Source.Path, id, check.Supersedes)
		}
		if target.Scope != check.Scope {
			return fmt.Errorf("%w: %s: check %s supersedes %s with a different scope", ErrV2CheckSupersedesConflict, check.Source.Path, id, check.Supersedes)
		}
		if existing, ok := supersededBy[check.Supersedes]; ok {
			return fmt.Errorf("%w: %s: checks %s and %s both supersede %s", ErrV2CheckSupersedesConflict, check.Source.Path, existing, id, check.Supersedes)
		}
		supersededBy[check.Supersedes] = id
		edges[id] = append(edges[id], check.Supersedes)
	}

	if cycle := findV2Cycle(ids, edges); cycle != nil {
		return fmt.Errorf("%w: check supersedes cycle %s", ErrV2CheckSupersedesConflict, describeV2Cycle(cycle, func(id string) string { return index.Checks[id].Source.Path }))
	}
	return nil
}
