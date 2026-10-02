package codehealth

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Report file layout, below the health directory.
const (
	reportFile    = "report.md"
	gitignoreFile = ".gitignore"
)

// ReportPath is where the report lives, relative to the project root, as the
// popover names it.
const ReportPath = ".savepoint/health/report.md"

// reportExists is one file check: true when the report is a regular file.
func (s Store) reportExists() bool {
	dir, ok, err := s.dir(false, healthDir)
	if err != nil || !ok {
		return false
	}
	info, err := os.Lstat(filepath.Join(dir, reportFile))
	return err == nil && info.Mode().IsRegular()
}

// Refusals RenderReport gives instead of a report. Callers match them with
// errors.Is.
var (
	ErrReportNotConfigured = errors.New("Code Health is not set up in this project, so there is no report to write")
	ErrReportNoSnapshot    = errors.New("Code Health has no saved snapshot yet, so there is no report to write")
)

// Every phrase the report says lives here (STYLE-09). The signals' own
// wording comes from the dashboard rows, never from this table.
const (
	reportTitle    = "# Code Health report"
	reportBrief    = "Investigate each signal below that is not Good, blocking ones first. For each, find the cause in the files listed, propose a fix, then apply it. When you are done, re-run the check with `%s` and confirm the signal improved."
	reportRerun    = "savepoint health check %s"
	reportMeasured = "Measured %s (%s)."
	reportNoFiles  = "No affected files were recorded."
	reportQuestion = "Question: %s"
	reportLabel    = "Label: %s"
	reportNumber   = "Number: %s"
	reportAim      = "Aim: %s"
	reportMeaning  = "Meaning: %s"
	reportNext     = "Next step: %s"
	reportSignOff  = "Sign-off: %s"
	reportTrend    = "Trend: %s"
	reportFiles    = "Affected files:"
)

// reportRank orders labels for the report: what needs a look comes first.
var reportRank = map[Classification]int{
	ClassificationNeedsAttention: 0,
	ClassificationWatch:          1,
	ClassificationUnknown:        2,
	ClassificationGood:           3,
}

// RenderReport describes a loaded dashboard as one plain-text report an agent
// can act on. objective names the Objective whose check the brief tells the
// agent to re-run. It reads no file and runs nothing.
func RenderReport(d Dashboard, objective string) (string, error) {
	switch d.State {
	case DashboardNotConfigured:
		return "", ErrReportNotConfigured
	case DashboardMeasured:
	default:
		return "", ErrReportNoSnapshot
	}
	rows := slices.Clone(d.Rows)
	slices.SortStableFunc(rows, func(a, b DashboardRow) int { return reportRank[a.Label] - reportRank[b.Label] })

	var b strings.Builder
	line := func(format string, args ...any) { fmt.Fprintf(&b, format+"\n", args...) }
	line("%s", reportTitle)
	line("")
	line(reportBrief, fmt.Sprintf(reportRerun, objective))
	line("")
	line(reportMeasured, d.MeasuredText, strings.ToLower(d.OriginText))
	line("%s", d.SignOff)
	line("%s", d.Headline)
	for _, row := range rows {
		line("")
		writeReportRow(line, row)
	}
	return b.String(), nil
}

func writeReportRow(line func(string, ...any), row DashboardRow) {
	heading := row.CapabilityText
	if row.Name != "" {
		heading += " (" + row.Name + ")"
	}
	line("## %s", heading)
	line("")
	line(reportQuestion, row.Question)
	line(reportLabel, row.LabelText)
	line(reportNumber, row.Value)
	line(reportAim, row.Aim)
	line(reportMeaning, row.Meaning)
	line(reportNext, row.NextStep)
	if row.SignOff != "" {
		line(reportSignOff, row.SignOff)
	}
	if row.SparkWord != "" {
		line(reportTrend, row.SparkWord)
	}
	line("")
	if len(row.Evidence) == 0 {
		line("%s", reportNoFiles)
		return
	}
	line("%s", reportFiles)
	for _, e := range row.Evidence {
		ref := e.Path
		if e.Line > 0 {
			ref = fmt.Sprintf("%s:%d", e.Path, e.Line)
		}
		if e.Note != "" {
			ref += " " + e.Note
		}
		line("- %s", ref)
	}
}

// WriteReport stores the report at .savepoint/health/report.md, replacing an
// older one atomically so a reader never sees half a file. It reports whether
// the file changed; identical text is left untouched. The first write also
// adds a .gitignore beside it that ignores the report, since the report is
// derived; an existing .gitignore there is never touched.
func (s Store) WriteReport(text string) (changed bool, err error) {
	dir, _, err := s.dir(true, healthDir)
	if err != nil {
		return false, err
	}
	if err := ensureReportIgnore(dir); err != nil {
		return false, err
	}
	path := filepath.Join(dir, reportFile)
	old, exists, err := readRecord(path)
	if err != nil {
		return false, err
	}
	if exists && string(old) == text {
		return false, nil
	}
	tmp, err := writeTemp(dir, []byte(text))
	if err != nil {
		return false, err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return false, err
	}
	return true, nil
}

// ensureReportIgnore creates the ignore file when none exists. It links rather
// than renames, so a file the owner edited, or one that appears meanwhile, is
// never replaced.
func ensureReportIgnore(dir string) error {
	path := filepath.Join(dir, gitignoreFile)
	if _, err := os.Lstat(path); err == nil {
		return nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	tmp, err := writeTemp(dir, []byte(reportFile+"\n"))
	if err != nil {
		return err
	}
	defer os.Remove(tmp)
	if err := os.Link(tmp, path); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	return nil
}
