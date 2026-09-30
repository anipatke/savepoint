package codehealth

import (
	"errors"
	"testing"
)

// Pinned vectors: a change here means persisted identities would change.
const (
	fixtureSnapshotID   = "sha256:67ad865be44dbc76d9e228ccbe7ae21f910130e86322d593fd3387e7cb457aa4"
	fixtureConfigDigest = "sha256:0ac92110764657552d88cf10cc65539a8da93ff79aa30822512801f60a8f22d6"
)

func TestSnapshotIDVector(t *testing.T) {
	if got := validSnapshot(t).ComputeID(); got != fixtureSnapshotID {
		t.Fatalf("ComputeID() = %s, want %s", got, fixtureSnapshotID)
	}
}

func TestSnapshotIDIgnoresListOrder(t *testing.T) {
	s := validSnapshot(t)
	shuffled := s
	shuffled.Results = []CapabilityResult{s.Results[2], s.Results[0], s.Results[1]}
	shuffled.Summary.Capabilities = []CapabilitySummary{s.Summary.Capabilities[1], s.Summary.Capabilities[2], s.Summary.Capabilities[0]}
	shuffled.Results[1].Details = []Detail{s.Results[0].Details[1], s.Results[0].Details[0]}
	if err := shuffled.Validate(); err != nil {
		t.Fatalf("reordered snapshot should keep its ID: %v", err)
	}
}

func TestSnapshotIDChangesWithObservation(t *testing.T) {
	base := validSnapshot(t)
	tests := map[string]func(*Snapshot){
		"commit":    func(s *Snapshot) { s.Repository.Commit = "0000000000000000000000000000000000000001" },
		"time":      func(s *Snapshot) { s.CreatedAt = "2026-10-01T10:00:00Z" },
		"value":     func(s *Snapshot) { s.Results[0].Value = &Value{Number: 90, Unit: UnitPercent} },
		"origin":    func(s *Snapshot) { s.Origin, s.Retention = OriginManual, RetentionPrunable },
		"evidence":  func(s *Snapshot) { s.Results[0].Evidence = nil },
		"freshness": func(s *Snapshot) { s.Results[0].Freshness = FreshnessStale },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			s := validSnapshot(t)
			mutate(&s)
			if s.ComputeID() == base.ID {
				t.Fatalf("ID did not change when %s changed", name)
			}
		})
	}
}

// TestMutationAfterValidationIsCaught edits a validated snapshot in place and
// checks that a fresh Validate call notices.
func TestMutationAfterValidationIsCaught(t *testing.T) {
	s := validSnapshot(t)
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.Results[0].Value.Number = 99
	if err := s.Validate(); !errors.Is(err, ErrIdentityMismatch) {
		t.Fatalf("Validate() after mutation = %v, want ErrIdentityMismatch", err)
	}
}

func TestCanonicalDoesNotMutateInput(t *testing.T) {
	s := validSnapshot(t)
	s.Results[0].Scope = []string{"b", "a"}
	_ = s.Canonical()
	if s.Results[0].Scope[0] != "b" {
		t.Fatal("Canonical sorted the caller's slice")
	}
}

func seriesFixture(t *testing.T) CapabilityResult {
	return validSnapshot(t).Results[0]
}

func TestSeriesIDIgnoresObservation(t *testing.T) {
	base := seriesFixture(t)
	want := base.SeriesID()
	observed := base
	observed.Value = &Value{Number: 12, Unit: UnitPercent}
	observed.CollectedAt = "2027-01-01T00:00:00Z"
	observed.Freshness = FreshnessStale
	observed.Evidence = nil
	observed.Reason = "different"
	if observed.SeriesID() != want {
		t.Fatal("series ID changed with value, time, freshness, or evidence")
	}
	reordered := base
	reordered.Scope = []string{"z", "a"}
	other := base
	other.Scope = []string{"a", "z"}
	if reordered.SeriesID() != other.SeriesID() {
		t.Fatal("series ID depends on scope order")
	}
}

func TestSeriesIDChangesWithCompatibility(t *testing.T) {
	base := seriesFixture(t)
	tests := map[string]func(*CapabilityResult){
		"provider":               func(r *CapabilityResult) { r.Provenance.Provider = ProviderCoveragePyJSON },
		"provider version":       func(r *CapabilityResult) { r.Provenance.ProviderVersion = "go1.27.0" },
		"report schema":          func(r *CapabilityResult) { r.Provenance.ReportSchema = "cover-profile/atomic" },
		"measurement definition": func(r *CapabilityResult) { r.Provenance.MeasurementDefinition = "statements/2" },
		"effective configuration": func(r *CapabilityResult) {
			r.ConfigDigest = "sha256:" + "9999999999999999999999999999999999999999999999999999999999999999"
		},
		"exclusions": func(r *CapabilityResult) { r.Exclusions = append(r.Exclusions, "generated/**") },
		"scope":      func(r *CapabilityResult) { r.Scope = []string{"cmd/**"} },
		"capability": func(r *CapabilityResult) { r.Capability = CapabilityDuplication },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			r := base
			mutate(&r)
			if r.SeriesID() == base.SeriesID() {
				t.Fatalf("series ID did not change with %s", name)
			}
		})
	}
}

func TestSeriesIDExcludesRepositoryCommit(t *testing.T) {
	a, b := validSnapshot(t), validSnapshot(t)
	b.Repository.Commit = "0000000000000000000000000000000000000002"
	b.Repository.InputFingerprint = "sha256:" + "4444444444444444444444444444444444444444444444444444444444444444"
	if a.Results[0].SeriesID() != b.Results[0].SeriesID() {
		t.Fatal("a new commit must not start a new series")
	}
}

func TestConfigDigest(t *testing.T) {
	cfg := validConfig(t)
	base := cfg.Capabilities[0]
	same := base
	same.TimeoutSeconds = 30
	same.Thresholds = &Threshold{Good: 1, Watch: 4}
	same.Executable = "/usr/local/go/bin/go"
	same.Report = "other.json"
	if base.Digest() != same.Digest() {
		t.Error("digest changed with timeout, thresholds, executable location, or report path")
	}
	changed := base
	changed.Args = []string{"test", "-json", "-race", "./..."}
	if base.Digest() == changed.Digest() {
		t.Error("digest ignored a change to measurement arguments")
	}
	if err := validateDigest("d", base.Digest()); err != nil {
		t.Errorf("digest is malformed: %v", err)
	}
	if got := base.Digest(); got != fixtureConfigDigest {
		t.Errorf("Digest() = %s, want pinned vector %s", got, fixtureConfigDigest)
	}
}
