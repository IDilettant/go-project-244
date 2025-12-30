package domain

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func assertSortedStrings(t *testing.T, keys []string) {
	t.Helper()

	require.True(t, sort.StringsAreSorted(keys))
}

func TestNode_KeysSorted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   Node
		want []string
	}{
		{name: "empty", in: Node{}, want: []string{}},
		{
			name: "sorted",
			in: Node{
				"b": 1,
				"a": 2,
				"z": 3,
			},
			want: []string{"a", "b", "z"},
		},
		{name: "single", in: Node{"a": 1}, want: []string{"a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.in.KeysSorted()
			require.Equal(t, tt.want, got)
			assertSortedStrings(t, got)
		})
	}
}

func TestNode_UnionKeysSorted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  Node
		right Node
		want  []string
	}{
		{name: "both empty", left: Node{}, right: Node{}, want: []string{}},
		{
			name: "disjoint",
			left: Node{"b": 1, "a": 2},
			right: Node{
				"c": 3,
				"d": 4,
			},
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "overlap",
			left: Node{
				"a": 1,
				"z": Node{"b": 1},
			},
			right: Node{
				"c": 3,
				"z": Node{"c": 3},
			},
			want: []string{"a", "c", "z"},
		},
		{
			name:  "same keys",
			left:  Node{"b": 1, "a": 2},
			right: Node{"a": 3, "b": 4},
			want:  []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.left.UnionKeysSorted(tt.right)
			require.Equal(t, tt.want, got)
			assertSortedStrings(t, got)
		})
	}
}

func TestNode_IsEmpty(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   Node
		want bool
	}{
		{name: "empty", in: Node{}, want: true},
		{name: "non-empty with nil value", in: Node{"a": nil}, want: false},
		{name: "non-empty", in: Node{"a": 1}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, tt.in.IsEmpty())
		})
	}
}

func TestAsNode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   any
		want Node
		ok   bool
	}{
		{name: "Node", in: Node{"a": 1}, want: Node{"a": 1}, ok: true},
		{name: "map[string]any", in: map[string]any{"a": 1}, want: Node{"a": 1}, ok: true},
		{name: "not a node", in: 123, want: nil, ok: false},
		{name: "nil", in: nil, want: nil, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := AsNode(tt.in)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNode_GetNode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		n    Node
		key  string
		want Node
		ok   bool
	}{
		{name: "missing key", n: Node{"a": 1}, key: "missing", want: nil, ok: false},
		{
			name: "value is Node",
			n:    Node{"child": Node{"x": 1}},
			key:  "child",
			want: Node{"x": 1},
			ok:   true,
		},
		{
			name: "value is map[string]any",
			n:    Node{"child": map[string]any{"x": 1}},
			key:  "child",
			want: Node{"x": 1},
			ok:   true,
		},
		{name: "value is not a node", n: Node{"child": "x"}, key: "child", want: nil, ok: false},
		{name: "empty node", n: Node{}, key: "child", want: nil, ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := tt.n.GetNode(tt.key)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.want, got)
		})
	}
}

