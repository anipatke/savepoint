package codehealth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Reading is what a Reader measured from one report. Normalizing real tool
// reports is the readers' work; Collect only assembles and classifies.
type Reading struct {
	// Partial marks a measurement that does not cover everything; Reason must
	// then say what is missing.
	Partial    bool
	Reason     string
	Value      *Value
	Provenance Provenance
	Details    []Detail
	Evidence   []EvidenceRef
}

// ReportInput is the report bytes handed to a Reader, with the scope the
// instance is configured to measure.
type ReportInput struct {
	Provider   ProviderKey
	Data       []byte
	Scope      []string
	Exclusions []string
}

// Reader turns one provider's report into a measured Reading. An error means
// the report was malformed or unusable; it is recorded as a failure of that
// instance only.
type Reader interface {
	Read(ctx context.Context, in ReportInput) (Reading, error)
}

// Readers registers the Reader of each provider. A provider without one is
// reported as unsupported.
type Readers map[ProviderKey]Reader

// reportPlaceholder is the argument that is replaced by the path of a fresh
// temporary report file.
const reportPlaceholder = "{report}"

// findingsExitCodes lists exit codes that mean "ran fine, findings present" for
// a provider; any other non-zero exit is a failure.
var findingsExitCodes = map[ProviderKey][]int{
	ProviderOSVScannerJSON: {1},
}

// CollectRequest is everything Collect needs. Runner and Git default to the
// real process runners; Clock defaults to time.Now.
type CollectRequest struct {
	Root    string
	Origin  Origin
	Config  Config
	Readers Readers
	Runner  ToolRunner
	Git     CommandRunner
	Clock   func() time.Time
}

// Collected is one instance's result and whether the project marked it
// required. Collect itself blocks nothing; policy belongs to the caller.
type Collected struct {
	Result   CapabilityResult
	Required bool
}

// Collection is the outcome of one Collect call: a result per configured
// instance in configuration order, then a not_configured result for each
// capability with no instance, and the saved snapshot.
type Collection struct {
	SnapshotID string
	Created    bool
	Results    []Collected
}

// Collect runs the configured instances one at a time, assembles one result for
// each, classifies them, and saves a single immutable snapshot. Executed tools
// run without a shell; tests and coverage only read reports that already exist.
// One instance's failure never changes another's result, and nothing is pruned.
// An error means no snapshot was saved.
func Collect(ctx context.Context, req CollectRequest) (Collection, error) {
	retention, ok := map[Origin]Retention{OriginOfficial: RetentionPermanent, OriginManual: RetentionPrunable}[req.Origin]
	if !ok {
		return Collection{}, fieldError(ErrInvalidOrigin, "origin", "%q", req.Origin)
	}
	if err := req.Config.Validate(); err != nil {
		return Collection{}, err
	}
	c := &collector{req: req, store: NewStore(req.Root)}
	if c.req.Runner == nil {
		c.req.Runner = ExecRunner{}
	}
	if c.req.Git == nil {
		c.req.Git = GitRunner{}
	}
	if c.req.Clock == nil {
		c.req.Clock = time.Now
	}

	history, err := c.history()
	if err != nil {
		return Collection{}, err
	}
	obs, err := ObserveRepository(ctx, c.req.Git, req.Root, InputScope{})
	if err != nil {
		return Collection{}, err
	}

	var collected []Collected
	configured := map[Capability]bool{}
	for _, cc := range req.Config.Capabilities {
		configured[cc.Capability] = true
		collected = append(collected, Collected{Result: c.instance(ctx, cc), Required: cc.Required})
	}
	for _, capability := range Capabilities() {
		if !configured[capability] {
			collected = append(collected, Collected{Result: CapabilityResult{
				Capability:  capability,
				Outcome:     OutcomeNotConfigured,
				Freshness:   FreshnessUnknown,
				CollectedAt: c.now(),
			}})
		}
	}

	thresholds := make(map[instanceKey]*Threshold)
	for _, cc := range req.Config.Capabilities {
		thresholds[instanceKey{cc.Capability, cc.Provider, cc.Name}] = cc.Thresholds
	}
	snap := Snapshot{
		Version:    SnapshotVersion,
		Origin:     req.Origin,
		Retention:  retention,
		CreatedAt:  c.now(),
		Repository: obs.Identity,
	}
	var assessments []Assessment
	for _, col := range collected {
		r := col.Result
		a := Assess(req.Origin, r, thresholds[r.key()], history[r.key()])
		assessments = append(assessments, a)
		snap.Results = append(snap.Results, r)
		snap.Summary.Capabilities = append(snap.Summary.Capabilities, a.Summary())
	}
	snap.Summary.Overall = Overall(assessments)
	snap.ID = snap.ComputeID()

	created, err := c.store.SaveSnapshot(snap)
	if err != nil {
		return Collection{}, err
	}
	return Collection{SnapshotID: snap.ID, Created: created, Results: collected}, nil
}

type collector struct {
	req   CollectRequest
	store Store
}

func (c *collector) now() string { return c.req.Clock().UTC().Format(timestampLayout) }

