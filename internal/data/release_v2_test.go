package data

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

const validReleaseBody = `
# Release

## Outcome

Ship the first-class release boundary.

## Why

The delivery promise needs an identity.

## Success Conditions

- Every member Objective is complete.

## Boundaries

Release does not own Tasks.
`

func validReleaseContent(id, title, status string) string {
	return "---\nid: " + id + "\ntitle: \"" + title + "\"\nstatus: " + status + "\n---\n" + validReleaseBody
}

func writeV2ReleaseFixture(t *testing.T, root, dirName, id, title string) {
	t.Helper()
	testutil.WriteFile(t, filepath.Join(root, v2ReleasesDirName, dirName, v2ReleaseFileName), validReleaseContent(id, title, "planned"))
}

func TestDecodeReleaseV2_valid(t *testing.T) {
	release, err := DecodeReleaseV2("releases/R001-first/Release.md", validReleaseContent("R001", "First-class releases", "in_progress"))
	if err != nil {
		t.Fatalf("DecodeReleaseV2() error = %v", err)
	}
	if release.ID != "R001" || release.Title != "First-class releases" || release.Status != ColumnInProgress {
		t.Fatalf("Release = %+v, want R001/title/in_progress", release)
	}
	if !strings.Contains(release.Outcome, "first-class release boundary") {
		t.Errorf("Outcome = %q, want decoded Outcome section", release.Outcome)
	}
	if !strings.Contains(release.SuccessConditions, "member Objective") {
		t.Errorf("SuccessConditions = %q, want decoded Success Conditions section", release.SuccessConditions)
	}
	if !strings.Contains(release.Source.Body, "## Boundaries") {
		t.Error("Source.Body did not retain the authored Release sections")
	}
}

func TestDecodeReleaseV2_rejectsMalformedIdentityLifecycleAndBody(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr error
	}{
		{
			name:    "malformed identity",
			content: validReleaseContent("release-1", "Release", "planned"),
			wantErr: ErrV2InvalidID,
		},
		{
			name:    "invalid lifecycle",
			content: validReleaseContent("R001", "Release", "review"),
			wantErr: ErrV2InvalidLifecycle,
		},
		{
			name:    "missing body section",
			content: strings.Replace(validReleaseContent("R001", "Release", "planned"), "## Why\n\nThe delivery promise needs an identity.\n", "", 1),
			wantErr: ErrV2ReleaseMissingSection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodeReleaseV2("releases/R001-release/Release.md", tt.content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeReleaseV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscoverV2Releases_absentDirectoryIsEmpty(t *testing.T) {
	releases, err := DiscoverV2Releases(t.TempDir())
	if err != nil {
		t.Fatalf("DiscoverV2Releases() error = %v", err)
	}
	if len(releases) != 0 {
		t.Fatalf("DiscoverV2Releases() = %v, want empty", releases)
	}
}

func TestLoadV2Index_releasesDeriveObjectiveMembership(t *testing.T) {
	root := t.TempDir()
	writeV2ReleaseFixture(t, root, "R001-first", "R001", "First release")
	writeV2ReleaseFixture(t, root, "R002-second", "R002", "Second release")

	writeV2ObjectiveWithRelease(t, root, "O001-first", "O001", "First objective", "R001")
	writeV2ObjectiveWithRelease(t, root, "O002-second", "O002", "Second objective", "R002")
	writeV2ObjectiveFixture(t, root, "O003-unassigned", "O003", "Unassigned objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-first.md", "T001", "First task", "O001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if len(index.Releases) != 2 {
		t.Fatalf("Releases = %d, want 2", len(index.Releases))
	}
	if got := index.ReleaseObjectives["R001"]; len(got) != 1 || got[0] != "O001" {
		t.Errorf("ReleaseObjectives[R001] = %v, want [O001]", got)
	}
	if got := index.ReleaseObjectives["R002"]; len(got) != 1 || got[0] != "O002" {
		t.Errorf("ReleaseObjectives[R002] = %v, want [O002]", got)
	}
	if got := index.ReleaseObjectives["R003"]; got != nil {
		t.Errorf("ReleaseObjectives[R003] = %v, want nil for unknown Release", got)
	}
	if index.Objectives["O003"].Release != "" {
		t.Errorf("unassigned Objective Release = %q, want empty", index.Objectives["O003"].Release)
	}
	if index.Tasks["T001"].Objective != "O001" {
		t.Errorf("Task Objective = %q, want O001", index.Tasks["T001"].Objective)
	}
}

func TestLoadV2Index_releaseIdentitySurvivesDirectorySlugChange(t *testing.T) {
	root := t.TempDir()
	writeV2ReleaseFixture(t, root, "R001-before", "R001", "Stable identity")

	first, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("first LoadV2Index() error = %v", err)
	}
	oldPath := filepath.Join(root, v2ReleasesDirName, "R001-before")
	newPath := filepath.Join(root, v2ReleasesDirName, "R001-after")
	if err := os.Rename(oldPath, newPath); err != nil {
		t.Fatalf("os.Rename() error = %v", err)
	}

	second, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("second LoadV2Index() error = %v", err)
	}
	if first.Releases["R001"].ID != second.Releases["R001"].ID {
		t.Errorf("Release identity changed after slug edit: %q -> %q", first.Releases["R001"].ID, second.Releases["R001"].ID)
	}
	if second.Releases["R001"].Source.Path != filepath.Join(v2ReleasesDirName, "R001-after", v2ReleaseFileName) {
		t.Errorf("Release Source.Path = %q, want moved path", second.Releases["R001"].Source.Path)
	}
}

func TestDiscoverV2Releases_rejectsDuplicateAndPathMismatch(t *testing.T) {
	t.Run("duplicate identity", func(t *testing.T) {
		root := t.TempDir()
		writeV2ReleaseFixture(t, root, "R001-first", "R001", "First")
		writeV2ReleaseFixture(t, root, "R001-second", "R001", "Second")
		_, err := DiscoverV2Releases(root)
		if !errors.Is(err, ErrV2DuplicateID) || !strings.Contains(err.Error(), "R001") {
			t.Fatalf("DiscoverV2Releases() error = %v, want duplicate R001 diagnostic", err)
		}
	})

	t.Run("directory identity mismatch", func(t *testing.T) {
		root := t.TempDir()
		writeV2ReleaseFixture(t, root, "R002-wrong", "R001", "Release")
		_, err := DiscoverV2Releases(root)
		if !errors.Is(err, ErrV2PathMismatch) || !strings.Contains(err.Error(), "R001") {
			t.Fatalf("DiscoverV2Releases() error = %v, want path/id diagnostic", err)
		}
	})
}

func TestDiscoverV2Releases_rejectsSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires elevated privileges on windows")
	}

	root := t.TempDir()
	outside := t.TempDir()
	writeV2ReleaseFixture(t, outside, "R001-first", "R001", "Escaped release")
	testutil.MkdirAll(t, filepath.Join(root, v2ReleasesDirName))
	if err := os.Symlink(filepath.Join(outside, v2ReleasesDirName, "R001-first"), filepath.Join(root, v2ReleasesDirName, "R001-first")); err != nil {
		t.Fatalf("os.Symlink() error = %v", err)
	}

	_, err := DiscoverV2Releases(root)
	if !errors.Is(err, ErrV2UnsafePath) || !strings.Contains(err.Error(), "R001-first") {
		t.Fatalf("DiscoverV2Releases() error = %v, want confined-path diagnostic", err)
	}
}

