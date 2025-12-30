package parser

import (
	"errors"
	"testing"

	"code/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestParseJSON(t *testing.T) {
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
			input: `{}`,
			want:  domain.Node{},
		},
		{
			name:  "types/primitives are parsed",
			input: `{"s":"x","b":true,"n":null}`,
			want: domain.Node{
				"s": "x",
				"b": true,
				"n": nil,
			},
		},
		{
			name:  "numbers/integer is decoded as float64",
			input: `{"a":1}`,
			want: domain.Node{
				"a": float64(1),
			},
		},
		{
			name:  "numbers/fractional is decoded as float64 with fraction",
			input: `{"a":1.5}`,
			want: domain.Node{
				"a": float64(1.5),
			},
		},
		{
			name:  "numbers/equal magnitude numbers have equal decoded representation",
			input: `{"a":1,"b":1.0,"c":1e0}`,
			want: domain.Node{
				"a": float64(1),
				"b": float64(1),
				"c": float64(1),
			},
		},
		{
			name:  "error/invalid json",
			input: `{"a":`,
			errAssert: func(err error) bool {
				return errors.Is(err, ErrInvalidJSON)
			},
		},
		{
			name:  "error/top-level must be object for domain.Node",
			input: `[]`,
			errAssert: func(err error) bool {
				return errors.Is(err, ErrInvalidJSON)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseJSON([]byte(tt.input))

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
