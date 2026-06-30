package data

import "fmt"

// AuditValidationCode classifies a cross-record audit-register problem so doctor
// can react to a class without string matching the message. Unlike the load-time
// FindingDiagnostic codes, these are not healable: they describe a finding whose
// lifecycle field or work-item link cannot be reconciled against the rest of the
// register and the discovered project, so they are reported and left for the
// author to fix.
type AuditValidationCode string

const (
	// Lifecycle-required fields: a status the finding has reached implies a
	// field that is missing.
	AuditVerifiedMissingProof          AuditValidationCode = "verified_missing_proof"
	AuditDuplicateMissingTarget        AuditValidationCode = "duplicate_missing_target"
	AuditDeferredMissingRationale      AuditValidationCode = "deferred_missing_rationale"
	AuditOwnerDecisionMissingRationale AuditValidationCode = "owner_decision_missing_rationale"
	AuditWaivedMissingRationale        AuditValidationCode = "waived_missing_rationale"

	// Broken references: a release, epic, task, defect, or duplicate-of link
	// that does not resolve against the discovered project or register.
	AuditReleaseRefMissing   AuditValidationCode = "release_ref_missing"
	AuditEpicRefMissing      AuditValidationCode = "epic_ref_missing"
	AuditTaskRefMissing      AuditValidationCode = "task_ref_missing"
	AuditDefectRefMissing    AuditValidationCode = "defect_ref_missing"
	AuditDuplicateRefMissing AuditValidationCode = "duplicate_ref_missing"
)

// AuditFindingValidation is one cross-record problem found in a finding. It names
// the finding so doctor can group results by file, carries a typed Code for
// programmatic handling, and an actionable Message for display. A clean register
// produces an empty slice.
type AuditFindingValidation struct {
	FindingID string
	Code      AuditValidationCode
	Message   string
}

// AuditWorkItems is the set of discovered work-item IDs a finding may reference.
// Doctor builds it by flattening discovery (releases, epics, tasks, defects) so
// the validator stays decoupled from filesystem traversal and is trivial to test
// with literal IDs. Task and epic references also resolve by their short ID
// (T###, E##) via the shared dependency helpers; releases and defects match by
// exact ID only.
type AuditWorkItems struct {
	Releases []string
	Epics    []string
	Tasks    []string
	Defects  []string
}

// ValidateAuditFindings checks every finding against the lifecycle-required-field
// rules and resolves its work-item and duplicate links against items and the
// finding set itself. Results are returned in finding order, lifecycle problems
// before reference problems, and references in release → epic → task → defect →
// duplicate order; references within a list keep their source order. The input is
// never mutated. A register whose findings are internally consistent and whose
// links all resolve returns nil.
func ValidateAuditFindings(findings []AuditFinding, items AuditWorkItems) []AuditFindingValidation {
	findingIDs := idSet(findingIDsOf(findings))
	releaseIDs := idSet(items.Releases)
	defectIDs := idSet(items.Defects)
	taskIDs := idSet(items.Tasks)
	epicIDs := idSet(items.Epics)
	taskShortIDs := shortIDSet(items.Tasks, taskShortID)
	epicShortIDs := shortIDSet(items.Epics, epicShortID)

	var results []AuditFindingValidation
	add := func(f AuditFinding, code AuditValidationCode, message string) {
		results = append(results, AuditFindingValidation{FindingID: f.ID, Code: code, Message: message})
	}

	for _, f := range findings {
		// Lifecycle-required fields.
		switch f.Status {
		case FindingVerified:
			if f.VerifiedProof == "" {
				add(f, AuditVerifiedMissingProof,
					"verified finding has no named proof; set verified_proof or record the passing regression test in ## Proof")
			}
		case FindingDuplicate:
			if f.DuplicateOf == "" {
				add(f, AuditDuplicateMissingTarget,
					"duplicate finding has no duplicate_of; point it at the canonical F### finding")
			}
		case FindingDeferred:
			if !hasRationale(f) {
				add(f, AuditDeferredMissingRationale,
					"deferred finding has no rationale; set deferral_reason")
			}
		case FindingOwnerDecision:
			if !hasRationale(f) {
				add(f, AuditOwnerDecisionMissingRationale,
					"owner_decision finding has no rationale; set deferral_reason or waiver_reason")
			}
		case FindingWaived:
			if !hasRationale(f) {
				add(f, AuditWaivedMissingRationale,
					"waived finding has no rationale; set waiver_reason or deferral_reason")
			}
		}

		// Broken references.
		for _, ref := range f.Releases {
			if !releaseIDs[ref] {
				add(f, AuditReleaseRefMissing, fmt.Sprintf("references unknown release %q", ref))
			}
		}
		for _, ref := range f.Epics {
			if !resolvesEpicRef(ref, epicIDs, epicShortIDs) {
				add(f, AuditEpicRefMissing, fmt.Sprintf("references unknown epic %q", ref))
			}
		}
		for _, ref := range f.Tasks {
			if !resolvesTaskRef(ref, taskIDs, taskShortIDs) {
				add(f, AuditTaskRefMissing, fmt.Sprintf("references unknown task %q", ref))
			}
		}
		for _, ref := range f.Defects {
			if !defectIDs[ref] {
				add(f, AuditDefectRefMissing, fmt.Sprintf("references unknown defect %q", ref))
			}
		}
		if f.DuplicateOf != "" {
			switch {
			case f.DuplicateOf == f.ID:
				add(f, AuditDuplicateRefMissing,
					fmt.Sprintf("duplicate_of %q points at the finding itself; point it at a different canonical F###", f.DuplicateOf))
			case !findingIDs[f.DuplicateOf]:
				add(f, AuditDuplicateRefMissing, fmt.Sprintf("duplicate_of references unknown finding %q", f.DuplicateOf))
			}
		}
	}

	return results
}

