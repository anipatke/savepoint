package data

import (
	"go/build"
	"os"
	"strings"
	"testing"
)

// resolveSelectionWithTestGoal and resolveNextWithTestGoal put the older
// in-memory resolver fixtures into a valid Goal context. Tests for a missing
// router Goal call ResolveSelection or ResolveNext directly.
func resolveSelectionWithTestGoal(index *V2Index, router *RouterStateV2) (Selection, *SelectionDiagnostic) {
	addTestGoalContext(index, router, true)
	return ResolveSelection(index, router)
}

func resolveNextWithTestGoal(input NextInput) Next {
	addTestGoalContext(input.Index, input.Router, false)
	return ResolveNext(input)
}

func addTestGoalContext(index *V2Index, router *RouterStateV2, includeEmptySelection bool) {
	if index == nil || router == nil || router.Release != "" {
		return
	}
	if !includeEmptySelection && router.Objective == "" && router.Task == "" && router.Issue == "" {
		return
	}
	const goalID = "R-001"
	router.Release = goalID
	if index.Releases == nil {
		index.Releases = map[string]*ReleaseV2{}
	}
	if _, ok := index.Releases[goalID]; !ok {
		index.Releases[goalID] = &ReleaseV2{ID: goalID, Title: "Test Goal"}
	}
	for _, objective := range index.Objectives {
		if objective.Release == "" {
			objective.Release = goalID
		}
	}
}

// TestResolveSelection_exactMatch proves a router naming both an Objective
// and a Task that resolve, and agree on ownership, returns both records and
// no diagnostic.
func TestResolveSelection_exactMatch(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Ship it"}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Title: "Write the code", Objective: "O-001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic != nil {
		t.Fatalf("resolveSelectionWithTestGoal() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Objective == nil || selection.Objective.ID != "O-001" {
		t.Fatalf("Selection.Objective = %+v, want O-001", selection.Objective)
	}
	if selection.Task == nil || selection.Task.ID != "T-001" {
		t.Fatalf("Selection.Task = %+v, want T-001", selection.Task)
	}
}

func TestResolveSelection_issueAloneBecomesNextAndResolvedIssueIsStale(t *testing.T) {
	index := newV2TestIndex()
	index.Issues = map[string]*IssueV2{}
	openIssue := &IssueV2{ID: "I-042", Title: "Repair the parser", Status: IssueStatusOpen}
	index.Issues[openIssue.ID] = openIssue
	router := &RouterStateV2{State: RouterPhaseTask, Issue: openIssue.ID}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic != nil {
		t.Fatalf("resolveSelectionWithTestGoal() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Issue != openIssue || selection.Objective != nil || selection.Task != nil {
		t.Fatalf("resolveSelectionWithTestGoal() = %+v, want Issue I-042 alone", selection)
	}
	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextIssue || next.Issue != openIssue {
		t.Fatalf("resolveNextWithTestGoal() = %+v, want NextIssue carrying I-042", next)
	}

	resolved := *openIssue
	resolved.Status = IssueStatusResolved
	index.Issues[resolved.ID] = &resolved
	selection, diagnostic = resolveSelectionWithTestGoal(index, router)
	if selection.Issue != &resolved || diagnostic == nil || diagnostic.Kind != SelectionDone || diagnostic.RecordKind != SelectionRecordIssue || diagnostic.ID != resolved.ID {
		t.Fatalf("resolveSelectionWithTestGoal() = (%+v, %+v), want resolved Issue SelectionDone", selection, diagnostic)
	}
	next = resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextIssue || next.Issue != &resolved || next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionDone {
		t.Fatalf("resolveNextWithTestGoal() = %+v, want resolved Issue line with stale-selection diagnostic", next)
	}
}

func TestResolveSelection_unknownIssueNamesIssueAndKeepsValidTaskNext(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-014"] = &ObjectiveV2{ID: "O-014", Title: "Router Issue target", Status: ColumnInProgress}
	index.Tasks["T-028"] = &TaskV2{ID: "T-028", Title: "Copy the line", Objective: "O-014", Status: ColumnInProgress, Stage: StageBuild}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-014", Task: "T-028", Issue: "I-042"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if selection.Task == nil || selection.Task.ID != "T-028" || diagnostic == nil || diagnostic.Kind != SelectionNotFound || diagnostic.RecordKind != SelectionRecordIssue || diagnostic.ID != "I-042" {
		t.Fatalf("resolveSelectionWithTestGoal() = (%+v, %+v), want selected Task plus named Issue not-found diagnostic", selection, diagnostic)
	}
	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Task == nil || next.Task.ID != "T-028" || next.Kind == NextNothingSelected || next.SelectionDiagnostic == nil || next.SelectionDiagnostic.RecordKind != SelectionRecordIssue {
		t.Fatalf("resolveNextWithTestGoal() = %+v, want T-028's Next with Issue diagnostic", next)
	}
}

func TestResolveSelection_issueTravelsAsTaskContext(t *testing.T) {
	index := newV2TestIndex()
	index.Issues = map[string]*IssueV2{}
	index.Objectives["O-014"] = &ObjectiveV2{ID: "O-014", Title: "Router Issue target", Status: ColumnInProgress}
	index.Tasks["T-028"] = &TaskV2{ID: "T-028", Title: "Copy the line", Objective: "O-014", Status: ColumnInProgress, Stage: StageBuild}
	issue := &IssueV2{ID: "I-042", Title: "Repair the parser", Status: IssueStatusInProgress}
	index.Issues[issue.ID] = issue
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-014", Task: "T-028", Issue: issue.ID}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Task == nil || next.Task.ID != "T-028" || next.Issue != issue {
		t.Fatalf("resolveNextWithTestGoal() = %+v, want Task T-028 with Issue I-042 in context", next)
	}
}

