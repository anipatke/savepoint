package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

// ManifestV1ToV2 is the model behind .savepoint/migrations/v1-to-v2.yml: the
// source hashes, the legacy-to-global identity map, the archive references,
// the typed legacy prerequisites active work
// depends on, and the owner decisions that resolved any ambiguity. It is a
// pure projection of a ConversionPlan plus the inventory Plan already read;
// building one never touches the filesystem, and writing one is create-only
// (see WriteManifestCreateOnly).
type ManifestV1ToV2 struct {
	// ManifestSchemaVersion is this manifest file's own shape version. It is
	// independent of the project's config.yml schema_version and must never
	// be confused with it (see data.ReadSchemaVersion's doc comment).
	ManifestSchemaVersion int                          `yaml:"manifest_schema_version"`
	Sources               []ManifestSource             `yaml:"sources"`
	Identities            []ManifestIdentity           `yaml:"identities"`
	Archives              []ManifestArchive            `yaml:"archives"`
	LegacyPrerequisites   []ManifestLegacyPrerequisite `yaml:"legacy_prerequisites,omitempty"`
	WaivedReferences      []ManifestWaivedReference    `yaml:"waived_references,omitempty"`
	Decisions             []ManifestDecision           `yaml:"decisions,omitempty"`
}

const manifestSchemaVersion = 1

// ManifestSource is one inventoried V1 file's exact source hash, retained as
// provenance alongside the identity and archive mappings.
type ManifestSource struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
}

// ManifestIdentity is one legacy-to-global identity mapping: which
// source-qualified V1 record became which new V2 record, and where.
type ManifestIdentity struct {
	GlobalID   string `yaml:"global_id"`
	Kind       string `yaml:"kind"`
	Release    string `yaml:"release,omitempty"`
	ReleaseID  string `yaml:"release_id,omitempty"`
	Epic       string `yaml:"epic,omitempty"`
	Path       string `yaml:"path"`
	OriginalID string `yaml:"original_id"`
	TargetPath string `yaml:"target_path"`
}

// ManifestArchive maps one V1 source to its byte-preserved archive copy. A
// converted record can have both an identity mapping and an archive mapping.
type ManifestArchive struct {
	SourcePath  string `yaml:"source_path"`
	ArchivePath string `yaml:"archive_path"`
	SHA256      string `yaml:"sha256,omitempty"`
	Role        string `yaml:"role,omitempty"`
	Release     string `yaml:"release,omitempty"`
	Epic        string `yaml:"epic,omitempty"`
	OriginalID  string `yaml:"original_id,omitempty"`
}

// ManifestLegacyPrerequisite is one typed record of an active Task's
// dependency on archived, completed V1 work with no V2 identity.
type ManifestLegacyPrerequisite struct {
	Task        string `yaml:"task"`
	ArchivePath string `yaml:"archive_path"`
	Evidence    string `yaml:"evidence"`
}

// ManifestWaivedReference is one typed record of an active converted Task
// whose V1 source named a waived audit finding. A waiver is a closed owner
// decision, never a fabricated Issue; this is what keeps the fact resolvable
// instead of silently dropped.
type ManifestWaivedReference struct {
	Task        string `yaml:"task"`
	ArchivePath string `yaml:"archive_path"`
	Reason      string `yaml:"reason"`
}

// ManifestDecision is one owner-supplied resolution recorded against a named
// ambiguity, with its provenance: the decision value itself, the
// --decisions FILE path it came from, and when it was read.
type ManifestDecision struct {
	AmbiguityID string    `yaml:"ambiguity_id"`
	Decision    string    `yaml:"decision"`
	SourceFile  string    `yaml:"source_file,omitempty"`
	DecidedAt   time.Time `yaml:"decided_at,omitempty"`
}

