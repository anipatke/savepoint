package migrate

import "testing"

// TestClassify_matchesFixtureManifests reads both frozen fixture manifests as
// the expectation and asserts Classify agrees file by file, naming any file
// the manifest covers that classification misses.
func TestClassify_matchesFixtureManifests(t *testing.T) {
	for _, fixture := range []string{"v1-basic", "v1-history"} {
		t.Run(fixture, func(t *testing.T) {
			manifest := loadFixtureManifest(t, fixture)
			if len(manifest.Files) == 0 {
				t.Fatalf("manifest for fixture %s lists no files", fixture)
			}
			for _, f := range manifest.Files {
				relPath := fixtureRelPath(f.Path)
				got := Classify(relPath)
				if string(got) != f.Role {
					t.Errorf("Classify(%s) = %s, want %s (from %s manifest.yml)", relPath, got, f.Role, fixture)
				}
			}
		})
	}
}

// TestClassify_unclassifiedIsReportedNotDropped proves a file matching no
// known role still gets a value back rather than being silently guessed at.
func TestClassify_unclassifiedIsReportedNotDropped(t *testing.T) {
	cases := []string{
		".savepoint/scratch-notes.txt",
		"README.md",
		".savepoint/releases/v1",
		".savepoint/releases/v1/random.md",
	}
	for _, path := range cases {
		want := RoleUnclassified
		if path == ".savepoint/releases/v1" {
			want = RoleRelease
		}
		if got := Classify(path); got != want {
			t.Errorf("Classify(%s) = %s, want %s", path, got, want)
		}
	}
}

// TestClassify_shippedSkill covers the one role vocabulary entry neither
// fixture's project happens to carry: a shipped skill under agent-skills/.
func TestClassify_shippedSkill(t *testing.T) {
	if got := Classify("agent-skills/savepoint-build-task/SKILL.md"); got != RoleSkill {
		t.Errorf("Classify(agent-skills/.../SKILL.md) = %s, want %s", got, RoleSkill)
	}
}

// TestClassify_roleVocabularyIsComplete pins the full frozen role vocabulary
// the task's acceptance criteria enumerate, so a future accidental rename or
// removal fails a test rather than silently narrowing the vocabulary.
func TestClassify_roleVocabularyIsComplete(t *testing.T) {
	want := []Role{
		RoleConfig, RoleRouter, RoleProductPRD, RoleArchitecture, RoleHealthCheck,
		RoleRelease, RoleReleasePRD, RoleEpicDetail, RoleEpicAudit, RoleTask, RoleDefect,
		RoleAuditPrompt, RoleAuditRegister, RoleFinding, RoleAuditRun,
		RoleManagedGuide, RoleSkill,
	}
	have := map[Role]bool{}
	for _, rule := range classifyRules {
		have[rule.role] = true
	}
	for _, role := range want {
		if !have[role] {
			t.Errorf("role %s has no classifyRules entry", role)
		}
	}
	if len(have) != len(want) {
		t.Errorf("classifyRules has %d distinct roles, want exactly %d: %v", len(have), len(want), want)
	}
}
