package migrate

import (
	"errors"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
)

func fixedClock(t time.Time) Clock { return func() time.Time { return t } }

func fixedOperationID(id string) OperationIDSource { return func() string { return id } }

func mustPlan(t *testing.T, projectRoot string) *ConversionPlan {
	t.Helper()
	p, err := Plan(projectRoot, nil, fixedClock(time.Date(2026, 9, 18, 0, 0, 0, 0, time.UTC)), fixedOperationID("op-test"))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	return p
}

func targetByPath(p *ConversionPlan, path string) (PlannedTarget, bool) {
	for _, t := range p.Targets {
		if t.Legacy.Path == path {
			return t, true
		}
	}
	return PlannedTarget{}, false
}

func archiveByPath(p *ConversionPlan, path string) (ArchiveEntry, bool) {
	for _, a := range p.Archives {
		if a.SourcePath == path {
			return a, true
		}
	}
	return ArchiveEntry{}, false
}

// --- write-free -------------------------------------------------------

// TestPlan_performsNoWrite snapshots both frozen fixtures' bytes and mtimes,
// runs Plan, and asserts nothing changed: Plan opens no file for writing and
// creates no directory.
func TestPlan_performsNoWrite(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			root := fixtureProjectRoot(fixture)
			before := snapshotTree(t, root)

			if _, err := Plan(root, nil, fixedClock(time.Now()), fixedOperationID("op-1")); err != nil {
				t.Fatalf("Plan(%s) error = %v", fixture, err)
			}

			assertSnapshotsEqual(t, before, snapshotTree(t, root))
		})
	}
}

// --- determinism --------------------------------------------------------

// TestPlan_deterministicAcrossRepeatedRuns proves that with an injected
// clock and operation ID, planning the same fixture twice produces an
// identical plan: identical targets, identical allocated IDs, identical
// ordering.
func TestPlan_deterministicAcrossRepeatedRuns(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			root := fixtureProjectRoot(fixture)
			clock := fixedClock(time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC))
			opID := fixedOperationID("op-fixed")

			first, err := Plan(root, nil, clock, opID)
			if err != nil {
				t.Fatalf("Plan() first call error = %v", err)
			}
			second, err := Plan(root, nil, clock, opID)
			if err != nil {
				t.Fatalf("Plan() second call error = %v", err)
			}

			m1, err := BuildManifest(first).Marshal()
			if err != nil {
				t.Fatalf("Marshal() first error = %v", err)
			}
			m2, err := BuildManifest(second).Marshal()
			if err != nil {
				t.Fatalf("Marshal() second error = %v", err)
			}
			if string(m1) != string(m2) {
				t.Errorf("Plan(%s) produced different manifests across repeated runs:\n--- first ---\n%s\n--- second ---\n%s", fixture, m1, m2)
			}
		})
	}
}

// --- v1-basic: archived-done-task + legacy prerequisite -----------------

