package codehealth

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func lizard(name string, scope ...string) CapabilityConfig {
	return CapabilityConfig{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: name, Scope: scope}
}

func TestConfigInstances(t *testing.T) {
	tests := []struct {
		name  string
		items []CapabilityConfig
		want  error
	}{
		{"one unnamed instance", []CapabilityConfig{lizard("")}, nil},
		{"two named instances with distinct scope", []CapabilityConfig{lizard("api", "api/**"), lizard("web", "web/**")}, nil},
		{"same name and provider is a duplicate", []CapabilityConfig{lizard("api", "api/**"), lizard("api", "web/**")}, ErrDuplicateInstance},
		{"two unnamed instances are a duplicate", []CapabilityConfig{lizard(""), lizard("")}, ErrDuplicateInstance},
		{"unnamed beside a named instance", []CapabilityConfig{lizard("", "api/**"), lizard("web", "web/**")}, ErrInvalidConfig},
		{"same scope refused", []CapabilityConfig{lizard("api", "src/**"), lizard("web", "src/**")}, ErrInvalidConfig},
		{"same scope in another order refused", []CapabilityConfig{lizard("api", "a/**", "b/**"), lizard("web", "b/**", "a/**")}, ErrInvalidConfig},
		{"both without scope refused", []CapabilityConfig{lizard("api"), lizard("web")}, ErrInvalidConfig},
		{"name must be a token", []CapabilityConfig{lizard("two\nlines")}, ErrUnboundedDetail},
		{"name length bounded", []CapabilityConfig{lizard(strings.Repeat("n", MaxTokenLen+1))}, ErrUnboundedDetail},
		{"other providers need no name", []CapabilityConfig{
			lizard(""),
			{Capability: CapabilityDuplication, Provider: ProviderJscpdJSON},
		}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Config{Version: ConfigVersion, Capabilities: tt.items}.Validate()
			if !errors.Is(err, tt.want) || (tt.want == nil && err != nil) {
				t.Fatalf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestConfigNameAndRequiredRoundTrip(t *testing.T) {
	in := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "api", Required: true, Scope: []string{"api/**"}},
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "web", Scope: []string{"web/**"}},
	}}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), `"required":false`) {
		t.Fatalf("optional instance should omit required: %s", b)
	}
	out, err := DecodeConfig(b)
	if err != nil {
		t.Fatal(err)
	}
	if out.Capabilities[0].Name != "api" || !out.Capabilities[0].Required || out.Capabilities[1].Required {
		t.Fatalf("round trip lost name or required: %+v", out.Capabilities)
	}
	if _, err := DecodeConfig([]byte(`{"version":1,"capabilities":[{"capability":"tests","provider":"go-test-json","extra":1}]}`)); err == nil {
		t.Fatal("unknown field was accepted")
	}
}

