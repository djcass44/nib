package build

import (
	"context"
	"fmt"
	"github.com/djcass44/ci-tools/pkg/ociutil"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"log"
	"strings"
)

const nibDataPath = "/var/run/nib"

func Append(ctx context.Context, appPath, baseRef string, platform *v1.Platform) (v1.Image, error) {
	// pull the base image
	log.Printf("pulling base image: %s", baseRef)
	base, err := crane.Pull(baseRef, crane.WithContext(ctx), crane.WithAuthFromKeychain(ociutil.KeyChain(ociutil.Auth{})))
	if err != nil {
		return nil, fmt.Errorf("pulling %s: %w", baseRef, err)
	}

	// create our new layer
	log.Printf("containerising directory: %s", appPath)
	layer, err := NewLayer(appPath, platform)
	if err != nil {
		return nil, err
	}

	// append our layer
	layers := []mutate.Addendum{
		{
			Layer: layer,
			History: v1.History{
				Author:    "nib",
				CreatedBy: "nib build",
				Created:   v1.Time{},
				Comment:   "nibdata contents, at $NIB_DATA_PATH",
			},
		},
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
	cfg.Author = "github.com/djcass44/nib"
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
