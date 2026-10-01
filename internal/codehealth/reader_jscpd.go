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

type jscpdReport struct {
	Version    string `json:"version"`
	Statistics *struct {
		Total *jscpdTotal `json:"total"`
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

	var evidence []RankedEvidence
	var outside int
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
		if !in.inputScope().relevant(first) && !in.inputScope().relevant(second) {
			continue
		}
		evidence = append(evidence, RankedEvidence{
			Ref: EvidenceRef{
				Path: first,
				Line: c.FirstFile.Start,
				Note: fmt.Sprintf("%d lines repeated at %s:%d", c.Lines, second, c.SecondFile.Start),
			},
			Rank: float64(c.Lines),
		})
	}

	version := report.Version
	if version == "" {
		version = unknownVersion
	}
	reason := skippedFilesReason(outside, "clone(s) in files outside the project")
	return Reading{
		Partial:    outside > 0,
		Reason:     reason,
		Value:      &Value{Number: *t.Percentage, Unit: UnitPercent},
		Provenance: Provenance{ProviderVersion: version, MeasurementDefinition: duplicationDefinition},
		Details: []Detail{
			{Key: DetailDuplicatedLines, Number: float64(*t.DuplicatedLines)},
			{Key: DetailTotalLines, Number: float64(*t.Lines)},
			{Key: DetailClones, Number: float64(t.Clones)},
			{Key: DetailSources, Number: float64(t.Sources)},
		},
		Evidence: WorstEvidence(evidence),
	}, nil
}
