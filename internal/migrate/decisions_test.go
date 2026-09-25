package migrate

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writeMinimalV1Project writes the smallest config.yml, router.md, and live
// Release PRD every ad hoc temp-dir test in this package builds on, so each
// test below only authors the source files its own scenario actually needs.
func writeMinimalV1Project(t *testing.T, root string) {
	t.Helper()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "v1-PRD.md"), "---\nstatus: in_progress\n---\n\n# V1\n")
}

func writeTaskWithStatus(t *testing.T, root, release, epic, taskID, status, dependsOn string) string {
	t.Helper()
	epicDir := filepath.Join(root, ".savepoint", "releases", release, "epics", epic)
	prefix := epic
	if i := strings.IndexByte(epic, '-'); i >= 0 {
		prefix = epic[:i]
	}
	writeFile(t, filepath.Join(epicDir, prefix+"-Detail.md"), "---\nstatus: in_progress\n---\n\n# "+epic+"\n")
	taskPath := ".savepoint/releases/" + release + "/epics/" + epic + "/tasks/" + taskID + ".md"
	depends := "[]"
	if dependsOn != "" {
		depends = "[" + dependsOn + "]"
	}
	writeFile(t, filepath.Join(root, filepath.FromSlash(taskPath)),
		"---\nid: "+epic+"/"+taskID+"\nstatus: "+status+"\ndepends_on: "+depends+"\n---\n\n# "+taskID+"\n")
	return taskPath
}

func writeFinding(t *testing.T, root, id, title, status, duplicateOf string) string {
	t.Helper()
	fm := "---\nid: " + id + "\ntitle: " + title + "\nstatus: " + status + "\nseverity: medium\nconfidence: medium\n"
	if duplicateOf != "" {
		fm += "duplicate_of: " + duplicateOf + "\n"
	}
	fm += "---\n\n## Summary\n\n" + title + "\n"
	path := ".savepoint/audit/findings/" + id + "-finding.md"
	writeFile(t, filepath.Join(root, filepath.FromSlash(path)), fm)
	return path
}

// --- AC1/AC2: stable, self-describing ambiguities --------------------------

func TestAmbiguity_idStableAcrossRepeatedPreviews(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")

	first := mustPlan(t, root)
	second := mustPlan(t, root)

	a1, ok := findAmbiguityByPath(first, taskPath)
	if !ok {
		t.Fatalf("first Plan() raised no ambiguity for %s", taskPath)
	}
	a2, ok := findAmbiguityByPath(second, taskPath)
	if !ok {
		t.Fatalf("second Plan() raised no ambiguity for %s", taskPath)
	}
	if a1.ID != a2.ID {
		t.Errorf("ambiguity ID changed across repeated previews: %q vs %q", a1.ID, a2.ID)
	}
}

func findAmbiguityByPath(p *ConversionPlan, path string) (Ambiguity, bool) {
	for _, a := range p.Ambiguities {
		if a.Path == path {
			return a, true
		}
	}
	return Ambiguity{}, false
}

func TestAmbiguity_statesPathKindDetailAndChoices(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")

	p := mustPlan(t, root)
	a, ok := findAmbiguityByPath(p, taskPath)
	if !ok {
		t.Fatalf("no ambiguity raised for %s", taskPath)
	}
	if a.Kind != AmbiguityUnrecognizedLifecycle {
		t.Errorf("Kind = %v, want AmbiguityUnrecognizedLifecycle", a.Kind)
	}
	if a.Detail == "" {
		t.Error("Detail is empty, want an explanation of why this is ambiguous")
	}
	if len(a.Choices) == 0 {
		t.Error("Choices is empty, want the concrete values that would resolve it")
	}
	if !containsString(a.Choices, "planned") {
		t.Errorf("Choices = %v, want the canonical task statuses including planned", a.Choices)
	}
}

// --- AC3: blocking coverage --------------------------------------------------

func TestPlan_duplicateSourceIdentity_isBlockingAmbiguity(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	epicDir := filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-x")
	writeFile(t, filepath.Join(epicDir, "E01-Detail.md"), "---\nstatus: in_progress\n---\n\n# E01\n")

	firstPath := ".savepoint/releases/v1/epics/E01-x/tasks/T001-a.md"
	secondPath := ".savepoint/releases/v1/epics/E01-x/tasks/T001-b.md"
	writeFile(t, filepath.Join(root, filepath.FromSlash(firstPath)),
		"---\nid: E01-x/T001-dup\nstatus: planned\ndepends_on: []\n---\n\n# T001 (a)\n")
	writeFile(t, filepath.Join(root, filepath.FromSlash(secondPath)),
		"---\nid: E01-x/T001-dup\nstatus: planned\ndepends_on: []\n---\n\n# T001 (b)\n")

	p := mustPlan(t, root)

	a, ok := findAmbiguityByPath(p, secondPath)
	if !ok {
		t.Fatalf("no ambiguity raised for the second declaration of a duplicated task id; Ambiguities = %+v", p.Ambiguities)
	}
	if a.Kind != AmbiguityDuplicateSourceIdentity {
		t.Errorf("Kind = %v, want AmbiguityDuplicateSourceIdentity", a.Kind)
	}
	if !a.Blocking {
		t.Error("Blocking = false, want true: a duplicate source identity must block apply")
	}
	if _, ok := targetByPath(p, secondPath); ok {
		t.Error("the duplicated second declaration was planned as a target, want it withheld pending a decision")
	}
	if p.Appliable {
		t.Error("Appliable = true, want false: an unresolved blocking ambiguity is present")
	}
	if len(p.UnresolvedBlockingIDs) != 1 || p.UnresolvedBlockingIDs[0] != a.ID {
		t.Errorf("UnresolvedBlockingIDs = %v, want exactly [%s]", p.UnresolvedBlockingIDs, a.ID)
	}
}

