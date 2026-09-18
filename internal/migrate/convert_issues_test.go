package migrate

import (
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

// --- round trip over the frozen fixtures ----------------------------------

// TestConvertIssue_roundTripDecodesFixtureIssues proves every Issue target
// planned over both frozen fixtures renders content that decodes through the
// real strict V2 decoder with no diagnostic.
func TestConvertIssue_roundTripDecodesFixtureIssues(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			root := fixtureProjectRoot(fixture)
			p := mustPlan(t, root)

			var sawIssue bool
			for _, target := range p.Targets {
				if target.Kind != TargetIssue {
					continue
				}
				sawIssue = true
				content, err := ConvertIssue(root, p, target)
				if err != nil {
					t.Fatalf("ConvertIssue(%s) error = %v", target.GlobalID, err)
				}
				issue, err := data.DecodeIssueV2(target.TargetPath, content)
				if err != nil {
					t.Fatalf("DecodeIssueV2(%s) error = %v", target.GlobalID, err)
				}
				if issue.Origin.Kind != data.IssueOriginMigration {
					t.Errorf("issue %s: Origin.Kind = %q, want migration", target.GlobalID, issue.Origin.Kind)
				}
				if issue.Origin.Check != "" {
					t.Errorf("issue %s: Origin.Check = %q, want empty: migration never borrows a Check's authority", target.GlobalID, issue.Origin.Check)
				}
			}
			if fixture == "v1-history" && !sawIssue {
				t.Fatal("v1-history planned no Issue targets; expected at least the open defect and findings")
			}
		})
	}
}

// TestConvertIssue_deterministicAcrossRepeatedRuns mirrors
// TestConvert_deterministicAcrossRepeatedRuns for Issue targets.
func TestConvertIssue_deterministicAcrossRepeatedRuns(t *testing.T) {
	root := fixtureProjectRoot("v1-history")
	p := mustPlan(t, root)

	for _, target := range p.Targets {
		if target.Kind != TargetIssue {
			continue
		}
		first, err := ConvertIssue(root, p, target)
		if err != nil {
			t.Fatalf("ConvertIssue(%s) error = %v", target.GlobalID, err)
		}
		second, err := ConvertIssue(root, p, target)
		if err != nil {
			t.Fatalf("ConvertIssue(%s) second call error = %v", target.GlobalID, err)
		}
		if first != second {
			t.Errorf("ConvertIssue(%s) not deterministic:\nfirst=%q\nsecond=%q", target.GlobalID, first, second)
		}
	}
}

// --- defect disposition mapping -------------------------------------------

// TestConvertIssue_defectDispositions proves an open V1 defect becomes an
// open Issue with type defect, and an in_progress one becomes in_progress,
// each linking the Task its reference named.
func TestConvertIssue_defectDispositions(t *testing.T) {
	cases := []struct {
		name       string
		status     string
		wantStatus data.IssueStatus
	}{
		{"open", "open", data.IssueStatusOpen},
		{"in_progress", "in_progress", data.IssueStatusInProgress},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFiles(t, root, map[string]string{
				".savepoint/releases/v1/epics/E01-x/E01-Detail.md": "---\nstatus: in_progress\n---\n\n# E01\n",
				".savepoint/releases/v1/epics/E01-x/tasks/T001-fix.md": "---\n" +
					"id: E01-x/T001-fix\nstatus: planned\ndepends_on: []\n---\n\n# T001\n",
				".savepoint/releases/v1/defects/D001-broken.md": "---\n" +
					"id: v1/D001-broken\nrelease: v1\nstatus: " + tc.status + "\nseverity: medium\n" +
					"title: \"Something broke\"\nreference: E01-x/T001-fix\n---\n\n# D001\n\nSymptom text.\n",
			})

			p := mustPlan(t, root)
			target, ok := targetByPath(p, ".savepoint/releases/v1/defects/D001-broken.md")
			if !ok || target.Kind != TargetIssue {
				t.Fatalf("defect was not planned as an Issue target: %+v", target)
			}

			content, err := ConvertIssue(root, p, target)
			if err != nil {
				t.Fatalf("ConvertIssue error = %v", err)
			}
			issue, err := data.DecodeIssueV2(target.TargetPath, content)
			if err != nil {
				t.Fatalf("DecodeIssueV2 error = %v", err)
			}

			if issue.Type != data.IssueTypeDefect {
				t.Errorf("Type = %q, want defect", issue.Type)
			}
			if issue.Status != tc.wantStatus {
				t.Errorf("Status = %q, want %q", issue.Status, tc.wantStatus)
			}
			taskTarget, _ := targetByPath(p, ".savepoint/releases/v1/epics/E01-x/tasks/T001-fix.md")
			if len(issue.Tasks) != 1 || issue.Tasks[0] != taskTarget.GlobalID {
				t.Errorf("Tasks = %v, want [%s] (resolved from the defect's reference)", issue.Tasks, taskTarget.GlobalID)
			}
		})
	}
}