func TestPlan_v1Basic_archivesCompletedTaskAndReservesNoID(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-basic"))

	const t1Path = ".savepoint/releases/v1/epics/E01-example/tasks/T001-original.md"
	const t2Path = ".savepoint/releases/v1/epics/E01-example/tasks/T002-follow-up.md"

	if _, ok := targetByPath(p, t1Path); ok {
		t.Errorf("done task %s was planned as a target, want archive-only", t1Path)
	}
	archived, ok := archiveByPath(p, t1Path)
	if !ok {
		t.Fatalf("done task %s was not archived", t1Path)
	}
	if archived.Role != RoleTask {
		t.Errorf("archived %s Role = %v, want RoleTask", t1Path, archived.Role)
	}
	if archived.ArchivePath != ".savepoint/archive/v1/"+t1Path {
		t.Errorf("archived %s ArchivePath = %q, want .savepoint/archive/v1/%s", t1Path, archived.ArchivePath, t1Path)
	}

	active, ok := targetByPath(p, t2Path)
	if !ok {
		t.Fatalf("in-progress task %s was not planned as an active target", t2Path)
	}
	if active.Kind != TargetTask {
		t.Errorf("active target Kind = %v, want TargetTask", active.Kind)
	}
	if active.GlobalID != "T-001" {
		t.Errorf("active task GlobalID = %q, want T-001 (the done sibling reserves nothing)", active.GlobalID)
	}
	if len(active.DependsOn) != 0 {
		t.Errorf("active task DependsOn = %v, want empty: its dependency is archived, not converted", active.DependsOn)
	}

	if len(p.Prereqs) != 1 {
		t.Fatalf("Prereqs = %+v, want exactly one legacy prerequisite", p.Prereqs)
	}
	prereq := p.Prereqs[0]
	if prereq.Task != active.GlobalID {
		t.Errorf("Prereqs[0].Task = %q, want %q", prereq.Task, active.GlobalID)
	}
	if prereq.ArchivePath != archived.ArchivePath {
		t.Errorf("Prereqs[0].ArchivePath = %q, want %q", prereq.ArchivePath, archived.ArchivePath)
	}
	if prereq.Evidence == "" {
		t.Error("Prereqs[0].Evidence is empty, want the original recorded completion evidence")
	}

	// The active epic itself becomes an Objective; the epic's own status
	// being in_progress (not done) makes it the Objective owner even though
	// one of its tasks is already archived.
	const epicPath = ".savepoint/releases/v1/epics/E01-example/E01-Detail.md"
	objective, ok := targetByPath(p, epicPath)
	if !ok || objective.Kind != TargetObjective {
		t.Fatalf("epic %s was not planned as an Objective target: %+v", epicPath, objective)
	}
	if objective.GlobalID != "O-001" {
		t.Errorf("Objective GlobalID = %q, want O-001", objective.GlobalID)
	}
	// Original epic bytes are also archived alongside conversion.
	if _, ok := archiveByPath(p, epicPath); !ok {
		t.Errorf("epic detail %s was not archived alongside its Objective conversion", epicPath)
	}
}

func TestPlan_routerSelectsLiveGoalOrPlansContinuationGoal(t *testing.T) {
	tests := []struct {
		name          string
		routerRelease string
		releaseStatus string
		wantGenerated bool
		wantReason    string
	}{
		{name: "live selection is kept", routerRelease: "v1", releaseStatus: "in_progress"},
		{name: "missing selection", routerRelease: "", releaseStatus: "in_progress", wantReason: "only existing live Goal"},
		{name: "unresolvable selection", routerRelease: "v9-missing", releaseStatus: "in_progress", wantReason: "only existing live Goal"},
		{name: "archived selection", routerRelease: "v1", releaseStatus: "done", wantGenerated: true, wantReason: "historical Goal"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
			writeFile(t, filepath.Join(root, ".savepoint", "router.md"), routerFixtureContent(
				"task-building", tc.routerRelease, "", "", `"Continue migrated work."`))
			writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "v1-PRD.md"),
				"---\nname: V1\nstatus: "+tc.releaseStatus+"\n---\n\n# V1\n")

			p := mustPlan(t, root)
			selection := p.GoalSelection
			if selection.GoalID == "" {
				t.Fatal("GoalSelection.GoalID is empty")
			}
			if selection.Generated != tc.wantGenerated {
				t.Fatalf("GoalSelection.Generated = %t, want %t: %+v", selection.Generated, tc.wantGenerated, selection)
			}
			if tc.wantReason != "" && !strings.Contains(selection.Reason, tc.wantReason) {
				t.Errorf("GoalSelection.Reason = %q, want it to contain %q", selection.Reason, tc.wantReason)
			}

			generated := 0
			for _, target := range p.Targets {
				if target.Generated {
					generated++
					if target.GlobalID != selection.GoalID || target.Kind != TargetRelease || target.ReleaseStatus != string(data.ColumnInProgress) {
						t.Errorf("generated target = %+v, want the selected live Goal", target)
					}
					if target.Legacy.Path != "" {
						t.Errorf("generated Goal Legacy.Path = %q, want no V1 source", target.Legacy.Path)
					}
				}
			}
			wantGeneratedCount := 0
			if tc.wantGenerated {
				wantGeneratedCount = 1
				if selection.GoalID != "R-002" {
					t.Errorf("fallback GoalID = %q, want next free R-002", selection.GoalID)
				}
			} else {
				if selection.GoalID != "R-001" {
					t.Errorf("kept GoalID = %q, want mapped V1 Goal R-001", selection.GoalID)
				}
			}
			if generated != wantGeneratedCount {
				t.Errorf("generated Goal targets = %d, want %d", generated, wantGeneratedCount)
			}

			preview := FormatPreview(p)
			if tc.wantReason != "" && (!strings.Contains(preview, tc.wantReason) || !strings.Contains(preview, selection.GoalID)) {
				t.Errorf("preview does not explain Goal selection %s (%s):\n%s", selection.GoalID, tc.wantReason, preview)
			}
			if strings.Contains(preview, "router now selects no Release") {
				t.Errorf("preview retains obsolete no-Release outcome:\n%s", preview)
			}

			manifest := BuildManifest(p)
			for _, identity := range manifest.Identities {
				if identity.GlobalID == selection.GoalID && tc.wantGenerated {
					t.Errorf("generated Goal has a fabricated V1 identity mapping: %+v", identity)
				}
			}
		})
	}
}

