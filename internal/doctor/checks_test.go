package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/testutil"
)

// --- CheckConfig ---

func TestCheckConfigMissing(t *testing.T) {
	root := t.TempDir()
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("CheckConfig() = %v, want not found error", err)
	}
}

func TestCheckConfigInvalidYAML(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "theme: [broken")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "invalid YAML") {
		t.Fatalf("CheckConfig() = %v, want invalid YAML error", err)
	}
}

func TestCheckConfigMissingQualityGates(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "theme:\n  bg: \"#000\"\n")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "quality_gates") {
		t.Fatalf("CheckConfig() = %v, want quality_gates error", err)
	}
}

func TestCheckConfigMissingTheme(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "quality_gates:\n  block_on_failure: true\n")
	err := CheckConfig(root)
	if err == nil || !strings.Contains(err.Error(), "theme") {
		t.Fatalf("CheckConfig() = %v, want theme error", err)
	}
}

func TestCheckConfigValid(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "quality_gates:\n  block_on_failure: true\ntheme:\n  bg: \"#000\"\n")
	if err := CheckConfig(root); err != nil {
		t.Fatalf("CheckConfig() = %v, want nil", err)
	}
}

// --- CheckRouter ---

// --- CheckStructure ---

// --- CheckDependencies ---

// --- CheckAuditState ---

// --- CheckOrphans ---

// helpers

// --- CheckStructure complexity ---

// --- CheckDefects ---

// --- CheckAuditRegister ---

// --- RunV2Checks ---

func writeV2Objective(t *testing.T, root, dirName, id, title string) string {
	t.Helper()
	ensureV2TestGoal(t, root)
	return writeV2ObjectiveWithRelease(t, root, dirName, id, title, "R-001")
}

func writeV2ObjectiveWithoutGoal(t *testing.T, root, dirName, id, title string) string {
	t.Helper()
	path := filepath.Join(root, "objectives", dirName, "Objective.md")
	testutil.WriteFile(t, path, "---\nid: "+id+"\ntitle: \""+title+"\"\nstatus: planned\n---\n\n# "+title+"\n")
	return path
}

// writeV2ObjectiveInProgress writes an Objective whose Tasks have started, so
// fixtures about other diagnostics are not also flagged as a planned
// Objective with started work.
func writeV2ObjectiveInProgress(t *testing.T, root, dirName, id, title string) string {
	t.Helper()
	ensureV2TestGoal(t, root)
	path := filepath.Join(root, "objectives", dirName, "Objective.md")
	testutil.WriteFile(t, path, "---\nid: "+id+"\ntitle: \""+title+"\"\nstatus: in_progress\nrelease: R-001\n---\n\n# "+title+"\n")
	return path
}

func writeV2Task(t *testing.T, root, objDirName, fileName, id, title, objective string) string {
	t.Helper()
	path := filepath.Join(root, "objectives", objDirName, "tasks", fileName)
	testutil.WriteFile(t, path, "---\nid: "+id+"\ntitle: \""+title+"\"\nobjective: "+objective+"\nplanned_by: {role: planner, session: planning-001}\nstatus: planned\n---\n\n# "+title+"\n")
	return path
}

func writeV2Release(t *testing.T, root, dirName, id, status, extraFrontmatter string) string {
	t.Helper()
	path := filepath.Join(root, "releases", dirName, "Release.md")
	content := "---\nid: " + id + "\ntitle: \"Release " + id + "\"\nstatus: " + status + "\n" + extraFrontmatter + "---\n\n" +
		"## Outcome\n\nDeliver the recorded promise.\n\n" +
		"## Why\n\nThe delivery boundary needs a stable identity.\n\n" +
		"## Success Conditions\n\n- All member Objectives are complete.\n\n" +
		"## Boundaries\n\nMembership is derived from Objective records.\n"
	testutil.WriteFile(t, path, content)
	return path
}

func ensureV2TestGoal(t *testing.T, root string) {
	t.Helper()
	goalPath := filepath.Join(root, "releases", "R-001", "Release.md")
	if _, err := os.Stat(goalPath); os.IsNotExist(err) {
		writeV2Release(t, root, "R-001", "R-001", "planned", "")
	} else if err != nil {
		t.Fatalf("Stat(%s) error = %v", goalPath, err)
	}
}

func writeV2ObjectiveWithRelease(t *testing.T, root, dirName, id, title, release string) string {
	t.Helper()
	path := filepath.Join(root, "objectives", dirName, "Objective.md")
	testutil.WriteFile(t, path, "---\nid: "+id+"\ntitle: \""+title+"\"\nstatus: planned\nrelease: "+release+"\n---\n\n# "+title+"\n")
	return path
}

func TestCheckProject_v2ValidNoProblems(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")

	if problems := RunV2Checks(root).Project; len(problems) != 0 {
		t.Fatalf("RunV2Checks().Project = %v, want no problems for a valid V2 project", problems)
	}
}

