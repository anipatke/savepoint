package migrate

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
	"gopkg.in/yaml.v3"
)

// --- round trip -----------------------------------------------------------

// TestConvert_roundTripDecodesAllActiveTargets proves every Objective and
// Task target planned over both frozen fixtures renders content that decodes
// through the real strict V2 decoders with no diagnostic, and carries no
// fabricated evidence.
func TestConvert_roundTripDecodesAllActiveTargets(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			root := fixtureProjectRoot(fixture)
			p := mustPlan(t, root)

			for _, target := range p.Targets {
				switch target.Kind {
				case TargetObjective:
					content, err := ConvertObjective(root, p, target)
					if err != nil {
						t.Fatalf("ConvertObjective(%s) error = %v", target.GlobalID, err)
					}
					obj, err := data.DecodeObjectiveV2(target.TargetPath, content)
					if err != nil {
						t.Fatalf("DecodeObjectiveV2(%s) error = %v", target.GlobalID, err)
					}
					if obj.Evidence != nil {
						t.Errorf("objective %s: Evidence = %+v, want nil (no fabrication)", target.GlobalID, obj.Evidence)
					}
					if obj.Status != data.ColumnPlanned && obj.Status != data.ColumnInProgress {
						t.Errorf("objective %s: Status = %q, want planned or in_progress", target.GlobalID, obj.Status)
					}
				case TargetTask:
					content, err := ConvertTask(root, p, target)
					if err != nil {
						t.Fatalf("ConvertTask(%s) error = %v", target.GlobalID, err)
					}
					task, err := data.DecodeTaskV2(target.TargetPath, content)
					if err != nil {
						t.Fatalf("DecodeTaskV2(%s) error = %v", target.GlobalID, err)
					}
					if task.Evidence != nil {
						t.Errorf("task %s: Evidence = %+v, want nil (no fabrication)", target.GlobalID, task.Evidence)
					}
				case TargetIssue:
					content, err := ConvertIssue(root, p, target)
					if err != nil {
						t.Fatalf("ConvertIssue(%s) error = %v", target.GlobalID, err)
					}
					issue, err := data.DecodeIssueV2(target.TargetPath, content)
					if err != nil {
						t.Fatalf("DecodeIssueV2(%s) error = %v", target.GlobalID, err)
					}
					if issue.Origin.Check != "" {
						t.Errorf("issue %s: Origin.Check = %q, want empty (no fabrication)", target.GlobalID, issue.Origin.Check)
					}
				}
			}
		})
	}
}

// --- determinism ------------------------------------------------------

func TestConvert_deterministicAcrossRepeatedRuns(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			root := fixtureProjectRoot(fixture)
			p := mustPlan(t, root)

			for _, target := range p.Targets {
				var first, second string
				var err error
				switch target.Kind {
				case TargetObjective:
					first, err = ConvertObjective(root, p, target)
					if err == nil {
						second, err = ConvertObjective(root, p, target)
					}
				case TargetTask:
					first, err = ConvertTask(root, p, target)
					if err == nil {
						second, err = ConvertTask(root, p, target)
					}
				case TargetIssue:
					first, err = ConvertIssue(root, p, target)
					if err == nil {
						second, err = ConvertIssue(root, p, target)
					}
				default:
					continue
				}
				if err != nil {
					t.Fatalf("convert %s error = %v", target.GlobalID, err)
				}
				if first != second {
					t.Errorf("convert %s not deterministic:\nfirst=%q\nsecond=%q", target.GlobalID, first, second)
				}
			}
		})
	}
}

// --- Objective field mapping -----------------------------------------

