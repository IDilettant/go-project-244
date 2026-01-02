package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/formatters/common"
)

func TestRunReturnsDiffOutput(t *testing.T) {
	t.Parallel()

	left := writeTempJSONFile(t, "left.json", `{"host":"hexlet.io","timeout":50}`)
	right := writeTempJSONFile(t, "right.json", `{"host":"hexlet.io","timeout":20}`)
	missing := filepath.Join(t.TempDir(), "no_such_file.json")

	type testCase struct {
		name      string
		run       func() (string, error)
		wantErrIs error
		wantCode  int
		wantOut   bool
	}

	tests := []testCase{
		{
			name: "ok/run returns diff output",
			run: func() (string, error) {
				return Run(left, right, common.FormatStylish)
			},
			wantCode: 0,
			wantOut:  true,
		},
		{
			name: "error/unknown format is usage error",
			run: func() (string, error) {
				return Run(left, right, "unknown")
			},
			wantErrIs: ErrUsage,
			wantCode:  exitCodeUsageError,
		},
		{
			name: "error/missing file is runtime error",
			run: func() (string, error) {
				return Run(missing, right, common.FormatStylish)
			},
			wantErrIs: ErrRuntime,
			wantCode:  exitCodeRuntimeError,
		},
		{
			name: "error/invalid args is usage error",
			run: func() (string, error) {
				return "", invalidArgsError()
			},
			wantErrIs: ErrUsage,
			wantCode:  exitCodeUsageError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			out, err := tt.run()

			if tt.wantErrIs == nil {
				require.NoError(t, err)
				require.NotEmpty(t, out)

				return
			}

			require.Error(t, err)
			require.ErrorIs(t, err, tt.wantErrIs)
			require.Equal(t, tt.wantCode, exitCodeFrom(err))
			require.Empty(t, out)
		})
	}
}

func TestNew_RunSuccess_WritesOutput(t *testing.T) {
	left := writeTempJSONFile(t, "left.json", `{"host":"hexlet.io","timeout":50}`)
	right := writeTempJSONFile(t, "right.json", `{"host":"hexlet.io","timeout":20}`)

	cmd := New()
	buf := &bytes.Buffer{}
	cmd.Writer = buf

	err := cmd.Run(context.Background(), []string{
		"gendiff",
		left,
		right,
	})

	require.NoError(t, err)
	require.NotEmpty(t, buf.String())
}

func writeTempJSONFile(t *testing.T, name, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, name)

	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	return path
}
