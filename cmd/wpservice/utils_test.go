package wpservice

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestParseDimensionsStringValid(t *testing.T) {
	point, err := ParseDimensionsString("1024x768")

	assert.NoError(t, err)
	assert.Equal(t, point.X, 1024)
	assert.Equal(t, point.Y, 768)
}

func TestParseDimensionsStringInvalid(t *testing.T) {
	point, err := ParseDimensionsString("3D")

	assert.Equal(t, point, image.ZP)
	assert.Equal(t, err.Error(), "Provided dimension string (3D) is not valid")
}

func TestParseDimensionsStringInvalidWidth(t *testing.T) {
	point, err := ParseDimensionsString("0x768")

	assert.Equal(t, point, image.ZP)
	assert.Equal(t, err.Error(), "Provided width is not a valid positive integer")
}

func TestParseDimensionsStringInvalidHeight(t *testing.T) {
	point, err := ParseDimensionsString("1024x0")

	assert.Equal(t, point, image.ZP)
	assert.Equal(t, err.Error(), "Provided height is not a valid positive integer")
}

func TestExtractGravitiesFromSourceImageScaled(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	doImageMagick = func(args ...string) (string, error) {
		assert.Equal(t, []string{"abc", "-gravity", "Center", "-scale", "64x64^", "-extent", "64x64", "images/abc_scaled_center"}, args)
		return "", nil
	}

	err := ExtractGravitiesFromSourceImage("abc", true, []string{"Center"}, "64x64", "images")
	assert.NoError(t, err)
}

func TestExtractGravitiesFromSourceImageUnscaled(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	doImageMagick = func(args ...string) (string, error) {
		assert.Equal(t, []string{"abc", "-gravity", "Center", "-extent", "64x64", "images/abc_center"}, args)
		return "", nil
	}

	err := ExtractGravitiesFromSourceImage("abc", false, []string{"Center"}, "64x64", "images")
	assert.NoError(t, err)
}

func TestExtractFromLocalImageSameAspectRatio(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "square.jpg"))
	assert.NoError(t, err)

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	outputDir, err := filepath.Abs(path.Join(tempDir, "64x64"))
	assert.NoError(t, err)

	expectedCalls := [][]string{
		[]string{sourceImage, "-gravity", "Center", "-scale", "64x64^", "-extent", "64x64", path.Join(outputDir, "square_scaled_center.jpg")},
		[]string{sourceImage, "-gravity", "North", "-extent", "64x64", path.Join(outputDir, "square_north.jpg")},
		[]string{sourceImage, "-gravity", "NorthEast", "-extent", "64x64", path.Join(outputDir, "square_northeast.jpg")},
		[]string{sourceImage, "-gravity", "East", "-extent", "64x64", path.Join(outputDir, "square_east.jpg")},
		[]string{sourceImage, "-gravity", "SouthEast", "-extent", "64x64", path.Join(outputDir, "square_southeast.jpg")},
		[]string{sourceImage, "-gravity", "South", "-extent", "64x64", path.Join(outputDir, "square_south.jpg")},
		[]string{sourceImage, "-gravity", "SouthWest", "-extent", "64x64", path.Join(outputDir, "square_southwest.jpg")},
		[]string{sourceImage, "-gravity", "West", "-extent", "64x64", path.Join(outputDir, "square_west.jpg")},
		[]string{sourceImage, "-gravity", "NorthWest", "-extent", "64x64", path.Join(outputDir, "square_northwest.jpg")},
		[]string{sourceImage, "-gravity", "Center", "-extent", "64x64", path.Join(outputDir, "square_center.jpg")},
	}

	doImageMagick = func(args ...string) (string, error) {
		// Fail now, rather than assert. Assert will continue, and crash at [0]
		if len(expectedCalls) == 0 {
			fmt.Fprintln(os.Stderr, "doImageMagick called when not expected")
			t.FailNow()
		}

		assert.Equal(t, expectedCalls[0], args)
		expectedCalls = expectedCalls[1:]
		return "", nil
	}

	err = ExtractFromLocalImage("64x64", tempDir, sourceImage)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(expectedCalls))
}

