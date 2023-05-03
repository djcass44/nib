package build

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"log"
	"strings"
)

const nipDataPath = "/var/run/nip"

func Append(ctx context.Context, appPath, baseRef string, platform *v1.Platform) (v1.Image, error) {
	// pull the base image
	log.Printf("pulling base image: %s", baseRef)
	base, err := crane.Pull(baseRef, crane.WithContext(ctx))
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
				Author:    "nip",
				CreatedBy: "nip build",
				Created:   v1.Time{},
				Comment:   "nipdata contents, at $NIP_DATA_PATH",
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
		cfg.Config.Env = append(cfg.Config.Env, "NIP_DATA_PATH=C:"+strings.ReplaceAll(nipDataPath, "/", `\`))
	} else {
		cfg.Config.Env = append(cfg.Config.Env, "NIP_DATA_PATH="+nipDataPath)
	}
	cfg.Author = "github.com/djcass44/nip"
	cfg.Config.WorkingDir = nipDataPath
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

func Push(ctx context.Context, img v1.Image, dst string) error {
	// push the image
	if err := crane.Push(img, dst, crane.WithContext(ctx)); err != nil {
		return fmt.Errorf("pushing image %s: %w", dst, err)
	}
	// parse what we just pushed, so we can show
	// the user
	ref, err := name.ParseReference(dst)
	if err != nil {
		return fmt.Errorf("parsing reference %s: %w", dst, err)
	}
	d, err := img.Digest()
	if err != nil {
		return fmt.Errorf("digest: %w", err)
	}
	fmt.Println(ref.Context().Digest(d.String()))
	return nil
}
