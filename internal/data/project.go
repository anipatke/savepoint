package data

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// V2Index is the identity-keyed project index for a V2 project's Release,
// Objective, and Task records, built by LoadV2Index. Records are looked up by
// their global ID, never by directory-derived numbering or path. Release
// membership is derived from Objective.Release; ReleaseObjectives is a
// reverse lookup only and is never read as an authored membership list. Task
// ownership in ObjectiveTasks is explicit: it comes from each Task's
// objective field, not from the directory a Task file happens to live in, so
// moving a Task file retains both its identity and its ownership.
//
// V2Index construction validates structure only: duplicate identity, path
// safety, ownership, and reference-graph diagnostics. It does not infer
// Check clearance, owner acceptance, or dependency satisfaction; those gates
// belong to the data gate resolvers.
type V2Index struct {
	Releases   map[string]*ReleaseV2
	Objectives map[string]*ObjectiveV2
	Tasks      map[string]*TaskV2
	Checks     map[string]*CheckV2
	// Issues holds every durable follow-up record keyed by its global I-###
	// identity, the identity that survives repair, recheck, and reopening.
	Issues map[string]*IssueV2
	// ObjectiveTasks maps an Objective ID to the sorted IDs of the Tasks
	// that declare it as their owner.
	ObjectiveTasks map[string][]string
	// ReleaseObjectives maps a Release ID to the sorted IDs of Objectives
	// whose records declare that Release. It is derived during load; Release
	// records do not contain a second, mutable membership list.
	ReleaseObjectives map[string][]string
	// DuplicateObjectiveRanks records non-fatal duplicate ranks within a
	// Goal and priority group. Members are sorted by Objective ID.
	DuplicateObjectiveRanks []DuplicateObjectiveRankV2
	// ObjectivesWithoutGoal contains the sorted IDs of live Objectives whose
	// records omit release. These records remain loadable so callers can report
	// and repair the missing Goal reference.
	ObjectivesWithoutGoal []string
	// TaskIssues maps a Task ID to the sorted IDs of the Issues that name it
	// as carrying their repair.
	TaskIssues map[string][]string
	// CheckIssues maps a Check ID to the sorted IDs of the Issues that name
	// it. Both maps are built from the Issue side: an Issue may legally name
	// a recheck or proof Check that never recorded it, while the pairing rule
	// guarantees every Check-side link appears there too.
	CheckIssues map[string][]string
	// ScopeChecks maps a Task, Objective, or Release ID to the IDs of the Checks that
	// name it as their scope target, in recorded (ascending C-### ID) order.
	ScopeChecks map[string][]string
	// LatestCheck maps a Task, Objective, or Release ID to the most recently recorded
	// Check ID for that target — the last entry of ScopeChecks[id].
	LatestCheck map[string]string
}

// HasLiveGoal reports whether the project has a Goal that is not an archived
// legacy completion.
func (index *V2Index) HasLiveGoal() bool {
	if index == nil {
		return false
	}
	for _, goal := range index.Releases {
		if goal != nil && goal.LegacyCompletion == nil {
			return true
		}
	}
	return false
}

