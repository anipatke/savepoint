package data

import "testing"

func TestMatchesV2IdentityRequiresHyphenatedKindAndAtLeastThreeDigits(t *testing.T) {
	for _, kind := range []byte("ROTCI") {
		for _, id := range []string{string(kind) + "-001", string(kind) + "-1234"} {
			if !matchesV2Identity(id, kind) {
				t.Errorf("matchesV2Identity(%q, %q) = false, want true", id, kind)
			}
		}

		for _, id := range []string{
			string(kind) + "001",
			string(kind) + "-01",
			string(kind) + "--001",
			string(kind) + "-001x",
			"X-001",
		} {
			if matchesV2Identity(id, kind) {
				t.Errorf("matchesV2Identity(%q, %q) = true, want false", id, kind)
			}
		}
	}
}