func TestPlan_unresolvedNarrativeFinding_isBlockingAmbiguity(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	path := writeFinding(t, root, "F099", "Orphaned duplicate report", "duplicate", "F404")

	p := mustPlan(t, root)

	a, ok := findAmbiguityByPath(p, path)
	if !ok {
		t.Fatalf("no ambiguity raised for a duplicate finding naming an unknown canonical; Ambiguities = %+v", p.Ambiguities)
	}
	if a.Kind != AmbiguityUnresolvedNarrativeFind {
		t.Errorf("Kind = %v, want AmbiguityUnresolvedNarrativeFind", a.Kind)
	}
	if !a.Blocking {
		t.Error("Blocking = false, want true")
	}
	if _, ok := targetByPath(p, path); ok {
		t.Error("finding with an unresolved duplicate_of was planned as a target, want it blocked instead")
	}
	if len(a.Choices) == 0 {
		t.Error("Choices is empty, want concrete resolutions (e.g. treat_as_original, archive)")
	}
}

// --- AC4: advisory does not block -------------------------------------------

func TestPlan_unclassifiedFile_isAdvisoryNotBlocking(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	junkPath := ".savepoint/notes-nobody-recognizes.txt"
	writeFile(t, filepath.Join(root, filepath.FromSlash(junkPath)), "scratch notes\n")

	p := mustPlan(t, root)

	a, ok := findAmbiguityByPath(p, junkPath)
	if !ok {
		t.Fatalf("no ambiguity reported for an unclassified file; Ambiguities = %+v", p.Ambiguities)
	}
	if a.Kind != AmbiguityUnclassifiedFile {
		t.Errorf("Kind = %v, want AmbiguityUnclassifiedFile", a.Kind)
	}
	if a.Blocking {
		t.Error("Blocking = true, want false: an unclassified file that is archived intact must not block")
	}
	if _, ok := archiveByPath(p, junkPath); !ok {
		t.Error("unclassified file was not archived, want it archived intact regardless of the advisory ambiguity")
	}
	if !p.Appliable {
		t.Errorf("Appliable = false, want true: only a blocking ambiguity should affect appliability; UnresolvedBlockingIDs = %v", p.UnresolvedBlockingIDs)
	}
}

// --- AC5/AC6/AC7: decisions resolve, validate, and gate appliability -------

func TestPlan_decisionResolvesBlockingAmbiguity_andMakesPlanAppliable(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")

	before := mustPlan(t, root)
	a, ok := findAmbiguityByPath(before, taskPath)
	if !ok {
		t.Fatalf("no ambiguity raised for %s", taskPath)
	}
	if before.Appliable {
		t.Fatal("Appliable = true before any decision, want false")
	}

	at := time.Date(2026, 9, 18, 9, 0, 0, 0, time.UTC)
	decisions := Decisions{a.ID: DecisionValue{Value: "planned", SourceFile: "decisions.yml", DecidedAt: at}}

	after, err := Plan(root, decisions, fixedClock(time.Now()), fixedOperationID("op"))
	if err != nil {
		t.Fatalf("Plan() with a valid decision error = %v", err)
	}
	resolved, ok := findAmbiguityByPath(after, taskPath)
	if !ok {
		t.Fatalf("ambiguity vanished after a decision was supplied, want it still reported as resolved")
	}
	if !resolved.Resolved {
		t.Error("Resolved = false, want true once a valid decision names this ambiguity")
	}
	if resolved.Decision != "planned" {
		t.Errorf("Decision = %q, want %q", resolved.Decision, "planned")
	}
	if resolved.DecisionSourceFile != "decisions.yml" || !resolved.DecisionAt.Equal(at) {
		t.Errorf("decision provenance = (%q, %v), want (%q, %v)", resolved.DecisionSourceFile, resolved.DecisionAt, "decisions.yml", at)
	}
	if !after.Appliable {
		t.Errorf("Appliable = false, want true once every blocking ambiguity is resolved; UnresolvedBlockingIDs = %v", after.UnresolvedBlockingIDs)
	}
	if len(after.UnresolvedBlockingIDs) != 0 {
		t.Errorf("UnresolvedBlockingIDs = %v, want none", after.UnresolvedBlockingIDs)
	}
}

