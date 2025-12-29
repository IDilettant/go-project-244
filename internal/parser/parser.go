package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	extJSON = ".json"
	extYAML = ".yaml"
	extYML  = ".yml"
)

// ParseFile reads config from path and parses it according to file extension
func ParseFile(path string) (Node, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path %q: %w", path, err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", absPath, err)
	}

	ext := strings.ToLower(filepath.Ext(absPath))

	var node Node
	switch ext {
	case extJSON:
		node, err = parseJSON(data)
	case extYAML, extYML:
		node, err = parseYAML(data)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, ext)
	}
	if err != nil {
		return nil, err
	}

	return normalizeNode(node)
}
