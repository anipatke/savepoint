package doctor

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/resume"
)

// RunV2Checks is the live doctor entry point. It accepts a .savepoint root
// only after the command cutover preflight has established that the project
// is V2, then keeps the runtime on strict V2 readers: no schema-dispatched
// Project load, V1 router reader, directory-shaped task discovery, or audit
// register parser is reachable from this path.
func RunV2Checks(root string) *DiagnosticReport {
	report := &DiagnosticReport{}
	report.ConfigCheck = CheckConfig(root)
	router, routerErr := checkRouterV2(root)
	report.RouterCheck = routerErr
	report.Gates.Results = runQualityGatesV2(root)

	version, err := data.ReadSchemaVersion(filepath.Join(root, "config.yml"))
	if err != nil {
		name := "schema-version-malformed"
		if errors.Is(err, data.ErrUnsupportedSchemaVersion) {
			name = "schema-version-unsupported"
		}
		report.Project = []Problem{{
			File:    filepath.Join(root, "config.yml"),
			Message: fmt.Sprintf("[%s] %v", name, err),
			Repair:  V2ProblemRepair(name),
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
	if router != nil && router.HasRetiredNextAction {
		report.Project = append(report.Project, Problem{
			File:     filepath.Join(root, "router.md"),
			Message:  "[router-next-action-retired] router.md has a retired next_action field",
			Repair:   "Delete the next_action line from router.md.",
			Category: HealthPendingReview,
		})
	}

	index, err := data.LoadV2Index(root)
	if err != nil {
		name := v2DiagnosticName(err)
		report.Project = append(report.Project, Problem{
			File:    root,
			Message: fmt.Sprintf("[%s] %v", name, err),
			Repair:  V2ProblemRepair(name),
		})
		return report
	}

	report.Project = append(report.Project, v2ConsistencyProblems(index)...)
	for _, duplicate := range index.DuplicateObjectiveRanks {
		file := root
		if goal := index.Releases[duplicate.GoalID]; goal != nil {
			file = filepath.Join(root, goal.Source.Path)
		}
		report.Project = append(report.Project, Problem{
			File: file,
			Message: fmt.Sprintf("[v2-objective-rank-duplicate] Goal %s has duplicate rank %d at %s priority for Objectives %s",
				duplicate.GoalID, duplicate.Rank, duplicate.Priority, strings.Join(duplicate.ObjectiveIDs, ", ")),
			Repair:   fmt.Sprintf("Assign these Objectives distinct positive rank values within Goal %s at %s priority.", duplicate.GoalID, duplicate.Priority),
			Category: HealthPendingReview,
		})
	}
	if _, diagnostic := data.ResolveSelection(index, router); diagnostic != nil {
		switch diagnostic.Kind {
		case data.SelectionReleaseMissing, data.SelectionReleaseNotFound, data.SelectionReleaseArchived:
			report.Project = append(report.Project, Problem{
				File:    filepath.Join(root, "router.md"),
				Message: "[router-goal-missing] " + resume.SelectionPhrase(diagnostic),
				Repair:  routerGoalRepair(index),
			})
		case data.SelectionDone:
			report.Project = append(report.Project, Problem{
				File:     filepath.Join(root, "router.md"),
				Message:  resume.SelectionPhrase(diagnostic),
				Repair:   "Use p on an unfinished Task in the board, or edit the router's objective/task selection in router.md",
				Category: HealthPendingReview,
			})
		}
	}
	for _, objectiveID := range index.ObjectivesWithoutGoal {
		objective := index.Objectives[objectiveID]
		report.Project = append(report.Project, Problem{
			File:    filepath.Join(root, objective.Source.Path),
			Message: fmt.Sprintf("[objective-goal-missing] Objective %s has no Goal reference", objectiveID),
			Repair:  objectiveGoalRepair(index),
		})
	}
	releaseDiagnostics := releaseDiagnosticsForIndex(root, index)
	report.Releases = releaseDiagnostics.Problems
	report.ReleaseNotes = releaseDiagnostics.Notes
	report.Issues = issuePostureForIndex(index)
	return report
}

func routerGoalRepair(index *data.V2Index) string {
	if !index.HasLiveGoal() {
		return "Create a Goal first, then choose it with g on the board."
	}
	return "Choose a Goal with g on the board."
}

func objectiveGoalRepair(index *data.V2Index) string {
	if !index.HasLiveGoal() {
		return "Create a Goal first, then add its `release:` field using an R-### or G-### Goal ID."
	}
	return "Add a `release:` field to this Objective file using an existing R-### or G-### Goal ID."
}

func checkRouterV2(root string) (*data.RouterStateV2, error) {
	path := filepath.Join(root, "router.md")
	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("router.md not found: %w", data.ErrConfigNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("router.md unreadable: %w", err)
	}
	router, err := data.NewRouterReader().ReadStateV2(string(raw))
	if err != nil {
		return nil, fmt.Errorf("router.md invalid V2 state block: %w", err)
	}
	return router, nil
}

// runQualityGatesV2 reads only configuration and executes the configured
// quality gates.
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
