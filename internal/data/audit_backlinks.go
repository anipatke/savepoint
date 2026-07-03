package data

// FindingsForTask returns the findings whose Tasks link back to the given full
// task ID, preserving the input order (findings arrive pre-sorted from the
// loader). A link matches when it equals the task ID exactly or is the task's
// short T### shorthand, mirroring how ResolveDependency and ValidateAuditFindings
// resolve task references. An empty ID or finding set yields no results, so an
// absent audit register never blocks a caller.
func FindingsForTask(findings []AuditFinding, taskID string) []AuditFinding {
	if taskID == "" {
		return nil
	}
	short := taskShortID(taskID)
	var linked []AuditFinding
	for _, f := range findings {
		for _, ref := range f.Tasks {
			if ref == taskID || (isShortTaskRef(ref) && taskShortID(ref) == short) {
				linked = append(linked, f)
				break
			}
		}
	}
	return linked
}

// FindingsForEpic returns the findings whose Epics link back to the given epic
// ID, preserving the input order. A link matches when it equals the epic ID
// exactly or is the epic's short E## shorthand, mirroring epic dependency
// resolution. An empty ID or finding set yields no results.
func FindingsForEpic(findings []AuditFinding, epicID string) []AuditFinding {
	if epicID == "" {
		return nil
	}
	short := epicShortID(epicID)
	var linked []AuditFinding
	for _, f := range findings {
		for _, ref := range f.Epics {
			if ref == epicID || (isShortEpicRef(ref) && ref == short) {
				linked = append(linked, f)
				break
			}
		}
	}
	return linked
}
