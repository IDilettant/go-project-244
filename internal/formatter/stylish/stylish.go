package stylish

import (
	"fmt"
	"strconv"
	"strings"

	"code/internal/diff"
	"code/internal/domain"
	"code/internal/formatter"
)

type Formatter struct{}

func New() *Formatter { return &Formatter{} }

const (
	rootDepth  = 1
	indentStep = 4

	plainShift  = 2
	signedShift = 4
)

func (f *Formatter) Format(changes []diff.Change) string {
	var b strings.Builder
	b.WriteString(formatter.OpeningBrace)
	b.WriteString(formatter.NewLine)

	f.writeChanges(&b, changes, rootDepth)

	b.WriteString(formatter.ClosingBrace)

	return b.String()
}

func (f *Formatter) writeChanges(b *strings.Builder, changes []diff.Change, depth int) {
	for _, ch := range changes {
		f.writeChange(b, ch, depth)
	}
}

func (f *Formatter) writeChange(b *strings.Builder, ch diff.Change, depth int) {
	if ch.IsContainer() {
		f.writeNested(b, ch, depth)

		return
	}

	switch ch.Type {
	case diff.Unchanged:
		f.writeLine(b, depth, "", ch.Key, ch.OldValue)

	case diff.Removed:
		f.writeLine(b, depth, formatter.SignRemoved, ch.Key, ch.OldValue)

	case diff.Added:
		f.writeLine(b, depth, formatter.SignAdded, ch.Key, ch.NewValue)

	case diff.Updated:
		f.writeLine(b, depth, formatter.SignRemoved, ch.Key, ch.OldValue)
		f.writeLine(b, depth, formatter.SignAdded, ch.Key, ch.NewValue)
	}
}

func (f *Formatter) writeNested(b *strings.Builder, ch diff.Change, depth int) {
	b.WriteString(f.keyIndent(depth))
	b.WriteString(ch.Key)
	b.WriteString(formatter.ColonSpace)
	b.WriteString(formatter.OpeningBrace)
	b.WriteString(formatter.NewLine)

	f.writeChanges(b, ch.Children, depth+1)

	b.WriteString(f.keyIndent(depth))
	b.WriteString(formatter.ClosingBrace)
	b.WriteString(formatter.NewLine)
}

func (f *Formatter) writeLine(b *strings.Builder, depth int, sign string, key string, value any) {
	b.WriteString(f.linePrefix(depth, sign))
	b.WriteString(key)
	b.WriteString(formatter.ColonSpace)
	b.WriteString(f.renderValue(value, depth))
	b.WriteString(formatter.NewLine)
}

func (f *Formatter) linePrefix(depth int, sign string) string {
	if sign == "" {
		return f.keyIndent(depth)
	}

	base := strings.Repeat(formatter.Space, depth*indentStep-signedShift)

	return base + sign + formatter.Space
}

func (f *Formatter) keyIndent(depth int) string {
	return strings.Repeat(formatter.Space, depth*indentStep-plainShift)
}

func (f *Formatter) renderValue(v any, depth int) string {
	switch x := v.(type) {
	case nil:
		return formatter.NullString
	case string:
		return x
	case bool:
		return strconv.FormatBool(x)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)

	case domain.Node:
		return f.renderNode(x, depth)

	case map[string]any:
		return f.renderNode(domain.Node(x), depth)

	default:
		return fmt.Sprintf("%v", x)
	}
}

func (f *Formatter) renderNode(obj domain.Node, depth int) string {
	if len(obj) == 0 {
		return formatter.OpeningBrace + formatter.ClosingBrace
	}

	next := depth + 1

	var b strings.Builder
	b.WriteString(formatter.OpeningBrace)
	b.WriteString(formatter.NewLine)

	for _, k := range obj.KeysSorted() {
		b.WriteString(f.keyIndent(next))
		b.WriteString(k)
		b.WriteString(formatter.ColonSpace)
		b.WriteString(f.renderValue(obj[k], next))
		b.WriteString(formatter.NewLine)
	}

	b.WriteString(f.keyIndent(depth))
	b.WriteString(formatter.ClosingBrace)

	return b.String()
}

