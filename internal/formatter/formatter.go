package formatter

import (
	"code/internal/diff"
)

type Formatter interface {
	Format(changes []diff.Change) string
}
