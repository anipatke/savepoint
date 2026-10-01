package codehealth

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
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

// Read implements Reader.
func (JUnitReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	dec := xml.NewDecoder(bytes.NewReader(in.Data))
	var total, skipped, failures, errs int
	var evidence []RankedEvidence
	sawRoot := false
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
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		switch start.Name.Local {
		case "testsuites", "testsuite":
			sawRoot = true
		case "testcase":
			var c junitCase
			if err := dec.DecodeElement(&c, &start); err != nil {
				return Reading{}, fmt.Errorf("JUnit XML could not be read: %w", err)
			}
			total++
			switch {
			case c.Error != nil:
				errs++
				evidence = appendJUnitEvidence(evidence, in, c, "errored")
			case c.Failure != nil:
				failures++
				evidence = appendJUnitEvidence(evidence, in, c, "failed")
			case c.Skipped != nil:
				skipped++
			}
		}
	}
	if !sawRoot {
		return Reading{}, errors.New("JUnit XML has no testsuite element")
	}

	rd := Reading{
		Value:      &Value{Number: float64(failures + errs), Unit: UnitCount},
		Provenance: Provenance{ProviderVersion: unknownVersion, MeasurementDefinition: testsDefinitionJUnit},
		Details: []Detail{
			{Key: "total_tests", Number: float64(total)},
			{Key: "skipped_tests", Number: float64(skipped)},
			{Key: "errors", Number: float64(errs)},
		},
		Evidence: WorstEvidence(evidence),
	}
	if total == 0 {
		rd.Reason = noTestsRanReason
	}
	return rd, nil
}

func appendJUnitEvidence(items []RankedEvidence, in ReportInput, c junitCase, verb string) []RankedEvidence {
	p, ok := in.Path(junitCasePath(in.Provider, c))
	if !ok {
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