// history groups earlier results by instance, oldest first. A damaged history
// fails collection before any tool runs: the store never repairs itself, and
// classifying without it would quietly lose baselines.
func (c *collector) history() (map[instanceKey][]HistoryEntry, error) {
	snaps, err := c.store.LoadSnapshots()
	if err != nil {
		return nil, err
	}
	out := make(map[instanceKey][]HistoryEntry)
	for _, s := range snaps {
		for _, r := range s.Results {
			out[r.key()] = append(out[r.key()], HistoryEntry{Origin: s.Origin, Result: r})
		}
	}
	return out, nil
}

// failure is an outcome without a measurement, with the reason to record.
type failure struct {
	outcome Outcome
	reason  string
}

func fail(o Outcome, format string, args ...any) *failure {
	return &failure{o, boundText(fmt.Sprintf(format, args...))}
}

// instance produces the result for one configured instance. It never returns
// an error: every problem becomes the instance's own outcome.
func (c *collector) instance(ctx context.Context, cc CapabilityConfig) CapabilityResult {
	r := c.measure(ctx, cc)
	r.CollectedAt = c.now()
	return r
}

func (c *collector) measure(ctx context.Context, cc CapabilityConfig) CapabilityResult {
	r := CapabilityResult{
		Capability:   cc.Capability,
		Name:         cc.Name,
		Outcome:      OutcomeFailed,
		Freshness:    FreshnessUnknown,
		Provenance:   Provenance{Provider: cc.Provider},
		ConfigDigest: cc.Digest(),
		Scope:        cc.Scope,
		Exclusions:   cc.Exclusions,
	}

	apply := func(f *failure) CapabilityResult {
		r.Outcome, r.Reason = f.outcome, f.reason
		return r
	}
	if ctx.Err() != nil {
		return apply(fail(OutcomeCancelled, "collection was cancelled before this instance ran"))
	}
	reader, ok := c.req.Readers[cc.Provider]
	if !ok || reader == nil {
		return apply(fail(OutcomeUnsupported, "no reader is registered for %s", cc.Provider))
	}

	var (
		data      []byte
		truncated bool
		freshness = FreshnessFresh
		f         *failure
	)
	if _, executed := DefaultTimeoutSeconds(cc.Provider); executed {
		data, truncated, f = c.execute(ctx, cc)
	} else {
		data, truncated, freshness, f = c.readReportFile(ctx, cc)
	}
	if f != nil {
		return apply(f)
	}

	if truncated {
		r.Outcome, r.Freshness = OutcomePartial, freshness
		r.Provenance.ProviderVersion, r.Provenance.MeasurementDefinition = "unknown", "unknown"
		r.Reason = fmt.Sprintf("report exceeded %d MiB and was not read", MaxReportBytes>>20)
		return r
	}
	reading, err := reader.Read(ctx, ReportInput{Provider: cc.Provider, Data: data, Scope: cc.Scope, Exclusions: cc.Exclusions})
	switch {
	case ctx.Err() != nil:
		return apply(fail(OutcomeCancelled, "collection was cancelled while reading the report"))
	case err != nil:
		return apply(fail(OutcomeFailed, "%s report could not be read: %s", cc.Provider, sanitizeLine(err.Error())))
	}

	measured := r
	measured.CollectedAt = c.now()
	measured.Outcome, measured.Freshness = OutcomeAvailable, freshness
	if reading.Partial {
		measured.Outcome = OutcomePartial
	}
	measured.Reason, measured.Value, measured.Details, measured.Evidence = reading.Reason, reading.Value, reading.Details, reading.Evidence
	measured.Provenance = reading.Provenance
	measured.Provenance.Provider = cc.Provider
	if err := measured.validate("result"); err != nil {
		return apply(fail(OutcomeFailed, "%s reader returned an unusable result: %s", cc.Provider, sanitizeLine(err.Error())))
	}
	return measured
}

// execute runs the configured tool and returns its report bytes.
func (c *collector) execute(ctx context.Context, cc CapabilityConfig) (data []byte, truncated bool, f *failure) {
	if cc.Executable == "" {
		return nil, false, fail(OutcomeUnavailable, "no executable is configured")
	}
	timeout, _ := cc.EffectiveTimeout()
	ictx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	args, reportFile, cleanup, err := c.withReportFile(cc.Args)
	if err != nil {
		return nil, false, fail(OutcomeFailed, "cannot prepare a temporary report file: %s", sanitizeLine(err.Error()))
	}
	defer cleanup()

	res, err := c.req.Runner.Run(ictx, ToolSpec{Dir: c.req.Root, Executable: cc.Executable, Args: args})
	switch {
	case ctx.Err() != nil:
		return nil, false, fail(OutcomeCancelled, "collection was cancelled while %s was running", filepath.Base(cc.Executable))
	case errors.Is(err, context.DeadlineExceeded), ictx.Err() != nil:
		return nil, false, fail(OutcomeTimedOut, "%s did not finish within %s", filepath.Base(cc.Executable), timeout)
	case errors.Is(err, ErrToolUnavailable):
		return nil, false, fail(OutcomeUnavailable, "%s was not found or cannot be run", filepath.Base(cc.Executable))
	case err != nil:
		return nil, false, fail(OutcomeFailed, "%s could not be run: %s", filepath.Base(cc.Executable), sanitizeLine(err.Error()))
	}
	if res.ExitCode != 0 && !slices.Contains(findingsExitCodes[cc.Provider], res.ExitCode) {
		reason := fmt.Sprintf("%s exited with status %d", filepath.Base(cc.Executable), res.ExitCode)
		if res.Stderr != "" {
			reason += ": " + res.Stderr
		}
		return nil, false, fail(OutcomeFailed, "%s", reason)
	}
	if reportFile == "" {
		return res.Stdout, res.Truncated, nil
	}
	data, truncated, _, rf := readBounded(reportFile)
	if rf != nil && rf.outcome == OutcomeAbsent {
		return nil, false, fail(OutcomeFailed, "%s exited cleanly but wrote no report", filepath.Base(cc.Executable))
	}
	return data, truncated, rf
}

