package main

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

const binPath string = "build/wp"

func TestHelpExecutes(t *testing.T) {
	cmd := exec.Command(binPath, "--help")

	_, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	cmd = exec.Command(binPath, "extract", "--help")

	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)
}

func TestExtractOneImage(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binPath, "extract", "128x128", tempDir, sourceImage)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	filenameSuffixes := []string{
		"scaled_center",
		"north",
		"northeast",
		"east",
		"southeast",
		"south",
		"southwest",
		"west",
		"northwest",
		"center",
	}

	expectedOutput := ""
	for _, str := range filenameSuffixes {
		expectedOutput += path.Join(tempDir, "128x128", "square_"+str) + ".jpg\n"
	}

	assert.Equal(t, expectedOutput, string(output))
}

func TestExtractOneImageDimensionError(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binPath, "extract", "1024x1024", tempDir, sourceImage)

	err = cmd.Start()
	assert.NoError(t, err)

	err = cmd.Wait()
	assert.NoError(t, err)
}

func TestExtractMultipleImages(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage1, _ := filepath.Abs(path.Join(cwd, "test_images", "tall.jpg"))
	sourceImage2, _ := filepath.Abs(path.Join(cwd, "test_images", "wide.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binPath, "extract", "128x128", tempDir, sourceImage1, sourceImage2)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	tallFilenameSuffixes := []string{
		"scaled_north",
		"scaled_center",
		"scaled_south",
	}

	wideFilenameSuffixes := []string{
		"scaled_west",
		"scaled_center",
		"scaled_east",
	}

	bothFilenameSuffixes := []string{
		"north",
		"northeast",
		"east",
		"southeast",
		"south",
		"southwest",
		"west",
		"northwest",
		"center",
	}

	expectedOutput := ""
	for _, str := range append(tallFilenameSuffixes, bothFilenameSuffixes...) {
		expectedOutput += path.Join(tempDir, "128x128", "tall_"+str) + ".jpg\n"
	}

	for _, str := range append(wideFilenameSuffixes, bothFilenameSuffixes...) {
		expectedOutput += path.Join(tempDir, "128x128", "wide_"+str) + ".jpg\n"
	}

	assert.Equal(t, expectedOutput, string(output))
}

// This test exists for historical purposes.
// There was once an issue where image extractions where the source image is
//   in the current working directory lead to the image being removed.
func TestExtractFromThisDirectory(t *testing.T) {
	cwd, _ := os.Getwd()

	originalImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))
	sourceImage, _ := filepath.Abs(path.Join(cwd, "square.jpg"))

	// Copy square.jpg from test images to the current working dir.
	source, err := os.Open(originalImage)
	assert.NoError(t, err)

	destination, err := os.Create(sourceImage)
	assert.NoError(t, err)

	_, err = io.Copy(destination, source)
	source.Close()
	destination.Close()

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binPath, "extract", "128x128", tempDir, sourceImage)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	filenameSuffixes := []string{
		"scaled_center",
		"north",
		"northeast",
		"east",
		"southeast",
		"south",
		"southwest",
		"west",
		"northwest",
		"center",
	}

	expectedOutput := ""
	for _, str := range filenameSuffixes {
		expectedOutput += path.Join(tempDir, "128x128", "square_"+str) + ".jpg\n"
	}

	assert.Equal(t, expectedOutput, string(output))

	_, err = os.Stat(sourceImage)
	assert.NoError(t, err)

	os.Remove(sourceImage)
}

func TestPickImage(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binPath, "pick", "128x128", tempDir, "north", sourceImage)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	filenameSuffixes := []string{
		"north",
	}

	expectedOutput := ""
	for _, str := range filenameSuffixes {
		expectedOutput += path.Join(tempDir, "128x128", "square_"+str) + ".jpg\n"
	}

	assert.Equal(t, expectedOutput, string(output))
}

func TestPickMultipleImages(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage1, _ := filepath.Abs(path.Join(cwd, "test_images", "tall.jpg"))
	sourceImage2, _ := filepath.Abs(path.Join(cwd, "test_images", "wide.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binPath, "pick", "128x128", tempDir, "north", sourceImage1, sourceImage2)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	tallFilenameSuffixes := []string{
		"north",
	}

	wideFilenameSuffixes := []string{
		"north",
	}

	expectedOutput := ""
	for _, str := range tallFilenameSuffixes {
		expectedOutput += path.Join(tempDir, "128x128", "tall_"+str) + ".jpg\n"
	}
	for _, str := range wideFilenameSuffixes {
		expectedOutput += path.Join(tempDir, "128x128", "wide_"+str) + ".jpg\n"
	}

	assert.Equal(t, expectedOutput, string(output))
}

func TestPickImageScaled(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binPath, "pick", "128x128", tempDir, "north", "--scaled", sourceImage)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	filenameSuffixes := []string{
		"scaled_north",
	}

	expectedOutput := ""
	for _, str := range filenameSuffixes {
		expectedOutput += path.Join(tempDir, "128x128", "square_"+str) + ".jpg\n"
	}

	assert.Equal(t, expectedOutput, string(output))
}

// This test exists for historical purposes.
// There was once an issue where image extractions where the source image is
//   in the current working directory lead to the image being removed.
func TestPickFromThisDirectory(t *testing.T) {
	cwd, _ := os.Getwd()

	originalImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))
	sourceImage, _ := filepath.Abs(path.Join(cwd, "square.jpg"))

	// Copy square.jpg from test images to the current working dir.
	source, err := os.Open(originalImage)
	assert.NoError(t, err)

	destination, err := os.Create(sourceImage)
	assert.NoError(t, err)

	_, err = io.Copy(destination, source)
	source.Close()
	destination.Close()

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binPath, "pick", "128x128", tempDir, "north", sourceImage)

	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)

	filenameSuffixes := []string{
		"north",
	}

	expectedOutput := ""
	for _, str := range filenameSuffixes {
		expectedOutput += path.Join(tempDir, "128x128", "square_"+str) + ".jpg\n"
	}

	assert.Equal(t, expectedOutput, string(output))

	_, err = os.Stat(sourceImage)
	assert.NoError(t, err)

	os.Remove(sourceImage)
}
