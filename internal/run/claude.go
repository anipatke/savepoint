package run

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// ProviderClaude labels runs of Claude Code and selects ParseClaudeLine.
const ProviderClaude = "claude"

// Claude Code's `--output-format stream-json` is documented at
// https://code.claude.com/docs/en/headless. Every fact about its shape lives in
// this file; the shapes below were checked against Claude Code 2.1.296
// (testdata/claude-2.1.296).

// claudeToolKinds maps a tool name to what the agent did with it, and
// claudeTargetKeys lists the input fields that name its target, in order.
var (
	claudeToolKinds = map[string]EventKind{
		"Read": KindRead, "Edit": KindEdit, "Write": KindEdit, "MultiEdit": KindEdit,
		"NotebookEdit": KindEdit, "Bash": KindCommand,
	}
	claudeTargetKeys = []string{"file_path", "notebook_path", "command"}
)

// exitCodeLine matches the first line of a failed command's result text.
var exitCodeLine = regexp.MustCompile(`^Exit code (-?\d+)`)

const maxTargetLen = 512

type claudeLine struct {
	Type      string          `json:"type"`
	Subtype   string          `json:"subtype"`
	SessionID string          `json:"session_id"`
	Model     string          `json:"model"`
	Message   json.RawMessage `json:"message"` // an object for assistant/user, text for system

	ToolName       string `json:"tool_name"`
	ToolUseID      string `json:"tool_use_id"`
	DecisionReason string `json:"decision_reason"`

	Result       string          `json:"result"`
	Usage        json.RawMessage `json:"usage"`
	TotalCostUSD *float64        `json:"total_cost_usd"`
	Denials      []struct {
		ToolName  string         `json:"tool_name"`
		ToolUseID string         `json:"tool_use_id"`
		ToolInput map[string]any `json:"tool_input"`
	} `json:"permission_denials"`
}

type claudeMessage struct {
	Content json.RawMessage `json:"content"` // a list of blocks, or plain text
}

type claudeBlock struct {
	Type      string          `json:"type"`
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Input     map[string]any  `json:"input"`
	ToolUseID string          `json:"tool_use_id"`
	IsError   bool            `json:"is_error"`
	Content   json.RawMessage `json:"content"`
}

var unparsedLine = []Event{{Kind: KindUnparsed}}

// ParseClaudeLine interprets one stream-json line. It ignores the partial
// `stream_event` deltas and housekeeping `system` lines on purpose; malformed
// JSON and line types it does not know come back as KindUnparsed. The line
// itself stays in events.jsonl either way.
func ParseClaudeLine(line []byte) []Event {
	line = bytes.TrimSpace(line)
	if len(line) == 0 {
		return nil
	}
	var l claudeLine
	if err := json.Unmarshal(line, &l); err != nil || l.Type == "" {
		return unparsedLine
	}
	switch l.Type {
	case "system":
		return claudeSystem(l)
	case "assistant":
		return claudeBlocks(l.Message, claudeToolUses)
	case "user":
		return claudeBlocks(l.Message, claudeToolResults)
	case "result":
		return claudeFinal(l)
	case "stream_event":
		return nil
	}
	return unparsedLine
}

func claudeSystem(l claudeLine) []Event {
	switch l.Subtype {
	case "init":
		return []Event{{Kind: KindInit, SessionID: oneLine(l.SessionID, maxTargetLen), Model: oneLine(l.Model, maxTargetLen)}}
	case "permission_denied":
		return []Event{{Kind: KindDenied, Tool: oneLine(l.ToolName, maxTargetLen), ToolUseID: l.ToolUseID, Text: oneLine(l.DecisionReason, maxTargetLen)}}
	}
	if claudeHousekeeping[l.Subtype] {
		return nil // hooks, thinking tokens, retries: nothing the record keeps
	}
	return unparsedLine
}

