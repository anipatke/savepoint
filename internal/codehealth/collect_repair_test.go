package codehealth

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCollectReadsTheReportASuggestedToolWrites(t *testing.T) {
	for _, provider := range []ProviderKey{ProviderLizardCSV, ProviderJscpdJSON, ProviderOSVScannerJSON} {
		t.Run(string(provider), func(t *testing.T) {
			root := project(t)
			write(t, root, "go.mod", "module example\n")
			proposals, err := Discover(context.Background(), root, func(s string) (string, error) { return s, nil })
			if err != nil {
				t.Fatal(err)
			}
			var cc CapabilityConfig
			for _, p := range proposals {
				if p.Config.Provider == provider {
					cc = p.Config
				}
			}
			if cc.Report == "" {
				t.Fatalf("no proposal for %s", provider)
			}
			reportPath := filepath.Join(root, filepath.FromSlash(cc.Report))
			var got string
			readers := Readers{provider: readerFunc(func(_ context.Context, in ReportInput) (Reading, error) {
				got = string(in.Data)
				return goodReading(1, capabilityUnits[cc.Capability]), nil
			})}
			writes := func(body string, code int) *fakeTools {
				return &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){
					cc.Executable: func(context.Context, ToolSpec) (ToolResult, error) {
						if body != "" {
							if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
								return ToolResult{}, err
							}
							if err := os.WriteFile(reportPath, []byte(body), 0o644); err != nil {
								return ToolResult{}, err
							}
						}
						return ToolResult{Stdout: []byte("console-status"), ExitCode: code}, nil
					},
				}}
			}

			if r := collect(t, root, cfgOf(cc), readers, writes("actual-report", 0)).Results[0].Result; r.Outcome != OutcomeAvailable || got != "actual-report" {
				t.Fatalf("outcome %s, reader got %q", r.Outcome, got)
			}
			if _, err := os.Stat(reportPath); !os.IsNotExist(err) {
				t.Errorf("report was left behind: %v", err)
			}

			// An old report must not stand in for a run that wrote none.
			write(t, root, cc.Report, "old-report")
			r := collect(t, root, cfgOf(cc), readers, writes("", 0)).Results[0].Result
			if r.Outcome != OutcomeFailed || !strings.Contains(r.Reason, "wrote no report") {
				t.Errorf("silent run: %s %q", r.Outcome, r.Reason)
			}

			// A failing run does not leave its partial report either.
			if r := collect(t, root, cfgOf(cc), readers, writes("partial", 2)).Results[0].Result; r.Outcome != OutcomeFailed {
				t.Errorf("exit 2: %s", r.Outcome)
			}
			if _, err := os.Stat(reportPath); !os.IsNotExist(err) {
				t.Errorf("failed run left a report: %v", err)
			}
		})
	}
}

func TestCollectKeepsHistoryAcrossARename(t *testing.T) {
	steps := []struct {
		name        string
		first, then CapabilityConfig
		kept        bool
	}{
		{"unnamed to named", lizardInstance("", "fake", "api/**"), lizardInstance("backend", "fake", "api/**"), true},
		{"named to named", lizardInstance("api", "fake", "api/**"), lizardInstance("backend", "fake", "api/**"), true},
		{"changed scope", lizardInstance("api", "fake", "api/**"), lizardInstance("api", "fake", "web/**"), false},
	}
	for _, s := range steps {
		t.Run(s.name, func(t *testing.T) {
			root := project(t)
			tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){"fake": stdout("report")}}
			run := func(cc CapabilityConfig, n int) string {
				got, err := Collect(context.Background(), CollectRequest{
					Root: root, Origin: OriginOfficial, Config: cfgOf(cc), Readers: Readers{cc.Provider: okReader(5, UnitCCN)},
					Runner: tools, Clock: func() time.Time { return testClock().Add(time.Duration(n) * time.Second) },
				})
				if err != nil {
					t.Fatal(err)
				}
				return collectedFor(t, got, CapabilityComplexity, cc.Name).Result.CollectedAt
			}
			for i := 0; i < 3; i++ {
				run(s.first, i)
			}
			run(s.then, 4)
			snaps, err := NewStore(root).LoadSnapshots()
			if err != nil {
				t.Fatal(err)
			}
			last := snaps[len(snaps)-1]
			for _, c := range last.Summary.Capabilities {
				if c.Capability != CapabilityComplexity {
					continue
				}
				if enough := !strings.Contains(c.Explanation, "Not enough comparable"); enough != s.kept {
					t.Errorf("history kept = %v, want %v: %s", enough, s.kept, c.Explanation)
				}
			}
		})
	}
}

