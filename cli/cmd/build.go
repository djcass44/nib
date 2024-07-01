package cmd

import (
	"fmt"
	"github.com/Snakdy/container-build-engine/pkg/builder"
	"github.com/Snakdy/container-build-engine/pkg/containers"
	"github.com/Snakdy/container-build-engine/pkg/pipelines"
	"github.com/djcass44/nib/cli/internal/dotenv"
	"github.com/djcass44/nib/cli/internal/packager"
	"github.com/djcass44/nib/cli/internal/pathfinder"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/spf13/cobra"
	"log"
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
	buildCmd.Flags().StringSlice(flagBuildPath, pathfinder.DefaultBuildPaths, "Which folders to check for compiled static files")
	buildCmd.Flags().Bool(flagSkipDotEnv, false, "Skip copying of the .env file")
	buildCmd.Flags().String(flagPlatform, "linux/amd64", "build platform")
	buildCmd.Flags().String(flagSave, "", "path to save the image as a tar archive")
}

func buildExec(cmd *cobra.Command, args []string) error {
	workingDir := args[0]
	localPath, _ := cmd.Flags().GetString(flagSave)
	cacheDir := os.Getenv(EnvCache)
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), ".nib-cache")
	}
	buildDirs, _ := cmd.Flags().GetStringSlice(flagBuildPath)
	if len(buildDirs) == 0 {
		log.Printf("--build-path list is empty, using default: %v", pathfinder.DefaultBuildPaths)
		buildDirs = pathfinder.DefaultBuildPaths
	}
	skipDotEnv, _ := cmd.Flags().GetBool(flagSkipDotEnv)

	platform, _ := cmd.Flags().GetString(flagPlatform)
	imgPlatform, err := v1.ParsePlatform(platform)
	if err != nil {
		return err
	}
	dataPath := "${NIB_DATA_PATH:-/var/run/nib}"

	statements := []pipelines.OrderedPipelineStatement{
		{
			ID: "set-env",
			Options: map[string]any{
				"NPM_CONFIG_LOGLEVEL": "error",
				"YARN_CACHE_FOLDER":   cacheDir,
				"PATH":                os.Getenv("PATH"),
			},
			Statement: &pipelines.Env{},
		},
		{
			ID: packager.StatementNodePackage,
			Options: map[string]any{
				"cache-dir": cacheDir,
			},
			Statement: &packager.NodePackager{},
			DependsOn: []string{"set-env"},
		},
		{
			ID: pathfinder.StatementPathfinder,
			Options: map[string]any{
				"build-dirs": buildDirs,
			},
			Statement: &pathfinder.Pathfinder{},
			DependsOn: []string{packager.StatementNodePackage},
		},
		{
			ID: "copy-build-dir",
			Options: map[string]any{
				"dst": dataPath,
			},
			Statement: &pipelines.Dir{},
			DependsOn: []string{
				packager.StatementNodePackage,
				pathfinder.StatementPathfinder,
			},
		},
		{
			ID: dotenv.StatementDotenv,
			Options: map[string]any{
				"skip": skipDotEnv,
				"path": dataPath,
			},
			Statement: &dotenv.Dotenv{},
			DependsOn: []string{packager.StatementNodePackage, "copy-build-dir"},
		},
		{
			ID: "set-runtime-env",
			Options: map[string]any{
				"PATH":          "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/home/somebody/.local/bin:/home/somebody/bin:/ko-app",
				"NIB_DATA_PATH": dataPath,
			},
			Statement: &pipelines.Env{},
			DependsOn: []string{"copy-build-dir"},
		},
	}

	// 4. add static files to base image
	baseImage := os.Getenv(EnvBaseImage)
	if baseImage == "" {
		baseImage = "ghcr.io/djcass44/nib/srv"
	}
	b, err := builder.NewBuilder(cmd.Context(), baseImage, statements, builder.Options{
		WorkingDir: workingDir,
		Metadata: builder.MetadataOptions{
			CreatedBy: "nib",
		},
	})
	if err != nil {
		return err
	}
	img, err := b.Build(cmd.Context(), imgPlatform)
	if err != nil {
		return err
	}
	if localPath != "" {
		return containers.Save(cmd.Context(), img, "image", localPath)
	}
	tags, _ := cmd.Flags().GetStringSlice(flagTag)
	for _, tag := range tags {
		if err := containers.Push(cmd.Context(), img, fmt.Sprintf("%s:%s", os.Getenv(EnvDockerRepo), tag)); err != nil {
			return err
		}
	}

	return nil
}
