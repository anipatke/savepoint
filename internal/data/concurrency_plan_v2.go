package data

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// This file holds the optional, advisory parallel-planning metadata of V2
// records: Objective lanes, Task lane references and read/write manifests,
// and Task-pair independence explanations. Nothing here is required for a
// record to load or for a lifecycle gate to decide. Malformed or unresolvable
// advisory metadata becomes a PlanDiagnosticV2 that names the offending
// record and makes only the affected plan unusable; required record and schema
// validation elsewhere is unchanged. Values are immutable to callers, and the
// helpers are pure: no filesystem probing, scheduling, or lifecycle input.

var planLaneKeyPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// PlanDiagnosticV2 describes one unusable piece of advisory planning metadata.
type PlanDiagnosticV2 struct {
	Record  string // O-### or T-### that carries the metadata
	Path    string // source file of that record
	Field   string // frontmatter field, for example planned_writes
	Message string
}

func (d PlanDiagnosticV2) String() string {
	return fmt.Sprintf("%s: %s %s: %s", d.Path, d.Record, d.Field, d.Message)
}

// PlanScopeV2 is a Task's planned read or write list. An omitted list is
// unknown (Declared false); an explicit empty list is a reviewed empty scope.
type PlanScopeV2 struct {
	declared bool
	paths    []string
}

// Declared reports whether the list was authored, including as empty.
func (s PlanScopeV2) Declared() bool { return s.declared }

// Paths returns a copy of the authored paths in authored order.
func (s PlanScopeV2) Paths() []string { return slices.Clone(s.paths) }

// Contains reports whether path is listed exactly.
func (s PlanScopeV2) Contains(path string) bool { return slices.Contains(s.paths, path) }

// TaskPlanV2 is a Task's optional planning metadata. Lane is an Objective-local
// key and is the only membership record; it is resolved against the owning
// Objective by the index, never here.
type TaskPlanV2 struct {
	Lane        string
	Reads       PlanScopeV2
	Writes      PlanScopeV2
	diagnostics []PlanDiagnosticV2
}

// Diagnostics returns a copy of the decode-time diagnostics. A plan with any
// diagnostic must not feed a recommendation.
func (p TaskPlanV2) Diagnostics() []PlanDiagnosticV2 { return slices.Clone(p.diagnostics) }

// ObjectiveLaneV2 is one named lane declared by an Objective.
type ObjectiveLaneV2 struct {
	Key   string
	Title string
}

// PlanOverlapV2 names one exact write/read overlap between two Tasks.
type PlanOverlapV2 struct {
	Writer string
	Reader string
	Path   string
}

// PlanIndependenceV2 explains why the listed exact overlaps between two
// Tasks do not couple them. Reviewed binds the explanation to the manifests
// that were reviewed (see PlanReviewDigest). It never explains shared writes
// or dependencies.
type PlanIndependenceV2 struct {
	Tasks    [2]string
	Overlaps []PlanOverlapV2
	Reason   string
	Reviewed string
}

// ObjectivePlanV2 is an Objective's optional planning metadata.
type ObjectivePlanV2 struct {
	Lanes       []ObjectiveLaneV2
	Independent []PlanIndependenceV2
	diagnostics []PlanDiagnosticV2
}

// Diagnostics returns a copy of the decode-time diagnostics.
func (p ObjectivePlanV2) Diagnostics() []PlanDiagnosticV2 { return slices.Clone(p.diagnostics) }

// Lane returns the declared lane with key.
func (p ObjectivePlanV2) Lane(key string) (ObjectiveLaneV2, bool) {
	for _, lane := range p.Lanes {
		if lane.Key == key {
			return lane, true
		}
	}
	return ObjectiveLaneV2{}, false
}

