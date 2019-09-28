package main;

import (
    "fmt"
    "os"
)

import (
    "github.com/spf13/cobra"
)

func main() {
    baseCommand := &cobra.Command{
        Use: "wp",
        Short: "Wallpaper generator CLI",
        Long: "Create many slices of an image passed in",
        RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("I did a thing")
			return nil
        },
    }

    if err := baseCommand.Execute(); err != nil {
        os.Exit(1)
    }
    os.Exit(0)
}