// BuildManifest projects plan into the manifest model. It is a pure
// transform: every value it needs already lives on plan (Sources was
// recorded by Plan's own inventory read), so building a manifest never reads
// the filesystem again.
func BuildManifest(plan *ConversionPlan) *ManifestV1ToV2 {
	m := &ManifestV1ToV2{
		ManifestSchemaVersion: manifestSchemaVersion,
	}

	for _, s := range plan.Sources {
		m.Sources = append(m.Sources, ManifestSource{Path: s.Path, SHA256: s.SHA256})
	}

	for _, t := range plan.Targets {
		if t.Generated {
			continue // generated Goals have no V1 source identity to map
		}
		releaseID := ""
		if t.Kind == TargetObjective {
			releaseID = t.ReleaseID
		}
		m.Identities = append(m.Identities, ManifestIdentity{
			GlobalID:   t.GlobalID,
			Kind:       string(t.Kind),
			Release:    t.Legacy.Release,
			ReleaseID:  releaseID,
			Epic:       t.Legacy.Epic,
			Path:       t.Legacy.Path,
			OriginalID: t.Legacy.OriginalID,
			TargetPath: t.InstallPath(),
		})
	}
	sort.Slice(m.Identities, func(i, j int) bool { return m.Identities[i].GlobalID < m.Identities[j].GlobalID })

	for _, a := range plan.Archives {
		entry := ManifestArchive{SourcePath: a.SourcePath, ArchivePath: a.ArchivePath, SHA256: a.SourceSHA256, Role: string(a.Role)}
		if a.Legacy != nil {
			entry.Release = a.Legacy.Release
			entry.Epic = a.Legacy.Epic
			entry.OriginalID = a.Legacy.OriginalID
		}
		m.Archives = append(m.Archives, entry)
	}
	sort.Slice(m.Archives, func(i, j int) bool { return m.Archives[i].SourcePath < m.Archives[j].SourcePath })

	for _, p := range plan.Prereqs {
		m.LegacyPrerequisites = append(m.LegacyPrerequisites, ManifestLegacyPrerequisite{
			Task:        p.Task,
			ArchivePath: p.ArchivePath,
			Evidence:    p.Evidence,
		})
	}
	sort.Slice(m.LegacyPrerequisites, func(i, j int) bool {
		return m.LegacyPrerequisites[i].Task < m.LegacyPrerequisites[j].Task
	})

	for _, w := range plan.WaivedRefs {
		m.WaivedReferences = append(m.WaivedReferences, ManifestWaivedReference{
			Task:        w.Task,
			ArchivePath: w.ArchivePath,
			Reason:      w.Reason,
		})
	}
	sort.Slice(m.WaivedReferences, func(i, j int) bool {
		return m.WaivedReferences[i].Task < m.WaivedReferences[j].Task
	})

	for _, amb := range plan.Ambiguities {
		if !amb.Resolved {
			continue
		}
		m.Decisions = append(m.Decisions, ManifestDecision{
			AmbiguityID: amb.ID,
			Decision:    amb.Decision,
			SourceFile:  amb.DecisionSourceFile,
			DecidedAt:   amb.DecisionAt,
		})
	}
	sort.Slice(m.Decisions, func(i, j int) bool { return m.Decisions[i].AmbiguityID < m.Decisions[j].AmbiguityID })

	return m
}

// Marshal renders the manifest as YAML.
func (m *ManifestV1ToV2) Marshal() ([]byte, error) {
	return yaml.Marshal(m)
}

// UnmarshalManifest parses the provenance recorded in v1-to-v2.yml.
func UnmarshalManifest(content []byte) (*ManifestV1ToV2, error) {
	var m ManifestV1ToV2
	if err := yaml.Unmarshal(content, &m); err != nil {
		return nil, fmt.Errorf("parse v1-to-v2.yml: %w", err)
	}
	return &m, nil
}

// WriteManifestCreateOnly writes the manifest to path, refusing outright if
// anything already exists there. This is the mechanical half of "never
// overwrite a prior backup or user sidecar": the manifest is written exactly
// once, by the apply step that activates the migration, and every later
// write attempt against the same path is a caller error, not a silent
// replace.
func WriteManifestCreateOnly(path string, m *ManifestV1ToV2) error {
	content, err := m.Marshal()
	if err != nil {
		return fmt.Errorf("marshal v1-to-v2.yml: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
