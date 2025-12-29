package app

import (
	"errors"
	"fmt"
)

const (
	exitCodeUsageError   = 2
	exitCodeRuntimeError = 1
)

var (
	ErrUsage   = errors.New("usage error")
	ErrRuntime = errors.New("runtime error")
)

func invalidArgsError() error {
	return fmt.Errorf(
		"%w: expected %d arguments: <filepath1> <filepath2>",
		ErrUsage,
		requiredArgumentsCount,
	)
}

func wrap(kind, err error) error {
	return fmt.Errorf("%w: %w", kind, err)
}

