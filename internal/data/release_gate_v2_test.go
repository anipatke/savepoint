package data

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func releaseGateIndex() *V2Index {
	checkedAt := time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC)
	return &V2Index{
		Releases: map[string]*ReleaseV2{
			"R-001": {
				ID: "R-001", Status: ColumnInProgress,
				Evidence: &Evidence{Freshness: &Freshness{
					State: FreshnessCurrent, Check: "C-002",
					AssessedBy: Actor{Role: ActorRoleChecker, Session: "checker-release"},
					AssessedAt: checkedAt, Basis: "Release integration reviewed",
				}},
			},
		},
		Objectives: map[string]*ObjectiveV2{
			"O-001": {
				ID: "O-001", Status: ColumnDone,
				Evidence: &Evidence{Freshness: &Freshness{
					State: FreshnessCurrent, Check: "C-001",
					AssessedBy: Actor{Role: ActorRoleChecker, Session: "checker-objective"},
					AssessedAt: checkedAt, Basis: "Objective integration reviewed",
				}},
			},
		},
		Tasks: map[string]*TaskV2{
			"T-001": {ID: "T-001", Objective: "O-001", Status: ColumnDone},
		},
		Checks: map[string]*CheckV2{
			"C-001": {ID: "C-001", Scope: CheckScope{Kind: CheckScopeObjective, ID: "O-001"}, Result: CheckResultClear, CheckedBy: Actor{Role: ActorRoleChecker, Session: "checker-objective"}, ExecutedSession: "build-objective"},
			"C-002": {ID: "C-002", Scope: CheckScope{Kind: CheckScopeRelease, ID: "R-001"}, Result: CheckResultClear, CheckedBy: Actor{Role: ActorRoleChecker, Session: "checker-release"}, ExecutedSession: "build-release"},
		},
		Issues:            map[string]*IssueV2{},
		ObjectiveTasks:    map[string][]string{"O-001": {"T-001"}},
		ReleaseObjectives: map[string][]string{"R-001": {"O-001"}},
		CheckIssues:       map[string][]string{},
		ScopeChecks:       map[string][]string{"O-001": {"C-001"}, "R-001": {"C-002"}},
		LatestCheck:       map[string]string{"O-001": "C-001", "R-001": "C-002"},
	}
}

func TestDecodeCheckV2_releaseScope(t *testing.T) {
	content := `---
id: C-010
scope: {kind: release, id: R-001}
result: CLEAR
checked_by: {role: checker, session: release-checker}
executed_session: release-build
checked_at: '2026-09-15T00:00:00Z'
---

# Release Check
`

	check, err := DecodeCheckV2("checks/C-010-release.md", content)
	if err != nil {
		t.Fatalf("DecodeCheckV2() error = %v", err)
	}
	if check.Scope != (CheckScope{Kind: CheckScopeRelease, ID: "R-001"}) {
		t.Fatalf("Scope = %+v, want release/R-001", check.Scope)
	}
}

func TestResolveReleaseCompletion_requiresMemberObjectivesAndCurrentAcceptance(t *testing.T) {
	index := releaseGateIndex()
	index.ReleaseObjectives["R-001"] = nil
	got := ResolveReleaseCompletion(index, "R-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockReleaseNoObjectives {
		t.Fatalf("no-member decision = %+v, want release-no-objectives blocker", got)
	}

	index = releaseGateIndex()
	index.Objectives["O-001"].Evidence = nil
	got = ResolveReleaseCompletion(index, "R-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockReleaseObjectiveIncomplete {
		t.Fatalf("incomplete-objective decision = %+v, want objective blocker", got)
	}

	index = releaseGateIndex()
	got = ResolveReleaseCompletion(index, "R-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockOwnerAcceptance {
		t.Fatalf("unaccepted decision = %+v, want owner blocker", got)
	}

	index.Releases["R-001"].Evidence.OwnerValidation = &OwnerValidation{
		AcceptedCheck: "C-002",
		AcceptedBy:    Actor{Role: ActorRoleOwner, Session: "owner-1"},
	}
	got = ResolveReleaseCompletion(index, "R-001")
	if !got.Allowed || got.Actor != ActorRoleChecker {
		t.Fatalf("accepted decision = %+v, want allowed by current Release Check", got)
	}
}

