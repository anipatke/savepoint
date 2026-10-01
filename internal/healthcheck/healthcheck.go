// Package healthcheck runs one official Code Health collection for an
// Objective: it validates the Objective against the strictly loaded V2 index,
// collects through codehealth, and renders the plain blocking verdict.
package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/opencode/savepoint/internal/codehealth"
	"github.com/opencode/savepoint/internal/data"
)

// Request is one official collection. Runner and Git default to the real
// process runners; tests inject fakes. There is deliberately no origin field:
// this path only ever produces official snapshots.
type Request struct {
	Dir       string
	Objective string
	Runner    codehealth.ToolRunner
	Git       codehealth.CommandRunner
	// Stderr receives a note when the report cannot be written after the
	// snapshot was saved. Nil discards it.
	Stderr io.Writer
}

// Run collects one official snapshot and writes the Objective, snapshot ID,
// whether it was newly created, and the rendered verdict to stdout. It returns
// an error only when no snapshot was saved: a blocking verdict is still a
// successful collection, and so is a report that cannot be written afterwards;
// that is noted on req.Stderr instead. Without health configuration it says so and saves
// nothing.
func Run(ctx context.Context, req Request, stdout io.Writer) error {
	root, err := data.ResolveTarget(req.Dir)
	switch {
	case errors.Is(err, data.ErrTargetMissing):
		return fmt.Errorf("health check: target directory does not exist: %s", req.Dir)
	case errors.Is(err, data.ErrTargetNotSavepoint):
		return fmt.Errorf("health check: target directory is not a Savepoint project: %s has no .savepoint directory", req.Dir)
	case err != nil:
		return fmt.Errorf("health check: %w", err)
	}
	if err := data.CheckRuntimeSchema(root); err != nil {
		return fmt.Errorf("health check: %w", err)
	}
	index, err := data.LoadV2Index(filepath.Join(root, ".savepoint"))
	if err != nil {
		return fmt.Errorf("health check: %w", err)
	}
	if _, ok := index.Objectives[req.Objective]; !ok {
		return fmt.Errorf("health check: objective %s does not exist in this project", req.Objective)
	}

	store := codehealth.NewStore(root)
	cfg, err := store.LoadConfig()
	if errors.Is(err, codehealth.ErrConfigNotFound) {
		_, err = fmt.Fprintln(stdout, "Code Health is not configured for this project; nothing was collected. Run `savepoint health setup` to configure it.")
		return err
	}
	if err != nil {
		return fmt.Errorf("health check: %w", err)
	}

	collection, err := codehealth.Collect(ctx, codehealth.CollectRequest{
		Root:    root,
		Origin:  codehealth.OriginOfficial,
		Config:  cfg,
		Readers: codehealth.DefaultReaders(),
		Runner:  req.Runner,
		Git:     req.Git,
	})
	if err != nil {
		return fmt.Errorf("health check: no snapshot was saved: %w", err)
	}

	snapshots, err := store.LoadSnapshots()
	if err != nil {
		return fmt.Errorf("health check: snapshot %s was saved but cannot be read back: %w", collection.SnapshotID, err)
	}
	var verdict codehealth.Verdict
	found := false
	for _, snapshot := range snapshots {
		if snapshot.ID != collection.SnapshotID {
			continue
		}
		if verdict, err = codehealth.Evaluate(snapshot, cfg); err != nil {
			return fmt.Errorf("health check: snapshot %s was saved but cannot be evaluated: %w", collection.SnapshotID, err)
		}
		found = true
	}
	if !found {
		return fmt.Errorf("health check: snapshot %s was saved but is missing from the store", collection.SnapshotID)
	}

	stored := "created"
	if !collection.Created {
		stored = "already stored"
	}
	if _, err := fmt.Fprintf(stdout, "Objective: %s\nSnapshot: %s (%s)\n%s", req.Objective, collection.SnapshotID, stored, verdict.Render()); err != nil && req.Stderr != nil {
		fmt.Fprintf(req.Stderr, "health check: snapshot %s was saved, but the report could not be written: %v\n", collection.SnapshotID, err)
	}
	return nil
}
