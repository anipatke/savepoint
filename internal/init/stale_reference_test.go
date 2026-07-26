package init

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// genericAuditName is the retired skill that split into savepoint-audit-task and
// savepoint-audit-epic. It may survive only where a record of the past is the
// point: released planning history, and the migration archive an upgrade writes.
const genericAuditName = "savepoint-audit"

// auditSkillToken matches the retired name plus any suffix, so a live source can
// keep the successor names while the bare alias is rejected.
var auditSkillToken = regexp.MustCompile(`savepoint-audit[a-z-]*`)

// liveSuccessorNames are the audit skill names live guidance may still use.
var liveSuccessorNames = map[string]bool{
	"savepoint-audit-task":     true,
	"savepoint-audit-epic":     true,
	"savepoint-audit-register": true,
}

// staleAuditReferences returns each line of content that names the retired
// generic skill rather than one of its successors.
func staleAuditReferences(content string) []string {
	var stale []string
	for i, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		for _, match := range auditSkillToken.FindAllString(line, -1) {
			// A trailing hyphen means the name was split across a line wrap;
			// treat it as a successor prefix rather than the bare alias.
			if match == genericAuditName+"-" || liveSuccessorNames[match] {
				continue
			}
			if match == genericAuditName {
				stale = append(stale, strings.TrimSpace(line)+" (line "+itoa(i+1)+")")
			}
		}
	}
	return stale
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

// liveGuidanceFiles walks the markdown sources an agent actually loads for
// routing and skill selection.
//
// Two areas are deliberately out of scope. This repository's own
// .savepoint/releases/ holds task files, audit records, and release PRDs that
// are historical records of when the generic skill existed; rewriting them
// would falsify the history. README.md is user documentation that must name the
// retired skill to explain what upgrade does to a project that still has it.
func liveGuidanceFiles(t *testing.T) []string {
	t.Helper()

	root := filepath.Join("..", "..")
	roots := []string{
		filepath.Join(root, "AGENTS.md"),
		filepath.Join(root, ".savepoint", "router.md"),
		filepath.Join(root, ".savepoint", "Design.md"),
		filepath.Join(root, "agent-skills"),
		filepath.Join(root, "templates"),
	}

	var files []string
	for _, target := range roots {
		info, err := os.Stat(target)
		if err != nil {
			t.Fatalf("stat live source %s: %v", target, err)
		}
		if !info.IsDir() {
			files = append(files, target)
			continue
		}

		err = filepath.WalkDir(target, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			t.Fatalf("walk live source %s: %v", target, err)
		}
	}

	if len(files) == 0 {
		t.Fatal("no live guidance files found")
	}
	return files
}

func TestLiveSourcesDoNotReferenceGenericAuditSkill(t *testing.T) {
	for _, path := range liveGuidanceFiles(t) {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		for _, line := range staleAuditReferences(string(data)) {
			t.Errorf("%s references the retired %q skill: %s", path, genericAuditName, line)
		}
	}
}

func TestScaffoldedProjectDoesNotReferenceGenericAuditSkill(t *testing.T) {
	target := scaffoldFromRealTemplates(t)

	err := filepath.WalkDir(target, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		rel, _ := filepath.Rel(target, path)
		for _, line := range staleAuditReferences(string(data)) {
			t.Errorf("generated %s references the retired %q skill: %s", rel, genericAuditName, line)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk generated project: %v", err)
	}
}

// The exclusion itself is a contract: historical release records are allowed to
// name the retired skill, and the check must not be tightened into rewriting
// them.
func TestStaleReferenceCheckExcludesHistoricalReleaseRecords(t *testing.T) {
	// This repository's own release records are history, not guidance. The
	// scaffold stub under templates/ is not history and stays in scope.
	historical := filepath.ToSlash(filepath.Join("..", "..", ".savepoint", "releases")) + "/"
	for _, path := range liveGuidanceFiles(t) {
		if strings.HasPrefix(filepath.ToSlash(path), historical) {
			t.Errorf("stale-reference scan must exclude historical release record %s", path)
		}
	}

	// Guard the exclusion itself: those records exist and do name the retired
	// skill, so a scan that started walking them would fail loudly.
	records, err := filepath.Glob(filepath.Join("..", "..", ".savepoint", "releases", "*", "epics", "*", "tasks", "*.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(records) == 0 {
		t.Fatal("no historical release records found; the exclusion is untested")
	}

	// The scan still rejects the bare alias, and still accepts the successors.
	if got := staleAuditReferences("use `savepoint-audit` for review"); len(got) != 1 {
		t.Errorf("bare alias not rejected: %v", got)
	}
	for _, ok := range []string{
		"| audit-pending | savepoint-audit-epic |",
		"use `savepoint-audit-task` for one task",
		"follow the `savepoint-audit-register` skill",
	} {
		if got := staleAuditReferences(ok); len(got) != 0 {
			t.Errorf("successor name %q rejected: %v", ok, got)
		}
	}
}
