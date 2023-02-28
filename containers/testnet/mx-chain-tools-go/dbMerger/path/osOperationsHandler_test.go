package path

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDataDirectory = "./testdata"

func TestNewOperationsHandler(t *testing.T) {
	t.Parallel()

	handler := NewOsOperationsHandler()
	assert.False(t, check.IfNil(handler))
}

func TestOperationsHandler_CopyDirectory(t *testing.T) {
	t.Parallel()

	destDir := "copyDirectoryTest"
	workingDir := path.Join(testDataDirectory, destDir)
	cleanupDirectory(t, workingDir)
	defer cleanupDirectory(t, workingDir)

	handler := NewOsOperationsHandler()
	err := handler.CopyDirectory(workingDir, "./testdata/srcDir")
	assert.Nil(t, err)

	assert.Equal(t, readFileContent(t, "./testdata/srcDir/a/1"), readFileContent(t, path.Join(workingDir, "a/1")))
	assert.Equal(t, readFileContent(t, "./testdata/srcDir/a/file2.file"), readFileContent(t, path.Join(workingDir, "a/file2.file")))
	assert.Equal(t, readFileContent(t, "./testdata/srcDir/b/a.txt"), readFileContent(t, path.Join(workingDir, "b/a.txt")))
	assert.Equal(t, readFileContent(t, "./testdata/srcDir/c.log"), readFileContent(t, path.Join(workingDir, "c.log")))
}

func cleanupDirectory(t *testing.T, workingDir string) {
	err := os.RemoveAll(workingDir)
	assert.Nil(t, err)
}

func TestOsOperationsHandler_CheckIfDirectoryIsEmpty(t *testing.T) {
	t.Parallel()

	destDir := "checkIfDirectoryIsEmptyTest"
	workingDir := path.Join(testDataDirectory, destDir)
	cleanupDirectory(t, workingDir)
	defer cleanupDirectory(t, workingDir)

	handler := NewOsOperationsHandler()

	t.Run("directory does not exists should error", func(t *testing.T) {
		err := handler.CheckIfDirectoryIsEmpty(workingDir)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "no such file or directory while reading the directory"))
	})
	t.Run("directory exists and it's empty", func(t *testing.T) {
		err := os.MkdirAll(workingDir, os.ModePerm)
		assert.Nil(t, err)

		err = handler.CheckIfDirectoryIsEmpty(workingDir)
		assert.Nil(t, err)
	})
	t.Run("directory exists but it contains a directory", func(t *testing.T) {
		err := os.MkdirAll(path.Join(workingDir, "test"), os.ModePerm)
		assert.Nil(t, err)

		err = handler.CheckIfDirectoryIsEmpty(workingDir)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, errDirectoryIsNotEmpty))
	})
	t.Run("directory exists but it contains a file", func(t *testing.T) {
		err := os.RemoveAll(workingDir)
		assert.Nil(t, err)

		err = os.MkdirAll(path.Join(workingDir, "test"), os.ModePerm)
		assert.Nil(t, err)

		err = ioutil.WriteFile(path.Join(workingDir, "test.log"), []byte("data"), os.ModePerm)
		assert.Nil(t, err)

		err = handler.CheckIfDirectoryIsEmpty(workingDir)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, errDirectoryIsNotEmpty))
	})
}

func readFileContent(tb testing.TB, path string) string {
	in, err := os.Open(path)
	require.Nil(tb, err)

	defer func() {
		errClose := in.Close()
		require.Nil(tb, errClose)
	}()

	buff, err := io.ReadAll(in)
	require.Nil(tb, err)
	contents := string(buff)

	log.Info("read file", "path", path, "contents", contents)

	return contents
}
