package data

import "strings"

type DependencyKind string

const (
	DependencyMissing DependencyKind = ""
	DependencyTask    DependencyKind = "task"
	DependencyEpic    DependencyKind = "epic"
)

type DependencyResolution struct {
	Kind       DependencyKind
	ID         string
	TaskStatus ColumnType
	EpicStatus string
}

func ResolveDependency(ref string, dependent Task, tasks []Task, epicStatuses map[string]string) DependencyResolution {
	for _, task := range tasks {
		if task.ID == ref && sameRelease(dependent.Release, task.Release) {
			return DependencyResolution{Kind: DependencyTask, ID: task.ID, TaskStatus: task.Column}
		}
	}

	if isShortTaskRef(ref) {
		shortRef := taskShortID(ref)
		for _, task := range tasks {
			if taskShortID(task.ID) == shortRef && sameRelease(dependent.Release, task.Release) && sameEpic(dependent.Epic, task.Epic) {
				return DependencyResolution{Kind: DependencyTask, ID: task.ID, TaskStatus: task.Column}
			}
		}
	}

	if status, ok := epicStatuses[ref]; ok {
		return DependencyResolution{Kind: DependencyEpic, ID: ref, EpicStatus: status}
	}

	if isShortEpicRef(ref) {
		for epicID, status := range epicStatuses {
			if epicShortID(epicID) == ref {
				return DependencyResolution{Kind: DependencyEpic, ID: epicID, EpicStatus: status}
			}
		}
	}

	return DependencyResolution{}
}

func sameRelease(a, b string) bool {
	return a == "" || b == "" || a == b
}

func sameEpic(a, b string) bool {
	return a == "" || b == "" || a == b
}

func isShortTaskRef(ref string) bool {
	if strings.Contains(ref, "/") {
		return false
	}
	short := taskShortID(ref)
	return len(short) == 4 && short[0] == 'T' && allDigits(short[1:])
}

func isShortEpicRef(ref string) bool {
	return len(ref) == 3 && ref[0] == 'E' && allDigits(ref[1:])
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return s != ""
}

func taskShortID(id string) string {
	if slash := strings.LastIndexByte(id, '/'); slash >= 0 {
		id = id[slash+1:]
	}
	if dash := strings.IndexByte(id, '-'); dash >= 0 {
		return id[:dash]
	}
	return id
}

func epicShortID(id string) string {
	if dash := strings.IndexByte(id, '-'); dash >= 0 {
		return id[:dash]
	}
	return id
}
