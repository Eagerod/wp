package main

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
	"path"
	"path/filepath"
)

import (
	"github.com/spf13/cobra"
)

import (
	"gitea.internal.aleemhaji.com/aleem/wp/cmd/wpservice"
)

func main() {
	epsilon := math.Nextafter(1, 2) - 1

	baseCommand := &cobra.Command{
		Use:   "wp <DesiredDimensions> <DestinationDir> <ImagePath> [ImagePath...] ",
		Short: "Wallpaper Generator CLI",
		Long:  "Create many different  of an image passed in",
		Args:  cobra.ExactArgs(3),
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
			imageOrigin := imageBoundingRect.Min
			imageSize := imageBoundingRect.Max
			if imageOrigin.X != 0 || imageOrigin.Y != 0 {
				return errors.New("Don't know how to deal with non-origin-point images")
			}

			desiredSize, err := wpservice.ParseDimensionsString(desiredDimensions)
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

			errs := wpservice.ExtractGravitiesFromSourceImage(
				imagePath,
				true,
				scaledGravities,
				desiredDimensions,
				destinationDirComplete,
			)

			for err := range errs {
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

			errs2 := wpservice.ExtractGravitiesFromSourceImage(
				imagePath,
				false,
				unscaledGravities,
				desiredDimensions,
				destinationDirComplete,
			)

			for err := range errs2 {
				fmt.Fprintln(os.Stderr, err)
			}

			if len(errs) != 0 || len(errs2) != 0 {
				return errors.New("Some export step failed. See above for more information")
			}

			return nil
		},
	}

	if err := baseCommand.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
