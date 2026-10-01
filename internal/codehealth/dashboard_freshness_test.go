package codehealth

import (
	"context"
	"strings"
	"testing"
)

// gitScript answers the read-only Git queries the observation and relation make.
type gitScript struct {
	head      string
	shallow   bool
	untracked string          // NUL-separated paths
	ancestors map[string]bool // "a>b": a is an ancestor of b
	counts    string
	err       error
}

func (g gitScript) reply(args []string) (CommandResult, error) {
	if g.err != nil {
		return CommandResult{}, g.err
	}
	out := func(s string) (CommandResult, error) { return CommandResult{Stdout: []byte(s)}, nil }
	switch strings.Join(args[:2], " ") {
	case "rev-parse --is-inside-work-tree":
		return out("true\n")
	case "rev-parse --verify":
		return out(g.head + "\n")
	case "rev-parse --is-shallow-repository":
		if g.shallow {
			return out("true\n")
		}
		return out("false\n")
	case "ls-files -z":
		if args[2] == "--others" {
			return out(g.untracked)
		}
		return out("")
	case "diff --relative":
		return out("")
	case "cat-file -e":
		return out("")
	case "merge-base --is-ancestor":
		if g.ancestors[args[2]+">"+args[3]] {
			return out("")
		}
		return CommandResult{ExitCode: 1}, nil
	case "rev-list --left-right":
		return out(g.counts + "\n")
	}
	return CommandResult{}, nil
}

func (g gitScript) Run(_ context.Context, _ string, args ...string) (CommandResult, error) {
	return g.reply(args)
}

func recordedAt(t *testing.T, root string, g gitScript) RepositoryIdentity {
	t.Helper()
	obs, err := ObserveRepository(context.Background(), g, root, InputScope{})
	if err != nil {
		t.Fatal(err)
	}
	return obs.Identity
}

func TestDashboardFreshness(t *testing.T) {
	c1, c2 := strings.Repeat("a", 40), strings.Repeat("b", 40)
	cases := map[string]struct {
		recorded func(t *testing.T, root string) RepositoryIdentity
		now      gitScript
		want     CodeState
		commits  int
		counted  bool
		differs  bool
		text     string
	}{
		"same commit": {
			recorded: func(t *testing.T, root string) RepositoryIdentity { return recordedAt(t, root, gitScript{head: c1}) },
			now:      gitScript{head: c1}, want: CodeMatches, text: "current commit",
		},
		"dirty since": {
			recorded: func(t *testing.T, root string) RepositoryIdentity { return recordedAt(t, root, gitScript{head: c1}) },
			now:      gitScript{head: c1, untracked: "new.go\x00"}, want: CodeMatches, differs: true, text: "have changed",
		},
		"behind with count": {
			recorded: func(t *testing.T, root string) RepositoryIdentity { return recordedAt(t, root, gitScript{head: c1}) },
			now:      gitScript{head: c2, ancestors: map[string]bool{c1 + ">" + c2: true}, counts: "0\t3"},
			want:     CodeMovedOn, commits: 3, counted: true, text: "3 commits behind",
		},
		"behind in a shallow clone has no count": {
			recorded: func(t *testing.T, root string) RepositoryIdentity { return recordedAt(t, root, gitScript{head: c1}) },
			now:      gitScript{head: c2, shallow: true, ancestors: map[string]bool{c1 + ">" + c2: true}},
			want:     CodeMovedOn,
		},
		"ahead is another branch": {
			recorded: func(t *testing.T, root string) RepositoryIdentity { return recordedAt(t, root, gitScript{head: c2}) },
			now:      gitScript{head: c1, ancestors: map[string]bool{c1 + ">" + c2: true}, counts: "2\t0"},
			want:     CodeOtherBranch, text: "2 commits ahead",
		},
		"diverged is another branch": {
			recorded: func(t *testing.T, root string) RepositoryIdentity { return recordedAt(t, root, gitScript{head: c1}) },
			now:      gitScript{head: c2, counts: "2\t4"},
			want:     CodeOtherBranch, text: "diverged",
		},
		"unborn repository is unknown": {
			recorded: func(t *testing.T, root string) RepositoryIdentity { return recordedAt(t, root, gitScript{head: c1}) },
			now:      gitScript{head: ""}, want: CodeUnknown, text: "without a commit",
		},
		"git unavailable is unknown": {
			recorded: func(t *testing.T, root string) RepositoryIdentity { return recordedAt(t, root, gitScript{head: c1}) },
			now:      gitScript{err: ErrGitMissing}, want: CodeUnknown, text: "Git is not installed",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			if tc.now.untracked != "" {
				writeFile(t, root+"/new.go", "package x\n")
			}
			recorded := tc.recorded(t, root)
			got := DashboardFreshness(context.Background(), root, tc.now, recorded)
			if got.State != tc.want || got.Commits != tc.commits || got.CommitsKnown != tc.counted || got.WorkingTreeDiffers != tc.differs {
				t.Fatalf("freshness = %+v, want state %s commits %d known %v differs %v", got, tc.want, tc.commits, tc.counted, tc.differs)
			}
			if !strings.Contains(got.Text, tc.text) {
				t.Errorf("text = %q, want it to contain %q", got.Text, tc.text)
			}
			assertNoBannedClaims(t, name, strings.ToLower(got.Text))
		})
	}
}

func TestDashboardFreshnessCancelledIsUnknown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	g := gitScript{err: context.Canceled}
	got := DashboardFreshness(ctx, t.TempDir(), g, RepositoryIdentity{})
	if got.State != CodeUnknown || got.Text == "" {
		t.Fatalf("freshness = %+v, want unknown with text", got)
	}
}
