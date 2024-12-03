package executor

import (
	"github.com/Snakdy/container-build-engine/pkg/pipelines"
)

const EnvExtraArgs = "BUILD_TOOL_EXTRA_ARGS"

type PackageManager interface {
	Detect(ctx BuildContext) bool
	Install(ctx BuildContext) error
	Build(ctx BuildContext) error
}

// Deprecated
type BuildContext struct {
	CacheDir string
	Ctx      pipelines.BuildContext
}
