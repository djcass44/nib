package packager

import (
	"context"
	"github.com/paketo-buildpacks/packit/chronos"
	"github.com/paketo-buildpacks/packit/scribe"
)

const EnvExtraArgs = "BUILD_EXTRA_ARGS"

type PackageManager interface {
	Detect(ctx context.Context, bctx BuildContext) bool
	Install(ctx context.Context, bctx BuildContext) error
	Build(ctx context.Context, bctx BuildContext) error
}

type BuildContext struct {
	WorkingDir string
	CacheDir   string

	Clock  chronos.Clock
	Logger scribe.Logger
}