func TestPlan_activeEpicWithoutResolvableReleaseFailsWithNamedError(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), routerFixtureContent(
		"epic-task-breakdown", "v1", "E01-orphan", "", `"Plan the epic."`))
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-orphan", "E01-Detail.md"),
		"---\nstatus: in_progress\n---\n\n# E01 Orphan\n")

	_, err := Plan(root, nil, fixedClock(time.Now()), fixedOperationID("op"))
	if err == nil {
		t.Fatal("Plan() error = nil, want a named failure for the active epic's missing Goal")
	}
	if !strings.Contains(err.Error(), "V1 epic v1/E01-orphan") || !strings.Contains(err.Error(), "no resolvable V2 Goal") {
		t.Fatalf("Plan() error = %q, want it to identify the V1 epic and missing Goal", err)
	}
}

// --- v1-history: source-qualified allocation -----------------------------

// TestPlan_v1History_sourceQualifiedAllocation proves the epic/short-task-ID
// collision across v1 and v1.1 resolves to two distinct global Task IDs, and
// that v1's fully-done epic is archived intact rather than given an
// Objective, while v1.1's active epic is not.
func TestPlan_v1History_sourceQualifiedAllocation(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-history"))

	const v1T001 = ".savepoint/releases/v1/epics/E01-example/tasks/T001-shared.md"
	const v11T001 = ".savepoint/releases/v1.1/epics/E01-example/tasks/T001-shared.md"
	const v11T002 = ".savepoint/releases/v1.1/epics/E01-example/tasks/T002-follow-up.md"

	if _, ok := targetByPath(p, v1T001); ok {
		t.Errorf("v1's done T001-shared was planned as an active target, want archive-only (its epic is done)")
	}
	v1Archived, ok := archiveByPath(p, v1T001)
	if !ok {
		t.Fatalf("v1's T001-shared was not archived")
	}

	v11Active, ok := targetByPath(p, v11T001)
	if !ok {
		t.Fatalf("v1.1's in-progress T001-shared was not planned as an active target")
	}
	if v11Active.GlobalID == "" {
		t.Fatal("v1.1 T001-shared has no allocated global ID")
	}

	followUp, ok := targetByPath(p, v11T002)
	if !ok {
		t.Fatalf("v1.1's T002-follow-up was not planned as an active target")
	}
	if len(followUp.DependsOn) != 1 || followUp.DependsOn[0] != v11Active.GlobalID {
		t.Errorf("T002-follow-up DependsOn = %v, want [%s] (v1.1's own T001-shared, never the archived v1 one)", followUp.DependsOn, v11Active.GlobalID)
	}
	if len(p.Prereqs) != 0 {
		t.Errorf("Prereqs = %v, want none: T002-follow-up's dependency resolved to an active task, not an archived one", p.Prereqs)
	}

	// v1's epic is done: fully archived, no Objective.
	const v1EpicPath = ".savepoint/releases/v1/epics/E01-example/E01-Detail.md"
	if _, ok := targetByPath(p, v1EpicPath); ok {
		t.Errorf("v1's done epic was planned as an Objective target, want archive-only")
	}
	if _, ok := archiveByPath(p, v1EpicPath); !ok {
		t.Errorf("v1's done epic detail was not archived")
	}

	// v1.1's epic is in_progress: gets an Objective.
	const v11EpicPath = ".savepoint/releases/v1.1/epics/E01-example/E01-Detail.md"
	if _, ok := targetByPath(p, v11EpicPath); !ok {
		t.Errorf("v1.1's in_progress epic was not planned as an Objective target")
	}

	if v1Archived.SourcePath == v11Active.Legacy.Path {
		t.Fatal("test setup: v1 and v1.1 T001-shared unexpectedly share a source path")
	}
}

