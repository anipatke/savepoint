package migrate

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// These tests prove an owner decision the plan reports as resolved is
// actually applied (I-067): the record it names is planned or archived as
// decided, never left behind unconverted.

func planWithDecision(t *testing.T, root, path string, value string) *ConversionPlan {
	t.Helper()
	before := mustPlan(t, root)
	a, ok := findAmbiguityByPath(before, path)
	if !ok {
		t.Fatalf("no ambiguity raised for %s; Ambiguities = %+v", path, before.Ambiguities)
	}
	after, err := Plan(root, Decisions{a.ID: {Value: value, SourceFile: "decisions.yml", DecidedAt: time.Now()}}, fixedClock(time.Now()), fixedOperationID("op"))
	if err != nil {
		t.Fatalf("Plan() with decision %s=%s error = %v", a.ID, value, err)
	}
	if !after.Appliable {
		t.Fatalf("Appliable = false after deciding %s; UnresolvedBlockingIDs = %v", a.ID, after.UnresolvedBlockingIDs)
	}
	assertResolvedDecisionsPlaceTheirSources(t, after)
	return after
}

// assertResolvedDecisionsPlaceTheirSources is the invariant I-067 broke: a
// resolved blocking ambiguity on a source file leaves that file either planned
// as a target or archived.
func assertResolvedDecisionsPlaceTheirSources(t *testing.T, p *ConversionPlan) {
	t.Helper()
	for _, a := range p.Ambiguities {
		if !a.Blocking || !a.Resolved || !strings.HasSuffix(a.Path, ".md") {
			continue
		}
		_, planned := targetByPath(p, a.Path)
		_, archived := archiveByPath(p, a.Path)
		if !planned && !archived {
			t.Errorf("ambiguity %s is resolved as %q but %s is neither planned nor archived", a.ID, a.Decision, a.Path)
		}
	}
}

func TestPlan_decidedEpicStatusPlansItsObjectiveAndTasks(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-open", "planned", "")
	detailPath := ".savepoint/releases/v1/epics/E01-x/E01-Detail.md"
	writeFile(t, filepath.Join(root, filepath.FromSlash(detailPath)), "---\nstatus: deferred\n---\n\n# E01\n")

	p := planWithDecision(t, root, detailPath, "planned")

	objective, ok := targetByPath(p, detailPath)
	if !ok || objective.Kind != TargetObjective || objective.DecidedStatus != "planned" {
		t.Fatalf("decided epic target = %+v (found %v), want a planned Objective", objective, ok)
	}
	if task, ok := targetByPath(p, taskPath); !ok || task.Kind != TargetTask {
		t.Fatalf("task under the decided epic = %+v (found %v), want a Task target", task, ok)
	}
	content, err := ConvertObjective(root, p, objective)
	if err != nil {
		t.Fatalf("ConvertObjective() error = %v", err)
	}
	if !strings.Contains(content, "status: planned") {
		t.Errorf("converted Objective does not carry the decided status:\n%s", content)
	}
}

func TestPlan_decidedTaskStatusPlansTheTask(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")
	// The shared fixture omits the title the converter requires.
	writeFile(t, filepath.Join(root, filepath.FromSlash(taskPath)),
		"---\nid: E01-x/T001-weird\ntitle: Weird status\nstatus: escalated\ndepends_on: []\n---\n\n# T001\n")

	p := planWithDecision(t, root, taskPath, "planned")

	task, ok := targetByPath(p, taskPath)
	if !ok || task.Kind != TargetTask || task.DecidedStatus != "planned" {
		t.Fatalf("decided task target = %+v (found %v), want a planned Task", task, ok)
	}
	content, err := ConvertTask(root, p, task)
	if err != nil {
		t.Fatalf("ConvertTask() error = %v", err)
	}
	if !strings.Contains(content, "status: planned") {
		t.Errorf("converted Task does not carry the decided status:\n%s", content)
	}
}