// claudeHousekeeping lists the system subtypes ignored on purpose; any other
// subtype is counted as unparsed.
var claudeHousekeeping = map[string]bool{
	"hook_started": true, "hook_response": true, "hook_progress": true,
	"thinking_tokens": true, "status": true, "compact_boundary": true, "api_retry": true,
}

// claudeBlocks decodes a message's content blocks and maps each with fn. A
// message whose content is plain text has no blocks and yields no events.
func claudeBlocks(raw json.RawMessage, fn func(claudeBlock) (Event, bool)) []Event {
	var msg claudeMessage
	if err := json.Unmarshal(raw, &msg); err != nil || string(bytes.TrimSpace(raw)) == "null" {
		return unparsedLine
	}
	var text string
	if string(bytes.TrimSpace(msg.Content)) != "null" && json.Unmarshal(msg.Content, &text) == nil {
		return nil // plain text has no blocks
	}
	var blocks []claudeBlock
	if json.Unmarshal(msg.Content, &blocks) != nil || blocks == nil {
		return unparsedLine
	}
	var events []Event
	for _, b := range blocks {
		if ev, ok := fn(b); ok {
			events = append(events, ev)
		}
	}
	return events
}

func claudeToolUses(b claudeBlock) (Event, bool) {
	if b.Type != "tool_use" {
		return Event{}, false
	}
	kind, target := claudeTool(b.Name, b.Input)
	if target == "" {
		kind, target = KindTool, b.Name
	}
	return Event{Kind: kind, ToolUseID: b.ID, Tool: oneLine(b.Name, maxTargetLen), Text: oneLine(target, maxTargetLen)}, true
}

func claudeToolResults(b claudeBlock) (Event, bool) {
	if b.Type != "tool_result" {
		return Event{}, false
	}
	ev := Event{Kind: KindToolResult, ToolUseID: b.ToolUseID, Failed: b.IsError}
	if m := exitCodeLine.FindStringSubmatch(claudeText(b.Content)); m != nil {
		if code, err := strconv.Atoi(m[1]); err == nil {
			ev.ExitCode = &code
		}
	}
	return ev, true
}

// claudeText is a tool result's text: the string itself, or the text blocks
// joined.
func claudeText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var parts []struct {
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &parts) != nil {
		return ""
	}
	var out []string
	for _, p := range parts {
		out = append(out, p.Text)
	}
	return strings.Join(out, "\n")
}

// claudeTool says what a tool call did and to what. Its target is empty when
// the input names none.
func claudeTool(name string, input map[string]any) (EventKind, string) {
	kind, known := claudeToolKinds[name]
	if !known {
		return KindTool, ""
	}
	for _, key := range claudeTargetKeys {
		if s, ok := input[key].(string); ok && s != "" {
			return kind, s
		}
	}
	return kind, ""
}

func claudeFinal(l claudeLine) []Event {
	ev := Event{
		Kind:         KindFinal,
		SessionID:    oneLine(l.SessionID, maxTargetLen),
		Usage:        looseObject(l.Usage),
		TotalCostUSD: l.TotalCostUSD,
		Text:         cleanText(l.Result),
	}
	for _, d := range l.Denials {
		_, target := claudeTool(d.ToolName, d.ToolInput)
		ev.Denials = append(ev.Denials, Denial{Tool: oneLine(d.ToolName, maxTargetLen), Target: oneLine(target, maxTargetLen), ToolUseID: d.ToolUseID})
	}
	return []Event{ev}
}

// looseObject decodes a JSON object into plain Go values with whole numbers kept
// whole (a float64 would write a large count as 1e+06) and every string cleaned.
func looseObject(raw json.RawMessage) map[string]any {
	if len(raw) == 0 {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var m map[string]any
	if dec.Decode(&m) != nil {
		return nil
	}
	return tidy(m).(map[string]any)
}

func tidy(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, e := range x {
			x[k] = tidy(e)
		}
		return x
	case []any:
		for i, e := range x {
			x[i] = tidy(e)
		}
		return x
	case string:
		return cleanText(x)
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return x.String()
	}
	return v
}