func TestCheckProject_missingReleaseNamesFileAndIDs(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Release(t, root, "R-001-first", "R-001", "planned", "")
	path := filepath.Join(root, "objectives", "O-001-ship", "Objective.md")
	testutil.WriteFile(t, path, "---\nid: O-001\ntitle: \"Ship it\"\nstatus: planned\nrelease: R-999\n---\n\n# Ship it\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 {
		t.Fatalf("RunV2Checks().Project = %v, want one missing Release problem", problems)
	}
	for _, want := range []string{"[v2-missing-release]", filepath.FromSlash("objectives/O-001-ship/Objective.md"), "O-001", "R-999"} {
		if !strings.Contains(problems[0].Message, want) {
			t.Errorf("problem message = %q, want %q", problems[0].Message, want)
		}
	}
	if !strings.Contains(problems[0].Repair, "referenced R-### Release") {
		t.Errorf("problem repair = %q, want manual Release-reference guidance", problems[0].Repair)
	}
}

func TestCheckReleaseReadiness_ignoresEmptyActiveGoalAndKeepsCanonicalFindings(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Release(t, root, "R-001-empty", "R-001", "in_progress", "")
	writeV2Release(t, root, "R-002-active", "R-002", "in_progress", "")
	writeV2Release(t, root, "R-003-done-empty", "R-003", "done", "")
	writeV2ObjectiveWithRelease(t, root, "O-002-member", "O-002", "Member", "R-002")
	writeV2ObjectiveWithoutGoal(t, root, "O-003-unassigned", "O-003", "Independent")

	problems := RunV2Checks(root).Releases
	if len(problems) != 2 {
		t.Fatalf("RunV2Checks().Releases = %v, want only actionable readiness findings", problems)
	}
	wants := []string{"[v2-release-objective-incomplete] release R-002", "[v2-release-no-objectives] release R-003"}
	paths := []string{
		filepath.Join("releases", "R-002-active", "Release.md"),
		filepath.Join("releases", "R-003-done-empty", "Release.md"),
	}
	for i, want := range wants {
		if !strings.Contains(problems[i].Message, want) {
			t.Errorf("problems[%d].Message = %q, want %q", i, problems[i].Message, want)
		}
		if problems[i].File != paths[i] {
			t.Errorf("problems[%d].File = %q, want the Release source path", i, problems[i].File)
		}
		if problems[i].Repair == "" || !strings.Contains(problems[i].Repair, "doctor") {
			t.Errorf("problems[%d].Repair = %q, want non-destructive guidance", i, problems[i].Repair)
		}
	}
}

func TestCheckReleaseReadiness_reportsMissingUnknownStaleNeedsWorkAndAcceptance(t *testing.T) {
	tests := []struct {
		name       string
		releaseID  string
		checkBlock string
		releaseFM  string
		want       string
	}{
		{name: "missing", releaseID: "R-001", want: "[v2-release-clearance-missing]"},
		{
			name:       "needs work",
			releaseID:  "R-002",
			checkBlock: "needs-work",
			want:       "[v2-release-clearance-needs-work]",
		},
		{
			name:       "unknown",
			releaseID:  "R-003",
			checkBlock: "unknown",
			releaseFM:  "freshness: {state: unknown, check: C-030, assessed_by: {role: checker, session: freshness-3}, assessed_at: '2026-09-14T00:00:00Z', basis: not reassessed}\n",
			want:       "[v2-release-clearance-unknown]",
		},
		{
			name:       "stale",
			releaseID:  "R-004",
			checkBlock: "stale",
			releaseFM:  "freshness: {state: stale, check: C-041, assessed_by: {role: checker, session: freshness-4}, assessed_at: '2026-09-14T00:00:00Z', basis: code changed}\n",
			want:       "[v2-release-clearance-stale]",
		},
		{
			name:       "owner acceptance",
			releaseID:  "R-005",
			checkBlock: "owner",
			releaseFM:  "freshness: {state: current, check: C-050, assessed_by: {role: checker, session: freshness-5}, assessed_at: '2026-09-14T00:00:00Z', basis: checked}\n",
			want:       "[v2-release-owner-acceptance-missing]",
		},
		{
			name:       "stale owner acceptance",
			releaseID:  "R-006",
			checkBlock: "owner-stale",
			releaseFM:  "freshness: {state: current, check: C-061, assessed_by: {role: checker, session: freshness-6}, assessed_at: '2026-09-14T00:00:00Z', basis: checked}\nowner_validation: {required: true, accepted_check: C-060, accepted_by: {role: owner, session: owner-6}}\n",
			want:       "[v2-release-owner-acceptance-stale]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
			writeV2Release(t, root, tt.releaseID+"-active", tt.releaseID, "in_progress", tt.releaseFM)
			objectiveID := "O" + strings.TrimPrefix(tt.releaseID, "R")
			objectiveDir := objectiveID + "-member"
			testutil.WriteFile(t, filepath.Join(root, "objectives", objectiveDir, "Objective.md"),
				"---\nid: "+objectiveID+"\ntitle: \"Member\"\nstatus: done\nrelease: "+tt.releaseID+"\nfreshness: {state: current, check: C"+strings.TrimPrefix(tt.releaseID, "R")+"0, assessed_by: {role: checker, session: objective-checker}, assessed_at: '2026-09-14T00:00:00Z', basis: checked}\n---\n\n# Member\n")
			objectiveCheck := "C" + strings.TrimPrefix(tt.releaseID, "R") + "0"
			writeV2Check(t, root, objectiveCheck, "{kind: objective, id: "+objectiveID+"}", "CLEAR", "")

			switch tt.checkBlock {
			case "needs-work":
				writeV2Check(t, root, "C-020", "{kind: release, id: R-002}", "NEEDS WORK", "")
			case "unknown":
				writeV2Check(t, root, "C-030", "{kind: release, id: R-003}", "CLEAR", "")
			case "stale":
				writeV2Check(t, root, "C-040", "{kind: release, id: R-004}", "CLEAR", "")
				writeV2Check(t, root, "C-041", "{kind: release, id: R-004}", "CLEAR", "C-040")
			case "owner":
				writeV2Check(t, root, "C-050", "{kind: release, id: R-005}", "CLEAR", "")
			case "owner-stale":
				writeV2Check(t, root, "C-060", "{kind: release, id: R-006}", "CLEAR", "")
				writeV2Check(t, root, "C-061", "{kind: release, id: R-006}", "CLEAR", "C-060")
			}

			problems := RunV2Checks(root).Releases
			if len(problems) != 1 || !strings.Contains(problems[0].Message, tt.want) {
				t.Fatalf("RunV2Checks().Releases = %v, want %s", problems, tt.want)
			}
		})
	}
}