// TestPlan_v1History_defectsStayReleaseScoped proves the two D001-shared
// defects (release-scoped, same short filename) are planned independently by
// their disposition: v1's resolved defect is archived, v1.1's open defect
// becomes an Issue.
func TestPlan_v1History_defectsStayReleaseScoped(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-history"))

	const v1Defect = ".savepoint/releases/v1/defects/D001-shared.md"
	const v11Defect = ".savepoint/releases/v1.1/defects/D001-shared.md"

	if _, ok := archiveByPath(p, v1Defect); !ok {
		t.Errorf("v1's resolved defect was not archived")
	}
	target, ok := targetByPath(p, v11Defect)
	if !ok || target.Kind != TargetIssue {
		t.Fatalf("v1.1's open defect was not planned as an Issue target: %+v", target)
	}
}

// TestPlan_v1History_findingDispositions proves each finding disposition in
// the fixture is planned per the epic's mapping: fixed becomes an Issue,
// waived is archived, and a duplicate whose canonical converts also becomes
// an Issue.
func TestPlan_v1History_findingDispositions(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-history"))

	const f001 = ".savepoint/audit/findings/F001-awaiting-proof.md"
	const f002 = ".savepoint/audit/findings/F002-owner-waiver.md"
	const f003 = ".savepoint/audit/findings/F003-duplicate.md"

	f001Target, ok := targetByPath(p, f001)
	if !ok || f001Target.Kind != TargetIssue {
		t.Fatalf("F001 (fixed) was not planned as an Issue: %+v", f001Target)
	}

	if _, ok := archiveByPath(p, f002); !ok {
		t.Errorf("F002 (waived) was not archived")
	}
	if _, ok := targetByPath(p, f002); ok {
		t.Errorf("F002 (waived) was planned as a target, want archive-only")
	}

	f003Target, ok := targetByPath(p, f003)
	if !ok || f003Target.Kind != TargetIssue {
		t.Fatalf("F003 (duplicate of a converting canonical) was not planned as an Issue: %+v", f003Target)
	}
	if f003Target.GlobalID == f001Target.GlobalID {
		t.Errorf("F003 and F001 share global ID %q, want distinct Issue identities", f003Target.GlobalID)
	}
}

