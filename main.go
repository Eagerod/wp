package main

// Need to load up the image libraries for them to be registered for decoding.
// Yay side effects!
import (
	"fmt"
	"os"
	"regexp"
)

import (
	"github.com/spf13/cobra"
)

import (
	"gitea.internal.aleemhaji.com/aleem/wp/cmd/wpservice"
)

var softErrorRegexp *regexp.Regexp = regexp.MustCompile(`^(?:Image .*? is not (?:tall|wide) enough to produce quality output\n?)+$`)

func main() {
	var scaledFlag bool 

	baseCommand := &cobra.Command{
		Use:   "wp",
		Short: "Wallpaper Generator CLI",
		Long:  "Manipulate images for use as desktop wallpapers",
	}

	extractCommand := &cobra.Command{
		Use: "extract desired_dimensions destination_dir image_path [image_path...]",
		Short: "Extract image slices",
		Long:  "Create many different slices of an image passed in",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			desiredDimensions := args[0]
			destinationDir := args[1]
			imagePaths := args[2:]

			var errs []error
			for _, imagePath := range imagePaths {
				err := wpservice.ExtractFromImage(desiredDimensions, destinationDir, imagePath)

				if err != nil {
					errs = append(errs, err)
				}
			}

			// If the only thing the error is is a series of soft errors, don't
			//   exit with failure.
			multiError := wpservice.MultiErrorFromErrors(errs)
			if multiError.Exists() {
				if softErrorRegexp.FindStringSubmatch(multiError.Error()) == nil {
					return multiError
				}

				fmt.Fprintln(os.Stderr, multiError.Error())
			}

			return nil
		},
	}

	pickCommand := &cobra.Command{
		Use:   "pick desired_dimensions destination_dir gravity [--scaled] image_path",
		Short: "Pick a single image slice",
		Long:  "Extract a single slice of an image with the given parameters",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			desiredDimensions := args[0]
			destinationDir := args[1]
			gravity := args[2]
			imagePath := args[3]

			err := wpservice.PickFromImage(desiredDimensions, destinationDir, imagePath, scaledFlag, gravity)

			if err != nil {
				if softErrorRegexp.FindStringSubmatch(err.Error()) == nil {
					return err
				}

				fmt.Fprintln(os.Stderr, err.Error())
			}

			return nil
		},
	}

    baseCommand.AddCommand(extractCommand)
    baseCommand.AddCommand(pickCommand)

    pickCommand.Flags().BoolVarP(&scaledFlag, "scaled", "", false, "Scale the image to the desired dimensions, rather than maintaining scale")

	if err := baseCommand.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