func TestCheckReleaseReadiness_reportsMaterialUnresolvedIssueFromCanonicalGate(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Release(t, root, "R-007-active", "R-007", "in_progress",
		"freshness: {state: current, check: C-070, assessed_by: {role: checker, session: freshness-7}, assessed_at: '2026-09-14T00:00:00Z', basis: checked}\nowner_validation: {required: true, accepted_check: C-070, accepted_by: {role: owner, session: owner-7}}\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-007-member", "Objective.md"),
		"---\nid: O-007\ntitle: \"Member\"\nstatus: done\nrelease: R-007\nfreshness: {state: current, check: C-0070, assessed_by: {role: checker, session: objective-checker}, assessed_at: '2026-09-14T00:00:00Z', basis: checked}\n---\n\n# Member\n")
	writeV2Check(t, root, "C-0070", "{kind: objective, id: O-007}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "checks", "C-070.md"),
		"---\nid: C-070\nscope: {kind: release, id: R-007}\nresult: CLEAR\nchecked_by: {role: checker, session: release-checker}\nexecuted_session: build-7\nchecked_at: '2026-09-14T00:00:00Z'\nissues: [I-001]\n---\n\n# Check\n")
	writeV2Issue(t, root, "I-001-blocker.md", "I-001", "open", "defect", "checks: [C-070]\n")

	problems := RunV2Checks(root).Releases
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-release-issue-unresolved]") || !strings.Contains(problems[0].Message, "I-001") {
		t.Fatalf("RunV2Checks().Releases = %v, want one named material Issue blocker", problems)
	}
}

func TestCheckProject_SchemaVersionMalformed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: nope\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 {
		t.Fatalf("RunV2Checks().Project = %v, want exactly 1 problem", problems)
	}
	if !strings.Contains(problems[0].Message, "[schema-version-malformed]") {
		t.Errorf("RunV2Checks().Project[0].Message = %q, want the schema-version-malformed diagnostic name", problems[0].Message)
	}
	if want := filepath.Join(root, "config.yml"); problems[0].File != want {
		t.Errorf("RunV2Checks().Project[0].File = %q, want the config file %q", problems[0].File, want)
	}
	if problems[0].Repair == "" {
		t.Error("RunV2Checks().Project[0].Repair is empty, want a typed repair suggestion")
	}
}

// TestCheckProject_MissingTaskTitleNotBackfilledFromObjective proves doctor
// reports a missing V2 Task title as an actionable schema error rather than
// silently accepting the Task's objective owner reference as a display
// title substitute.
func TestCheckProject_MissingTaskTitleNotBackfilledFromObjective(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-write.md"),
		"---\nid: T-001\nobjective: O-001\nstatus: planned\n---\n\n# Write it\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 {
		t.Fatalf("RunV2Checks().Project = %v, want exactly 1 problem", problems)
	}
	if !strings.Contains(problems[0].Message, "[v2-missing-field]") {
		t.Errorf("RunV2Checks().Project[0].Message = %q, want the v2-missing-field diagnostic name", problems[0].Message)
	}
	if !strings.Contains(problems[0].Message, "missing required field title") {
		t.Errorf("RunV2Checks().Project[0].Message = %q, want it to name the missing title field, not accept objective O-001 as a substitute", problems[0].Message)
	}
}

