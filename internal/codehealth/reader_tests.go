package codehealth

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
)

// Readers for the tests capability. The value is the failed-test count; the
// totals, skips, and errors are supporting details, and failing tests are named
// as affected items. Raw test output is never read into a result.

const (
	testsDefinitionGo    = "failed tests plus packages that failed to build"
	testsDefinitionJUnit = "failed tests plus errored tests"
	noTestsRanReason     = "no tests ran"
	unknownVersion       = "unknown"
	maxGoEventBytes      = 1 << 20
)

// GoTestReader reads `go test -json` event streams.
type GoTestReader struct{}

type goTestEvent struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
}

type goPackageRun struct {
	started map[string]bool   // tests that began
	results map[string]string // test -> pass, fail, or skip
	final   string            // the package's own pass, fail, or skip action
}

// Read implements Reader.
func (GoTestReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	packages, order, err := decodeGoTestEvents(ctx, in.Data)
	if err != nil {
		return Reading{}, err
	}
	modules := LoadGoModules(in.Root, in.inputScope())

	var total, passed, skipped, failed, buildFailures int
	var evidence []RankedEvidence
	unfinished := 0
	note := func(pkg, text string) {
		if dir, ok := modules.Resolve(pkg); ok {
			evidence = append(evidence, RankedEvidence{Ref: EvidenceRef{Path: dir, Note: text}, Rank: 1})
		}
	}
	for _, name := range order {
		p := packages[name]
		pkgFailed := 0
		tests := make([]string, 0, len(p.started))
		for t := range p.started {
			tests = append(tests, t)
		}
		sort.Strings(tests)
		for _, t := range tests {
			switch p.results[t] {
			case "pass":
				total++
				passed++
			case "skip":
				total++
				skipped++
			case "fail":
				total++
				failed++
				pkgFailed++
				note(name, "test "+t+" failed")
			default:
				// Started but never finished: a cut-off stream, or a panic or
				// timeout that killed the package.
				total++
				if p.final == "" {
					unfinished++
				} else {
					failed++
					pkgFailed++
					note(name, "test "+t+" did not finish")
				}
			}
		}
		switch {
		case p.final == "" && len(tests) > 0:
			// Counted above as unfinished.
		case p.final == "fail" && pkgFailed == 0:
			buildFailures++
			failed++
			note(name, "package "+name+" failed to build")
		}
	}

	rd := Reading{
		Value:      &Value{Number: float64(failed), Unit: UnitCount},
		Provenance: Provenance{ProviderVersion: unknownVersion, MeasurementDefinition: testsDefinitionGo},
		Details: []Detail{
			{Key: "total_tests", Number: float64(total)},
			{Key: "passed_tests", Number: float64(passed)},
			{Key: "skipped_tests", Number: float64(skipped)},
			{Key: "build_failures", Number: float64(buildFailures)},
		},
		Evidence: WorstEvidence(evidence),
	}
	switch {
	case unfinished > 0:
		rd.Partial = true
		rd.Reason = fmt.Sprintf("the test stream ended before %d started test(s) finished", unfinished)
	case total == 0 && buildFailures == 0:
		rd.Reason = noTestsRanReason
	}
	return rd, nil
}

func decodeGoTestEvents(ctx context.Context, data []byte) (map[string]*goPackageRun, []string, error) {
	packages := map[string]*goPackageRun{}
	var order []string
	pkg := func(name string) *goPackageRun {
		p, ok := packages[name]
		if !ok {
			p = &goPackageRun{started: map[string]bool{}, results: map[string]string{}}
			packages[name] = p
			order = append(order, name)
		}
		return p
	}

	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(nil, maxGoEventBytes)
	lines := 0
	for sc.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, nil, err
		}
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		var ev goTestEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			return nil, nil, fmt.Errorf("go test JSON line %d is not an event: %w", lines+1, err)
		}
		lines++
		if ev.Package == "" {
			continue
		}
		p := pkg(ev.Package)
		switch ev.Action {
		case "run":
			if ev.Test != "" {
				p.started[ev.Test] = true
			}
		case "pass", "fail", "skip":
			if ev.Test == "" {
				p.final = ev.Action
			} else {
				p.started[ev.Test] = true
				p.results[ev.Test] = ev.Action
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, nil, fmt.Errorf("go test JSON could not be read: %w", err)
	}
	if lines == 0 {
		return nil, nil, errors.New("go test JSON report is empty")
	}
	return packages, order, nil
}
