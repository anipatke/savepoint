package codehealth

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
)

// Readers for the coverage capability. The value is statement coverage on every
// stack, so a Go profile, a Vitest report, and a coverage.py report mean the
// same thing. Branch and function counts are supporting details where the
// report has them, and the least-covered files are the affected items.

const (
	coverageDefinition = "covered statements as a share of all statements"
	goProfileMode      = "mode:"
	maxGoProfileLine   = 1 << 20

	DetailCoveredStatements = "covered_statements"
	DetailTotalStatements   = "total_statements"
	DetailCoveredFunctions  = "covered_functions"
	DetailTotalFunctions    = "total_functions"
	DetailCoveredBranches   = "covered_branches"
	DetailTotalBranches     = "total_branches"
)

// stmtCount is how many statements a file has and how many ran.
type stmtCount struct{ covered, total int }

// statementReading turns per-file statement counts into a Reading. Files with
// no statements are ignored. No statements at all is an error: the report
// measured nothing, which is neither 0% nor 100%. extra carries supporting
// details, and a non-empty partialReason marks the reading partial.
func statementReading(files map[string]stmtCount, version, partialReason string, extra []Detail) (Reading, error) {
	var covered, total int
	for _, f := range files {
		covered += f.covered
		total += f.total
	}
	if total == 0 {
		return Reading{}, errors.New("coverage report has no statements to measure")
	}
	details := append([]Detail{
		{Key: DetailCoveredStatements, Number: float64(covered)},
		{Key: DetailTotalStatements, Number: float64(total)},
	}, extra...)
	return Reading{
		Partial:    partialReason != "",
		Reason:     partialReason,
		Value:      &Value{Number: percent(covered, total), Unit: UnitPercent},
		Provenance: Provenance{ProviderVersion: version, MeasurementDefinition: coverageDefinition},
		Details:    details,
		Evidence:   leastCovered(files),
	}, nil
}

// leastCovered names the files with the lowest statement coverage. Files that
// are fully covered, or have no statements, are not affected items.
func leastCovered(files map[string]stmtCount) []EvidenceRef {
	var evidence []RankedEvidence
	for name, f := range files {
		if f.total == 0 || f.covered >= f.total {
			continue
		}
		pct := percent(f.covered, f.total)
		evidence = append(evidence, RankedEvidence{
			Ref:  EvidenceRef{Path: name, Note: fmt.Sprintf("%.1f%% of %d statements covered", pct, f.total)},
			Rank: 100 - pct,
		})
	}
	return WorstEvidence(evidence)
}

func percent(covered, total int) float64 { return float64(covered) * 100 / float64(total) }

// skippedFilesReason says how many report files could not be placed inside the
// repository, or "" when none were.
func skippedFilesReason(n int, what string) string {
	if n == 0 {
		return ""
	}
	return fmt.Sprintf("%d %s could not be mapped to repository paths", n, what)
}

// GoCoverReader reads Go cover profiles.
type GoCoverReader struct{}

// goBlock is one profile row after merging: the same block listed by several
// packages' test runs is one block, covered when any run covered it.
type goBlock struct {
	file    string
	stmts   int
	covered bool
}

// Read implements Reader.
func (GoCoverReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	blocks, err := decodeGoProfile(ctx, in.Data)
	if err != nil {
		return Reading{}, err
	}
	// Resolution ignores the instance scope so that a file outside it is dropped
	// quietly; only a file no module owns makes the reading partial.
	modules := LoadGoModules(in.Root, InputScope{})
	files := map[string]stmtCount{}
	unresolved := map[string]bool{}
	for _, b := range blocks {
		dir, ok := modules.Resolve(path.Dir(b.file))
		if !ok {
			unresolved[b.file] = true
			continue
		}
		rel, ok := in.Path(path.Join(dir, path.Base(b.file)))
		if !ok {
			continue
		}
		f := files[rel]
		f.total += b.stmts
		if b.covered {
			f.covered += b.stmts
		}
		files[rel] = f
	}
	return statementReading(files, unknownVersion, skippedFilesReason(len(unresolved), "profile file(s)"), nil)
}

func decodeGoProfile(ctx context.Context, data []byte) ([]goBlock, error) {
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(nil, maxGoProfileLine)
	seenMode := false
	index := map[string]int{}
	var blocks []goBlock
	for n := 1; sc.Scan(); n++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if !seenMode {
			if !strings.HasPrefix(line, goProfileMode) {
				return nil, errors.New("go cover profile must start with a mode: line")
			}
			seenMode = true
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("go cover profile line %d is not a coverage block", n)
		}
		colon := strings.LastIndex(fields[0], ":")
		stmts, err1 := strconv.Atoi(fields[1])
		count, err2 := strconv.ParseInt(fields[2], 10, 64)
		if colon <= 0 || !strings.Contains(fields[0][colon:], ",") || err1 != nil || err2 != nil || stmts < 0 || count < 0 {
			return nil, fmt.Errorf("go cover profile line %d is not a coverage block", n)
		}
		if i, dup := index[fields[0]]; dup {
			blocks[i].covered = blocks[i].covered || count > 0
			continue
		}
		index[fields[0]] = len(blocks)
		blocks = append(blocks, goBlock{file: fields[0][:colon], stmts: stmts, covered: count > 0})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("go cover profile could not be read: %w", err)
	}
	if !seenMode {
		return nil, errors.New("go cover profile is empty")
	}
	return blocks, nil
}

