package codehealth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

// Gap says why a proposal cannot be used as it stands. The empty Gap means the
// proposal is ready to confirm.
type Gap string

const (
	GapNone              Gap = ""
	GapMissingExecutable Gap = "missing_executable"
	GapUnsupportedStack  Gap = "unsupported_stack"
	GapNoReportYet       Gap = "no_report_yet"
	GapUnreadableFile    Gap = "unreadable_file"
	GapOversizedFile     Gap = "oversized_file"
	GapReadLimit         Gap = "read_limit_reached"
)

// Proposal is one confirmation-ready suggestion. Component is the repository
// relative directory it came from ("." for the root). A proposal with an empty
// Config.Capability is a component-level note that carries only a Gap and a
// Reason. GateFlag is the command the owner would add to their own gate for a
// report-only provider; Savepoint never edits their scripts.
type Proposal struct {
	Component string
	Config    CapabilityConfig
	Reason    string
	Gap       Gap
	GateFlag  string
}

// LookPath resolves an executable name, like exec.LookPath. Discover calls it
// and nothing else outside the project tree.
type LookPath func(file string) (string, error)

// Discover inspects a bounded set of files under root and returns proposals in
// a deterministic order. It never runs a process, installs, writes, or uses the
// network.
func Discover(ctx context.Context, root string, lookPath LookPath) ([]Proposal, error) {
	d := &discoverer{root: root, lookPath: lookPath}
	if err := d.init(); err != nil {
		return nil, err
	}
	comps, err := d.findComponents(ctx)
	if err != nil {
		return nil, err
	}
	var out []Proposal
	for _, c := range comps {
		out = append(out, d.propose(c)...)
	}
	if len(comps) == 0 {
		out = append(out, Proposal{Component: ".", Gap: GapUnsupportedStack,
			Reason: "no supported project found: expected go.mod, package.json, pyproject.toml, pytest.ini, or setup.cfg within " +
				fmt.Sprint(MaxDiscoveryDepth) + " levels"})
	}
	out = append(out, d.notes...)
	nameInstances(out)
	out = d.validated(out)
	sortProposals(out)
	return out, nil
}

type discoverer struct {
	root     string // as given
	realRoot string // symlinks resolved
	lookPath LookPath
	reads    int
	dirs     int
	notes    []Proposal
}

// component is one directory holding at least one manifest.
type component struct {
	dir    string // slash-separated, "." for the root
	stacks []stack
	files  map[string]bool // directory entries that are regular files
	dirSet map[string]bool // directory entries that are directories
}

func (d *discoverer) init() error {
	real, err := filepath.EvalSymlinks(d.root)
	if err != nil {
		return fmt.Errorf("discover: %w", err)
	}
	info, err := os.Stat(real)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("discover: %s is not a directory", d.root)
	}
	d.realRoot = real
	return nil
}

func (d *discoverer) abs(rel string) string {
	return filepath.Join(d.root, filepath.FromSlash(rel))
}

// skipDir reports directories the search never enters: hidden ones and the
// directory-style default exclusions.
func skipDir(name string) bool {
	if strings.HasPrefix(name, ".") {
		return true
	}
	for _, p := range defaultExclusions {
		if dir, ok := strings.CutSuffix(p, "/**"); ok && dir == name {
			return true
		}
	}
	return false
}