// TestResolveSelection_nearMissIDNeverSubstituted proves a project holding
// both T-014 and T-140 resolves a router naming T-014 to exactly T-014: no
// numeric-proximity or prefix fallback ever selects the other record.
func TestResolveSelection_nearMissIDNeverSubstituted(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Ship it"}
	index.Tasks["T-014"] = &TaskV2{ID: "T-014", Title: "The real target", Objective: "O-001"}
	index.Tasks["T-140"] = &TaskV2{ID: "T-140", Title: "A similarly numbered task", Objective: "O-001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-014"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic != nil {
		t.Fatalf("resolveSelectionWithTestGoal() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Task == nil || selection.Task.ID != "T-014" {
		t.Fatalf("Selection.Task = %+v, want exactly T-014, never T-140", selection.Task)
	}
}

// TestResolveSelection_absentTask proves a router-named Task ID missing from
// the live index returns a typed not-found diagnostic naming it, rather than
// a partial selection.
func TestResolveSelection_absentTask(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Ship it"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-999"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic == nil {
		t.Fatal("resolveSelectionWithTestGoal() diagnostic = nil, want SelectionNotFound for the absent task")
	}
	if diagnostic.Kind != SelectionNotFound {
		t.Errorf("diagnostic.Kind = %q, want SelectionNotFound", diagnostic.Kind)
	}
	if diagnostic.RecordKind != SelectionRecordTask || diagnostic.ID != "T-999" {
		t.Errorf("diagnostic = %+v, want RecordKind task, ID T-999", diagnostic)
	}
	if selection.Objective != nil || selection.Task != nil {
		t.Errorf("Selection = %+v, want zero value alongside a diagnostic", selection)
	}
}

// TestResolveSelection_absentObjective proves the same for a router-named
// Objective ID missing from the live index.
func TestResolveSelection_absentObjective(t *testing.T) {
	index := newV2TestIndex()
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O-999"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic == nil {
		t.Fatal("resolveSelectionWithTestGoal() diagnostic = nil, want SelectionNotFound for the absent objective")
	}
	if diagnostic.Kind != SelectionNotFound {
		t.Errorf("diagnostic.Kind = %q, want SelectionNotFound", diagnostic.Kind)
	}
	if diagnostic.RecordKind != SelectionRecordObjective || diagnostic.ID != "O-999" {
		t.Errorf("diagnostic = %+v, want RecordKind objective, ID O-999", diagnostic)
	}
	if selection.Objective != nil || selection.Task != nil {
		t.Errorf("Selection = %+v, want zero value alongside a diagnostic", selection)
	}
}

// TestResolveSelection_mismatch proves a Task that resolves but is owned by
// a different Objective than the router names returns a typed mismatch
// diagnostic naming both, and does not substitute the Task's own Objective
// for the router's.
func TestResolveSelection_mismatch(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Router's guess"}
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Title: "Task's real owner"}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Title: "Moved task", Objective: "O-002"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic == nil {
		t.Fatal("resolveSelectionWithTestGoal() diagnostic = nil, want SelectionMismatch")
	}
	if diagnostic.Kind != SelectionMismatch {
		t.Errorf("diagnostic.Kind = %q, want SelectionMismatch", diagnostic.Kind)
	}
	if diagnostic.RouterObjective != "O-001" || diagnostic.Task != "T-001" || diagnostic.TaskObjective != "O-002" {
		t.Errorf("diagnostic = %+v, want RouterObjective O-001, Task T-001, TaskObjective O-002", diagnostic)
	}
	if selection.Objective != nil || selection.Task != nil {
		t.Errorf("Selection = %+v, want zero value alongside a diagnostic", selection)
	}
}

// TestResolveSelection_objectiveOnly proves a router naming an Objective
// with no Task resolves cleanly to that Objective alone, and is not a
// diagnostic.
func TestResolveSelection_objectiveOnly(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Plan me"}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O-001"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic != nil {
		t.Fatalf("resolveSelectionWithTestGoal() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Objective == nil || selection.Objective.ID != "O-001" {
		t.Fatalf("Selection.Objective = %+v, want O-001", selection.Objective)
	}
	if selection.Task != nil {
		t.Errorf("Selection.Task = %+v, want nil", selection.Task)
	}
}

func TestResolveSelection_doneTaskKeepsSelectionAndAddsDiagnostic(t *testing.T) {
	index := &V2Index{
		Objectives:     map[string]*ObjectiveV2{"O-001": {ID: "O-001", Title: "Ship it", Status: ColumnInProgress}},
		Tasks:          map[string]*TaskV2{"T-001": {ID: "T-001", Title: "Write it", Objective: "O-001", Status: ColumnDone}},
		ObjectiveTasks: map[string][]string{"O-001": {"T-001"}},
	}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic == nil || diagnostic.Kind != SelectionDone {
		t.Fatalf("resolveSelectionWithTestGoal() diagnostic = %+v, want SelectionDone", diagnostic)
	}
	if diagnostic.RecordKind != SelectionRecordTask || diagnostic.ID != "T-001" {
		t.Errorf("diagnostic = %+v, want finished Task T-001", diagnostic)
	}
	if selection.Objective == nil || selection.Objective.ID != "O-001" || selection.Task == nil || selection.Task.ID != "T-001" {
		t.Fatalf("resolveSelectionWithTestGoal() selection = %+v, want the original Objective O-001 / Task T-001", selection)
	}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind == NextNothingSelected || next.Objective == nil || next.Objective.ID != "O-001" || next.Task != nil {
		t.Errorf("resolveNextWithTestGoal() = %+v, want T-029's Objective rung for O-001 while retaining the warning", next)
	}
	if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionDone || next.SelectionDiagnostic.ID != "T-001" {
		t.Errorf("resolveNextWithTestGoal().SelectionDiagnostic = %+v, want SelectionDone for T-001", next.SelectionDiagnostic)
	}
}

func TestResolveSelection_doneObjectiveWithoutTaskAddsDiagnostic(t *testing.T) {
	index := &V2Index{
		Objectives: map[string]*ObjectiveV2{"O-002": {ID: "O-002", Title: "Ship it", Status: ColumnDone}},
	}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O-002"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic == nil || diagnostic.Kind != SelectionDone {
		t.Fatalf("resolveSelectionWithTestGoal() diagnostic = %+v, want SelectionDone", diagnostic)
	}
	if diagnostic.RecordKind != SelectionRecordObjective || diagnostic.ID != "O-002" {
		t.Errorf("diagnostic = %+v, want finished Objective O-002", diagnostic)
	}
	if selection.Objective == nil || selection.Objective.ID != "O-002" || selection.Task != nil {
		t.Fatalf("resolveSelectionWithTestGoal() selection = %+v, want the original Objective O-002 without a Task", selection)
	}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextPlanObjective || next.Objective == nil || next.Objective.ID != "O-002" {
		t.Errorf("resolveNextWithTestGoal() = %+v, want the same selected Objective rung as before the diagnostic", next)
	}
	if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionDone || next.SelectionDiagnostic.ID != "O-002" {
		t.Errorf("resolveNextWithTestGoal().SelectionDiagnostic = %+v, want SelectionDone for O-002", next.SelectionDiagnostic)
	}
}

