package packager

import (
	cbev1 "github.com/Snakdy/container-build-engine/pkg/api/v1"
	"github.com/Snakdy/container-build-engine/pkg/pipelines"
	"github.com/Snakdy/container-build-engine/pkg/pipelines/utils"
	"github.com/djcass44/nib/cli/pkg/executor"
	"github.com/go-logr/logr"
)

const StatementNodePackage = "node-package"

var buildEngines = []executor.PackageManager{
	NewNPM(commandNPM),
	&Yarn{},
}

type NodePackager struct {
	options cbev1.Options
}

func (p *NodePackager) Run(ctx *pipelines.BuildContext, _ ...cbev1.Options) (cbev1.Options, error) {
	log := logr.FromContextOrDiscard(ctx.Context)
	log.V(7).Info("running statement node package", "options", p.options)

	cacheDir, err := cbev1.GetRequired[string](p.options, "cache-dir")
	if err != nil {
		return cbev1.Options{}, err
	}

	buildContext := executor.BuildContext{
		Ctx:      *ctx,
		CacheDir: cacheDir,
	}

	pkg := buildEngines[0]
	for _, engine := range buildEngines {
		ok := engine.Detect(buildContext)
		if ok {
			pkg = engine
			break
		}
	}

	// 1. install
	err = pkg.Install(buildContext)
	if err != nil {
		return cbev1.Options{}, err
	}

	// 2. build
	err = pkg.Build(buildContext)
	if err != nil {
		return cbev1.Options{}, err
	}
	return cbev1.Options{}, nil
}

func (p *NodePackager) Name() string {
	return StatementNodePackage
}

func (p *NodePackager) SetOptions(options cbev1.Options) {
	if p.options == nil {
		p.options = map[string]any{}
	}
	utils.CopyMap(options, p.options)
}
