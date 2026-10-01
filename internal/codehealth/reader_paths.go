package codehealth

import (
	"bufio"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Shared groundwork for the report readers: tools write paths in many forms,
// and every reader must hand back only safe repository-relative ones.

// RelPath maps a path from a tool report to a repository-relative slash path.
// It accepts absolute paths inside root and relative paths, with either slash
// style and a leading "./". A path outside root, one that traverses out of it,
// or one that evidence may not name (see validateRelativePath and
// sensitiveSegment) is rejected so the reader can drop it. The mapping is
// lexical: it never touches the file system.
func RelPath(root, p string) (string, bool) {
	p = strings.ReplaceAll(p, `\`, "/")
	root = strings.TrimRight(strings.ReplaceAll(root, `\`, "/"), "/")
	if isAbsSlash(p) {
		rest, ok := strings.CutPrefix(path.Clean(p), root+"/")
		if root == "" || !ok {
			return "", false
		}
		p = rest
	}
	for strings.HasPrefix(p, "./") {
		p = p[2:]
	}
	p = path.Clean(p)
	if p == "." || p == ".." || strings.HasPrefix(p, "../") {
		return "", false
	}
	if validateRelativePath(ErrUnsafeEvidence, "path", p) != nil {
		return "", false
	}
	if slicesAnyFunc(splitPath(p), sensitiveSegment) {
		return "", false
	}
	return p, true
}

func isAbsSlash(p string) bool {
	return strings.HasPrefix(p, "/") || (len(p) >= 3 && p[1] == ':' && p[2] == '/')
}

func slicesAnyFunc(s []string, f func(string) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

// inputScope is the instance scope as repository matching understands it.
func (in ReportInput) inputScope() InputScope {
	return InputScope{Include: in.Scope, Exclude: in.Exclusions}
}

// Path maps a report path to a repository-relative path inside the instance
// scope. A path outside the root or outside the scope is rejected.
func (in ReportInput) Path(p string) (string, bool) {
	rel, ok := RelPath(in.Root, p)
	if !ok || !in.inputScope().relevant(rel) {
		return "", false
	}
	return rel, true
}

// Bounds for finding Go modules, so a huge tree cannot make a read unbounded.
const (
	maxModuleDirs   = 20000
	maxGoModBytes   = 64 << 10
	goModFile       = "go.mod"
	goModuleKeyword = "module"
)

// GoModules maps Go module paths to the repository directory holding their
// go.mod, so a report's import paths become repository paths.
type GoModules struct {
	scope InputScope
	dirs  map[string]string // module path -> repository-relative directory
}

// LoadGoModules finds every go.mod at or under root, nested modules of a
// monorepo included. Hidden, vendor, node_modules, and testdata directories
// and links are not searched. A module whose go.mod cannot be read is skipped.
func LoadGoModules(root string, scope InputScope) GoModules {
	m := GoModules{scope: scope, dirs: map[string]string{}}
	seen := 0
	filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if p != root && (strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" || name == "testdata") {
				return fs.SkipDir
			}
			if seen++; seen > maxModuleDirs {
				return fs.SkipAll
			}
			return nil
		}
		if d.Name() != goModFile || !d.Type().IsRegular() {
			return nil
		}
		mod := readModulePath(p)
		rel, err := filepath.Rel(root, filepath.Dir(p))
		if mod == "" || err != nil {
			return nil
		}
		if _, dup := m.dirs[mod]; !dup {
			m.dirs[mod] = filepath.ToSlash(rel)
		}
		return nil
	})
	return m
}

// readModulePath returns the module path declared by a go.mod, or "".
func readModulePath(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(nil, maxGoModBytes)
	for sc.Scan() {
		line, _, _ := strings.Cut(sc.Text(), "//")
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == goModuleKeyword {
			return strings.Trim(fields[1], "\"`")
		}
	}
	return ""
}

// Resolve returns the repository-relative directory of a Go import path, using
// the module with the longest matching path. It reports false for an import
// path no module owns, such as the standard library or a dependency, and for a
// directory outside the instance scope; it never guesses.
func (m GoModules) Resolve(importPath string) (string, bool) {
	best, bestDir := "", ""
	for mod, dir := range m.dirs {
		if (importPath == mod || strings.HasPrefix(importPath, mod+"/")) && len(mod) > len(best) {
			best, bestDir = mod, dir
		}
	}
	if best == "" {
		return "", false
	}
	rel := path.Join(bestDir, strings.TrimPrefix(strings.TrimPrefix(importPath, best), "/"))
	if rel == "." {
		// The repository root has no path of its own; scope is judged by go.mod.
		return rel, m.scope.relevant(goModFile)
	}
	if _, ok := RelPath("/r", "/r/"+rel); !ok || !m.scope.relevant(rel) {
		return "", false
	}
	return rel, true
}

// RankedEvidence is an affected item with how bad it is; a larger Rank is worse.
type RankedEvidence struct {
	Ref  EvidenceRef
	Rank float64
}

// WorstEvidence keeps the worst MaxEvidence items in a deterministic order:
// worst first, ties broken by path, line, then note. Notes are reduced to one
// bounded line, and an item that evidence may not name is dropped.
func WorstEvidence(items []RankedEvidence) []EvidenceRef {
	kept := make([]RankedEvidence, 0, len(items))
	for _, it := range items {
		it.Ref.Note = sanitizeLine(it.Ref.Note)
		if it.Ref.validate("evidence") == nil {
			kept = append(kept, it)
		}
	}
	sort.SliceStable(kept, func(i, j int) bool {
		a, b := kept[i], kept[j]
		switch {
		case a.Rank != b.Rank:
			return a.Rank > b.Rank
		case a.Ref.Path != b.Ref.Path:
			return a.Ref.Path < b.Ref.Path
		case a.Ref.Line != b.Ref.Line:
			return a.Ref.Line < b.Ref.Line
		}
		return a.Ref.Note < b.Ref.Note
	})
	if len(kept) > MaxEvidence {
		kept = kept[:MaxEvidence]
	}
	out := make([]EvidenceRef, len(kept))
	for i, k := range kept {
		out[i] = k.Ref
	}
	return out
}
