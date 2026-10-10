package doctor

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/opencode/savepoint/internal/run"
)

const runnerLost = "runner_lost"

// lostRunProblems names each run whose runner died without a final state. It
// only reads; a project without a runs directory reports nothing. It is advice
// for the owner, so it never counts as a structural problem.
func lostRunProblems(root string) []Problem {
	var problems []Problem
	for _, name := range run.LostRuns(filepath.Dir(root), time.Now()) {
		problems = append(problems, Problem{
			File:     filepath.Join(root, "runs", name, "run.yml"),
			Message:  fmt.Sprintf("[%s] run %s stopped updating and its process is gone", runnerLost, name),
			Repair:   "Read events.jsonl in that run directory to see how far it got, then start the Task again if needed.",
			Category: HealthPendingReview,
		})
	}
	return problems
}
