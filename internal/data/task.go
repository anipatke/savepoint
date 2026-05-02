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

type TaskStatus string

const (
	StatusPlanned    TaskStatus = "planned"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
	StatusAudited    TaskStatus = "audited"
)

type Progress struct {
	Stage   ProgressStage `yaml:"stage"`
	Started bool          `yaml:"started"`
}

type Task struct {
	ID          string        `yaml:"id"`
	Title       string        `yaml:"title"`
	Description string        `yaml:"description,omitempty"`
	Status      string        `yaml:"status,omitempty"`
	Epic        string        `yaml:"epic"`
	Release     string        `yaml:"release"`
	Column      ColumnType    `yaml:"column"`
	Stage       ProgressStage `yaml:"stage,omitempty"`
	Priority    string        `yaml:"priority,omitempty"`
	Points      int           `yaml:"points,omitempty"`
	Tags        []string      `yaml:"tags,omitempty"`
	Acceptance  []string      `yaml:"acceptance,omitempty"`
	Checklist   []CheckItem   `yaml:"checklist,omitempty"`
	Notes       string        `yaml:"notes,omitempty"`
	DependsOn   []string      `yaml:"depends_on,omitempty"`
	Progress    Progress      `yaml:"progress,omitempty"`
	Path        string        `yaml:"-"`
	Mtime       time.Time     `yaml:"-"`
}

func (t Task) String() string {
	return fmt.Sprintf("Task(%s)", t.ID)
}