// TestPlan_v1ShippedAuditReadmesAreArchivedNotParsed proves the READMEs V1
// init shipped into audit/findings and audit/runs neither abort planning nor
// become Issues: both are archived byte-for-byte, and real findings still
// convert.
func TestPlan_v1ShippedAuditReadmesAreArchivedNotParsed(t *testing.T) {
	root := copyFixtureProject(t, "v1-history")
	const findingsReadme = ".savepoint/audit/findings/README.md"
	const runsReadme = ".savepoint/audit/runs/README.md"
	writeFile(t, filepath.Join(root, filepath.FromSlash(findingsReadme)), "# Audit Findings\n\nOne file per finding.\n")
	writeFile(t, filepath.Join(root, filepath.FromSlash(runsReadme)), "# Audit Runs\n\nOne file per run.\n")

	p := mustPlan(t, root)
	for _, path := range []string{findingsReadme, runsReadme} {
		if _, ok := archiveByPath(p, path); !ok {
			t.Errorf("%s was not archived", path)
		}
		if target, ok := targetByPath(p, path); ok {
			t.Errorf("%s was planned as a target: %+v", path, target)
		}
	}
	if target, ok := targetByPath(p, ".savepoint/audit/findings/F001-awaiting-proof.md"); !ok || target.Kind != TargetIssue {
		t.Errorf("F001 no longer converts to an Issue alongside the README: %+v", target)
	}
}

// TestPlan_v1History_epicAuditAndRegisterArchived proves the epic audit file
// and the project-level audit prompt/register/run are archived as immutable
// history, independent of any finding's disposition.
func TestPlan_v1History_epicAuditAndRegisterArchived(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-history"))

	for _, path := range []string{
		".savepoint/releases/v1.1/epics/E01-example/E01-Audit.md",
		".savepoint/audit/prompt.md",
		".savepoint/audit/register.md",
		".savepoint/audit/runs/2026-07-01-example.md",
	} {
		if _, ok := archiveByPath(p, path); !ok {
			t.Errorf("%s was not archived", path)
		}
	}
}

// --- schema gate ----------------------------------------------------------

func TestPlan_schemaAlreadyV2_reportsAndDoesNoWork(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "schema_version: 2\n")

	p, err := Plan(root, nil, fixedClock(time.Now()), fixedOperationID("op"))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if !p.SchemaAlreadyV2 {
		t.Fatal("SchemaAlreadyV2 = false, want true")
	}
	if len(p.Targets) != 0 || len(p.Archives) != 0 || len(p.Conflicts) != 0 {
		t.Errorf("Plan() = %+v, want no work planned for an already-V2 project", p)
	}
}

func TestPlan_unsupportedSchemaVersion_returnsExistingDiagnostic(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "schema_version: 99\n")

	_, err := Plan(root, nil, fixedClock(time.Now()), fixedOperationID("op"))
	if !errors.Is(err, data.ErrUnsupportedSchemaVersion) {
		t.Fatalf("Plan() error = %v, want data.ErrUnsupportedSchemaVersion", err)
	}
}

// --- existing manifest conflict -------------------------------------------