func TestLoadV2Index_danglingObjectiveReleaseNamesPathAndIDs(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveWithRelease(t, root, "O001-first", "O001", "First objective", "R999")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2MissingRelease) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2MissingRelease", err)
	}
	for _, want := range []string{"objectives/O001-first/Objective.md", "O001", "R999"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("LoadV2Index() error = %v, want %q named", err, want)
		}
	}
}

func TestLoadV2Index_rejectsLegacyPackagingTextWhenReleasesExist(t *testing.T) {
	root := t.TempDir()
	writeV2ReleaseFixture(t, root, "R001-first", "R001", "First release")
	writeV2ObjectiveWithRelease(t, root, "O001-first", "O001", "First objective", "v2")

	_, err := LoadV2Index(root)
	if !errors.Is(err, ErrV2InvalidReleaseReference) {
		t.Fatalf("LoadV2Index() error = %v, want ErrV2InvalidReleaseReference", err)
	}
	for _, want := range []string{"objectives/O001-first/Objective.md", "O001", "v2"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("LoadV2Index() error = %v, want %q named", err, want)
		}
	}
}

func TestLoadV2Index_noReleaseRecordsPreservesLegacyReleaseFreeProject(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O001-first", "O001", "First objective")
	writeV2TaskFixture(t, root, "O001-first", "T001-first.md", "T001", "First task", "O001")

	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	if len(index.Releases) != 0 || len(index.ReleaseObjectives) != 0 {
		t.Fatalf("release index = %v / %v, want both empty", index.Releases, index.ReleaseObjectives)
	}
}

func writeV2ObjectiveWithRelease(t *testing.T, root, dirName, id, title, release string) {
	t.Helper()
	content := "---\nid: " + id + "\ntitle: \"" + title + "\"\nstatus: planned\nrelease: " + release + "\n---\n\n# " + title + "\n"
	testutil.WriteFile(t, filepath.Join(root, v2ObjectivesDirName, dirName, v2ObjectiveFileName), content)
}