// LoadV2Index discovers and validates the optional Release records and the V2
// Objective and Task records under root (a .savepoint directory), returning
// one authoritative, identity-keyed index. It fails closed on the first
// structural diagnostic encountered — duplicate/malformed IDs, unsafe paths,
// missing ownership, dangling Release references, missing dependency targets,
// self-dependencies, or Task/Objective cycles — in deterministic ID order, so
// repeated runs against the same project report the same diagnostic first.
func LoadV2Index(root string) (*V2Index, error) {
	objectives, tasks, err := DiscoverV2Records(root)
	if err != nil {
		return nil, err
	}

	releases, err := DiscoverV2Releases(root)
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
		Releases:          releases,
		Objectives:        objectives,
		Tasks:             tasks,
		Checks:            checks,
		Issues:            issues,
		ObjectiveTasks:    map[string][]string{},
		ReleaseObjectives: map[string][]string{},
		TaskIssues:        map[string][]string{},
		CheckIssues:       map[string][]string{},
		ScopeChecks:       map[string][]string{},
		LatestCheck:       map[string]string{},
	}

	for id := range releases {
		index.ReleaseObjectives[id] = nil
	}

	if err := indexReleaseObjectives(index); err != nil {
		return nil, err
	}
	index.DuplicateObjectiveRanks = duplicateObjectiveRankFacts(index)

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

	if err := validateIssueEscalationTargets(index); err != nil {
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

// indexReleaseObjectives validates typed Objective release references and
// derives the only reverse membership view. Any present reference must use the
// R-### identity vocabulary and resolve to a discovered Release.
func indexReleaseObjectives(index *V2Index) error {
	for _, id := range slices.Sorted(maps.Keys(index.Objectives)) {
		objective := index.Objectives[id]
		if objective.Release == "" {
			index.ObjectivesWithoutGoal = append(index.ObjectivesWithoutGoal, objective.ID)
			continue
		}

		ref := objective.Release
		if !matchesV2Identity(ref, 'R') {
			return fmt.Errorf("%w: %s: objective %s release %q must match R-###", ErrV2InvalidReleaseReference, objective.Source.Path, objective.ID, ref)
		}

		releaseID := string(ref)
		if _, ok := index.Releases[releaseID]; !ok {
			return fmt.Errorf("%w: %s: objective %s references missing release %s", ErrV2MissingRelease, objective.Source.Path, objective.ID, releaseID)
		}
		index.ReleaseObjectives[releaseID] = append(index.ReleaseObjectives[releaseID], objective.ID)
	}
	return nil
}

// DuplicateObjectiveRankV2 is the deterministic diagnostic fact for two or
// more Objectives sharing a rank within one Goal and priority group.
type DuplicateObjectiveRankV2 struct {
	GoalID       string
	Priority     ObjectivePriority
	Rank         int
	ObjectiveIDs []string
}

// OrderedObjectiveIDsForGoal returns the Goal's Objective IDs in planning
// order: priority, ranked before unranked, ascending rank, then Objective ID.
// Goal membership remains derived from each Objective's release field.
func OrderedObjectiveIDsForGoal(index *V2Index, goalID string) []string {
	if index == nil {
		return nil
	}
	ids := append([]string(nil), index.ReleaseObjectives[goalID]...)
	slices.SortFunc(ids, func(a, b string) int {
		objectiveA := index.Objectives[a]
		objectiveB := index.Objectives[b]
		priorityA, rankA := ObjectivePriorityMedium, 0
		priorityB, rankB := ObjectivePriorityMedium, 0
		if objectiveA != nil {
			priorityA, rankA = objectiveA.Priority, objectiveA.Rank
		}
		if objectiveB != nil {
			priorityB, rankB = objectiveB.Priority, objectiveB.Rank
		}
		if orderA, orderB := objectivePriorityOrder(priorityA), objectivePriorityOrder(priorityB); orderA != orderB {
			return orderA - orderB
		}
		rankedA, rankedB := rankA > 0, rankB > 0
		if rankedA != rankedB {
			if rankedA {
				return -1
			}
			return 1
		}
		if rankedA && rankA != rankB {
			return rankA - rankB
		}
		return strings.Compare(a, b)
	})
	return ids
}

func duplicateObjectiveRankFacts(index *V2Index) []DuplicateObjectiveRankV2 {
	type groupKey struct {
		goalID   string
		priority ObjectivePriority
		rank     int
	}
	groups := make(map[groupKey][]string)
	for _, goalID := range slices.Sorted(maps.Keys(index.ReleaseObjectives)) {
		for _, objectiveID := range index.ReleaseObjectives[goalID] {
			objective := index.Objectives[objectiveID]
			if objective == nil || objective.Rank <= 0 {
				continue
			}
			key := groupKey{goalID: goalID, priority: objective.Priority, rank: objective.Rank}
			groups[key] = append(groups[key], objectiveID)
		}
	}

	facts := make([]DuplicateObjectiveRankV2, 0)
	for key, objectiveIDs := range groups {
		if len(objectiveIDs) < 2 {
			continue
		}
		slices.Sort(objectiveIDs)
		facts = append(facts, DuplicateObjectiveRankV2{
			GoalID:       key.goalID,
			Priority:     key.priority,
			Rank:         key.rank,
			ObjectiveIDs: objectiveIDs,
		})
	}
	slices.SortFunc(facts, func(a, b DuplicateObjectiveRankV2) int {
		if goalOrder := strings.Compare(a.GoalID, b.GoalID); goalOrder != 0 {
			return goalOrder
		}
		if priorityOrder := objectivePriorityOrder(a.Priority) - objectivePriorityOrder(b.Priority); priorityOrder != 0 {
			return priorityOrder
		}
		return a.Rank - b.Rank
	})
	return facts
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

// validateIssueEscalationTargets validates every escalated_to reference by
// global ID: the named Objective must exist. Unlike duplicate_of, no cycle is
// possible here — Issue and Objective are different record families — so this
// is an existence check only. It walks Issue IDs in sorted order so a project
// with more than one dangling reference reports the same one first on every
// run.
func validateIssueEscalationTargets(index *V2Index) error {
	for _, id := range slices.Sorted(maps.Keys(index.Issues)) {
		issue := index.Issues[id]
		if issue.EscalatedTo == "" {
			continue
		}
		if _, ok := index.Objectives[issue.EscalatedTo]; !ok {
			return fmt.Errorf("%w: %s: issue %s escalates to missing objective %s", ErrV2IssueMissingEscalationTarget, issue.Source.Path, id, issue.EscalatedTo)
		}
	}
	return nil
}

// validateEvidenceReferences resolves every Check reference named in Task,
// Objective, and Release evidence blocks — last_check, freshness.check,
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
	for _, id := range slices.Sorted(maps.Keys(index.Releases)) {
		release := index.Releases[id]
		if err := checkEvidenceReferences(index, release.Source.Path, "release", release.ID, release.Evidence); err != nil {
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
	digits := strings.TrimPrefix(id, "C-")
	digits = strings.TrimLeft(digits, "0")
	if digits == "" {
		return "0"
	}
	return digits
}

func checkScopeTargetExists(index *V2Index, scope CheckScope) bool {
	if index == nil {
		return false
	}
	switch scope.Kind {
	case CheckScopeTask:
		_, ok := index.Tasks[scope.ID]
		return ok
	case CheckScopeObjective:
		_, ok := index.Objectives[scope.ID]
		return ok
	case CheckScopeRelease:
		_, ok := index.Releases[scope.ID]
		return ok
	default:
		return false
	}
}

// validateCheckSupersedesChains validates, for every Check with a
// supersedes reference: the target Check exists, it shares the same scope,
// no two Checks supersede the same target, and the resulting graph has no
// cycle. It also requires each scope's numeric Check order to be one linear
// supersession chain, so LatestCheck cannot select a head unrelated to the
// explicit supersedes history. ids is the sorted Check ID order indexChecks
// already computed, so diagnostics and cycle discovery stay deterministic
// across runs.
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

	for _, targetID := range slices.Sorted(maps.Keys(index.ScopeChecks)) {
		checkIDs := index.ScopeChecks[targetID]
		first := index.Checks[checkIDs[0]]
		if first.Supersedes != "" {
			return fmt.Errorf("%w: checks for scope %s are not an ID-ordered chain: first check %s supersedes %s", ErrV2CheckSupersedesConflict, targetID, first.ID, first.Supersedes)
		}

		for i := 1; i < len(checkIDs); i++ {
			check := index.Checks[checkIDs[i]]
			want := checkIDs[i-1]
			if check.Supersedes != want {
				return fmt.Errorf("%w: checks for scope %s are not an ID-ordered chain: check %s supersedes %s, want %s", ErrV2CheckSupersedesConflict, targetID, check.ID, check.Supersedes, want)
			}
		}
	}

	return nil
}
