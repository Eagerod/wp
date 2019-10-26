package main

import (
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

func TestHelpExecutes(t *testing.T) {
	_, err := exec.LookPath("wp")
	assert.NoError(t, err)

	cmd := exec.Command("wp", "--help")

	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)

	cmd = exec.Command("wp", "extract", "--help")

	_, err = cmd.CombinedOutput()
	assert.NoError(t, err)
}

func TestExtractOneImage(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("wp", "extract", "128x128", tempDir, sourceImage)

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

	cmd := exec.Command("wp", "extract", "1024x1024", tempDir, sourceImage)

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

	cmd := exec.Command("wp", "extract", "128x128", tempDir, sourceImage1, sourceImage2)

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

func TestPickImage(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("wp", "pick", "128x128", tempDir, "north", sourceImage)

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

func TestPickImageScaled(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, _ := filepath.Abs(path.Join(cwd, "test_images", "square.jpg"))

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("wp", "pick", "128x128", tempDir, "north", "--scaled", sourceImage)

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
