package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func parseYAML(data []byte) (Node, error) {
	var out Node

	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidYAML, err)
	}

	return out, nil
}
