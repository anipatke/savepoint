// Package healthcheck runs one official Code Health collection for an
// Objective: it validates the Objective against the strictly loaded V2 index,
// collects through codehealth, and renders the plain blocking verdict.
package healthcheck

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
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
	root, err := resolveHealthTarget(req)
	if err != nil {
		return err
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
		Root:      root,
		Origin:    codehealth.OriginOfficial,
		Config:    cfg,
		Readers:   codehealth.DefaultReaders(),
		Runner:    req.Runner,
		Git:       req.Git,
		Objective: req.Objective,
	})
	if errors.Is(err, codehealth.ErrCollectionCancelled) {
		return fmt.Errorf("health check: cancelled; no snapshot was saved: %w", err)
	}
	if err != nil {
		return fmt.Errorf("health check: no snapshot was saved: %w", err)
	}

	verdict, err := evaluateSaved(store, cfg, collection.SnapshotID)
	if err != nil {
		return err
	}

	if collection.ReportErr != nil && req.Stderr != nil {
		fmt.Fprintf(req.Stderr, "warning: snapshot %s was saved, but the report could not be written: %v\n", collection.SnapshotID, collection.ReportErr)
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

// resolveHealthTarget returns the project root once it is a Savepoint project at
// a supported schema that holds the requested Objective.
func resolveHealthTarget(req Request) (string, error) {
	root, err := data.ResolveTarget(req.Dir)
	switch {
	case errors.Is(err, data.ErrTargetMissing):
		return "", fmt.Errorf("health check: target directory does not exist: %s", req.Dir)
	case errors.Is(err, data.ErrTargetNotSavepoint):
		return "", fmt.Errorf("health check: target directory is not a Savepoint project: %s has no .savepoint directory", req.Dir)
	case err != nil:
		return "", fmt.Errorf("health check: %w", err)
	}
	if err := data.CheckRuntimeSchema(root); err != nil {
		return "", fmt.Errorf("health check: %w", err)
	}
	index, err := data.LoadV2Index(filepath.Join(root, ".savepoint"))
	if err != nil {
		return "", fmt.Errorf("health check: %w", err)
	}
	if _, ok := index.Objectives[req.Objective]; !ok {
		return "", fmt.Errorf("health check: objective %s does not exist in this project", req.Objective)
	}
	return root, nil
}

// evaluateSaved reads the just-saved snapshot back from the store and
// evaluates it against the configuration.
func evaluateSaved(store codehealth.Store, cfg codehealth.Config, id string) (codehealth.Verdict, error) {
	snapshots, err := store.LoadSnapshots()
	if err != nil {
		return codehealth.Verdict{}, fmt.Errorf("health check: snapshot %s was saved but cannot be read back: %w", id, err)
	}
	var verdict codehealth.Verdict
	found := false
	for _, snapshot := range snapshots {
		if snapshot.ID != id {
			continue
		}
		if verdict, err = codehealth.Evaluate(snapshot, cfg); err != nil {
			return codehealth.Verdict{}, fmt.Errorf("health check: snapshot %s was saved but cannot be evaluated: %w", id, err)
		}
		found = true
	}
	if !found {
		return codehealth.Verdict{}, fmt.Errorf("health check: snapshot %s was saved but is missing from the store", id)
	}
	return verdict, nil
}

// ReportRequest is one rewrite of the Code Health report.
type ReportRequest struct {
	Dir string
}

// RunReport rewrites .savepoint/health/report.md from the newest saved snapshot
// and prints its path. It runs no tool and saves no snapshot. Without Code
// Health setup or a snapshot it writes nothing and returns that as an error.
// The re-run line names the router's Objective, or O-### when none is selected.
func RunReport(req ReportRequest, stdout io.Writer) error {
	root, err := data.ResolveTarget(req.Dir)
	switch {
	case errors.Is(err, data.ErrTargetMissing):
		return fmt.Errorf("health report: target directory does not exist: %s", req.Dir)
	case errors.Is(err, data.ErrTargetNotSavepoint):
		return fmt.Errorf("health report: target directory is not a Savepoint project: %s has no .savepoint directory", req.Dir)
	case err != nil:
		return fmt.Errorf("health report: %w", err)
	}
	if err := data.CheckRuntimeSchema(root); err != nil {
		return fmt.Errorf("health report: %w", err)
	}
	path, _, err := codehealth.RefreshReport(root, routerObjective(root))
	if err != nil {
		return fmt.Errorf("health report: %w", err)
	}
	_, err = fmt.Fprintf(stdout, "Report written to %s\n", path)
	return err
}

// routerObjective is the router's selected Objective, or empty when the router
// cannot be read or selects none.
func routerObjective(root string) string {
	content, err := os.ReadFile(filepath.Join(root, ".savepoint", "router.md"))
	if err != nil {
		return ""
	}
	state, err := data.NewRouterReader().ReadStateV2(string(content))
	if err != nil {
		return ""
	}
	return state.Objective
}
