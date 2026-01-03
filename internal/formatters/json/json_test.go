package json_test

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"code/internal/diff"
	jsonfmt "code/internal/formatters/json"
	"code/internal/parser"
)

const (
	testdataDir = "testdata"
	fixtureDir  = "fixture"

	leftJSONFileName  = "filepath1.json"
	rightJSONFileName = "filepath2.json"
)

func TestJSONFormatter_Smoke(t *testing.T) {
	t.Parallel()

	td := getFixtureDirPath(t)

	leftPath := filepath.Join(td, leftJSONFileName)
	rightPath := filepath.Join(td, rightJSONFileName)

	left, err := parser.ParseFile(leftPath)
	require.NoError(t, err)

	right, err := parser.ParseFile(rightPath)
	require.NoError(t, err)

	changes := diff.Compare(left, right)

	f := jsonfmt.New()
	output, err := f.Format(changes)
	require.NoError(t, err)
	require.NotEmpty(t, output)

	var decoded map[string]any
	err = json.Unmarshal([]byte(output), &decoded)
	require.NoError(t, err)

	require.Contains(t, decoded, "common")
	require.Contains(t, decoded, "group1")
}

func getFixtureDirPath(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)

	jsonDir := filepath.Dir(thisFile)
	formatterDir := filepath.Dir(jsonDir)

	return filepath.Join(
		formatterDir,
		testdataDir,
		fixtureDir,
	)
}
