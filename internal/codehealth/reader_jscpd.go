package codehealth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
			Sources map[string]jscpdSource `json:"sources"`
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
	if report.Statistics == nil || report.Statistics.Total == nil {
		return Reading{}, errors.New("jscpd JSON has no statistics total")
	}
	t := report.Statistics.Total
	if t.Lines == nil || t.DuplicatedLines == nil || t.Percentage == nil {
		return Reading{}, errors.New("jscpd JSON total lacks lines, duplicatedLines, or percentage")
	}
	if *t.Lines <= 0 {
		return Reading{}, errors.New("jscpd scanned no lines, so there is nothing to measure")
	}
	if *t.DuplicatedLines < 0 || *t.DuplicatedLines > *t.Lines {
		return Reading{}, errors.New("jscpd JSON totals are inconsistent")
	}
	if want := percent(*t.DuplicatedLines, *t.Lines); math.Abs(want-*t.Percentage) > jscpdPercentTolerance {
		return Reading{}, fmt.Errorf("jscpd JSON percentage %v does not match %d of %d lines", *t.Percentage, *t.DuplicatedLines, *t.Lines)
	}

	scope := in.inputScope()
	scoped := len(scope.Include) > 0 || len(scope.Exclude) > 0
	version := report.Version
	if version == "" {
		version = unknownVersion
	}
	prov := Provenance{ProviderVersion: version, MeasurementDefinition: duplicationDefinition}
	lines, duplicated, sources := *t.Lines, *t.DuplicatedLines, t.Sources
	percentage := *t.Percentage
	if scoped {
		var found bool
		lines, duplicated, sources, found = jscpdScopedCounts(in, report)
		if !found {
			return Reading{
				Partial:    true,
				Reason:     "jscpd JSON has no per-file line counts, so its totals cannot be limited to the instance scope",
				Provenance: prov,
			}, nil
		}
		if lines == 0 {
			return Reading{}, errors.New("jscpd scanned no lines inside the instance scope, so there is nothing to measure")
		}
		percentage = percent(duplicated, lines)
	}

	var evidence []RankedEvidence
	var outside, clones int
	for _, c := range report.Duplicates {
		if err := ctx.Err(); err != nil {
			return Reading{}, err
		}
		first, ok1 := RelPath(in.Root, c.FirstFile.Name)
		second, ok2 := RelPath(in.Root, c.SecondFile.Name)
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
	if !scoped {
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

// jscpdScopedCounts sums the per-file line counts of files inside the instance
// scope. found is false when the report carries no per-file counts at all.
func jscpdScopedCounts(in ReportInput, report jscpdReport) (lines, duplicated, sources int, found bool) {
	scope := in.inputScope()
	for _, format := range report.Statistics.Formats {
		for name, src := range format.Sources {
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
