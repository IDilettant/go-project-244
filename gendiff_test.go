package code_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code"
	"code/internal/diff"
	"code/internal/parser"
)

type stubFormatter struct {
	called bool
	got    []diff.Change
	out    string
}

func (s *stubFormatter) Format(changes []diff.Change) string {
	s.called = true
	s.got = changes

	return s.out
}

func TestGenDiff_smoke(t *testing.T) {
	t.Parallel()

	left := parser.Node{"a": float64(1)}
	right := parser.Node{"a": float64(2)}

	sf := &stubFormatter{out: "RESULT"}

	got := code.GenDiff(left, right, sf)

	require.True(t, sf.called, "formatter must be called")
	require.NotEmpty(t, sf.got, "formatter must receive changes")
	require.Equal(t, "RESULT", got, "GenDiff must return formatter output")
}
