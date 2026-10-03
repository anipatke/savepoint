package codehealth

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

// SharedReport is one report file that more than one configured instance
// reads. Savepoint reads the file for each, so one run's output is counted as
// every instance's result, and two runners writing it overwrite each other.
type SharedReport struct {
	Path      string
	Instances []CapabilityConfig
}

// SharedReports lists report files read by more than one instance, in the order
// each is first configured. Setup warns about them and never edits a configured
// entry; doctor reports them too.
func SharedReports(cfg Config) []SharedReport {
	var out []SharedReport
	at := map[string]int{}
	for _, cc := range cfg.Capabilities {
		if cc.Report == "" {
			continue
		}
		key := path.Clean(cc.Report)
		if i, ok := at[key]; ok {
			out[i].Instances = append(out[i].Instances, cc)
			continue
		}
		at[key] = len(out)
		out = append(out, SharedReport{Path: cc.Report, Instances: []CapabilityConfig{cc}})
	}
	return slices.DeleteFunc(out, func(s SharedReport) bool { return len(s.Instances) < 2 })
}

// Warning says what is wrong and how to fix it.
func (s SharedReport) Warning() string {
	names := make([]string, len(s.Instances))
	for i, cc := range s.Instances {
		names[i] = label(cc, "")
	}
	each := "both read"
	if len(names) > 2 {
		each = "all read"
	}
	return fmt.Sprintf("%s %s %s, so each would be judged on the other's results. Give each its own report file (for example %s) in %s.",
		strings.Join(names, " and "), each, s.Path, suggestSeparate(s), ConfigPath)
}

func suggestSeparate(s SharedReport) string {
	last := s.Instances[len(s.Instances)-1]
	runner := strings.SplitN(string(last.Provider), "-", 2)[0]
	return path.Join(path.Dir(s.Path), runner+"-"+path.Base(s.Path))
}

// reportCommandFor adapts a conventional report command to the path actually
// configured, so a renamed report is not suggested with the default name.
func reportCommandFor(gate string, provider ProviderKey, report string) string {
	conventional := reportOnlyProviders[provider].Report
	if conventional == "" || report == "" || path.Base(conventional) == path.Base(report) {
		return gate
	}
	return strings.ReplaceAll(gate, path.Base(conventional), path.Base(report))
}

// unignoredReports lists the report files or top-level directories that
// configured report-only providers write and the project's root .gitignore
// does not cover. It reads only that one file and does nothing when it cannot.
func unignoredReports(root string, cfg Config) []string {
	patterns := gitignorePatterns(filepath.Join(root, ".gitignore"))
	var out []string
	for _, cc := range cfg.Capabilities {
		if cc.Report == "" || cc.Executable != "" || strings.HasPrefix(cc.Report, savepointDir+"/") {
			continue
		}
		target := path.Clean(cc.Report)
		if ignored(patterns, target) {
			continue
		}
		// A report inside a directory is ignored as the directory.
		if dir, _, nested := strings.Cut(target, "/"); nested {
			if ignored(patterns, dir) {
				continue
			}
			target = dir + "/"
		}
		if !slices.Contains(out, target) {
			out = append(out, target)
		}
	}
	return out
}

func gitignorePatterns(file string) []string {
	f, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}
		out = append(out, strings.TrimSuffix(strings.TrimPrefix(line, "/"), "/"))
	}
	return out
}

func ignored(patterns []string, target string) bool {
	for _, pat := range patterns {
		if ok, _ := path.Match(pat, target); ok {
			return true
		}
		if ok, _ := path.Match(pat, path.Base(target)); ok && !strings.Contains(pat, "/") {
			return true
		}
	}
	return false
}
