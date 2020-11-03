package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print the version number and exit",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(os.Args[0] + ": " + VersionBuild)
		return nil
	},
}
