package cmd

import "github.com/spf13/cobra"

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "",
	RunE:  build,
}

func build(cmd *cobra.Command, _ []string) error {
	return nil
}
