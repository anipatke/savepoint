package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fencedAwareHeadings returns the "## " and "### " headings of a Markdown
// file in order, ignoring lines inside fenced code blocks.
func fencedAwareHeadings(content string) (h2, h3 []string) {
	fenced := false
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(line, "```") {
			fenced = !fenced
			continue
		}
		switch {
		case fenced:
		case strings.HasPrefix(line, "## "):
			h2 = append(h2, strings.TrimPrefix(line, "## "))
		case strings.HasPrefix(line, "### "):
			h3 = append(h3, strings.TrimPrefix(line, "### "))
		}
	}
	return h2, h3
}

// contentsList returns the bullet items under the "**Contents**" label.
func contentsList(content string) []string {
	_, rest, found := strings.Cut(strings.ReplaceAll(content, "\r\n", "\n"), "\n**Contents**\n")
	if !found {
		return nil
	}
	var items []string
	for _, line := range strings.Split(rest, "\n") {
		switch {
		case strings.HasPrefix(line, "- "):
			items = append(items, strings.TrimPrefix(line, "- "))
		case strings.TrimSpace(line) == "":
		default:
			return items
		}
	}
	return items
}

func TestLongReferencesOpenWithContentsListInFileOrder(t *testing.T) {
	for _, name := range []string{"check-method.md", "issue-capture.md"} {
		forEachSkillFile(t, func(root string) string { return filepath.Join(root, "references", name) }, func(tree, path, content string) {
			want, _ := fencedAwareHeadings(content)
			got := contentsList(content)
			if len(got) == 0 {
				t.Fatalf("%s: %s has no contents list", tree, path)
			}
			if strings.Join(got, "\n") != strings.Join(want, "\n") {
				t.Errorf("%s: %s contents list = %q, want every ## heading in file order %q", tree, path, got, want)
			}
			if first := strings.Index(content, "\n## "); strings.Index(content, "\n**Contents**\n") > first {
				t.Errorf("%s: %s contents list comes after the first ## section", tree, path)
			}
		})
	}
}

func TestCheckMethodFullChecklistNamesRealHeadingsInOrder(t *testing.T) {
	forEachSkillFile(t, func(root string) string { return filepath.Join(root, "references", "check-method.md") }, func(tree, path, content string) {
		h2, h3 := fencedAwareHeadings(content)
		index := map[string]int{}
		for _, h := range h2 {
			index[h] = strings.Index(content, "\n## "+h+"\n")
		}
		for _, h := range h3 {
			index[h] = strings.Index(content, "\n### "+h+"\n")
		}
		if !strings.Contains(content, "\n### Full Check Progress Checklist\n") {
			t.Fatalf("%s: %s has no Full Check Progress Checklist", tree, path)
		}
		if _, ok := sectionBody(content, "## Quick And Full Evidence Modes"); !ok {
			t.Fatalf("%s: %s lacks Quick And Full Evidence Modes", tree, path)
		}
		modes := content[strings.Index(content, "\n## Quick And Full Evidence Modes\n"):]
		modes = modes[:strings.Index(modes[1:], "\n## ")+1]
		if strings.Index(modes, "### Full Check Progress Checklist") < 0 {
			t.Errorf("%s: %s checklist is not under Quick And Full Evidence Modes", tree, path)
		}
		if !strings.Contains(modes, "**Full**") || !strings.Contains(modes, "Full Check Progress Checklist, below") {
			t.Errorf("%s: %s Full bullet does not point to the checklist", tree, path)
		}

		_, block, _ := strings.Cut(content, "### Full Check Progress Checklist\n")
		_, block, _ = strings.Cut(block, "```markdown\n")
		block, _, _ = strings.Cut(block, "```")
		steps := 0
		for _, line := range strings.Split(block, "\n") {
			if line == "" {
				continue
			}
			name, ok := strings.CutPrefix(line, "- [ ] ")
			if !ok {
				t.Errorf("%s: %s checklist line %q is not an unticked checkbox", tree, path, line)
				continue
			}
			steps++
			if _, ok := index[name]; !ok {
				t.Errorf("%s: %s checklist step %q names no existing heading", tree, path, name)
			}
		}
		// Every ## section after the modes section, except Re-check wording
		// that is conditional, must appear as a step: Full applies them all.
		for _, h := range h2 {
			if h == "Task Check And Objective Check Depth" || h == "Quick And Full Evidence Modes" {
				continue
			}
			if !strings.Contains(block, "- [ ] "+h+"\n") {
				t.Errorf("%s: %s checklist omits section %q", tree, path, h)
			}
		}
		last := -1
		for _, line := range strings.Split(block, "\n") {
			if name, ok := strings.CutPrefix(line, "- [ ] "); ok {
				if i := index[name]; i < last {
					t.Errorf("%s: %s checklist step %q is out of method order", tree, path, name)
				} else {
					last = i
				}
			}
		}
		if steps == 0 {
			t.Errorf("%s: %s checklist has no steps", tree, path)
		}
	})
}

func TestUpgradeDeliversRevisedReferencesToUneditedProjects(t *testing.T) {
	root := filepath.Join("..", "..")
	templates := os.DirFS(filepath.Join(root, "templates", "project-v2"))
	target := savepointProject(t)

	for _, name := range []string{"check-method.md", "issue-capture.md"} {
		dir := filepath.Join(target, "agent-skills", "references")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		// An unedited project holds the previous release: no contents list.
		current := readTemplate(t, root, "agent-skills", "references", name)
		before, _, _ := strings.Cut(current, "\n**Contents**\n")
		_, after, _ := strings.Cut(current, "\n## ")
		old := before + "\n## " + after
		if err := os.WriteFile(filepath.Join(dir, name), []byte(old), 0644); err != nil {
			t.Fatal(err)
		}
	}

	report := mustUpgrade(t, templates, target)

	for _, name := range []string{"check-method.md", "issue-capture.md"} {
		path := "agent-skills/references/" + name
		if got := upgradeActionFor(t, report, path); got != ActionUpdated {
			t.Errorf("%s action = %v, want updated", path, got)
		}
		got := readTemplate(t, target, "agent-skills", "references", name)
		want := readTemplate(t, root, "templates", "project-v2", "agent-skills", "references", name)
		if got != want {
			t.Errorf("%s was not refreshed to the revised reference", path)
		}
		if len(contentsList(got)) == 0 {
			t.Errorf("%s has no contents list after upgrade", path)
		}
	}
}

