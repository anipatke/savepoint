package init

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ManifestVersion is the schema version written into new manifests. It exists
// so a later shape change can be detected rather than mis-parsed.
const ManifestVersion = 1

// manifestRelPath is the manifest location inside a project, relative to the
// project root. It is deliberately outside the template tree: upgrade skips
// .savepoint/, so the manifest is never itself an upgradable asset.
var manifestRelPath = filepath.Join(".savepoint", ".upgrade-manifest.yml")

// Manifest records the SHA-256 of each skill file as Savepoint last wrote it.
// Comparing an on-disk file to these hashes is what distinguishes a user edit
// from a merely outdated copy; template comparison alone cannot.
type Manifest struct {
	Version int               `yaml:"version"`
	Skills  map[string]string `yaml:"skills"`
}

// NewManifest returns an empty manifest at the current schema version.
func NewManifest() *Manifest {
	return &Manifest{Version: ManifestVersion, Skills: map[string]string{}}
}

// LoadManifest reads the manifest for the project rooted at dir. A missing
// manifest is not an error: a project from before the manifest existed simply
// has no recorded provenance yet.
func LoadManifest(dir string) (*Manifest, error) {
	path := filepath.Join(dir, manifestRelPath)

	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return NewManifest(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", path, err)
	}

	var manifest Manifest
	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %w", path, err)
	}

	if manifest.Version == 0 {
		manifest.Version = ManifestVersion
	}
	if manifest.Skills == nil {
		manifest.Skills = map[string]string{}
	}
	return &manifest, nil
}

// Record stores the hash of content for path, which must be a template-relative
// path. Paths outside the manifest scope are ignored, so callers may record
// every asset they write without repeating the scope rule.
func (m *Manifest) Record(path string, content []byte) {
	key := filepath.ToSlash(path)
	if !isManifestPath(key) {
		return
	}
	if m.Skills == nil {
		m.Skills = map[string]string{}
	}
	m.Skills[key] = hashContent(content)
}

// Hash returns the recorded hash for path and whether the path is recorded.
func (m *Manifest) Hash(path string) (string, bool) {
	hash, ok := m.Skills[filepath.ToSlash(path)]
	return hash, ok
}

// Save writes the manifest for the project rooted at dir.
func (m *Manifest) Save(dir string) error {
	return m.save(dir, AtomicWrite)
}

// save writes the manifest through write. An upgrade that changed nothing must
// change nothing on disk, so identical content returns without writing: the
// manifest keeps its file identity and modification time.
func (m *Manifest) save(dir string, write assetWriter) error {
	path := filepath.Join(dir, manifestRelPath)

	if m.Version == 0 {
		m.Version = ManifestVersion
	}
	if m.Skills == nil {
		m.Skills = map[string]string{}
	}

	content, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("encode manifest %s: %w", path, err)
	}

	existing, readErr := os.ReadFile(path)
	switch {
	case readErr == nil && bytes.Equal(existing, content):
		return nil
	case readErr != nil && !os.IsNotExist(readErr):
		return fmt.Errorf("read existing manifest %s: %w", path, readErr)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create parent dir for %s: %w", path, err)
	}
	if err := write(path, content); err != nil {
		return fmt.Errorf("write manifest %s: %w", path, err)
	}
	return nil
}

// isManifestPath reports whether a template-relative slash path is covered by
// the manifest: skill entrypoints only. AGENTS.md carries its own ownership
// signal in the marker pair, shared references under agent-skills/references/
// stay package-owned, and no .savepoint/ path is wholesale-owned by Savepoint.
func isManifestPath(path string) bool {
	parts := strings.Split(path, "/")
	return len(parts) == 3 && parts[0] == "agent-skills" && parts[2] == "SKILL.md"
}

// hashContent hashes exact file bytes, with no line-ending normalization, so a
// recorded hash matches the file Savepoint wrote byte for byte.
func hashContent(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}