func TestCheckProject_MissingOwner(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-999")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-missing-owner]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-missing-owner", problems)
	}
}

func TestCheckProject_DependencyCycle(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-alpha.md"),
		"---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: planned\ndepends_on: [{task: T-002}]\n---\n\n# Alpha\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-002-beta.md"),
		"---\nid: T-002\ntitle: \"Beta\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: planned\ndepends_on: [{task: T-001}]\n---\n\n# Beta\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-dependency-cycle]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-dependency-cycle", problems)
	}
}

// TestCheckProject_ReadOnly proves RunV2Checks never writes to the project
// it inspects, for a valid V1 project with quality gates configured, a valid
// V2 project, and an invalid V2 project.
func TestCheckProject_ReadOnly(t *testing.T) {
	t.Run("valid V1 project with quality gates configured", func(t *testing.T) {
		root := t.TempDir()
		testutil.SetupMinimalProject(t, root, "v1", "E01-foo")
		configPath := filepath.Join(root, "config.yml")
		before, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}

		RunV2Checks(root)

		after, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) != string(after) {
			t.Fatalf("RunV2Checks() changed config.yml:\nbefore: %s\nafter: %s", before, after)
		}
	})

	t.Run("valid V2 project", func(t *testing.T) {
		root := t.TempDir()
		testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
		writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
		taskPath := writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
		before, err := os.ReadFile(taskPath)
		if err != nil {
			t.Fatal(err)
		}

		if problems := RunV2Checks(root).Project; len(problems) != 0 {
			t.Fatalf("RunV2Checks().Project = %v, want no problems", problems)
		}

		after, err := os.ReadFile(taskPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) != string(after) {
			t.Fatalf("RunV2Checks().Project changed %s:\nbefore: %s\nafter: %s", taskPath, before, after)
		}
	})

	t.Run("invalid V2 project", func(t *testing.T) {
		root := t.TempDir()
		testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
		writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
		taskPath := writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-999")
		before, err := os.ReadFile(taskPath)
		if err != nil {
			t.Fatal(err)
		}

		if problems := RunV2Checks(root).Project; len(problems) == 0 {
			t.Fatal("RunV2Checks().Project = no problems, want the missing-owner diagnostic")
		}

		after, err := os.ReadFile(taskPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(before) != string(after) {
			t.Fatalf("RunV2Checks().Project changed %s:\nbefore: %s\nafter: %s", taskPath, before, after)
		}
	})
}

// --- RunV2Checks: Check and evidence diagnostics ---

func writeV2Check(t *testing.T, root, id, scope, result, supersedes string) string {
	t.Helper()
	path := filepath.Join(root, "checks", id+".md")
	body := "---\nid: " + id + "\nscope: " + scope + "\nresult: " + result + "\n" +
		"checked_by: {role: checker, session: sess-1}\nexecuted_session: build-1\nchecked_at: '2026-09-14T00:00:00Z'\n"
	if supersedes != "" {
		body += "supersedes: " + supersedes + "\n"
	}
	body += "---\n\n# Check\n"
	testutil.WriteFile(t, path, body)
	return path
}

func TestCheckProject_CheckMalformed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "BOGUS", "")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-check-malformed]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-check-malformed", problems)
	}
	if problems[0].Repair == "" {
		t.Error("RunV2Checks().Project[0].Repair is empty, want a typed repair suggestion")
	}
}

func TestCheckProject_CheckMissingScopeTarget(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	writeV2Check(t, root, "C-001", "{kind: task, id: T-999}", "CLEAR", "")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-check-missing-scope-target]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-check-missing-scope-target", problems)
	}
}

func TestCheckProject_CheckMissingReference(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	writeV2Check(t, root, "C-002", "{kind: task, id: T-001}", "CLEAR", "C-001")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-check-missing-reference]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-check-missing-reference", problems)
	}
}

func TestCheckProject_CheckSupersedesConflict(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "CLEAR", "")
	writeV2Check(t, root, "C-002", "{kind: task, id: T-001}", "CLEAR", "C-001")
	writeV2Check(t, root, "C-003", "{kind: task, id: T-001}", "CLEAR", "C-001")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-check-supersedes-conflict]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-check-supersedes-conflict", problems)
	}
}

func TestCheckProject_EvidenceMalformed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-write.md"),
		"---\nid: T-001\ntitle: \"Write it\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: planned\n"+
			"freshness: {state: bogus, check: C-001, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: reviewed}\n"+
			"---\n\n# Write it\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-evidence-malformed]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-evidence-malformed", problems)
	}
}

func TestCheckProject_EvidenceMissingReference(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-write.md"),
		"---\nid: T-001\ntitle: \"Write it\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: planned\nlast_check: C-999\n---\n\n# Write it\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-evidence-missing-reference]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-evidence-missing-reference", problems)
	}
}

