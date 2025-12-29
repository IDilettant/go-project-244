package parser

import "errors"

var (
	ErrUnsupportedFormat = errors.New("unsupported file format")

	ErrInvalidJSON = errors.New("invalid json")
	ErrInvalidYAML = errors.New("invalid yaml")
)
