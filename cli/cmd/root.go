package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var command = &cobra.Command{
	Use:   "nib",
	Short: "",
}

func init() {
	command.AddCommand(buildCmd)
}

func Execute(version string) {
	command.Version = version
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
