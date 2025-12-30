package code

import (
	"code/internal/diff"
	"code/internal/domain"
	"code/internal/formatter"
)

func GenDiff(leftNode, rightNode domain.Node, f formatter.Formatter) string {
	changes := diff.Compare(leftNode, rightNode)

	return f.Format(changes)
}
