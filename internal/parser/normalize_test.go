package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeNode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     Node
		want      Node
		errAssert func(error) bool
	}{
		{
			name: "ok/primitives are preserved",
			input: Node{
				"s": "x",
				"b": true,
				"n": nil,
				"f": float64(1.5),
			},
			want: Node{
				"s": "x",
				"b": true,
				"n": nil,
				"f": float64(1.5),
			},
		},
		{
			name: "ok/ints are converted to float64",
			input: Node{
				"i": int(1),
				"j": int64(2),
			},
			want: Node{
				"i": float64(1),
				"j": float64(2),
			},
		},
		{
			name: "ok/nested structures are normalized",
			input: Node{
				"outer": Node{
					"inner": Node{
						"i": int(10),
						"f": float64(2.25),
					},
					"arr": []any{
						int(1),
						float64(2.5),
						Node{"x": int64(3)},
						[]any{int(4), float64(5.75)},
					},
				},
			},
			want: Node{
				"outer": Node{
					"inner": Node{
						"i": float64(10),
						"f": float64(2.25),
					},
					"arr": []any{
						float64(1),
						float64(2.5),
						Node{"x": float64(3)},
						[]any{float64(4), float64(5.75)},
					},
				},
			},
		},
		{
			name: "ok/yaml map any any is converted to node",
			input: Node{
				"outer": map[any]any{
					"inner": map[any]any{
						"i": int(1),
					},
				},
			},
			want: Node{
				"outer": Node{
					"inner": Node{
						"i": float64(1),
					},
				},
			},
		},
		{
			name: "error/yaml map with non string key",
			input: Node{
				"bad": map[any]any{
					1: "x",
				},
			},
			errAssert: func(err error) bool {
				return err != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeNode(tt.input)

			if tt.errAssert != nil {
				require.Error(t, err)
				require.True(t, tt.errAssert(err))

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
