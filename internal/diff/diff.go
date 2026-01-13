package diff

import (
	"reflect"

	"code/internal/domain"
)

// Compare computes a diff between two nodes (supports nesting via Children).
func Compare(leftNode, rightNode domain.Node) []Change {
	sortedKeys := leftNode.UnionKeysSorted(rightNode)

	changes := make([]Change, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		changes = append(changes, buildChangeForKey(key, leftNode, rightNode))
	}

	return changes
}

func buildChangeForKey(key string, leftNode, rightNode domain.Node) Change {
	leftValue, leftExists := leftNode[key]
	rightValue, rightExists := rightNode[key]

	if hasOnlyInLeft(leftExists, rightExists) {
		return Change{
			Key:      key,
			Type:     Removed,
			OldValue: leftValue,
			NewValue: nil,
		}
	}

	if hasOnlyInRight(leftExists, rightExists) {
		return Change{
			Key:      key,
			Type:     Added,
			OldValue: nil,
			NewValue: rightValue,
		}
	}

	// If both values are containers, diff them recursively.
	leftChild, isNodeLeft := leftNode.GetNode(key)
	rightChild, isNodeRight := rightNode.GetNode(key)

	if isNodeLeft && isNodeRight {
		return Change{
			Key:      key,
			Type:     Nested,
			Children: Compare(leftChild, rightChild),
		}
	}

	if isEqualValues(leftValue, rightValue) {
		return Change{
			Key:      key,
			Type:     Unchanged,
			OldValue: leftValue,
			NewValue: rightValue,
		}
	}

	return Change{
		Key:      key,
		Type:     Updated,
		OldValue: leftValue,
		NewValue: rightValue,
	}
}

func isEqualValues(leftValue, rightValue any) bool {
	return reflect.DeepEqual(leftValue, rightValue)
}

func hasOnlyInLeft(leftExists, rightExists bool) bool {
	return leftExists && !rightExists
}

func hasOnlyInRight(leftExists, rightExists bool) bool {
	return !leftExists && rightExists
}
