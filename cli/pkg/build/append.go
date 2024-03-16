package build

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/djcass44/all-your-base/pkg/containerutil"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

const (
	nibDataPath = "/var/run/nib"
	nibAuthor   = "github.com/djcass44/nib"
)

func Append(ctx context.Context, baseRef string, platform *v1.Platform, appPaths ...LayerPath) (v1.Image, error) {
	// pull the base image
	log.Printf("pulling base image: %s", baseRef)

	base, err := containerutil.Get(ctx, baseRef)
	if err != nil {
		return nil, fmt.Errorf("pulling %s: %w", baseRef, err)
	}

	// create our new layers
	var layers []mutate.Addendum
	for i, path := range appPaths {
		log.Printf("containerising directory %d: %s", i, path)
		layer, err := NewLayer(path.Path, path.Chroot, platform)
		if err != nil {
			return nil, err
		}

		// append our layer
		layers = append(layers, mutate.Addendum{
			MediaType: types.OCILayer,
			Layer:     layer,
			History: v1.History{
				Author:    "nib",
				CreatedBy: "nib build",
				Created:   v1.Time{},
				Comment:   "nibdata contents, at $NIB_DATA_PATH",
			},
		})
	}
	withData, err := mutate.Append(base, layers...)
	if err != nil {
		return nil, fmt.Errorf("appending layers: %w", err)
	}
	// grab a copy of the base image's config file, and set
	// our entrypoint and env vars
	cfg, err := withData.ConfigFile()
	if err != nil {
		return nil, err
	}
	cfg = cfg.DeepCopy()
	if platform.OS == "windows" {
		cfg.Config.Env = append(cfg.Config.Env, "NIB_DATA_PATH=C:"+strings.ReplaceAll(nibDataPath, "/", `\`))
	} else {
		cfg.Config.Env = append(cfg.Config.Env, "NIB_DATA_PATH="+nibDataPath)
	}
	cfg.Author = nibAuthor
	cfg.Config.WorkingDir = nibDataPath
	if cfg.Config.Labels == nil {
		cfg.Config.Labels = map[string]string{}
	}

	// package everything up
	img, err := mutate.ConfigFile(withData, cfg)
	if err != nil {
		return nil, err
	}
	return img, nil
}
