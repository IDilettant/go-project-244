package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/domain"
)

func TestParseFile(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		path      string
		setup     func(t *testing.T) string
		want      domain.Node
		errAssert func(t *testing.T, err error)
	}

	tests := []testCase{
		{
			name: "ok/json file parsed",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				p := filepath.Join(dir, "file.json")
				require.NoError(t, os.WriteFile(p, []byte(`{"a":1,"b":"x"}`), 0o644))

				return p
			},
			want: domain.Node{
				"a": float64(1),
				"b": "x",
			},
		},
		{
			name: "ok/yaml file parsed and numbers normalized",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				p := filepath.Join(dir, "file.yaml")
				require.NoError(t, os.WriteFile(p, []byte("a: 1\nb: x\n"), 0o644))

				return p
			},
			want: domain.Node{
				"a": float64(1),
				"b": "x",
			},
		},
		{
			name: "error/unsupported extension",
			setup: func(t *testing.T) string {
				t.Helper()

				dir := t.TempDir()
				p := filepath.Join(dir, "file.txt")
				require.NoError(t, os.WriteFile(p, []byte(`{"a":1}`), 0o644))

				return p
			},
			errAssert: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrUnsupportedFormat)
			},
		},
		{
			name: "error/file does not exist",
			path: filepath.Join("no", "such", "file.json"),
			errAssert: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				require.Contains(t, err.Error(), "read file")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			path := tt.path
			if tt.setup != nil {
				path = tt.setup(t)
			}

			got, err := ParseFile(path)

			if tt.errAssert != nil {
				require.Error(t, err)
				tt.errAssert(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
