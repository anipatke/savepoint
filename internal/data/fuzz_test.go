package data

import (
	"testing"
)

func FuzzExtractFrontmatter(f *testing.F) {
	seeds := []string{
		"---\nid: E01/T001\nstatus: planned\n---\nbody",
		"---\n---\n",
		"---\n\n---\n",
		"---\nid: test\n---",
		"",
		"# no frontmatter",
		"---\nid: [broken\n---\n",
		"---\nname: héllo wörld\n---\n",
		"---\nid: test\nstatus: in_progress\nphase: build\n---\nbody content",
		"---\r\nid: test\r\n---\r\nbody",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, content string) {
		_, _ = extractFrontmatter(content)
	})
}

func FuzzParseFrontmatter(f *testing.F) {
	seeds := []string{
		"---\nid: E01/T001\nstatus: planned\n---\nbody",
		"---\n---\n",
		"---\nid: [broken\n---\n",
		"---\nname: héllo\n---\n",
		"",
		"no frontmatter",
		"---\ntags: [a, b, c]\n---\n",
		"---\nnested:\n  key: val\n---\n",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, content string) {
		p := NewParser()
		_, _ = p.ParseFrontmatter(content)
	})
}

func FuzzSplitFrontmatterBody(f *testing.F) {
	seeds := []string{
		"---\nid: E01/T001\nstatus: planned\n---\nbody",
		"---\n---\n",
		"---\nkey: value\n---",
		"",
		"# no frontmatter",
		"---\nid: test\nstatus: in_progress\n---\n\n## Section\n\nContent.",
		"---\nid: test\n---\n\nbody with unicode: 日本語",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, content string) {
		yamlStr, body, err := SplitFrontmatterBody(content)
		if err != nil {
			return
		}
		reconstructed := "---\n" + yamlStr + "\n---" + body
		_, body2, err2 := SplitFrontmatterBody(reconstructed)
		if err2 != nil {
			t.Errorf("round-trip SplitFrontmatterBody failed on reconstructed: %v", err2)
		}
		if body2 != body {
			t.Errorf("round-trip body mismatch: got %q, want %q", body2, body)
		}
	})
}

func FuzzParseFindingFile(f *testing.F) {
	seeds := []string{
		"---\nid: F001\ntitle: \"Finding\"\nstatus: open\nseverity: medium\nconfidence: high\nproof_needed: test\nfirst_seen: \"2026-07-01\"\nlast_seen: \"2026-07-01\"\n---\n\n# Finding\n",
		"---\nid: F001\nstatus: bogus\nseverity: blocker\nconfidence: certain\n---\n",
		"---\nid: [broken\n---\n",
		"---\n---\n",
		"",
		"# no frontmatter",
		"---\nid: F999\nduplicate_of: F999\nstatus: duplicate\n---\n",
		"---\ntasks: [\"E01-foo/T001\"]\ndefects: [\"D001-bug\"]\nreleases: [v1]\n---\n",
		"---\nid: 12345\nfirst_seen: 2026-07-01\n---\n",
	}
	for _, s := range seeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, content string) {
		p := NewParser()
		finding, err := p.ParseFindingFile("F001-fuzz.md", content)
		if err != nil {
			return
		}
		// Load-time healing must leave only canonical enum values.
		if !isValidFindingStatus(finding.Status) {
			t.Errorf("healed finding has non-canonical status %q", finding.Status)
		}
		if !isValidFindingSeverity(finding.Severity) {
			t.Errorf("healed finding has non-canonical severity %q", finding.Severity)
		}
		if !isValidFindingConfidence(finding.Confidence) {
			t.Errorf("healed finding has non-canonical confidence %q", finding.Confidence)
		}
		if finding.ID != "F001" {
			t.Errorf("healed finding ID = %q, want filename id F001", finding.ID)
		}
	})
}
