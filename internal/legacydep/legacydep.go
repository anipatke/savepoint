// Package legacydep finds an old Savepoint copy that a project's own
// package.json would run in place of the current release. `npx savepoint`
// prefers the project's installed copy, so a leftover 1.x dependency makes a
// migrated V2 project look broken. Detection is read-only: it never edits the
// user's package files.
package legacydep

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const packageName = "savepoint"

// sections are the package.json fields npm installs from.
var sections = []string{"dependencies", "devDependencies", "optionalDependencies", "peerDependencies"}

// Finding is one old Savepoint dependency in a project.
type Finding struct {
	File      string // package.json path
	Section   string // package.json field that declares it, or empty if only installed
	Range     string // declared range, or empty
	Installed string // version in node_modules, or empty
}

// Message says what was found and why it matters.
func (f Finding) Message() string {
	var parts []string
	if f.Range != "" {
		parts = append(parts, fmt.Sprintf("package.json pins savepoint %s in %s", f.Range, f.Section))
	}
	if f.Installed != "" {
		parts = append(parts, fmt.Sprintf("node_modules holds savepoint %s", f.Installed))
	}
	return strings.Join(parts, "; ") + ". `npx savepoint` runs that older copy, not this release, so a V2 project can fail with errors from the old version."
}

// Repair says how to fix it.
func (f Finding) Repair() string {
	return "Run `npm install -D savepoint@latest` (savepoint is a dev tool, so devDependencies is the right place), or remove the dependency and use `npx savepoint@latest`."
}

// Detect looks for an old Savepoint dependency under projectDir. It reports
// one when the declared range admits only versions below 2, or when the
// installed copy is below 2. No package.json, an unreadable one, or a current
// dependency is not a finding.
func Detect(projectDir string) (Finding, bool) {
	file := filepath.Join(projectDir, "package.json")
	raw, err := os.ReadFile(file)
	if err != nil {
		return Finding{}, false
	}
	var pkg map[string]json.RawMessage
	if json.Unmarshal(raw, &pkg) != nil {
		return Finding{}, false
	}

	f := Finding{File: file}
	for _, section := range sections {
		var deps map[string]string
		if json.Unmarshal(pkg[section], &deps) != nil {
			continue
		}
		if r, ok := deps[packageName]; ok && rangeIsLegacy(r) {
			f.Section, f.Range = section, r
			break
		}
	}
	if v := installedVersion(projectDir); v != "" && majorOf(v) >= 0 && majorOf(v) < 2 {
		f.Installed = v
	}
	if f.Range == "" && f.Installed == "" {
		return Finding{}, false
	}
	return f, true
}

func installedVersion(projectDir string) string {
	raw, err := os.ReadFile(filepath.Join(projectDir, "node_modules", packageName, "package.json"))
	if err != nil {
		return ""
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if json.Unmarshal(raw, &pkg) != nil {
		return ""
	}
	return pkg.Version
}

var leadingNumber = regexp.MustCompile(`\d+`)

// rangeIsLegacy reports whether every way the range can be satisfied is below
// major 2. A range with an alternative that admits 2.x or later (">=1",
// "^1 || ^2") is not legacy, and one that names no version ("latest", "*",
// a path or git URL) is left alone.
func rangeIsLegacy(r string) bool {
	alternatives := strings.Split(r, "||")
	legacy := false
	for _, alt := range alternatives {
		alt = strings.TrimSpace(alt)
		if strings.HasPrefix(alt, ">") {
			return false // open-ended upward: admits 2.x
		}
		m := leadingNumber.FindString(alt)
		if m == "" {
			return false
		}
		major := majorOf(m)
		if major >= 2 {
			return false
		}
		legacy = true
	}
	return legacy
}

// majorOf returns the leading major version number of v, or -1.
func majorOf(v string) int {
	m := leadingNumber.FindString(v)
	if m == "" {
		return -1
	}
	n := 0
	for _, c := range m {
		n = n*10 + int(c-'0')
	}
	return n
}
