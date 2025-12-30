package parser

import (
	"bytes"
	"encoding/json"
	"fmt"

	"code/internal/domain"
)

func parseJSON(data []byte) (domain.Node, error) {
	dec := json.NewDecoder(bytes.NewReader(data))

	var out domain.Node

	if err := dec.Decode(&out); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidJSON, err)
	}

	return out, nil
}
