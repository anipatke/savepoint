package run

import (
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v3"
)

// LostRuns names the run directories under projectRoot whose runner died
// without a final state (see RunnerLost). It only reads: a project with no
// runs directory, or a run.yml that cannot be read, yields nothing.
func LostRuns(projectRoot string, now time.Time) []string {
	runs := filepath.Join(projectRoot, ".savepoint", "runs")
	entries, err := os.ReadDir(runs)
	if err != nil {
		return nil
	}
	var lost []string
	for _, entry := range entries {
		content, err := os.ReadFile(filepath.Join(runs, entry.Name(), recordFile))
		if err != nil {
			continue
		}
		var rec Record
		if yaml.Unmarshal(content, &rec) == nil && RunnerLost(rec, now) {
			lost = append(lost, entry.Name())
		}
	}
	sort.Strings(lost)
	return lost
}
