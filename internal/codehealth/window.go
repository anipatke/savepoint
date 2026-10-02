package codehealth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// headBytes is how much of a snapshot file the window lookup reads to learn its
// time and origin. A canonical file keeps both within its first 200 bytes.
const headBytes = 512

// windowOfficial is how many of the newest official snapshots a bounded load
// keeps: enough for the longest reader of history, the sparkline.
const windowOfficial = MaxSparkPoints

// snapshotHead is what the window lookup knows about a file before decoding it.
type snapshotHead struct {
	name      string
	id        string
	createdAt string
	origin    Origin
	// snap is set only when the head could not be read cheaply and the whole file
	// was decoded instead.
	snap *Snapshot
}

// LoadWindow returns the snapshots a dashboard needs, oldest first, in the order
// LoadSnapshots uses: the newest windowOfficial official snapshots, every
// snapshot newer than the oldest of them, and the newest snapshot. Each file is
// listed and read only for its creation time and origin; just the window is
// decoded and checked against its name. cut is true when older snapshots were
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
		snap, err := loadHeadSnapshot(dir, h)
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

// loadHeadSnapshot decodes the file behind h exactly as LoadSnapshots does.
func loadHeadSnapshot(dir string, h snapshotHead) (Snapshot, error) {
	if h.snap != nil {
		return *h.snap, nil
	}
	return decodeNamed(dir, h.name, h.id)
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

// readHead learns a file's creation time and origin from its first bytes. A
// head that does not give both (another field order, a long first field) falls
// back to decoding the whole file, which fails the load if the file is damaged.
func readHead(dir, name, id string) (snapshotHead, error) {
	path := filepath.Join(dir, name)
	head, err := readPrefix(path)
	if err != nil {
		return snapshotHead{}, err
	}
	if createdAt, origin, ok := parseHead(head); ok {
		return snapshotHead{name: name, id: id, createdAt: createdAt, origin: origin}, nil
	}
	snap, err := decodeNamed(dir, name, id)
	if err != nil {
		return snapshotHead{}, err
	}
	return snapshotHead{name: name, id: id, createdAt: snap.CreatedAt, origin: snap.Origin, snap: &snap}, nil
}

// readPrefix reads at most headBytes of a regular file under the same size and
// type limits readRecord applies.
func readPrefix(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%w: %s is not a regular file", ErrUnsafePath, path)
	}
	if info.Size() > maxRecordBytes {
		return nil, fieldError(ErrUnboundedDetail, filepath.Base(path), "is %d bytes; at most %d", info.Size(), maxRecordBytes)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]byte, headBytes)
	n, err := io.ReadFull(f, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, err
	}
	return buf[:n], nil
}

// parseHead reads the top-level fields of a snapshot's first bytes for
// created_at and origin. ok is false when either is missing, is not a valid time
// or origin, or appears more than once in the head: the full decoder keeps the
// last of a repeated field, so the caller decodes the file instead. A value cut
// off by the end of the head ends the scan.
func parseHead(head []byte) (createdAt string, origin Origin, ok bool) {
	dec := json.NewDecoder(bytes.NewReader(head))
	if t, err := dec.Token(); err != nil || t != json.Delim('{') {
		return "", "", false
	}
	var seenTime, seenOrigin bool
	for {
		key, err := dec.Token()
		if err != nil {
			break
		}
		k, isString := key.(string)
		if !isString {
			break // the closing brace of a short file
		}
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			if (k == "created_at" && seenTime) || (k == "origin" && seenOrigin) {
				return "", "", false
			}
			break
		}
		var v string
		isText := json.Unmarshal(raw, &v) == nil
		switch k {
		case "created_at":
			if seenTime || !isText {
				return "", "", false
			}
			seenTime, createdAt = true, v
		case "origin":
			if seenOrigin || !isText {
				return "", "", false
			}
			seenOrigin, origin = true, Origin(v)
		}
	}
	if !seenTime || !seenOrigin || validateTimestamp("created_at", createdAt) != nil || (origin != OriginOfficial && origin != OriginManual) {
		return "", "", false
	}
	return createdAt, origin, true
}