func TestConvertObjective_activeEpicFields(t *testing.T) {
	cases := []struct {
		fixture       string
		wantTitle     string
		wantStatus    data.ColumnType
		wantRelease   string
		wantDependsOn []string
	}{
		{fixture: "v1-basic", wantTitle: "Epic E01: Example", wantStatus: data.ColumnInProgress, wantRelease: "R-001"},
		{fixture: "v1-history", wantTitle: "Epic E01: Example (v1.1 continuation)", wantStatus: data.ColumnInProgress, wantRelease: "R-002"},
	}

	for _, tc := range cases {
		t.Run(tc.fixture, func(t *testing.T) {
			root := fixtureProjectRoot(tc.fixture)
			p := mustPlan(t, root)

			var target PlannedTarget
			var found bool
			for _, tg := range p.Targets {
				if tg.Kind == TargetObjective {
					target, found = tg, true
					break
				}
			}
			if !found {
				t.Fatalf("no objective target planned for fixture %s", tc.fixture)
			}

			content, err := ConvertObjective(root, p, target)
			if err != nil {
				t.Fatalf("ConvertObjective error = %v", err)
			}
			obj, err := data.DecodeObjectiveV2(target.TargetPath, content)
			if err != nil {
				t.Fatalf("DecodeObjectiveV2 error = %v", err)
			}

			if obj.Title != tc.wantTitle {
				t.Errorf("Title = %q, want %q", obj.Title, tc.wantTitle)
			}
			if obj.Status != tc.wantStatus {
				t.Errorf("Status = %q, want %q", obj.Status, tc.wantStatus)
			}
			if obj.Release != tc.wantRelease {
				t.Errorf("Release = %q, want %q", obj.Release, tc.wantRelease)
			}
			if len(obj.DependsOn) != len(tc.wantDependsOn) {
				t.Errorf("DependsOn = %v, want %v", obj.DependsOn, tc.wantDependsOn)
			}
		})
	}
}

// TestConvertObjective_unauditedEpicWithAllTasksDone proves an epic whose
// every Task is done, but whose own status was never advanced to done or
// audited, converts to an in_progress Objective with no clearance — Objective
// status must come from the epic's own recorded status, never be computed
// from child Task completion.
func TestConvertObjective_unauditedEpicWithAllTasksDone(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		".savepoint/releases/v1/epics/E02-allDone/E02-Detail.md": "---\n" +
			"type: epic-design\n" +
			"status: in_progress\n" +
			"---\n\n# Epic E02: All Tasks Done\n\n## Purpose\n\nEvery Task is done but the epic itself was never audited.\n",
		".savepoint/releases/v1/epics/E02-allDone/tasks/T001-only.md": "---\n" +
			"id: E02-allDone/T001-only\n" +
			"status: done\n" +
			"objective: \"The only task, already done\"\n" +
			"depends_on: []\n" +
			"---\n\n# T001: Only task\n\n## Acceptance Criteria\n\n- [x] Done\n",
	})

	p := mustPlan(t, root)

	var target PlannedTarget
	var found bool
	for _, tg := range p.Targets {
		if tg.Kind == TargetObjective {
			target, found = tg, true
		}
	}
	if !found {
		t.Fatalf("expected an Objective target for the unaudited epic; targets = %+v", p.Targets)
	}

	content, err := ConvertObjective(root, p, target)
	if err != nil {
		t.Fatalf("ConvertObjective error = %v", err)
	}
	obj, err := data.DecodeObjectiveV2(target.TargetPath, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2 error = %v", err)
	}
	if obj.Status != data.ColumnInProgress {
		t.Errorf("Status = %q, want in_progress (epic's own recorded status, not derived from task completion)", obj.Status)
	}
	if obj.Evidence != nil {
		t.Errorf("Evidence = %+v, want nil: an unaudited epic must still earn an integration Check", obj.Evidence)
	}
}

// TestConvertObjective_dependsOnMapsToAllocatedObjective proves a declared V1
// epic dependency resolves to the depended-on epic's allocated O-###, and that
// an unresolvable reference is refused rather than silently dropped.
func TestConvertObjective_dependsOnMapsToAllocatedObjective(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		".savepoint/releases/v1/epics/E01-base/E01-Detail.md": "---\n" +
			"type: epic-design\nstatus: in_progress\n---\n\n# Epic E01: Base\n\n## Purpose\n\nBase epic.\n",
		".savepoint/releases/v1/epics/E02-dependent/E02-Detail.md": "---\n" +
			"type: epic-design\nstatus: planned\ndepends_on: [E01-base]\n---\n\n# Epic E02: Dependent\n\n## Purpose\n\nDepends on E01.\n",
	})

	p := mustPlan(t, root)

	base, ok := targetByPath(p, ".savepoint/releases/v1/epics/E01-base/E01-Detail.md")
	if !ok {
		t.Fatalf("expected E01-base to be planned as an objective")
	}
	dependent, ok := targetByPath(p, ".savepoint/releases/v1/epics/E02-dependent/E02-Detail.md")
	if !ok {
		t.Fatalf("expected E02-dependent to be planned as an objective")
	}

	content, err := ConvertObjective(root, p, dependent)
	if err != nil {
		t.Fatalf("ConvertObjective error = %v", err)
	}
	obj, err := data.DecodeObjectiveV2(dependent.TargetPath, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2 error = %v", err)
	}
	if len(obj.DependsOn) != 1 || obj.DependsOn[0] != base.GlobalID {
		t.Errorf("DependsOn = %v, want [%s]", obj.DependsOn, base.GlobalID)
	}
}

