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
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE:  buildExec,
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
	pkg := packager.NPM{}
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
		baseImage = "harbor.dcas.dev/public.ecr.aws/bitnami/nginx:1.23"
	}
	img, err := build.Append(cmd.Context(), appPath, baseImage, platform)
	if err != nil {
		return err
	}
	tags := []string{"latest"}
	for _, tag := range tags {
		if err := build.Push(cmd.Context(), img, fmt.Sprintf("%s:%s", os.Getenv(EnvDockerRepo), tag)); err != nil {
			return err
		}
	}

	return nil
}
