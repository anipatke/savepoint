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
// overall label and how many of the signals are Good, one count per signal
// however many instances it has. Counts are set only when
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
	d, err := LoadDashboard(projectPath)
	if err != nil {
		return Chip{State: ChipNotSetUp}
	}
	return d.Chip()
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
	c := Chip{State: ChipMeasured, Overall: d.Overall, Label: d.OverallText}
	worst := map[Capability]DashboardRow{}
	var order []Capability
	for _, row := range d.Rows {
		at, seen := worst[row.Capability]
		if !seen {
			order = append(order, row.Capability)
		}
		if !seen || WorseInstance(row, at) {
			worst[row.Capability] = row
		}
	}
	c.Signals = len(order)
	for _, capability := range order {
		if worst[capability].shownLabel() == ClassificationGood {
			c.Good++
		}
	}
	return c
}

// WorseInstance orders instances of one signal: blocking first, then by label
// severity. A signal with several instances stands for its worst one, on the
// chip and in the popover alike.
func WorseInstance(a, b DashboardRow) bool {
	if a.BlocksSignOff() != b.BlocksSignOff() {
		return a.BlocksSignOff()
	}
	return instanceSeverity[a.Label] > instanceSeverity[b.Label]
}

var instanceSeverity = map[Classification]int{
	ClassificationGood:           0,
	ClassificationWatch:          1,
	ClassificationUnknown:        2,
	ClassificationNeedsAttention: 3,
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
