package wp

import (
	"errors"
	"fmt"
	"image"
	"strconv"
)

/*
  Parse a string in the form <x>x<y> and return a Point specifying the extents
*/
func ParseDimensionsString(str string) (image.Point, error) {
	dimensionsMatch := dimensionsRegexp.FindStringSubmatch(str)

	if len(dimensionsMatch) == 0 {
		return image.ZP, errors.New(fmt.Sprintf("Provided dimension string (%s) is not valid", str))
	}

	width, err := strconv.Atoi(dimensionsMatch[1])
	if err != nil || width <= 0 {
		return image.ZP, errors.New("Provided width is not a valid positive integer")
	}

	height, err := strconv.Atoi(dimensionsMatch[2])
	if err != nil || height <= 0 {
		return image.ZP, errors.New("Provided height is not a valid positive integer")
	}

	return image.Pt(width, height), nil
}
