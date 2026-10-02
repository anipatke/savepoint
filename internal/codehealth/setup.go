package codehealth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Probe answers the two questions reconciliation asks about a configured
// instance: can its executable be found, and does a project file exist. Both are
// read-only.
type Probe struct {
	LookPath LookPath
	Exists   func(rel string) bool
}

// ProjectProbe checks files under root and executables with lookPath.
func ProjectProbe(root string, lookPath LookPath) Probe {
	return Probe{
		LookPath: lookPath,
		Exists: func(rel string) bool {
			_, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
			return err == nil
		},
	}
}

// Attention is a configured instance whose executable or inputs are gone. Setup
// reports it and never removes it.
type Attention struct {
	Config  CapabilityConfig
	Problem string
}

// SetupPlan is what setup found and what applying it would change. Configured
// entries are never edited: Unchanged and Attention both keep them as they are.
type SetupPlan struct {
	HasConfig    bool
	New          []Proposal
	Unchanged    []CapabilityConfig
	Attention    []Attention
	NotSuggested []Proposal

	merged  Config
	invalid error
}

// Plan reconciles an existing configuration (nil when there is none) with
// discovery's proposals. A proposal that matches a configured instance by key,
// or by capability, provider, and scope, is already covered. A proposal that
// cannot be used as it stands is listed but never added.
func Plan(existing *Config, proposals []Proposal, probe Probe) SetupPlan {
	plan := SetupPlan{merged: Config{Version: ConfigVersion}}
	if existing != nil {
		plan.HasConfig = true
		plan.merged.Capabilities = slices.Clone(existing.Capabilities)
	}
	for _, cc := range plan.merged.Capabilities {
		if problem := probe.problem(cc); problem != "" {
			plan.Attention = append(plan.Attention, Attention{Config: cc, Problem: problem})
		} else {
			plan.Unchanged = append(plan.Unchanged, cc)
		}
	}
	configured := slices.Clone(plan.merged.Capabilities)
	for _, p := range proposals {
		switch {
		case p.Config.Capability == "" || p.Gap == GapUnsupportedStack:
			plan.NotSuggested = append(plan.NotSuggested, p)
		case slices.ContainsFunc(configured, func(cc CapabilityConfig) bool { return covers(cc, p.Config) }):
		default:
			plan.New = append(plan.New, p)
			plan.merged.Capabilities = append(plan.merged.Capabilities, p.Config)
		}
	}
	if len(plan.New) > 0 {
		plan.invalid = plan.merged.Validate()
	}
	return plan
}

func covers(configured, proposed CapabilityConfig) bool {
	if configured.Capability != proposed.Capability || configured.Provider != proposed.Provider {
		return false
	}
	if configured.Name != "" && configured.Name == proposed.Name {
		return true
	}
	return sameScope(configured.Scope, proposed.Scope)
}

func (pr Probe) problem(cc CapabilityConfig) string {
	switch {
	case cc.Executable != "":
		if _, err := pr.LookPath(cc.Executable); err != nil {
			return cc.Executable + " is not on PATH"
		}
	case cc.Report != "" && !pr.Exists(cc.Report):
		return "no report found at " + cc.Report + " yet; your own gate writes it"
	}
	for _, pattern := range cc.Scope {
		if dir := scopeDir(pattern); dir != "" && !pr.Exists(dir) {
			return "scope " + dir + " no longer exists"
		}
	}
	return ""
}

// scopeDir is the literal directory a scope pattern starts with, if any.
func scopeDir(pattern string) string {
	prefix := pattern
	if i := strings.IndexAny(pattern, "*?["); i >= 0 {
		prefix = pattern[:i]
		if j := strings.LastIndex(prefix, "/"); j >= 0 {
			prefix = prefix[:j]
		} else {
			prefix = ""
		}
	}
	return strings.TrimSuffix(prefix, "/")
}

// PlanProject discovers proposals under root, loads any existing configuration,
// and reconciles them. It writes nothing.
func PlanProject(ctx context.Context, root string, lookPath LookPath) (SetupPlan, error) {
	proposals, err := Discover(ctx, root, lookPath)
	if err != nil {
		return SetupPlan{}, err
	}
	var existing *Config
	cfg, err := NewStore(root).LoadConfig()
	switch {
	case err == nil:
		existing = &cfg
	case errors.Is(err, ErrConfigNotFound):
	default:
		return SetupPlan{}, fmt.Errorf("existing health configuration cannot be used; fix or remove it first: %w", err)
	}
	return Plan(existing, proposals, ProjectProbe(root, lookPath)), nil
}