// TestConvertIssue_resolvedDefectNeverReachesConvert proves a resolved
// defect is archived rather than planned, so ConvertIssue is never asked to
// render one.
func TestConvertIssue_resolvedDefectNeverReachesConvert(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-history"))
	const resolved = ".savepoint/releases/v1/defects/D001-shared.md"
	if _, ok := targetByPath(p, resolved); ok {
		t.Errorf("resolved defect was planned as a target, want archive-only")
	}
	if _, ok := archiveByPath(p, resolved); !ok {
		t.Errorf("resolved defect was not archived")
	}
}

// --- finding disposition mapping ------------------------------------------

func findingProject(t *testing.T, extraFrontmatter, status string) (root string, findingPath string) {
	t.Helper()
	root = t.TempDir()
	findingPath = ".savepoint/audit/findings/F001-example.md"
	writeFiles(t, root, map[string]string{
		findingPath: "---\n" +
			"id: F001\ntitle: Example finding\nstatus: " + status + "\nseverity: medium\nconfidence: high\n" +
			"first_seen: 2026-01-01\nlast_seen: 2026-02-01\nproof_needed: a regression test\n" +
			extraFrontmatter +
			"---\n\n## Summary\n\nSomething was found.\n",
	})
	return root, findingPath
}

// TestConvertIssue_findingOpenTriagedMapped proves the three non-terminal
// findings become open Issues carrying their guardrail_ids and severity
// forward, with no history entry invented.
func TestConvertIssue_findingOpenTriagedMapped(t *testing.T) {
	for _, status := range []string{"open", "triaged", "mapped"} {
		t.Run(status, func(t *testing.T) {
			root, findingPath := findingProject(t, "guardrail_ids: [TEST-05]\n", status)
			p := mustPlan(t, root)
			target, ok := targetByPath(p, findingPath)
			if !ok || target.Kind != TargetIssue {
				t.Fatalf("finding was not planned as an Issue target: %+v", target)
			}

			content, err := ConvertIssue(root, p, target)
			if err != nil {
				t.Fatalf("ConvertIssue error = %v", err)
			}
			issue, err := data.DecodeIssueV2(target.TargetPath, content)
			if err != nil {
				t.Fatalf("DecodeIssueV2 error = %v", err)
			}
			if issue.Status != data.IssueStatusOpen {
				t.Errorf("Status = %q, want open", issue.Status)
			}
			if len(issue.GuardrailIDs) != 1 || issue.GuardrailIDs[0] != "TEST-05" {
				t.Errorf("GuardrailIDs = %v, want [TEST-05]", issue.GuardrailIDs)
			}
			if len(issue.History) != 0 {
				t.Errorf("History = %+v, want none for a plain %s finding", issue.History, status)
			}
			if !strings.Contains(content, "Something was found.") {
				t.Errorf("rendered content dropped the authored body")
			}
			if !strings.Contains(content, "a regression test") {
				t.Errorf("rendered content dropped the finding's proof_needed requirement")
			}
		})
	}
}

