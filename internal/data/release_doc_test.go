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
			name: "both present",
			write: map[string]string{
				"PRD.md":    "# Product\nvision",
				"Design.md": "# Design\narchitecture",
			},
			wantBody: map[ReleaseDocID]string{
				ReleaseDocPRD:    "# Product\nvision",
				ReleaseDocDesign: "# Design\narchitecture",
			},
			wantAvail: map[ReleaseDocID]bool{ReleaseDocPRD: true, ReleaseDocDesign: true},
		},
		{
			name:      "both missing",
			write:     map[string]string{},
			wantAvail: map[ReleaseDocID]bool{ReleaseDocPRD: false, ReleaseDocDesign: false},
		},
		{
			name:  "prd present design missing",
			write: map[string]string{"PRD.md": "only prd"},
			wantBody: map[ReleaseDocID]string{
				ReleaseDocPRD: "only prd",
			},
			wantAvail: map[ReleaseDocID]bool{ReleaseDocPRD: true, ReleaseDocDesign: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			for name, content := range tt.write {
				testutil.WriteFile(t, filepath.Join(root, name), content)
			}

			docs, err := LoadReleaseDocs(root)
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
				if !wantAvail && doc.Body != "" {
					t.Errorf("doc %q unavailable but Body = %q, want empty", id, doc.Body)
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

func TestLoadReleaseDocsLabelsAndPaths(t *testing.T) {
	docs, err := LoadReleaseDocs(t.TempDir())
	if err != nil {
		t.Fatalf("LoadReleaseDocs() error = %v", err)
	}

	want := map[ReleaseDocID]struct{ label, path string }{
		ReleaseDocPRD:    {label: "PRD", path: "PRD.md"},
		ReleaseDocDesign: {label: "Design", path: "Design.md"},
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

	_, err := LoadReleaseDocs(root)
	if err == nil {
		t.Fatal("LoadReleaseDocs() error = nil, want read error")
	}
	if !strings.Contains(err.Error(), "PRD.md") {
		t.Errorf("LoadReleaseDocs() error = %v, want path context containing PRD.md", err)
	}
}
