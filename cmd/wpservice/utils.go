package wpservice

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type ImageMagickRunner func(args ...string) (string, error)

const imagemagickBin string = "convert"

var dimensionsRegexp *regexp.Regexp = regexp.MustCompile(`^(\d+)x(\d+)$`)

var doImageMagick ImageMagickRunner = func(args ...string) (string, error) {
	cmd := exec.Command(imagemagickBin, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

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

/*
  Run imagemagick against the provided source path and generate crops or
  rescales of the image.
*/
func ExtractGravitiesFromSourceImage(
	sourcePath string,
	scaled bool,
	gravities []string,
	dimensions string,
	output string,
) error {
	sourceImageBasename := path.Base(sourcePath)
	sourceImageExtension := path.Ext(sourcePath)
	destImagePrefix := sourceImageBasename[:len(sourceImageBasename)-len(sourceImageExtension)]

	var errs []error
	for _, gravity := range gravities {
		var outputFilename string
		if scaled {
			outputFilename = destImagePrefix + "_scaled_" + strings.ToLower(gravity) + sourceImageExtension
		} else {
			outputFilename = destImagePrefix + "_" + strings.ToLower(gravity) + sourceImageExtension
		}
		outputPath := path.Join(output, outputFilename)

		if _, err := os.Stat(outputPath); err == nil {
			continue
		}

		imagemagickArgs := []string{
			sourcePath,
			"-gravity", gravity,
			"-extent", dimensions,
			outputPath,
		}

		if scaled {
			scaledArgs := []string{"-scale", dimensions + "^"}
			imagemagickArgs = append(imagemagickArgs[:3], append(scaledArgs, imagemagickArgs[3:]...)...)
		}

		fmt.Fprintln(os.Stderr, outputPath)

		output, err := doImageMagick(imagemagickArgs...)

		if err != nil {
			errs = append(errs, errors.New(output))
			errs = append(errs, err)
		}
	}

	multiError := MultiErrorFromErrors(errs)
	if multiError.Exists() {
		return multiError
	}

	return nil
}

func ExtractFromLocalImage(intendedDimensions string, destination string, localPath string) error {
	epsilon := math.Nextafter(1, 2) - 1

	destinationDirComplete, err := filepath.Abs(path.Join(destination, intendedDimensions))
	if err != nil {
		return err
	}

	dirFilemode := os.ModeDir | 0755
	if err := os.Mkdir(destinationDirComplete, dirFilemode); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	// Get provided image dimensions
	sourceImage, err := os.Open(localPath)
	if err != nil {
		return err
	}

	image, _, err := image.Decode(sourceImage)
	if err != nil {
		return err
	}

	imageBoundingRect := image.Bounds()
	imageOrigin := imageBoundingRect.Min
	imageSize := imageBoundingRect.Max
	if imageOrigin.X != 0 || imageOrigin.Y != 0 {
		return errors.New("Don't know how to deal with non-origin-point images")
	}

	desiredSize, err := ParseDimensionsString(intendedDimensions)
	if err != nil {
		return err
	}

	if imageSize.X < desiredSize.X {
		return errors.New(fmt.Sprintf("Image (%s) is not wide enough to produce quality output", path.Base(localPath)))
	}

	if imageSize.Y < desiredSize.Y {
		return errors.New(fmt.Sprintf("Image (%s) is not tall enough to produce quality output", path.Base(localPath)))
	}

	// Check aspect ratio to know which direction scaled images will be
	//   sliced.
	// There will be a lot of duplicates without this step.
	desiredAspectRatio := float64(imageSize.X) / float64(imageSize.Y)
	imageAspectRatio := float64(desiredSize.X) / float64(desiredSize.Y)

	var scaledGravities []string = nil
	if math.Abs(desiredAspectRatio-imageAspectRatio) < epsilon {
		scaledGravities = []string{
			"Center",
		}
	} else if desiredAspectRatio > imageAspectRatio {
		scaledGravities = []string{
			"West",
			"Center",
			"East",
		}
	} else {
		scaledGravities = []string{
			"North",
			"Center",
			"South",
		}
	}

	err1 := ExtractGravitiesFromSourceImage(
		localPath,
		true,
		scaledGravities,
		intendedDimensions,
		destinationDirComplete,
	)

	unscaledGravities := []string{
		"North",
		"NorthEast",
		"East",
		"SouthEast",
		"South",
		"SouthWest",
		"West",
		"NorthWest",
		"Center",
	}

	err2 := ExtractGravitiesFromSourceImage(
		localPath,
		false,
		unscaledGravities,
		intendedDimensions,
		destinationDirComplete,
	)

	if err1 != nil || err2 != nil {
		return MultiErrorFromErrors([]error{err1, err2})
	}

	return nil
}
