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
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			desiredDimensions := args[0]
			destinationDir := args[1]
			imagePath := args[2]

			err := wpservice.ExtractFromLocalImage(desiredDimensions, destinationDir, imagePath)

			// If the only thing the error is is a series of soft errors, don't
			//   exit with failure.
			if err != nil {
				if softErrorRegexp.FindStringSubmatch(err.Error()) == nil {
					return err
				}

				fmt.Fprintln(os.Stderr, err.Error())
			}

			return nil
		},
	}

	if err := baseCommand.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