// withReportFile replaces a literal {report} argument with a path in a fresh
// temporary directory outside the project. cleanup removes the directory.
func (c *collector) withReportFile(in []string) (args []string, file string, cleanup func(), err error) {
	cleanup = func() {}
	if !slices.Contains(in, reportPlaceholder) {
		return in, "", cleanup, nil
	}
	dir, err := os.MkdirTemp("", "savepoint-health-")
	if err != nil {
		return nil, "", cleanup, err
	}
	cleanup = func() { os.RemoveAll(dir) }
	if inside(c.req.Root, dir) {
		cleanup()
		return nil, "", func() {}, errors.New("the temporary directory is inside the project")
	}
	file = filepath.Join(dir, "report")
	args = slices.Clone(in)
	for i, a := range args {
		if a == reportPlaceholder {
			args[i] = file
		}
	}
	return args, file, cleanup, nil
}

// readReportFile reads an existing report for a report-only provider. It never
// runs anything. The report is stale when older than the newest relevant input.
func (c *collector) readReportFile(ctx context.Context, cc CapabilityConfig) (data []byte, truncated bool, freshness Freshness, f *failure) {
	if cc.Report == "" {
		return nil, false, FreshnessUnknown, fail(OutcomeUnavailable, "no report path is configured")
	}
	full := filepath.Join(c.req.Root, filepath.FromSlash(cc.Report))
	resolved, err := filepath.EvalSymlinks(full)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, false, FreshnessUnknown, fail(OutcomeAbsent, "no report at %s", cc.Report)
	case err != nil:
		return nil, false, FreshnessUnknown, fail(OutcomeFailed, "report at %s cannot be read", cc.Report)
	}
	root, err := filepath.EvalSymlinks(c.req.Root)
	if err != nil || !inside(root, resolved) {
		return nil, false, FreshnessUnknown, fail(OutcomeFailed, "report at %s resolves outside the project", cc.Report)
	}
	data, truncated, mtime, f := readBounded(resolved)
	if f != nil {
		if f.outcome == OutcomeAbsent {
			f = fail(OutcomeAbsent, "no report at %s", cc.Report)
		}
		return nil, false, FreshnessUnknown, f
	}

	scope := InputScope{Include: cc.Scope, Exclude: append(slices.Clone(cc.Exclusions), cc.Report)}
	obs, err := ObserveRepository(ctx, c.req.Git, c.req.Root, scope)
	switch {
	case ctx.Err() != nil:
		return nil, false, FreshnessUnknown, fail(OutcomeCancelled, "collection was cancelled while checking report freshness")
	case err != nil:
		return data, truncated, FreshnessUnknown, nil
	case mtime.Before(obs.Newest):
		return data, truncated, FreshnessStale, nil
	}
	return data, truncated, FreshnessFresh, nil
}

// readBounded reads a regular file up to MaxReportBytes. For a missing file it
// returns an absent failure. The returned time is the file's modification time.
func readBounded(path string) (data []byte, truncated bool, mtime time.Time, f *failure) {
	info, err := os.Stat(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, false, time.Time{}, &failure{outcome: OutcomeAbsent}
	case err != nil:
		return nil, false, time.Time{}, fail(OutcomeFailed, "report cannot be read")
	case !info.Mode().IsRegular():
		return nil, false, time.Time{}, fail(OutcomeFailed, "report is not a regular file")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, false, time.Time{}, fail(OutcomeFailed, "report cannot be opened")
	}
	defer file.Close()
	data, err = io.ReadAll(io.LimitReader(file, MaxReportBytes+1))
	if err != nil {
		return nil, false, time.Time{}, fail(OutcomeFailed, "report cannot be read")
	}
	if len(data) > MaxReportBytes {
		return data[:MaxReportBytes], true, info.ModTime(), nil
	}
	return data, false, info.ModTime(), nil
}

// inside reports whether p is at or below dir.
func inside(dir, p string) bool {
	rel, err := filepath.Rel(dir, p)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}
