package wpservice;

import (
    "errors"
    "image"
    "regexp"
    "strconv"
)

/*
  Parse a string in the form <x>x<y> and return a Point specifying the extents
*/
func ParseDimensionsString(str string) (image.Point, error) {
    dimensionsRegexp := regexp.MustCompile(`^(\d+)x(\d+)$`)

    dimensionsMatch := dimensionsRegexp.FindStringSubmatch(str)

    width, err := strconv.Atoi(dimensionsMatch[1])
    if err != nil {
        return image.ZP, errors.New("Provided width is not a valid integer")
    }

    height, err := strconv.Atoi(dimensionsMatch[2])
    if err != nil {
        return image.ZP, errors.New("Provided height is not a valid integer")
    }

    return image.Pt(width, height), nil
}
