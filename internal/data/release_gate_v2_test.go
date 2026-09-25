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

func TestResolveReleaseCompletion_requiresMemberObjectivesOnly(t *testing.T) {
	index := releaseGateIndex()
	index.ReleaseObjectives["R-001"] = nil
	got := ResolveReleaseCompletion(index, "R-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockReleaseNoObjectives {
		t.Fatalf("no-member decision = %+v, want release-no-objectives blocker", got)
	}

	index = releaseGateIndex()
	index.Objectives["O-001"].Evidence.Freshness.State = FreshnessStale
	got = ResolveReleaseCompletion(index, "R-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockReleaseObjectiveIncomplete {
		t.Fatalf("incomplete-objective decision = %+v, want objective blocker", got)
	}

	index = releaseGateIndex()
	index.Checks["C-002"].Result = CheckResultNeedsWork
	got = ResolveReleaseCompletion(index, "R-001")
	if !got.Allowed || got.Actor != ActorRoleOwner {
		t.Fatalf("member-complete decision with a NEEDS WORK Goal Check = %+v, want owner completion allowed", got)
	}
}

func TestResolveReleaseCompletion_ignoresEveryGoalCheckClearanceState(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*V2Index)
	}{
		{
			name: "missing",
			configure: func(index *V2Index) {
				index.LatestCheck["R-001"] = ""
			},
		},
		{
			name: "unknown",
			configure: func(index *V2Index) {
				index.Releases["R-001"].Evidence.Freshness.State = FreshnessUnknown
			},
		},
		{
			name: "stale",
			configure: func(index *V2Index) {
				index.Releases["R-001"].Evidence.Freshness.State = FreshnessStale
			},
		},
		{
			name: "needs work",
			configure: func(index *V2Index) {
				index.Checks["C-002"].Result = CheckResultNeedsWork
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := releaseGateIndex()
			tt.configure(index)
			got := ResolveReleaseCompletion(index, "R-001")
			if !got.Allowed || got.Actor != ActorRoleOwner || len(got.Blockers) != 0 {
				t.Fatalf("decision = %+v, want member-based Goal completion regardless of Goal Check state", got)
			}
		})
	}
}

func TestResolveReleaseCompletion_ignoresGoalCheckAndItsIssues(t *testing.T) {
	index := releaseGateIndex()
	index.Checks["C-002"].Issues = []string{"I-001"}
	index.Issues["I-001"] = &IssueV2{ID: "I-001", Status: IssueStatusOpen}
	index.CheckIssues["C-002"] = []string{"I-001"}
	index.Checks["C-002"].Result = CheckResultNeedsWork

	got := ResolveReleaseCompletion(index, "R-001")
	index.Releases["R-001"].Evidence.Exception = &Exception{
		Requirements: []string{"I-001"}, Reason: "owner accepted the known risk",
		Owner: "owner-1", Check: "C-002",
	}
	got = ResolveReleaseCompletion(index, "R-001")
	if !got.Allowed || got.AllowedByException || got.Exception != nil || got.Actor != ActorRoleOwner {
		t.Fatalf("decision with a NEEDS WORK Goal Check, unresolved linked Issue, and Goal exception = %+v, want ordinary member-based completion", got)
	}
}

func TestResolveReleaseCompletion_allowsCompleteMembersWithoutGoalCheck(t *testing.T) {
	index := releaseGateIndex()
	delete(index.Checks, "C-002")
	delete(index.ScopeChecks, "R-001")
	delete(index.LatestCheck, "R-001")
	index.Releases["R-001"].Evidence = nil

	got := ResolveReleaseCompletion(index, "R-001")
	if !got.Allowed || got.Actor != ActorRoleOwner {
		t.Fatalf("decision without Goal Check evidence = %+v, want allowed", got)
	}
}

func TestResolveReleaseCompletion_legacyCompletionWithLiveMembersIsUnchanged(t *testing.T) {
	index := releaseGateIndex()
	index.Releases["R-001"].Status = ColumnDone
	index.Releases["R-001"].Evidence = nil
	index.Releases["R-001"].LegacyCompletion = &LegacyCompletionReference{
		SourcePath: "releases/v1/PRD.md", ArchivePath: ".savepoint/archive/v1/releases/v1/PRD.md",
		SHA256: strings.Repeat("a", 64),
	}

	got := ResolveReleaseCompletion(index, "R-001")
	if !got.Allowed || !got.AllowedByLegacyCompletion || got.LegacyCompletion == nil {
		t.Fatalf("legacy member-complete decision = %+v, want preserved historical completion", got)
	}

	index.Objectives["O-001"].Status = ColumnPlanned
	got = ResolveReleaseCompletion(index, "R-001")
	if got.Allowed || len(got.Blockers) != 1 || got.Blockers[0].Kind != GateBlockReleaseObjectiveIncomplete {
		t.Fatalf("legacy incomplete-member decision = %+v, want member blocker", got)
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
