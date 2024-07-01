package packager

import (
	"github.com/djcass44/nib/cli/pkg/executor"
	"github.com/go-logr/logr"
	"os"
	"path/filepath"
	"strings"
)

const commandNPM = "npm"
const lockfileNPM = "package-lock.json"

func NewNPM(command string) *NPM {
	return &NPM{command: command}
}

type NPM struct {
	command string
}

// Detect checks to see if the build directory contains
// an NPM lock file
func (*NPM) Detect(ctx executor.BuildContext) bool {
	log := logr.FromContextOrDiscard(ctx.Ctx.Context)
	log.Info("checking for NPM lockfile")

	_, err := os.Stat(filepath.Join(ctx.Ctx.WorkingDirectory, lockfileNPM))
	return err == nil
}

// Install installs packages using NPM
func (n *NPM) Install(ctx executor.BuildContext) error {
	log := logr.FromContextOrDiscard(ctx.Ctx.Context)
	log.Info("executing install process")

	var extraArgs []string
	if val := os.Getenv(executor.EnvExtraArgs); val != "" {
		extraArgs = strings.Split(val, " ")
	}

	return executor.Exec(ctx, executor.Options{
		Command: n.command,
		Args:    append([]string{"ci", "--include=dev", "--unsafe-perm", "--cache", ctx.CacheDir}, extraArgs...),
	})
}

// Build runs the NPM build script
func (n *NPM) Build(ctx executor.BuildContext) error {
	log := logr.FromContextOrDiscard(ctx.Ctx.Context)
	log.Info("executing build process")

	return executor.Exec(ctx, executor.Options{
		Command: n.command,
		Args:    []string{"run", "build"},
	})
}
