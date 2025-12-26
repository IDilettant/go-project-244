package stylish

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"code/internal/diff"
	"code/internal/formatter"
)

type Formatter struct{}

func New() *Formatter {
	return &Formatter{}
}

func (f *Formatter) Format(changes []diff.Change) string {
	var builder strings.Builder

	builder.WriteString(formatter.OpeningBrace)
	builder.WriteString(formatter.NewLine)

	for _, change := range changes {
		f.writeChange(&builder, change)
	}

	builder.WriteString(formatter.ClosingBrace)

	return builder.String()
}

func (f *Formatter) writeChange(builder *strings.Builder, change diff.Change) {
	switch change.Type {
	case diff.Unchanged:
		f.writeUnchangedLine(builder, change.Key, change.OldValue)

	case diff.Removed:
		f.writeSignedLine(builder, formatter.SignRemoved, change.Key, change.OldValue)

	case diff.Added:
		f.writeSignedLine(builder, formatter.SignAdded, change.Key, change.NewValue)

	case diff.Updated:
		f.writeSignedLine(builder, formatter.SignRemoved, change.Key, change.OldValue)
		f.writeSignedLine(builder, formatter.SignAdded, change.Key, change.NewValue)
	}
}

func (f *Formatter) writeUnchangedLine(builder *strings.Builder, key string, value any) {
	builder.WriteString(formatter.IndentForUnchanged)
	builder.WriteString(key)
	builder.WriteString(formatter.ColonSpace)
	builder.WriteString(f.formatValue(value))
	builder.WriteString(formatter.NewLine)
}

func (f *Formatter) writeSignedLine(builder *strings.Builder, sign string, key string, value any) {
	builder.WriteString(formatter.IndentBase)
	builder.WriteString(sign)
	builder.WriteString(formatter.Space)
	builder.WriteString(key)
	builder.WriteString(formatter.ColonSpace)
	builder.WriteString(f.formatValue(value))
	builder.WriteString(formatter.NewLine)
}

func (f *Formatter) formatValue(value any) string {
	switch typedValue := value.(type) {
	case nil:
		return formatter.NullString

	case string:
		return typedValue

	case bool:
		return strconv.FormatBool(typedValue)

	case json.Number:
		return typedValue.String()

	default:
		return fmt.Sprintf("%v", typedValue)
	}
}
