package wpservice;

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
) []error {
    sourceImageBasename := path.Base(sourcePath)
    sourceImageExtension := path.Ext(sourcePath)
    destImagePrefix := sourceImageBasename[:len(sourceImageBasename) - len(sourceImageExtension)]

    var errors []error
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
            scaledArgs := []string{"-scaled", dimensions + "^"}
            args := append(imagemagickArgs[0:3], scaledArgs...)
            imagemagickArgs = append(args, imagemagickArgs[3:]...)
        }

        cmd := exec.Command("convert", imagemagickArgs...)

        err := cmd.Run()

        if err != nil {
            errors = append(errors, err)
        }
    }

    return errors
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
        return errors.New("Image is not wide enough to produce quality output")
    }

    if imageSize.Y < desiredSize.Y {
        return errors.New("Image is not tall enough to produce quality output")
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
            "East",
            "Center",
            "West",
        }
    } else {
        scaledGravities = []string{
            "North",
            "Center",
            "South",
        }
    }

    errs := ExtractGravitiesFromSourceImage(
        localPath,
        true,
        scaledGravities,
        intendedDimensions,
        destinationDirComplete,
    )

    for _, err := range errs {
        fmt.Fprintln(os.Stderr, err)
    }

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

    errs2 := ExtractGravitiesFromSourceImage(
        localPath,
        false,
        unscaledGravities,
        intendedDimensions,
        destinationDirComplete,
    )

    for _, err := range errs2 {
        fmt.Fprintln(os.Stderr, err)
    }

    if len(errs) != 0 || len(errs2) != 0 {
        return errors.New("Some export step failed. See above for more information")
    }

    return nil
}
