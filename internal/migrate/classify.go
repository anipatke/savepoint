package migrate

import "regexp"

// Role names what an inventoried source path is, so migration can give it an
// accountable destination or archive reference. The vocabulary is frozen
// against the two V1 fixture manifests in internal/data/testdata/migration/.
type Role string

const (
	RoleConfig        Role = "config"
	RoleRouter        Role = "router"
	RoleProductPRD    Role = "product-prd"
	RoleArchitecture  Role = "architecture"
	RoleHealthCheck   Role = "health-check"
	RoleRelease       Role = "release"
	RoleReleasePRD    Role = "release-prd"
	RoleEpicDetail    Role = "epic-detail"
	RoleEpicAudit     Role = "epic-audit"
	RoleTask          Role = "task"
	RoleDefect        Role = "defect"
	RoleAuditPrompt   Role = "audit-prompt"
	RoleAuditRegister Role = "audit-register"
	RoleFinding       Role = "finding"
	RoleAuditRun      Role = "audit-run"
	RoleManagedGuide  Role = "managed-guide"
	RoleSkill         Role = "skill"

	// RoleUnclassified marks a file that matched no known shape. It is still
	// inventoried and still reported — never dropped, never guessed at, and
	// never silently archived without appearing in the report.
	RoleUnclassified Role = "unclassified"
)

// classifyRule pairs a role with the path shape that assigns it. Keeping the
// vocabulary as data here, rather than as a chain of conditionals, is what
// lets the fixture-driven tests enumerate every rule and lets a new V1 shape
// be added as one more row.
type classifyRule struct {
	role    Role
	pattern *regexp.Regexp
}

// classifyRules matches against a project-relative, forward-slash normalized
// path — the same form Inventory returns in SourceFile.Path. Every pattern is
// anchored start-to-end so a role never matches a path it merely contains.
var classifyRules = []classifyRule{
	{RoleManagedGuide, regexp.MustCompile(`^AGENTS\.md$`)},
	{RoleSkill, regexp.MustCompile(`^agent-skills/[^/]+/SKILL\.md$`)},
	{RoleConfig, regexp.MustCompile(`^\.savepoint/config\.yml$`)},
	{RoleRouter, regexp.MustCompile(`^\.savepoint/router\.md$`)},
	{RoleProductPRD, regexp.MustCompile(`^\.savepoint/PRD\.md$`)},
	{RoleArchitecture, regexp.MustCompile(`^\.savepoint/Design\.md$`)},
	{RoleHealthCheck, regexp.MustCompile(`^\.savepoint/Health-Check\.md$`)},
	{RoleRelease, regexp.MustCompile(`^\.savepoint/releases/[^/]+$`)},
	{RoleAuditPrompt, regexp.MustCompile(`^\.savepoint/audit/prompt\.md$`)},
	{RoleAuditRegister, regexp.MustCompile(`^\.savepoint/audit/register\.md$`)},
	{RoleFinding, regexp.MustCompile(`^\.savepoint/audit/findings/[^/]+\.md$`)},
	{RoleAuditRun, regexp.MustCompile(`^\.savepoint/audit/runs/[^/]+\.md$`)},
	{RoleReleasePRD, regexp.MustCompile(`^\.savepoint/releases/[^/]+/[^/]+-PRD\.md$`)},
	{RoleEpicDetail, regexp.MustCompile(`^\.savepoint/releases/[^/]+/epics/[^/]+/[^/]+-Detail\.md$`)},
	{RoleEpicAudit, regexp.MustCompile(`^\.savepoint/releases/[^/]+/epics/[^/]+/[^/]+-Audit\.md$`)},
	{RoleTask, regexp.MustCompile(`^\.savepoint/releases/[^/]+/epics/[^/]+/tasks/[^/]+\.md$`)},
	{RoleDefect, regexp.MustCompile(`^\.savepoint/releases/[^/]+/defects/[^/]+\.md$`)},
}

// Classify assigns path (project-relative, forward-slash form) the role of
// the first matching rule, or RoleUnclassified when none match.
func Classify(path string) Role {
	for _, rule := range classifyRules {
		if rule.pattern.MatchString(path) {
			return rule.role
		}
	}
	return RoleUnclassified
}
