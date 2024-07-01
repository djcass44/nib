package pathfinder

import (
	cbev1 "github.com/Snakdy/container-build-engine/pkg/api/v1"
	"github.com/Snakdy/container-build-engine/pkg/pipelines"
	"github.com/Snakdy/container-build-engine/pkg/pipelines/utils"
	"github.com/go-logr/logr"
)

const StatementPathfinder = "pathfinder"

type Pathfinder struct {
	options cbev1.Options
}

func (p *Pathfinder) Run(ctx *pipelines.BuildContext, _ ...cbev1.Options) (cbev1.Options, error) {
	log := logr.FromContextOrDiscard(ctx.Context)
	log.Info("running pathfinder statement")

	buildDirs, err := cbev1.GetRequired[[]string](p.options, "build-dirs")
	if err != nil {
		return cbev1.Options{}, err
	}

	log.V(4).Info("checking for build directories", "dirs", buildDirs)
	appDir, err := FindBuildDir(ctx.WorkingDirectory, buildDirs)
	if err != nil {
		return cbev1.Options{}, err
	}

	log.V(4).Info("resolved build directory", "dir", appDir)
	return cbev1.Options{
		"src": appDir,
	}, nil
}

func (*Pathfinder) Name() string {
	return StatementPathfinder
}

func (p *Pathfinder) SetOptions(options cbev1.Options) {
	if p.options == nil {
		p.options = map[string]any{}
	}
	utils.CopyMap(options, p.options)
}