func writeV2Issue(t *testing.T, root, fileName, id, status, issueType, extra string) string {
	t.Helper()
	path := filepath.Join(root, "issues", fileName)
	body := "---\nid: " + id + "\ntitle: \"Follow-up\"\ntype: " + issueType + "\nstatus: " + status +
		"\nsource: {kind: report, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z'}\n" + extra +
		"---\n\n# Issue\n"
	testutil.WriteFile(t, path, body)
	return path
}

func TestCheckProject_IssueMalformed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-bad.md", "I-001", "open", "bogus", "")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-malformed]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-malformed", problems)
	}
}

func TestCheckProject_IssueMissingDuplicateTarget(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "open", "defect", "duplicate_of: I-999\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-missing-duplicate-target]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-missing-duplicate-target", problems)
	}
}

func TestCheckProject_IssueSelfDuplicate(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "open", "defect", "duplicate_of: I-001\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-self-duplicate]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-self-duplicate", problems)
	}
}

func TestCheckProject_IssueDuplicateCycle(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "open", "defect", "duplicate_of: I-002\n")
	writeV2Issue(t, root, "I-002-y.md", "I-002", "open", "defect", "duplicate_of: I-001\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-duplicate-cycle]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-duplicate-cycle", problems)
	}
}

func TestCheckProject_IssueMissingLinkTarget(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "open", "defect", "tasks: [T-999]\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-missing-link-target]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-missing-link-target", problems)
	}
}

func TestCheckProject_IssueUnpairedCheckLink(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	testutil.WriteFile(t, filepath.Join(root, "checks", "C-001.md"),
		"---\nid: C-001\nscope: {kind: task, id: T-001}\nresult: CLEAR\n"+
			"checked_by: {role: checker, session: sess-1}\nexecuted_session: build-1\nchecked_at: '2026-09-14T00:00:00Z'\nissues: [I-001]\n---\n\n# Check\n")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "open", "defect", "")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-unpaired-check-link]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-unpaired-check-link", problems)
	}
}

func TestCheckProject_IssueResolutionRequired(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "resolved", "defect", "")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-required]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-resolution-required", problems)
	}
}

func TestCheckProject_IssueResolutionNotAllowed(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "open", "defect",
		"resolution: {disposition: accepted, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z', reason: \"risk accepted\"}\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-not-allowed]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-resolution-not-allowed", problems)
	}
}

func TestCheckProject_IssueResolutionMissingProof(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "resolved", "defect",
		"resolution: {disposition: verified, actor: {role: checker, session: sess-1}, at: '2026-09-15T00:00:00Z'}\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-missing-proof]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-resolution-missing-proof", problems)
	}
}

func TestCheckProject_IssueResolutionUnusableProof(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "CLEAR", "")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "resolved", "defect",
		"resolution: {disposition: verified, check: C-001, actor: {role: checker, session: sess-1}, at: '2026-09-15T00:00:00Z'}\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-unusable-proof]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-resolution-unusable-proof", problems)
	}
}

func TestCheckProject_IssueResolutionFieldMismatch(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "CLEAR", "")
	writeV2Issue(t, root, "I-001-x.md", "I-001", "resolved", "defect",
		"checks: [C-001]\nresolution: {disposition: accepted, check: C-001, actor: {role: owner, session: owner-1}, at: '2026-09-15T00:00:00Z', reason: \"n/a\"}\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-resolution-field-mismatch]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming v2-issue-resolution-field-mismatch", problems)
	}
}

// TestV2DiagnosticName_noNewIssueSentinelFallsThrough proves every new E44
// data sentinel — including the write-time-only ones RunV2Checks can never
// actually surface, since doctor never writes — maps to a name other than the
// generic fallback, so v2DiagnosticName stays a complete registry of every
// named diagnostic.
func TestV2DiagnosticName_noNewIssueSentinelFallsThrough(t *testing.T) {
	sentinels := []error{
		data.ErrV2InvalidReleaseReference,
		data.ErrV2MissingRelease,
		data.ErrV2ReleaseMissingSection,
		data.ErrV2ReleaseLegacyMalformed,
		data.ErrV2IssueMalformed,
		data.ErrV2IssueMissingDuplicateTarget,
		data.ErrV2IssueSelfDuplicate,
		data.ErrV2IssueDuplicateCycle,
		data.ErrV2IssueMissingLinkTarget,
		data.ErrV2IssueUnpairedCheckLink,
		data.ErrV2IssueResolutionRequired,
		data.ErrV2IssueResolutionNotAllowed,
		data.ErrV2IssueResolutionMissingProof,
		data.ErrV2IssueResolutionUnusableProof,
		data.ErrV2IssueResolutionFieldMismatch,
		data.ErrV2IssueAlreadyExists,
		data.ErrV2IssueHistoryNotAppendOnly,
	}
	for _, sentinel := range sentinels {
		name := v2DiagnosticName(fmt.Errorf("wrap: %w", sentinel))
		if name == "v2-project-error" {
			t.Errorf("v2DiagnosticName(%v) fell through to the generic fallback", sentinel)
		}
		if V2ProblemRepair(name) == "Review the V2 project diagnostic and fix the reported record" {
			t.Errorf("V2ProblemRepair(%q) fell through to the generic repair suggestion", name)
		}
	}
}

