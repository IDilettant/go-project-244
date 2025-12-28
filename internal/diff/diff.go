package diff

import (
	"reflect"
	"sort"

	"code/internal/parser"
)

// Compare computes a diff between two flat JSON nodes
func Compare(leftNode, rightNode parser.Node) []Change {
	sortedKeys := getSortedUnionKeys(leftNode, rightNode)

	changes := make([]Change, 0, len(sortedKeys))
	for _, key := range sortedKeys {
		change := buildChangeForKey(key, leftNode, rightNode)
		changes = append(changes, change)
	}

	return changes
}

func getSortedUnionKeys(leftNode, rightNode parser.Node) []string {
	uniqueKeys := make(map[string]struct{}, len(leftNode)+len(rightNode))

	for key := range leftNode {
		uniqueKeys[key] = struct{}{}
	}

	for key := range rightNode {
		uniqueKeys[key] = struct{}{}
	}

	keys := make([]string, 0, len(uniqueKeys))

	for key := range uniqueKeys {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func buildChangeForKey(key string, leftNode, rightNode parser.Node) Change {
	leftValue, leftExists := leftNode[key]
	rightValue, rightExists := rightNode[key]

	switch {
	case hasOnlyInLeft(leftExists, rightExists):
		return Change{
			Key:      key,
			Type:     Removed,
			OldValue: leftValue,
			NewValue: nil,
		}

	case hasOnlyInRight(leftExists, rightExists):
		return Change{
			Key:      key,
			Type:     Added,
			OldValue: nil,
			NewValue: rightValue,
		}

	case isEqualValues(leftValue, rightValue):
		return Change{
			Key:      key,
			Type:     Unchanged,
			OldValue: leftValue,
			NewValue: rightValue,
		}

	default:
		return Change{
			Key:      key,
			Type:     Updated,
			OldValue: leftValue,
			NewValue: rightValue,
		}
	}
}

func hasOnlyInLeft(leftExists, rightExists bool) bool {
	return leftExists && !rightExists
}

func hasOnlyInRight(leftExists, rightExists bool) bool {
	return !leftExists && rightExists
}

func isEqualValues(leftValue, rightValue any) bool {
	return reflect.DeepEqual(leftValue, rightValue)
}
