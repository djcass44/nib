package packager

import (
	"context"
	"fmt"
	"github.com/paketo-buildpacks/packit/pexec"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type NPM struct{}

func (*NPM) Install(_ context.Context, bctx BuildContext) error {
	bctx.Logger.Process("Executing install process")

	// assemble the information our executable needs
	exec := pexec.NewExecutable("npm")
	args := []string{"ci", "--include=dev", "--unsafe-perm", "--cache", filepath.Join(bctx.WorkingDir, ".cache")}
	bctx.Logger.Subprocess("Running 'npm %s'", strings.Join(args, " "))
	// shell out
	duration, err := bctx.Clock.Measure(func() error {
		return exec.Execute(pexec.Execution{
			Args: args,
			Dir:  bctx.WorkingDir,
			Env: append(
				os.Environ(),
				"NPM_CONFIG_LOGLEVEL=error",
			),
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})
	})
	if err != nil {
		return fmt.Errorf("npm ci failed: %w", err)
	}
	bctx.Logger.Action("Completed in %s", duration.Round(time.Millisecond))
	bctx.Logger.Break()
	return nil
}

func (*NPM) Build(ctx context.Context, bctx BuildContext) error {
	bctx.Logger.Process("Executing build process")

	// assemble the information our executable needs
	exec := pexec.NewExecutable("npm")
	args := []string{"run", "build"}
	bctx.Logger.Subprocess("Running 'npm %s'", strings.Join(args, " "))
	// shell out
	duration, err := bctx.Clock.Measure(func() error {
		return exec.Execute(pexec.Execution{
			Args: args,
			Dir:  bctx.WorkingDir,
			Env: append(
				os.Environ(),
				"NPM_CONFIG_LOGLEVEL=error",
			),
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})
	})
	if err != nil {
		return fmt.Errorf("npm build failed: %w", err)
	}
	bctx.Logger.Action("Completed in %s", duration.Round(time.Millisecond))
	bctx.Logger.Break()
	return nil
}