func TestResolveSelection_notDoneSelectionsHaveNoDoneDiagnostic(t *testing.T) {
	index := &V2Index{
		Objectives: map[string]*ObjectiveV2{
			"O-003": {ID: "O-003", Title: "Active", Status: ColumnInProgress},
		},
		Tasks: map[string]*TaskV2{
			"T-003": {ID: "T-003", Title: "Build it", Objective: "O-003", Status: ColumnPlanned},
		},
	}
	for _, router := range []*RouterStateV2{
		{State: RouterPhaseDesign, Objective: "O-003"},
		{State: RouterPhaseTask, Objective: "O-003", Task: "T-003"},
	} {
		if selection, diagnostic := resolveSelectionWithTestGoal(index, router); diagnostic != nil {
			t.Errorf("resolveSelectionWithTestGoal(%+v) = (%+v, %+v), want no done diagnostic", router, selection, diagnostic)
		}
	}
}

func TestResolveSelection_releaseContextDiagnostics(t *testing.T) {
	tests := []struct {
		name        string
		configure   func(*V2Index)
		router      RouterStateV2
		wantKind    SelectionDiagnosticKind
		wantRelease bool
	}{
		{
			name: "valid release",
			configure: func(index *V2Index) {
				index.Releases["R-001"] = &ReleaseV2{ID: "R-001", Title: "First"}
			},
			router:      RouterStateV2{Release: "R-001"},
			wantRelease: true,
		},
		{
			name:     "missing release",
			router:   RouterStateV2{Release: "R-999"},
			wantKind: SelectionReleaseNotFound,
		},
		{
			name: "archived release",
			configure: func(index *V2Index) {
				index.Releases["R-001"] = &ReleaseV2{ID: "R-001", LegacyCompletion: &LegacyCompletionReference{SourcePath: "legacy.md"}}
			},
			router:      RouterStateV2{Release: "R-001"},
			wantKind:    SelectionReleaseArchived,
			wantRelease: true,
		},
		{
			name: "unassigned objective",
			configure: func(index *V2Index) {
				index.Releases["R-001"] = &ReleaseV2{ID: "R-001"}
				index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001"}
			},
			router:      RouterStateV2{Release: "R-001", Objective: "O-001"},
			wantKind:    SelectionObjectiveUnassigned,
			wantRelease: true,
		},
		{
			name: "objective belongs to another release",
			configure: func(index *V2Index) {
				index.Releases["R-001"] = &ReleaseV2{ID: "R-001"}
				index.Releases["R-002"] = &ReleaseV2{ID: "R-002"}
				index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Release: "R-002"}
			},
			router:      RouterStateV2{Release: "R-001", Objective: "O-001"},
			wantKind:    SelectionReleaseMismatch,
			wantRelease: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := newV2TestIndex()
			index.Releases = map[string]*ReleaseV2{}
			if tt.configure != nil {
				tt.configure(index)
			}

			selection, diagnostic := resolveSelectionWithTestGoal(index, &tt.router)
			if tt.wantKind == "" {
				if diagnostic != nil {
					t.Fatalf("diagnostic = %+v, want nil", diagnostic)
				}
			} else if diagnostic == nil || diagnostic.Kind != tt.wantKind {
				t.Fatalf("diagnostic = %+v, want %q", diagnostic, tt.wantKind)
			}
			if (selection.Release != nil) != tt.wantRelease {
				t.Fatalf("selection.Release = %+v, want present=%t", selection.Release, tt.wantRelease)
			}
		})
	}
}

func TestResolveSelection_releaseMismatchNeverSubstitutesObjective(t *testing.T) {
	index := newV2TestIndex()
	index.Releases = map[string]*ReleaseV2{
		"R-001": {ID: "R-001", Title: "Selected release"},
		"R-002": {ID: "R-002", Title: "Other release"},
	}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Same title", Release: "R-002"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, &RouterStateV2{Release: "R-001", Objective: "O-001"})
	if diagnostic == nil || diagnostic.Kind != SelectionReleaseMismatch {
		t.Fatalf("diagnostic = %+v, want release mismatch", diagnostic)
	}
	if selection.Objective != nil || selection.Release == nil || selection.Release.ID != "R-001" {
		t.Fatalf("selection = %+v, want only the exact selected Release and no substituted Objective", selection)
	}
}

func TestResolveSelection_missingGoalPreemptsOtherSelections(t *testing.T) {
	index := newV2TestIndex()
	index.Releases = map[string]*ReleaseV2{
		"R-001": {ID: "R-001", Title: "Available Goal"},
	}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Selected Objective", Release: "R-001"}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Title: "Selected Task", Objective: "O-001", Status: ColumnPlanned}
	index.Issues = map[string]*IssueV2{
		"I-001": {ID: "I-001", Title: "Selected Issue", Type: IssueTypeDefect, Status: IssueStatusOpen},
	}
	cases := []struct {
		name   string
		router RouterStateV2
	}{
		{name: "no record", router: RouterStateV2{State: RouterPhaseIdea}},
		{name: "objective and task", router: RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}},
		{name: "issue", router: RouterStateV2{State: RouterPhaseTask, Issue: "I-001"}},
		{name: "objective and issue", router: RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Issue: "I-001"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			selection, diagnostic := ResolveSelection(index, &tc.router)
			if diagnostic == nil || diagnostic.Kind != SelectionReleaseMissing || diagnostic.RecordKind != SelectionRecordRelease {
				t.Fatalf("ResolveSelection() diagnostic = %+v, want missing Goal diagnostic", diagnostic)
			}
			if selection != (Selection{}) {
				t.Fatalf("ResolveSelection() = %+v, want no record resolved before Goal selection", selection)
			}

			next := ResolveNext(NextInput{Index: index, Router: &tc.router})
			if next.Kind != NextNothingSelected || next.Release != nil || next.Objective != nil || next.Task != nil || next.Issue != nil || len(next.Issues) != 0 {
				t.Fatalf("ResolveNext() = %+v, want only the missing Goal diagnostic and no selected work", next)
			}
			if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionReleaseMissing {
				t.Errorf("ResolveNext().SelectionDiagnostic = %+v, want missing Goal diagnostic", next.SelectionDiagnostic)
			}
		})
	}
}

