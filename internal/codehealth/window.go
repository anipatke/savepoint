package codehealth

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// windowOfficial is how many of the newest official snapshots a bounded load
// keeps: enough for the longest reader of history, the sparkline.
const windowOfficial = MaxSparkPoints

// snapshotHead is what the window lookup knows about a file before decoding it.
type snapshotHead struct {
	name      string
	id        string
	createdAt string
	origin    Origin
}

// LoadWindow returns the snapshots a dashboard needs, oldest first, in the order
// LoadSnapshots uses: the newest windowOfficial official snapshots, every
// snapshot newer than the oldest of them, and the newest snapshot. Each file is
// read only for its creation time and origin, by the same JSON rules the full
// decoder uses (a repeated field keeps its last value); just the window is
// validated and checked against its name. cut is true when older snapshots were
// left unread.
//
// A name that is not a snapshot name, a file whose time and origin cannot be
// read, and any damaged snapshot in the window fail the load. A damaged body
// outside the window is not seen here; Collect, Prune and LoadSnapshots still
// validate every file. Nothing is repaired or rewritten.
func (s Store) LoadWindow() (snaps []Snapshot, cut bool, err error) {
	dir, ok, err := s.dir(false, healthDir, snapshotsDir)
	if err != nil || !ok {
		return nil, false, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, false, err
	}
	heads := make([]snapshotHead, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if isTempName(name) {
			continue
		}
		id, ok := snapshotIDFromFileName(name)
		if !ok {
			return nil, false, fieldError(ErrMalformedRecord, name, "is not a snapshot file name")
		}
		h, err := readHead(dir, name, id)
		if err != nil {
			return nil, false, err
		}
		heads = append(heads, h)
	}
	slices.SortFunc(heads, func(a, b snapshotHead) int {
		if c := strings.Compare(a.createdAt, b.createdAt); c != 0 {
			return c
		}
		return strings.Compare(a.id, b.id)
	})
	first := windowStart(heads)
	for _, h := range heads[first:] {
		snap, err := decodeNamed(dir, h.name, h.id)
		if err != nil {
			return nil, false, err
		}
		snaps = append(snaps, snap)
	}
	return snaps, first > 0, nil
}

// windowStart is the index of the oldest head in the window: that of the
// windowOfficial-th newest official snapshot, or 0 when there are fewer.
func windowStart(heads []snapshotHead) int {
	seen := 0
	for i := len(heads) - 1; i >= 0; i-- {
		if heads[i].origin != OriginOfficial {
			continue
		}
		if seen++; seen < windowOfficial {
			continue
		}
		// With no official snapshot before this one, the older files are only
		// manual: nothing is beyond the window, so the whole history loads.
		for _, h := range heads[:i] {
			if h.origin == OriginOfficial {
				return i
			}
		}
		return 0
	}
	return 0
}

// decodeNamed reads, decodes and validates one snapshot file and checks that its
// content has the identity its name claims.
func decodeNamed(dir, name, id string) (Snapshot, error) {
	data, _, err := readRecord(filepath.Join(dir, name))
	if err != nil {
		return Snapshot{}, err
	}
	snap, err := DecodeSnapshot(data)
	if err != nil {
		return Snapshot{}, fmt.Errorf("%s: %w", name, err)
	}
	if snap.ID != id {
		return Snapshot{}, fieldError(ErrIdentityMismatch, name, "file name says %s, content is %s", id, snap.ID)
	}
	return snap, nil
}

// readHead learns a file's creation time and origin without validating the
// rest of it. A file that is not a bounded JSON object with a valid time and
// origin fails the load.
func readHead(dir, name, id string) (snapshotHead, error) {
	data, exists, err := readRecord(filepath.Join(dir, name))
	if err != nil {
		return snapshotHead{}, err
	}
	if !exists {
		return snapshotHead{}, fmt.Errorf("%s: %w", name, fs.ErrNotExist)
	}
	var f struct {
		CreatedAt string `json:"created_at"`
		Origin    Origin `json:"origin"`
	}
	if err := json.Unmarshal(data, &f); err != nil {
		return snapshotHead{}, fieldError(ErrMalformedRecord, name, "is not readable JSON: %v", err)
	}
	if err := validateTimestamp("created_at", f.CreatedAt); err != nil {
		return snapshotHead{}, fmt.Errorf("%s: %w", name, err)
	}
	if f.Origin != OriginOfficial && f.Origin != OriginManual {
		return snapshotHead{}, fieldError(ErrMalformedRecord, name, "has origin %q", f.Origin)
	}
	return snapshotHead{name: name, id: id, createdAt: f.CreatedAt, origin: f.Origin}, nil
}