// TestConvertIssue_deferredAndOwnerDecision proves each becomes an open
// Issue with no resolution, carrying a dated history entry of its own kind.
func TestConvertIssue_deferredAndOwnerDecision(t *testing.T) {
	cases := []struct {
		status   string
		wantKind data.IssueHistoryKind
	}{
		{"deferred", data.IssueHistoryDeferred},
		{"owner_decision", data.IssueHistoryOwnerDecision},
	}

	for _, tc := range cases {
		t.Run(tc.status, func(t *testing.T) {
			root, findingPath := findingProject(t, "deferral_reason: Not worth fixing yet.\n", tc.status)
			p := mustPlan(t, root)
			target, _ := targetByPath(p, findingPath)

			content, err := ConvertIssue(root, p, target)
			if err != nil {
				t.Fatalf("ConvertIssue error = %v", err)
			}
			issue, err := data.DecodeIssueV2(target.TargetPath, content)
			if err != nil {
				t.Fatalf("DecodeIssueV2 error = %v", err)
			}

			if issue.Status != data.IssueStatusOpen {
				t.Errorf("Status = %q, want open: %s is still unresolved follow-up", issue.Status, tc.status)
			}
			if issue.Resolution != nil {
				t.Errorf("Resolution = %+v, want nil: %s never closes an Issue", issue.Resolution, tc.status)
			}
			if len(issue.History) != 1 || issue.History[0].Kind != tc.wantKind {
				t.Fatalf("History = %+v, want one %s entry", issue.History, tc.wantKind)
			}
			if !strings.Contains(issue.History[0].Note, "Not worth fixing yet.") {
				t.Errorf("History note = %q, want it to carry the recorded deferral_reason", issue.History[0].Note)
			}
			wantAt := "2026-02-01T00:00:00Z" // last_seen from findingProject
			if issue.History[0].At.Format("2006-01-02T15:04:05Z") != wantAt {
				t.Errorf("History[0].At = %v, want the finding's own last_seen date %s", issue.History[0].At, wantAt)
			}
		})
	}
}

// TestConvertIssue_inProgressAndFixedNeverResolved proves both dispositions
// land as an in_progress Issue with a repair_attempted entry, and that
// neither is ever laundered into resolved: fixed still awaits independent
// proof.
func TestConvertIssue_inProgressAndFixedNeverResolved(t *testing.T) {
	for _, status := range []string{"in_progress", "fixed"} {
		t.Run(status, func(t *testing.T) {
			root, findingPath := findingProject(t, "", status)
			p := mustPlan(t, root)
			target, _ := targetByPath(p, findingPath)

			content, err := ConvertIssue(root, p, target)
			if err != nil {
				t.Fatalf("ConvertIssue error = %v", err)
			}
			issue, err := data.DecodeIssueV2(target.TargetPath, content)
			if err != nil {
				t.Fatalf("DecodeIssueV2 error = %v", err)
			}

			if issue.Status != data.IssueStatusInProgress {
				t.Errorf("Status = %q, want in_progress", issue.Status)
			}
			if issue.Status == data.IssueStatusResolved {
				t.Fatalf("a %s finding converted to resolved; it still awaits independent verification", status)
			}
			if issue.Resolution != nil {
				t.Errorf("Resolution = %+v, want nil", issue.Resolution)
			}
			if len(issue.History) != 1 || issue.History[0].Kind != data.IssueHistoryRepairAttempted {
				t.Fatalf("History = %+v, want one repair_attempted entry", issue.History)
			}
		})
	}
}

// TestConvertIssue_verifiedFindingProducesNoIssue proves a verified finding
// is archived and never reaches ConvertIssue at all.
func TestConvertIssue_verifiedFindingProducesNoIssue(t *testing.T) {
	root, findingPath := findingProject(t, "verified_proof: regression test landed\n", "verified")
	p := mustPlan(t, root)
	if _, ok := targetByPath(p, findingPath); ok {
		t.Errorf("verified finding was planned as a target, want archive-only")
	}
	if _, ok := archiveByPath(p, findingPath); !ok {
		t.Errorf("verified finding was not archived")
	}
}

// --- waived findings and their manifest reference -------------------------