func TestResolveNext_carriesObjectivesWithoutGoalFacts(t *testing.T) {
	index := newV2TestIndex()
	index.Releases = map[string]*ReleaseV2{"R-001": {ID: "R-001", Title: "Available Goal"}}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Unassigned Objective"}
	index.ObjectivesWithoutGoal = []string{"O-001"}

	next := ResolveNext(NextInput{
		Index:  index,
		Router: &RouterStateV2{State: RouterPhaseDesign, Release: "R-001", Objective: "O-001"},
	})
	if len(next.ObjectivesWithoutGoal) != 1 || next.ObjectivesWithoutGoal[0] != "O-001" {
		t.Fatalf("ObjectivesWithoutGoal = %v, want [O-001] from the project index", next.ObjectivesWithoutGoal)
	}
}

// TestResolveSelection_noSelection proves a valid Goal context with no
// Objective and no Task resolves cleanly to no selection, not a diagnostic.
func TestResolveSelection_noSelection(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Title: "Unrelated"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic != nil {
		t.Fatalf("resolveSelectionWithTestGoal() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Objective != nil || selection.Task != nil {
		t.Errorf("Selection = %+v, want zero value for no selection", selection)
	}
}

// TestResolveSelection_readsNoFilesystem builds its V2Index directly in
// memory, with no discovery or file read anywhere in the call, proving
// ResolveSelection consults only the two values it is handed.
func TestResolveSelection_readsNoFilesystem(t *testing.T) {
	index := &V2Index{
		Objectives: map[string]*ObjectiveV2{
			"O-001": {ID: "O-001", Title: "In-memory only"},
		},
		Tasks: map[string]*TaskV2{
			"T-001": {ID: "T-001", Title: "In-memory only", Objective: "O-001"},
		},
	}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	selection, diagnostic := resolveSelectionWithTestGoal(index, router)
	if diagnostic != nil {
		t.Fatalf("resolveSelectionWithTestGoal() diagnostic = %+v, want nil", diagnostic)
	}
	if selection.Objective == nil || selection.Task == nil {
		t.Fatalf("Selection = %+v, want both records resolved from the in-memory index", selection)
	}
}

// TestResolveNext_replanOutranksExecution proves a replan flag reports
// NextReplan rather than NextExecute, even though the Task has no
// dependency and would otherwise be ready to start.
func TestResolveNext_replanOutranksExecution(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnPlanned,
		Evidence: &Evidence{Replan: &Replan{Reason: "scope changed", RecordedBy: Actor{Role: ActorRoleOwner, Session: "sess-1"}}},
	}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextReplan {
		t.Fatalf("Kind = %q, want replan", next.Kind)
	}
}

func TestResolveNext_selectedReleaseReplanOutranksReleaseExecution(t *testing.T) {
	index := newV2TestIndex()
	index.Releases = map[string]*ReleaseV2{"R-001": {ID: "R-001"}}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Release: "R-001"}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageBuild,
		Evidence: &Evidence{Replan: &Replan{Reason: "scope changed", RecordedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}}},
	}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.ReleaseObjectives = map[string][]string{"R-001": {"O-001"}}

	next := resolveNextWithTestGoal(NextInput{
		Index: index, Router: &RouterStateV2{Release: "R-001", Objective: "O-001", Task: "T-001"},
	})
	if next.Kind != NextReplan || next.Task == nil || next.Task.ID != "T-001" {
		t.Fatalf("next = %+v, want selected Release task replan", next)
	}
}

func TestResolveNext_selectedReleaseWithoutObjectiveDoesNotChooseMemberTask(t *testing.T) {
	index := newV2TestIndex()
	index.Releases = map[string]*ReleaseV2{
		"R-001": {ID: "R-001", Title: "Selected"},
		"R-002": {ID: "R-002", Title: "Unrelated"},
	}
	index.ReleaseObjectives = map[string][]string{}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Release: "R-001", Status: ColumnPlanned}
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Release: "R-002", Status: ColumnPlanned}
	index.Objectives["O-003"] = &ObjectiveV2{ID: "O-003", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Objective: "O-002", Status: ColumnPlanned}
	index.Tasks["T-003"] = &TaskV2{ID: "T-003", Objective: "O-003", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.ObjectiveTasks["O-002"] = []string{"T-002"}
	index.ObjectiveTasks["O-003"] = []string{"T-003"}
	index.ReleaseObjectives["R-001"] = []string{"O-001"}
	index.ReleaseObjectives["R-002"] = []string{"O-002"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: &RouterStateV2{Release: "R-001"}})
	if next.Kind != NextNothingSelected || next.Task != nil || next.Objective != nil {
		t.Fatalf("next = %+v, want nothing selected without substituting member work", next)
	}
	if next.Release == nil || next.Release.ID != "R-001" {
		t.Fatalf("next.Release = %+v, want R-001", next.Release)
	}
}

func TestResolveNext_selectedReleaseWithoutObjectiveDoesNotChooseActiveTask(t *testing.T) {
	index := newV2TestIndex()
	index.Releases = map[string]*ReleaseV2{"R-001": {ID: "R-001"}}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Release: "R-001", Status: ColumnInProgress}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageBuild}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.ReleaseObjectives = map[string][]string{"R-001": {"O-001"}}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: &RouterStateV2{Release: "R-001"}})
	if next.Kind != NextNothingSelected || next.Task != nil || next.Objective != nil {
		t.Fatalf("next = %+v, want nothing selected without substituting active member work", next)
	}
	if next.Release == nil || next.Release.ID != "R-001" {
		t.Fatalf("next.Release = %+v, want selected R-001 context retained", next.Release)
	}
}

func TestResolveNext_selectedReleaseWithCompleteMembersIsReadyRegardlessOfGoalCheck(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*V2Index)
	}{
		{
			name: "no Goal Check",
			configure: func(index *V2Index) {
				configured := releaseGateIndex()
				delete(configured.Checks, "C-002")
				delete(configured.ScopeChecks, "R-001")
				delete(configured.LatestCheck, "R-001")
				configured.Releases["R-001"].Evidence = nil
				*index = *configured
			},
		},
		{
			name: "NEEDS WORK Goal Check",
			configure: func(index *V2Index) {
				configured := releaseGateIndex()
				configured.Checks["C-002"].Result = CheckResultNeedsWork
				*index = *configured
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := newV2TestIndex()
			tt.configure(index)
			next := resolveNextWithTestGoal(NextInput{Index: index, Router: &RouterStateV2{Release: "R-001"}})
			if next.Kind != NextReleaseReady {
				t.Fatalf("next.Kind = %q, want %q; next = %+v", next.Kind, NextReleaseReady, next)
			}
			if next.Release == nil || next.Release.ID != "R-001" || next.GateDecision == nil || !next.GateDecision.Allowed || next.Clearance != nil {
				t.Fatalf("next = %+v, want ready Release decision with no Goal Check clearance", next)
			}
		})
	}
}

