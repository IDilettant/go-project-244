package diff

import (
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/parser"
)

func TestCompareSemantics(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name  string
		left  parser.Node
		right parser.Node
		want  []Change
	}

	tests := []testCase{
		{
			name:  "empty/both empty returns empty result",
			left:  parser.Node{},
			right: parser.Node{},
			want:  []Change{},
		},
		{
			name:  "unchanged/identical keys have unchanged type",
			left:  parser.Node{"a": float64(1), "b": "x"},
			right: parser.Node{"a": float64(1), "b": "x"},
			want: []Change{
				{Key: "a", Type: Unchanged, OldValue: float64(1), NewValue: float64(1)},
				{Key: "b", Type: Unchanged, OldValue: "x", NewValue: "x"},
			},
		},
		{
			name:  "comparison/numeric value rs by fractional part marks as updated",
			left:  parser.Node{"a": float64(1)},
			right: parser.Node{"a": float64(1.5)},
			want: []Change{
				{Key: "a", Type: Updated, OldValue: float64(1), NewValue: float64(1.5)},
			},
		},
		{
			name:  "removed/key exists only on the left",
			left:  parser.Node{"a": float64(1)},
			right: parser.Node{},
			want: []Change{
				{Key: "a", Type: Removed, OldValue: float64(1)},
			},
		},
		{
			name:  "added/key exists only on the right",
			left:  parser.Node{},
			right: parser.Node{"a": float64(1)},
			want: []Change{
				{Key: "a", Type: Added, NewValue: float64(1)},
			},
		},
		{
			name:  "updated/same key rent value",
			left:  parser.Node{"a": float64(1)},
			right: parser.Node{"a": float64(2)},
			want: []Change{
				{Key: "a", Type: Updated, OldValue: float64(1), NewValue: float64(2)},
			},
		},
		{
			name:  "mixed/added removed unchanged together",
			left:  parser.Node{"a": float64(1), "b": float64(2)},
			right: parser.Node{"a": float64(1), "c": float64(3)},
			want: []Change{
				{Key: "a", Type: Unchanged, OldValue: float64(1), NewValue: float64(1)},
				{Key: "b", Type: Removed, OldValue: float64(2)},
				{Key: "c", Type: Added, NewValue: float64(3)},
			},
		},
		{
			name:  "types/string and number are rent",
			left:  parser.Node{"a": "1"},
			right: parser.Node{"a": float64(1)},
			want: []Change{
				{Key: "a", Type: Updated, OldValue: "1", NewValue: float64(1)},
			},
		},
		{
			name:  "null/both null values are unchanged",
			left:  parser.Node{"a": nil},
			right: parser.Node{"a": nil},
			want: []Change{
				{Key: "a", Type: Unchanged, OldValue: nil, NewValue: nil},
			},
		},
		{
			name:  "null/missing key is removed",
			left:  parser.Node{"a": nil},
			right: parser.Node{},
			want: []Change{
				{Key: "a", Type: Removed, OldValue: nil},
			},
		},
		{
			name:  "null/missing key is added",
			left:  parser.Node{},
			right: parser.Node{"a": nil},
			want: []Change{
				{Key: "a", Type: Added, NewValue: nil},
			},
		},
		{
			name:  "keys/key comparison is case-sensitive",
			left:  parser.Node{"A": float64(1)},
			right: parser.Node{"a": float64(1)},
			want: []Change{
				{Key: "A", Type: Removed, OldValue: float64(1)},
				{Key: "a", Type: Added, NewValue: float64(1)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Compare(tt.left, tt.right)

			require.Len(t, got, len(tt.want))
			for i := range tt.want {
				assertChange(t, got[i], tt.want[i])
			}
		})
	}
}

func TestCompareOrderingSortsKeysAscending(t *testing.T) {
	t.Parallel()

	left := parser.Node{
		"b":  float64(2),
		"a":  float64(1),
		"c":  float64(3),
		"уж": float64(4),
	}
	right := parser.Node{
		"b":  float64(2),
		"d":  float64(4),
		"уж": float64(5),
	}

	got := Compare(left, right)
	require.NotEmpty(t, got)

	for i := 1; i < len(got); i++ {
		require.LessOrEqual(t, got[i-1].Key, got[i].Key, "changes must be sorted by key")
	}
}

func assertChange(t *testing.T, got Change, want Change) {
	t.Helper()

	require.Equal(t, want.Key, got.Key, "key mismatch")
	require.Equal(t, want.Type, got.Type, "type mismatch for key %q", want.Key)

	switch want.Type {
	case Added:
		require.Equal(t, want.NewValue, got.NewValue)

	case Removed:
		require.Equal(t, want.OldValue, got.OldValue)

	case Updated:
		require.Equal(t, want.OldValue, got.OldValue)
		require.Equal(t, want.NewValue, got.NewValue)

	case Unchanged:
		require.Equal(t, want.OldValue, got.OldValue)

	default:
		t.Fatalf("unknown ChangeType: %v", want.Type)
	}
}
