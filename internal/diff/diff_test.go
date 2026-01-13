package diff

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/domain"
)

func TestCompareSemantics(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  domain.Node
		right domain.Node
		want  []Change
	}{
		{
			name:  "empty/both empty returns empty result",
			left:  domain.Node{},
			right: domain.Node{},
			want:  []Change{},
		},
		{
			name: "nested/all change types inside container",
			left: domain.Node{
				"common": domain.Node{
					"unchanged": float64(1),
					"removed":   float64(2),
					"updated":   float64(3),
					"keepNull":  nil,
					"typeLeft":  "1",
				},
			},
			right: domain.Node{
				"common": domain.Node{
					"unchanged": float64(1),
					"added":     float64(10),
					"updated":   float64(3.5),
					"keepNull":  nil,
					"typeLeft":  float64(1),
				},
			},
			want: []Change{
				{
					Key:  "common",
					Type: Nested,
					Children: []Change{
						{Key: "added", Type: Added, OldValue: nil, NewValue: float64(10)},
						{Key: "keepNull", Type: Unchanged, OldValue: nil, NewValue: nil},
						{Key: "removed", Type: Removed, OldValue: float64(2), NewValue: nil},
						{Key: "typeLeft", Type: Updated, OldValue: "1", NewValue: float64(1)},
						{
							Key:      "unchanged",
							Type:     Unchanged,
							OldValue: float64(1),
							NewValue: float64(1),
						},
						{
							Key:      "updated",
							Type:     Updated,
							OldValue: float64(3),
							NewValue: float64(3.5),
						},
					},
				},
			},
		},
		{
			name: "nested/container vs primitive becomes updated",
			left: domain.Node{
				"common": domain.Node{
					"nodeToValue": domain.Node{"k": "v"},
					"valueToNode": "x",
				},
			},
			right: domain.Node{
				"common": domain.Node{
					"nodeToValue": "str",
					"valueToNode": domain.Node{"k": "v"},
				},
			},
			want: []Change{
				{
					Key:  "common",
					Type: Nested,
					Children: []Change{
						{
							Key:      "nodeToValue",
							Type:     Updated,
							OldValue: domain.Node{"k": "v"},
							NewValue: "str",
						},
						{
							Key:      "valueToNode",
							Type:     Updated,
							OldValue: "x",
							NewValue: domain.Node{"k": "v"},
						},
					},
				},
			},
		},
		{
			name: "nested/containers build children recursively",
			left: domain.Node{
				"root": domain.Node{
					"child": domain.Node{
						"a": float64(1),
					},
				},
			},
			right: domain.Node{
				"root": domain.Node{
					"child": domain.Node{
						"a": float64(2),
						"b": float64(3),
					},
				},
			},
			want: []Change{
				{
					Key:  "root",
					Type: Nested,
					Children: []Change{
						{
							Key:  "child",
							Type: Nested,
							Children: []Change{
								{
									Key:      "a",
									Type:     Updated,
									OldValue: float64(1),
									NewValue: float64(2),
								},
								{Key: "b", Type: Added, OldValue: nil, NewValue: float64(3)},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Compare(tt.left, tt.right)
			require.Equal(t, tt.want, got)
		})
	}
}