// ValidatePlanPath reports why path is not an exact portable project-relative
// file path. Globs, directories, absolute, drive and UNC paths, traversal and
// Windows-ambiguous names are rejected. A path that does not exist yet is
// valid; nothing touches the filesystem.
func ValidatePlanPath(path string) error {
	switch {
	case path == "":
		return fmt.Errorf("path is empty")
	case path != strings.TrimSpace(path):
		return fmt.Errorf("path %q has surrounding whitespace", path)
	case strings.ContainsAny(path, "*?[]{}"):
		return fmt.Errorf("path %q is a glob; list exact files", path)
	case strings.Contains(path, `\`):
		return fmt.Errorf("path %q uses a backslash; use forward slashes", path)
	case strings.Contains(path, ":"):
		return fmt.Errorf("path %q has a colon (drive or stream); use a project-relative path", path)
	case strings.HasPrefix(path, "/"):
		return fmt.Errorf("path %q is absolute or UNC; use a project-relative path", path)
	case strings.HasSuffix(path, "/"):
		return fmt.Errorf("path %q is a directory; list exact files", path)
	}
	for _, r := range path {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("path %q has a control character", path)
		}
	}
	for _, segment := range strings.Split(path, "/") {
		switch {
		case segment == "":
			return fmt.Errorf("path %q has an empty segment", path)
		case segment == "." || segment == "..":
			return fmt.Errorf("path %q has a %q segment; use a normalized path", path, segment)
		case strings.HasSuffix(segment, ".") || strings.HasSuffix(segment, " "):
			return fmt.Errorf("path %q has a segment ending in a dot or space, which Windows aliases", path)
		}
	}
	return nil
}

// PlanPathRelation is how two valid plan paths compare.
type PlanPathRelation int

const (
	PlanPathsDistinct PlanPathRelation = iota
	PlanPathsSame
	// PlanPathsMayAlias means the paths differ only by case, so a
	// case-insensitive checkout may treat them as one file. Lexical
	// comparison never certifies them as distinct.
	PlanPathsMayAlias
)

// ComparePlanPaths relates two valid plan paths conservatively.
func ComparePlanPaths(a, b string) PlanPathRelation {
	switch {
	case a == b:
		return PlanPathsSame
	case strings.ToLower(a) == strings.ToLower(b):
		return PlanPathsMayAlias
	default:
		return PlanPathsDistinct
	}
}

// PlanReviewDigest binds an independence explanation to the reviewed
// manifests of its two Tasks, so editing either manifest invalidates it.
func PlanReviewDigest(idA string, planA TaskPlanV2, idB string, planB TaskPlanV2) string {
	type side struct {
		id   string
		plan TaskPlanV2
	}
	sides := []side{{idA, planA}, {idB, planB}}
	slices.SortFunc(sides, func(a, b side) int { return strings.Compare(a.id, b.id) })

	var canonical strings.Builder
	for _, s := range sides {
		fmt.Fprintf(&canonical, "task\x00%s\n", s.id)
		for _, scope := range []struct {
			name  string
			value PlanScopeV2
		}{{"reads", s.plan.Reads}, {"writes", s.plan.Writes}} {
			paths := scope.value.Paths()
			slices.Sort(paths)
			fmt.Fprintf(&canonical, "%s\x00%t\x00%s\n", scope.name, scope.value.declared, strings.Join(paths, "\x00"))
		}
	}
	sum := sha256.Sum256([]byte(canonical.String()))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func decodeTaskPlan(path, taskID string, lane, reads, writes yaml.Node) TaskPlanV2 {
	var plan TaskPlanV2
	diagnose := func(field, format string, args ...any) {
		plan.diagnostics = append(plan.diagnostics, PlanDiagnosticV2{Record: taskID, Path: path, Field: field, Message: fmt.Sprintf(format, args...)})
	}

	if lane.Kind != 0 {
		if lane.Kind != yaml.ScalarNode || lane.Tag != "!!str" || !planLaneKeyPattern.MatchString(lane.Value) {
			diagnose("lane", "must be a lane key of lowercase letters, digits and hyphens")
		} else {
			plan.Lane = lane.Value
		}
	}
	for _, field := range []struct {
		name string
		node yaml.Node
		dst  *PlanScopeV2
	}{{"planned_reads", reads, &plan.Reads}, {"planned_writes", writes, &plan.Writes}} {
		scope, problems := decodePlanScope(field.node)
		for _, problem := range problems {
			diagnose(field.name, "%s", problem)
		}
		if len(problems) == 0 {
			*field.dst = scope
		}
	}
	for _, problem := range planScopeAliasProblems(plan.Reads, plan.Writes) {
		diagnose("planned_writes", "%s", problem)
	}
	return plan
}

// decodePlanScope decodes an optional list of exact paths. Exact duplicates
// collapse; case-only variants are reported because they may be one file.
func decodePlanScope(node yaml.Node) (PlanScopeV2, []string) {
	if node.Kind == 0 {
		return PlanScopeV2{}, nil
	}
	if node.Kind != yaml.SequenceNode {
		return PlanScopeV2{}, []string{"must be a list of exact project-relative file paths"}
	}
	scope := PlanScopeV2{declared: true, paths: []string{}}
	var problems []string
	for _, item := range node.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
			problems = append(problems, "entries must be path strings")
			continue
		}
		if err := ValidatePlanPath(item.Value); err != nil {
			problems = append(problems, err.Error())
			continue
		}
		if !slices.Contains(scope.paths, item.Value) {
			scope.paths = append(scope.paths, item.Value)
		}
	}
	problems = append(problems, planScopeAliasProblems(scope)...)
	return scope, problems
}

// planScopeAliasProblems reports distinct paths that differ only by case
// across the given scopes.
func planScopeAliasProblems(scopes ...PlanScopeV2) []string {
	seen := map[string]string{}
	var problems []string
	for _, scope := range scopes {
		for _, path := range scope.paths {
			key := strings.ToLower(path)
			if first, ok := seen[key]; ok && first != path {
				problems = append(problems, fmt.Sprintf("paths %q and %q differ only by case and may be one file", first, path))
				continue
			}
			seen[key] = path
		}
	}
	return problems
}

func decodeObjectivePlan(path, objectiveID string, lanes, independence yaml.Node) ObjectivePlanV2 {
	var plan ObjectivePlanV2
	diagnose := func(field, format string, args ...any) {
		plan.diagnostics = append(plan.diagnostics, PlanDiagnosticV2{Record: objectiveID, Path: path, Field: field, Message: fmt.Sprintf(format, args...)})
	}

	if lanes.Kind != 0 {
		if lanes.Kind != yaml.SequenceNode {
			diagnose("lanes", "must be a list of {key, title} entries")
		}
		seen := map[string]bool{}
		ambiguous := map[string]bool{}
		for i, item := range lanes.Content {
			var raw struct {
				Key   string `yaml:"key"`
				Title string `yaml:"title"`
			}
			if lanes.Kind != yaml.SequenceNode || item.Kind != yaml.MappingNode || item.Decode(&raw) != nil {
				diagnose("lanes", "entry %d must be a {key, title} mapping", i+1)
				continue
			}
			switch {
			case !planLaneKeyPattern.MatchString(raw.Key):
				diagnose("lanes", "entry %d key %q must use lowercase letters, digits and hyphens", i+1, raw.Key)
			case strings.TrimSpace(raw.Title) == "":
				diagnose("lanes", "lane %q needs a readable title", raw.Key)
			case seen[raw.Key]:
				diagnose("lanes", "lane key %q is declared more than once; none of its declarations is used", raw.Key)
				ambiguous[raw.Key] = true
			default:
				seen[raw.Key] = true
				plan.Lanes = append(plan.Lanes, ObjectiveLaneV2{Key: raw.Key, Title: raw.Title})
			}
		}
		// A key declared twice is ambiguous: neither title can be trusted, so the
		// lane is dropped and its Tasks become unusable for recommendations, the
		// same as a lane that was never declared.
		plan.Lanes = slices.DeleteFunc(plan.Lanes, func(l ObjectiveLaneV2) bool { return ambiguous[l.Key] })
	}

	if independence.Kind != 0 {
		if independence.Kind != yaml.SequenceNode {
			diagnose("independence", "must be a list of explanations")
		}
		for i, item := range independence.Content {
			explanation, problem := decodeIndependence(item)
			if independence.Kind != yaml.SequenceNode {
				break
			}
			if problem != "" {
				diagnose("independence", "entry %d %s", i+1, problem)
				continue
			}
			plan.Independent = append(plan.Independent, explanation)
		}
	}
	return plan
}

func decodeIndependence(node *yaml.Node) (PlanIndependenceV2, string) {
	var raw struct {
		Tasks    []string `yaml:"tasks"`
		Overlaps []struct {
			Writer string `yaml:"writer"`
			Reader string `yaml:"reader"`
			Path   string `yaml:"path"`
		} `yaml:"overlaps"`
		Reason   string `yaml:"reason"`
		Reviewed string `yaml:"reviewed"`
	}
	if node.Kind != yaml.MappingNode || node.Decode(&raw) != nil {
		return PlanIndependenceV2{}, "must be a mapping of tasks, overlaps, reason and reviewed"
	}
	if len(raw.Tasks) != 2 || !matchesV2Identity(raw.Tasks[0], 'T') || !matchesV2Identity(raw.Tasks[1], 'T') {
		return PlanIndependenceV2{}, "tasks must name exactly two T-### Tasks"
	}
	if raw.Tasks[0] == raw.Tasks[1] {
		return PlanIndependenceV2{}, "tasks must name two different Tasks"
	}
	if strings.TrimSpace(raw.Reason) == "" {
		return PlanIndependenceV2{}, "needs a non-empty reason"
	}
	if len(raw.Overlaps) == 0 {
		return PlanIndependenceV2{}, "must list the exact write/read overlaps it explains"
	}
	explanation := PlanIndependenceV2{Reason: raw.Reason, Reviewed: raw.Reviewed}
	explanation.Tasks = [2]string{raw.Tasks[0], raw.Tasks[1]}
	if explanation.Tasks[0] > explanation.Tasks[1] {
		explanation.Tasks[0], explanation.Tasks[1] = explanation.Tasks[1], explanation.Tasks[0]
	}
	for _, overlap := range raw.Overlaps {
		explanation.Overlaps = append(explanation.Overlaps, PlanOverlapV2{Writer: overlap.Writer, Reader: overlap.Reader, Path: overlap.Path})
	}
	return explanation, ""
}

// validateV2Planning resolves advisory metadata against the indexed records
// and records what it cannot use. It never returns an error: a plan that does
// not resolve is only unusable for recommendations.
func validateV2Planning(index *V2Index) {
	index.PlanDiagnostics = nil
	index.UnusablePlans = map[string]bool{}
	index.PlanIndependence = map[string][]PlanIndependenceV2{}

	for _, id := range slices.Sorted(maps.Keys(index.Objectives)) {
		index.PlanDiagnostics = append(index.PlanDiagnostics, index.Objectives[id].Plan.diagnostics...)
	}
	for _, id := range slices.Sorted(maps.Keys(index.Tasks)) {
		task := index.Tasks[id]
		plan := task.Plan
		index.PlanDiagnostics = append(index.PlanDiagnostics, plan.diagnostics...)
		unusable := len(plan.diagnostics) > 0
		if plan.Lane != "" {
			if _, ok := index.Objectives[task.Objective].Plan.Lane(plan.Lane); !ok {
				index.PlanDiagnostics = append(index.PlanDiagnostics, PlanDiagnosticV2{
					Record: id, Path: task.Source.Path, Field: "lane",
					Message: fmt.Sprintf("lane %q is not declared by owning objective %s", plan.Lane, task.Objective),
				})
				unusable = true
			}
		}
		if unusable {
			index.UnusablePlans[id] = true
		}
	}

	for _, objectiveID := range slices.Sorted(maps.Keys(index.Objectives)) {
		objective := index.Objectives[objectiveID]
		pairs := map[[2]string]int{}
		for _, explanation := range objective.Plan.Independent {
			pairs[explanation.Tasks]++
		}
		for _, explanation := range objective.Plan.Independent {
			problem := independenceProblem(index, objectiveID, explanation, pairs[explanation.Tasks])
			if problem == "" {
				index.PlanIndependence[objectiveID] = append(index.PlanIndependence[objectiveID], explanation)
				continue
			}
			index.PlanDiagnostics = append(index.PlanDiagnostics, PlanDiagnosticV2{
				Record: objectiveID, Path: objective.Source.Path, Field: "independence",
				Message: fmt.Sprintf("explanation for %s and %s %s", explanation.Tasks[0], explanation.Tasks[1], problem),
			})
		}
	}
}

// independenceProblem returns why an explanation cannot be used, or "".
func independenceProblem(index *V2Index, objectiveID string, e PlanIndependenceV2, declarations int) string {
	if declarations > 1 {
		return "is declared more than once; none of the duplicates is used"
	}
	var plans [2]TaskPlanV2
	for i, id := range e.Tasks {
		task, ok := index.Tasks[id]
		switch {
		case !ok:
			return fmt.Sprintf("names missing Task %s", id)
		case task.Objective != objectiveID:
			return fmt.Sprintf("names %s, which belongs to %s", id, task.Objective)
		case index.UnusablePlans[id]:
			return fmt.Sprintf("depends on %s, whose planning metadata is unusable", id)
		}
		plans[i] = task.Plan
	}
	for _, overlap := range e.Overlaps {
		var writer, reader TaskPlanV2
		switch {
		case overlap.Writer == e.Tasks[0] && overlap.Reader == e.Tasks[1]:
			writer, reader = plans[0], plans[1]
		case overlap.Writer == e.Tasks[1] && overlap.Reader == e.Tasks[0]:
			writer, reader = plans[1], plans[0]
		default:
			return fmt.Sprintf("overlap on %q must name one Task as writer and the other as reader", overlap.Path)
		}
		if !writer.Writes.Contains(overlap.Path) || !reader.Reads.Contains(overlap.Path) {
			return fmt.Sprintf("overlap on %q is not a write by %s and a read by %s", overlap.Path, overlap.Writer, overlap.Reader)
		}
		if reader.Writes.Contains(overlap.Path) {
			return fmt.Sprintf("overlap on %q is also a shared write, which cannot be explained", overlap.Path)
		}
	}
	if e.Reviewed != PlanReviewDigest(e.Tasks[0], plans[0], e.Tasks[1], plans[1]) {
		return "is stale or unreviewed; its Task manifests changed since the review"
	}
	return ""
}
