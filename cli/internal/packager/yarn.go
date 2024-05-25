package packager

import (
	"github.com/djcass44/nib/cli/pkg/executor"
	"github.com/go-logr/logr"
	"os"
	"path/filepath"
	"strings"
)

const commandYarn = "yarn"
const lockfileYarn = "yarn.lock"

type Yarn struct{}

// Detect checks to see if the build directory contains
// a Yarn lock file
func (*Yarn) Detect(ctx executor.BuildContext) bool {
	log := logr.FromContextOrDiscard(ctx.Ctx.Context)
	log.Info("checking for Yarn lockfile")

	_, err := os.Stat(filepath.Join(ctx.Ctx.WorkingDirectory, lockfileYarn))
	return err == nil
}

// Install installs packages using Yarn
func (*Yarn) Install(ctx executor.BuildContext) error {
	log := logr.FromContextOrDiscard(ctx.Ctx.Context)
	log.Info("executing install process")

	var extraArgs []string
	if val := os.Getenv(executor.EnvExtraArgs); val != "" {
		extraArgs = strings.Split(val, " ")
	}

	return executor.Exec(ctx, executor.Options{
		Command: commandYarn,
		Args:    append([]string{"install", "--immutable"}, extraArgs...),
	})
}

// Build runs the Yarn build script
func (*Yarn) Build(ctx executor.BuildContext) error {
	log := logr.FromContextOrDiscard(ctx.Ctx.Context)
	log.Info("executing build process")

	return executor.Exec(ctx, executor.Options{
		Command: commandYarn,
		Args:    []string{"run", "build"},
	})
}