const writeResumeFixRule = "After writing or editing any `.savepoint/` record"

func TestWriteResumeFixRuleHasOneHomeAndSkillsPointToIt(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, path := range []string{
		filepath.Join(root, "templates", "project-v2", "AGENTS.md"),
		filepath.Join(root, "AGENTS.md"),
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		_, cli, found := strings.Cut(content, "\n## CLI Rules\n")
		if !found || !strings.Contains(cli, writeResumeFixRule) ||
			!strings.Contains(cli, "run `savepoint resume` again until it loads") ||
			!strings.Contains(cli, "Do not substitute other work") {
			t.Errorf("%s: CLI Rules lacks the write → resume → fix rule", path)
		}
		if strings.Count(content, "until it loads") != 1 {
			t.Errorf("%s: rule should appear exactly once", path)
		}
	}
	for _, skill := range v2Skills {
		forEachSkillFile(t, func(root string) string { return filepath.Join(root, skill, "SKILL.md") }, func(tree, path, content string) {
			if !strings.Contains(content, "follow AGENTS.md's CLI Rules (write → resume → fix)") {
				t.Errorf("%s: %s lacks the pointer to the CLI Rules rule", tree, path)
			}
			for _, restated := range []string{"until it loads", "strict-load the complete index", "strict loading of the full index"} {
				if strings.Contains(content, restated) {
					t.Errorf("%s: %s restates the rule (%q)", tree, path, restated)
				}
			}
		})
	}
}

func TestSkillDescriptionsNameTheirNextWords(t *testing.T) {
	root := filepath.Join("..", "..")
	data, err := os.ReadFile(filepath.Join(root, "templates", "project-v2", "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, workflow, _ := strings.Cut(string(data), "2. The Next line's first word")
	workflow, _, _ = strings.Cut(workflow, "\n")

	routed := map[string][]string{
		"savepoint-design": {"Plan", "Replan"},
		"savepoint-task":   {"Start", "Build", "Test", "Pick a Task in", "Fix"},
		"savepoint-check":  {"Check", "Assess"},
	}
	for skill, words := range routed {
		for _, w := range words {
			if !strings.Contains(workflow, "`"+w+"`") {
				t.Errorf("Workflow no longer routes %q (expected for %s)", w, skill)
			}
		}
	}

	known := []string{"Start", "Build", "Test", "Check", "Assess", "Plan", "Replan", "Pick a Task in", "Fix"}
	for _, skill := range []string{"savepoint-design", "savepoint-task", "savepoint-check"} {
		forEachSkillFile(t, func(root string) string { return filepath.Join(root, skill, "SKILL.md") }, func(tree, path, content string) {
			desc := frontmatterField(content, "description")
			if len(desc) >= 1024 {
				t.Errorf("%s: %s description is %d characters", tree, path, len(desc))
			}
			want := map[string]bool{}
			for _, w := range routed[skill] {
				want[w] = true
			}
			if len(want) == 0 {
				t.Fatalf("%s: no Next words found for %s in the Workflow", tree, skill)
			}
			for _, w := range known {
				has := strings.Contains(desc, "`"+w+"`")
				if want[w] && !has {
					t.Errorf("%s: %s description is missing Next word %q", tree, path, w)
				}
				if !want[w] && has {
					t.Errorf("%s: %s description names %q, which Workflow routes elsewhere", tree, path, w)
				}
			}
		})
	}
	forEachSkillFile(t, func(root string) string { return filepath.Join(root, "savepoint-idea", "SKILL.md") }, func(tree, path, content string) {
		if !strings.Contains(frontmatterField(content, "description"), "router state is idea") {
			t.Errorf("%s: %s description lost 'router state is idea'", tree, path)
		}
	})
}

func TestUpgradeDeliversRevisedSkillsAndRuleBlock(t *testing.T) {
	root := filepath.Join("..", "..")
	templates := os.DirFS(filepath.Join(root, "templates", "project-v2"))
	target := savepointProject(t)
	for _, skill := range v2Skills {
		dir := filepath.Join(target, "agent-skills", skill)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# previous release\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	old := managedBegin + "\n# Agents Guide\n\n## CLI Rules\n\nold\n" + managedEnd + "\n"
	if err := os.WriteFile(filepath.Join(target, "AGENTS.md"), []byte(old), 0644); err != nil {
		t.Fatal(err)
	}

	report := mustUpgrade(t, templates, target)

	for _, skill := range v2Skills {
		path := "agent-skills/" + skill + "/SKILL.md"
		if got := upgradeActionFor(t, report, path); got != ActionUpdated {
			t.Errorf("%s action = %v, want updated", path, got)
		}
		got := readTemplate(t, target, "agent-skills", skill, "SKILL.md")
		if got != readTemplate(t, root, "templates", "project-v2", "agent-skills", skill, "SKILL.md") {
			t.Errorf("%s not refreshed", path)
		}
	}
	assertContains(t, readTemplate(t, target, "AGENTS.md"), writeResumeFixRule)
}
