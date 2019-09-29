package main;

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
    "strings"
)

import (
    "github.com/spf13/cobra"
)

import (
    "gitea.internal.aleemhaji.com/aleem/wp/cmd/wpservice"
)

func main() {
    // dimensionsRegexp := regexp.MustCompile(`^(\d+)x(\d+)$`)
    epsilon := math.Nextafter(1, 2) - 1

    baseCommand := &cobra.Command{
        Use: "wp <DesiredDimensions> <DestinationDir> <ImagePath> [ImagePath...] ",
        Short: "Wallpaper Generator CLI",
        Long: "Create many different  of an image passed in",
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            desiredDimensions := args[0]
            destinationDir := args[1]
            imagePath := args[2]

            destinationDirComplete, err := filepath.Abs(path.Join(destinationDir, desiredDimensions))
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
            sourceImage, err := os.Open(imagePath)
            if err != nil {
                return err
            }

            image, _, err := image.Decode(sourceImage)
            if err != nil {
                return err
            }

            imageBoundingRect := image.Bounds()
            if imageBoundingRect.Min.X != 0 || imageBoundingRect.Min.Y != 0 {
                return errors.New("Don't know how to deal with non-origin-point images")
            }

            desiredBoundingRect, err := wpservice.ParseDimensionsString(desiredDimensions)
            if err != nil {
                return err
            }

            if imageBoundingRect.Max.X < desiredBoundingRect.X {
                return errors.New("Image is not wide enough to produce quality output")
            }

            if imageBoundingRect.Max.Y < desiredBoundingRect.Y {
                return errors.New("Image is not tall enough to produce quality output")
            }

            // Check aspect ratio to know which direction scaled images will be
            //   sliced.
            // There will be a lot of duplicates without this step.
            desiredAspectRatio := float64(imageBoundingRect.Max.X) / float64(imageBoundingRect.Max.Y)
            imageAspectRatio := float64(desiredBoundingRect.X) / float64(desiredBoundingRect.Y)

            var scaledGravities []string = nil
            if math.Abs(desiredAspectRatio - imageAspectRatio) < epsilon {
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

            sourceImageBasename := path.Base(imagePath)
            sourceImageExtension := path.Ext(imagePath)
            destImagePrefix := sourceImageBasename[:len(sourceImageBasename) - len(sourceImageExtension)]

            for _, gravity := range scaledGravities {
                outputFilename := destImagePrefix + "_scaled_" + strings.ToLower(gravity) + sourceImageExtension
                outputPath := path.Join(destinationDirComplete, outputFilename)

                if _, err := os.Stat(outputPath); err == nil {
                  continue
                }

                cmd := exec.Command(
                    "convert", 
                    imagePath,
                    "-gravity", gravity,
                    "-scale", desiredDimensions + "^",
                    "-extent", desiredDimensions,
                    outputPath,
                )

                err := cmd.Run()

                if err != nil {
                    fmt.Fprintln(os.Stderr, "Error creating", outputFilename, err)
                }
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

            for _, gravity := range unscaledGravities {
                outputFilename := destImagePrefix + "_" + strings.ToLower(gravity) + sourceImageExtension
                outputPath := path.Join(destinationDirComplete, outputFilename)

                if _, err := os.Stat(outputPath); err == nil {
                  continue
                }

                cmd := exec.Command(
                    "convert", 
                    imagePath,
                    "-gravity", gravity,
                    "-extent", desiredDimensions,
                    outputPath,
                )

                err := cmd.Run()

                if err != nil {
                    fmt.Fprintln(os.Stderr, "Error creating", outputFilename, err)
                }
            }

			return nil
        },
    }

    if err := baseCommand.Execute(); err != nil {
        os.Exit(1)
    }
    os.Exit(0)
}
