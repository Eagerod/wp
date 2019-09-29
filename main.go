package main

// Need to load up the image libraries for them to be registered for decoding.
// Yay side effects!
import (
	"os"
)

import (
	"github.com/spf13/cobra"
)

import (
	"gitea.internal.aleemhaji.com/aleem/wp/cmd/wpservice"
)

func main() {
	baseCommand := &cobra.Command{
		Use:   "wp <DesiredDimensions> <DestinationDir> <ImagePath> [ImagePath...] ",
		Short: "Wallpaper Generator CLI",
		Long:  "Create many different  of an image passed in",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			desiredDimensions := args[0]
			destinationDir := args[1]
			imagePath := args[2]

			return wpservice.ExtractFromLocalImage(desiredDimensions, destinationDir, imagePath)
		},
	}

	if err := baseCommand.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
