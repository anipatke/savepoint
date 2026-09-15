package data

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// --- v1-history fixture helpers -----------------------------------------
//
// v1-history extends v1-basic's shape: the same epic and short Task ID
// recur across two releases, and a full audit/ register is present. These
// tests reuse the migration fixture helpers defined in
// migration_source_test.go (loadMigrationManifest, readMigrationFile,
// assertFixtureBytesMatchManifest, migrationFixtureDir) rather than
// redefining them.

func rawFrontmatterMap(t *testing.T, content string) map[string]any {
	t.Helper()
	fm, _, err := SplitFrontmatterBody(normalizeLineEndings(content))
	if err != nil {
		t.Fatalf("SplitFrontmatterBody() error = %v", err)
	}
	var m map[string]any
	if err := yaml.Unmarshal([]byte(fm), &m); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}
	return m
}

// loadHistoryTask reads and parses one v1-history task file, then applies the
// explicit release/epic scope discovery would normally supply (mirrors
// migration_source_test.go's TestMigrationSourceBasicInterpretation).
func loadHistoryTask(t *testing.T, relPath, release, epic string) Task {
	t.Helper()
	const fixture = "v1-history"
	raw := readMigrationFile(t, fixture, relPath)
	parser := NewParser()
	task, err := parser.ParseTaskFile(relPath, string(raw))
	if err != nil {
		t.Fatalf("ParseTaskFile(%s) error = %v", relPath, err)
	}
	task.Release = release
	task.Epic = epic
	return *task
}

func loadHistoryDefect(t *testing.T, relPath string) Defect {
	t.Helper()
	const fixture = "v1-history"
	raw := readMigrationFile(t, fixture, relPath)
	parser := NewParser()
	defect, err := parser.ParseDefectFile(relPath, string(raw))
	if err != nil {
		t.Fatalf("ParseDefectFile(%s) error = %v", relPath, err)
	}
	return *defect
}

const (
	relV1T001Shared    = "project/.savepoint/releases/v1/epics/E01-example/tasks/T001-shared.md"
	relV1D001Shared    = "project/.savepoint/releases/v1/defects/D001-shared.md"
	relV11T001Shared   = "project/.savepoint/releases/v1.1/epics/E01-example/tasks/T001-shared.md"
	relV11T002FollowUp = "project/.savepoint/releases/v1.1/epics/E01-example/tasks/T002-follow-up.md"
	relV11D001Shared   = "project/.savepoint/releases/v1.1/defects/D001-shared.md"
	relF001            = "project/.savepoint/audit/findings/F001-awaiting-proof.md"
	relF002            = "project/.savepoint/audit/findings/F002-owner-waiver.md"
	relF003            = "project/.savepoint/audit/findings/F003-duplicate.md"
	relRun2026         = "project/.savepoint/audit/runs/2026-07-01-example.md"
)