func (d *discoverer) findComponents(ctx context.Context) ([]component, error) {
	var comps []component
	var walk func(dir string, depth int) error
	walk = func(dir string, depth int) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if d.dirs++; d.dirs > MaxDiscoveryDirs {
			return nil
		}
		entries, err := os.ReadDir(d.abs(dir))
		if err != nil {
			d.note(dir, GapUnreadableFile, "cannot list directory: "+err.Error())
			return nil
		}
		c := component{dir: dir, files: map[string]bool{}, dirSet: map[string]bool{}}
		var subs []string
		for _, e := range entries {
			switch {
			case e.Type()&os.ModeSymlink != 0:
				// Symlinks are never followed during the search; a manifest
				// that is one is resolved and checked when it is read.
				if _, ok := manifestStacks[e.Name()]; ok {
					c.files[e.Name()] = true
				}
			case e.IsDir():
				c.dirSet[e.Name()] = true
				if depth < MaxDiscoveryDepth && !skipDir(e.Name()) {
					subs = append(subs, e.Name())
				}
			default:
				c.files[e.Name()] = true
			}
		}
		for _, name := range sortedKeys(c.files) {
			if s, ok := manifestStacks[name]; ok && !slices.Contains(c.stacks, s) {
				c.stacks = append(c.stacks, s)
			}
		}
		if len(c.stacks) > 0 {
			comps = append(comps, c)
		}
		slices.Sort(subs)
		for _, s := range subs {
			if err := walk(path.Join(dir, s), depth+1); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk(".", 0); err != nil {
		return nil, err
	}
	return comps, nil
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func (d *discoverer) note(dir string, gap Gap, reason string) {
	d.notes = append(d.notes, Proposal{Component: dir, Gap: gap, Reason: reason})
}

// readFile reads one project file within the bounds. A problem becomes a named
// gap note and a false result, never an error.
func (d *discoverer) readFile(dir, name string) ([]byte, bool) {
	rel := path.Join(dir, name)
	if d.reads >= MaxDiscoveryFiles {
		d.note(dir, GapReadLimit, fmt.Sprintf("%s not read: at most %d files are inspected", rel, MaxDiscoveryFiles))
		return nil, false
	}
	d.reads++
	resolved, err := filepath.EvalSymlinks(d.abs(rel))
	if err == nil {
		var within string
		within, err = filepath.Rel(d.realRoot, resolved)
		if err == nil && (within == ".." || strings.HasPrefix(within, ".."+string(filepath.Separator))) {
			d.note(dir, GapUnreadableFile, rel+" is a link that leaves the project and was not followed")
			return nil, false
		}
	}
	if err != nil {
		d.note(dir, GapUnreadableFile, rel+" could not be read: "+err.Error())
		return nil, false
	}
	f, err := os.Open(resolved)
	if err != nil {
		d.note(dir, GapUnreadableFile, rel+" could not be read: "+err.Error())
		return nil, false
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, MaxDiscoveryFileBytes+1))
	switch {
	case err != nil:
		d.note(dir, GapUnreadableFile, rel+" could not be read: "+err.Error())
		return nil, false
	case len(data) > MaxDiscoveryFileBytes:
		d.note(dir, GapOversizedFile, fmt.Sprintf("%s is larger than %d bytes and was not read", rel, MaxDiscoveryFileBytes))
		return nil, false
	}
	return data, true
}

// evidence is what a component's files say about its tooling.
type evidence struct {
	vitest, vitestV8, pytest, coveragePy bool
}

func (d *discoverer) inspect(c component) evidence {
	var ev evidence
	for _, cfg := range vitestConfigs {
		ev.vitest = ev.vitest || c.files[cfg]
	}
	ev.pytest = c.files["pytest.ini"]
	ev.coveragePy = c.files[".coveragerc"]
	if c.files["package.json"] {
		if data, ok := d.readFile(c.dir, "package.json"); ok {
			d.scanPackageJSON(c, data, &ev)
		}
	}
	if c.files["pyproject.toml"] {
		if data, ok := d.readFile(c.dir, "pyproject.toml"); ok {
			for _, line := range lines(data) {
				ev.pytest = ev.pytest || strings.Contains(line, "pytest")
				ev.coveragePy = ev.coveragePy || strings.Contains(line, "coverage")
			}
		}
	}
	if c.files["setup.cfg"] {
		if data, ok := d.readFile(c.dir, "setup.cfg"); ok {
			for _, line := range lines(data) {
				ev.pytest = ev.pytest || strings.HasPrefix(line, "[tool:pytest]")
				ev.coveragePy = ev.coveragePy || strings.HasPrefix(line, "[coverage:")
			}
		}
	}
	return ev
}

func (d *discoverer) scanPackageJSON(c component, data []byte, ev *evidence) {
	var pkg struct {
		Scripts         map[string]string `json:"scripts"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&pkg); err != nil {
		d.note(c.dir, GapUnreadableFile, path.Join(c.dir, "package.json")+" is not valid JSON")
		return
	}
	for _, deps := range []map[string]string{pkg.Dependencies, pkg.DevDependencies} {
		_, v := deps["vitest"]
		_, v8 := deps["@vitest/coverage-v8"]
		ev.vitest = ev.vitest || v
		ev.vitestV8 = ev.vitestV8 || v8
	}
	for _, s := range pkg.Scripts {
		ev.vitest = ev.vitest || strings.Contains(s, "vitest")
	}
}

// lines returns trimmed, lower-cased, comment-free lines of a small text file.
func lines(data []byte) []string {
	var out []string
	for _, l := range strings.Split(string(data), "\n") {
		l = strings.ToLower(strings.TrimSpace(l))
		if l != "" && !strings.HasPrefix(l, "#") && !strings.HasPrefix(l, ";") {
			out = append(out, l)
		}
	}
	return out
}

func (c component) hasStack(s stack) bool { return slices.Contains(c.stacks, s) }

func (c component) exists(name string) bool { return c.files[name] || c.dirSet[name] }

// slug names a component in file names and instance names.
func (c component) slug() string {
	if c.dir == "." {
		return "root"
	}
	return strings.ReplaceAll(c.dir, "/", "-")
}

// join makes a component-relative path repository-relative.
func (c component) join(rel string) string {
	if c.dir == "." {
		return rel
	}
	return path.Join(c.dir, rel)
}

func (c component) scope() []string {
	if c.dir == "." {
		return nil
	}
	return []string{c.dir + "/**"}
}

// exclusions returns the default exclusions that apply to the component, with
// directory patterns anchored to it.
func (c component) exclusions() []string {
	var out []string
	for _, p := range defaultExclusions {
		dir, isDir := strings.CutSuffix(p, "/**")
		switch {
		case !isDir:
			out = append(out, p)
		case c.dirSet[dir]:
			out = append(out, c.join(p))
		}
	}
	return out
}

func (d *discoverer) propose(c component) []Proposal {
	ev := d.inspect(c)
	var out []Proposal
	add := func(p Proposal) { p.Component = c.dir; out = append(out, p) }
	gapNote := func(reason string) { add(Proposal{Gap: GapUnsupportedStack, Reason: reason}) }

	if c.hasStack(stackGo) {
		add(d.reportProposal(c, CapabilityTests, ProviderGoTestJSON, "Go's own test runner emits a JSON event stream."))
		add(d.reportProposal(c, CapabilityCoverage, ProviderGoCoverProfile, "Go's own test runner writes a coverage profile."))
	}
	if c.hasStack(stackJS) {
		switch {
		case !ev.vitest:
			gapNote("no supported JavaScript test runner found; Vitest is the supported one")
		default:
			add(d.reportProposal(c, CapabilityTests, ProviderVitestJUnit, "Vitest is the project's test runner and can write JUnit XML."))
			if ev.vitestV8 {
				add(d.reportProposal(c, CapabilityCoverage, ProviderVitestV8, "@vitest/coverage-v8 is installed, so Vitest can write a JSON coverage report."))
			} else {
				add(Proposal{Config: CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderVitestV8},
					Gap: GapUnsupportedStack, Reason: "Vitest coverage needs @vitest/coverage-v8, which is not in package.json"})
			}
		}
	}
	if c.hasStack(stackPython) {
		if !ev.pytest {
			gapNote("no supported Python test runner found; pytest is the supported one")
		} else {
			add(d.reportProposal(c, CapabilityTests, ProviderPytestJUnit, "pytest is the project's test runner and can write JUnit XML."))
			if ev.coveragePy {
				add(d.reportProposal(c, CapabilityCoverage, ProviderCoveragePyJSON, "coverage.py is configured, so it can write a JSON coverage report."))
			} else {
				add(Proposal{Config: CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderCoveragePyJSON},
					Gap: GapUnsupportedStack, Reason: "no coverage.py configuration found"})
			}
		}
	}
	add(d.executedProposal(c, CapabilityComplexity, ProviderLizardCSV))
	add(d.executedProposal(c, CapabilityDuplication, ProviderJscpdJSON))
	osv := d.executedProposal(c, CapabilityDependencyVulnerability, ProviderOSVScannerJSON)
	if !c.hasLockfile() {
		osv.Gap = GapUnsupportedStack
		osv.Reason = "no supported lockfile found for OSV-Scanner; " + osv.Reason
	}
	add(osv)
	return out
}

func (c component) hasLockfile() bool {
	for _, s := range c.stacks {
		for _, l := range lockfiles[s] {
			if c.files[l] {
				return true
			}
		}
	}
	return false
}

func (d *discoverer) reportProposal(c component, cap Capability, p ProviderKey, why string) Proposal {
	ro := reportOnlyProviders[p]
	prop := Proposal{
		Config: CapabilityConfig{Capability: cap, Provider: p, Report: c.join(ro.Report),
			Scope: c.scope(), Exclusions: c.exclusions()},
		Reason:   why,
		GateFlag: ro.Gate,
	}
	if c.dir != "." {
		prop.GateFlag = "cd " + c.dir + " && " + ro.Gate
	}
	if _, err := os.Stat(d.abs(prop.Config.Report)); err == nil {
		prop.Reason += " A report already exists at " + prop.Config.Report + "."
	} else {
		prop.Gap = GapNoReportYet
		prop.Reason += " Add the command to your own gate so it writes " + prop.Config.Report + "."
	}
	return prop
}

func (d *discoverer) executedProposal(c component, cap Capability, p ProviderKey) Proposal {
	ex := executedProviders[p]
	report := ex.Report(c.slug())
	target := c.dir
	exclusions := c.exclusions()
	prop := Proposal{
		Config: CapabilityConfig{Capability: cap, Provider: p, Executable: ex.Executable,
			Args: ex.Args(target, report, exclusions), Report: report,
			Scope: c.scope(), Exclusions: exclusions},
		Reason: ex.Reason,
	}
	for _, name := range providerConfigs[p] {
		if c.files[name] {
			prop.Reason += " Existing " + name + " found."
		}
	}
	if _, err := d.lookPath(ex.Executable); err != nil {
		prop.Gap = GapMissingExecutable
		prop.Reason += " " + ex.Executable + " is not on PATH; install it yourself, Savepoint does not."
	}
	return prop
}

// nameInstances gives every instance that shares a capability and provider a
// name from its component, as the configuration requires.
func nameInstances(ps []Proposal) {
	type pair struct {
		c Capability
		p ProviderKey
	}
	groups := map[pair][]int{}
	for i, p := range ps {
		if p.Config.Capability != "" {
			k := pair{p.Config.Capability, p.Config.Provider}
			groups[k] = append(groups[k], i)
		}
	}
	for _, idx := range groups {
		if len(idx) < 2 {
			continue
		}
		for _, i := range idx {
			ps[i].Config.Name = component{dir: ps[i].Component}.slug()
		}
	}
}

// validated keeps the promise that every proposed configuration is valid: one
// that is not becomes a component-level note saying why.
func (d *discoverer) validated(ps []Proposal) []Proposal {
	out := ps[:0:0]
	for _, p := range ps {
		if p.Config.Capability != "" {
			cfg := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{p.Config}}
			if err := cfg.Validate(); err != nil {
				p = Proposal{Component: p.Component, Gap: GapUnsupportedStack,
					Reason: fmt.Sprintf("%s cannot be configured: %v", p.Config.Provider, err)}
			}
		}
		out = append(out, p)
	}
	return out
}

func sortProposals(ps []Proposal) {
	order := map[Capability]int{}
	for i, c := range Capabilities() {
		order[c] = i + 1
	}
	slices.SortStableFunc(ps, func(a, b Proposal) int {
		if c := strings.Compare(a.Component, b.Component); c != 0 {
			// The root sorts first, then paths alphabetically.
			switch {
			case a.Component == ".":
				return -1
			case b.Component == ".":
				return 1
			}
			return c
		}
		if c := order[a.Config.Capability] - order[b.Config.Capability]; c != 0 {
			return c
		}
		return strings.Compare(string(a.Config.Provider), string(b.Config.Provider))
	})
}
