package wp

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestPrepareImageFromSourceLocal(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "square.jpg"))
	assert.NoError(t, err)

	is, err := PrepareImageFromSource(sourceImage, "")
	assert.NoError(t, err)

	assert.Equal(t, "square.jpg", path.Base(is.LocalPath))
	CleanupImageSource(is)
}

func TestPrepareImageFromSourceLocalCached(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cwd, _ := os.Getwd()
	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "square.jpg"))
	assert.NoError(t, err)

	is, err := PrepareImageFromSource(sourceImage, tempDir)
	assert.NoError(t, err)

	assert.Equal(t, "square.jpg", path.Base(is.LocalPath))
	CleanupImageSource(is)

	_, err = os.Stat(is.LocalPath)
	assert.NoError(t, err)
}

func TestPrepareImageFromSourceLocalWithProtocol(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "square.jpg"))
	assert.NoError(t, err)

	is, err := PrepareImageFromSource("file://"+sourceImage, "")
	assert.NoError(t, err)

	assert.Equal(t, "square.jpg", path.Base(is.LocalPath))
	CleanupImageSource(is)
}

func TestPrepareImageFromSourceRemote(t *testing.T) {
	g := downloadFile
	defer func() {
		downloadFile = g
	}()

	cwd, _ := os.Getwd()
	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "square.jpg"))
	assert.NoError(t, err)

	downloadFile = func(dest, url string) error {
		input, err := ioutil.ReadFile(sourceImage)
		assert.NoError(t, err)

		err = ioutil.WriteFile(dest, input, 0644)
		assert.NoError(t, err)
		return nil
	}

	is, err := PrepareImageFromSource("http://"+sourceImage, "")
	assert.NoError(t, err)

	assert.Equal(t, "square.jpg", path.Base(is.LocalPath))
	CleanupImageSource(is)
}

func TestCleanupImageSource(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "square.jpg"))
	assert.NoError(t, err)

	is, err := PrepareImageFromSource(sourceImage, "")
	assert.NoError(t, err)

	assert.Equal(t, "square.jpg", path.Base(is.LocalPath))
	CleanupImageSource(is)

	_, err = os.Stat(filepath.Dir(is.LocalPath))

	e, ok := err.(*os.PathError)
	assert.True(t, ok)
	assert.NotNil(t, e)
	assert.True(t, os.IsNotExist(e))
}
