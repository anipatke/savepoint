package doctor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/opencode/savepoint/internal/data"
)

// RunV2Checks is the live doctor entry point. It accepts a .savepoint root
// only after the command cutover preflight has established that the project
// is V2, then keeps the runtime on strict V2 readers: no schema-dispatched
// Project load, V1 router reader, directory-shaped task discovery, or audit
// register parser is reachable from this path.
func RunV2Checks(root string) *DiagnosticReport {
	report := &DiagnosticReport{}
	report.ConfigCheck = CheckConfig(root)
	report.RouterCheck = checkRouterV2(root)
	report.Migration = CheckMigration(root)
	report.Gates.Results = runQualityGatesV2(root)

	version, err := data.ReadSchemaVersion(filepath.Join(root, "config.yml"))
	if err != nil {
		report.Project = []Problem{{
			File:    filepath.Join(root, "config.yml"),
			Message: fmt.Sprintf("[schema-version-malformed] %v", err),
			Repair:  V2ProblemRepair("schema-version-malformed"),
		}}
		return report
	}
	if version != data.SchemaVersionV2 {
		report.Project = []Problem{{
			File:    filepath.Join(root, "config.yml"),
			Message: "[schema-version-unsupported] live doctor accepts schema_version: 2 only; run `savepoint migrate --dry-run` for legacy input",
			Repair:  V2ProblemRepair("schema-version-unsupported"),
		}}
		return report
	}

	index, err := data.LoadV2Index(root)
	if err != nil {
		name := v2DiagnosticName(err)
		report.Project = []Problem{{
			File:    root,
			Message: fmt.Sprintf("[%s] %v", name, err),
			Repair:  V2ProblemRepair(name),
		}}
		return report
	}

	report.Project = v2ConsistencyProblems(index)
	releaseDiagnostics := releaseDiagnosticsForIndex(root, index)
	report.Releases = releaseDiagnostics.Problems
	report.ReleaseNotes = releaseDiagnostics.Notes
	report.Issues = issuePostureForIndex(index)
	return report
}

func checkRouterV2(root string) error {
	path := filepath.Join(root, "router.md")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("router.md not found: %w", data.ErrConfigNotFound)
	}
	if err != nil {
		return fmt.Errorf("router.md unreadable: %w", err)
	}
	if _, err := data.NewRouterReader().ReadStateV2(string(raw)); err != nil {
		return fmt.Errorf("router.md invalid V2 state block: %w", err)
	}
	return nil
}

// runQualityGatesV2 is the V2-only counterpart to the historical injected
// gate runner. It reads only configuration and executes the same gate
// primitive; it does not construct the legacy discoverer/parser dependencies
// held by DoctorDependencies for frozen tests.
func runQualityGatesV2(root string) []GateResult {
	configPath := filepath.Join(root, "config.yml")
	cfg, err := data.NewConfigReader().Read(configPath)
	if err != nil {
		return []GateResult{{
			Name:   "config",
			Passed: false,
			Output: fmt.Sprintf("cannot read config: %v", err),
		}}
	}

	return runConfiguredQualityGates(root, cfg)
}
