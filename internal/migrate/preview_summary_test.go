package migrate

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFormatSummaryPreview_readyPlanIsShortAndCollapsed(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-history"))

	got := FormatSummaryPreview(p)
	for _, want := range []string{"Will create", "Objectives", "Will archive", "Status:", "--verbose"} {
		if !strings.Contains(got, want) {
			t.Errorf("summary missing %q:\n%s", want, got)
		}
	}
	for _, archive := range p.Archives {
		if strings.Contains(got, archive.ArchivePath) {
			t.Fatalf("summary lists archive path %s; archives must be collapsed to a count:\n%s", archive.ArchivePath, got)
		}
	}
	if full := FormatPreview(p); len(got) >= len(full) {
		t.Errorf("summary (%d bytes) is not shorter than the full listing (%d bytes)", len(got), len(full))
	}
	if again := FormatSummaryPreview(p); again != got {
		t.Error("FormatSummaryPreview is not deterministic for the same plan")
	}
}

func TestFormatSummaryPreview_blockedPlanListsDecisionsWithEntries(t *testing.T) {
	root := t.TempDir()
	writeMinimalV1Project(t, root)
	writeTaskWithStatus(t, root, "v1", "E01-x", "T001-open", "planned", "")
	detail := ".savepoint/releases/v1/epics/E01-x/E01-Detail.md"
	writeFile(t, filepath.Join(root, filepath.FromSlash(detail)), "---\nstatus: deferred\n---\n\n# E01\n")

	p := mustPlan(t, root)
	if p.Appliable {
		t.Fatal("plan with a deferred epic is appliable, want blocked")
	}
	a, ok := findAmbiguityByPath(p, detail)
	if !ok {
		t.Fatalf("no ambiguity for %s", detail)
	}

	got := FormatSummaryPreview(p)
	for _, want := range []string{
		"Decisions needed (1)",
		detail,
		"choices: planned, in_progress, done, audited",
		"    - id: " + a.ID + "\n      value: <planned | in_progress | done | audited>",
		"Status: blocked until 1 decision(s)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("blocked summary missing %q:\n%s", want, got)
		}
	}

	decided, err := Plan(root, Decisions{a.ID: {Value: "planned", SourceFile: "decisions.yml", DecidedAt: time.Now()}}, fixedClock(time.Now()), fixedOperationID("op"))
	if err != nil {
		t.Fatalf("Plan() with decision error = %v", err)
	}
	after := FormatSummaryPreview(decided)
	for _, want := range []string{"Decisions applied (1)", detail + " -> planned  (from decisions.yml)", "status planned by your decision", "Status: ready"} {
		if !strings.Contains(after, want) {
			t.Errorf("decided summary missing %q:\n%s", want, after)
		}
	}
	if strings.Contains(after, "Decisions needed") {
		t.Errorf("decided summary still lists decisions needed:\n%s", after)
	}
}
