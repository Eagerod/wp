package wpservice;

import (
    "errors"
    "image"
    "os"
    "os/exec"
    "path"
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

        var cmd *exec.Cmd 
        if scaled {
            cmd = exec.Command(
                "convert",
                sourcePath,
                "-gravity", gravity,
                "-scale", dimensions + "^",
                "-extent", dimensions,
                outputPath,
            )
        } else {
            cmd = exec.Command(
                "convert",
                sourcePath,
                "-gravity", gravity,
                "-extent", dimensions,
                outputPath,
            )
        }

        err := cmd.Run()

        if err != nil {
            errors = append(errors, err)
        }
    }

    return errors
}
