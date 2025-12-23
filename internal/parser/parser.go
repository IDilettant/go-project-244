package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	switch ext {
	case ".json":
		return parseJSON(data)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, ext)
	}
}
