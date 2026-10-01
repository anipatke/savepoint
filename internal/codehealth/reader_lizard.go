package codehealth

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// Reader for the complexity capability. Lizard writes a headerless CSV with one
// row per function; the value is the highest cyclomatic complexity (CCN), and
// the most complex functions are the affected items.

const (
	complexityDefinition = "highest cyclomatic complexity (CCN) of any function"

	// Lizard CSV columns, by position: NLOC, CCN, token, PARAM, length,
	// location, file, function, long name, start, end.
	lizardColumns  = 11
	lizardCCN      = 1
	lizardFile     = 6
	lizardFunction = 7
	lizardStart    = 9

	DetailFunctions       = "functions"
	DetailAverageCCN      = "average_ccn"
	DetailFunctionsOver10 = "functions_over_10"
	DetailFunctionsOver20 = "functions_over_20"
)

// LizardReader reads Lizard CSV reports.
type LizardReader struct{}

// lizardFunc is one function row that falls inside the root and the scope.
type lizardFunc struct {
	file string
	name string
	line int
	ccn  int
}

// Read implements Reader. A row with the wrong column count or a number that
// does not parse makes the whole report unusable. A function outside the root
// is dropped and makes the reading partial; one outside the instance scope is
// dropped quietly, as scope is the project's own choice.
func (LizardReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	funcs, outside, err := decodeLizard(ctx, in)
	if err != nil {
		return Reading{}, err
	}
	reason := skippedFilesReason(outside, "function(s) in files outside the project")
	prov := Provenance{
		ProviderVersion:       unknownVersion,
		MeasurementDefinition: complexityDefinition,
	}
	if len(funcs) == 0 {
		if outside > 0 {
			// Nothing was measured, so there is no complexity to report, not zero.
			return Reading{Partial: true, Reason: reason, Provenance: prov}, nil
		}
		return Reading{
			Reason:     "Lizard reported no functions; Lizard's CSV does not state its version",
			Value:      &Value{Number: 0, Unit: UnitCCN},
			Provenance: prov,
			Details:    complexityDetails(nil),
		}, nil
	}

	highest := 0
	for _, f := range funcs {
		highest = max(highest, f.ccn)
	}
	if reason == "" {
		reason = "Lizard's CSV does not state its version"
	} else {
		reason += "; Lizard's CSV does not state its version"
	}
	return Reading{
		Partial:    outside > 0,
		Reason:     reason,
		Value:      &Value{Number: float64(highest), Unit: UnitCCN},
		Provenance: prov,
		Details:    complexityDetails(funcs),
		Evidence:   mostComplex(funcs),
	}, nil
}

func decodeLizard(ctx context.Context, in ReportInput) (funcs []lizardFunc, outside int, err error) {
	cr := csv.NewReader(bytes.NewReader(in.Data))
	cr.FieldsPerRecord = lizardColumns
	scope := in.inputScope()
	for {
		if err := ctx.Err(); err != nil {
			return nil, 0, err
		}
		row, err := cr.Read()
		if errors.Is(err, io.EOF) {
			return funcs, outside, nil
		}
		if err != nil {
			return nil, 0, fmt.Errorf("lizard csv: %w", err)
		}
		line, _ := cr.FieldPos(0)
		ccn, err := strconv.Atoi(row[lizardCCN])
		if err != nil || ccn < 0 {
			return nil, 0, fmt.Errorf("lizard csv line %d: CCN %q is not a whole number", line, row[lizardCCN])
		}
		start, err := strconv.Atoi(row[lizardStart])
		if err != nil || start < 0 {
			return nil, 0, fmt.Errorf("lizard csv line %d: start line %q is not a whole number", line, row[lizardStart])
		}
		rel, ok := RelPath(in.Root, row[lizardFile])
		if !ok {
			outside++
			continue
		}
		if !scope.relevant(rel) {
			continue
		}
		funcs = append(funcs, lizardFunc{file: rel, name: row[lizardFunction], line: start, ccn: ccn})
	}
}

// complexityDetails reports how many functions there are, their mean CCN, and
// how many are over the Good and Watch limits of the default thresholds.
func complexityDetails(funcs []lizardFunc) []Detail {
	var sum, over10, over20 int
	limits := defaultThresholds[CapabilityComplexity]
	for _, f := range funcs {
		sum += f.ccn
		if float64(f.ccn) > limits.Good {
			over10++
		}
		if float64(f.ccn) > limits.Watch {
			over20++
		}
	}
	average := 0.0
	if len(funcs) > 0 {
		average = float64(sum) / float64(len(funcs))
	}
	return []Detail{
		{Key: DetailFunctions, Number: float64(len(funcs))},
		{Key: DetailAverageCCN, Number: average},
		{Key: DetailFunctionsOver10, Number: float64(over10)},
		{Key: DetailFunctionsOver20, Number: float64(over20)},
	}
}

// mostComplex names the functions with the highest CCN. Functions within the
// Good limit are not affected items.
func mostComplex(funcs []lizardFunc) []EvidenceRef {
	good := defaultThresholds[CapabilityComplexity].Good
	var evidence []RankedEvidence
	for _, f := range funcs {
		if float64(f.ccn) <= good {
			continue
		}
		evidence = append(evidence, RankedEvidence{
			Ref:  EvidenceRef{Path: f.file, Line: f.line, Note: fmt.Sprintf("%s has complexity %d", f.name, f.ccn)},
			Rank: float64(f.ccn),
		})
	}
	return WorstEvidence(evidence)
}
