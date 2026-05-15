package data

import (
	"fmt"
	"time"
)

type CheckItem struct {
	Text string
	Done bool
}

type ColumnType string

const (
	ColumnPlanned    ColumnType = "planned"
	ColumnInProgress ColumnType = "in_progress"
	ColumnDone       ColumnType = "done"
)

type ProgressStage string

const (
	StageBuild ProgressStage = "build"
	StageTest  ProgressStage = "test"
	StageAudit ProgressStage = "audit"
)

type ComplexityTier string

const (
	ComplexityLow    ComplexityTier = "low"
	ComplexityMedium ComplexityTier = "medium"
	ComplexityHigh   ComplexityTier = "high"
	ComplexitySpike  ComplexityTier = "spike"
)

const MaxComplexityReasonLen = 120

type Progress struct {
	Stage   ProgressStage `yaml:"stage"`
	Started bool          `yaml:"started"`
}

type Task struct {
	ID               string         `yaml:"id"`
	Title            string         `yaml:"title"`
	Description      string         `yaml:"description,omitempty"`
	Epic             string         `yaml:"epic"`
	Release          string         `yaml:"release"`
	Column           ColumnType     `yaml:"column"`
	Stage            ProgressStage  `yaml:"stage,omitempty"`
	Priority         string         `yaml:"priority,omitempty"`
	Points           int            `yaml:"points,omitempty"`
	Tags             []string       `yaml:"tags,omitempty"`
	Acceptance       []string       `yaml:"acceptance,omitempty"`
	Checklist        []CheckItem    `yaml:"checklist,omitempty"`
	Notes            string         `yaml:"notes,omitempty"`
	DependsOn        []string       `yaml:"depends_on,omitempty"`
	Progress         Progress       `yaml:"progress,omitempty"`
	ComplexityTier   ComplexityTier `yaml:"complexity_tier,omitempty"`
	ComplexityReason string         `yaml:"complexity_reason,omitempty"`
	Path             string         `yaml:"-"`
	Mtime            time.Time      `yaml:"-"`
}

func (t Task) String() string {
	return fmt.Sprintf("Task(%s)", t.ID)
}