// TestResolveNext_replanOutranksCheckNeeded proves the same at the audit
// stage, where the Task would otherwise need a fresh Check.
func TestResolveNext_replanOutranksCheckNeeded(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Replan: &Replan{Reason: "scope changed", RecordedBy: Actor{Role: ActorRoleOwner, Session: "sess-1"}}},
	}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextReplan {
		t.Fatalf("Kind = %q, want replan (no check is recorded at all, which would otherwise be check_needed)", next.Kind)
	}
}

// TestResolveNext_taskDependencyBlocks proves an unsatisfied Task
// dependency reaches NextDependency carrying ResolveTaskStart's typed
// DependencyBlock.
func TestResolveNext_taskDependencyBlocks(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Objective: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnPlanned,
		DependsOn: []TaskDependencyV2{{Task: "T-002", Requires: TaskDependencyClear}},
	}
	index.ObjectiveTasks["O-001"] = []string{"T-001", "T-002"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextDependency {
		t.Fatalf("Kind = %q, want dependency", next.Kind)
	}
	if next.GateDecision == nil || len(next.GateDecision.Blockers) == 0 || next.GateDecision.Blockers[0].Kind != GateBlockDependency {
		t.Fatalf("GateDecision = %+v, want a GateBlockDependency blocker", next.GateDecision)
	}
	if dep := next.GateDecision.Blockers[0].Dependency; dep == nil || dep.Target != "T-002" {
		t.Errorf("Dependency block = %+v, want target T-002", dep)
	}
}

// TestResolveNext_objectiveDependencyBlocks proves an unsatisfied Objective
// dependency of the Task's owning Objective is distinguishable from a Task
// dependency: same NextDependency rung, but carrying
// ResolveObjectiveDependency's typed ObjectiveDependencyBlock instead.
func TestResolveNext_objectiveDependencyBlocks(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-002"] = &ObjectiveV2{ID: "O-002", Status: ColumnPlanned}
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned, DependsOn: []string{"O-002"}}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextDependency {
		t.Fatalf("Kind = %q, want dependency", next.Kind)
	}
	if next.GateDecision == nil || len(next.GateDecision.Blockers) == 0 || next.GateDecision.Blockers[0].Kind != GateBlockObjectiveDependency {
		t.Fatalf("GateDecision = %+v, want a GateBlockObjectiveDependency blocker", next.GateDecision)
	}
	if dep := next.GateDecision.Blockers[0].ObjectiveDependency; dep == nil || dep.Target != "O-002" {
		t.Errorf("ObjectiveDependency block = %+v, want target O-002", dep)
	}
}

// TestResolveNext_executeWhenTaskPlannedAndReady proves a planned Task with
// satisfied dependencies reaches NextExecute under executor authority.
func TestResolveNext_executeWhenTaskPlannedAndReady(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextExecute {
		t.Fatalf("Kind = %q, want execute", next.Kind)
	}
	if next.GateDecision == nil || !next.GateDecision.Allowed || next.GateDecision.Actor != ActorRoleExecutor {
		t.Errorf("GateDecision = %+v, want allowed under executor authority", next.GateDecision)
	}
}

// TestResolveNext_executeWhenTaskInProgress proves the same for a Task mid
// build/test, read through ResolveTaskAdvance rather than ResolveTaskStart.
func TestResolveNext_executeWhenTaskInProgress(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageBuild}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextExecute {
		t.Fatalf("Kind = %q, want execute", next.Kind)
	}
}

// TestResolveNext_checkNeeded_missing, _needsWork, _stale, and _unknown
// prove each of the four non-current clearance states reaches the same
// NextCheckNeeded rung while carrying its own distinguishable
// Clearance.State — no restatement of freshness logic in next.go, only the
// value ResolveClearance already computed.
func TestResolveNext_checkNeeded_missing(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded || next.Clearance == nil || next.Clearance.State != ClearanceMissing {
		t.Fatalf("next = %+v, want check_needed/missing", next)
	}
}

func TestResolveNext_checkNeeded_needsWork(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded || next.Clearance == nil || next.Clearance.State != ClearanceNeedsWork {
		t.Fatalf("next = %+v, want check_needed/needs_work", next)
	}
}

func TestResolveNext_checkNeeded_stale(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessStale, Check: "C-001", AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"}, Basis: "code changed after the check"}},
	}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded || next.Clearance == nil || next.Clearance.State != ClearanceStale {
		t.Fatalf("next = %+v, want check_needed/stale", next)
	}
}

func TestResolveNext_checkNeeded_unknown(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Freshness: &Freshness{State: FreshnessUnknown, Check: "C-001", Basis: "not reassessed"}},
	}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded || next.Clearance == nil || next.Clearance.State != ClearanceUnknown {
		t.Fatalf("next = %+v, want check_needed/unknown", next)
	}
}

// TestResolveNext_checkNeededOutranksOwnerValidation proves a Task needing
// both a fresh Check and owner validation reports the Check rung first: the
// switch inside ResolveTaskCompletion only ever adds an owner-acceptance
// blocker once clearance is already current.
func TestResolveNext_checkNeededOutranksOwnerValidation(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{OwnerValidation: &OwnerValidation{Required: true}},
	}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextCheckNeeded {
		t.Fatalf("Kind = %q, want check_needed (clearance is missing, so owner validation is not yet in view)", next.Kind)
	}
}

// TestResolveNext_ownerValidationRequiredAfterClearanceCurrent proves owner
// validation is reported only once clearance is current.
func TestResolveNext_ownerValidationRequiredAfterClearanceCurrent(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultClear)
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = mustCurrentTask("T-001", "C-001", &Evidence{OwnerValidation: &OwnerValidation{Required: true}})
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextOwnerValidationRequired {
		t.Fatalf("Kind = %q, want owner_validation_required", next.Kind)
	}
	if next.Clearance == nil || next.Clearance.State != ClearanceCurrent {
		t.Errorf("Clearance = %+v, want current", next.Clearance)
	}
}

