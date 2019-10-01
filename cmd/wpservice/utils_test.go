package wpservice

import (
	"image"
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

	ExtractGravitiesFromSourceImage("abc", true, []string{"Center"}, "64x64", "images")
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

	ExtractGravitiesFromSourceImage("abc", false, []string{"Center"}, "64x64", "images")
}
