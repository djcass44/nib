package dotenv

import (
	"errors"
	cbev1 "github.com/Snakdy/container-build-engine/pkg/api/v1"
	"github.com/Snakdy/container-build-engine/pkg/envs"
	"github.com/Snakdy/container-build-engine/pkg/pipelines"
	"github.com/Snakdy/container-build-engine/pkg/pipelines/utils"
	"github.com/go-logr/logr"
	"io"
	"os"
	"path/filepath"
)

const StatementDotenv = "dotenv"

type Dotenv struct {
	options cbev1.Options
}

func (p *Dotenv) Run(ctx *pipelines.BuildContext) error {
	log := logr.FromContextOrDiscard(ctx.Context)
	log.Info("running dotenv statement", "options", p.options)

	skip, err := cbev1.GetOptional[bool](p.options, "skip")
	if err != nil {
		return err
	}
	path, err := cbev1.GetRequired[string](p.options, "path")
	if err != nil {
		return err
	}

	if skip {
		log.V(5).Info("skipping dotenv generation")
		return nil
	}

	path = filepath.Clean(envs.ExpandEnvFunc(path, pipelines.ExpandList(ctx.ConfigFile.Config.Env)))

	// copy the .env file if it exists
	srcPath := filepath.Join(ctx.WorkingDirectory, ".env")
	dstPath := filepath.Join(path, ".env")

	log.V(4).Info("checking for dotenv file", "path", srcPath)

	if _, err := os.Stat(srcPath); !errors.Is(err, os.ErrNotExist) {
		log.Info("detected .env file")
		err = func() error {
			src, err := os.Open(srcPath)
			if err != nil {
				return err
			}
			defer src.Close()
			log.V(4).Info("writing .env file", "path", dstPath)
			dst, err := ctx.FS.Create(dstPath)
			if err != nil {
				return err
			}
			defer dst.Close()
			_, err = io.Copy(dst, src)
			return err
		}()
		if err != nil {
			log.Error(err, "failed to copy .env file")
			return err
		}
	}

	return nil
}

func (p *Dotenv) Name() string {
	return StatementDotenv
}

func (p *Dotenv) SetOptions(options cbev1.Options) {
	if p.options == nil {
		p.options = map[string]any{}
	}
	utils.CopyMap(options, p.options)
}
