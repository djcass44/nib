package executor

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/paketo-buildpacks/packit/pexec"
	"time"
)

type Options struct {
	Command string
	Args    []string
}

// Exec runs an external process
func Exec(ctx BuildContext, opts Options) error {
	log := logr.FromContextOrDiscard(ctx.Ctx.Context)

	// assemble the information our executable needs
	executor := pexec.NewExecutable(opts.Command)
	log.Info("running executable process", "command", opts.Command, "args", opts.Args)
	// shell out
	start := time.Now()
	err := executor.Execute(pexec.Execution{
		Args: opts.Args,
		Dir:  ctx.Ctx.WorkingDirectory,
		Env:  ctx.Ctx.ConfigFile.Config.Env,
	})
	if err != nil {
		return fmt.Errorf("%s failed: %w", opts.Command, err)
	}
	log.Info("completed execution", "duration", time.Since(start))
	return nil
}