func TestConvertObjective_unresolvedDependsOnRefused(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, map[string]string{
		".savepoint/releases/v1/epics/E02-dependent/E02-Detail.md": "---\n" +
			"type: epic-design\nstatus: planned\ndepends_on: [E99-missing]\n---\n\n# Epic E02: Dependent\n\n## Purpose\n\nDepends on a nonexistent epic.\n",
	})

	p := mustPlan(t, root)
	dependent, ok := targetByPath(p, ".savepoint/releases/v1/epics/E02-dependent/E02-Detail.md")
	if !ok {
		t.Fatalf("expected E02-dependent to be planned as an objective")
	}

	_, err := ConvertObjective(root, p, dependent)
	if !errors.Is(err, ErrUnresolvedDependency) {
		t.Fatalf("ConvertObjective error = %v, want ErrUnresolvedDependency", err)
	}
}

// --- Task field mapping ------------------------------------------------

// TestConvertTask_activeFields covers ownership, status/stage mapping
// (including the legacy phase: implementation alias), and depends_on shape
// across both fixtures' active Tasks.
func TestConvertTask_activeFields(t *testing.T) {
	t.Run("v1-basic T002-follow-up: phase alias, legacy prerequisite not a dependency", func(t *testing.T) {
		root := fixtureProjectRoot("v1-basic")
		p := mustPlan(t, root)
		target, ok := targetByPath(p, ".savepoint/releases/v1/epics/E01-example/tasks/T002-follow-up.md")
		if !ok {
			t.Fatalf("expected T002-follow-up to be planned as an active task")
		}

		content, err := ConvertTask(root, p, target)
		if err != nil {
			t.Fatalf("ConvertTask error = %v", err)
		}
		task, err := data.DecodeTaskV2(target.TargetPath, content)
		if err != nil {
			t.Fatalf("DecodeTaskV2 error = %v", err)
		}

		if task.Status != data.ColumnInProgress {
			t.Errorf("Status = %q, want in_progress", task.Status)
		}
		if task.Stage != data.StageBuild {
			t.Errorf("Stage = %q, want build (mapped from legacy phase: implementation)", task.Stage)
		}
		if len(task.DependsOn) != 0 {
			t.Errorf("DependsOn = %v, want empty: the V1 dependency was archived, done work, which becomes a legacy prerequisite, not a depends_on entry", task.DependsOn)
		}
		if task.Objective == "" {
			t.Errorf("Objective owner is empty")
		}
		if task.Release != "v1" {
			t.Errorf("Release = %q, want v1", task.Release)
		}
	})

	t.Run("v1-history T002-follow-up: planned status carries no stage, depends_on names the sibling active task", func(t *testing.T) {
		root := fixtureProjectRoot("v1-history")
		p := mustPlan(t, root)
		t001, ok := targetByPath(p, ".savepoint/releases/v1.1/epics/E01-example/tasks/T001-shared.md")
		if !ok {
			t.Fatalf("expected v1.1 T001-shared to be planned as an active task")
		}
		t002, ok := targetByPath(p, ".savepoint/releases/v1.1/epics/E01-example/tasks/T002-follow-up.md")
		if !ok {
			t.Fatalf("expected v1.1 T002-follow-up to be planned as an active task")
		}

		content, err := ConvertTask(root, p, t002)
		if err != nil {
			t.Fatalf("ConvertTask error = %v", err)
		}
		task, err := data.DecodeTaskV2(t002.TargetPath, content)
		if err != nil {
			t.Fatalf("DecodeTaskV2 error = %v", err)
		}

		if task.Status != data.ColumnPlanned {
			t.Errorf("Status = %q, want planned", task.Status)
		}
		if task.Stage != "" {
			t.Errorf("Stage = %q, want empty: stage is only meaningful when status is in_progress", task.Stage)
		}
		if len(task.DependsOn) != 1 || task.DependsOn[0].Task != t001.GlobalID {
			t.Errorf("DependsOn = %v, want [%s] (the active v1.1 T001-shared, not the done v1 one of the same short id)", task.DependsOn, t001.GlobalID)
		}
	})
}