func TestExtractFromLocalImageWideAspectRatio(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "wide.jpg"))
	assert.NoError(t, err)

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	outputDir, _ := filepath.Abs(path.Join(tempDir, "64x64"))

	expectedCalls := [][]string{
		[]string{sourceImage, "-gravity", "West", "-scale", "64x64^", "-extent", "64x64", path.Join(outputDir, "wide_scaled_west.jpg")},
		[]string{sourceImage, "-gravity", "Center", "-scale", "64x64^", "-extent", "64x64", path.Join(outputDir, "wide_scaled_center.jpg")},
		[]string{sourceImage, "-gravity", "East", "-scale", "64x64^", "-extent", "64x64", path.Join(outputDir, "wide_scaled_east.jpg")},
		[]string{sourceImage, "-gravity", "North", "-extent", "64x64", path.Join(outputDir, "wide_north.jpg")},
		[]string{sourceImage, "-gravity", "NorthEast", "-extent", "64x64", path.Join(outputDir, "wide_northeast.jpg")},
		[]string{sourceImage, "-gravity", "East", "-extent", "64x64", path.Join(outputDir, "wide_east.jpg")},
		[]string{sourceImage, "-gravity", "SouthEast", "-extent", "64x64", path.Join(outputDir, "wide_southeast.jpg")},
		[]string{sourceImage, "-gravity", "South", "-extent", "64x64", path.Join(outputDir, "wide_south.jpg")},
		[]string{sourceImage, "-gravity", "SouthWest", "-extent", "64x64", path.Join(outputDir, "wide_southwest.jpg")},
		[]string{sourceImage, "-gravity", "West", "-extent", "64x64", path.Join(outputDir, "wide_west.jpg")},
		[]string{sourceImage, "-gravity", "NorthWest", "-extent", "64x64", path.Join(outputDir, "wide_northwest.jpg")},
		[]string{sourceImage, "-gravity", "Center", "-extent", "64x64", path.Join(outputDir, "wide_center.jpg")},
	}

	doImageMagick = func(args ...string) (string, error) {
		// Fail now, rather than assert. Assert will continue, and crash at [0]
		if len(expectedCalls) == 0 {
			fmt.Fprintln(os.Stderr, "doImageMagick called when not expected")
			t.FailNow()
		}

		assert.Equal(t, expectedCalls[0], args)
		expectedCalls = expectedCalls[1:]
		return "", nil
	}

	err = ExtractFromLocalImage("64x64", tempDir, sourceImage)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(expectedCalls))
}

func TestExtractFromLocalImageTallAspectRatio(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "tall.jpg"))
	assert.NoError(t, err)

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	outputDir, err := filepath.Abs(path.Join(tempDir, "64x64"))
	assert.NoError(t, err)

	expectedCalls := [][]string{
		[]string{sourceImage, "-gravity", "North", "-scale", "64x64^", "-extent", "64x64", path.Join(outputDir, "tall_scaled_north.jpg")},
		[]string{sourceImage, "-gravity", "Center", "-scale", "64x64^", "-extent", "64x64", path.Join(outputDir, "tall_scaled_center.jpg")},
		[]string{sourceImage, "-gravity", "South", "-scale", "64x64^", "-extent", "64x64", path.Join(outputDir, "tall_scaled_south.jpg")},
		[]string{sourceImage, "-gravity", "North", "-extent", "64x64", path.Join(outputDir, "tall_north.jpg")},
		[]string{sourceImage, "-gravity", "NorthEast", "-extent", "64x64", path.Join(outputDir, "tall_northeast.jpg")},
		[]string{sourceImage, "-gravity", "East", "-extent", "64x64", path.Join(outputDir, "tall_east.jpg")},
		[]string{sourceImage, "-gravity", "SouthEast", "-extent", "64x64", path.Join(outputDir, "tall_southeast.jpg")},
		[]string{sourceImage, "-gravity", "South", "-extent", "64x64", path.Join(outputDir, "tall_south.jpg")},
		[]string{sourceImage, "-gravity", "SouthWest", "-extent", "64x64", path.Join(outputDir, "tall_southwest.jpg")},
		[]string{sourceImage, "-gravity", "West", "-extent", "64x64", path.Join(outputDir, "tall_west.jpg")},
		[]string{sourceImage, "-gravity", "NorthWest", "-extent", "64x64", path.Join(outputDir, "tall_northwest.jpg")},
		[]string{sourceImage, "-gravity", "Center", "-extent", "64x64", path.Join(outputDir, "tall_center.jpg")},
	}

	doImageMagick = func(args ...string) (string, error) {
		// Fail now, rather than assert. Assert will continue, and crash at [0]
		if len(expectedCalls) == 0 {
			fmt.Fprintln(os.Stderr, "doImageMagick called when not expected")
			t.FailNow()
		}

		assert.Equal(t, expectedCalls[0], args)
		expectedCalls = expectedCalls[1:]
		return "", nil
	}

	err = ExtractFromLocalImage("64x64", tempDir, sourceImage)
	assert.NoError(t, err)

	assert.Equal(t, 0, len(expectedCalls))
}

func TestExtractFromImageLocal(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "tall.jpg"))
	assert.NoError(t, err)

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	doImageMagick = func(args ...string) (string, error) {
		return "", nil
	}

	err = ExtractFromImage("64x64", tempDir, sourceImage)
	assert.NoError(t, err)
}

func TestExtractFromImageUsingFileProtocol(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "tall.jpg"))
	assert.NoError(t, err)

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	doImageMagick = func(args ...string) (string, error) {
		return "", nil
	}

	err = ExtractFromImage("64x64", tempDir, "file://"+sourceImage)
	assert.NoError(t, err)
}

func TestExtractFromImageUsingRemoteFile(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	g := downloadFile
	defer func() {
		downloadFile = g
	}()

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "tall.jpg"))
	assert.NoError(t, err)

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	doImageMagick = func(args ...string) (string, error) {
		return "", nil
	}

	downloadFile = func(dest, url string) error {
		input, err := ioutil.ReadFile(sourceImage)
		assert.NoError(t, err)

		err = ioutil.WriteFile(dest, input, 0644)
		assert.NoError(t, err)
		return nil
	}

	// Download file is mocked, so url can be invalid
	err = ExtractFromImage("64x64", tempDir, "http://"+sourceImage)
	assert.NoError(t, err)
}
