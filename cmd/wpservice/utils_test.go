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
