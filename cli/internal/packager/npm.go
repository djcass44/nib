package packager

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

const commandNPM = "npm"
const lockfileNPM = "package-lock.json"

type NPM struct{}

func (*NPM) Detect(_ context.Context, bctx BuildContext) bool {
	bctx.Logger.Process("Checking for NPM lockfile")

	_, err := os.Stat(filepath.Join(bctx.WorkingDir, lockfileNPM))
	return err == nil
}

func (*NPM) Install(_ context.Context, bctx BuildContext) error {
	bctx.Logger.Process("Executing install process")

	var extraArgs []string
	if val := os.Getenv(EnvExtraArgs); val != "" {
		extraArgs = strings.Split(val, " ")
	}

	return exec(bctx, options{
		extraEnv: []string{"NPM_CONFIG_LOGLEVEL=error"},
		command:  commandNPM,
		args:     append([]string{"ci", "--include=dev", "--unsafe-perm", "--cache", bctx.CacheDir}, extraArgs...),
	})
}

func (*NPM) Build(_ context.Context, bctx BuildContext) error {
	bctx.Logger.Process("Executing build process")

	return exec(bctx, options{
		extraEnv: []string{"NPM_CONFIG_LOGLEVEL=error"},
		command:  commandNPM,
		args:     []string{"run", "build"},
	})
}