// TestConvertIssue_waivedFindingReferencedByActiveWork proves the fixture's
// F002 (waived, work_item naming v1.1's active T001-shared) is archived with
// no Issue, and that a typed WaivedReference lands in the manifest naming the
// converted Task rather than an Issue being invented to hold it.
func TestConvertIssue_waivedFindingReferencedByActiveWork(t *testing.T) {
	root := fixtureProjectRoot("v1-history")
	p := mustPlan(t, root)

	const f002 = ".savepoint/audit/findings/F002-owner-waiver.md"
	if _, ok := targetByPath(p, f002); ok {
		t.Fatalf("waived finding was planned as an Issue target, want archive-only")
	}
	archived, ok := archiveByPath(p, f002)
	if !ok {
		t.Fatalf("waived finding was not archived")
	}

	activeTask, ok := targetByPath(p, ".savepoint/releases/v1.1/epics/E01-example/tasks/T001-shared.md")
	if !ok {
		t.Fatalf("expected v1.1's T001-shared to be planned active")
	}

	var found *WaivedReference
	for i := range p.WaivedRefs {
		if p.WaivedRefs[i].Task == activeTask.GlobalID {
			found = &p.WaivedRefs[i]
		}
	}
	if found == nil {
		t.Fatalf("WaivedRefs = %+v, want one naming task %s", p.WaivedRefs, activeTask.GlobalID)
	}
	if found.ArchivePath != archived.ArchivePath {
		t.Errorf("WaivedReference.ArchivePath = %q, want %q", found.ArchivePath, archived.ArchivePath)
	}
	if found.Reason == "" {
		t.Errorf("WaivedReference.Reason is empty, want the finding's recorded waiver_reason")
	}

	m := BuildManifest(p)
	var inManifest bool
	for _, w := range m.WaivedReferences {
		if w.Task == activeTask.GlobalID {
			inManifest = true
		}
	}
	if !inManifest {
		t.Errorf("manifest WaivedReferences = %+v, want an entry for task %s", m.WaivedReferences, activeTask.GlobalID)
	}
}

// TestConvertIssue_waivedFindingWithNoActiveReference proves a waived
// finding that names no active work produces neither an Issue nor a
// WaivedReference: nothing to invent, nothing to record.
func TestConvertIssue_waivedFindingWithNoActiveReference(t *testing.T) {
	root, findingPath := findingProject(t, "waiver_reason: Owner accepted the gap.\n", "waived")
	p := mustPlan(t, root)

	if _, ok := targetByPath(p, findingPath); ok {
		t.Errorf("waived finding was planned as a target, want archive-only")
	}
	if len(p.WaivedRefs) != 0 {
		t.Errorf("WaivedRefs = %+v, want none: the finding named no active work", p.WaivedRefs)
	}
}

// --- duplicate findings, both directions -----------------------------------

// TestConvertIssue_duplicateCanonicalConverts proves the fixture's F003
// resolves as a duplicate Issue naming F001's converted Issue, with no
// independent proof check and no laundering into an unrelated disposition.
func TestConvertIssue_duplicateCanonicalConverts(t *testing.T) {
	root := fixtureProjectRoot("v1-history")
	p := mustPlan(t, root)

	f001Target, ok := targetByPath(p, ".savepoint/audit/findings/F001-awaiting-proof.md")
	if !ok {
		t.Fatalf("F001 was not planned as a target")
	}
	f003Target, ok := targetByPath(p, ".savepoint/audit/findings/F003-duplicate.md")
	if !ok {
		t.Fatalf("F003 was not planned as a target")
	}
	if f003Target.DuplicateOfGlobalID != f001Target.GlobalID {
		t.Fatalf("F003 DuplicateOfGlobalID = %q, want %q (F001's converted Issue)", f003Target.DuplicateOfGlobalID, f001Target.GlobalID)
	}

	content, err := ConvertIssue(root, p, f003Target)
	if err != nil {
		t.Fatalf("ConvertIssue error = %v", err)
	}
	issue, err := data.DecodeIssueV2(f003Target.TargetPath, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2 error = %v", err)
	}

	if issue.Status != data.IssueStatusResolved {
		t.Errorf("Status = %q, want resolved", issue.Status)
	}
	if issue.DuplicateOf != f001Target.GlobalID {
		t.Errorf("DuplicateOf = %q, want %q", issue.DuplicateOf, f001Target.GlobalID)
	}
	if issue.Resolution == nil || issue.Resolution.Disposition != data.IssueDispositionDuplicate {
		t.Fatalf("Resolution = %+v, want disposition duplicate", issue.Resolution)
	}
	if issue.Resolution.Check != "" {
		t.Errorf("Resolution.Check = %q, want empty: a duplicate proves nothing itself", issue.Resolution.Check)
	}
}

