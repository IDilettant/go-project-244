package parser

import (
	"encoding/json"
	"fmt"
)

type Node map[string]any

func parseJSON(data []byte) (Node, error) {
	var out Node

	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidJSON, err)
	}

	return out, nil
}