// TestCheckProject_ObjectiveConsistencyDiagnostics proves doctor reports every
// evaluation-level inconsistency InspectObjectiveConsistency finds — a done
// Objective without current integration clearance, and a done Objective with
// an incomplete owned Task — named against the Objective's own source record.
func TestCheckProject_ObjectiveConsistencyDiagnostics(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	ensureV2TestGoal(t, root)
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "Objective.md"),
		"---\nid: O-001\ntitle: \"Ship it\"\nstatus: done\nrelease: R-001\n---\n\n# Ship it\n")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-write.md"),
		"---\nid: T-001\ntitle: \"Write it\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: in_progress\nstage: build\n---\n\n# Write it\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 2 {
		t.Fatalf("RunV2Checks().Project = %+v, want 2 objective consistency problems", problems)
	}
	if !strings.Contains(problems[0].Message, "[v2-objective-done-without-clearance]") {
		t.Errorf("problems[0].Message = %q, want v2-objective-done-without-clearance", problems[0].Message)
	}
	if !strings.Contains(problems[1].Message, "[v2-objective-done-with-incomplete-task]") {
		t.Errorf("problems[1].Message = %q, want v2-objective-done-with-incomplete-task", problems[1].Message)
	}
	for _, p := range problems {
		if p.Repair == "" {
			t.Errorf("problem %+v missing repair", p)
		}
		if !strings.HasSuffix(p.File, "Objective.md") {
			t.Errorf("problem File = %q, want the Objective's source record path", p.File)
		}
	}
}

// TestCheckProject_IssueVerifiedProofSuperseded proves doctor reports a
// verified Issue whose proof Check has been superseded by a later Check
// recorded against the same scope — an inconsistency a successful load
// cannot refuse by itself.
func TestCheckProject_IssueVerifiedProofSuperseded(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "CLEAR", "")
	writeV2Check(t, root, "C-002", "{kind: task, id: T-001}", "CLEAR", "C-001")
	writeV2Issue(t, root, "I-001-flaky.md", "I-001", "resolved", "defect",
		"checks: [C-001]\nresolution: {disposition: verified, check: C-001, actor: {role: checker, session: sess-1}, at: '2026-09-14T00:00:00Z'}\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[v2-issue-verified-proof-superseded]") {
		t.Fatalf("RunV2Checks().Project = %+v, want 1 problem naming v2-issue-verified-proof-superseded", problems)
	}
	if !strings.HasSuffix(problems[0].File, "I-001-flaky.md") {
		t.Errorf("problems[0].File = %q, want the Issue's source record path", problems[0].File)
	}
}

// TestCheckProject_OpenAdvisoryIssueDoesNotFailLoad proves an open Issue with
// no resolution is purely advisory: structural and gate evidence decide
// health, so a project carrying only an open Issue alongside otherwise-clean
// records reports no problems.
func TestCheckProject_OpenAdvisoryIssueDoesNotFailLoad(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	writeV2Issue(t, root, "I-001-flaky.md", "I-001", "open", "defect", "")

	if problems := RunV2Checks(root).Project; len(problems) != 0 {
		t.Fatalf("RunV2Checks().Project = %v, want no problems: an open advisory Issue must not make the project unhealthy", problems)
	}
}

// TestIssuePostureReport_derivedCountsNoStoredSummary proves Issue posture is
// computed fresh from the loaded index every call, with no summary cached
// between them.
func TestIssuePostureReport_derivedCountsNoStoredSummary(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Issue(t, root, "I-001-a.md", "I-001", "open", "defect", "")

	posture := RunV2Checks(root).Issues
	if posture == nil {
		t.Fatal("RunV2Checks().Issues = nil, want a posture for a V2 project")
	}
	if posture.StatusCounts[data.IssueStatusOpen] != 1 {
		t.Errorf("StatusCounts[open] = %d, want 1", posture.StatusCounts[data.IssueStatusOpen])
	}
	if posture.TypeCounts[data.IssueTypeDefect] != 1 {
		t.Errorf("TypeCounts[defect] = %d, want 1", posture.TypeCounts[data.IssueTypeDefect])
	}

	writeV2Issue(t, root, "I-002-b.md", "I-002", "open", "drift", "")
	posture2 := RunV2Checks(root).Issues
	if posture2.StatusCounts[data.IssueStatusOpen] != 2 {
		t.Errorf("StatusCounts[open] after a second Issue = %d, want 2 (recomputed, not cached)", posture2.StatusCounts[data.IssueStatusOpen])
	}
}

