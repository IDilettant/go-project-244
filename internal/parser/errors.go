package parser

import "errors"

var (
	ErrUnsupportedFormat = errors.New("unsupported file format")
	ErrUnsupportedType   = errors.New("unsupported type")

	ErrInvalidJSON = errors.New("invalid json")
	ErrInvalidYAML = errors.New("invalid yaml")
)
