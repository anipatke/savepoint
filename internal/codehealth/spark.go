package codehealth

import (
	"fmt"
	"math"
	"slices"
)

// MaxSparkPoints bounds how many official checks one sparkline draws.
const MaxSparkPoints = 10

// sparkBlocks run from the smallest value in the window to the largest.
var sparkBlocks = []rune("▁▂▃▄▅▆▇█")

// sparkFlat is drawn for every point when the window moved less than the
// signal's material-movement size, so noise does not look like movement.
const sparkFlat = 3

// Every phrase of the sparkline wording lives here (STYLE-09).
var sparkWordText = map[TrendDirection]string{
	TrendImproving: "better",
	TrendSteady:    "steady",
	TrendDeclining: "worse",
}

const (
	textSparkThin      = "not enough history yet (%d of %d)"
	textSparkEarly     = "early"
	textSparkRestarted = "restarted because settings changed"
)

// sparkValues are the comparable official values the trend uses, oldest first,
// plus current when it is itself a complete official one, newest MaxSparkPoints.
func sparkValues(current CapabilityResult, origin Origin, s series) []float64 {
	values := slices.Clone(s.values)
	if origin == OriginOfficial && usableValue(current) {
		values = append(values, current.Value.Number)
	}
	if len(values) > MaxSparkPoints {
		values = values[len(values)-MaxSparkPoints:]
	}
	return values
}

// sparkline draws values (oldest first) as one block each, scaled to the
// window's own range, and says in words which way the trend went. It returns
// the drawing, the better/worse/steady word, and a note about thin or
// restarted history. The drawing is empty below minObservations points, and
// the word is empty when there is no trend.
func sparkline(c Capability, values []float64, dir TrendDirection, restarted bool) (spark, word, note string) {
	word = sparkWordText[dir]
	var notes []string
	switch n := len(values); {
	case n < minObservations:
		notes = append(notes, fmt.Sprintf(textSparkThin, n, minObservations))
	case n < minObservations+2:
		notes = append(notes, textSparkEarly)
	}
	if restarted {
		notes = append(notes, textSparkRestarted)
	}
	note = join(notes)
	if len(values) < minObservations {
		return "", word, note
	}
	lo, hi := slices.Min(values), slices.Max(values)
	spread := math.Round((hi-lo)*1e6) / 1e6
	out := make([]rune, len(values))
	for i, v := range values {
		switch {
		case spread < materialMovement[c], hi == lo:
			out[i] = sparkBlocks[sparkFlat]
		default:
			out[i] = sparkBlocks[int(math.Round((v-lo)/(hi-lo)*float64(len(sparkBlocks)-1)))]
		}
	}
	return string(out), word, note
}