// TestResolveNext_exceptionAllowedCompletionReportsExecuteNotClearance
// proves a completion allowed only by a recorded exception reaches
// NextExecute carrying AllowedByException and the Exception itself, and
// carries no Clearance — it is never presented as a CLEAR result.
func TestResolveNext_exceptionAllowedCompletionReportsExecuteNotClearance(t *testing.T) {
	index := newV2TestIndex()
	mustCheck(index, "C-001", "T-001", CheckResultNeedsWork)
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{
		ID: "T-001", Objective: "O-001", Status: ColumnInProgress, Stage: StageAudit,
		Evidence: &Evidence{Exception: &Exception{
			Requirements: []string{"R1"},
			Reason:       "known risk accepted",
			Owner:        "owner-1",
			Check:        "C-001",
		}},
	}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextExecute {
		t.Fatalf("Kind = %q, want execute (allowed by exception)", next.Kind)
	}
	if next.GateDecision == nil || !next.GateDecision.AllowedByException || next.GateDecision.Exception == nil {
		t.Fatalf("GateDecision = %+v, want AllowedByException with Exception carried", next.GateDecision)
	}
	if next.Clearance != nil {
		t.Errorf("Clearance = %+v, want nil (an exception-allowed completion is never presented as clearance)", next.Clearance)
	}
}

func TestResolveNext_selectedObjectiveDoesNotChooseUnrelatedReadyTask(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-013"] = &ObjectiveV2{ID: "O-013", Status: ColumnInProgress}
	index.Objectives["O-018"] = &ObjectiveV2{ID: "O-018", Title: "Selected objective", Status: ColumnInProgress}
	index.Tasks["T-006"] = &TaskV2{ID: "T-006", Objective: "O-013", Status: ColumnPlanned}
	index.ObjectiveTasks["O-013"] = []string{"T-006"}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O-018"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextPlanObjective {
		t.Fatalf("Kind = %q, want plan_objective for the selected Task-less Objective", next.Kind)
	}
	if next.Objective == nil || next.Objective.ID != "O-018" || next.Task != nil {
		t.Fatalf("next = %+v, want the selected O-018 and no Task", next)
	}
}

func TestResolveNext_selectedTasklessObjectiveWinsWhenOtherObjectiveHasActiveTask(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-013"] = &ObjectiveV2{ID: "O-013", Status: ColumnInProgress}
	index.Objectives["O-018"] = &ObjectiveV2{ID: "O-018", Title: "Selected objective", Status: ColumnInProgress}
	index.Tasks["T-006"] = &TaskV2{ID: "T-006", Objective: "O-013", Status: ColumnInProgress, Stage: StageBuild}
	index.ObjectiveTasks["O-013"] = []string{"T-006"}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O-018"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextPlanObjective {
		t.Fatalf("Kind = %q, want plan_objective for the selected Task-less Objective", next.Kind)
	}
	if next.Objective == nil || next.Objective.ID != "O-018" || next.Task != nil {
		t.Fatalf("next = %+v, want the selected O-018 and no Task", next)
	}
}

func TestResolveNext_doneSelectedTaskKeepsObjectiveAheadOfReleaseActiveTask(t *testing.T) {
	index := newV2TestIndex()
	index.Releases = map[string]*ReleaseV2{}
	index.ReleaseObjectives = map[string][]string{}
	index.Releases["R-006"] = &ReleaseV2{ID: "R-006", Title: "Selected release"}
	index.Objectives["O-013"] = &ObjectiveV2{ID: "O-013", Status: ColumnInProgress, Release: "R-006"}
	index.Objectives["O-018"] = &ObjectiveV2{ID: "O-018", Title: "Selected objective", Status: ColumnInProgress, Release: "R-006"}
	index.Tasks["T-006"] = &TaskV2{ID: "T-006", Objective: "O-013", Status: ColumnInProgress, Stage: StageBuild}
	index.Tasks["T-020"] = &TaskV2{ID: "T-020", Objective: "O-018", Status: ColumnDone}
	index.ObjectiveTasks["O-013"] = []string{"T-006"}
	index.ObjectiveTasks["O-018"] = []string{"T-020"}
	index.ReleaseObjectives["R-006"] = []string{"O-013", "O-018"}
	router := &RouterStateV2{
		State: RouterPhaseCheck, Release: "R-006", Objective: "O-018", Task: "T-020",
	}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextObjectiveIntegration {
		t.Fatalf("Kind = %q, want objective_integration for O-018", next.Kind)
	}
	if next.Objective == nil || next.Objective.ID != "O-018" || next.Task != nil {
		t.Fatalf("next = %+v, want O-018 integration with no selected Task", next)
	}
}

// TestResolveNext_objectiveIntegrationRung proves an Objective whose owned
// Tasks are all done, but whose own integration clearance is not current,
// reaches NextObjectiveIntegration carrying ResolveObjectiveCompletion's
// decision and the Objective's own Clearance.
func TestResolveNext_objectiveIntegrationRung(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextObjectiveIntegration {
		t.Fatalf("Kind = %q, want objective_integration", next.Kind)
	}
	if next.Objective == nil || next.Objective.ID != "O-001" {
		t.Fatalf("Objective = %+v, want O-001", next.Objective)
	}
	if next.GateDecision == nil || next.GateDecision.Allowed {
		t.Fatalf("GateDecision = %+v, want blocked (no objective check recorded)", next.GateDecision)
	}
	if next.Clearance == nil || next.Clearance.State != ClearanceMissing {
		t.Errorf("Clearance = %+v, want missing", next.Clearance)
	}
}

func TestResolveNext_selectedObjectiveReadyUsesCompletionResolver(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-001", nil)
	index.Objectives["O-001"].Status = ColumnInProgress
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: &RouterStateV2{State: RouterPhaseCheck, Objective: "O-001"}})
	if next.Kind != NextObjectiveReady {
		t.Fatalf("Kind = %q, want objective_ready", next.Kind)
	}
	if next.Objective == nil || next.Objective.ID != "O-001" || next.Task != nil {
		t.Fatalf("next = %+v, want O-001 ready for owner closure with no Task", next)
	}
	if next.GateDecision == nil || !next.GateDecision.Allowed {
		t.Fatalf("GateDecision = %+v, want the Objective completion resolver's allowed decision", next.GateDecision)
	}
	if next.Clearance == nil || next.Clearance.State != ClearanceCurrent || next.Clearance.Check != "C-001" {
		t.Fatalf("Clearance = %+v, want current C-001", next.Clearance)
	}
}