func TestIssuePostureReport_v1ProjectReturnsNil(t *testing.T) {
	root := t.TempDir()
	testutil.SetupMinimalProject(t, root, "v1", "E01-foo")

	if got := RunV2Checks(root).Issues; got != nil {
		t.Errorf("RunV2Checks().Issues = %+v, want nil for a V1 project", got)
	}
}

// TestCheckProject_v2ValidWithChecksNoProblems proves a clean V2 project
// carrying real Check and evidence records — current clearance on a
// technical Task and an owner-accepted Task — reports no problems.
func TestCheckProject_v2ValidWithChecksNoProblems(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2ObjectiveInProgress(t, root, "O-001-ship", "O-001", "Ship it")
	writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-write.md"),
		"---\nid: T-001\ntitle: \"Write it\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: done\n"+
			"freshness: {state: current, check: C-001, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: reviewed}\n"+
			"---\n\n# Write it\n")

	if problems := RunV2Checks(root).Project; len(problems) != 0 {
		t.Fatalf("RunV2Checks().Project = %v, want no problems for a clean V2 project with checks", problems)
	}
}

// TestCheckProject_ConsistencyDiagnostics proves doctor reports every
// evaluation-level inconsistency InspectTaskConsistency finds — done without
// current clearance, owner acceptance naming a superseded check, and
// evidence that clears completion while status lags — alongside structural
// load failures, in deterministic (task ID) order, without failing the load.
func TestCheckProject_ConsistencyDiagnostics(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2ObjectiveInProgress(t, root, "O-001-ship", "O-001", "Ship it")

	// T-001: marked done but no Check was ever recorded.
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-alpha.md"),
		"---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: done\n---\n\n# Alpha\n")

	// T-002: owner accepted C-002, but C-003 has since superseded it as latest.
	writeV2Check(t, root, "C-002", "{kind: task, id: T-002}", "CLEAR", "")
	writeV2Check(t, root, "C-003", "{kind: task, id: T-002}", "CLEAR", "C-002")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-002-beta.md"),
		"---\nid: T-002\ntitle: \"Beta\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: in_progress\nstage: audit\n"+
			"freshness: {state: current, check: C-003, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: rechecked}\n"+
			"owner_validation: {required: true, accepted_check: C-002, accepted_by: {role: owner, session: owner-1}}\n"+
			"---\n\n# Beta\n")

	// T-003: evidence clears completion (current clearance, no owner
	// validation required) but status was never advanced to done.
	writeV2Check(t, root, "C-004", "{kind: task, id: T-003}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-003-gamma.md"),
		"---\nid: T-003\ntitle: \"Gamma\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: in_progress\nstage: audit\n"+
			"freshness: {state: current, check: C-004, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: reviewed}\n"+
			"---\n\n# Gamma\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 3 {
		t.Fatalf("RunV2Checks().Project = %+v, want 3 consistency problems", problems)
	}

	wantOrder := []string{"[v2-done-without-clearance]", "[v2-acceptance-superseded]", "[v2-evidence-contradicts-status]"}
	for i, want := range wantOrder {
		if !strings.Contains(problems[i].Message, want) {
			t.Errorf("problems[%d].Message = %q, want containing %q (deterministic task-ID order)", i, problems[i].Message, want)
		}
		if problems[i].Repair == "" {
			t.Errorf("problems[%d].Repair is empty, want a typed repair suggestion", i)
		}
		if !strings.HasSuffix(problems[i].File, ".md") {
			t.Errorf("problems[%d].File = %q, want the Task's source record path", i, problems[i].File)
		}
	}
}

// TestCheckProject_MissingCheckVersusStaleVersusUnknownEvidence proves a done
// Task with no recorded Check, one whose freshness assessment marks its Check
// unknown, and one whose freshness assessment marks its Check stale all land under the missing-evidence health category rather than malformed
// data, and that their messages distinguish which of the three it is (E48
// T-006 AC: "A target with no recorded Check and a target with a stale or
// unknown freshness assessment report under missing evidence, distinguished
// from one another in the output").
func TestCheckProject_MissingCheckVersusStaleVersusUnknownEvidence(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2ObjectiveInProgress(t, root, "O-001-ship", "O-001", "Ship it")

	// T-001: done, no Check ever recorded — clearance missing.
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-alpha.md"),
		"---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: done\n---\n\n# Alpha\n")

	// T-002: done, latest Check is CLEAR but freshness marks it unknown — clearance unknown.
	writeV2Check(t, root, "C-001", "{kind: task, id: T-002}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-002-beta.md"),
		"---\nid: T-002\ntitle: \"Beta\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: done\n"+
			"freshness: {state: unknown, check: C-001, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: not reassessed}\n"+
			"---\n\n# Beta\n")

	// T-003: done, latest Check is CLEAR but freshness marks it stale — clearance stale.
	writeV2Check(t, root, "C-003", "{kind: task, id: T-003}", "CLEAR", "")
	testutil.WriteFile(t, filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-003-gamma.md"),
		"---\nid: T-003\ntitle: \"Gamma\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: done\n"+
			"freshness: {state: stale, check: C-003, assessed_by: {role: checker, session: s}, assessed_at: '2026-09-14T00:00:00Z', basis: code changed}\n"+
			"---\n\n# Gamma\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 3 {
		t.Fatalf("RunV2Checks().Project = %+v, want 3 problems (T-001 missing, T-002 unknown, T-003 stale)", problems)
	}

	wantDetail := []string{"clearance is missing", "clearance is unknown", "clearance is stale"}
	for i, want := range wantDetail {
		if !strings.Contains(problems[i].Message, want) {
			t.Errorf("problems[%d].Message = %q, want it to name %q", i, problems[i].Message, want)
		}
		if problems[i].Category != HealthMissingEvidence {
			t.Errorf("problems[%d].Category = %q, want %q: nothing is broken, evidence just isn't current", i, problems[i].Category, HealthMissingEvidence)
		}
	}
}

