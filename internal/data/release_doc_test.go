package data

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/testutil"
)

func docByID(docs []ReleaseDoc, id ReleaseDocID) (ReleaseDoc, bool) {
	for _, doc := range docs {
		if doc.ID == id {
			return doc, true
		}
	}
	return ReleaseDoc{}, false
}

func TestLoadReleaseDocs(t *testing.T) {
	tests := []struct {
		name      string
		write     map[string]string
		wantBody  map[ReleaseDocID]string
		wantAvail map[ReleaseDocID]bool
	}{
		{
			name: "all present",
			write: map[string]string{
				"releases/v1.2/v1.2-PRD.md": "# Release PRD\nv1.2 scope",
				"PRD.md":                    "# Overall PRD\nproduct vision",
				"Design.md":                 "# Overall Design\narchitecture",
			},
			wantBody: map[ReleaseDocID]string{
				ReleaseDocReleasePRD:    "# Release PRD\nv1.2 scope",
				ReleaseDocOverallPRD:    "# Overall PRD\nproduct vision",
				ReleaseDocOverallDesign: "# Overall Design\narchitecture",
			},
			wantAvail: map[ReleaseDocID]bool{
				ReleaseDocReleasePRD:    true,
				ReleaseDocOverallPRD:    true,
				ReleaseDocOverallDesign: true,
			},
		},
		{
			name:  "release prd missing is not fatal",
			write: map[string]string{"PRD.md": "# Overall PRD", "Design.md": "# Overall Design"},
			wantAvail: map[ReleaseDocID]bool{
				ReleaseDocReleasePRD:    false,
				ReleaseDocOverallPRD:    true,
				ReleaseDocOverallDesign: true,
			},
		},
		{
			name:  "only release prd present",
			write: map[string]string{"releases/v1.2/v1.2-PRD.md": "# Release PRD"},
			wantBody: map[ReleaseDocID]string{
				ReleaseDocReleasePRD: "# Release PRD",
			},
			wantAvail: map[ReleaseDocID]bool{
				ReleaseDocReleasePRD:    true,
				ReleaseDocOverallPRD:    false,
				ReleaseDocOverallDesign: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			for name, content := range tt.write {
				testutil.WriteFile(t, filepath.Join(root, filepath.FromSlash(name)), content)
			}

			docs, err := LoadReleaseDocs(root, "v1.2")
			if err != nil {
				t.Fatalf("LoadReleaseDocs() error = %v", err)
			}
			if len(docs) != len(releaseDocSpecs) {
				t.Fatalf("LoadReleaseDocs() returned %d docs, want %d", len(docs), len(releaseDocSpecs))
			}

			for id, wantAvail := range tt.wantAvail {
				doc, ok := docByID(docs, id)
				if !ok {
					t.Fatalf("LoadReleaseDocs() missing doc %q", id)
				}
				if doc.Available != wantAvail {
					t.Errorf("doc %q Available = %v, want %v", id, doc.Available, wantAvail)
				}
			}
			for id, wantBody := range tt.wantBody {
				doc, _ := docByID(docs, id)
				if doc.Body != wantBody {
					t.Errorf("doc %q Body = %q, want %q", id, doc.Body, wantBody)
				}
			}
		})
	}
}

// TestLoadReleaseDocsReleaseScoped proves the Release PRD tracks the requested
// release: v1.2's PRD must not appear when loading v1.3.
func TestLoadReleaseDocsReleaseScoped(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(t, filepath.Join(root, "releases", "v1.2", "v1.2-PRD.md"), "v1.2 prd")

	docs, err := LoadReleaseDocs(root, "v1.3")
	if err != nil {
		t.Fatalf("LoadReleaseDocs() error = %v", err)
	}
	prd, _ := docByID(docs, ReleaseDocReleasePRD)
	if prd.Available {
		t.Errorf("Release PRD Available = true for v1.3, want false (only v1.2 exists)")
	}
}

func TestLoadReleaseDocsLabelsAndPaths(t *testing.T) {
	docs, err := LoadReleaseDocs(t.TempDir(), "v1.2")
	if err != nil {
		t.Fatalf("LoadReleaseDocs() error = %v", err)
	}

	want := map[ReleaseDocID]struct{ label, path string }{
		ReleaseDocReleasePRD:    {label: "Release PRD", path: filepath.FromSlash("releases/v1.2/v1.2-PRD.md")},
		ReleaseDocOverallPRD:    {label: "Overall PRD", path: "PRD.md"},
		ReleaseDocOverallDesign: {label: "Overall Design", path: "Design.md"},
	}
	for id, w := range want {
		doc, ok := docByID(docs, id)
		if !ok {
			t.Fatalf("LoadReleaseDocs() missing doc %q", id)
		}
		if doc.Label != w.label {
			t.Errorf("doc %q Label = %q, want %q", id, doc.Label, w.label)
		}
		if doc.Path != w.path {
			t.Errorf("doc %q Path = %q, want %q", id, doc.Path, w.path)
		}
	}
}

func TestLoadReleaseDocsReadError(t *testing.T) {
	root := t.TempDir()
	// A directory in place of the expected file triggers a read error that is
	// not os.IsNotExist, exercising the fatal-error path portably.
	testutil.MkdirAll(t, filepath.Join(root, "PRD.md"))

	_, err := LoadReleaseDocs(root, "v1.2")
	if err == nil {
		t.Fatal("LoadReleaseDocs() error = nil, want read error")
	}
	if !strings.Contains(err.Error(), "PRD.md") {
		t.Errorf("LoadReleaseDocs() error = %v, want path context containing PRD.md", err)
	}
}
