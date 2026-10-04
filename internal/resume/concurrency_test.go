package resume

import (
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestParallelLinesAreEmptyWhenOffOrMissing(t *testing.T) {
	index := &data.V2Index{}
	withheld := &data.ConcurrencyV2{Enabled: true, Withheld: &data.ConcurrencyReasonV2{Detail: "x"}}
	for name, got := range map[string][]string{
		"nil index":      ParallelLines(nil, withheld, ""),
		"nil projection": ParallelLines(index, nil, ""),
		"off":            ParallelLines(index, &data.ConcurrencyV2{Withheld: withheld.Withheld}, ""),
		"nothing to say": ParallelLines(index, &data.ConcurrencyV2{Enabled: true}, ""),
	} {
		if len(got) != 0 {
			t.Errorf("%s: ParallelLines() = %q, want nothing", name, got)
		}
	}
}

func TestParallelLinesStateWithheldAdviceAsOptionalAndSanitised(t *testing.T) {
	c := &data.ConcurrencyV2{Enabled: true, Withheld: &data.ConcurrencyReasonV2{Detail: "T-1 records a replan\x1b[2J\x07,\nso no launch is suggested"}}
	for _, focus := range []string{"", "T-1"} {
		text := strings.Join(ParallelLines(&data.V2Index{}, c, focus), "\n")
		if !strings.Contains(text, "optional") || !strings.Contains(text, "Withheld: T-1 records a replan, so no launch is suggested") {
			t.Errorf("focus %q: ParallelLines() = %q, want optional framing and the cleaned reason", focus, text)
		}
		if strings.ContainsAny(text, "\x1b\x07") {
			t.Errorf("focus %q: output keeps a terminal control: %q", focus, text)
		}
	}
}