// VitestCoverageReader reads Vitest V8 coverage-final.json (Istanbul format).
type VitestCoverageReader struct{}

type istanbulFile struct {
	S map[string]int64   `json:"s"`
	F map[string]int64   `json:"f"`
	B map[string][]int64 `json:"b"`
}

// Read implements Reader.
func (VitestCoverageReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	var report map[string]istanbulFile
	if err := json.Unmarshal(in.Data, &report); err != nil {
		return Reading{}, fmt.Errorf("vitest coverage JSON is not a coverage-final.json report: %w", err)
	}
	names := make([]string, 0, len(report))
	for name := range report {
		names = append(names, name)
	}
	sort.Strings(names)

	files := map[string]stmtCount{}
	var coveredFns, totalFns, coveredBranches, totalBranches, unmapped int
	for _, name := range names {
		if err := ctx.Err(); err != nil {
			return Reading{}, err
		}
		rel, ok := RelPath(in.Root, name)
		if !ok {
			unmapped++
			continue
		}
		if !in.inputScope().relevant(rel) {
			continue
		}
		f := report[name]
		sc := files[rel]
		for _, n := range f.S {
			sc.total++
			if n > 0 {
				sc.covered++
			}
		}
		files[rel] = sc
		for _, n := range f.F {
			totalFns++
			if n > 0 {
				coveredFns++
			}
		}
		for _, arms := range f.B {
			for _, n := range arms {
				totalBranches++
				if n > 0 {
					coveredBranches++
				}
			}
		}
	}
	return statementReading(files, unknownVersion, skippedFilesReason(unmapped, "report file(s)"), []Detail{
		{Key: DetailCoveredFunctions, Number: float64(coveredFns)},
		{Key: DetailTotalFunctions, Number: float64(totalFns)},
		{Key: DetailCoveredBranches, Number: float64(coveredBranches)},
		{Key: DetailTotalBranches, Number: float64(totalBranches)},
	})
}

// CoveragePyReader reads coverage.py JSON reports.
type CoveragePyReader struct{}

type coveragePySummary struct {
	CoveredLines    *int `json:"covered_lines"`
	NumStatements   *int `json:"num_statements"`
	NumBranches     int  `json:"num_branches"`
	CoveredBranches int  `json:"covered_branches"`
}

type coveragePyReport struct {
	Meta struct {
		Version        string `json:"version"`
		BranchCoverage bool   `json:"branch_coverage"`
	} `json:"meta"`
	Files map[string]struct {
		Summary coveragePySummary `json:"summary"`
	} `json:"files"`
	Totals coveragePySummary `json:"totals"`
}

// Read implements Reader. The value comes from the report's own totals, which
// count statements, not physical lines.
func (CoveragePyReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	var report coveragePyReport
	if err := json.Unmarshal(in.Data, &report); err != nil {
		return Reading{}, fmt.Errorf("coverage.py JSON is not a coverage report: %w", err)
	}
	t := report.Totals
	if t.CoveredLines == nil || t.NumStatements == nil {
		return Reading{}, errors.New("coverage.py JSON has no totals for statements")
	}
	if *t.NumStatements <= 0 {
		return Reading{}, errors.New("coverage report has no statements to measure")
	}
	if *t.CoveredLines < 0 || *t.CoveredLines > *t.NumStatements {
		return Reading{}, errors.New("coverage.py JSON totals are inconsistent")
	}

	// The files only rank the affected items; they never change the headline.
	files := map[string]stmtCount{}
	for name, f := range report.Files {
		if f.Summary.CoveredLines == nil || f.Summary.NumStatements == nil {
			continue
		}
		if rel, ok := in.Path(name); ok {
			files[rel] = stmtCount{covered: *f.Summary.CoveredLines, total: *f.Summary.NumStatements}
		}
	}
	version := report.Meta.Version
	if version == "" {
		version = unknownVersion
	}
	rd := Reading{
		Value:      &Value{Number: percent(*t.CoveredLines, *t.NumStatements), Unit: UnitPercent},
		Provenance: Provenance{ProviderVersion: version, MeasurementDefinition: coverageDefinition},
		Details: []Detail{
			{Key: DetailCoveredStatements, Number: float64(*t.CoveredLines)},
			{Key: DetailTotalStatements, Number: float64(*t.NumStatements)},
		},
		Evidence: leastCovered(files),
	}
	if report.Meta.BranchCoverage {
		rd.Details = append(rd.Details,
			Detail{Key: DetailCoveredBranches, Number: float64(t.CoveredBranches)},
			Detail{Key: DetailTotalBranches, Number: float64(t.NumBranches)})
	}
	return rd, nil
}
