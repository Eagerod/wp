package cmd

import (
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var softErrorRegexp *regexp.Regexp = regexp.MustCompile(`^(?:Image .*? is not (?:tall|wide) enough to produce quality output\n?)+$`)

var scaledFlag bool
var cacheDir string

var baseCommand = &cobra.Command{
	Use:   os.Args[0],
	Short: "Wallpaper Generator CLI",
	Long:  "Manipulate images for use as desktop wallpapers",
}

func Execute() {
	baseCommand.AddCommand(extractCommand)
	baseCommand.AddCommand(pickCommand)
	baseCommand.AddCommand(versionCommand)

	extractCommand.Flags().StringVarP(&cacheDir, "cache", "", "", "Source image cache; used to prevent repeated downloads")

	pickCommand.Flags().BoolVarP(&scaledFlag, "scaled", "", false, "Scale the image to the desired dimensions, rather than maintaining scale")
	pickCommand.Flags().StringVarP(&cacheDir, "cache", "", "", "Source image cache; used to prevent repeated downloads")

	if err := baseCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
