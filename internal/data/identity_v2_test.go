package data

import "testing"

func TestMatchesV2IdentityRequiresHyphenatedKindAndAtLeastThreeDigits(t *testing.T) {
	for _, kind := range []byte("RGOTCI") {
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

func TestMatchesGoalIdentityV2AcceptsHistoricalAndNewPrefixes(t *testing.T) {
	for _, id := range []string{"R-001", "R-1234", "G-001", "G-1234"} {
		if !matchesGoalIdentityV2(id) {
			t.Errorf("matchesGoalIdentityV2(%q) = false, want true", id)
		}
	}
	for _, id := range []string{"O-001", "G-01", "g-001", "G-001x"} {
		if matchesGoalIdentityV2(id) {
			t.Errorf("matchesGoalIdentityV2(%q) = true, want false", id)
		}
	}
}
