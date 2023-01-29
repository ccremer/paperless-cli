package archive

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnzip(t *testing.T) {
	testFilePath := "testdata/unzip.zip"
	testDir := "testdata/run"

	// cleanup previous test files in case of failure
	require.NoError(t, os.RemoveAll(testDir))

	err := Unzip(context.TODO(), testFilePath, testDir)
	assert.NoError(t, err, "unzip failed with error")

	assert.FileExists(t, filepath.Join(testDir, "toplevel.file"))
	assert.FileExists(t, filepath.Join(testDir, "Dir In Archive", "Sub Dir.file"))

	// cleanup
	require.NoError(t, os.RemoveAll(testDir))
}