func TestPlan_decidedDoneTaskStatusArchivesTheTask(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	taskPath := writeTaskWithStatus(t, root, "v1", "E01-x", "T001-weird", "escalated", "")

	p := planWithDecision(t, root, taskPath, "done")

	if _, ok := archiveByPath(p, taskPath); !ok {
		t.Fatalf("task decided done was not archived")
	}
	if _, ok := targetByPath(p, taskPath); ok {
		t.Errorf("task decided done was also planned as a target")
	}
}

func TestPlan_unresolvedNarrativeFindingDecisions(t *testing.T) {
	t.Run("archive", func(t *testing.T) {
		root := t.TempDir()
		writeMinimalV1Project(t, root)
		path := writeFinding(t, root, "F099", "Orphaned duplicate report", "duplicate", "F404")

		p := planWithDecision(t, root, path, "archive")
		if _, ok := archiveByPath(p, path); !ok {
			t.Fatal("finding decided archive was not archived")
		}
	})
	t.Run("treat_as_original", func(t *testing.T) {
		root := t.TempDir()
		writeMinimalV1Project(t, root)
		path := writeFinding(t, root, "F099", "Orphaned duplicate report", "duplicate", "F404")

		p := planWithDecision(t, root, path, "treat_as_original")
		issue, ok := targetByPath(p, path)
		if !ok || issue.Kind != TargetIssue {
			t.Fatalf("finding decided treat_as_original = %+v (found %v), want an Issue target", issue, ok)
		}
		content, err := ConvertIssue(root, p, issue)
		if err != nil {
			t.Fatalf("ConvertIssue() error = %v", err)
		}
		if !strings.Contains(content, "status: open") || strings.Contains(content, "duplicate_of") {
			t.Errorf("converted Issue is not an open original:\n%s", content)
		}
	})
}

func TestPlan_duplicateTaskIdentityDecisions(t *testing.T) {
	setup := func(t *testing.T) (string, string, string) {
		root := t.TempDir()
		writeMinimalV1Project(t, root)
		epicDir := filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-x")
		writeFile(t, filepath.Join(epicDir, "E01-Detail.md"), "---\nstatus: in_progress\n---\n\n# E01\n")
		first := ".savepoint/releases/v1/epics/E01-x/tasks/T001-a.md"
		second := ".savepoint/releases/v1/epics/E01-x/tasks/T001-b.md"
		writeFile(t, filepath.Join(root, filepath.FromSlash(first)), "---\nid: E01-x/T001-dup\nstatus: planned\ndepends_on: []\n---\n\n# T001 (a)\n")
		writeFile(t, filepath.Join(root, filepath.FromSlash(second)), "---\nid: E01-x/T001-dup\nstatus: planned\ndepends_on: []\n---\n\n# T001 (b)\n")
		return root, first, second
	}

	t.Run("keep_first", func(t *testing.T) {
		root, first, second := setup(t)
		p := planWithDecision(t, root, second, "keep_first")
		if _, ok := targetByPath(p, first); !ok {
			t.Error("keep_first did not plan the first declaration")
		}
		if _, ok := archiveByPath(p, second); !ok {
			t.Error("keep_first did not archive the second declaration")
		}
	})
	t.Run("keep_second", func(t *testing.T) {
		root, first, second := setup(t)
		before := mustPlan(t, root)
		firstTarget, _ := targetByPath(before, first)

		p := planWithDecision(t, root, second, "keep_second")
		if _, ok := targetByPath(p, first); ok {
			t.Error("keep_second still plans the first declaration")
		}
		if _, ok := archiveByPath(p, first); !ok {
			t.Error("keep_second did not archive the first declaration")
		}
		kept, ok := targetByPath(p, second)
		if !ok {
			t.Fatal("keep_second did not plan the second declaration")
		}
		if kept.GlobalID != firstTarget.GlobalID {
			t.Errorf("keep_second planned %s, want the identity %s the first declaration held", kept.GlobalID, firstTarget.GlobalID)
		}
	})
}
