package codehealth

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ChipState says which of the three header chips the board shows.
type ChipState string

const (
	ChipNotSetUp ChipState = "not_set_up"
	ChipNoCheck  ChipState = "no_check"
	ChipMeasured ChipState = "measured"
)

// Chip is the small glance at overall health: the newest saved snapshot's
// overall label and how many of the signals are Good. Counts are set only when
// State is ChipMeasured.
type Chip struct {
	State   ChipState
	Overall Classification
	Label   string
	Good    int
	Signals int
}

// LoadChip reads the saved configuration and the newest snapshot of either
// origin, the one the dashboard describes. It runs no tool and no Git command.
// Storage that cannot be read reads as not set up: a glance must not fail the
// board, and doctor reports the damage.
func LoadChip(projectPath string) Chip {
	store := NewStore(projectPath)
	cfg, err := store.LoadConfig()
	if err != nil {
		return Chip{State: ChipNotSetUp}
	}
	snap, found, err := store.LatestSnapshot()
	if err != nil {
		return Chip{State: ChipNotSetUp}
	}
	if !found {
		return Chip{State: ChipNoCheck}
	}
	good := 0
	for _, cs := range snap.Summary.Capabilities {
		if cs.Classification == ClassificationGood {
			good++
		}
	}
	return Chip{
		State:   ChipMeasured,
		Overall: snap.Summary.Overall,
		Label:   classificationText[snap.Summary.Overall],
		Good:    good,
		Signals: signalCount(cfg),
	}
}

// Chip describes a loaded dashboard the same way LoadChip describes storage, so
// a screen that already holds the dashboard needs no second read.
func (d Dashboard) Chip() Chip {
	switch d.State {
	case DashboardNotConfigured:
		return Chip{State: ChipNotSetUp}
	case DashboardFirstRun:
		return Chip{State: ChipNoCheck}
	}
	c := Chip{State: ChipMeasured, Overall: d.Overall, Label: d.OverallText, Signals: len(d.Rows)}
	for _, row := range d.Rows {
		if row.Label == ClassificationGood {
			c.Good++
		}
	}
	return c
}

// signalCount is the number of rows the dashboard lists: each configured
// instance, or one placeholder for a capability that has none.
func signalCount(cfg Config) int {
	n := 0
	for _, c := range Capabilities() {
		instances := 0
		for _, cc := range cfg.Capabilities {
			if cc.Capability == c {
				instances++
			}
		}
		n += max(instances, 1)
	}
	return n
}

// LatestSnapshot returns the newest stored snapshot by the order LoadSnapshots
// uses: stored creation time, then identity. File times are not consulted, so a
// restored or copied history keeps its order. Each file is read only for its
// creation time; only the newest is fully decoded and checked against its
// name, and a file that cannot be read for its time is left to doctor.
func (s Store) LatestSnapshot() (snap Snapshot, found bool, err error) {
	dir, ok, err := s.dir(false, healthDir, snapshotsDir)
	if err != nil || !ok {
		return Snapshot{}, false, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Snapshot{}, false, err
	}
	var newest, id, newestAt string
	for _, e := range entries {
		name := e.Name()
		if isTempName(name) {
			continue
		}
		nameID, ok := snapshotIDFromFileName(name)
		if !ok {
			return Snapshot{}, false, fieldError(ErrMalformedRecord, name, "is not a snapshot file name")
		}
		data, exists, err := readRecord(filepath.Join(dir, name))
		if err != nil {
			return Snapshot{}, false, err
		}
		var stamp struct {
			CreatedAt string `json:"created_at"`
		}
		if !exists || json.Unmarshal(data, &stamp) != nil {
			continue
		}
		if newest == "" || stamp.CreatedAt > newestAt || (stamp.CreatedAt == newestAt && nameID > id) {
			newest, id, newestAt = name, nameID, stamp.CreatedAt
		}
	}
	if newest == "" {
		return Snapshot{}, false, nil
	}
	data, _, err := readRecord(filepath.Join(dir, newest))
	if err != nil {
		return Snapshot{}, false, err
	}
	snap, err = DecodeSnapshot(data)
	if err != nil {
		return Snapshot{}, false, err
	}
	if snap.ID != id {
		return Snapshot{}, false, fieldError(ErrIdentityMismatch, newest, "file name says %s, content is %s", id, snap.ID)
	}
	return snap, true, nil
}
