package data

import "regexp"

// v2IdentityPattern is the single syntax rule for numbered V2 record IDs.
// matchesV2Identity additionally checks the expected record kind at each use.
var v2IdentityPattern = regexp.MustCompile(`^[ROTCI]-[0-9]{3,}$`)

func matchesV2Identity(id string, kind byte) bool {
	return len(id) > 1 && id[0] == kind && v2IdentityPattern.MatchString(id)
}