// TestConvertIssue_duplicateCanonicalArchiveOnly proves a duplicate finding
// whose canonical is archive-only (verified) is itself archived, producing
// no Issue in either direction.
func TestConvertIssue_duplicateCanonicalArchiveOnly(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		".savepoint/audit/findings/F001-canonical.md": "---\n" +
			"id: F001\ntitle: Canonical finding\nstatus: verified\nseverity: medium\nconfidence: high\n" +
			"first_seen: 2026-01-01\nlast_seen: 2026-02-01\nproof_needed: n/a\nverified_proof: fixed and proven\n" +
			"---\n\n## Summary\n\nAlready verified.\n",
		".savepoint/audit/findings/F002-dup.md": "---\n" +
			"id: F002\ntitle: Duplicate of the canonical\nstatus: duplicate\nseverity: medium\nconfidence: high\n" +
			"first_seen: 2026-01-02\nlast_seen: 2026-02-01\nproof_needed: n/a\nduplicate_of: F001\n" +
			"---\n\n## Summary\n\nSame issue, filed twice.\n",
	})

	p := mustPlan(t, root)
	const dupPath = ".savepoint/audit/findings/F002-dup.md"
	if _, ok := targetByPath(p, dupPath); ok {
		t.Errorf("duplicate of an archive-only canonical was planned as a target, want archive-only")
	}
	if _, ok := archiveByPath(p, dupPath); !ok {
		t.Errorf("duplicate of an archive-only canonical was not archived")
	}
}

// TestConvertIssue_duplicateOfUnknownFinding_isBlockingAmbiguity proves a
// duplicate_of naming no known finding in the project is refused as a
// blocking ambiguity rather than silently archived or guessed at.
func TestConvertIssue_duplicateOfUnknownFinding_isBlockingAmbiguity(t *testing.T) {
	root, findingPath := findingProject(t, "duplicate_of: F999\n", "duplicate")
	p := mustPlan(t, root)

	if _, ok := targetByPath(p, findingPath); ok {
		t.Errorf("duplicate naming an unknown canonical was planned as a target, want blocked as an ambiguity")
	}
	if _, ok := archiveByPath(p, findingPath); ok {
		t.Errorf("duplicate naming an unknown canonical was archived, want blocked as an ambiguity instead of guessed")
	}

	var found *Ambiguity
	for i := range p.Ambiguities {
		if p.Ambiguities[i].Path == findingPath {
			found = &p.Ambiguities[i]
		}
	}
	if found == nil {
		t.Fatalf("Ambiguities = %+v, want one naming %s", p.Ambiguities, findingPath)
	}
	if found.Kind != AmbiguityUnresolvedNarrativeFind {
		t.Errorf("Ambiguity.Kind = %v, want AmbiguityUnresolvedNarrativeFind", found.Kind)
	}
}

// --- source provenance and history date seeding ----------------------------

// TestConvertIssue_sourceRecordsMigrationKindAndLegacyKey proves every
// converted Issue's source names kind: migration with an actor session
// naming the exact V1 source path — the source-qualified legacy key — and
// that source.at is the migration time, never a fabricated date.
func TestConvertIssue_sourceRecordsMigrationKindAndLegacyKey(t *testing.T) {
	root, findingPath := findingProject(t, "", "open")
	p := mustPlan(t, root)
	target, _ := targetByPath(p, findingPath)

	content, err := ConvertIssue(root, p, target)
	if err != nil {
		t.Fatalf("ConvertIssue error = %v", err)
	}
	issue, err := data.DecodeIssueV2(target.TargetPath, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2 error = %v", err)
	}

	if issue.Origin.Kind != data.IssueOriginMigration {
		t.Errorf("Origin.Kind = %q, want migration", issue.Origin.Kind)
	}
	if issue.Origin.Actor.Session != findingPath {
		t.Errorf("Origin.Actor.Session = %q, want the source-qualified legacy path %q", issue.Origin.Actor.Session, findingPath)
	}
	if !issue.Origin.At.Equal(p.GeneratedAt) {
		t.Errorf("Origin.At = %v, want the migration's own GeneratedAt %v", issue.Origin.At, p.GeneratedAt)
	}
}

