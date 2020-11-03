package cmd

import (
	"fmt"
	"os"

	"github.com/Eagerod/wp/pkg/wp"
	"github.com/spf13/cobra"
)

var extractCommand = &cobra.Command{
	Use:   "extract desired_dimensions destination_dir image_path [image_path...]",
	Short: "Extract image slices",
	Long:  "Create many different slices of an image passed in",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		desiredDimensions := args[0]
		destinationDir := args[1]
		imagePaths := args[2:]

		var errs []error
		for _, imagePath := range imagePaths {
			is, err := wp.PrepareImageFromSource(imagePath, cacheDir)
			if err != nil {
				return err
			}
			defer wp.CleanupImageSource(is)

			err = wp.ExtractFromImage(desiredDimensions, destinationDir, is)

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