func TestEffectiveTimeout(t *testing.T) {
	tests := []struct {
		name     string
		provider ProviderKey
		capab    Capability
		override int
		want     time.Duration
		ok       bool
	}{
		{"lizard default", ProviderLizardCSV, CapabilityComplexity, 0, 60 * time.Second, true},
		{"jscpd default", ProviderJscpdJSON, CapabilityDuplication, 0, 60 * time.Second, true},
		{"osv default", ProviderOSVScannerJSON, CapabilityDependencyVulnerability, 0, 120 * time.Second, true},
		{"override wins", ProviderOSVScannerJSON, CapabilityDependencyVulnerability, 600, 600 * time.Second, true},
		{"override below default", ProviderLizardCSV, CapabilityComplexity, 5, 5 * time.Second, true},
		{"report-only has none", ProviderGoTestJSON, CapabilityTests, 0, 0, false},
		{"report-only coverage has none", ProviderGoCoverProfile, CapabilityCoverage, 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := CapabilityConfig{Capability: tt.capab, Provider: tt.provider, TimeoutSeconds: tt.override}.EffectiveTimeout()
			if got != tt.want || ok != tt.ok {
				t.Fatalf("EffectiveTimeout() = %v, %v; want %v, %v", got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestEveryExecutedProviderHasADefaultTimeout(t *testing.T) {
	for p, c := range providerCapability {
		_, ok := DefaultTimeoutSeconds(p)
		executed := c == CapabilityComplexity || c == CapabilityDuplication || c == CapabilityDependencyVulnerability
		if ok != executed {
			t.Errorf("%s: default timeout present = %v, want %v", p, ok, executed)
		}
	}
}

func TestTimeoutBounds(t *testing.T) {
	for _, tt := range []struct {
		seconds int
		want    error
	}{{0, nil}, {1, nil}, {600, nil}, {601, ErrInvalidConfig}, {-1, ErrInvalidConfig}} {
		c := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{lizard("")}}
		c.Capabilities[0].TimeoutSeconds = tt.seconds
		if err := c.Validate(); !errors.Is(err, tt.want) || (tt.want == nil && err != nil) {
			t.Errorf("timeout %d: Validate() = %v, want %v", tt.seconds, err, tt.want)
		}
	}
	if MaxTimeoutSecond != 600 {
		t.Fatalf("MaxTimeoutSecond = %d, want 600", MaxTimeoutSecond)
	}
}

// namedSnapshot holds two complexity instances of one provider.
func namedSnapshot(t *testing.T) Snapshot {
	t.Helper()
	s := validSnapshot(t)
	mk := func(name, scope string, v float64) CapabilityResult {
		r := result(CapabilityComplexity, v)
		r.Name, r.Scope = name, []string{scope}
		return r
	}
	s.Results = []CapabilityResult{mk("api", "api/**", 4), mk("web", "web/**", 12)}
	s.Summary = Summary{Overall: ClassificationWatch, Capabilities: []CapabilitySummary{
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "web", Classification: ClassificationWatch},
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Name: "api", Classification: ClassificationWatch},
	}}
	return reseal(s)
}

func TestSnapshotNamedInstances(t *testing.T) {
	s := namedSnapshot(t)
	if err := s.Validate(); err != nil {
		t.Fatalf("two named instances should validate: %v", err)
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	back, err := DecodeSnapshot(b)
	if err != nil {
		t.Fatalf("round trip: %v", err)
	}
	if back.Results[0].Name != "api" || back.Summary.Capabilities[0].Name != "web" {
		t.Fatalf("name lost on round trip: %+v", back)
	}
	if s.ID == fixtureSnapshotID {
		t.Fatal("a named snapshot must hash differently from the pinned unnamed one")
	}

	tests := []struct {
		name   string
		mutate func(*Snapshot)
		want   error
	}{
		{"duplicate named result", func(s *Snapshot) { s.Results[1].Name = "api" }, ErrDuplicateInstance},
		{"bad name", func(s *Snapshot) { s.Results[0].Name = "a\nb" }, ErrUnboundedDetail},
		{"summary for unknown name", func(s *Snapshot) { s.Summary.Capabilities[0].Name = "mobile" }, ErrInvalidSummary},
		{"summary drops the name", func(s *Snapshot) { s.Summary.Capabilities[0].Name = "" }, ErrInvalidSummary},
		{"not configured with a name", func(s *Snapshot) {
			s.Results[0] = CapabilityResult{Capability: CapabilityComplexity, Name: "api", Outcome: OutcomeNotConfigured, Freshness: FreshnessUnknown, CollectedAt: "2026-10-01T00:00:00Z"}
		}, ErrInvalidOutcome},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := namedSnapshot(t)
			tt.mutate(&c)
			if err := reseal(c).Validate(); !errors.Is(err, tt.want) {
				t.Fatalf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSnapshotOmitsEmptyName(t *testing.T) {
	b, err := json.Marshal(validSnapshot(t))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), `"name"`) {
		t.Fatalf("unnamed snapshot serialized a name: %s", b)
	}
}

func TestSnapshotIDIgnoresInstanceOrder(t *testing.T) {
	s := namedSnapshot(t)
	swapped := s
	swapped.Results = []CapabilityResult{s.Results[1], s.Results[0]}
	if swapped.ComputeID() != s.ComputeID() {
		t.Fatal("result order changed the snapshot ID")
	}
}

func TestRenameKeepsSeriesButScopeStartsOne(t *testing.T) {
	base := result(CapabilityComplexity, 5)
	base.Name, base.Scope = "api", []string{"api/**"}
	renamed := base
	renamed.Name = "backend"
	if base.SeriesID() != renamed.SeriesID() {
		t.Fatal("renaming an instance started a new series")
	}
	rescoped := base
	rescoped.Scope = []string{"api/v2/**"}
	if base.SeriesID() == rescoped.SeriesID() {
		t.Fatal("changing scope must start a new series")
	}
	unnamed := base
	unnamed.Name = ""
	if base.SeriesID() != unnamed.SeriesID() {
		t.Fatal("series identity depends on name")
	}
}

func TestNameStaysOutOfConfigDigest(t *testing.T) {
	a, b := lizard("api", "api/**"), lizard("backend", "api/**")
	b.Required = true
	if a.Digest() != b.Digest() {
		t.Fatal("name or required changed the config digest")
	}
}

func TestInstancesAreAssessedSeparately(t *testing.T) {
	scoped := func(name, scope string, v float64) CapabilityResult {
		r := result(CapabilityComplexity, v)
		r.Name, r.Scope = name, []string{scope}
		return r
	}
	var history []HistoryEntry
	for _, v := range []float64{5, 5, 5} {
		history = append(history,
			HistoryEntry{Origin: OriginOfficial, Result: scoped("api", "api/**", v)},
			HistoryEntry{Origin: OriginOfficial, Result: scoped("web", "web/**", v+30)})
	}
	api := Assess(OriginOfficial, scoped("api", "api/**", 5), nil, history)
	web := Assess(OriginOfficial, scoped("web", "web/**", 35), nil, history)

	if api.Name != "api" || web.Name != "web" || api.Summary().Name != "api" {
		t.Fatalf("assessment lost the instance name: %q %q", api.Name, web.Name)
	}
	if api.Trend.Observations != 4 || web.Trend.Observations != 4 {
		t.Fatalf("each instance should see only its own 3 earlier checks plus current: api %d, web %d", api.Trend.Observations, web.Trend.Observations)
	}
	if api.Trend.Max != 5 || web.Trend.Min != 35 {
		t.Fatalf("values mixed across instances: api max %v, web min %v", api.Trend.Max, web.Trend.Min)
	}
	if api.Classification != ClassificationGood {
		t.Fatalf("api = %s, want good", api.Classification)
	}
	if web.Classification != ClassificationNeedsAttention {
		t.Fatalf("web = %s, want needs_attention", web.Classification)
	}
	if got := Overall([]Assessment{api, web}); got != ClassificationNeedsAttention {
		t.Fatalf("Overall() = %s; a healthy instance must not hide an unhealthy one", got)
	}
	if !strings.Contains(web.Explanation, "earlier official result") {
		t.Fatalf("the other instance's history should be reported as not compared, got %q", web.Explanation)
	}
}
