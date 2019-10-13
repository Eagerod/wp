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

var softErrorRegexp *regexp.Regexp = regexp.MustCompile(`^(?:Image is not (?:tall|wide) enough to produce quality output\n?)+$`)

func main() {
	baseCommand := &cobra.Command{
		Use:   "wp <DesiredDimensions> <DestinationDir> <ImagePath> [ImagePath...] ",
		Short: "Wallpaper Generator CLI",
		Long:  "Create many different slices of an image passed in",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			desiredDimensions := args[0]
			destinationDir := args[1]
			imagePaths := args[2:]

			var errs []error
			for _, imagePath := range imagePaths {
				err := wpservice.ExtractFromLocalImage(desiredDimensions, destinationDir, imagePath)

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

	if err := baseCommand.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