func TestOwnerObjectiveID_extractsHyphenatedIDs(t *testing.T) {
	tests := []struct {
		path string
		want string
		ok   bool
	}{
		{"objectives/O-001-example/tasks/T-001-work", "O-001", true},
		{"objectives/O-1234-example/tasks/T-001-work", "O-1234", true},
		{"objectives/O-01-example/tasks/T-001-work", "", false},
		{"objectives/O-abc-example/tasks/T-001-work", "", false},
		{"tasks/T-001-work", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got, ok := ownerObjectiveID(tt.path)
			if got != tt.want || ok != tt.ok {
				t.Errorf("ownerObjectiveID(%q) = (%q, %t), want (%q, %t)", tt.path, got, ok, tt.want, tt.ok)
			}
		})
	}
}

// TestConvertTask_legacyPrerequisiteLine proves a Task whose V1 dependency
// was an archived, completed Task carries an authored line naming the
// archive path and the recorded completion evidence.
func TestConvertTask_legacyPrerequisiteLine(t *testing.T) {
	root := fixtureProjectRoot("v1-basic")
	p := mustPlan(t, root)
	target, ok := targetByPath(p, ".savepoint/releases/v1/epics/E01-example/tasks/T002-follow-up.md")
	if !ok {
		t.Fatalf("expected T002-follow-up to be planned as an active task")
	}

	content, err := ConvertTask(root, p, target)
	if err != nil {
		t.Fatalf("ConvertTask error = %v", err)
	}

	wantArchive := ".savepoint/archive/v1/.savepoint/releases/v1/epics/E01-example/tasks/T001-original.md"
	if !strings.Contains(content, wantArchive) {
		t.Errorf("rendered content does not name the archive path %q:\n%s", wantArchive, content)
	}
	if !strings.Contains(content, "Baseline output exists") {
		t.Errorf("rendered content does not carry the recorded completion evidence:\n%s", content)
	}
	if !strings.Contains(content, "## Legacy Prerequisite") {
		t.Errorf("rendered content missing the Legacy Prerequisite heading:\n%s", content)
	}
}

// --- unrecognized lifecycle: no healing --------------------------------

// TestConvertTask_unrecognizedStatusRefusesRatherThanHeals proves that an
// unrecognized raw status blocks rendering with ErrAmbiguousLifecycle, even
// though data.ParseTaskFile's own load-time path would silently heal the
// same content to status: planned.
func TestConvertTask_unrecognizedStatusRefusesRatherThanHeals(t *testing.T) {
	root := t.TempDir()
	rel := ".savepoint/releases/v1/epics/E01-x/tasks/T001-bad.md"
	content := "---\nid: E01-x/T001-bad\nstatus: mid-review\nobjective: \"bad status\"\n---\n\n# T001\n"
	writeFiles(t, root, map[string]string{rel: content})

	// Prove the healing path this must not be reachable from actually exists
	// and would have healed the same content.
	parsed, err := data.NewParser().ParseTaskFile(rel, content)
	if err != nil {
		t.Fatalf("ParseTaskFile error = %v", err)
	}
	if parsed.Column != data.ColumnPlanned {
		t.Fatalf("setup invariant broken: ParseTaskFile healed %q to %q, want planned", "mid-review", parsed.Column)
	}

	target := PlannedTarget{
		Kind:       TargetTask,
		GlobalID:   "T-001",
		Legacy:     LegacyKey{Release: "v1", Epic: "E01-x", Path: rel, OriginalID: "E01-x/T001-bad"},
		TargetPath: "objectives/O-001-x/tasks/T-001-bad",
	}
	_, err = ConvertTask(root, &ConversionPlan{}, target)
	if !errors.Is(err, ErrAmbiguousLifecycle) {
		t.Fatalf("ConvertTask error = %v, want ErrAmbiguousLifecycle", err)
	}
}

// TestConvertTask_unrecognizedStageRefusesRatherThanHeals covers the one
// lifecycle dimension plan.go itself does not screen: an in_progress Task
// with a stage value that is neither canonical nor the legacy
// implementation alias.
func TestConvertTask_unrecognizedStageRefusesRatherThanHeals(t *testing.T) {
	root := t.TempDir()
	rel := ".savepoint/releases/v1/epics/E01-x/tasks/T001-bad.md"
	content := "---\nid: E01-x/T001-bad\nstatus: in_progress\nstage: reviewing\nobjective: \"bad stage\"\n---\n\n# T001\n"
	writeFiles(t, root, map[string]string{rel: content})

	target := PlannedTarget{
		Kind:       TargetTask,
		GlobalID:   "T-001",
		Legacy:     LegacyKey{Release: "v1", Epic: "E01-x", Path: rel, OriginalID: "E01-x/T001-bad"},
		TargetPath: "objectives/O-001-x/tasks/T-001-bad",
	}
	_, err := ConvertTask(root, &ConversionPlan{}, target)
	if !errors.Is(err, ErrAmbiguousLifecycle) {
		t.Fatalf("ConvertTask error = %v, want ErrAmbiguousLifecycle", err)
	}
}

