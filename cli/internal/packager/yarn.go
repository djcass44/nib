package packager

import (
	"context"
	"fmt"
	"github.com/djcass44/nib/cli/pkg/executor"
	"os"
	"path/filepath"
	"strings"
)

const commandYarn = "yarn"
const lockfileYarn = "yarn.lock"

type Yarn struct{}

// Detect checks to see if the build directory contains
// a Yarn lock file
func (*Yarn) Detect(_ context.Context, bctx executor.BuildContext) bool {
	bctx.Logger.Process("Checking for Yarn lockfile")

	_, err := os.Stat(filepath.Join(bctx.WorkingDir, lockfileYarn))
	return err == nil
}

// Install installs packages using Yarn
func (*Yarn) Install(_ context.Context, bctx executor.BuildContext) error {
	bctx.Logger.Process("Executing install process")

	var extraArgs []string
	if val := os.Getenv(executor.EnvExtraArgs); val != "" {
		extraArgs = strings.Split(val, " ")
	}

	return executor.Exec(bctx, executor.Options{
		ExtraEnv: []string{fmt.Sprintf("YARN_CACHE_FOLDER=%s", bctx.CacheDir)},
		Command:  commandYarn,
		Args:     append([]string{"install", "--immutable"}, extraArgs...),
	})
}

// Build runs the Yarn build script
func (*Yarn) Build(_ context.Context, bctx executor.BuildContext) error {
	bctx.Logger.Process("Executing build process")

	return executor.Exec(bctx, executor.Options{
		ExtraEnv: []string{fmt.Sprintf("YARN_CACHE_FOLDER=%s", bctx.CacheDir)},
		Command:  commandYarn,
		Args:     []string{"run", "build"},
	})
}
