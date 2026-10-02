package codehealth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
)

// Reader for the duplication capability. The value is jscpd's own share of
// duplicated lines, cross-checked against its line counts; the largest repeated
// blocks are the affected items.

const (
	duplicationDefinition = "share of lines that are part of a repeated block"

	// jscpd rounds its percentage to two decimals.
	jscpdPercentTolerance = 0.5

	DetailDuplicatedLines = "duplicated_lines"
	DetailTotalLines      = "total_lines"
	DetailClones          = "clones"
	DetailSources         = "sources"
)

// JscpdReader reads jscpd JSON reports.
type JscpdReader struct{}

type jscpdLocation struct {
	Name  string `json:"name"`
	Start int    `json:"start"`
}

type jscpdTotal struct {
	Lines           *int     `json:"lines"`
	Sources         int      `json:"sources"`
	Clones          int      `json:"clones"`
	DuplicatedLines *int     `json:"duplicatedLines"`
	Percentage      *float64 `json:"percentage"`
}

// jscpdSource is one file's line counts, as jscpd writes them per format.
type jscpdSource struct {
	Lines           int `json:"lines"`
	DuplicatedLines int `json:"duplicatedLines"`
}

type jscpdReport struct {
	Version    string `json:"version"`
	Statistics *struct {
		Total   *jscpdTotal `json:"total"`
		Formats map[string]struct {
			// Older jscpd writes per-file counts here; jscpd 5 writes a plain
			// source count, which carries no per-file lines.
			Sources json.RawMessage `json:"sources"`
		} `json:"formats"`
	} `json:"statistics"`
	Duplicates []struct {
		Lines      int           `json:"lines"`
		FirstFile  jscpdLocation `json:"firstFile"`
		SecondFile jscpdLocation `json:"secondFile"`
	} `json:"duplicates"`
}

// Read implements Reader. The headline is the report's own total; a report
// without totals, or whose totals disagree with themselves, is unusable. A
// clone with a file outside the root is left out of the affected items and
// makes the reading partial; one outside the instance scope is left out quietly.
// With an instance scope the totals cover the whole project, so the headline is
// rebuilt from the per-file counts of files inside the scope; a report without
// them cannot be scoped and gives a partial reading with no value.
func (JscpdReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	var report jscpdReport
	if err := json.Unmarshal(in.Data, &report); err != nil {
		return Reading{}, fmt.Errorf("jscpd JSON is not a jscpd report: %w", err)
	}
	if err := validateJscpdTotals(report); err != nil {
		return Reading{}, err
	}
	t := report.Statistics.Total

	scope := in.inputScope()
	scoped := len(scope.Include) > 0 || len(scope.Exclude) > 0
	version := report.Version
	if version == "" {
		version = unknownVersion
	}
	prov := Provenance{ProviderVersion: version, MeasurementDefinition: duplicationDefinition}
	lines, duplicated, sources := *t.Lines, *t.DuplicatedLines, t.Sources
	percentage := *t.Percentage
	// Without per-file counts (jscpd 5) the report's own total stands: the
	// tool was already given the exclusions and target, so it is the best
	// measurement available.
	rebuilt := false
	if scoped {
		var found bool
		var l, d, n int
		l, d, n, found = jscpdScopedCounts(in, report)
		if found {
			if l == 0 {
				return Reading{}, errors.New("jscpd scanned no lines inside the instance scope, so there is nothing to measure")
			}
			lines, duplicated, sources, rebuilt = l, d, n, true
			percentage = percent(duplicated, lines)
		}
	}

	evidence, outside, clones, err := jscpdCloneEvidence(ctx, in, report, scope)
	if err != nil {
		return Reading{}, err
	}
	if !rebuilt {
		clones = t.Clones
	}

	reason := skippedFilesReason(outside, "clone(s) in files outside the project")
	return Reading{
		Partial:    outside > 0,
		Reason:     reason,
		Value:      &Value{Number: percentage, Unit: UnitPercent},
		Provenance: prov,
		Details: []Detail{
			{Key: DetailDuplicatedLines, Number: float64(duplicated)},
			{Key: DetailTotalLines, Number: float64(lines)},
			{Key: DetailClones, Number: float64(clones)},
			{Key: DetailSources, Number: float64(sources)},
		},
		Evidence: WorstEvidence(evidence),
	}, nil
}

