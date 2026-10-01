package codehealth

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestRelPath(t *testing.T) {
	tests := []struct {
		name, root, in, want string
		ok                   bool
	}{
		{"relative", "/r/p", "a/b.go", "a/b.go", true},
		{"dot prefix", "/r/p", "./a/b.go", "a/b.go", true},
		{"backslashes", "/r/p", `a\b\c.go`, "a/b/c.go", true},
		{"absolute inside", "/r/p", "/r/p/a/b.go", "a/b.go", true},
		{"absolute with dots", "/r/p", "/r/p/x/../a.go", "a.go", true},
		{"windows absolute", `C:\work\p`, `C:\work\p\src\a.ts`, "src/a.ts", true},
		{"windows absolute forward", "C:/work/p/", "C:/work/p/src/a.ts", "src/a.ts", true},
		{"absolute outside", "/r/p", "/r/other/a.go", "", false},
		{"sibling prefix", "/r/p", "/r/p2/a.go", "", false},
		{"root itself", "/r/p", "/r/p", "", false},
		{"traversal", "/r/p", "../a.go", "", false},
		{"traversal after clean", "/r/p", "a/../../a.go", "", false},
		{"windows outside", `C:\work\p`, `D:\work\p\a.go`, "", false},
		{"url form", "/r/p", "https://x/a.go", "", false},
		{"sensitive", "/r/p", ".env", "", false},
		{"sensitive nested", "/r/p", "a/.git/config", "", false},
		{"empty", "/r/p", "", "", false},
		{"dot", "/r/p", ".", "", false},
		{"control", "/r/p", "a\n.go", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := RelPath(tt.root, tt.in)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("RelPath(%q, %q) = %q, %v; want %q, %v", tt.root, tt.in, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestReportInputPathHonoursScope(t *testing.T) {
	in := ReportInput{Root: "/r/p", Scope: []string{"src/**"}, Exclusions: []string{"src/gen/**"}}
	for p, want := range map[string]bool{
		"/r/p/src/a.go":     true,
		"src/gen/a.go":      false,
		"other/a.go":        false,
		"/elsewhere/a.go":   false,
		`src\sub\a.go`:      true,
		"/r/p/src/gen/a.go": false,
	} {
		if _, ok := in.Path(p); ok != want {
			t.Errorf("Path(%q) ok = %v, want %v", p, ok, want)
		}
	}
	if _, ok := (ReportInput{Root: "/r/p"}).Path("a.go"); !ok {
		t.Error("an empty scope should admit every path")
	}
}

func writeIn(t *testing.T, root, rel, content string) {
	t.Helper()
	writeFile(t, filepath.Join(root, filepath.FromSlash(rel)), content)
}

func TestGoModulesResolve(t *testing.T) {
	root := t.TempDir()
	writeIn(t, root, "go.mod", "module example.com/mono // root\n\ngo 1.22\n")
	writeIn(t, root, "svc/api/go.mod", "module \"example.com/mono/svc/api\"\n")
	writeIn(t, root, "vendor/x/go.mod", "module example.com/vendored\n")
	writeIn(t, root, "node_modules/y/go.mod", "module example.com/node\n")
	writeIn(t, root, "broken/go.mod", "go 1.22\n")

	m := LoadGoModules(root, InputScope{})
	tests := []struct {
		name, in, want string
		ok             bool
	}{
		{"root package", "example.com/mono", ".", true},
		{"package in root module", "example.com/mono/internal/x", "internal/x", true},
		{"nested module wins", "example.com/mono/svc/api/handlers", "svc/api/handlers", true},
		{"nested module itself", "example.com/mono/svc/api", "svc/api", true},
		{"unresolvable", "github.com/other/lib/pkg", "", false},
		{"standard library", "net/http", "", false},
		{"module prefix is not a path prefix", "example.com/monorepo/x", "", false},
		{"vendored module ignored", "example.com/vendored/p", "", false},
		{"module without a name ignored", "broken", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := m.Resolve(tt.in)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("Resolve(%q) = %q, %v; want %q, %v", tt.in, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestGoModulesResolveKeepsToScope(t *testing.T) {
	root := t.TempDir()
	writeIn(t, root, "svc/api/go.mod", "module example.com/api\n")
	writeIn(t, root, "svc/web/go.mod", "module example.com/web\n")

	m := LoadGoModules(root, InputScope{Include: []string{"svc/api/**"}})
	if got, ok := m.Resolve("example.com/api/h"); !ok || got != "svc/api/h" {
		t.Errorf("in-scope module = %q, %v", got, ok)
	}
	if got, ok := m.Resolve("example.com/web/h"); ok {
		t.Errorf("out-of-scope module resolved to %q", got)
	}
}

func TestWorstEvidence(t *testing.T) {
	var items []RankedEvidence
	for i := 0; i < MaxEvidence+5; i++ {
		items = append(items, RankedEvidence{Ref: EvidenceRef{Path: fmt.Sprintf("f%02d.go", i)}, Rank: float64(i)})
	}
	got := WorstEvidence(items)
	if len(got) != MaxEvidence {
		t.Fatalf("kept %d, want %d", len(got), MaxEvidence)
	}
	if got[0].Path != "f24.go" || got[MaxEvidence-1].Path != "f05.go" {
		t.Errorf("not the worst first: %v ... %v", got[0], got[MaxEvidence-1])
	}

	tie := WorstEvidence([]RankedEvidence{
		{Ref: EvidenceRef{Path: "b.go", Line: 1}, Rank: 1},
		{Ref: EvidenceRef{Path: "a.go", Line: 9}, Rank: 1},
		{Ref: EvidenceRef{Path: "a.go", Line: 2}, Rank: 1},
	})
	if tie[0] != (EvidenceRef{Path: "a.go", Line: 2}) || tie[1].Line != 9 || tie[2].Path != "b.go" {
		t.Errorf("ties are not ordered by path then line: %v", tie)
	}
}

func TestWorstEvidenceBoundsNotesAndDropsUnsafeItems(t *testing.T) {
	got := WorstEvidence([]RankedEvidence{
		{Ref: EvidenceRef{Path: "a.go", Note: "line one\nline two " + strings.Repeat("x", 500)}, Rank: 2},
		{Ref: EvidenceRef{Path: "../escape.go"}, Rank: 9},
		{Ref: EvidenceRef{Path: ".env"}, Rank: 8},
		{Ref: EvidenceRef{Path: "neg.go", Line: -1}, Rank: 7},
	})
	if len(got) != 1 || got[0].Path != "a.go" {
		t.Fatalf("got %v, want only a.go", got)
	}
	if len(got[0].Note) > MaxNoteLen || strings.ContainsAny(got[0].Note, "\n") {
		t.Errorf("note not bounded to one line: %q", got[0].Note)
	}
	if err := got[0].validate("e"); err != nil {
		t.Errorf("kept item fails validation: %v", err)
	}
	if WorstEvidence(nil) == nil {
		t.Error("an empty result should be an empty slice, not nil")
	}
}

func TestCollectHandsReadersTheRoot(t *testing.T) {
	dir := project(t)
	var gotRoot string
	rd := readerFunc(func(_ context.Context, in ReportInput) (Reading, error) {
		gotRoot = in.Root
		return goodReading(1, UnitCCN), nil
	})
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){"lizard": stdout("x")}}
	collect(t, dir, cfgOf(lizardInstance("a", "lizard")), Readers{ProviderLizardCSV: rd}, tools)
	if gotRoot != dir {
		t.Fatalf("reader saw root %q, want %q", gotRoot, dir)
	}
}
