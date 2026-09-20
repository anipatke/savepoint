// Package migrate owns the operational side of the V2 cutover boundary.
//
// PreflightCutover is deliberately read-only. It composes the migration
// planner, the durable-operation recovery checks, V2 loading, and the
// canonical data.ResolveReleaseCutover result, but it never applies a plan,
// probes writeability, or invents a second Release policy.
package migrate

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/opencode/savepoint/internal/data"
)

// CutoverBlockKind is the stable category of one operational refusal. The
// categories intentionally distinguish a project that needs migration from a
// project whose migration plan, recovery record, schema, or V2 records are
// themselves unsafe to proceed with.
type CutoverBlockKind string

const (
	CutoverBlockTarget             CutoverBlockKind = "target"
	CutoverBlockPendingOperation   CutoverBlockKind = "pending_operation"
	CutoverBlockRecoveryConflict   CutoverBlockKind = "recovery_conflict"
	CutoverBlockUnrecoverablePlan  CutoverBlockKind = "unrecoverable_plan"
	CutoverBlockUnsupportedSchema  CutoverBlockKind = "unsupported_schema"
	CutoverBlockMalformedSchema    CutoverBlockKind = "malformed_schema"
	CutoverBlockSchemaRead         CutoverBlockKind = "schema_read"
	CutoverBlockMigrationRequired  CutoverBlockKind = "migration_required"
	CutoverBlockMigrationPlan      CutoverBlockKind = "migration_plan"
	CutoverBlockMigrationConflict  CutoverBlockKind = "migration_conflict"
	CutoverBlockMigrationAmbiguity CutoverBlockKind = "migration_ambiguity"
	CutoverBlockInvalidV2          CutoverBlockKind = "invalid_v2"
	CutoverBlockRelease            CutoverBlockKind = "release"
)

// CutoverBlocker names one condition that prevents a cutover. The optional
// fields retain the typed evidence behind the condition so callers can render
// actionable diagnostics without parsing Detail. Release blockers carry the
// exact GateBlocker returned by data.ResolveReleaseCutover.
type CutoverBlocker struct {
	Kind        CutoverBlockKind
	Detail      string
	Path        string
	ReleaseID   string
	OperationID string
	AmbiguityID string
	Gate        data.GateBlocker
}

// CutoverPreflightOptions supplies the deterministic inputs needed when a V1
// project is planned. Zero options are valid for production callers; the
// defaults are the same clock and operation-ID source used by migrate's
// command path.
type CutoverPreflightOptions struct {
	Decisions      Decisions
	Now            Clock
	NewOperationID OperationIDSource
}

// CutoverOptions is retained as a short alias for callers that already name
// this operation simply "cutover".
type CutoverOptions = CutoverPreflightOptions

// CutoverPreflightResult is the complete, fail-closed answer for one project.
// Allowed is true exactly when Blockers is empty. Plan is populated for a V1
// project so the caller can show the reviewed migration result; Index is
// populated for a valid V2 project so downstream consumers can reuse the same
// loaded interpretation. Neither pointer is a promise that any write took
// place.
type CutoverPreflightResult struct {
	Allowed     bool
	ProjectRoot string
	Schema      data.SchemaVersion
	Plan        *ConversionPlan
	Index       *data.V2Index
	Pending     *PendingOperationReport
	Blockers    []CutoverBlocker
}

// CutoverResult is a short alias for consumers that do not need the longer
// preflight name in their local code.
type CutoverResult = CutoverPreflightResult