// --- verbatim body, CRLF, and unknown-field preservation ----------------

// TestConvertTask_verbatimCRLFBodyAndUnknownFields proves the CRLF fixture's
// authored body (including its multiline text and body comment) is relocated
// byte for byte, and that its unrecognized frontmatter fields (legacy_note,
// the nested metadata block) survive on the converted record.
func TestConvertTask_verbatimCRLFBodyAndUnknownFields(t *testing.T) {
	root := fixtureProjectRoot("v1-basic")
	rel := ".savepoint/releases/v1/epics/E01-example/tasks/T002-follow-up.md"
	target, ok := targetByPath(mustPlan(t, root), rel)
	if !ok {
		t.Fatalf("expected T002-follow-up to be planned as an active task")
	}

	rawBytes, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatalf("read fixture source: %v", err)
	}
	raw := string(rawBytes)
	if !strings.Contains(raw, "\r\n") {
		t.Fatalf("test setup invariant broken: fixture source is not CRLF")
	}

	p := mustPlan(t, root)
	content, err := ConvertTask(root, p, target)
	if err != nil {
		t.Fatalf("ConvertTask error = %v", err)
	}

	if !strings.Contains(content, "\r\n") {
		t.Errorf("rendered content lost the source's CRLF line endings")
	}

	// The exact authored comment and the exact multiline wrapped Acceptance
	// Criteria/Implementation Plan text, byte for byte including their
	// original \r\n, must appear unchanged in the rendered output.
	verbatimSnippets := []string{
		"<!-- Author note: this task intentionally keeps v1-era authoring shape,\r\n     including the phase field and this comment, as migration source\r\n     evidence. Do not \"fix\" it to match current templates. -->",
		"- [ ] Follow-up documented across more than\r\n      one line, describing the exact expected\r\n      behavior for the dependent step",
		"- [ ] Extend it with the following steps: first\r\n      gather inputs, then apply the follow-up\r\n      transform, then record results",
	}
	for _, snippet := range verbatimSnippets {
		if !strings.Contains(content, snippet) {
			t.Errorf("rendered content missing verbatim snippet:\n%s\n\nfull content:\n%s", snippet, content)
		}
	}

	task, err := data.DecodeTaskV2(target.TargetPath, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2 error = %v", err)
	}
	if !task.Source.CRLF {
		t.Errorf("decoded Source.CRLF = false, want true")
	}

	var fm struct {
		LegacyFields struct {
			LegacyNote string `yaml:"legacy_note"`
			Metadata   struct {
				Origin   string `yaml:"origin"`
				Reviewer struct {
					Name string `yaml:"name"`
					Role string `yaml:"role"`
				} `yaml:"reviewer"`
			} `yaml:"metadata"`
		} `yaml:"legacy_fields"`
	}
	fmYAML, _, err := data.SplitFrontmatterBody(content)
	if err != nil {
		t.Fatalf("SplitFrontmatterBody error = %v", err)
	}
	if err := yaml.Unmarshal([]byte(fmYAML), &fm); err != nil {
		t.Fatalf("parse rendered frontmatter: %v", err)
	}
	if fm.LegacyFields.LegacyNote != "kept from v1 authoring; not part of the current schema" {
		t.Errorf("legacy_fields.legacy_note = %q, not preserved", fm.LegacyFields.LegacyNote)
	}
	if fm.LegacyFields.Metadata.Origin != "hand-authored" {
		t.Errorf("legacy_fields.metadata.origin = %q, not preserved", fm.LegacyFields.Metadata.Origin)
	}
	if fm.LegacyFields.Metadata.Reviewer.Name != "sam" || fm.LegacyFields.Metadata.Reviewer.Role != "peer" {
		t.Errorf("legacy_fields.metadata.reviewer = %+v, not preserved", fm.LegacyFields.Metadata.Reviewer)
	}
}

// --- test helpers ---------------------------------------------------------

func writeFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("mkdir for %s: %v", rel, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
}
