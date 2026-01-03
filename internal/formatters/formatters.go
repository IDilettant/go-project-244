package formatters

import (
	"fmt"

	"code/internal/diff"
	"code/internal/formatters/common"
	"code/internal/formatters/plain"
	"code/internal/formatters/stylish"
)

const (
	emptyFlag = ""
)

type Formatter interface {
	Format(changes []diff.Change) string
}

func SelectFormatter(format string) (Formatter, error) {
	switch format {
	case emptyFlag, common.FormatStylish:
		return stylish.New(), nil
	case common.FormatPlain:
		return plain.New(), nil
	default:
		return nil, fmt.Errorf("%w: %s", common.ErrUnknownFormat, format)
	}
}
