package codehealth

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"
)

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func validConfig(t *testing.T) Config {
	t.Helper()
	c, err := DecodeConfig(readFixture(t, "valid-config-v1.json"))
	if err != nil {
		t.Fatalf("fixture config invalid: %v", err)
	}
	return c
}

func validSnapshot(t *testing.T) Snapshot {
	t.Helper()
	s, err := DecodeSnapshot(readFixture(t, "valid-snapshot-v1.json"))
	if err != nil {
		t.Fatalf("fixture snapshot invalid: %v", err)
	}
	return s
}

// reseal recomputes the ID after a deliberate edit so the test reaches the
// rule it targets instead of failing on the identity check.
func reseal(s Snapshot) Snapshot {
	s.ID = s.ComputeID()
	return s
}

func TestCapabilitiesAreTheFixedFive(t *testing.T) {
	want := []Capability{"tests", "coverage", "complexity", "duplication", "dependency_vulnerabilities"}
	if !slices.Equal(Capabilities(), want) {
		t.Fatalf("Capabilities() = %v, want %v", Capabilities(), want)
	}
	for _, c := range want {
		if _, ok := capabilityUnits[c]; !ok {
			t.Errorf("capability %s has no unit", c)
		}
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Config)
		want   error
	}{
		{"valid fixture", func(*Config) {}, nil},
		{"empty capability list is valid", func(c *Config) { c.Capabilities = nil }, nil},
		{"unsupported version", func(c *Config) { c.Version = 2 }, ErrUnsupportedVersion},
		{"unknown capability", func(c *Config) { c.Capabilities[0].Capability = "hotspots" }, ErrUnknownCapability},
		{"sixth capability rejected", func(c *Config) { c.Capabilities[0].Capability = "dead_code" }, ErrUnknownCapability},
		{"unknown provider", func(c *Config) { c.Capabilities[0].Provider = "made-up" }, ErrUnknownProvider},
		{"provider serves other capability", func(c *Config) { c.Capabilities[0].Provider = ProviderJscpdJSON }, ErrProviderMismatch},
		{"duplicate instance", func(c *Config) { c.Capabilities = append(c.Capabilities, c.Capabilities[0]) }, ErrDuplicateInstance},
		{"second provider for same capability is allowed", func(c *Config) {
			c.Capabilities = append(c.Capabilities, CapabilityConfig{Capability: CapabilityTests, Provider: ProviderPytestJUnit})
		}, nil},
		{"absolute report path", func(c *Config) { c.Capabilities[0].Report = "/etc/report.json" }, ErrInvalidConfig},
		{"traversing report path", func(c *Config) { c.Capabilities[0].Report = "../out.json" }, ErrInvalidConfig},
		{"windows drive report path", func(c *Config) { c.Capabilities[0].Report = `C:\out.json` }, ErrInvalidConfig},
		{"duplicate scope entry", func(c *Config) { c.Capabilities[0].Scope = []string{"a/**", "a/**"} }, ErrInvalidConfig},
		{"negative timeout", func(c *Config) { c.Capabilities[0].TimeoutSeconds = -1 }, ErrInvalidConfig},
		{"timeout over cap", func(c *Config) { c.Capabilities[0].TimeoutSeconds = MaxTimeoutSecond + 1 }, ErrInvalidConfig},
		{"too many args", func(c *Config) { c.Capabilities[0].Args = make([]string, MaxArgs+1) }, ErrUnboundedDetail},
		{"oversized arg", func(c *Config) { c.Capabilities[0].Args = []string{strings.Repeat("a", MaxArgLen+1)} }, ErrUnboundedDetail},
		{"too many scope entries", func(c *Config) {
			for i := 0; i <= MaxScopeEntries; i++ {
				c.Capabilities[0].Scope = append(c.Capabilities[0].Scope, "d"+string(rune('a'+i%26))+strings.Repeat("x", i))
			}
		}, ErrUnboundedDetail},
		{"coverage threshold order", func(c *Config) { c.Capabilities[1].Thresholds = &Threshold{Good: 50, Watch: 70} }, ErrInvalidThreshold},
		{"lower-is-better threshold order", func(c *Config) { c.Capabilities[0].Thresholds = &Threshold{Good: 5, Watch: 1} }, ErrInvalidThreshold},
		{"percent threshold over 100", func(c *Config) { c.Capabilities[1].Thresholds = &Threshold{Good: 101, Watch: 60} }, ErrInvalidThreshold},
		{"negative threshold", func(c *Config) { c.Capabilities[0].Thresholds = &Threshold{Good: -1, Watch: 2} }, ErrInvalidThreshold},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validConfig(t)
			tt.mutate(&c)
			err := c.Validate()
			if !errors.Is(err, tt.want) || (tt.want == nil && err != nil) {
				t.Fatalf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestDecodeConfigRejectsMalformedInput(t *testing.T) {
	for name, in := range map[string]string{
		"not json":      `{`,
		"unknown field": `{"version":1,"capabilities":[],"extra":true}`,
		"trailing data": `{"version":1,"capabilities":[]} {}`,
		"wrong type":    `{"version":"1"}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := DecodeConfig([]byte(in)); !errors.Is(err, ErrMalformedRecord) {
				t.Fatalf("DecodeConfig() = %v, want ErrMalformedRecord", err)
			}
		})
	}
}

func TestOutcomeAndFreshnessAreOrthogonal(t *testing.T) {
	for _, outcome := range Outcomes() {
		for _, fresh := range []Freshness{FreshnessFresh, FreshnessStale, FreshnessUnknown} {
			t.Run(string(outcome)+"/"+string(fresh), func(t *testing.T) {
				r := resultFor(outcome, fresh)
				err := r.validate("r")
				wantOK := fresh == FreshnessUnknown || outcome.Measured()
				if wantOK && err != nil {
					t.Fatalf("validate() = %v, want valid", err)
				}
				if !wantOK && !errors.Is(err, ErrInvalidFreshness) {
					t.Fatalf("validate() = %v, want ErrInvalidFreshness", err)
				}
			})
		}
	}
}

// resultFor builds a structurally valid result for any outcome.
func resultFor(o Outcome, f Freshness) CapabilityResult {
	r := CapabilityResult{
		Capability:  CapabilityCoverage,
		Outcome:     o,
		Freshness:   f,
		CollectedAt: "2026-09-30T09:59:00Z",
		Provenance: Provenance{
			Provider: ProviderGoCoverProfile, ProviderVersion: "go1.26.2", MeasurementDefinition: "statements/1",
		},
		ConfigDigest: "sha256:" + strings.Repeat("2", 64),
	}
	switch o {
	case OutcomeNotConfigured:
		r.Provenance, r.ConfigDigest = Provenance{}, ""
	case OutcomeAvailable:
		r.Value = &Value{Number: 50, Unit: UnitPercent}
	case OutcomePartial:
		r.Value = &Value{Number: 50, Unit: UnitPercent}
		r.Reason = "report truncated at 10000 rows"
	}
	return r
}

func TestOnlyMeasuredOutcomesCarryValues(t *testing.T) {
	for _, o := range Outcomes() {
		t.Run(string(o), func(t *testing.T) {
			r := resultFor(o, FreshnessUnknown)
			r.Value = &Value{Number: 0, Unit: UnitPercent} // a zero must not sneak in
			if o == OutcomePartial {
				r.Reason = "truncated"
			}
			err := r.validate("r")
			if o.Measured() && err != nil {
				t.Fatalf("measured outcome rejected a value: %v", err)
			}
			if !o.Measured() && !errors.Is(err, ErrValueNotAllowed) {
				t.Fatalf("validate() = %v, want ErrValueNotAllowed", err)
			}
		})
	}
}

func TestMeasuredOutcomeRules(t *testing.T) {
	t.Run("available without a value", func(t *testing.T) {
		r := resultFor(OutcomeAvailable, FreshnessFresh)
		r.Value = nil
		if err := r.validate("r"); !errors.Is(err, ErrMissingValue) {
			t.Fatalf("got %v, want ErrMissingValue", err)
		}
	})
	t.Run("partial without a reason", func(t *testing.T) {
		r := resultFor(OutcomePartial, FreshnessFresh)
		r.Reason = ""
		if err := r.validate("r"); !errors.Is(err, ErrMissingValue) {
			t.Fatalf("got %v, want ErrMissingValue", err)
		}
	})
	t.Run("partial may omit the value", func(t *testing.T) {
		r := resultFor(OutcomePartial, FreshnessUnknown)
		r.Value = nil
		if err := r.validate("r"); err != nil {
			t.Fatalf("got %v, want valid", err)
		}
	})
	t.Run("measured zero is valid", func(t *testing.T) {
		r := resultFor(OutcomeAvailable, FreshnessFresh)
		r.Value = &Value{Number: 0, Unit: UnitPercent}
		if err := r.validate("r"); err != nil {
			t.Fatalf("got %v, want valid", err)
		}
	})
}

func TestValueUnitCompatibility(t *testing.T) {
	tests := []struct {
		name string
		c    Capability
		v    Value
		ok   bool
	}{
		{"coverage percent", CapabilityCoverage, Value{84.5, UnitPercent}, true},
		{"coverage over 100", CapabilityCoverage, Value{100.1, UnitPercent}, false},
		{"coverage as count", CapabilityCoverage, Value{5, UnitCount}, false},
		{"tests count", CapabilityTests, Value{3, UnitCount}, true},
		{"tests fractional count", CapabilityTests, Value{2.5, UnitCount}, false},
		{"tests as percent", CapabilityTests, Value{3, UnitPercent}, false},
		{"complexity ccn", CapabilityComplexity, Value{12.5, UnitCCN}, true},
		{"complexity as percent", CapabilityComplexity, Value{12, UnitPercent}, false},
		{"duplication percent", CapabilityDuplication, Value{4, UnitPercent}, true},
		{"vulnerabilities count", CapabilityDependencyVulnerability, Value{0, UnitCount}, true},
		{"negative", CapabilityTests, Value{-1, UnitCount}, false},
		{"unit missing", CapabilityCoverage, Value{50, ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.validate("v", tt.c)
			if tt.ok && err != nil {
				t.Fatalf("got %v, want valid", err)
			}
			if !tt.ok && !errors.Is(err, ErrIncompatibleUnit) {
				t.Fatalf("got %v, want ErrIncompatibleUnit", err)
			}
		})
	}
}

func TestNonFiniteNumbersAreRejected(t *testing.T) {
	nan := zeroOverZero()
	if err := (Value{nan, UnitPercent}).validate("v", CapabilityCoverage); !errors.Is(err, ErrIncompatibleUnit) {
		t.Fatalf("NaN value: got %v", err)
	}
	if err := (Threshold{Good: nan, Watch: 1}).validate("t", CapabilityTests); !errors.Is(err, ErrInvalidThreshold) {
		t.Fatalf("NaN threshold: got %v", err)
	}
}

func zeroOverZero() float64 {
	z := 0.0
	return z / z
}

func TestEvidenceReferenceSafety(t *testing.T) {
	tests := []struct {
		path string
		ok   bool
	}{
		{"coverage.out", true},
		{"internal/data/gate_v2.go", true},
		{"/etc/passwd", false},
		{"~/secrets.txt", false},
		{"../outside", false},
		{"a/../../b", false},
		{"a//b", false},
		{"./a", false},
		{`a\b`, false},
		{"C:/Users/me/x", false},
		{"https://example.com/x", false},
		{"", false},
		{strings.Repeat("a/", MaxPathLen), false},
		{"dir/.env", false},
		{"dir/.env.production", false},
		{".git/config", false},
		{"home/.ssh/id_rsa", false},
		{"keys/server.PEM", false},
		{"deploy/credentials.json", false},
		{"bad\x00name", false},
		{"bad\nname", false},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := EvidenceRef{Path: tt.path}.validate("e")
			if tt.ok && err != nil {
				t.Fatalf("got %v, want valid", err)
			}
			if !tt.ok && !errors.Is(err, ErrUnsafeEvidence) {
				t.Fatalf("got %v, want ErrUnsafeEvidence", err)
			}
		})
	}
	if err := (EvidenceRef{Path: "a.go", Line: -1}).validate("e"); !errors.Is(err, ErrUnsafeEvidence) {
		t.Errorf("negative line: got %v", err)
	}
	if err := (EvidenceRef{Path: "a.go", Note: strings.Repeat("n", MaxNoteLen+1)}).validate("e"); !errors.Is(err, ErrUnboundedDetail) {
		t.Errorf("long note: got %v", err)
	}
}

func TestSnapshotValidation(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Snapshot)
		want   error
	}{
		{"valid fixture", func(*Snapshot) {}, nil},
		{"unsupported version", func(s *Snapshot) { s.Version = 9 }, ErrUnsupportedVersion},
		{"unknown origin", func(s *Snapshot) { s.Origin = "nightly" }, ErrInvalidOrigin},
		{"official but prunable", func(s *Snapshot) { s.Retention = RetentionPrunable }, ErrInvalidOrigin},
		{"manual but permanent", func(s *Snapshot) { s.Origin = OriginManual }, ErrInvalidOrigin},
		{"manual prunable is valid", func(s *Snapshot) { s.Origin, s.Retention = OriginManual, RetentionPrunable }, nil},
		{"created_at with offset", func(s *Snapshot) { s.CreatedAt = "2026-09-30T10:00:00+01:00" }, ErrMalformedTimestamp},
		{"created_at date only", func(s *Snapshot) { s.CreatedAt = "2026-09-30" }, ErrMalformedTimestamp},
		{"created_at fractional", func(s *Snapshot) { s.CreatedAt = "2026-09-30T10:00:00.5Z" }, ErrMalformedTimestamp},
		{"collected_at impossible date", func(s *Snapshot) { s.Results[0].CollectedAt = "2026-02-30T10:00:00Z" }, ErrMalformedTimestamp},
		{"short commit", func(s *Snapshot) { s.Repository.Commit = "9fdc11f" }, ErrMalformedIdentity},
		{"uppercase commit", func(s *Snapshot) { s.Repository.Commit = strings.ToUpper(s.Repository.Commit) }, ErrMalformedIdentity},
		{"bad fingerprint", func(s *Snapshot) { s.Repository.InputFingerprint = "abc" }, ErrMalformedIdentity},
		{"bad config digest", func(s *Snapshot) { s.Results[0].ConfigDigest = "sha256:xyz" }, ErrMalformedIdentity},
		{"no commit and clean", func(s *Snapshot) { s.Repository.Commit, s.Repository.Dirty = "", false }, ErrInvalidRepository},
		{"no commit and dirty is valid", func(s *Snapshot) { s.Repository.Commit, s.Repository.Dirty = "", true }, nil},
		{"dirty with commit is valid", func(s *Snapshot) { s.Repository.Dirty = true }, nil},
		{"unknown capability", func(s *Snapshot) { s.Results[0].Capability = "hotspots" }, ErrUnknownCapability},
		{"unknown outcome", func(s *Snapshot) { s.Results[0].Outcome = "ok" }, ErrInvalidOutcome},
		{"unknown freshness", func(s *Snapshot) { s.Results[0].Freshness = "recent" }, ErrInvalidFreshness},
		{"failed but fresh", func(s *Snapshot) { s.Results[1].Freshness = FreshnessFresh }, ErrInvalidFreshness},
		{"failed carrying a value", func(s *Snapshot) { s.Results[1].Value = &Value{Number: 0, Unit: UnitCount} }, ErrValueNotAllowed},
		{"not configured with provider", func(s *Snapshot) { s.Results[2].Provenance.Provider = ProviderLizardCSV }, ErrInvalidOutcome},
		{"provider from wrong capability", func(s *Snapshot) { s.Results[0].Provenance.Provider = ProviderJscpdJSON }, ErrProviderMismatch},
		{"measured without provider version", func(s *Snapshot) { s.Results[0].Provenance.ProviderVersion = "" }, ErrMalformedRecord},
		{"measured without definition", func(s *Snapshot) { s.Results[0].Provenance.MeasurementDefinition = "" }, ErrMalformedRecord},
		{"incompatible unit", func(s *Snapshot) { s.Results[0].Value.Unit = UnitCount }, ErrIncompatibleUnit},
		{"duplicate result instance", func(s *Snapshot) { s.Results = append(s.Results, s.Results[0]) }, ErrDuplicateInstance},
		{"absolute evidence", func(s *Snapshot) { s.Results[0].Evidence[0].Path = "/home/me/coverage.out" }, ErrUnsafeEvidence},
		{"sensitive evidence", func(s *Snapshot) { s.Results[0].Evidence[0].Path = "config/.env" }, ErrUnsafeEvidence},
		{"too many details", func(s *Snapshot) {
			for i := 0; i <= MaxDetails; i++ {
				s.Results[0].Details = append(s.Results[0].Details, Detail{Key: "k" + strings.Repeat("x", i), Number: 1})
			}
		}, ErrUnboundedDetail},
		{"duplicate detail key", func(s *Snapshot) { s.Results[0].Details[1].Key = s.Results[0].Details[0].Key }, ErrUnboundedDetail},
		{"too many evidence refs", func(s *Snapshot) {
			for i := 0; i <= MaxEvidence; i++ {
				s.Results[0].Evidence = append(s.Results[0].Evidence, EvidenceRef{Path: "f" + strings.Repeat("x", i)})
			}
		}, ErrUnboundedDetail},
		{"oversized reason", func(s *Snapshot) { s.Results[1].Reason = strings.Repeat("r", MaxReasonLen+1) }, ErrUnboundedDetail},
		{"multiline reason", func(s *Snapshot) { s.Results[1].Reason = "a\nb" }, ErrUnboundedDetail},
		{"summary count mismatch", func(s *Snapshot) { s.Summary.Capabilities = s.Summary.Capabilities[:1] }, ErrInvalidSummary},
		{"summary for unknown result", func(s *Snapshot) { s.Summary.Capabilities[0].Capability = CapabilityDuplication }, ErrInvalidSummary},
		{"unknown classification", func(s *Snapshot) { s.Summary.Overall = "fine" }, ErrInvalidSummary},
		{"good from failed evidence", func(s *Snapshot) { s.Summary.Capabilities[1].Classification = ClassificationGood }, ErrInvalidSummary},
		{"good from stale evidence", func(s *Snapshot) { s.Results[0].Freshness = FreshnessStale }, ErrInvalidSummary},
		{"good overall despite failure", func(s *Snapshot) { s.Summary.Overall = ClassificationGood }, ErrInvalidSummary},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := validSnapshot(t)
			tt.mutate(&s)
			err := reseal(s).Validate()
			if tt.want == nil && err != nil || !errors.Is(err, tt.want) {
				t.Fatalf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestDecodeSnapshotRejectsMalformedInput(t *testing.T) {
	for name, in := range map[string]string{
		"not json":      `[`,
		"unknown field": `{"version":1,"surprise":1}`,
		"wrong type":    `{"version":"one"}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := DecodeSnapshot([]byte(in)); !errors.Is(err, ErrMalformedRecord) {
				t.Fatalf("got %v, want ErrMalformedRecord", err)
			}
		})
	}
}
