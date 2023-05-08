package packager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const commandYarn = "yarn"
const lockfileYarn = "yarn.lock"

type Yarn struct{}

func (*Yarn) Detect(_ context.Context, bctx BuildContext) bool {
	bctx.Logger.Process("Checking for Yarn lockfile")

	_, err := os.Stat(filepath.Join(bctx.WorkingDir, lockfileYarn))
	return err == nil
}

func (*Yarn) Install(_ context.Context, bctx BuildContext) error {
	bctx.Logger.Process("Executing install process")

	return exec(bctx, options{
		extraEnv: []string{fmt.Sprintf("YARN_CACHE_FOLDER=%s", bctx.CacheDir)},
		command:  commandYarn,
		args:     []string{"install", "--immutable"},
	})
}

func (*Yarn) Build(_ context.Context, bctx BuildContext) error {
	bctx.Logger.Process("Executing build process")

	return exec(bctx, options{
		extraEnv: []string{fmt.Sprintf("YARN_CACHE_FOLDER=%s", bctx.CacheDir)},
		command:  commandYarn,
		args:     []string{"run", "build"},
	})
}
