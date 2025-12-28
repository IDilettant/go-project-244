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

func TestStylishFormatter_integration_flat_fixture(t *testing.T) {
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

	// current file: internal/formatter/stylish/stylish_integration_test.go
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")

	// .../internal/formatter/stylish
	stylishDir := filepath.Dir(thisFile)

	// .../internal/formatter
	formatterDir := filepath.Dir(stylishDir)

	// .../internal/formatter/testdata/fixture/flat
	return filepath.Join(
		formatterDir,
		testdataDir,
		fixtureDir,
		flatFixture,
	)
}