// TestCheckProject_ConsistencyReadOnly proves RunV2Checks never writes to a
// V2 project carrying Check and evidence records that trip a consistency
// diagnostic: every file's bytes and modification time stay unchanged.
func TestCheckProject_ConsistencyReadOnly(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	writeV2Objective(t, root, "O-001-ship", "O-001", "Ship it")
	taskPath := filepath.Join(root, "objectives", "O-001-ship", "tasks", "T-001-alpha.md")
	testutil.WriteFile(t, taskPath, "---\nid: T-001\ntitle: \"Alpha\"\nobjective: O-001\nstatus: done\n---\n\n# Alpha\n")

	type snapshot struct {
		bytes []byte
		mtime time.Time
	}
	before := map[string]snapshot{}
	for _, p := range []string{taskPath} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		before[p] = snapshot{bytes: raw, mtime: info.ModTime()}
	}

	if problems := RunV2Checks(root).Project; len(problems) == 0 {
		t.Fatal("RunV2Checks().Project = no problems, want the done-without-clearance diagnostic")
	}

	for p, want := range before {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if string(raw) != string(want.bytes) {
			t.Fatalf("RunV2Checks().Project changed %s bytes:\nbefore: %s\nafter: %s", p, want.bytes, raw)
		}
		if !info.ModTime().Equal(want.mtime) {
			t.Fatalf("RunV2Checks().Project changed %s mtime: before=%v after=%v", p, want.mtime, info.ModTime())
		}
	}
}

// TestCheckProject_ObjectiveAndIssueConsistencyReadOnly proves RunV2Checks
// never writes to a V2 project carrying Objective and Issue records that trip
// the new E44 consistency diagnostics: every file's bytes and modification
// time stay unchanged.
func TestCheckProject_ObjectiveAndIssueConsistencyReadOnly(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "schema_version: 2\n")
	objectivePath := filepath.Join(root, "objectives", "O-001-ship", "Objective.md")
	testutil.WriteFile(t, objectivePath, "---\nid: O-001\ntitle: \"Ship it\"\nstatus: done\n---\n\n# Ship it\n")
	taskPath := writeV2Task(t, root, "O-001-ship", "T-001-write.md", "T-001", "Write it", "O-001")
	checkC001Path := writeV2Check(t, root, "C-001", "{kind: task, id: T-001}", "CLEAR", "")
	checkC002Path := writeV2Check(t, root, "C-002", "{kind: task, id: T-001}", "CLEAR", "C-001")
	issuePath := writeV2Issue(t, root, "I-001-flaky.md", "I-001", "resolved", "defect",
		"checks: [C-001]\nresolution: {disposition: verified, check: C-001, actor: {role: checker, session: sess-1}, at: '2026-09-14T00:00:00Z'}\n")

	type snapshot struct {
		bytes []byte
		mtime time.Time
	}
	before := map[string]snapshot{}
	for _, p := range []string{objectivePath, taskPath, checkC001Path, checkC002Path, issuePath} {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		before[p] = snapshot{bytes: raw, mtime: info.ModTime()}
	}

	if problems := RunV2Checks(root).Project; len(problems) == 0 {
		t.Fatal("RunV2Checks().Project = no problems, want the objective and issue consistency diagnostics")
	}

	for p, want := range before {
		info, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if string(raw) != string(want.bytes) {
			t.Fatalf("RunV2Checks().Project changed %s bytes:\nbefore: %s\nafter: %s", p, want.bytes, raw)
		}
		if !info.ModTime().Equal(want.mtime) {
			t.Fatalf("RunV2Checks().Project changed %s mtime: before=%v after=%v", p, want.mtime, info.ModTime())
		}
	}
}

func TestCheckProject_SchemaVersionUnsupported(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "config.yml"), "quality_gates: {}\n")

	problems := RunV2Checks(root).Project
	if len(problems) != 1 || !strings.Contains(problems[0].Message, "[schema-version-unsupported]") {
		t.Fatalf("RunV2Checks().Project = %v, want 1 problem naming schema-version-unsupported", problems)
	}
}