func TestPlan_existingManifest_isRefusedAsConflictNotOverwritten(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	manifestPath := filepath.Join(root, filepath.FromSlash(migrationsDirRel), manifestFileName)
	writeFile(t, manifestPath, "manifest_schema_version: 1\n")

	before := snapshotTree(t, root)

	p, err := Plan(root, nil, fixedClock(time.Now()), fixedOperationID("op"))
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	if len(p.Conflicts) != 1 || p.Conflicts[0].Kind != ConflictExistingManifest {
		t.Fatalf("Conflicts = %+v, want exactly one ConflictExistingManifest", p.Conflicts)
	}
	if len(p.Targets) != 0 {
		t.Errorf("Targets = %v, want none planned once an existing manifest conflict is found", p.Targets)
	}

	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

// --- destination collision conflict -----------------------------------------
//
// Regression coverage for an audit finding against E45: a path this plan
// intends to create (a converted record, Idea.md, or an archive entry) that
// already exists on disk used to reach Apply undetected and crash with an
// internal invariant error and no way out. Plan must catch it first, as a
// named, reviewable ConflictDestinationExists, exactly like an existing
// manifest.

func TestPlan_archivePathLivesInsideSavepoint(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	writeFile(t, filepath.Join(root, ".savepoint", "PRD.md"), "# Idea\n")

	p := mustPlan(t, root)

	if len(p.Archives) == 0 {
		t.Fatal("no archives planned, want at least the PRD.md archive entry")
	}
	for _, a := range p.Archives {
		if !strings.HasPrefix(a.ArchivePath, ".savepoint/archive/v1/") {
			t.Errorf("archive %s ArchivePath = %q, want it to live under .savepoint/archive/v1/ so it is inventoried, path-confined, and collision-checked like every other write",
				a.SourcePath, a.ArchivePath)
		}
	}
}

func TestPlan_preExistingIdeaDestination_isNamedConflict(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	writeFile(t, filepath.Join(root, ".savepoint", "PRD.md"), "# Idea\n")
	// A user who read the V2 docs and started early: Idea.md already exists
	// exactly where PRD.md's DocumentIdea relocation plans to create it.
	writeFile(t, filepath.Join(root, ".savepoint", "Idea.md"), "my early idea\n")

	before := snapshotTree(t, root)

	p := mustPlan(t, root)

	want := ".savepoint/Idea.md"
	found := false
	for _, c := range p.Conflicts {
		if c.Kind == ConflictDestinationExists && c.Path == want {
			found = true
		}
	}
	if !found {
		t.Fatalf("Conflicts = %+v, want a ConflictDestinationExists naming %s", p.Conflicts, want)
	}

	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestPlan_preExistingArchiveDestination_isNamedConflict(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	writeFile(t, filepath.Join(root, ".savepoint", "PRD.md"), "# Idea\n")
	// Something already occupies PRD.md's own archive destination.
	collision := filepath.Join(root, ".savepoint", "archive", "v1", ".savepoint", "PRD.md")
	writeFile(t, collision, "not what migration expects here\n")

	before := snapshotTree(t, root)

	p := mustPlan(t, root)

	want := ".savepoint/archive/v1/.savepoint/PRD.md"
	found := false
	for _, c := range p.Conflicts {
		if c.Kind == ConflictDestinationExists && c.Path == want {
			found = true
		}
	}
	if !found {
		t.Fatalf("Conflicts = %+v, want a ConflictDestinationExists naming %s", p.Conflicts, want)
	}

	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestPlan_replacedRouterDocument_isNeverACollision(t *testing.T) {
	// router.md is rewritten in place (ActionReplace): its destination is
	// expected to already exist, and must never be reported as a collision.
	root := t.TempDir()
	writeMinimalV1Project(t, root)

	p := mustPlan(t, root)

	for _, c := range p.Conflicts {
		if c.Kind == ConflictDestinationExists && c.Path == ".savepoint/router.md" {
			t.Fatalf("router.md's own replace target was reported as a destination collision: %+v", c)
		}
	}
}

// --- preserved migrations-dir coexistence ----------------------------------

// TestPlan_preservesExistingMigrationsReadmeAndArchivedSkill proves
// internal/init's pre-existing .savepoint/migrations/README.md and archived
// skill copy are neither archived under archive/v1/ nor reported as
// unclassified: they are left exactly where they are.
func TestPlan_preservesExistingMigrationsReadmeAndArchivedSkill(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	writeFile(t, filepath.Join(root, ".savepoint", "migrations", "README.md"), "# Savepoint Migrations\n")
	writeFile(t, filepath.Join(root, ".savepoint", "migrations", "savepoint-audit-SKILL.md"), "# archived skill\n")

	p := mustPlan(t, root)

	for _, path := range []string{".savepoint/migrations/README.md", ".savepoint/migrations/savepoint-audit-SKILL.md"} {
		if _, ok := archiveByPath(p, path); ok {
			t.Errorf("%s was archived, want it preserved in place untouched", path)
		}
		if _, ok := targetByPath(p, path); ok {
			t.Errorf("%s was planned as a target, want it preserved in place untouched", path)
		}
	}
}

// --- unrecognized lifecycle ambiguity --------------------------------------

func TestPlan_unrecognizedTaskStatus_isBlockingAmbiguityNotHealed(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "v1-PRD.md"), "---\nstatus: in_progress\n---\n\n# V1\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-x", "E01-Detail.md"),
		"---\nstatus: in_progress\n---\n\n# E01\n")
	taskPath := ".savepoint/releases/v1/epics/E01-x/tasks/T001-weird.md"
	writeFile(t, filepath.Join(root, filepath.FromSlash(taskPath)),
		"---\nid: E01-x/T001-weird\nstatus: escalated\ndepends_on: []\n---\n\n# T001\n")

	p := mustPlan(t, root)

	if _, ok := targetByPath(p, taskPath); ok {
		t.Errorf("task with unrecognized status was planned as a target, want it blocked as an ambiguity")
	}
	if _, ok := archiveByPath(p, taskPath); ok {
		t.Errorf("task with unrecognized status was archived, want it blocked as an ambiguity instead of guessed")
	}

	var found *Ambiguity
	for i := range p.Ambiguities {
		if p.Ambiguities[i].Path == taskPath {
			found = &p.Ambiguities[i]
		}
	}
	if found == nil {
		t.Fatalf("Ambiguities = %+v, want one naming %s", p.Ambiguities, taskPath)
	}
	if found.Kind != AmbiguityUnrecognizedLifecycle {
		t.Errorf("Ambiguity.Kind = %v, want AmbiguityUnrecognizedLifecycle", found.Kind)
	}
	if found.ID != string(AmbiguityUnrecognizedLifecycle)+":"+taskPath {
		t.Errorf("Ambiguity.ID = %q, want a stable id derived from kind and path", found.ID)
	}
}

