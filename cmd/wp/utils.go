package wp

// Need to load up the image libraries for them to be registered for decoding.
// Yay side effects!
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
	"strings"
)

type ImageMagickRunner func(args ...string) (string, error)

const imagemagickBin string = "convert"

var dimensionsRegexp *regexp.Regexp = regexp.MustCompile(`^(\d+)x(\d+)$`)
var epsilon float64 = math.Nextafter(1, 2) - 1

// Gravity sets:
var equalAspectRatioGravities []string = []string{
	"Center",
}
var tallAspectRatioGravities []string = []string{
	"North",
	"Center",
	"South",
}
var wideAspectRatioGravities []string = []string{
	"West",
	"Center",
	"East",
}
var unscaledGravities []string = []string{
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

var doImageMagick ImageMagickRunner = func(args ...string) (string, error) {
	cmd := exec.Command(imagemagickBin, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Get the dimensions of an image at the path passed in.
func GetImageDimensions(imagePath string) (image.Point, error) {
	sourceImage, err := os.Open(imagePath)
	if err != nil {
		return image.ZP, err
	}

	img, _, err := image.Decode(sourceImage)
	if err != nil {
		return image.ZP, err
	}

	imageBoundingRect := img.Bounds()
	imageOrigin := imageBoundingRect.Min
	imageSize := imageBoundingRect.Max
	if imageOrigin.X != 0 || imageOrigin.Y != 0 {
		return image.ZP, errors.New("Don't know how to deal with non-origin-point images")
	}

	return imageSize, nil
}

// Get the final output filename of writing an image with the given parameters.
func GetOutputFilename(outputDir string, gravity string, scaled bool, sourcePath string) string {
	sourceImageBasename := path.Base(sourcePath)
	sourceImageExtension := path.Ext(sourcePath)
	destImagePrefix := sourceImageBasename[:len(sourceImageBasename)-len(sourceImageExtension)]

	var outputFilename string
	if scaled {
		outputFilename = destImagePrefix + "_scaled_" + strings.ToLower(gravity) + sourceImageExtension
	} else {
		outputFilename = destImagePrefix + "_" + strings.ToLower(gravity) + sourceImageExtension
	}
	return path.Join(outputDir, outputFilename)
}

// mkdir -p
func osMkdirp(p string, mode os.FileMode) error {
	// dirFilemode := os.ModeDir | mode
	if err := os.Mkdir(p, mode); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	return nil
}

/*
  Run imagemagick against the provided source path and generate crops or
  rescales of the image.
*/
func ExtractGravitiesFromLocalImage(
	sourcePath string,
	scaled bool,
	gravities []string,
	dimensions string,
	output string,
) error {
	var errs []error
	for _, gravity := range gravities {
		outputPath := GetOutputFilename(output, gravity, scaled, sourcePath)

		if _, err := os.Stat(outputPath); err == nil {
			fmt.Fprintln(os.Stderr, outputPath)
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

	if err := MultiErrorFromErrors(errs); err.Exists() {
		return err
	}

	return nil
}

func ExtractFromLocalImage(intendedDimensions string, destination string, localPath string) error {
	// Check to make sure the passed in output dimensions are valid before
	//   creating the directory.
	desiredSize, err := ParseDimensionsString(intendedDimensions)
	if err != nil {
		return err
	}

	imageSize, err := GetImageDimensions(localPath)
	if err != nil {
		return err
	}

	if imageSize.X < desiredSize.X {
		return errors.New(fmt.Sprintf("Image (%s) is not wide enough to produce quality output", path.Base(localPath)))
	}

	if imageSize.Y < desiredSize.Y {
		return errors.New(fmt.Sprintf("Image (%s) is not tall enough to produce quality output", path.Base(localPath)))
	}

	destinationDirComplete, err := filepath.Abs(path.Join(destination, intendedDimensions))
	if err != nil {
		return err
	}

	if err := osMkdirp(destinationDirComplete, 0755); err != nil {
		return err
	}

	// Check aspect ratio to know which direction scaled images will be
	//   sliced.
	// There will be a lot of duplicates without this step.
	desiredAspectRatio := float64(imageSize.X) / float64(imageSize.Y)
	imageAspectRatio := float64(desiredSize.X) / float64(desiredSize.Y)

	var scaledGravities []string = nil
	if math.Abs(desiredAspectRatio-imageAspectRatio) < epsilon {
		scaledGravities = equalAspectRatioGravities
	} else if desiredAspectRatio > imageAspectRatio {
		scaledGravities = wideAspectRatioGravities
	} else {
		scaledGravities = tallAspectRatioGravities
	}

	err1 := ExtractGravitiesFromLocalImage(
		localPath,
		true,
		scaledGravities,
		intendedDimensions,
		destinationDirComplete,
	)

	err2 := ExtractGravitiesFromLocalImage(
		localPath,
		false,
		unscaledGravities,
		intendedDimensions,
		destinationDirComplete,
	)

	if err := MultiErrorFromErrors([]error{err1, err2}); err.Exists() {
		return err
	}

	return nil
}

func ExtractFromImage(intendedDimensions string, destination string, imageSource *ImageSource) error {
	return ExtractFromLocalImage(intendedDimensions, destination, imageSource.LocalPath)
}

func PickFromImage(intendedDimensions string, destination string, imageSource *ImageSource, scaled bool, gravity string) error {
	destination = path.Join(destination, intendedDimensions)
	if err := osMkdirp(destination, 0755); err != nil {
		return err
	}

	return ExtractGravitiesFromLocalImage(imageSource.LocalPath, scaled, []string{gravity}, intendedDimensions, destination)
}