// TestMigrationHistoryScopedReferences proves the epic/short-Task-ID collision
// across v1 and v1.1 resolves release-scoped: v1's T001-shared is done, v1.1's
// is in progress, and v1.1's dependent Task resolves its short depends_on
// reference to the in-progress v1.1 record, never the done v1 one. It also
// proves the two D001-shared defect records stay distinct by release.
func TestMigrationHistoryScopedReferences(t *testing.T) {
	v1T001 := loadHistoryTask(t, relV1T001Shared, "v1", "E01-example")
	v11T001 := loadHistoryTask(t, relV11T001Shared, "v1.1", "E01-example")
	v11T002 := loadHistoryTask(t, relV11T002FollowUp, "v1.1", "E01-example")

	if v1T001.Column != ColumnDone {
		t.Errorf("v1 T001-shared.Column = %v, want done", v1T001.Column)
	}
	if v11T001.Column != ColumnInProgress || v11T001.Stage != StageBuild {
		t.Errorf("v1.1 T001-shared.Column/Stage = %v/%v, want in_progress/build", v11T001.Column, v11T001.Stage)
	}
	if v11T002.Column != ColumnPlanned {
		t.Errorf("v1.1 T002-follow-up.Column = %v, want planned", v11T002.Column)
	}
	if len(v11T002.DependsOn) != 1 || v11T002.DependsOn[0] != "T001-shared" {
		t.Fatalf("v1.1 T002-follow-up.DependsOn = %v, want [T001-shared]", v11T002.DependsOn)
	}

	// Both releases' T001-shared share the same short ID; the combined
	// candidate list mirrors what a project-wide (not release-filtered) load
	// would hand to the resolver.
	allTasks := []Task{v1T001, v11T001, v11T002}

	resolution := ResolveDependency(v11T002.DependsOn[0], v11T002, allTasks, map[string]string{})
	if resolution.Kind != DependencyTask {
		t.Fatalf("ResolveDependency() = %+v, want a task resolution", resolution)
	}
	if resolution.ID != v11T001.ID || resolution.TaskStatus != ColumnInProgress {
		t.Errorf("ResolveDependency() = %+v, want v1.1's in_progress T001-shared (%s), not the done v1 record", resolution, v11T001.ID)
	}

	v1D001 := loadHistoryDefect(t, relV1D001Shared)
	v11D001 := loadHistoryDefect(t, relV11D001Shared)
	if v1D001.ID == v11D001.ID {
		t.Fatalf("v1/v1.1 D001-shared share frontmatter id %q; want release-qualified distinct IDs", v1D001.ID)
	}
	if v1D001.Status != DefectResolved {
		t.Errorf("v1 D001-shared.Status = %v, want resolved", v1D001.Status)
	}
	if v11D001.Status != DefectOpen {
		t.Errorf("v1.1 D001-shared.Status = %v, want open", v11D001.Status)
	}
}

// TestMigrationHistoryDispositions proves the audit register's three findings
// carry their required raw dispositions: F001 is fixed but not verified
// (unmet proof), F002 is an owner-waived disposition rather than a technical
// fix, and F003 explicitly points at F001 as its canonical duplicate. It also
// proves the immutable run and epic audit retain their authored content, and
// the project's Health-Check.md carries its short custom procedure.
func TestMigrationHistoryDispositions(t *testing.T) {
	parser := NewParser()

	f1raw := readMigrationFile(t, "v1-history", relF001)
	f1, err := parser.ParseFindingFile(relF001, string(f1raw))
	if err != nil {
		t.Fatalf("ParseFindingFile(F001) error = %v", err)
	}
	if f1.Status != FindingFixed {
		t.Errorf("F001.Status = %v, want fixed", f1.Status)
	}
	if f1.ProofNeeded == "" {
		t.Error("F001.ProofNeeded = empty, want a named regression-test requirement")
	}
	if f1.VerifiedProof != "" {
		t.Errorf("F001.VerifiedProof = %q, want empty: raw fixed must not imply verified", f1.VerifiedProof)
	}

	f2raw := readMigrationFile(t, "v1-history", relF002)
	f2, err := parser.ParseFindingFile(relF002, string(f2raw))
	if err != nil {
		t.Fatalf("ParseFindingFile(F002) error = %v", err)
	}
	if f2.Status != FindingWaived {
		t.Errorf("F002.Status = %v, want waived", f2.Status)
	}
	if f2.WaiverReason == "" {
		t.Error("F002.WaiverReason = empty, want an explicit owner reason")
	}

	f3raw := readMigrationFile(t, "v1-history", relF003)
	f3, err := parser.ParseFindingFile(relF003, string(f3raw))
	if err != nil {
		t.Fatalf("ParseFindingFile(F003) error = %v", err)
	}
	if f3.Status != FindingDuplicate {
		t.Errorf("F003.Status = %v, want duplicate", f3.Status)
	}
	if f3.DuplicateOf != f1.ID {
		t.Errorf("F003.DuplicateOf = %q, want %q", f3.DuplicateOf, f1.ID)
	}

	runRaw := readMigrationFile(t, "v1-history", relRun2026)
	run, err := parser.ParseRunFile(relRun2026, string(runRaw))
	if err != nil {
		t.Fatalf("ParseRunFile() error = %v", err)
	}
	if run.Commit == "" {
		t.Error("run.Commit = empty, want the frozen historical commit reference")
	}
	if run.Label != "example" || run.Date != "2026-07-01" {
		t.Errorf("run.Label/Date = %q/%q, want example/2026-07-01", run.Label, run.Date)
	}
	if !strings.Contains(run.Body, "## Reconciliation") {
		t.Error("run.Body did not retain its authored Reconciliation section")
	}
	if diags := DiagnoseRun(run, relRun2026); len(diags) != 0 {
		t.Errorf("DiagnoseRun() = %v, want none for a well-formed frozen run", diags)
	}

	epicAudit := string(readMigrationFile(t, "v1-history", "project/.savepoint/releases/v1.1/epics/E01-example/E01-Audit.md"))
	if !strings.Contains(epicAudit, "## Main Findings") {
		t.Error("E01-Audit.md lost its authored Main Findings section")
	}

	healthCheck := string(readMigrationFile(t, "v1-history", "project/.savepoint/Health-Check.md"))
	if !strings.Contains(healthCheck, "## Quick Check") {
		t.Error("Health-Check.md lost its Quick Check section")
	}
	if strings.Contains(healthCheck, "## Full Check") || strings.Contains(healthCheck, "## Deep Check") {
		t.Error("Health-Check.md should carry only the short custom Quick procedure, not the full template shape")
	}

	assertHistoryInventoryMatchesManifest(t)
}

