package packager

import (
	"context"
	"github.com/djcass44/nib/cli/pkg/executor"
	"os"
	"path/filepath"
	"strings"
)

const commandNPM = "npm"
const lockfileNPM = "package-lock.json"

type NPM struct{}

// Detect checks to see if the build directory contains
// an NPM lock file
func (*NPM) Detect(_ context.Context, bctx executor.BuildContext) bool {
	bctx.Logger.Process("Checking for NPM lockfile")

	_, err := os.Stat(filepath.Join(bctx.WorkingDir, lockfileNPM))
	return err == nil
}

// Install installs packages using NPM
func (*NPM) Install(_ context.Context, bctx executor.BuildContext) error {
	bctx.Logger.Process("Executing install process")

	var extraArgs []string
	if val := os.Getenv(executor.EnvExtraArgs); val != "" {
		extraArgs = strings.Split(val, " ")
	}

	return executor.Exec(bctx, executor.Options{
		ExtraEnv: []string{"NPM_CONFIG_LOGLEVEL=error"},
		Command:  commandNPM,
		Args:     append([]string{"ci", "--include=dev", "--unsafe-perm", "--cache", bctx.CacheDir}, extraArgs...),
	})
}

// Build runs the NPM build script
func (*NPM) Build(_ context.Context, bctx executor.BuildContext) error {
	bctx.Logger.Process("Executing build process")

	return executor.Exec(bctx, executor.Options{
		ExtraEnv: []string{"NPM_CONFIG_LOGLEVEL=error"},
		Command:  commandNPM,
		Args:     []string{"run", "build"},
	})
}
