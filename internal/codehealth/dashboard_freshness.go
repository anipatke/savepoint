package codehealth

import "context"

// CodeState says how the measured code relates to the code now.
type CodeState string

const (
	CodeMatches     CodeState = "matches"
	CodeMovedOn     CodeState = "moved_on"
	CodeOtherBranch CodeState = "other_branch"
	CodeUnknown     CodeState = "unknown"
)

// CodeFreshness is the plain answer to "is this still about my code?".
type CodeFreshness struct {
	State CodeState
	// Commits is how many commits the code has moved on by; it is meaningful
	// only when CommitsKnown is true.
	Commits      int
	CommitsKnown bool
	// WorkingTreeDiffers is true when the files in scope differ from what was
	// measured.
	WorkingTreeDiffers bool
	Text               string
}

// DashboardFreshness compares the repository state a dashboard was measured at
// with the repository now. It runs only read-only local Git queries. Any Git
// failure, including cancellation, yields CodeUnknown: the screen never has an
// error to handle.
func DashboardFreshness(ctx context.Context, root string, git CommandRunner, recorded RepositoryIdentity) CodeFreshness {
	if git == nil {
		git = GitRunner{}
	}
	obs, err := ObserveRepository(ctx, git, root, InputScope{})
	if err != nil {
		return CodeFreshness{State: CodeUnknown, Text: DescribeRepositoryError(err)}
	}
	rel, err := Relate(ctx, git, root, recorded, obs)
	if err != nil {
		return CodeFreshness{State: CodeUnknown, Text: DescribeRepositoryError(err)}
	}
	f := CodeFreshness{WorkingTreeDiffers: !rel.InputsMatch, Text: rel.Describe()}
	switch rel.Kind {
	case RelationSameCommit:
		f.State = CodeMatches
	case RelationBehind:
		f.State = CodeMovedOn
		f.Commits, f.CommitsKnown = rel.CurrentOnly, rel.Counted
	case RelationAhead, RelationDiverged:
		f.State = CodeOtherBranch
	default:
		f.State = CodeUnknown
	}
	return f
}
