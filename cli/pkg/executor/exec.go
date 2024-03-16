package executor

import (
	"fmt"
	"github.com/paketo-buildpacks/packit/pexec"
	"os"
	"strings"
	"time"
)

type Options struct {
	Command  string
	Args     []string
	ExtraEnv []string
}

// Exec runs an external process
func Exec(ctx BuildContext, opts Options) error {
	// assemble the information our executable needs
	executor := pexec.NewExecutable(opts.Command)
	ctx.Logger.Subprocess("Running '%s %s'", opts.Command, strings.Join(opts.Args, " "))
	// shell out
	duration, err := ctx.Clock.Measure(func() error {
		return executor.Execute(pexec.Execution{
			Args: opts.Args,
			Dir:  ctx.WorkingDir,
			Env: append(
				os.Environ(),
				opts.ExtraEnv...,
			),
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})
	})
	if err != nil {
		return fmt.Errorf("%s failed: %w", opts.Command, err)
	}
	ctx.Logger.Action("Completed in %s", duration.Round(time.Millisecond))
	ctx.Logger.Break()
	return nil
}
