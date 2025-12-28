package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Node map[string]any

func parseJSON(data []byte) (Node, error) {
	dec := json.NewDecoder(bytes.NewReader(data))

	var out Node

	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidJSON, err)
	}

	return out, nil
}