// hasRationale reports whether a finding carries any human-readable rationale for
// a resolved-without-fix disposition. Either deferral_reason or waiver_reason
// satisfies it: the finding template documents deferral_reason as the rationale
// for deferred, owner_decision, and waived, while the model adds a dedicated
// waiver_reason, so accepting either avoids falsely flagging a finding that
// followed either convention.
func hasRationale(f AuditFinding) bool {
	return f.DeferralReason != "" || f.WaiverReason != ""
}

// findingIDsOf collects the IDs of every finding for duplicate-of resolution.
func findingIDsOf(findings []AuditFinding) []string {
	ids := make([]string, len(findings))
	for i, f := range findings {
		ids[i] = f.ID
	}
	return ids
}

// idSet builds a membership set from a slice of IDs, dropping the empty string so
// an unset reference never resolves against an unset ID.
func idSet(ids []string) map[string]bool {
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		if id != "" {
			set[id] = true
		}
	}
	return set
}

// shortIDSet indexes IDs by their short form (e.g. T### or E##) so a finding can
// reference a task or epic by short ID, mirroring dependency resolution. Empty
// short forms are dropped.
func shortIDSet(ids []string, short func(string) string) map[string]bool {
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		if s := short(id); s != "" {
			set[s] = true
		}
	}
	return set
}

// resolvesEpicRef reports whether an epic reference resolves against the
// discovered project. A full ID must match exactly; the short-ID fallback only
// applies to genuine shorthand (E##) references, mirroring ResolveDependency, so
// a full-looking but wrong reference like "E31-wrong-slug" is not silently
// accepted because some E31-* epic exists.
func resolvesEpicRef(ref string, fullIDs, shortIDs map[string]bool) bool {
	if fullIDs[ref] {
		return true
	}
	return isShortEpicRef(ref) && shortIDs[ref]
}

// resolvesTaskRef reports whether a task reference resolves against the
// discovered project. A full ID must match exactly; the short-ID fallback only
// applies to genuine shorthand (T###) references, mirroring ResolveDependency, so
// a broken full path like "E99-missing/T004-..." is not silently accepted because
// some T004 exists elsewhere.
func resolvesTaskRef(ref string, fullIDs, shortIDs map[string]bool) bool {
	if fullIDs[ref] {
		return true
	}
	if !isShortTaskRef(ref) {
		return false
	}
	return shortIDs[taskShortID(ref)]
}