// PreflightCutover performs the complete read-only cutover check for a
// project directory. It never returns an operational refusal as a Go error:
// every refusal is a typed Blocker so a command can report all conditions in
// one pass. The result's Allowed field is false for any target, schema,
// migration, recovery, V2, or Release blocker.
func PreflightCutover(projectRoot string, options CutoverPreflightOptions) CutoverPreflightResult {
	options = options.withDefaults()
	result := CutoverPreflightResult{Schema: data.SchemaVersionV1}

	root, err := ResolveTarget(projectRoot)
	if err != nil {
		result.add(CutoverBlocker{
			Kind:   CutoverBlockTarget,
			Detail: fmt.Sprintf("cutover target is not usable: %v; provide a project directory containing .savepoint", err),
		})
		return result
	}
	result.ProjectRoot = root

	if !result.checkPendingOperation(root) {
		return result.finish()
	}

	configPath := filepath.Join(root, ".savepoint", "config.yml")
	version, err := data.ReadSchemaVersion(configPath)
	result.Schema = version
	if err != nil {
		result.add(CutoverBlocker{
			Kind:   schemaBlockKind(err),
			Path:   filepath.ToSlash(filepath.Join(".savepoint", "config.yml")),
			Detail: schemaBlockDetail(configPath, err),
		})
		return result.finish()
	}

	if version != data.SchemaVersionV2 {
		return result.finishV1(root, options)
	}

	result.Plan = &ConversionPlan{SchemaAlreadyV2: true}
	project, err := data.LoadProject(filepath.Join(root, ".savepoint"))
	if err != nil {
		result.add(CutoverBlocker{
			Kind:   CutoverBlockInvalidV2,
			Path:   ".savepoint",
			Detail: fmt.Sprintf("V2 project cannot be loaded safely: %v; repair the named record before cutover", err),
		})
		return result.finish()
	}
	if project.SchemaVersion != data.SchemaVersionV2 || project.V2 == nil {
		result.add(CutoverBlocker{
			Kind:   CutoverBlockInvalidV2,
			Path:   ".savepoint",
			Detail: "V2 cutover preflight observed a project that did not resolve to a V2 index; reread the project and repair the schema before cutover",
		})
		return result.finish()
	}
	result.Index = project.V2

	// This is the only Release policy consulted here. The migration package
	// only carries the typed result across the operational boundary.
	releaseDecision := data.ResolveReleaseCutover(project.V2)
	for _, blocker := range releaseDecision.Blockers {
		detail := blocker.Gate.Detail
		if detail == "" {
			detail = fmt.Sprintf("canonical Release completion returned %s", blocker.Gate.Kind)
		}
		if blocker.ReleaseID != "" {
			detail = fmt.Sprintf("release %s: %s", blocker.ReleaseID, detail)
		}
		result.add(CutoverBlocker{
			Kind:      CutoverBlockRelease,
			Detail:    detail,
			ReleaseID: blocker.ReleaseID,
			Gate:      blocker.Gate,
		})
	}

	return result.finish()
}

// EvaluateCutover is an explicit synonym for callers that prefer a verb at
// the call site. It is intentionally the same read-only operation.
func EvaluateCutover(projectRoot string, options CutoverPreflightOptions) CutoverPreflightResult {
	return PreflightCutover(projectRoot, options)
}

func (o CutoverPreflightOptions) withDefaults() CutoverPreflightOptions {
	if o.Now == nil {
		o.Now = time.Now
	}
	if o.NewOperationID == nil {
		o.NewOperationID = NewOperationID
	}
	return o
}

// checkPendingOperation handles recovery before schema dispatch. A pending
// operation is authoritative: even when every recorded path is verified, the
// final schema activation still has to be completed by an explicit recovery.
// Revalidation and recoverability checks remain read-only and are included as
// separate blockers when they reveal a second problem.
func (r *CutoverPreflightResult) checkPendingOperation(root string) bool {
	pending, err := PendingOperation(root)
	if err != nil {
		r.add(CutoverBlocker{
			Kind:   CutoverBlockPendingOperation,
			Detail: fmt.Sprintf("pending migration operation cannot be inspected: %v; resolve the operation before cutover", err),
		})
		return false
	}
	if pending == nil {
		return true
	}
	r.Pending = pending
	r.add(CutoverBlocker{
		Kind:        CutoverBlockPendingOperation,
		OperationID: pending.OperationID,
		Detail:      pending.RecoveryGuidance(),
	})

	op, err := LoadOperation(pending.Dir)
	if err != nil {
		r.add(CutoverBlocker{
			Kind:        CutoverBlockUnrecoverablePlan,
			OperationID: pending.OperationID,
			Detail:      fmt.Sprintf("pending operation %s cannot be loaded for recovery: %v; restore its operation record or discard the interrupted migration with owner approval", pending.OperationID, err),
		})
		return false
	}

	if err := revalidatePendingOperation(root, op); err != nil {
		r.add(CutoverBlocker{
			Kind:        CutoverBlockRecoveryConflict,
			OperationID: pending.OperationID,
			Detail:      fmt.Sprintf("pending operation %s is not safe to resume: %v; restore the recorded source/output state before recovery", pending.OperationID, err),
		})
	}
	if _, err := recoverApplyBatch(root, op, nil); err != nil {
		kind := CutoverBlockUnrecoverablePlan
		if !errors.Is(err, ErrRecoveryPlanMissing) {
			kind = CutoverBlockRecoveryConflict
		}
		r.add(CutoverBlocker{
			Kind:        kind,
			OperationID: pending.OperationID,
			Detail:      fmt.Sprintf("pending operation %s has no verified recovery path: %v; run migration recovery only after repairing the named state", pending.OperationID, err),
		})
	}
	return false
}

