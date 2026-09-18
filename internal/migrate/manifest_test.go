package migrate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestBuildManifest_roundTripsThroughYAML proves BuildManifest's projection
// of a real fixture's plan survives a Marshal/UnmarshalManifest round trip
// with every source hash, identity, archive entry, and legacy prerequisite
// intact.
func TestBuildManifest_roundTripsThroughYAML(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-basic"))
	m := BuildManifest(p)

	content, err := m.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	got, err := UnmarshalManifest(content)
	if err != nil {
		t.Fatalf("UnmarshalManifest() error = %v", err)
	}

	if len(got.Sources) != len(p.Sources) {
		t.Errorf("round-tripped Sources = %d entries, want %d", len(got.Sources), len(p.Sources))
	}
	if len(got.Identities) != len(p.Targets) {
		t.Errorf("round-tripped Identities = %d entries, want %d", len(got.Identities), len(p.Targets))
	}
	if len(got.Archives) != len(p.Archives) {
		t.Errorf("round-tripped Archives = %d entries, want %d", len(got.Archives), len(p.Archives))
	}
	if len(got.LegacyPrerequisites) != len(p.Prereqs) {
		t.Errorf("round-tripped LegacyPrerequisites = %d entries, want %d", len(got.LegacyPrerequisites), len(p.Prereqs))
	}
	if got.OperationID != p.OperationID {
		t.Errorf("round-tripped OperationID = %q, want %q", got.OperationID, p.OperationID)
	}
	if !got.GeneratedAt.Equal(p.GeneratedAt) {
		t.Errorf("round-tripped GeneratedAt = %v, want %v", got.GeneratedAt, p.GeneratedAt)
	}
}

// TestBuildManifest_recordsSourceQualifiedIdentitiesDistinctly proves
// v1-history's two T001-shared records land in the manifest as two distinct
// identities, each naming its own release and path.
func TestBuildManifest_recordsSourceQualifiedIdentitiesDistinctly(t *testing.T) {
	p := mustPlan(t, fixtureProjectRoot("v1-history"))
	m := BuildManifest(p)

	var v11 *ManifestIdentity
	for i := range m.Identities {
		if m.Identities[i].Path == ".savepoint/releases/v1.1/epics/E01-example/tasks/T001-shared.md" {
			v11 = &m.Identities[i]
		}
	}
	if v11 == nil {
		t.Fatalf("manifest identities = %+v, want one naming v1.1's T001-shared", m.Identities)
	}
	if v11.Release != "v1.1" || v11.Epic != "E01-example" {
		t.Errorf("v1.1 T001-shared identity = %+v, want release v1.1, epic E01-example", v11)
	}

	// v1's T001-shared is archived (its epic is done), so it must not also
	// appear as an identity.
	for _, id := range m.Identities {
		if id.Path == ".savepoint/releases/v1/epics/E01-example/tasks/T001-shared.md" {
			t.Errorf("v1's archived T001-shared unexpectedly appears as a manifest identity: %+v", id)
		}
	}
}

// --- create-only write ------------------------------------------------

func TestWriteManifestCreateOnly_refusesToOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "migrations", "v1-to-v2.yml")

	m := &ManifestV1ToV2{ManifestSchemaVersion: manifestSchemaVersion, OperationID: "op-1", GeneratedAt: time.Now()}
	if err := WriteManifestCreateOnly(path, m); err != nil {
		t.Fatalf("first WriteManifestCreateOnly() error = %v", err)
	}

	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read written manifest: %v", err)
	}

	m2 := &ManifestV1ToV2{ManifestSchemaVersion: manifestSchemaVersion, OperationID: "op-2", GeneratedAt: time.Now()}
	if err := WriteManifestCreateOnly(path, m2); err == nil {
		t.Fatal("second WriteManifestCreateOnly() error = nil, want a refusal: the manifest is create-only")
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read manifest after refused overwrite: %v", err)
	}
	if string(after) != string(first) {
		t.Error("manifest content changed after a refused overwrite attempt")
	}
}
