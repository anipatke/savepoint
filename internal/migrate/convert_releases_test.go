package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/data"
)

func releaseTargets(p *ConversionPlan) []PlannedTarget {
	var targets []PlannedTarget
	for _, target := range p.Targets {
		if target.Kind == TargetRelease {
			targets = append(targets, target)
		}
	}
	return targets
}

func TestConvertRelease_activeAndHistoricalForms(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			root := fixtureProjectRoot(fixture)
			plan := mustPlan(t, root)
			targets := releaseTargets(plan)
			if len(targets) == 0 {
				t.Fatal("plan produced no Release targets")
			}

			for _, target := range targets {
				content, err := ConvertRelease(root, plan, target)
				if err != nil {
					t.Fatalf("ConvertRelease(%s) error = %v", target.GlobalID, err)
				}
				release, err := data.DecodeReleaseV2(target.TargetPath, content)
				if err != nil {
					t.Fatalf("DecodeReleaseV2(%s) error = %v\n%s", target.GlobalID, err, content)
				}
				if release.ID != target.GlobalID {
					t.Errorf("Release ID = %q, want %q", release.ID, target.GlobalID)
				}
				if !target.Generated && !strings.HasPrefix(release.ID, "R-") {
					t.Errorf("converted V1 Release ID = %q, want an unchanged R-### identity", release.ID)
				}
				if !strings.Contains(release.Source.Body, "## Legacy Source (verbatim)") {
					t.Error("Release body has no accountable legacy-source section")
				}

				if target.Legacy.Release == "v1" && fixture == "v1-history" {
					if release.Status != data.ColumnDone || release.LegacyCompletion == nil {
						t.Fatalf("historical Release = %+v, want done with typed legacy completion", release)
					}
					if release.Evidence != nil {
						t.Error("historical Release has fabricated V2 evidence")
					}
				} else if release.Status != data.ColumnInProgress || release.LegacyCompletion != nil {
					t.Errorf("active Release = %+v, want in_progress without legacy completion", release)
				}
			}
		})
	}
}

func TestConvertRelease_rendersGeneratedContinuationGoalWithoutV1Source(t *testing.T) {
	target := PlannedTarget{
		Kind:           TargetRelease,
		GlobalID:       "G-001",
		TargetPath:     "releases/G-001-continued-after-migration/Release.md",
		ReleaseStatus:  string(data.ColumnInProgress),
		Generated:      true,
		GeneratedTitle: continuationGoalTitle,
	}
	content, err := ConvertRelease("", nil, target)
	if err != nil {
		t.Fatalf("ConvertRelease() error = %v", err)
	}
	if !strings.Contains(content, "title: "+continuationGoalTitle) {
		t.Errorf("generated Goal title missing:\n%s", content)
	}
	for _, section := range []string{"## Outcome", "## Why", "## Success Conditions", "## Boundaries", "TODO:"} {
		if !strings.Contains(content, section) {
			t.Errorf("generated Goal missing stub section %q:\n%s", section, content)
		}
	}
	if strings.Contains(content, "Legacy Source") || strings.Contains(content, "legacy_completion:") {
		t.Errorf("generated Goal incorrectly claims a V1 source:\n%s", content)
	}
	if _, err := data.DecodeReleaseV2(target.TargetPath, content); err != nil {
		t.Fatalf("generated Goal does not strict-load: %v", err)
	}
}

func TestConvertRelease_preservesAuthoredHeadingsInsideLegacyFence(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v9", "v9-PRD.md"), "---\nname: Custom\nstatus: in_progress\n---\n\n# Custom\n\n## Authored Section\n\nKeep this section.\n\n```\n## Code heading\n```\n")
	p := mustPlan(t, root)
	if len(releaseTargets(p)) != 1 {
		t.Fatalf("Release targets = %+v, want one", p.Targets)
	}
	content, err := ConvertRelease(root, p, releaseTargets(p)[0])
	if err != nil {
		t.Fatalf("ConvertRelease() error = %v", err)
	}
	if !strings.Contains(content, "## Authored Section") || !strings.Contains(content, "## Code heading") {
		t.Fatalf("rendered Release lost authored headings:\n%s", content)
	}
	if _, err := data.DecodeReleaseV2(releaseTargets(p)[0].TargetPath, content); err != nil {
		t.Fatalf("DecodeReleaseV2() rejected preserved authored headings: %v", err)
	}
}

func TestPlan_releasePRDIdentityAndObjectiveReferenceAreStable(t *testing.T) {
	root := fixtureProjectRoot("v1-history")
	first := mustPlan(t, root)
	second := mustPlan(t, root)

	firstTargets := releaseTargets(first)
	secondTargets := releaseTargets(second)
	if len(firstTargets) != len(secondTargets) {
		t.Fatalf("release target count changed: %d -> %d", len(firstTargets), len(secondTargets))
	}
	for i := range firstTargets {
		if firstTargets[i].GlobalID != secondTargets[i].GlobalID || firstTargets[i].TargetPath != secondTargets[i].TargetPath {
			t.Errorf("release target %d changed across plans: %+v -> %+v", i, firstTargets[i], secondTargets[i])
		}
	}
	for _, target := range first.Targets {
		if target.Kind == TargetObjective && target.ReleaseID == "" {
			t.Errorf("Objective %s has no allocated ReleaseID: %+v", target.GlobalID, target)
		}
	}
}

