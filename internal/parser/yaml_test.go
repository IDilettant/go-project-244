package parser

import (
	"errors"
	"testing"

	"code/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestParseYAML(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		input     string
		want      domain.Node
		errAssert func(error) bool
	}

	tests := []testCase{
		{
			name:  "empty/object",
			input: "{}",
			want:  domain.Node{},
		},
		{
			name:  "types/primitives are parsed",
			input: "s: x\nb: true\nn: null\n",
			want: domain.Node{
				"s": "x",
				"b": true,
				"n": nil,
			},
		},
		{
			name:  "numbers/integer is decoded",
			input: "a: 1\n",
			want: domain.Node{
				"a": 1,
			},
		},
		{
			name:  "numbers/fractional is decoded",
			input: "a: 1.5\n",
			want: domain.Node{
				"a": 1.5,
			},
		},
		{
			name:  "numbers/equal magnitude values are decoded consistently",
			input: "a: 1\nb: 1.0\nc: 1e0\n",
			want: domain.Node{
				"a": 1,
				"b": 1.0,
				"c": 1e0,
			},
		},
		{
			name:  "error/invalid yaml",
			input: "a: [\n",
			errAssert: func(err error) bool {
				return errors.Is(err, ErrInvalidYAML)
			},
		},
		{
			name:  "error/top-level must be object for domain.Node",
			input: "- 1\n",
			errAssert: func(err error) bool {
				return errors.Is(err, ErrInvalidYAML)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseYAML([]byte(tt.input))

			if tt.errAssert != nil {
				require.Error(t, err)
				require.True(t, tt.errAssert(err), "unexpected error: %v", err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
