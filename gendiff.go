package code

import (
	"code/internal/diff"
	"code/internal/domain"
	"code/internal/formatters"
)

func GenDiff(leftNode, rightNode domain.Node, f formatters.Formatter) string {
	changes := diff.Compare(leftNode, rightNode)

	return f.Format(changes)
}