// TestPlan_missingDependencyTarget_isBlockingAmbiguity proves an active
// task's depends_on reference that resolves to nothing produces a named,
// stable ambiguity rather than silently dropping the reference.
func TestPlan_missingDependencyTarget_isBlockingAmbiguity(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), "# Router\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "v1-PRD.md"), "---\nstatus: in_progress\n---\n\n# V1\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-x", "E01-Detail.md"),
		"---\nstatus: in_progress\n---\n\n# E01\n")
	taskPath := ".savepoint/releases/v1/epics/E01-x/tasks/T001-orphan.md"
	writeFile(t, filepath.Join(root, filepath.FromSlash(taskPath)),
		"---\nid: E01-x/T001-orphan\nstatus: planned\ndepends_on: [T099-missing]\n---\n\n# T001\n")

	p := mustPlan(t, root)

	target, ok := targetByPath(p, taskPath)
	if !ok {
		t.Fatalf("task was not planned as an active target")
	}
	if len(target.DependsOn) != 0 {
		t.Errorf("DependsOn = %v, want empty: the reference never resolved", target.DependsOn)
	}

	var found bool
	for _, a := range p.Ambiguities {
		if a.Kind == AmbiguityMissingDependencyTarget && a.Path == taskPath {
			found = true
		}
	}
	if !found {
		t.Errorf("Ambiguities = %+v, want a AmbiguityMissingDependencyTarget for %s", p.Ambiguities, taskPath)
	}
}

// --- manifest coexistence with sorted output -------------------------------

func TestPlan_ambiguitiesSortedByID(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-basic"))
	ids := make([]string, len(p.Ambiguities))
	for i, a := range p.Ambiguities {
		ids[i] = a.ID
	}
	if !sort.StringsAreSorted(ids) {
		t.Errorf("Ambiguities not sorted by ID: %v", ids)
	}
}

func TestIDAllocator_allocateUsesHyphenatedIdentities(t *testing.T) {
	allocator := newIDAllocator()
	for _, prefix := range []string{"R", "O", "T", "I"} {
		if got, want := allocator.allocate(prefix), prefix+"-001"; got != want {
			t.Errorf("allocate(%q) = %q, want %q", prefix, got, want)
		}
	}

	allocator.reserve("I-002")
	if got, want := allocator.allocate("I"), "I-003"; got != want {
		t.Errorf("allocate(I) after reserving I-002 = %q, want %q", got, want)
	}
}
