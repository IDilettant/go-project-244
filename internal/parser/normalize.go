package parser

import (
	"fmt"

	"code/internal/domain"
)

func normalizeNode(node domain.Node) (domain.Node, error) {
	out := make(domain.Node, len(node))

	for key, value := range node {
		normalized, err := normalizeValue(value)
		if err != nil {
			return nil, err
		}
		out[key] = normalized
	}

	return out, nil
}

func normalizeValue(v any) (any, error) {
	switch x := v.(type) {
	case nil, string, bool, float64:
		return x, nil

	case int:
		return float64(x), nil
	case int64:
		return float64(x), nil

	case []any:
		return normalizeSlice(x)

	case map[any]any:
		n, err := convertToNode(x)
		if err != nil {
			return nil, err
		}
		
		return normalizeNode(n)

	default:
		if n, ok := domain.AsNode(x); ok {
			return normalizeNode(n)
		}

		return nil, fmt.Errorf("%w: %T", ErrUnsupportedType, v)
	}
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

func convertToNode(m map[any]any) (domain.Node, error) {
	out := make(domain.Node, len(m))

	for key, value := range m {
		keyString, ok := key.(string)
		if !ok {
			return nil, fmt.Errorf("%w: key type %T", ErrUnsupportedType, key)
		}

		normalized, err := normalizeValue(value)
		if err != nil {
			return nil, err
		}

		out[keyString] = normalized
	}

	return out, nil
}

