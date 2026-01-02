package plain

import (
	"fmt"
	"strconv"
	"strings"

	"code/internal/diff"
	"code/internal/domain"
	"code/internal/formatters/common"
)

type Formatter struct{}

func New() *Formatter { return &Formatter{} }

func (f *Formatter) Format(changes []diff.Change) string {
	var b strings.Builder
	f.writeChanges(&b, changes, "")
	return b.String()
}

func (f *Formatter) writeChanges(
	b *strings.Builder,
	changes []diff.Change,
	prefix string,
) {
	for _, ch := range changes {
		f.writeChange(b, ch, prefix)
	}
}

func (f *Formatter) writeChange(
	b *strings.Builder,
	ch diff.Change,
	prefix string,
) {
	path := joinPath(prefix, ch.Key)

	if ch.IsContainer() {
		f.writeChanges(b, ch.Children, path)
		return
	}

	switch ch.Type {

	case diff.Added:
		b.WriteString(common.Property)
		writeQuoted(b, path)
		b.WriteString(" was added with value: ")
		b.WriteString(renderValue(ch.NewValue))
		b.WriteString(common.NewLine)

	case diff.Removed:
		b.WriteString(common.Property)
		writeQuoted(b, path)
		b.WriteString(" was removed")
		b.WriteString(common.NewLine)

	case diff.Updated:
		b.WriteString(common.Property)
		writeQuoted(b, path)
		b.WriteString(" was updated. From ")
		b.WriteString(renderValue(ch.OldValue))
		b.WriteString(" to ")
		b.WriteString(renderValue(ch.NewValue))
		b.WriteString(common.NewLine)

	default:
		return
	}
}

func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + common.Dot + key
}

func writeQuoted(b *strings.Builder, s string) {
	b.WriteString(common.Quote)
	b.WriteString(s)
	b.WriteString(common.Quote)
}

func renderValue(v any) string {
	if v == nil {
		return common.NullString
	}

	switch x := v.(type) {
	case string:
		return common.Quote + x + common.Quote

	case bool:
		return strconv.FormatBool(x)

	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)

	case domain.Node:
		return common.ComplexValue

	case map[string]any:
		return common.ComplexValue

	default:
		if _, ok := domain.AsNode(v); ok {
			return common.ComplexValue
		}

		return fmt.Sprintf("%v", x)
	}
}
