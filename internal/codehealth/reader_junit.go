package codehealth

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// JUnitReader reads JUnit XML as written by Vitest and pytest. Only test case
// results are read; hostnames and captured output are ignored.
type JUnitReader struct{}

type junitCase struct {
	Name      string    `xml:"name,attr"`
	Classname string    `xml:"classname,attr"`
	File      string    `xml:"file,attr"`
	Line      int       `xml:"line,attr"`
	Failure   *struct{} `xml:"failure"`
	Error     *struct{} `xml:"error"`
	Skipped   *struct{} `xml:"skipped"`
}

// junitSuite tracks one suite element: how many test cases it holds against
// the count it declares, when it declares one.
type junitSuite struct {
	declared int // -1 when the suite states no count
	seen     int
}

// junitTally accumulates what a report's test cases and suites add up to.
type junitTally struct {
	total, skipped, failures, errs, mismatched int
	evidence                                   []RankedEvidence
}

// add counts one test case. rel is its repository path and mapped says the path
// is known and inside the instance scope, so it can carry evidence.
func (t *junitTally) add(c junitCase, rel string, mapped bool) {
	t.total++
	switch {
	case c.Error != nil:
		t.errs++
		t.evidence = appendJUnitEvidence(t.evidence, rel, mapped, c, "errored")
	case c.Failure != nil:
		t.failures++
		t.evidence = appendJUnitEvidence(t.evidence, rel, mapped, c, "failed")
	case c.Skipped != nil:
		t.skipped++
	}
}

func (t *junitTally) reading() Reading {
	rd := Reading{
		Value:      &Value{Number: float64(t.failures + t.errs), Unit: UnitCount},
		Provenance: Provenance{ProviderVersion: unknownVersion, MeasurementDefinition: testsDefinitionJUnit},
		Details: []Detail{
			{Key: "total_tests", Number: float64(t.total)},
			{Key: "skipped_tests", Number: float64(t.skipped)},
			{Key: "errors", Number: float64(t.errs)},
		},
		Evidence: WorstEvidence(t.evidence),
	}
	switch {
	case t.mismatched > 0:
		rd.Partial = true
		rd.Reason = fmt.Sprintf("%d suite(s) declare a test count that does not match the test cases reported", t.mismatched)
	case t.total == 0:
		rd.Reason = noTestsRanReason
	}
	return rd
}

// Read implements Reader. Test cases the instance scope excludes are not
// counted; a case whose file cannot be placed in the repository cannot be
// judged, so it is. A report with several document roots is unusable, and one
// whose suites declare more or fewer tests than they hold is partial.
func (JUnitReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	dec := xml.NewDecoder(bytes.NewReader(in.Data))
	var tally junitTally
	var suites []*junitSuite
	depth, roots := 0, 0
	sawSuite := false
	for {
		if err := ctx.Err(); err != nil {
			return Reading{}, err
		}
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return Reading{}, fmt.Errorf("JUnit XML could not be read: %w", err)
		}
		switch el := tok.(type) {
		case xml.EndElement:
			depth--
			if n := len(suites); n > 0 && isJUnitSuite(el.Name.Local) {
				if s := suites[n-1]; s.declared >= 0 && s.declared != s.seen {
					tally.mismatched++
				}
				suites = suites[:n-1]
			}
		case xml.StartElement:
			if depth == 0 {
				if roots++; roots > 1 {
					return Reading{}, errors.New("JUnit XML has more than one document root")
				}
			}
			depth++
			switch {
			case isJUnitSuite(el.Name.Local):
				sawSuite = true
				declared, err := declaredTests(el)
				if err != nil {
					return Reading{}, err
				}
				suites = append(suites, &junitSuite{declared: declared})
			case el.Name.Local == "testcase":
				depth--
				if err := tally.readCase(dec, el, in, suites); err != nil {
					return Reading{}, err
				}
			}
		}
	}
	if !sawSuite {
		return Reading{}, errors.New("JUnit XML has no testsuite element")
	}
	return tally.reading(), nil
}

// readCase decodes one test case, counts it toward every open suite, and adds it
// to the tally unless the instance scope excludes its file.
func (t *junitTally) readCase(dec *xml.Decoder, el xml.StartElement, in ReportInput, suites []*junitSuite) error {
	var c junitCase
	if err := dec.DecodeElement(&c, &el); err != nil {
		return fmt.Errorf("JUnit XML could not be read: %w", err)
	}
	for _, s := range suites {
		s.seen++
	}
	rel, known := RelPath(in.Root, junitCasePath(in.Provider, c))
	mapped := known && in.inputScope().relevant(rel)
	if known && !mapped {
		return nil
	}
	t.add(c, rel, mapped)
	return nil
}

func isJUnitSuite(name string) bool { return name == "testsuites" || name == "testsuite" }

// declaredTests is the test count a suite states, or -1 when it states none.
func declaredTests(el xml.StartElement) (int, error) {
	for _, a := range el.Attr {
		if a.Name.Local != "tests" {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(a.Value))
		if err != nil || n < 0 {
			return 0, fmt.Errorf("JUnit XML declares %q tests, which is not a count", a.Value)
		}
		return n, nil
	}
	return -1, nil
}

func appendJUnitEvidence(items []RankedEvidence, p string, mapped bool, c junitCase, verb string) []RankedEvidence {
	if !mapped {
		return items
	}
	name := c.Name
	if name == "" {
		name = c.Classname
	}
	return append(items, RankedEvidence{Ref: EvidenceRef{Path: p, Line: max(c.Line, 0), Note: "test " + name + " " + verb}, Rank: 1})
}

// junitCasePath is the file a test case lives in: its file attribute when the
// runner wrote one, otherwise one derived from the class name. Vitest writes a
// path as the class name; pytest writes a dotted module path, optionally ending
// in test class names, which are dropped.
func junitCasePath(provider ProviderKey, c junitCase) string {
	if c.File != "" {
		return c.File
	}
	cn := c.Classname
	switch {
	case cn == "":
		return ""
	case provider != ProviderPytestJUnit:
		return cn
	case strings.Contains(cn, "/"):
		return cn
	}
	parts := strings.Split(cn, ".")
	for len(parts) > 1 && startsUpper(parts[len(parts)-1]) {
		parts = parts[:len(parts)-1]
	}
	return strings.Join(parts, "/") + ".py"
}

func startsUpper(s string) bool {
	for _, r := range s {
		return unicode.IsUpper(r)
	}
	return false
}