// TestConvertIssue_seededHistoryDateFallsBackWhenAbsent proves a deferred
// finding recording no last_seen date seeds its history entry with the
// migration time and an explicit note, rather than inventing a date.
func TestConvertIssue_seededHistoryDateFallsBackWhenAbsent(t *testing.T) {
	root := t.TempDir()
	findingPath := ".savepoint/audit/findings/F001-nodate.md"
	writeFiles(t, root, map[string]string{
		findingPath: "---\n" +
			"id: F001\ntitle: No date recorded\nstatus: deferred\nseverity: low\nconfidence: medium\n" +
			"first_seen: 2026-01-01\nproof_needed: n/a\ndeferral_reason: Later.\n" +
			"---\n\n## Summary\n\nNo last_seen field at all.\n",
	})

	p := mustPlan(t, root)
	target, _ := targetByPath(p, findingPath)

	content, err := ConvertIssue(root, p, target)
	if err != nil {
		t.Fatalf("ConvertIssue error = %v", err)
	}
	issue, err := data.DecodeIssueV2(target.TargetPath, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2 error = %v", err)
	}

	if len(issue.History) != 1 {
		t.Fatalf("History = %+v, want one entry", issue.History)
	}
	if !issue.History[0].At.Equal(p.GeneratedAt) {
		t.Errorf("History[0].At = %v, want the migration time %v (no last_seen was recorded)", issue.History[0].At, p.GeneratedAt)
	}
	if !strings.Contains(issue.History[0].Note, "recorded no date") {
		t.Errorf("History[0].Note = %q, want an explicit note that the original date was absent", issue.History[0].Note)
	}
}

// --- task linking -----------------------------------------------------------

// TestConvertIssue_findingLinksConvertedTask proves the fixture's F001 names
// v1.1's T001-shared and that the converted Issue's tasks field carries its
// global Task ID, while a finding naming no task links to nothing.
func TestConvertIssue_findingLinksConvertedTask(t *testing.T) {
	root := fixtureProjectRoot("v1-history")
	p := mustPlan(t, root)

	f001Target, _ := targetByPath(p, ".savepoint/audit/findings/F001-awaiting-proof.md")
	content, err := ConvertIssue(root, p, f001Target)
	if err != nil {
		t.Fatalf("ConvertIssue error = %v", err)
	}
	issue, err := data.DecodeIssueV2(f001Target.TargetPath, content)
	if err != nil {
		t.Fatalf("DecodeIssueV2 error = %v", err)
	}

	taskTarget, ok := targetByPath(p, ".savepoint/releases/v1.1/epics/E01-example/tasks/T001-shared.md")
	if !ok {
		t.Fatalf("expected v1.1's T001-shared to be planned active")
	}
	if len(issue.Tasks) != 1 || issue.Tasks[0] != taskTarget.GlobalID {
		t.Errorf("Tasks = %v, want [%s]", issue.Tasks, taskTarget.GlobalID)
	}

	root2, findingPath2 := findingProject(t, "", "open")
	p2 := mustPlan(t, root2)
	target2, _ := targetByPath(p2, findingPath2)
	content2, err := ConvertIssue(root2, p2, target2)
	if err != nil {
		t.Fatalf("ConvertIssue error = %v", err)
	}
	issue2, err := data.DecodeIssueV2(target2.TargetPath, content2)
	if err != nil {
		t.Fatalf("DecodeIssueV2 error = %v", err)
	}
	if len(issue2.Tasks) != 0 {
		t.Errorf("Tasks = %v, want none: this finding named no task", issue2.Tasks)
	}
}

// --- legacy_fields preservation ---------------------------------------------

// TestConvertIssue_unknownFrontmatterFieldPreserved proves an unrecognized
// finding frontmatter key (like the fixture's reviewer_note) survives under
// legacy_fields rather than being silently dropped.
func TestConvertIssue_unknownFrontmatterFieldPreserved(t *testing.T) {
	root := fixtureProjectRoot("v1-history")
	p := mustPlan(t, root)
	target, ok := targetByPath(p, ".savepoint/audit/findings/F001-awaiting-proof.md")
	if !ok {
		t.Fatalf("F001 was not planned as a target")
	}

	content, err := ConvertIssue(root, p, target)
	if err != nil {
		t.Fatalf("ConvertIssue error = %v", err)
	}
	if !strings.Contains(content, "reviewer_note") {
		t.Errorf("rendered content dropped the unrecognized reviewer_note field")
	}
}

// TestConvertIssue_notAnIssueTarget proves ConvertIssue refuses a target of
// the wrong kind rather than rendering something meaningless.
func TestConvertIssue_notAnIssueTarget(t *testing.T) {
	root := fixtureProjectRoot("v1-basic")
	p := mustPlan(t, root)
	var objectiveTarget PlannedTarget
	for _, tg := range p.Targets {
		if tg.Kind == TargetObjective {
			objectiveTarget = tg
			break
		}
	}
	if _, err := ConvertIssue(root, p, objectiveTarget); err == nil {
		t.Fatal("ConvertIssue() on an Objective target error = nil, want a refusal")
	}
}
