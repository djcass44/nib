package build

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func Push(ctx context.Context, img v1.Image, dst string) error {
	// push the image
	if err := crane.Push(img, dst, crane.WithContext(ctx), crane.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
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
	fmt.Println(ref.String() + "@" + d.String())
	return nil
}
