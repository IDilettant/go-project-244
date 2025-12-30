package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"code/internal/domain"
)

func parseYAML(data []byte) (domain.Node, error) {
	var out domain.Node

	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidYAML, err)
	}

	return out, nil
}
