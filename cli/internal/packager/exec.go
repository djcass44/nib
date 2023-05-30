package packager

import (
	"fmt"
	"github.com/paketo-buildpacks/packit/pexec"
	"os"
	"strings"
	"time"
)

type options struct {
	command  string
	args     []string
	extraEnv []string
}

func exec(ctx BuildContext, opts options) error {
	// assemble the information our executable needs
	exec := pexec.NewExecutable(opts.command)
	ctx.Logger.Subprocess("Running '%s %s'", opts.command, strings.Join(opts.args, " "))
	// shell out
	duration, err := ctx.Clock.Measure(func() error {
		return exec.Execute(pexec.Execution{
			Args: opts.args,
			Dir:  ctx.WorkingDir,
			Env: append(
				os.Environ(),
				opts.extraEnv...,
			),
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})
	})
	if err != nil {
		return fmt.Errorf("%s failed: %w", opts.command, err)
	}
	ctx.Logger.Action("Completed in %s", duration.Round(time.Millisecond))
	ctx.Logger.Break()
	return nil
}
