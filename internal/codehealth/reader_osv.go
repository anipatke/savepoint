package codehealth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

// Reader for the dependency-vulnerability capability. The value is the number
// of vulnerability groups OSV-Scanner reports, one per package; each group's
// worst CVSS score places it in a severity bucket. Matching vulnerabilities to
// packages is the scanner's work, never this reader's.

const (
	vulnerabilityDefinition = "vulnerability groups per package, bucketed by CVSS"

	// CVSS score floors of the severity buckets.
	cvssCritical = 9.0
	cvssHigh     = 7.0
	cvssMedium   = 4.0
	cvssMax      = 10.0

	// An unknown-severity group blocks like a high one, so it ranks with it.
	unknownSeverityRank = cvssHigh

	maxAdvisoryIDs = 3
)

// OSVScannerReader reads OSV-Scanner JSON reports.
type OSVScannerReader struct{}

type osvPackage struct {
	Package struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"package"`
	Vulnerabilities []json.RawMessage `json:"vulnerabilities"`
	Groups          []struct {
		IDs         []string `json:"ids"`
		MaxSeverity string   `json:"max_severity"`
	} `json:"groups"`
	Error string `json:"error"`
}

type osvReport struct {
	Version          string `json:"version"`
	DatabaseSnapshot string `json:"database_snapshot"`
	Results          []struct {
		Source struct {
			Path string `json:"path"`
		} `json:"source"`
		Packages []osvPackage `json:"packages"`
		Error    string       `json:"error"`
	} `json:"results"`
}

// osvTally counts what could not be read exactly, so the reading can say so.
type osvTally struct {
	outside, unresolved, errored, ungrouped int
}

// osvScan accumulates the groups, evidence and limitations of one report.
type osvScan struct {
	counts   map[string]int
	total    int
	tally    osvTally
	evidence []RankedEvidence
}

// addPackage counts one package's groups. rel is its source's repository path;
// a source outside the project (inside false) cannot be named, so its groups
// count but get no affected item.
func (o *osvScan) addPackage(p osvPackage, rel string, inside bool) {
	switch {
	case p.Error != "":
		o.tally.errored++
	case p.Package.Version == "":
		o.tally.unresolved++
	}
	if len(p.Groups) == 0 && len(p.Vulnerabilities) > 0 {
		o.tally.ungrouped++
	}
	seen := map[string]bool{}
	for _, g := range p.Groups {
		// A group repeated within a package, even with its IDs reordered,
		// is one vulnerability group, not two.
		if key := osvGroupKey(g.IDs); key != "" {
			if seen[key] {
				continue
			}
			seen[key] = true
		}
		bucket, rank := osvBucket(g.MaxSeverity)
		o.counts[bucket]++
		o.total++
		if inside {
			o.evidence = append(o.evidence, RankedEvidence{
				Ref:  EvidenceRef{Path: rel, Note: osvNote(p, g.IDs)},
				Rank: rank,
			})
		}
	}
}

// Read implements Reader. Every group of a source inside the instance scope is
// counted. A source outside the project cannot be named, so its groups still
// count but get no affected item; a source scanned without resolved versions,
// with a scanner error, or whose vulnerabilities were not grouped makes the
// reading partial.
func (OSVScannerReader) Read(ctx context.Context, in ReportInput) (Reading, error) {
	var report osvReport
	if err := json.Unmarshal(in.Data, &report); err != nil {
		return Reading{}, fmt.Errorf("OSV-Scanner JSON is not an OSV-Scanner report: %w", err)
	}
	if report.Results == nil {
		return Reading{}, errors.New("OSV-Scanner JSON has no results list")
	}

	scan := osvScan{counts: map[string]int{}}
	for _, r := range report.Results {
		if err := ctx.Err(); err != nil {
			return Reading{}, err
		}
		rel, inside := RelPath(in.Root, r.Source.Path)
		if inside && !in.inputScope().relevant(rel) {
			continue
		}
		if !inside {
			scan.tally.outside++
		}
		if r.Error != "" {
			scan.tally.errored++
		}
		for _, p := range r.Packages {
			scan.addPackage(p, rel, inside)
		}
	}
	counts, total, tally, evidence := scan.counts, scan.total, scan.tally, scan.evidence

	snapshot := unknownVersion
	if s := sanitizeLine(report.DatabaseSnapshot); s != "" {
		snapshot = s
	}
	version := report.Version
	if version == "" {
		version = unknownVersion
	}
	reasons := tally.reasons()
	partial := len(reasons) > 0
	reasons = append(reasons, "Vulnerability database was queried at collection time; its snapshot is "+snapshot+".")
	details := make([]Detail, 0, len(vulnerabilitySeverityKeys))
	for _, key := range vulnerabilitySeverityKeys {
		details = append(details, Detail{Key: key, Number: float64(counts[key])})
	}
	return Reading{
		Partial:    partial,
		Reason:     boundText(strings.Join(reasons, "; ")),
		Value:      &Value{Number: float64(total), Unit: UnitCount},
		Provenance: Provenance{ProviderVersion: version, MeasurementDefinition: vulnerabilityDefinition},
		Details:    details,
		Evidence:   WorstEvidence(evidence),
	}, nil
}

// osvGroupKey identifies a group by its advisory IDs regardless of order, or
// "" for a group that names none and so cannot be told from another.
func osvGroupKey(ids []string) string {
	sorted := slices.Clone(ids)
	slices.Sort(sorted)
	return strings.Join(slices.Compact(sorted), "\x00")
}

func (t osvTally) reasons() []string {
	var out []string
	for _, c := range []struct {
		n    int
		what string
	}{
		{t.outside, "source(s) outside the project have no affected items"},
		{t.unresolved, "package(s) scanned without a resolved version"},
		{t.errored, "scanner error(s) in the report"},
		{t.ungrouped, "package(s) with vulnerabilities but no groups"},
	} {
		if c.n > 0 {
			out = append(out, fmt.Sprintf("%d %s", c.n, c.what))
		}
	}
	return out
}

// osvBucket places a group by its max_severity CVSS score and ranks it for
// worst-first ordering. A missing, unparsable, or non-positive score is unknown.
func osvBucket(maxSeverity string) (string, float64) {
	score, err := strconv.ParseFloat(strings.TrimSpace(maxSeverity), 64)
	switch {
	case err != nil || math.IsNaN(score) || score <= 0 || score > cvssMax:
		return DetailUnknownVulnerabilities, unknownSeverityRank
	case score >= cvssCritical:
		return DetailCriticalVulnerabilities, score
	case score >= cvssHigh:
		return DetailHighVulnerabilities, score
	case score >= cvssMedium:
		return DetailMediumVulnerabilities, score
	}
	return DetailLowVulnerabilities, score
}

// osvNote names the package, its version, and the first few advisory IDs.
func osvNote(p osvPackage, ids []string) string {
	version := p.Package.Version
	if version == "" {
		version = unknownVersion
	}
	more := ""
	if len(ids) > maxAdvisoryIDs {
		more = fmt.Sprintf(" +%d more", len(ids)-maxAdvisoryIDs)
		ids = ids[:maxAdvisoryIDs]
	}
	return fmt.Sprintf("%s %s: %s%s", p.Package.Name, version, strings.Join(ids, ", "), more)
}