// validateJscpdTotals refuses a report without totals, or whose totals disagree
// with themselves.
func validateJscpdTotals(report jscpdReport) error {
	if report.Statistics == nil || report.Statistics.Total == nil {
		return errors.New("jscpd JSON has no statistics total")
	}
	t := report.Statistics.Total
	if t.Lines == nil || t.DuplicatedLines == nil || t.Percentage == nil {
		return errors.New("jscpd JSON total lacks lines, duplicatedLines, or percentage")
	}
	if *t.Lines <= 0 {
		return errors.New("jscpd scanned no lines, so there is nothing to measure")
	}
	if *t.DuplicatedLines < 0 || *t.DuplicatedLines > *t.Lines {
		return errors.New("jscpd JSON totals are inconsistent")
	}
	if want := percent(*t.DuplicatedLines, *t.Lines); math.Abs(want-*t.Percentage) > jscpdPercentTolerance {
		return fmt.Errorf("jscpd JSON percentage %v does not match %d of %d lines", *t.Percentage, *t.DuplicatedLines, *t.Lines)
	}
	return nil
}

// jscpdCloneEvidence turns the report's clones into evidence for files inside
// the scope. It also counts clones that touch a file outside the project, which
// are left out, and the clones it kept.
func jscpdCloneEvidence(ctx context.Context, in ReportInput, report jscpdReport, scope InputScope) ([]RankedEvidence, int, int, error) {
	var evidence []RankedEvidence
	var outside, clones int
	for _, c := range report.Duplicates {
		if err := ctx.Err(); err != nil {
			return nil, 0, 0, err
		}
		first, ok1 := RelPath(in.Root, jscpdFileName(c.FirstFile.Name))
		second, ok2 := RelPath(in.Root, jscpdFileName(c.SecondFile.Name))
		if !ok1 || !ok2 {
			outside++
			continue
		}
		// Evidence names a file inside the scope. The other side of the clone may
		// lie outside it and stays in the note, since that is where the code repeats.
		here, there, hereLine, thereLine := first, second, c.FirstFile.Start, c.SecondFile.Start
		if !scope.relevant(first) {
			here, there, hereLine, thereLine = second, first, c.SecondFile.Start, c.FirstFile.Start
		}
		if !scope.relevant(here) {
			continue
		}
		clones++
		evidence = append(evidence, RankedEvidence{
			Ref: EvidenceRef{
				Path: here,
				Line: hereLine,
				Note: fmt.Sprintf("%d lines repeated at %s:%d", c.Lines, there, thereLine),
			},
			Rank: float64(c.Lines),
		})
	}
	return evidence, outside, clones, nil
}

// jscpdFileName drops the ":format" suffix jscpd 5 adds to a clone found in a
// code block inside another file, such as "docs/a.md:markdown".
func jscpdFileName(name string) string {
	i := strings.LastIndex(name, ":")
	if i <= 0 || i == len(name)-1 || strings.ContainsAny(name[i+1:], `/\:`) {
		return name
	}
	return name[:i]
}

// jscpdScopedCounts sums the per-file line counts of files inside the instance
// scope. found is false when the report carries no per-file counts at all.
func jscpdScopedCounts(in ReportInput, report jscpdReport) (lines, duplicated, sources int, found bool) {
	scope := in.inputScope()
	for _, format := range report.Statistics.Formats {
		var files map[string]jscpdSource
		if json.Unmarshal(format.Sources, &files) != nil {
			continue
		}
		for name, src := range files {
			found = true
			if rel, ok := RelPath(in.Root, name); ok && scope.relevant(rel) {
				lines += src.Lines
				duplicated += src.DuplicatedLines
				sources++
			}
		}
	}
	return lines, duplicated, sources, found
}