func (r *CutoverPreflightResult) finishV1(root string, options CutoverPreflightOptions) CutoverPreflightResult {
	plan, err := Plan(root, options.Decisions, options.Now, options.NewOperationID)
	if err != nil {
		r.add(CutoverBlocker{
			Kind:   CutoverBlockMigrationPlan,
			Detail: fmt.Sprintf("V1 project has no usable migration plan: %v; correct the named source or decision and rerun the migration preview", err),
		})
		return r.finish()
	}
	r.Plan = plan

	// A V1 project can never be cut over merely because its plan is clean: the
	// reviewed plan still has to be explicitly applied first.
	r.add(CutoverBlocker{
		Kind:   CutoverBlockMigrationRequired,
		Detail: "project is still V1; review `savepoint migrate --dry-run`, then apply the migration before cutover",
	})

	for _, conflict := range plan.Conflicts {
		r.add(CutoverBlocker{
			Kind:   CutoverBlockMigrationConflict,
			Path:   conflict.Path,
			Detail: fmt.Sprintf("migration conflict at %s: %s; resolve the existing destination or manifest before applying", conflict.Path, conflict.Detail),
		})
	}
	for _, ambiguity := range plan.Ambiguities {
		if !ambiguity.Blocking || ambiguity.Resolved {
			continue
		}
		r.add(CutoverBlocker{
			Kind:        CutoverBlockMigrationAmbiguity,
			Path:        ambiguity.Path,
			AmbiguityID: ambiguity.ID,
			Detail:      fmt.Sprintf("unresolved migration decision %s at %s: %s; record an explicit decision before applying", ambiguity.ID, ambiguity.Path, ambiguity.Detail),
		})
	}

	if !plan.Appliable && len(plan.UnresolvedBlockingIDs) == 0 && len(plan.Conflicts) == 0 {
		r.add(CutoverBlocker{
			Kind:   CutoverBlockMigrationPlan,
			Detail: "migration plan is marked non-appliable without a named blocking condition; regenerate the preview before cutover",
		})
	}
	return r.finish()
}

func (r *CutoverPreflightResult) add(blocker CutoverBlocker) {
	r.Blockers = append(r.Blockers, blocker)
}

func (r *CutoverPreflightResult) finish() CutoverPreflightResult {
	r.Allowed = len(r.Blockers) == 0
	return *r
}

func schemaBlockKind(err error) CutoverBlockKind {
	switch {
	case errors.Is(err, data.ErrUnsupportedSchemaVersion):
		return CutoverBlockUnsupportedSchema
	case errors.Is(err, data.ErrMalformedSchemaVersion):
		return CutoverBlockMalformedSchema
	default:
		return CutoverBlockSchemaRead
	}
}

func schemaBlockDetail(path string, err error) string {
	switch {
	case errors.Is(err, data.ErrUnsupportedSchemaVersion):
		return fmt.Sprintf("%v; set schema_version: 2 only after the reviewed V1 migration has completed", err)
	case errors.Is(err, data.ErrMalformedSchemaVersion):
		return fmt.Sprintf("%v; replace schema_version with an integer supported by this project", err)
	default:
		return fmt.Sprintf("cannot read schema at %s: %v; restore a readable config before cutover", path, err)
	}
}