func TestPlan_decisionsFile_unknownID_isError(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")

	decisions := Decisions{"unrecognized_lifecycle:no-such-path.md": {Value: "planned"}}
	_, err := Plan(root, decisions, fixedClock(time.Now()), fixedOperationID("op"))
	if !errors.Is(err, ErrUnknownAmbiguityID) {
		t.Fatalf("Plan() error = %v, want ErrUnknownAmbiguityID", err)
	}
}

func TestPlan_decisionValueOutsideChoices_isError(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")

	before := mustPlan(t, root)
	a, ok := findAmbiguityByPath(before, taskPath)
	if !ok {
		t.Fatalf("no ambiguity raised for %s", taskPath)
	}

	decisions := Decisions{a.ID: {Value: "orbiting"}}
	_, err := Plan(root, decisions, fixedClock(time.Now()), fixedOperationID("op"))
	if !errors.Is(err, ErrDecisionValueNotAllowed) {
		t.Fatalf("Plan() error = %v, want ErrDecisionValueNotAllowed", err)
	}
}

// --- AC8: decisions recorded into the manifest with provenance --------------

func TestBuildManifest_recordsDecisionProvenance(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")

	before := mustPlan(t, root)
	a, ok := findAmbiguityByPath(before, taskPath)
	if !ok {
		t.Fatalf("no ambiguity raised for %s", taskPath)
	}

	at := time.Date(2026, 9, 18, 9, 0, 0, 0, time.UTC)
	decisions := Decisions{a.ID: {Value: "planned", SourceFile: "decisions.yml", DecidedAt: at}}
	p, err := Plan(root, decisions, fixedClock(time.Now()), fixedOperationID("op"))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	m := BuildManifest(p)
	if len(m.Decisions) != 1 {
		t.Fatalf("manifest Decisions = %+v, want exactly one", m.Decisions)
	}
	got := m.Decisions[0]
	if got.AmbiguityID != a.ID || got.Decision != "planned" || got.SourceFile != "decisions.yml" || !got.DecidedAt.Equal(at) {
		t.Errorf("manifest decision = %+v, want ambiguity %s = planned from decisions.yml at %v", got, a.ID, at)
	}
}

// --- AC9: reading a decisions file, and previewing with one, writes nothing -

func TestReadDecisionsFile_readsWithoutWriting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "decisions.yml")
	writeFile(t, path, "decisions:\n  - id: \"unrecognized_lifecycle:some/path.md\"\n    value: planned\n")

	before := snapshotTree(t, dir)
	at := time.Date(2026, 9, 18, 9, 0, 0, 0, time.UTC)

	decisions, err := ReadDecisionsFile(path, at)
	if err != nil {
		t.Fatalf("ReadDecisionsFile() error = %v", err)
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, dir))

	dv, ok := decisions["unrecognized_lifecycle:some/path.md"]
	if !ok {
		t.Fatalf("decisions = %+v, want the one entry the file names", decisions)
	}
	if dv.Value != "planned" || dv.SourceFile != path || !dv.DecidedAt.Equal(at) {
		t.Errorf("decision = %+v, want value planned, source %s, decided %v", dv, path, at)
	}
}

func TestPlan_writeFreeWithDecisionsSupplied(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")

	before := mustPlan(t, root)
	a, ok := findAmbiguityByPath(before, taskPath)
	if !ok {
		t.Fatalf("no ambiguity raised for %s", taskPath)
	}
	decisions := Decisions{a.ID: {Value: "planned", SourceFile: "decisions.yml", DecidedAt: time.Now()}}

	snapshot := snapshotTree(t, root)
	if _, err := Plan(root, decisions, fixedClock(time.Now()), fixedOperationID("op")); err != nil {
		t.Fatalf("Plan() with decisions supplied error = %v", err)
	}
	assertSnapshotsEqual(t, snapshot, snapshotTree(t, root))
}

// --- AC10: never resolved by default, heuristic, or similarity -------------

func TestPlan_similarFindingsAreNeverAutoMerged(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	const title = "Token leak in request logs"
	f1 := writeFinding(t, root, "F001", title, "open", "")
	f2 := writeFinding(t, root, "F002", title, "open", "")

	p := mustPlan(t, root)

	t1, ok := targetByPath(p, f1)
	if !ok || t1.Kind != TargetIssue {
		t.Fatalf("F001 was not planned as its own Issue: %+v", t1)
	}
	t2, ok := targetByPath(p, f2)
	if !ok || t2.Kind != TargetIssue {
		t.Fatalf("F002 was not planned as its own Issue: %+v", t2)
	}
	if t1.GlobalID == t2.GlobalID {
		t.Errorf("F001 and F002 share global ID %q despite no duplicate_of link; identical text must never be auto-merged", t1.GlobalID)
	}
	if t2.DuplicateOfGlobalID != "" {
		t.Errorf("F002.DuplicateOfGlobalID = %q, want empty: nothing named it a duplicate", t2.DuplicateOfGlobalID)
	}
}