// TestResolveNext_selectedObjectiveWithIncompleteTasksRequestsTaskSelection
// proves that an Objective with unfinished owned Tasks remains selected
// until the router names which one to work on.
func TestResolveNext_selectedObjectiveWithIncompleteTasksRequestsTaskSelection(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnInProgress}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextSelectTask {
		t.Fatalf("Kind = %q, want select_task for the selected Objective", next.Kind)
	}
	if next.Objective == nil || next.Objective.ID != "O-001" || next.Task != nil {
		t.Fatalf("next = %+v, want O-001 with no Task", next)
	}
}

// TestResolveNext_noObjectiveSelectionDoesNotPickReadyTask proves a ready
// Task elsewhere is not a substitute for an empty router selection.
func TestResolveNext_noObjectiveSelectionDoesNotPickReadyTask(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-002"] = &TaskV2{ID: "T-002", Objective: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001", "T-002"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextNothingSelected || next.Task != nil || next.Objective != nil {
		t.Fatalf("next = %+v, want nothing selected and no substituted Task", next)
	}
	if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionReleaseMissing {
		t.Errorf("SelectionDiagnostic = %+v, want missing Goal", next.SelectionDiagnostic)
	}
}

// TestResolveNext_unselectedObjectiveIsNotChosen proves an Objective with no
// Tasks is not selected when the router names no Objective.
func TestResolveNext_unselectedObjectiveIsNotChosen(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextNothingSelected || next.Objective != nil || next.Task != nil {
		t.Fatalf("next = %+v, want nothing selected without choosing O-001", next)
	}
	if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionReleaseMissing {
		t.Errorf("SelectionDiagnostic = %+v, want missing Goal", next.SelectionDiagnostic)
	}
}

// TestResolveNext_missingGoalForEmptyProject proves a fresh project with no
// Goal points the owner to Goal selection before planning work.
func TestResolveNext_missingGoalForEmptyProject(t *testing.T) {
	index := newV2TestIndex()
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextNothingSelected {
		t.Fatalf("Kind = %q, want nothing_selected", next.Kind)
	}
	if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionReleaseMissing {
		t.Errorf("SelectionDiagnostic = %+v, want missing Goal for an empty project", next.SelectionDiagnostic)
	}
	if next.Objective != nil || next.Task != nil {
		t.Errorf("next = %+v, want no Objective or Task selected", next)
	}
}

// TestResolveNext_nothingSelectedWhenNothingSelected proves completed work
// elsewhere does not change an empty router selection.
func TestResolveNext_nothingSelectedWhenNothingSelected(t *testing.T) {
	index := newV2TestIndex()
	mustObjectiveCheck(index, "C-001", "O-001", CheckResultClear)
	index.Objectives["O-001"] = mustCurrentObjective("O-001", "C-001", nil)
	index.Objectives["O-001"].Status = ColumnDone
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextNothingSelected {
		t.Fatalf("Kind = %q, want nothing_selected when the router selects nothing", next.Kind)
	}
}

// TestResolveNext_unresolvedSelectionDoesNotSubstituteAvailableWork proves an
// unresolved router selection reports its diagnostic without selecting a
// different Task from the project.
func TestResolveNext_unresolvedSelectionDoesNotSubstituteAvailableWork(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-999"} // T-999 does not exist

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionNotFound {
		t.Fatalf("SelectionDiagnostic = %+v, want SelectionNotFound for T-999", next.SelectionDiagnostic)
	}
	if next.Kind != NextNothingSelected || next.Task != nil || next.Objective != nil {
		t.Fatalf("next = %+v, want nothing selected without substituting T-001", next)
	}
}

// TestResolveNext_noGoalCarriesDiagnostic proves a router naming no Goal is
// diagnosed even when it names no Objective.
func TestResolveNext_noGoalCarriesDiagnostic(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseIdea}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != SelectionReleaseMissing {
		t.Errorf("SelectionDiagnostic = %+v, want missing Goal", next.SelectionDiagnostic)
	}
	if next.Kind != NextNothingSelected || next.Task != nil || next.Objective != nil {
		t.Errorf("next = %+v, want nothing selected without choosing O-001 or T-001", next)
	}
}

// TestResolveNext_readsNoFilesystem builds its V2Index directly in memory,
// with no discovery or file read anywhere in the call, proving ResolveNext
// consults only the values it is handed.
func TestResolveNext_readsNoFilesystem(t *testing.T) {
	index := &V2Index{
		Objectives:     map[string]*ObjectiveV2{"O-001": {ID: "O-001", Status: ColumnPlanned}},
		Tasks:          map[string]*TaskV2{"T-001": {ID: "T-001", Objective: "O-001", Status: ColumnPlanned}},
		ObjectiveTasks: map[string][]string{"O-001": {"T-001"}},
		Checks:         map[string]*CheckV2{},
		ScopeChecks:    map[string][]string{},
		LatestCheck:    map[string]string{},
	}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Kind != NextExecute {
		t.Fatalf("Kind = %q, want execute, from an in-memory-only index", next.Kind)
	}
}

// TestResolveNext_issuesLinkedToSelectedTask proves a selected Task's own
// linked Issues (TaskIssues) land on Next.Issues, in the same sorted order
// the index already keeps them in.
func TestResolveNext_issuesLinkedToSelectedTask(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Issues = map[string]*IssueV2{
		"I-001": {ID: "I-001", Title: "Found during build", Type: IssueTypeDefect, Status: IssueStatusOpen},
		"I-002": {ID: "I-002", Title: "Unrelated", Type: IssueTypeDrift, Status: IssueStatusOpen},
	}
	index.TaskIssues = map[string][]string{"T-001": {"I-001"}}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if len(next.Issues) != 1 || next.Issues[0].ID != "I-001" {
		t.Fatalf("Issues = %+v, want exactly [I-001]", next.Issues)
	}
}

