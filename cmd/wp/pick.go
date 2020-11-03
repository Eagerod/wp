package cmd

import (
	"fmt"
	"os"

	"github.com/Eagerod/wp/pkg/wp"
	"github.com/spf13/cobra"
)

var pickCommand = &cobra.Command{
	Use:   "pick desired_dimensions destination_dir gravity [--scaled] image_path [image_path...]",
	Short: "Pick a single image slice",
	Long:  "Extract a single slice of an image with the given parameters",
	Args:  cobra.MinimumNArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		desiredDimensions := args[0]
		destinationDir := args[1]
		gravity := args[2]
		imagePaths := args[3:]

		var errs []error
		for _, imagePath := range imagePaths {
			is, err := wp.PrepareImageFromSource(imagePath, cacheDir)
			if err != nil {
				return err
			}
			defer wp.CleanupImageSource(is)

			err = wp.PickFromImage(desiredDimensions, destinationDir, is, scaledFlag, gravity)

			if err != nil {
				errs = append(errs, err)
			}
		}

		// If the only thing the error is is a series of soft errors, don't
		//   exit with failure.
		multiError := wp.MultiErrorFromErrors(errs)
		if multiError.Exists() {
			if softErrorRegexp.FindStringSubmatch(multiError.Error()) == nil {
				return multiError
			}

			fmt.Fprintln(os.Stderr, multiError.Error())
		}

		return nil
	},
}