func TestResolveReleaseCompletion_refusesEveryNonCurrentReleaseState(t *testing.T) {
	tests := []struct {
		name       string
		configure  func(*V2Index)
		wantKind   GateBlockKind
		wantDetail string
	}{
		{
			name: "missing",
			configure: func(index *V2Index) {
				index.LatestCheck["R-001"] = ""
			},
			wantKind: GateBlockClearanceMissing,
		},
		{
			name: "unknown",
			configure: func(index *V2Index) {
				index.Releases["R-001"].Evidence.Freshness = nil
			},
			wantKind: GateBlockClearanceUnknown,
		},
		{
			name: "stale",
			configure: func(index *V2Index) {
				index.Releases["R-001"].Evidence.Freshness.Check = "C-001"
			},
			wantKind: GateBlockClearanceStale,
		},
		{
			name: "needs work",
			configure: func(index *V2Index) {
				index.Checks["C-002"].Result = CheckResultNeedsWork
			},
			wantKind: GateBlockClearanceNeedsWork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := releaseGateIndex()
			tt.configure(index)
			got := ResolveReleaseCompletion(index, "R-001")
			if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != tt.wantKind {
				t.Fatalf("decision = %+v, want one %s blocker", got, tt.wantKind)
			}
		})
	}
}

func TestResolveReleaseCompletion_requiresIssueResolutionOrScopedException(t *testing.T) {
	index := releaseGateIndex()
	index.Checks["C-002"].Issues = []string{"I-001"}
	index.Issues["I-001"] = &IssueV2{ID: "I-001", Status: IssueStatusOpen}
	index.CheckIssues["C-002"] = []string{"I-001"}

	got := ResolveReleaseCompletion(index, "R-001")
	if got.Allowed || len(got.Blockers) != 2 {
		t.Fatalf("unresolved issue decision = %+v, want issue and owner blockers", got)
	}
	if got.Blockers[0].Kind != GateBlockReleaseIssueUnresolved || got.Blockers[0].Issue != "I-001" {
		t.Fatalf("first blocker = %+v, want unresolved I-001", got.Blockers[0])
	}

	index.Releases["R-001"].Evidence.Exception = &Exception{
		Requirements: []string{"I-001"}, Reason: "owner accepted the known risk",
		Owner: "owner-1", Check: "C-002",
	}
	index.Releases["R-001"].Evidence.OwnerValidation = &OwnerValidation{
		AcceptedCheck: "C-002", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"},
	}
	got = ResolveReleaseCompletion(index, "R-001")
	if !got.Allowed || !got.AllowedByException || got.Exception == nil {
		t.Fatalf("exception decision = %+v, want allowed by scoped owner exception", got)
	}
}

func TestResolveReleaseCompletion_supersededAcceptanceAndFreshnessDoNotCarryForward(t *testing.T) {
	index := releaseGateIndex()
	index.Checks["C-003"] = &CheckV2{
		ID: "C-003", Scope: CheckScope{Kind: CheckScopeRelease, ID: "R-001"},
		Result: CheckResultClear, CheckedBy: Actor{Role: ActorRoleChecker, Session: "checker-release-2"},
		ExecutedSession: "build-release-2", Supersedes: "C-002",
	}
	index.ScopeChecks["R-001"] = []string{"C-002", "C-003"}
	index.LatestCheck["R-001"] = "C-003"
	index.Releases["R-001"].Evidence.OwnerValidation = &OwnerValidation{
		AcceptedCheck: "C-002", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"},
	}

	got := ResolveReleaseCompletion(index, "R-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockClearanceStale {
		t.Fatalf("superseded decision = %+v, want stale Release clearance", got)
	}
}

func TestResolveReleaseCompletion_historicalReferenceIsSeparateFromClearance(t *testing.T) {
	index := releaseGateIndex()
	index.Releases["R-001"].Status = ColumnDone
	index.Releases["R-001"].Evidence = nil
	index.LatestCheck["R-001"] = ""
	index.Checks = map[string]*CheckV2{}
	index.ScopeChecks = map[string][]string{}
	index.Releases["R-001"].LegacyCompletion = &LegacyCompletionReference{
		SourcePath: "releases/v1/PRD.md", ArchivePath: ".savepoint/archive/v1/releases/v1/PRD.md",
		SHA256: strings.Repeat("a", 64),
	}

	decision := ResolveReleaseCompletion(index, "R-001")
	if !decision.Allowed || !decision.AllowedByLegacyCompletion || decision.LegacyCompletion == nil {
		t.Fatalf("historical decision = %+v, want separate historical allowance", decision)
	}
	if clearance := ResolveClearance(index, "R-001"); clearance.State != ClearanceMissing {
		t.Fatalf("historical clearance = %+v, want missing rather than current CLEAR", clearance)
	}
}