// TestResolveNext_issuesForObjectiveUnionOwnedTasksAndOwnChecks proves an
// Objective selected with no Task reports the union of every owned Task's
// linked Issues and any Issue linked to a Check scoped to the Objective
// itself — since IssueV2 carries no direct Objective reference — sorted and
// deduplicated.
func TestResolveNext_issuesForObjectiveUnionOwnedTasksAndOwnChecks(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnDone}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	index.Issues = map[string]*IssueV2{
		"I-001": {ID: "I-001", Title: "From owned task", Type: IssueTypeDefect, Status: IssueStatusOpen},
		"I-002": {ID: "I-002", Title: "From objective's own check", Type: IssueTypeGuardrail, Status: IssueStatusOpen},
	}
	index.TaskIssues = map[string][]string{"T-001": {"I-001"}}
	index.ScopeChecks = map[string][]string{"O-001": {"C-001"}}
	index.CheckIssues = map[string][]string{"C-001": {"I-002"}}
	router := &RouterStateV2{State: RouterPhaseDesign, Objective: "O-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if len(next.Issues) != 2 || next.Issues[0].ID != "I-001" || next.Issues[1].ID != "I-002" {
		t.Fatalf("Issues = %+v, want [I-001 I-002] sorted", next.Issues)
	}
}

// TestResolveNext_noIssuesLeavesNilNotEmpty proves a selected Task with no
// linked Issues reports a nil slice, so a rendering surface can tell "no
// Issues" apart from "an empty Issues section".
func TestResolveNext_noIssuesLeavesNilNotEmpty(t *testing.T) {
	index := newV2TestIndex()
	index.Objectives["O-001"] = &ObjectiveV2{ID: "O-001", Status: ColumnPlanned}
	index.Tasks["T-001"] = &TaskV2{ID: "T-001", Objective: "O-001", Status: ColumnPlanned}
	index.ObjectiveTasks["O-001"] = []string{"T-001"}
	router := &RouterStateV2{State: RouterPhaseTask, Objective: "O-001", Task: "T-001"}

	next := resolveNextWithTestGoal(NextInput{Index: index, Router: router})
	if next.Issues != nil {
		t.Fatalf("Issues = %+v, want nil", next.Issues)
	}
}

// TestNext_packageDoesNotImportMigrate proves internal/data does not import
// internal/migrate. internal/migrate already imports internal/data for the
// V2 records it converts into and the project-root helpers it shares, so the
// reverse import would close a cycle.
func TestNext_packageDoesNotImportMigrate(t *testing.T) {
	if _, err := os.Stat("next.go"); err != nil {
		t.Skipf("package source is not beside the test binary: %v", err)
	}
	pkg, err := build.ImportDir(".", build.IgnoreVendor)
	if err != nil {
		t.Fatalf("scan package imports: %v", err)
	}
	for _, imported := range append(pkg.Imports, pkg.TestImports...) {
		if strings.HasSuffix(imported, "/internal/migrate") {
			t.Errorf("internal/data imports %s, which already imports internal/data and would close a cycle", imported)
		}
	}
}

// TestDataPackage_staysBeneathEverySurfaceItFeeds proves internal/data
// imports none of internal/resume, internal/board, internal/doctor, or
// internal/migrate (E48 T-007). The projection this package computes is
// meant to be consumed by every one of those surfaces; a package that
// imported any of them back would mean the projection had started leaning
// on how one particular consumer renders or reports it, which is exactly
// the parallel-interpretation failure mode E48-Detail's "one projection,
// computed above the gates and below every surface" design commits against.
func TestDataPackage_staysBeneathEverySurfaceItFeeds(t *testing.T) {
	if _, err := os.Stat("next.go"); err != nil {
		t.Skipf("package source is not beside the test binary: %v", err)
	}
	pkg, err := build.ImportDir(".", build.IgnoreVendor)
	if err != nil {
		t.Fatalf("scan package imports: %v", err)
	}
	forbidden := []string{"/internal/resume", "/internal/board", "/internal/doctor", "/internal/migrate"}
	for _, imported := range append(pkg.Imports, pkg.TestImports...) {
		for _, suffix := range forbidden {
			if strings.HasSuffix(imported, suffix) {
				t.Errorf("internal/data imports %s, which internal/data must stay beneath rather than depend on", imported)
			}
		}
	}
}

// TestResolveNext_nilInputsReturnAValueRatherThanPanicking covers the
// boundary a long-running consumer reaches that a command does not: the
// board holds a NextInput across reloads, so it can call with a load that
// has not completed or did not succeed. An empty project is the reading;
// reporting the failed load stays the caller's job.
func TestResolveNext_nilInputsReturnAValueRatherThanPanicking(t *testing.T) {
	index := &V2Index{Objectives: map[string]*ObjectiveV2{}, Tasks: map[string]*TaskV2{}}
	router := &RouterStateV2{State: RouterPhaseDesign}

	cases := []struct {
		name           string
		input          NextInput
		want           NextKind
		wantDiagnostic SelectionDiagnosticKind
	}{
		{"nil index", NextInput{Router: router}, NextNothingSelected, SelectionReleaseMissing},
		{"nil router", NextInput{Index: index}, NextNothingSelected, SelectionReleaseMissing},
		{"both nil", NextInput{}, NextNothingSelected, SelectionReleaseMissing},
		{"zero value", NextInput{}, NextNothingSelected, SelectionReleaseMissing},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			next := resolveNextWithTestGoal(tc.input)
			if next.Kind != tc.want {
				t.Errorf("resolveNextWithTestGoal().Kind = %q, want %q", next.Kind, tc.want)
			}
			if next.Objective != nil || next.Task != nil {
				t.Errorf("resolveNextWithTestGoal() selected a record from an absent project: %+v", next)
			}
			if tc.wantDiagnostic == "" && next.SelectionDiagnostic != nil {
				t.Errorf("resolveNextWithTestGoal() diagnostic = %+v, want nil for an absent router", next.SelectionDiagnostic)
			} else if tc.wantDiagnostic != "" && (next.SelectionDiagnostic == nil || next.SelectionDiagnostic.Kind != tc.wantDiagnostic) {
				t.Errorf("resolveNextWithTestGoal() diagnostic = %+v, want %q", next.SelectionDiagnostic, tc.wantDiagnostic)
			}
			if len(next.Issues) != 0 {
				t.Errorf("resolveNextWithTestGoal() reported %d issues from an absent project", len(next.Issues))
			}
		})
	}
}

// TestResolveNext_nilInputMatchesAnEmptyProject pins the guard to the
// reading it claims: an absent index answers exactly as a project that
// loaded and has nothing in it.
func TestResolveNext_nilInputMatchesAnEmptyProject(t *testing.T) {
	empty := resolveNextWithTestGoal(NextInput{
		Index: &V2Index{},
	})
	absent := resolveNextWithTestGoal(NextInput{})

	if empty.Kind != absent.Kind {
		t.Errorf("empty project resolved %q, absent input resolved %q", empty.Kind, absent.Kind)
	}
}
