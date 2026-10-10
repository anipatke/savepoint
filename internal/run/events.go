package run

// EventKind names what one provider line meant, in words that do not depend on
// the provider.
type EventKind string

const (
	KindInit       EventKind = "init"        // the session started
	KindRead       EventKind = "read"        // the agent read a file
	KindEdit       EventKind = "edit"        // the agent changed a file
	KindCommand    EventKind = "command"     // the agent ran a command
	KindTool       EventKind = "tool"        // any other tool call; Text is its name
	KindToolResult EventKind = "tool_result" // a tool call came back
	KindDenied     EventKind = "denied"      // the provider refused an action
	KindFinal      EventKind = "final"       // the provider's closing report
	KindUnparsed   EventKind = "unparsed"    // a line the adapter did not understand
)

// maxListed bounds each list kept in run.yml.
const maxListed = 500

// Event is one interpreted provider line. Only the fields that belong to Kind
// are set. Text is already cleaned for display (ARCH-05).
type Event struct {
	Kind      EventKind
	ToolUseID string
	Tool      string
	Text      string // path, command, tool name, denial reason, or final message
	ExitCode  *int   // KindToolResult, when the provider states one
	Failed    bool   // KindToolResult: the provider marked the call as an error

	SessionID string // KindInit, KindFinal
	Model     string // KindInit

	Usage        map[string]any // KindFinal, verbatim
	TotalCostUSD *float64       // KindFinal
	Denials      []Denial       // KindFinal: the provider's own list
}

// Adapter turns one stdout line into zero or more events. It is a pure
// function: it reads no files and keeps no state. A line it ignores on purpose
// returns nil; a line it cannot interpret returns one KindUnparsed event.
type Adapter func(line []byte) []Event

// adapters maps a provider label to its adapter. Provider-specific knowledge
// lives in one file per provider.
var adapters = map[string]Adapter{ProviderClaude: ParseClaudeLine}

// observer folds events into a Record. It remembers which tool call each id
// named so a later result or denial can be tied back to it.
type observer struct {
	targets map[string]string // tool_use_id -> path or command
	command map[string]int    // tool_use_id -> index into Record.Commands
	denied  map[string]bool   // tool_use_id of every recorded denial
}

func newObserver() *observer {
	return &observer{targets: map[string]string{}, command: map[string]int{}, denied: map[string]bool{}}
}

// apply updates rec from ev.
func (o *observer) apply(rec *Record, ev Event) {
	switch ev.Kind {
	case KindInit:
		setIf(&rec.SessionID, ev.SessionID)
		setIf(&rec.Model, ev.Model)
	case KindRead:
		o.remember(ev)
		rec.FilesRead = appendUnique(rec, rec.FilesRead, ev.Text)
	case KindEdit:
		o.remember(ev)
		rec.FilesEdited = appendUnique(rec, rec.FilesEdited, ev.Text)
	case KindCommand:
		o.remember(ev)
		if len(rec.Commands) >= maxListed {
			rec.ListsTruncated = true
			return
		}
		rec.Commands = append(rec.Commands, CommandRun{Command: ev.Text})
		if ev.ToolUseID != "" {
			o.command[ev.ToolUseID] = len(rec.Commands) - 1
		}
	case KindToolResult:
		i, ok := o.command[ev.ToolUseID]
		if !ok {
			return
		}
		switch {
		case ev.ExitCode != nil:
			rec.Commands[i].ExitCode = ev.ExitCode
		case !ev.Failed:
			// The provider reports every non-zero exit as an error, so a result
			// that is not an error is an exit 0.
			zero := 0
			rec.Commands[i].ExitCode = &zero
		}
	case KindDenied:
		o.deny(rec, Denial{Tool: ev.Tool, Target: o.targets[ev.ToolUseID], Reason: ev.Text, ToolUseID: ev.ToolUseID})
	case KindFinal:
		setIf(&rec.SessionID, ev.SessionID)
		rec.Usage = ev.Usage
		rec.TotalCostUSD = ev.TotalCostUSD
		rec.AgentReported = ev.Text
		for _, d := range ev.Denials {
			if d.Target == "" {
				d.Target = o.targets[d.ToolUseID]
			}
			o.deny(rec, d)
		}
	case KindUnparsed:
		rec.Unparsed++
	}
}

func (o *observer) remember(ev Event) {
	if ev.ToolUseID != "" {
		o.targets[ev.ToolUseID] = ev.Text
	}
}

// deny records a denial once per tool call: the live event and the closing
// report both describe the same refusal.
func (o *observer) deny(rec *Record, d Denial) {
	if d.ToolUseID != "" {
		if o.denied[d.ToolUseID] {
			return
		}
		o.denied[d.ToolUseID] = true
	}
	if len(rec.Denials) >= maxListed {
		rec.ListsTruncated = true
		return
	}
	rec.Denials = append(rec.Denials, d)
}

func setIf(dst *string, v string) {
	if v != "" {
		*dst = v
	}
}

func appendUnique(rec *Record, list []string, v string) []string {
	for _, have := range list {
		if have == v {
			return list
		}
	}
	if len(list) >= maxListed {
		rec.ListsTruncated = true
		return list
	}
	return append(list, v)
}
