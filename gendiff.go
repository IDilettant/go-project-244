package code

import (
	"code/internal/diff"
	"code/internal/formatter"
	"code/internal/parser"
)

func GenDiff(leftNode, rightNode parser.Node, f formatter.Formatter) string {
	changes := diff.Compare(leftNode, rightNode)

	return f.Format(changes)
}