// assertHistoryInventoryMatchesManifest diffs the fixture's on-disk file set
// against manifest.yml's recorded paths in both directions (mirrors
// TestMigrationSourceBasicInventory), and confirms the one file this
// project-variant intentionally lacks is absent by design.
func assertHistoryInventoryMatchesManifest(t *testing.T) {
	t.Helper()
	const fixture = "v1-history"
	manifest := loadMigrationManifest(t, fixture)
	if len(manifest.Files) == 0 {
		t.Fatal("manifest.yml lists no files")
	}
	assertFixtureBytesMatchManifest(t, fixture)

	recorded := map[string]bool{}
	for _, f := range manifest.Files {
		recorded[filepath.ToSlash(f.Path)] = true
	}

	root := migrationFixtureDir(fixture)
	var onDisk []string
	err := filepath.Walk(filepath.Join(root, "project"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		onDisk = append(onDisk, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walk fixture project tree: %v", err)
	}

	if len(onDisk) != len(recorded) {
		t.Errorf("fixture has %d files on disk, manifest records %d: on disk = %v", len(onDisk), len(recorded), onDisk)
	}
	for _, rel := range onDisk {
		if !recorded[rel] {
			t.Errorf("fixture file %s is on disk but not recorded in manifest.yml", rel)
		}
	}

	for _, a := range manifest.AbsentByDesign {
		if _, err := os.Stat(filepath.Join(root, a.Path)); !os.IsNotExist(err) {
			t.Errorf("expected %s to be absent by design (%s), stat err = %v", a.Path, a.Reason, err)
		}
	}
}

// TestMigrationHistoryRawAndNormalized shows raw source evidence surviving
// even when the typed loader either drops an unknown field or heals an
// invalid one. It also exercises LoadAuditRegisterSet across the whole
// audit/ tree and shows it retains all three findings' canonical
// distinctions rather than collapsing them.
func TestMigrationHistoryRawAndNormalized(t *testing.T) {
	f1raw := string(readMigrationFile(t, "v1-history", relF001))

	// The frontmatter carries an unknown field that has no place in the
	// typed AuditFinding struct; it must still round-trip through raw YAML.
	rawFields := rawFrontmatterMap(t, f1raw)
	const wantNote = "kept from v1 audit tooling; not part of the current finding schema"
	if rawFields["reviewer_note"] != wantNote {
		t.Errorf("raw F001 frontmatter reviewer_note = %v, want %q", rawFields["reviewer_note"], wantNote)
	}

	parser := NewParser()
	f1, err := parser.ParseFindingFile(relF001, f1raw)
	if err != nil {
		t.Fatalf("ParseFindingFile(F001) error = %v", err)
	}
	// The typed model has no field for it; NormalizeFindingForLoad does not
	// invent one, so this is purely evidence the raw map above carried it.
	_ = f1

	// A temporary unknown-status variant demonstrates raw preservation
	// alongside the current loader's open default, without rewriting the
	// frozen fixture.
	unknownStatus := strings.Replace(f1raw, "status: fixed", "status: escalated", 1)
	if unknownStatus == f1raw {
		t.Fatal("test setup: status field not found in F001 fixture")
	}
	dest := filepath.Join(t.TempDir(), "F001-unknown-status.md")
	if err := os.WriteFile(dest, []byte(unknownStatus), 0644); err != nil {
		t.Fatalf("write unknown-status copy: %v", err)
	}
	unknownContent, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("read unknown-status copy: %v", err)
	}

	rawVariant, err := parser.ParseRawFindingFile(dest, string(unknownContent))
	if err != nil {
		t.Fatalf("ParseRawFindingFile() error = %v", err)
	}
	if rawVariant.Status != "escalated" {
		t.Errorf("ParseRawFindingFile().Status = %q, want the raw unrecognized value preserved", rawVariant.Status)
	}

	healedVariant, err := parser.ParseFindingFile(dest, string(unknownContent))
	if err != nil {
		t.Fatalf("ParseFindingFile() error = %v", err)
	}
	if healedVariant.Status != FindingOpen {
		t.Errorf("ParseFindingFile().Status = %q, want healed to open", healedVariant.Status)
	}

	diags := DiagnoseFinding(rawVariant, dest)
	if !hasDiagnostic(diags, FindingInvalidStatusCode) {
		t.Errorf("DiagnoseFinding() = %v, want an invalid_status diagnostic for the unknown-status variant", diags)
	}

	// LoadAuditRegisterSet must retain all three findings distinctly, sorted
	// by the canonical lifecycle order.
	root := filepath.Join(migrationFixtureDir("v1-history"), "project", ".savepoint")
	set, err := LoadAuditRegisterSet(root)
	if err != nil {
		t.Fatalf("LoadAuditRegisterSet() error = %v", err)
	}
	if !set.Prompt.Available || set.Prompt.Version != "v1" {
		t.Errorf("prompt available/version = %v/%q, want true/v1", set.Prompt.Available, set.Prompt.Version)
	}
	if !set.Register.Available || !set.Register.HasSummary {
		t.Errorf("register available/hasSummary = %v/%v, want true/true", set.Register.Available, set.Register.HasSummary)
	}
	if len(set.Findings) != 3 {
		t.Fatalf("len(set.Findings) = %d, want 3", len(set.Findings))
	}
	byID := map[string]AuditFinding{}
	for _, f := range set.Findings {
		byID[f.ID] = f
	}
	if byID["F001"].Status != FindingFixed || byID["F002"].Status != FindingWaived || byID["F003"].Status != FindingDuplicate {
		t.Errorf("set.Findings statuses = F001:%v F002:%v F003:%v, want fixed/waived/duplicate",
			byID["F001"].Status, byID["F002"].Status, byID["F003"].Status)
	}
	if len(set.Runs) != 1 || set.Runs[0].Label != "example" {
		t.Errorf("set.Runs = %v, want one run labeled example", set.Runs)
	}
}

// TestMigrationHistoryFailures exercises the reader's behavior on malformed
// and unresolved-reference variants using temporary copies only; the frozen
// fixture bytes must be unchanged before and after.
// TestMigrationHistoryLoadProjectDispatch proves the E42 schema-dispatch
// boundary (LoadProject) reaches the frozen v1-history fixture, with its two
// releases, through transitional V1 dispatch, and that doing so never
// touches the frozen source bytes.
func TestMigrationHistoryLoadProjectDispatch(t *testing.T) {
	const fixture = "v1-history"
	savepointRoot := filepath.Join(migrationFixtureDir(fixture), "project", ".savepoint")

	project, err := LoadProject(savepointRoot)
	if err != nil {
		t.Fatalf("LoadProject() error = %v", err)
	}
	if project.SchemaVersion != SchemaVersionV1 {
		t.Fatalf("LoadProject() SchemaVersion = %v, want SchemaVersionV1", project.SchemaVersion)
	}
	if project.V1 == nil {
		t.Fatal("LoadProject() V1 discover adapter = nil, want non-nil for transitional V1 dispatch")
	}

	releases, err := project.V1.ListReleases(savepointRoot)
	if err != nil {
		t.Fatalf("project.V1.ListReleases() error = %v", err)
	}
	if len(releases) != 2 || releases[0].ID != "v1" || releases[1].ID != "v1.1" {
		t.Fatalf("project.V1.ListReleases() = %v, want [v1 v1.1]", releases)
	}

	assertFixtureBytesMatchManifest(t, fixture)
}

func TestMigrationHistoryFailures(t *testing.T) {
	parser := NewParser()

	t.Run("malformed finding frontmatter produces the existing named parse error", func(t *testing.T) {
		raw := string(readMigrationFile(t, "v1-history", relF001))
		corrupted := strings.Replace(raw, "status: fixed", "status: [unterminated", 1)
		if corrupted == raw {
			t.Fatal("test setup: corruption target line not found in F001 fixture")
		}
		dest := filepath.Join(t.TempDir(), "F001-corrupt.md")
		if err := os.WriteFile(dest, []byte(corrupted), 0644); err != nil {
			t.Fatalf("write corrupted copy: %v", err)
		}
		content, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("read corrupted copy: %v", err)
		}
		if _, err := parser.ParseFindingFile(dest, string(content)); err == nil {
			t.Fatal("ParseFindingFile() error = nil, want a malformed-YAML parse error")
		} else if !strings.Contains(err.Error(), "parse error for") {
			t.Errorf("ParseFindingFile() error = %v, want the existing named parse-error shape", err)
		}
	})

	t.Run("malformed run frontmatter produces the existing named parse error", func(t *testing.T) {
		raw := string(readMigrationFile(t, "v1-history", relRun2026))
		corrupted := strings.Replace(raw, "mode: full", "mode: [unterminated", 1)
		if corrupted == raw {
			t.Fatal("test setup: corruption target line not found in run fixture")
		}
		dest := filepath.Join(t.TempDir(), "2026-07-01-corrupt.md")
		if err := os.WriteFile(dest, []byte(corrupted), 0644); err != nil {
			t.Fatalf("write corrupted copy: %v", err)
		}
		content, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("read corrupted copy: %v", err)
		}
		if _, err := parser.ParseRunFile(dest, string(content)); err == nil {
			t.Fatal("ParseRunFile() error = nil, want a malformed-YAML parse error")
		} else if !strings.Contains(err.Error(), "parse error for") {
			t.Errorf("ParseRunFile() error = %v, want the existing named parse-error shape", err)
		}
	})

	t.Run("unresolved same-epic dependency reference resolves to the existing missing result", func(t *testing.T) {
		v1T001 := loadHistoryTask(t, relV1T001Shared, "v1", "E01-example")
		v11T001 := loadHistoryTask(t, relV11T001Shared, "v1.1", "E01-example")

		raw := string(readMigrationFile(t, "v1-history", relV11T002FollowUp))
		broken := strings.Replace(raw, "depends_on: [T001-shared]", "depends_on: [T099-missing]", 1)
		if broken == raw {
			t.Fatal("test setup: depends_on line not found in T002-follow-up fixture")
		}
		dest := filepath.Join(t.TempDir(), "T002-broken-dep.md")
		if err := os.WriteFile(dest, []byte(broken), 0644); err != nil {
			t.Fatalf("write broken-dependency copy: %v", err)
		}
		content, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("read broken-dependency copy: %v", err)
		}
		t2broken, err := parser.ParseTaskFile(dest, string(content))
		if err != nil {
			t.Fatalf("ParseTaskFile() error = %v, want the broken-dependency copy to still parse", err)
		}
		t2broken.Release, t2broken.Epic = "v1.1", "E01-example"

		resolution := ResolveDependency(t2broken.DependsOn[0], *t2broken, []Task{v1T001, v11T001, *t2broken}, map[string]string{})
		if resolution != (DependencyResolution{}) {
			t.Errorf("ResolveDependency() = %+v, want the existing zero-value missing-reference result", resolution)
		}
	})

	// All three failure variants ran against temporary copies only; the
	// frozen fixture bytes must be exactly what the manifest recorded.
	assertFixtureBytesMatchManifest(t, "v1-history")
}
