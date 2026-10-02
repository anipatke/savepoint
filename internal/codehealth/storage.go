package codehealth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Layout of the health area below a project root.
const (
	savepointDir = ".savepoint"
	healthDir    = "health"
	snapshotsDir = "snapshots"
	configFile   = "config.json"
	snapshotExt  = ".json"
	tempPrefix   = ".tmp-"
	tempSuffix   = ".tmp"
)

// ManualRetention is how many manual snapshots explicit pruning keeps.
const ManualRetention = 10

// maxRecordBytes bounds what a loader will read, so a bad file cannot exhaust
// memory.
const maxRecordBytes = 1 << 20

// Storage errors. Callers match them with errors.Is.
var (
	ErrNotProject       = errors.New("not a Savepoint project")
	ErrUnsafePath       = errors.New("unsafe health path")
	ErrConfigNotFound   = errors.New("health configuration not found")
	ErrSnapshotConflict = errors.New("a different snapshot already has this identity")
)

// Store reads and writes the health area of one project. It holds only the
// project path, so nothing depends on the working directory or on state that
// outlives a call.
type Store struct{ projectPath string }

// NewStore returns a Store rooted at projectPath.
func NewStore(projectPath string) Store { return Store{projectPath: projectPath} }

// dir resolves a directory under .savepoint/ without following symlinks. It
// returns ok=false when an existing path is absent and create is false. With
// create, it makes only the health directories; it never makes .savepoint, so
// a directory that is not a project is refused rather than scaffolded.
func (s Store) dir(create bool, elems ...string) (path string, ok bool, err error) {
	path = s.projectPath
	for i, e := range append([]string{savepointDir}, elems...) {
		path = filepath.Join(path, e)
		info, err := os.Lstat(path)
		if errors.Is(err, fs.ErrNotExist) {
			if !create {
				return "", false, nil
			}
			if i == 0 {
				return "", false, fmt.Errorf("%w: %s has no %s directory", ErrNotProject, s.projectPath, savepointDir)
			}
			if err := os.Mkdir(path, 0o755); err != nil && !errors.Is(err, fs.ErrExist) {
				return "", false, err
			}
			if info, err = os.Lstat(path); err != nil {
				return "", false, err
			}
		} else if err != nil {
			return "", false, err
		}
		if !info.IsDir() {
			return "", false, fmt.Errorf("%w: %s is a symlink or not a directory", ErrUnsafePath, path)
		}
	}
	return path, true, nil
}

// readRecord reads a bounded regular file. exists is false when it is absent.
func readRecord(path string) (data []byte, exists bool, err error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if !info.Mode().IsRegular() {
		return nil, true, fmt.Errorf("%w: %s is not a regular file", ErrUnsafePath, path)
	}
	if info.Size() > maxRecordBytes {
		return nil, true, fieldError(ErrUnboundedDetail, filepath.Base(path), "is %d bytes; at most %d", info.Size(), maxRecordBytes)
	}
	data, err = os.ReadFile(path)
	return data, true, err
}

func encodeRecord(v any) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	return append(b, '\n'), err
}

// writeTemp writes data to a fresh temporary file beside its destination, so a
// later rename or link is atomic. The caller removes the returned path.
func writeTemp(dir string, data []byte) (string, error) {
	f, err := os.CreateTemp(dir, tempPrefix+"*"+tempSuffix)
	if err != nil {
		return "", err
	}
	name := f.Name()
	_, werr := f.Write(data)
	if werr == nil {
		werr = f.Sync()
	}
	if cerr := f.Close(); werr == nil {
		werr = cerr
	}
	if werr == nil {
		werr = os.Chmod(name, 0o644)
	}
	if werr != nil {
		os.Remove(name)
		return "", werr
	}
	return name, nil
}

// LoadConfig reads the configuration. It returns ErrConfigNotFound when none
// exists, leaving the caller to decide what an absent configuration means.
func (s Store) LoadConfig() (Config, error) {
	dir, ok, err := s.dir(false, healthDir)
	if err != nil {
		return Config{}, err
	}
	if !ok {
		return Config{}, ErrConfigNotFound
	}
	data, exists, err := readRecord(filepath.Join(dir, configFile))
	if err != nil {
		return Config{}, err
	}
	if !exists {
		return Config{}, ErrConfigNotFound
	}
	c, err := DecodeConfig(data)
	if err != nil {
		return Config{}, fmt.Errorf("%s: %w", configFile, err)
	}
	return c, nil
}

// SaveConfig validates and stores the configuration, replacing any existing
// file atomically. It reports whether the file changed; identical content is
// left untouched.
func (s Store) SaveConfig(c Config) (changed bool, err error) {
	if err := c.Validate(); err != nil {
		return false, err
	}
	data, err := encodeRecord(c)
	if err != nil {
		return false, err
	}
	dir, _, err := s.dir(true, healthDir)
	if err != nil {
		return false, err
	}
	path := filepath.Join(dir, configFile)
	old, exists, err := readRecord(path)
	if err != nil {
		return false, err
	}
	if exists && string(old) == string(data) {
		return false, nil
	}
	tmp, err := writeTemp(dir, data)
	if err != nil {
		return false, err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return false, err
	}
	return true, nil
}

// snapshotFileName is the only name a snapshot may have: its identity digest.
func snapshotFileName(id string) string {
	return strings.TrimPrefix(id, digestPrefix) + snapshotExt
}

func snapshotIDFromFileName(name string) (string, bool) {
	raw, ok := strings.CutSuffix(name, snapshotExt)
	if !ok {
		return "", false
	}
	id := digestPrefix + raw
	return id, validateDigest("name", id) == nil
}

