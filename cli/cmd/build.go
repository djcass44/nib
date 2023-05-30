package cmd

import (
	"fmt"
	"github.com/djcass44/nib/cli/internal/build"
	"github.com/djcass44/nib/cli/internal/packager"
	"github.com/djcass44/nib/cli/internal/pathfinder"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build and publish container images from the given directory.",
	Long:  "This sub-command builds the provided directory into static files, containerises them, and publishes them.",
	Args:  cobra.ExactArgs(1),
	RunE:  buildExec,
}

func init() {
	buildCmd.Flags().StringSliceP(flagTag, "t", []string{"latest"}, "Which tags to use for the produced image instead of the default 'latest' tag")
}

var buildEngines = []packager.PackageManager{
	&packager.NPM{},
	&packager.Yarn{},
}

func buildExec(cmd *cobra.Command, args []string) error {
	workingDir := args[0]
	cacheDir := os.Getenv(EnvCache)
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), ".nib-cache")
	}

	bctx := packager.BuildContext{
		WorkingDir: workingDir,
		CacheDir:   cacheDir,
		Clock:      chronos.DefaultClock,
		Logger:     scribe.NewLogger(os.Stdout),
	}
	// 1. install dependencies
	pkg := buildEngines[0]
	for _, engine := range buildEngines {
		ok := engine.Detect(cmd.Context(), bctx)
		if ok {
			pkg = engine
			break
		}
	}
	err := pkg.Install(cmd.Context(), bctx)
	if err != nil {
		return err
	}

	// 2. build
	err = pkg.Build(cmd.Context(), bctx)
	if err != nil {
		return err
	}

	// 3. figure out where our static files were
	// just put
	appPath, err := pathfinder.FindBuildDir(workingDir)
	if err != nil {
		return err
	}

	platform, err := v1.ParsePlatform("linux/amd64")
	if err != nil {
		return err
	}

	// 4. add static files to base image
	baseImage := os.Getenv(EnvBaseImage)
	if baseImage == "" {
		baseImage = "ghcr.io/djcass44/nib/srv"
	}
	img, err := build.Append(cmd.Context(), appPath, baseImage, platform)
	if err != nil {
		return err
	}
	tags, _ := cmd.Flags().GetStringSlice(flagTag)
	for _, tag := range tags {
		if err := build.Push(cmd.Context(), img, fmt.Sprintf("%s:%s", os.Getenv(EnvDockerRepo), tag)); err != nil {
			return err
		}
	}

	return nil
}
