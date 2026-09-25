package data

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// spliceV2Frontmatter rebuilds a record's frontmatter text after a managed
// patch while keeping the source bytes of every top-level field no patch
// touched. Re-encoding the whole document would restyle untouched fields
// (indentation, wrapped scalars), so an append-only history entry would
// still show every old entry as changed. A touched sequence that only gained
// items keeps its existing lines and appends the new items at the source's
// sequence indent. It reports false when the source text cannot be mapped
// onto the node, and callers then fall back to encoding the whole document.
func spliceV2Frontmatter(text string, original, patched *yaml.Node, touched map[string]bool, indent int) (string, bool) {
	if text == "" || !isV2MappingDocument(original) || !isV2MappingDocument(patched) {
		return "", false
	}
	originalMapping, patchedMapping := original.Content[0], patched.Content[0]
	lines := strings.Split(text, "\n")

	type fieldSpan struct {
		start, end int
		value      *yaml.Node
	}
	spans := make(map[string]fieldSpan)
	prefixEnd := len(lines)
	for i := 0; i+1 < len(originalMapping.Content); i += 2 {
		key := originalMapping.Content[i]
		if key.Line < 1 || key.Column != 1 {
			return "", false
		}
		start, end := key.Line-1, len(lines)
		if i+2 < len(originalMapping.Content) {
			end = originalMapping.Content[i+2].Line - 1
		}
		if _, duplicate := spans[key.Value]; duplicate || end <= start || end > len(lines) {
			return "", false
		}
		if i == 0 {
			prefixEnd = start
		}
		spans[key.Value] = fieldSpan{start: start, end: end, value: originalMapping.Content[i+1]}
	}

	out := append([]string{}, lines[:prefixEnd]...)
	for i := 0; i+1 < len(patchedMapping.Content); i += 2 {
		key, value := patchedMapping.Content[i].Value, patchedMapping.Content[i+1]
		span, existed := spans[key]
		if existed && !touched[key] {
			out = append(out, lines[span.start:span.end]...)
			continue
		}
		if existed {
			if appended, ok := appendedV2SequenceLines(key, span.value, value, lines[span.start:span.end]); ok {
				out = append(out, appended...)
				continue
			}
		}
		encoded, ok := encodeV2Field(key, value, indent)
		if !ok {
			return "", false
		}
		out = append(out, encoded...)
	}
	return strings.TrimSpace(strings.Join(out, "\n")), true
}

// appendedV2SequenceLines keeps a block sequence's source lines when the
// patch only appended items, adding the new items at the sequence's own
// indent. It reports false for any other change, a flow sequence, or an
// indentless sequence the encoder cannot reproduce.
func appendedV2SequenceLines(key string, original, patched *yaml.Node, span []string) ([]string, bool) {
	if original == nil || patched == nil || original.Kind != yaml.SequenceNode || patched.Kind != yaml.SequenceNode {
		return nil, false
	}
	if original.Style&yaml.FlowStyle != 0 || len(original.Content) == 0 || len(patched.Content) <= len(original.Content) {
		return nil, false
	}
	indent := original.Column - 1
	if indent < 2 {
		return nil, false
	}
	for i, item := range original.Content {
		if !nodesEqualByEncoding(item, patched.Content[i]) {
			return nil, false
		}
	}

	added := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: patched.Content[len(original.Content):]}
	encoded, ok := encodeV2Field(key, added, indent)
	if !ok || len(encoded) < 2 || encoded[0] != key+":" {
		return nil, false
	}

	kept := len(span)
	for kept > 0 && strings.TrimSpace(span[kept-1]) == "" {
		kept--
	}
	out := append([]string{}, span[:kept]...)
	out = append(out, encoded[1:]...)
	return append(out, span[kept:]...), true
}

// encodeV2Field encodes one top-level key and its value as frontmatter lines.
func encodeV2Field(key string, value *yaml.Node, indent int) ([]string, bool) {
	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{{
		Kind:    yaml.MappingNode,
		Tag:     "!!map",
		Content: []*yaml.Node{{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}, value},
	}}}
	text, err := encodeV2Document(doc, indent)
	if err != nil {
		return nil, false
	}
	return strings.Split(text, "\n"), true
}

// encodeV2Document encodes a whole frontmatter document at indent, trimmed
// of the encoder's surrounding whitespace.
func encodeV2Document(doc *yaml.Node, indent int) (string, error) {
	var encoded strings.Builder
	encoder := yaml.NewEncoder(&encoded)
	encoder.SetIndent(indent)
	if err := encoder.Encode(doc); err != nil {
		return "", err
	}
	if err := encoder.Close(); err != nil {
		return "", err
	}
	return strings.TrimSpace(encoded.String()), nil
}

func isV2MappingDocument(node *yaml.Node) bool {
	return node != nil && node.Kind == yaml.DocumentNode && len(node.Content) == 1 && node.Content[0].Kind == yaml.MappingNode
}
