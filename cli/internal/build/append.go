package build

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"io"
	"log"
	"os"
)

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
	// save the layer to a temporary file
	f, err := os.CreateTemp("", "layer-*.tar")
	if err != nil {
		return nil, err
	}
	r, err := layer.Compressed()
	if err != nil {
		return nil, fmt.Errorf("reading layer: %w", err)
	}
	_, err = io.Copy(f, r)
	if err != nil {
		return nil, fmt.Errorf("writing layer to %s: %w", f.Name(), err)
	}

	// append our layer
	img, err := crane.Append(base, f.Name())
	if err != nil {
		return nil, fmt.Errorf("appending %s: %w", f.Name(), err)
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
