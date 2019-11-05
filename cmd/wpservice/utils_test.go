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

func TestGetImageDimensions(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "square.jpg"))
	assert.NoError(t, err)

	dims, err := GetImageDimensions(sourceImage)
	assert.NoError(t, err)

	assert.Equal(t, dims, image.Point{128, 128})
}

func TestGetImageDimensionsNotFound(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images", "not-an-image.jpg"))
	assert.NoError(t, err)

	dims, err := GetImageDimensions(sourceImage)
	assert.Equal(t, image.ZP, dims)

	e, ok := err.(*os.PathError)
	assert.True(t, ok)
	assert.NotNil(t, e)
	assert.True(t, os.IsNotExist(e))
}

func TestGetImageDimensionsNotImage(t *testing.T) {
	cwd, _ := os.Getwd()
	sourceImage, err := filepath.Abs(path.Join(cwd, "..", "..", "test_images"))
	assert.NoError(t, err)

	dims, err := GetImageDimensions(sourceImage)
	assert.Equal(t, image.ZP, dims)

	assert.Equal(t, image.ErrFormat, err)
}

func TestGetOutputFilename(t *testing.T) {
	p := GetOutputFilename("/some/path", "north", false, "image.jpg")
	assert.Equal(t, "/some/path/image_north.jpg", p)

	p = GetOutputFilename("./some/path", "south", true, "image.png")
	assert.Equal(t, "some/path/image_scaled_south.png", p)
}

func TestOsMkdirp(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(d)

	mkdirp := path.Join(d, "new-dir")

	err = osMkdirp(mkdirp, 0755)
	assert.NoError(t, err)

	s, err := os.Stat(mkdirp)
	assert.NoError(t, err)

	var flm os.FileMode = os.ModeDir | 0755

	assert.True(t, s.IsDir())
	assert.Equal(t, flm, s.Mode())
}

func TestOsMkdirpFailsTooNested(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(d)

	mkdirp := path.Join(d, "new-dir", "cant-dir")

	err = osMkdirp(mkdirp, 0755)
	assert.True(t, os.IsNotExist(err))
}

func TestExtractGravitiesFromLocalImageScaled(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	doImageMagick = func(args ...string) (string, error) {
		assert.Equal(t, []string{"abc", "-gravity", "Center", "-scale", "64x64^", "-extent", "64x64", "images/abc_scaled_center"}, args)
		return "", nil
	}

	err := ExtractGravitiesFromLocalImage("abc", true, []string{"Center"}, "64x64", "images")
	assert.NoError(t, err)
}

func TestExtractGravitiesFromLocalImageUnscaled(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	doImageMagick = func(args ...string) (string, error) {
		assert.Equal(t, []string{"abc", "-gravity", "Center", "-extent", "64x64", "images/abc_center"}, args)
		return "", nil
	}

	err := ExtractGravitiesFromLocalImage("abc", false, []string{"Center"}, "64x64", "images")
	assert.NoError(t, err)
}

func TestExtractGravitiesFromLocalImageAlreadyExists(t *testing.T) {
	f := doImageMagick
	defer func() {
		doImageMagick = f
	}()

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	outputPath := path.Join(tempDir, "abc_center")

	_, err = os.Create(outputPath)
    assert.NoError(t, err)

	doImageMagick = func(args ...string) (string, error) {
		assert.Fail(t, "Imagemagick should not be called in this test")
		return "", nil
	}

	err = ExtractGravitiesFromLocalImage("abc", false, []string{"Center"}, "64x64", tempDir)
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

	is, err := PrepareImageFromSource(sourceImage, "")
	assert.NoError(t, err)

	tempDir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	err = ExtractFromImage("64x64", tempDir, is)
	assert.NoError(t, err)
}