func TestPlan_duplicateReleasePRDIsBlockingAndWriteFree(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v9", "v9-PRD.md"), "---\nname: One\nstatus: in_progress\n---\n\n# One\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v9", "alternate-PRD.md"), "---\nname: Two\nstatus: in_progress\n---\n\n# Two\n")
	before := snapshotTree(t, root)
	p := mustPlan(t, root)
	if p.Appliable {
		t.Fatal("Plan.Appliable = true, want false for duplicate release PRDs")
	}
	if len(p.UnresolvedBlockingIDs) != 1 || !strings.Contains(p.UnresolvedBlockingIDs[0], string(AmbiguityDuplicateReleaseSource)) {
		t.Fatalf("UnresolvedBlockingIDs = %v, want one duplicate release source", p.UnresolvedBlockingIDs)
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestPlan_ambiguousReleaseDispositionBlocksCutover(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	releasePath := filepath.Join(root, ".savepoint", "releases", "v9", "v9-PRD.md")
	writeFile(t, releasePath, "---\nname: Ambiguous\nstatus: audited\n---\n\n# Ambiguous\n")
	before := snapshotTree(t, root)

	plan := mustPlan(t, root)
	if plan.Appliable {
		t.Fatal("Plan.Appliable = true, want false for an ambiguous Release disposition")
	}
	if len(plan.UnresolvedBlockingIDs) != 1 || !strings.Contains(plan.UnresolvedBlockingIDs[0], string(AmbiguityReleaseCompletion)) {
		t.Fatalf("UnresolvedBlockingIDs = %v, want Release completion ambiguity", plan.UnresolvedBlockingIDs)
	}
	if _, err := Apply(root, plan); err == nil {
		t.Fatal("Apply() error = nil, want cutover refusal while Release disposition is ambiguous")
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))
}

func TestPlan_auditedReleaseWithActiveEpicShowsDecisionInPreview(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, ".savepoint", "config.yml"), "quality_gates: {}\n")
	writeFile(t, filepath.Join(root, ".savepoint", "router.md"), routerFixtureContent(
		"epic-task-breakdown", "v1", "E02-active", "", `"Plan the active epic."`))
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "v1-PRD.md"),
		"---\nname: Ambiguous\nstatus: audited\n---\n\n# Ambiguous\n")
	writeFile(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E02-active", "E02-Detail.md"),
		"---\nstatus: in_progress\n---\n\n# Active epic\n")

	plan, err := Plan(root, nil, fixedClock(time.Now()), fixedOperationID("op"))
	if err != nil {
		t.Fatalf("Plan() error = %v, want the lifecycle decision to appear in the preview", err)
	}
	if plan.Appliable || len(plan.UnresolvedBlockingIDs) == 0 {
		t.Fatalf("Plan = %+v, want a blocked plan with an unresolved release decision", plan)
	}
	if !strings.Contains(strings.Join(plan.UnresolvedBlockingIDs, " "), string(AmbiguityReleaseCompletion)) {
		t.Errorf("UnresolvedBlockingIDs = %v, want the audited release completion decision", plan.UnresolvedBlockingIDs)
	}

	preview := FormatPreview(plan)
	for _, want := range []string{"audited without a typed completion decision", "choices: in_progress, done", "owner lifecycle decision"} {
		if !strings.Contains(preview, want) {
			t.Errorf("preview is missing %q:\n%s", want, preview)
		}
	}
	if strings.Contains(preview, "no resolvable V2 Goal") {
		t.Errorf("preview masks the lifecycle decision with a Goal error:\n%s", preview)
	}
}

func TestSourceReleaseIdentity_reservesRNumbers(t *testing.T) {
	cases := map[string]string{
		"v1":       "",
		"R001-old": "R001",
		"R12-old":  "",
		"R001":     "R001",
	}
	for input, want := range cases {
		if got := sourceReleaseIdentity(input); got != want {
			t.Errorf("sourceReleaseIdentity(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestConvertRelease_historicalArchiveHashMatchesSource(t *testing.T) {
	root := fixtureProjectRoot("v1-history")
	p := mustPlan(t, root)
	var target PlannedTarget
	for _, candidate := range releaseTargets(p) {
		if candidate.Legacy.Release == "v1" {
			target = candidate
		}
	}
	content, err := ConvertRelease(root, p, target)
	if err != nil {
		t.Fatal(err)
	}
	release, err := data.DecodeReleaseV2(target.TargetPath, content)
	if err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(target.Legacy.Path)))
	if err != nil {
		t.Fatal(err)
	}
	if release.LegacyCompletion.SHA256 != hashBytes(original) {
		t.Fatalf("legacy completion hash = %s, want %s", release.LegacyCompletion.SHA256, hashBytes(original))
	}
}