func TestCollectKeepsTwoScopedInstancesSeparate(t *testing.T) {
	root := project(t)
	api, web := lizardInstance("api", "fake", "api/**"), lizardInstance("web", "fake", "web/**")
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){"fake": stdout("report")}}
	// api is steady at 5; web jumps to 40 and must not drag api's range.
	readers := Readers{ProviderLizardCSV: readerFunc(func(_ context.Context, in ReportInput) (Reading, error) {
		if in.Scope[0] == "web/**" {
			return goodReading(40, UnitCCN), nil
		}
		return goodReading(5, UnitCCN), nil
	})}
	for i := 0; i < 4; i++ {
		n := i
		if _, err := Collect(context.Background(), CollectRequest{Root: root, Origin: OriginOfficial, Config: cfgOf(api, web), Readers: readers, Runner: tools,
			Clock: func() time.Time { return testClock().Add(time.Duration(n) * time.Second) }}); err != nil {
			t.Fatal(err)
		}
	}
	snaps, err := NewStore(root).LoadSnapshots()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range snaps[len(snaps)-1].Summary.Capabilities {
		if c.Capability == CapabilityComplexity && c.Name == "api" && strings.Contains(c.Explanation, "Not enough comparable") {
			t.Errorf("api lost its own history: %s", c.Explanation)
		}
	}
}

func lizardInstances(n int) []CapabilityConfig {
	var items []CapabilityConfig
	for i := 0; i < n; i++ {
		items = append(items, lizardInstance(fmt.Sprintf("n%d", i), "fake", fmt.Sprintf("s%d/**", i)))
	}
	return items
}

func TestConfigCapacityCountsPlaceholdersForUnconfiguredCapabilities(t *testing.T) {
	missing := len(Capabilities()) - 1
	atLimit := MaxResults - missing
	if err := cfgOf(lizardInstances(atLimit)...).Validate(); err != nil {
		t.Errorf("%d instances of one capability should fit: %v", atLimit, err)
	}
	if err := cfgOf(lizardInstances(atLimit + 1)...).Validate(); err == nil {
		t.Errorf("%d instances of one capability leave no room for placeholders", atLimit+1)
	}

	all := Config{}
	for _, c := range Capabilities() {
		all.Capabilities = append(all.Capabilities, CapabilityConfig{Capability: c})
	}
	if got := all.resultCount(); got != len(Capabilities()) {
		t.Errorf("every capability configured once is %d results, got %d", len(Capabilities()), got)
	}
}

func TestCollectRefusesAnOverfullConfigurationBeforeRunningTools(t *testing.T) {
	root := project(t)
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){"fake": stdout("report")}}
	_, err := Collect(context.Background(), CollectRequest{Root: root, Origin: OriginManual, Config: cfgOf(lizardInstances(MaxResults)...),
		Readers: Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, Runner: tools, Clock: testClock})
	if err == nil || len(tools.calls) != 0 {
		t.Fatalf("err=%v, tools run=%d", err, len(tools.calls))
	}

	atLimit := MaxResults - (len(Capabilities()) - 1)
	cfg := cfgOf(lizardInstances(atLimit)...)
	got := collect(t, root, cfg, Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, tools)
	if len(got.Results) != MaxResults {
		t.Errorf("a full configuration saved %d results, want %d", len(got.Results), MaxResults)
	}
}