// Apply saves the configuration with the new proposals added. It reports
// whether the file changed. With nothing new it writes nothing, so a configured
// file keeps its bytes.
func (p SetupPlan) Apply(store Store) (changed bool, err error) {
	if len(p.New) == 0 {
		return false, nil
	}
	if p.invalid != nil {
		return false, fmt.Errorf("the suggested tools conflict with your configuration, so nothing was saved: %w", p.invalid)
	}
	return store.SaveConfig(p.merged)
}

// ConfigPath is the project-relative path setup writes.
const ConfigPath = savepointDir + "/" + healthDir + "/" + configFile

// Preview renders the plan for a person who has not applied it yet.
func (p SetupPlan) Preview() string {
	var b strings.Builder
	p.writeBody(&b)
	switch {
	case len(p.New) > 0:
		fmt.Fprintf(&b, "\nNothing was written. Run `savepoint health setup --apply` to add the new tools to %s.\n", ConfigPath)
	default:
		b.WriteString("\nNo changes. Nothing was written.\n")
	}
	return b.String()
}

// Applied renders the plan after Apply.
func (p SetupPlan) Applied(changed bool) string {
	var b strings.Builder
	p.writeBody(&b)
	if changed {
		fmt.Fprintf(&b, "\nSaved %d new health tool(s) to %s.\n", len(p.New), ConfigPath)
	} else {
		b.WriteString("\nNo changes.\n")
	}
	return b.String()
}

func label(cc CapabilityConfig, component string) string {
	s := fmt.Sprintf("%s via %s", cc.Capability, cc.Provider)
	if cc.Name != "" {
		s += " (" + cc.Name + ")"
	}
	if component != "" && component != "." {
		s += " in " + component
	}
	return s
}

func (p SetupPlan) writeBody(b *strings.Builder) {
	b.WriteString("Health tools\n")
	if len(p.New) > 0 {
		b.WriteString("\nNew:\n")
		for _, pr := range p.New {
			fmt.Fprintf(b, "  + %s\n      %s\n", label(pr.Config, pr.Component), pr.Reason)
			if pr.Gap != GapNone {
				fmt.Fprintf(b, "      gap: %s\n", strings.ReplaceAll(string(pr.Gap), "_", " "))
			}
			if pr.GateFlag != "" {
				fmt.Fprintf(b, "      add to your own gate: %s\n", pr.GateFlag)
			}
		}
	}
	if len(p.Unchanged) > 0 {
		b.WriteString("\nUnchanged:\n")
		for _, cc := range p.Unchanged {
			fmt.Fprintf(b, "  = %s\n", label(cc, ""))
		}
	}
	if len(p.Attention) > 0 {
		b.WriteString("\nNeeds attention (kept as configured):\n")
		for _, a := range p.Attention {
			fmt.Fprintf(b, "  ! %s: %s\n", label(a.Config, ""), a.Problem)
		}
	}
	if len(p.NotSuggested) > 0 {
		b.WriteString("\nNot suggested:\n")
		for _, pr := range p.NotSuggested {
			what := pr.Component
			if pr.Config.Capability != "" {
				what = label(pr.Config, pr.Component)
			}
			fmt.Fprintf(b, "  - %s: %s\n", what, pr.Reason)
		}
	}
	if ex := p.exclusions(); len(ex) > 0 {
		fmt.Fprintf(b, "\nDefault exclusions for new tools: %s\n", strings.Join(ex, ", "))
	}
	if p.usesOSV() {
		b.WriteString("\nOSV-Scanner contacts OSV.dev: the scanner, not Savepoint, sends your package names and versions there.\n")
	}
	if p.invalid != nil {
		fmt.Fprintf(b, "\nThe new tools conflict with your configuration and cannot be added: %v\n", p.invalid)
	}
}

func (p SetupPlan) exclusions() []string {
	var out []string
	for _, pr := range p.New {
		for _, e := range pr.Config.Exclusions {
			if !slices.Contains(out, e) {
				out = append(out, e)
			}
		}
	}
	slices.Sort(out)
	return out
}

func (p SetupPlan) usesOSV() bool {
	isOSV := func(cc CapabilityConfig) bool { return cc.Provider == ProviderOSVScannerJSON }
	return slices.ContainsFunc(p.Unchanged, isOSV) ||
		slices.ContainsFunc(p.Attention, func(a Attention) bool { return isOSV(a.Config) }) ||
		slices.ContainsFunc(p.New, func(pr Proposal) bool { return isOSV(pr.Config) })
}
