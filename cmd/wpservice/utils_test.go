package wpservice;

import (
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
