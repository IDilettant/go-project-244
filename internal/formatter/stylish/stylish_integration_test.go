package stylish_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	"code/internal/formatter/stylish"
	"code/internal/parser"
)

const (
	testdataDir = "testdata"
	fixtureDir  = "fixture"
	flatFixture = "flat"

	leftFileName  = "filepath1.json"
	rightFileName = "filepath2.json"
	expectedFile  = "expected_stylish.txt"
)

func TestStylishFormatterIntegrationFlatFixture(t *testing.T) {
	t.Parallel()

	td := flatFixtureDir(t)

	leftPath := filepath.Join(td, leftFileName)
	rightPath := filepath.Join(td, rightFileName)
	wantPath := filepath.Join(td, expectedFile)

	left, err := parser.ParseFile(leftPath)
	require.NoError(t, err)

	right, err := parser.ParseFile(rightPath)
	require.NoError(t, err)

	changes := diff.Compare(left, right)

	f := stylish.New()
	got := f.Format(changes)

	wantBytes, err := os.ReadFile(wantPath)
	require.NoError(t, err)

	require.Equal(t, string(wantBytes), got)
}

func flatFixtureDir(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")

	stylishDir := filepath.Dir(thisFile)
	formatterDir := filepath.Dir(stylishDir)

	return filepath.Join(
		formatterDir,
		testdataDir,
		fixtureDir,
		flatFixture,
	)
}
