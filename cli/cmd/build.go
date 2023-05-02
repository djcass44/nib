package cmd

import (
	"github.com/djcass44/nib/cli/internal/packager"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/spf13/cobra"
	"os"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE:  build,
}

func build(cmd *cobra.Command, args []string) error {
	workingDir := args[0]
	// 1. install dependencies
	pkg := packager.NPM{}
	err := pkg.Install(cmd.Context(), packager.BuildContext{
		WorkingDir: workingDir,
		Clock:      chronos.DefaultClock,
		Logger:     scribe.NewLogger(os.Stdout),
	})
	if err != nil {
		return err
	}

	// 2. build

	// 3. add static files to base image

	return nil
}