func isTempName(name string) bool {
	return strings.HasPrefix(name, tempPrefix) && strings.HasSuffix(name, tempSuffix)
}

// LoadSnapshots returns every stored snapshot, oldest first, ties broken by
// identity. Snapshots from different comparison series are all returned; use
// SeriesID to tell them apart. Any unreadable, unsupported, or inconsistent
// file fails the load, and nothing is repaired or rewritten. Leftovers of an
// interrupted write are ignored.
func (s Store) LoadSnapshots() ([]Snapshot, error) {
	dir, ok, err := s.dir(false, healthDir, snapshotsDir)
	if err != nil || !ok {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []Snapshot
	seen := make(map[string]string, len(entries))
	for _, e := range entries {
		name := e.Name()
		if isTempName(name) {
			continue
		}
		id, ok := snapshotIDFromFileName(name)
		if !ok {
			return nil, fieldError(ErrMalformedRecord, name, "is not a snapshot file name")
		}
		data, _, err := readRecord(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		snap, err := DecodeSnapshot(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if snap.ID != id {
			return nil, fieldError(ErrIdentityMismatch, name, "file name says %s, content is %s", id, snap.ID)
		}
		if first, dup := seen[id]; dup {
			return nil, fieldError(ErrDuplicateInstance, name, "repeats snapshot %s already in %s", id, first)
		}
		seen[id] = name
		out = append(out, snap)
	}
	slices.SortFunc(out, compareSnapshots)
	return out, nil
}

func compareSnapshots(a, b Snapshot) int {
	if c := strings.Compare(a.CreatedAt, b.CreatedAt); c != 0 {
		return c
	}
	return strings.Compare(a.ID, b.ID)
}

// SaveSnapshot stores a snapshot create-only and atomically, so it is safe to
// retry. It reports whether a file was created: an existing identical snapshot
// is left unchanged, and a file at the same identity with different content is
// refused with ErrSnapshotConflict.
func (s Store) SaveSnapshot(snap Snapshot) (created bool, err error) {
	if err := snap.Validate(); err != nil {
		return false, err
	}
	data, err := encodeRecord(snap.Canonical())
	if err != nil {
		return false, err
	}
	dir, _, err := s.dir(true, healthDir, snapshotsDir)
	if err != nil {
		return false, err
	}
	final := filepath.Join(dir, snapshotFileName(snap.ID))
	if exists, err := sameSnapshot(final, snap.ID); err != nil || exists {
		return false, err
	}
	tmp, err := writeTemp(dir, data)
	if err != nil {
		return false, err
	}
	defer os.Remove(tmp)
	// Link fails if the name exists, unlike rename, which would replace it.
	if err := os.Link(tmp, final); err != nil {
		// Lost a race with another writer: identical content is still success.
		if errors.Is(err, fs.ErrExist) {
			if exists, serr := sameSnapshot(final, snap.ID); serr != nil || exists {
				return false, serr
			}
		}
		return false, err
	}
	return true, nil
}

// sameSnapshot reports whether path already holds the snapshot with this
// identity. Anything else at that path is a conflict.
func sameSnapshot(path, id string) (bool, error) {
	data, exists, err := readRecord(path)
	if err != nil || !exists {
		return false, err
	}
	existing, err := DecodeSnapshot(data)
	if err != nil || existing.ID != id {
		return false, fmt.Errorf("%w: %s: existing file is not snapshot %s (%v)", ErrSnapshotConflict, filepath.Base(path), id, err)
	}
	return true, nil
}

// PruneReport lists snapshot identities in load order. Removed is what a plan
// would remove or what Prune actually removed; Retained is everything kept.
type PruneReport struct {
	Retained []string
	Removed  []string
}

// PlanPrune reports what Prune would remove without touching any file. Only
// manual snapshots are ever candidates: the newest ManualRetention stay and
// official snapshots are never listed for removal.
func (s Store) PlanPrune() (PruneReport, error) {
	snaps, err := s.LoadSnapshots()
	if err != nil {
		return PruneReport{}, err
	}
	return planPrune(snaps), nil
}

func planPrune(snaps []Snapshot) PruneReport {
	var manual []string
	for _, sn := range snaps {
		if sn.Origin == OriginManual {
			manual = append(manual, sn.ID)
		}
	}
	drop := map[string]bool{}
	if over := len(manual) - ManualRetention; over > 0 {
		for _, id := range manual[:over] {
			drop[id] = true
		}
	}
	var r PruneReport
	for _, sn := range snaps {
		if drop[sn.ID] {
			r.Removed = append(r.Removed, sn.ID)
		} else {
			r.Retained = append(r.Retained, sn.ID)
		}
	}
	return r
}

// Prune is the explicit maintenance operation: it removes the manual snapshots
// PlanPrune lists, oldest first. It refuses to run if any history file is
// invalid. On a removal failure it stops and reports what was removed so far;
// every remaining file is untouched and still valid. Running it again is safe.
func (s Store) Prune() (PruneReport, error) {
	snaps, err := s.LoadSnapshots()
	if err != nil {
		return PruneReport{}, err
	}
	plan := planPrune(snaps)
	if len(plan.Removed) == 0 {
		return plan, nil
	}
	dir, ok, err := s.dir(false, healthDir, snapshotsDir)
	if err != nil || !ok {
		return PruneReport{Retained: plan.Retained}, errors.Join(err, errors.New("snapshot directory vanished"))
	}
	report := PruneReport{Retained: plan.Retained}
	for _, id := range plan.Removed {
		err := os.Remove(filepath.Join(dir, snapshotFileName(id)))
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			report.Retained = append(report.Retained, plan.Removed[len(report.Removed):]...)
			return report, err
		}
		report.Removed = append(report.Removed, id)
	}
	return report, nil
}
