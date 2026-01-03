package stylish_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatters/stylish"
	"code/internal/parser"
)

const (
	testdataDir = "testdata"
	fixtureDir  = "fixture"

	leftJSONFileName  = "filepath1.json"
	rightJSONFileName = "filepath2.json"

	leftYMLFileName  = "filepath1.yml"
	rightYMLFileName = "filepath2.yml"

	expectedFile = "expected_stylish.txt"
)

func TestStylishFormatterIntegration(t *testing.T) {
	t.Parallel()

	td := getFixtureDirPath(t)
	wantPath := filepath.Join(td, expectedFile)

	wantBytes, err := os.ReadFile(wantPath)
	require.NoError(t, err)
	want := string(wantBytes)

	tests := []struct {
		name      string
		leftFile  string
		rightFile string
	}{
		{
			name:      "ok/json files",
			leftFile:  leftJSONFileName,
			rightFile: rightJSONFileName,
		},
		{
			name:      "ok/yml files",
			leftFile:  leftYMLFileName,
			rightFile: rightYMLFileName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			leftPath := filepath.Join(td, tt.leftFile)
			rightPath := filepath.Join(td, tt.rightFile)

			left, err := parser.ParseFile(leftPath)
			require.NoError(t, err)

			right, err := parser.ParseFile(rightPath)
			require.NoError(t, err)

			changes := diff.Compare(left, right)

			f := stylish.New()
			got, err := f.Format(changes)
			require.NoError(t, err)

			require.Equal(t, want, got)
		})
	}
}

func getFixtureDirPath(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")

	stylishDir := filepath.Dir(thisFile)
	formatterDir := filepath.Dir(stylishDir)

	return filepath.Join(
		formatterDir,
		testdataDir,
		fixtureDir,
	)
}
