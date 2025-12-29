package parser

import "fmt"

func normalizeValue(value any) (any, error) {
	switch typed := value.(type) {
	case nil, string, bool, float64:
		return typed, nil

	case int:
		return float64(typed), nil
	case int64:
		return float64(typed), nil

	case Node:
		return normalizeNode(typed)

	case map[any]any:
		node, err := convertToNode(typed)
		if err != nil {
			return nil, err
		}

		return normalizeNode(node)

	case []any:
		return normalizeSlice(typed)

	default:
		return value, nil
	}
}

func normalizeNode(node Node) (Node, error) {
	out := make(Node, len(node))

	for key, value := range node {
		normalized, err := normalizeValue(value)
		if err != nil {
			return nil, err
		}

		out[key] = normalized
	}

	return out, nil
}

func normalizeSlice(items []any) ([]any, error) {
	out := make([]any, len(items))

	for i, item := range items {
		normalized, err := normalizeValue(item)
		if err != nil {
			return nil, err
		}

		out[i] = normalized
	}

	return out, nil
}

func convertToNode(m map[any]any) (Node, error) {
	out := make(Node, len(m))

	for key, value := range m {
		keyString, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("%w: non-string key %T", ErrInvalidYAML, key)
		}

		out[keyString] = value
	}

	return out, nil
}

