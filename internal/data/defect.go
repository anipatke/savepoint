package data

import "time"

type DefectSeverity string

const (
	SeverityCritical DefectSeverity = "critical"
	SeverityHigh     DefectSeverity = "high"
	SeverityMedium   DefectSeverity = "medium"
	SeverityLow      DefectSeverity = "low"
)

type DefectStatus string

const (
	DefectOpen       DefectStatus = "open"
	DefectInProgress DefectStatus = "in_progress"
	DefectResolved   DefectStatus = "resolved"
)

type Defect struct {
	ID         string         `yaml:"id"`
	Release    string         `yaml:"release"`
	Status     DefectStatus   `yaml:"status"`
	Severity   DefectSeverity `yaml:"severity"`
	Introduced string         `yaml:"introduced,omitempty"`
	Reference  string         `yaml:"reference,omitempty"`
	Stage      ProgressStage  `yaml:"stage,omitempty"`
	Title      string         `yaml:"title"`
	Body       string         `yaml:"-"`
	Path       string         `yaml:"-"`
	Mtime      time.Time      `yaml:"-"`
}

type DefectInfo struct {
	ID   string
	Path string
}