func TestResolveReleaseCompletion_allowsArchivedHistoricalReleaseWithoutLiveObjectives(t *testing.T) {
	index := releaseGateIndex()
	index.Releases["R-001"].Status = ColumnDone
	index.Releases["R-001"].Evidence = nil
	index.ReleaseObjectives["R-001"] = nil
	index.LatestCheck["R-001"] = ""
	index.Checks = map[string]*CheckV2{}
	index.ScopeChecks = map[string][]string{}
	index.Releases["R-001"].LegacyCompletion = &LegacyCompletionReference{
		SourcePath: "releases/v1/PRD.md", ArchivePath: ".savepoint/archive/v1/releases/v1/PRD.md",
		SHA256: strings.Repeat("a", 64),
	}

	decision := ResolveReleaseCompletion(index, "R-001")
	if !decision.Allowed || !decision.AllowedByLegacyCompletion {
		t.Fatalf("historical no-member decision = %+v, want allowed historical completion", decision)
	}
}

func TestResolveReleaseCutoverComposesCanonicalReleaseDecisions(t *testing.T) {
	index := releaseGateIndex()

	blocked := ResolveReleaseCutover(index)
	if blocked.Allowed || len(blocked.Blockers) != 1 {
		t.Fatalf("blocked cutover decision = %+v, want one canonical blocker", blocked)
	}
	if blocked.Blockers[0].ReleaseID != "R-001" || blocked.Blockers[0].Gate.Kind != GateBlockOwnerAcceptance {
		t.Fatalf("blocked cutover blocker = %+v, want R-001 owner acceptance", blocked.Blockers[0])
	}

	index.Releases["R-001"].Evidence.OwnerValidation = &OwnerValidation{
		AcceptedCheck: "C-002",
		AcceptedBy:    Actor{Role: ActorRoleOwner, Session: "owner-1"},
	}
	allowed := ResolveReleaseCutover(index)
	if !allowed.Allowed || len(allowed.Blockers) != 0 {
		t.Fatalf("accepted cutover decision = %+v, want allowed", allowed)
	}
}

func TestResolveReleaseCutover_refusesTechnicalIssueAndOwnerStates(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*V2Index)
		wantKind  GateBlockKind
	}{
		{
			name: "technically unclear",
			configure: func(index *V2Index) {
				index.Releases["R-001"].Evidence = nil
			},
			wantKind: GateBlockClearanceUnknown,
		},
		{
			name: "material Issue",
			configure: func(index *V2Index) {
				index.Checks["C-002"].Issues = []string{"I-001"}
				index.CheckIssues["C-002"] = []string{"I-001"}
				index.Issues["I-001"] = &IssueV2{ID: "I-001", Status: IssueStatusOpen}
			},
			wantKind: GateBlockReleaseIssueUnresolved,
		},
		{
			name:      "owner acceptance",
			configure: func(*V2Index) {},
			wantKind:  GateBlockOwnerAcceptance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := releaseGateIndex()
			tt.configure(index)
			decision := ResolveReleaseCutover(index)
			if decision.Allowed || len(decision.Blockers) == 0 {
				t.Fatalf("cutover decision = %+v, want refusal", decision)
			}
			found := false
			for _, blocker := range decision.Blockers {
				if blocker.ReleaseID == "R-001" && blocker.Gate.Kind == tt.wantKind {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("cutover blockers = %+v, want R-001/%s", decision.Blockers, tt.wantKind)
			}
		})
	}
}

func TestResolveReleaseCutover_allowsOptionalReleaseModel(t *testing.T) {
	decision := ResolveReleaseCutover(&V2Index{Releases: map[string]*ReleaseV2{}})
	if !decision.Allowed || len(decision.Blockers) != 0 {
		t.Fatalf("no-Release cutover decision = %+v, want allowed optional model", decision)
	}
}

func TestDecodeReleaseV2_legacyCompletionIsTypedAndDoneOnly(t *testing.T) {
	content := `---
id: R-001
title: "Historical release"
status: done
legacy_completion:
  source_path: releases/v1/PRD.md
  archive_path: .savepoint/archive/v1/releases/v1/PRD.md
  sha256: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
---

# Release

## Outcome

Historical outcome.

## Why

Historical reason.

## Success Conditions

Historical conditions.

## Boundaries

Historical boundaries.
`
	release, err := DecodeReleaseV2("releases/R-001-historical/Release.md", content)
	if err != nil {
		t.Fatalf("DecodeReleaseV2() error = %v", err)
	}
	if release.LegacyCompletion == nil || release.LegacyCompletion.ArchivePath == "" {
		t.Fatalf("LegacyCompletion = %+v, want typed archive reference", release.LegacyCompletion)
	}

	bad := strings.Replace(content, "status: done", "status: in_progress", 1)
	if _, err := DecodeReleaseV2("Release.md", bad); !errors.Is(err, ErrV2ReleaseLegacyMalformed) {
		t.Fatalf("DecodeReleaseV2() error = %v, want legacy lifecycle diagnostic", err)
	}
}
